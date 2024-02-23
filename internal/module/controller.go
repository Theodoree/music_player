package module

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"github.com/Theodoree/music_player/internal/db"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/tool"
	"k8s.io/klog"
	"math"
	"slices"
	"time"
)

type MusicPlayerOperation interface {
	Play() bool
	Pause() bool
	Prev()
	Next()
	
	SelectMusic(id uint)
	SelectTable(id uint)
	
	SaveToTable(path string)
	GetMusicTableList() binding.DataList
	GetMusicList() binding.DataList
	GetVolume() binding.Float        // [0.01,1]
	GetMusicCurTime() binding.String // 00:00
	GetProcessBar() binding.Float    // [0,100]
	GetMusicEndTime() binding.String // 00:00
	GetSingerName() binding.String
	GetMusicName() binding.String
}

var _ MusicPlayerOperation = (*controller)(nil)

type controller struct {
	ctx    context.Context
	cancel context.CancelFunc
	store  db.MusicStore
	
	selectTableID uint
	list          list    // 音乐列表
	player        *player // 播放器
	
	table      normalDataList[model.MusicTable] //
	mode       BindingModel[int]                // [0,1,2]  Cycle SingleCycle Random
	volume     BindingModel[float64]            // 音量[0.01,1]
	processBar processBar                       // 进度条相关组件
	musicName  BindingModel[string]             // 歌曲
	SingerName BindingModel[string]             // 歌手名
	guiAlert   func(str string)                 // 前端提示函数
}

func NewController(ctx context.Context, alert func(str string)) MusicPlayerOperation {
	var c controller
	c.ctx, c.cancel = context.WithCancel(ctx)
	c.store = db.New(db.SqliteFactory(""), db.MemoryCacheFactory())
	c.selectTableID = db.DefaultTableID
	c.list.setMode(PlaybackModeCycle)
	c.refreshList()
	c.refreshTable()
	c.player = nil
	
	// 初始化绑定对象
	_ = c.mode.Set(int(PlaybackModeCycle))
	c.mode.AddListener(&dataListener{fn: func() {
		c.list.setMode(PlaybackMode(c.mode.val))
	}})
	
	_ = c.volume.Set(0.1)
	c.volume.AddListener(&dataListener{fn: func() {
		if c.player == nil {
			return
		}
		c.player.SetVolume(c.volume.val)
	}})
	
	c.processBar.UpdateTS(0)
	c.processBar.UpdateEnd(0)
	c.processBar.cur.AddListener(&dataListener{fn: func() {
		if c.player == nil {
			return
		}
		
		val, _ := c.processBar.cur.Get()
		if val == c.processBar.lastSet || val == 0 { // 自己设置的
			return
		}
		c.player.Seek(val)
	}})
	
	_ = c.musicName.Set("perfect")
	_ = c.SingerName.Set("Ed Sheeran")
	c.guiAlert = alert
	
	return &c
}

func (c *controller) refreshTable() {
	tables, _ := c.store.GetMusicTable(0, 1000)
	c.table.SetItems(tables)
}

func (c *controller) refreshList() {
	// 切换表格,停止音乐，切换新表格
	//if c.player != nil {
	//	c.player.Close()
	//	c.player = nil
	//}
	items, _ := c.store.GetMusicByMusicTableID(c.selectTableID)
	c.list.setItems(items)
}

func (c *controller) Play() bool {
	if c.player == nil {
		return false
	}
	c.player.Play()
	return true
}
func (c *controller) Pause() bool {
	if c.player == nil {
		return false
	}
	c.player.Pause()
	return true
}
func (c *controller) Prev() {
	if !c.list.valid() {
		c.guiAlert("No music")
		return
	}
	if c.player != nil {
		c.player.Close()
		c.player = nil
	}
	c._play(c.list.prev())
}
func (c *controller) Next() {
	if !c.list.valid() {
		c.guiAlert("No music")
		return
	}
	if c.player != nil {
		c.player.Close()
		c.player = nil
	}
	c._play(c.list.next())
}
func (c *controller) _play(music model.Music) {
	volume, _ := c.volume.Get()
	player, err := newPlayer(c.ctx, music, volume, func(cur time.Duration) {
		c.processBar.UpdateTS(cur)
	}, c.Next)
	if err != nil {
		klog.Error(err)
		c.guiAlert("播放失败,尝试下一首")
		time.Sleep(time.Second)
		c.Next()
		return
	}
	_ = c.musicName.Set(music.Name)
	_ = c.SingerName.Set(music.Singer)
	c.processBar.UpdateEnd(player.EndTime())
	c.player = player
}
func (c *controller) SelectMusic(musicID uint) {
	idx := c.list.SearchByID(musicID)
	c.list.setIdx(idx - 1)
	c.Next()
}
func (c *controller) SelectTable(selectTableID uint) {
	if c.selectTableID == selectTableID {
		return
	}
	c.selectTableID = selectTableID
	c.refreshList()
}

func (c *controller) SaveToTable(path string) {
	if c.selectTableID == 0 {
		return
	}
	musics := tool.SearchMusicFileByPath(path)
	for idx := range musics {
		musics[idx].MusicTableID = c.selectTableID
	}
	err := c.store.SaveMusics(musics)
	if err != nil {
		c.guiAlert(err.Error())
		return
	}
	
	c.refreshList()
}

func (c *controller) GetMusicTableList() binding.DataList {
	return &c.table
}
func (c *controller) GetMusicList() binding.DataList {
	return &c.list
}
func (c *controller) GetVolume() binding.Float {
	return &c.volume
}
func (c *controller) GetMusicCurTime() binding.String {
	return &c.processBar.curStr
}
func (c *controller) GetProcessBar() binding.Float {
	return &c.processBar.cur
}
func (c *controller) GetMusicEndTime() binding.String {
	return &c.processBar.end
}
func (c *controller) GetSingerName() binding.String {
	return &c.SingerName
}
func (c *controller) GetMusicName() binding.String {
	return &c.musicName
}

type BindingModel[T any] struct {
	val           T
	DataListeners []binding.DataListener
}

func (b *BindingModel[T]) Get() (T, error) {
	return b.val, nil
}
func (b *BindingModel[T]) Set(t T) error {
	b.val = t
	b.Signal()
	return nil
}
func (b *BindingModel[T]) AddListener(listener binding.DataListener) {
	b.DataListeners = append(b.DataListeners, listener)
}
func (b *BindingModel[T]) RemoveListener(ls binding.DataListener) {
	b.DataListeners = slices.DeleteFunc(b.DataListeners, func(listener binding.DataListener) bool {
		return listener == ls
	})
}
func (b *BindingModel[T]) Signal() {
	for _, listener := range b.DataListeners {
		listener.DataChanged()
	}
}

type BindingDataList[T binding.DataItem] struct {
	items         []T
	DataListeners []binding.DataListener
}

func (b *BindingDataList[T]) AddListener(listener binding.DataListener) {
	b.DataListeners = append(b.DataListeners, listener)
}
func (b *BindingDataList[T]) RemoveListener(ls binding.DataListener) {
	b.DataListeners = slices.DeleteFunc(b.DataListeners, func(listener binding.DataListener) bool {
		return listener == ls
	})
}
func (b *BindingDataList[T]) GetItem(index int) (binding.DataItem, error) {
	return b.items[index], nil
}
func (b *BindingDataList[T]) Length() int {
	return len(b.items)
}

type dataListener struct {
	fn func()
}

func (d *dataListener) DataChanged() {
	d.fn()
}

type processBar struct {
	curStr BindingModel[string]  // 00:00
	cur    BindingModel[float64] // 百分比
	end    BindingModel[string]  // 00:00
	
	lastSet float64
	endTs   time.Duration
}

func (p *processBar) UpdateEnd(ts time.Duration) {
	p.endTs = ts
	_ = p.end.Set(fmt.Sprintf("%02d:%02d", int(ts.Minutes()), int(ts.Seconds())%60))
}
func (p *processBar) UpdateTS(ts time.Duration) {
	_ = p.curStr.Set(fmt.Sprintf("%02d:%02d", int(ts.Minutes()), int(ts.Seconds())%60))
	if ts == 0 {
		_ = p.cur.Set(0)
		return
	}
	p.lastSet = math.Round(float64(ts) / float64(p.endTs) * 100)
	_ = p.cur.Set(p.lastSet)
}

type normalDataList[T binding.DataItem] struct {
	BindingDataList[T]
	EQ func(a, b T) bool
}

func (b *normalDataList[T]) Signal() {
	for _, listener := range b.DataListeners {
		listener.DataChanged()
	}
}
func (b *normalDataList[T]) AddItem(items T) {
	b.items = append(b.items, items)
	b.Signal()
}
func (b *normalDataList[T]) DeleteItem(_item T) {
	b.items = slices.DeleteFunc(b.items, func(item T) bool {
		return b.EQ(item, _item)
	})
	b.Signal()
}
func (b *normalDataList[T]) SetItems(items []T) {
	b.items = items
	b.Signal()
}
