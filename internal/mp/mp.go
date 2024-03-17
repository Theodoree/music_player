package mp

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	
	"fyne.io/fyne/v2/data/binding"
	"github.com/Theodoree/music_player/internal/db"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/music"
	"github.com/Theodoree/music_player/internal/music/local"
	"github.com/Theodoree/music_player/internal/music/netease"
	"github.com/Theodoree/music_player/internal/tool"
	"gorm.io/gorm"
	"k8s.io/klog"
)

var _ MusicPlayer = (*musicPlayer)(nil)

type musicPlayerData struct {
	tableList        bindingTable[model.MusicTable]
	picList          bindingTable[model.Picture]
	volume           BindingModel[float64]
	singerName       BindingModel[string]
	musicName        BindingModel[string]
	PlayStatus       BindingModel[bool]
	streamMusicTable bindingTable[*BindingModel[string]]
	processBar       processBar
	mode             BindingModel[PlayMode]
}

type bindingTable[T binding.DataItem] struct {
	items BindingDataList[T]
	index BindingModel[int]
}

type musicPlayer struct {
	ctx           context.Context
	cancel        context.CancelFunc
	store         db.MusicStore
	alert         func(str string)
	cb            music.Callback
	localSource   music.Source
	neteaseSource music.Source
	
	settings
	musicPlayerData musicPlayerData
	selectList      *list
	list            list
	neteaseList     list
	curMusic        music.Music
}

type settings struct {
	SavePath string
}

func NewMusicPlayer(ctx context.Context, alert func(str string)) (MusicPlayer, error) {
	var s musicPlayer
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.store = db.New(db.SqliteFactory("/Users/ted/workspace/go/music_player"), db.MemoryCacheFactory())
	if err := s.InitSettings("/Users/ted/workspace/go/music_player"); err != nil {
		return nil, err
	}
	s.alert = alert
	s.cb = music.Callback{
		CurTime: func(duration time.Duration) {
			s.musicPlayerData.processBar.UpdateTS(duration)
		},
		DoneFn: func(status model.Status) {
			if status == model.StatusPlayDone {
				s.Next()
			}
		},
	}
	s.localSource = local.Source(s.ctx, s.store)
	s.neteaseSource = netease.Source(s.ctx, s.settings.SavePath)
	s.selectList = &s.list
	if err := s.init(); err != nil {
		return nil, err
	}
	return &s, nil
}

// 绑定数据
func (m *musicPlayer) init() error {
	if err := m.refreshTable(); err != nil {
		return err
	}
	_ = m.musicPlayerData.volume.Set(10)
	
	// 本地列表
	m.musicPlayerData.tableList.index.AddListener(&DataListener{func() {
		idx := m.musicPlayerData.tableList.index.get()
		item, _ := m.musicPlayerData.tableList.items.GetItem(idx)
		id := item.(model.MusicTable).ID
		musics, _ := m.localSource.List(id)
		m.list.setItems(musics, id, true)
		m.selectList = &m.list
	}})
	_ = m.musicPlayerData.tableList.index.Set(0)
	
	// 流媒体表格
	m.musicPlayerData.streamMusicTable.items.AddItem(&BindingModel[string]{val: "网易云"})
	m.list.searchKey.AddListener(&DataListener{Fn: func() {
		m.list.Search(m.list.searchKey.get())
	}})
	
	// 流媒体列表搜索
	m.neteaseList.searchKey.AddListener(&DataListener{func() {
		index := m.musicPlayerData.streamMusicTable.index.get()
		searchKey := m.neteaseList.searchKey.get()
		switch index {
		case 0:
			musics, _ := m.neteaseSource.SearchMusic(0, searchKey)
			m.neteaseList.setItems(musics, 0, false)
			m.selectList = &m.neteaseList
		}
		
	}})
	
	// 播放条拖动
	m.musicPlayerData.processBar.cur.AddListener(&DataListener{func() {
		cur := m.musicPlayerData.processBar.cur.get()
		if cur == m.musicPlayerData.processBar.lastSet || m.curMusic == nil {
			return
		}
		m.curMusic.Seek(cur)
	}})
	
	// 音量条拖动
	m.musicPlayerData.volume.AddListener(&DataListener{func() {
		volume, _ := m.musicPlayerData.volume.Get()
		if m.curMusic != nil {
			m.curMusic.SetVolume(volume / 100)
		}
	}})
	
	return nil
}

// InitSettings 初始化基础配置
func (m *musicPlayer) InitSettings(basePath string) error {
	if basePath == "" {
		var err error
		basePath, err = filepath.Abs(filepath.Dir(filepath.Join(os.Args[0])))
		if err != nil {
			klog.Error(err)
			return err
		}
	}
	m.settings.SavePath = filepath.Join(basePath, "save")
	_ = os.MkdirAll(m.settings.SavePath, 0766)
	return nil
}

// Play Implementation MusicPlayerFrontend
func (m *musicPlayer) Play() {
	if !m.selectList.valid() {
		m.alert("No music")
		return
	}
	if m.curMusic == nil {
		m.curMusic = m.selectList.next(m.musicPlayerData.mode.get())
	}
	m._play(m.curMusic)
}
func (m *musicPlayer) Pause() {
	if !m.selectList.valid() || m.curMusic == nil {
		m.alert("No music")
		return
	}
	err := m.curMusic.Pause()
	if err != nil {
		klog.Error(err)
		m.alert(err.Error())
		return
	}
	_ = m.musicPlayerData.PlayStatus.Set(false)
}
func (m *musicPlayer) Stop() {
	if !m.selectList.valid() || m.curMusic == nil {
		m.alert("No music")
		return
	}
	err := m.curMusic.Stop()
	if err != nil {
		klog.Error(err)
		m.alert(err.Error())
		return
	}
	m.curMusic = nil
	m.resetMusicPlayerData()
	_ = m.musicPlayerData.PlayStatus.Set(false)
	return
	
}

func (m *musicPlayer) resetMusicPlayerData() {
	m.musicPlayerData.processBar.UpdateTS(0)
	m.musicPlayerData.processBar.UpdateEnd(0)
	_ = m.musicPlayerData.singerName.Set("")
	_ = m.musicPlayerData.musicName.Set("")
}
func (m *musicPlayer) Prev() {
	if !m.selectList.valid() {
		m.alert("No music")
		return
	}
	if m.curMusic != nil {
		err := m.curMusic.Stop()
		if err != nil {
			klog.Error(err)
		}
		m.curMusic = nil
		m.resetMusicPlayerData()
	}
	m._play(m.selectList.prev(m.musicPlayerData.mode.get()))
}
func (m *musicPlayer) Next() {
	if !m.selectList.valid() {
		m.alert("No music")
		return
	}
	if m.curMusic != nil {
		err := m.curMusic.Stop()
		if err != nil {
			klog.Error(err)
		}
		m.curMusic = nil
		m.resetMusicPlayerData()
	}
	m._play(m.selectList.next(m.musicPlayerData.mode.get()))
}
func (m *musicPlayer) _play(music music.Music) bool {
	volume, _ := m.musicPlayerData.volume.Get()
	
	if err := music.Play(&m.cb, volume/1e2); err != nil {
		klog.Error(err)
		m.alert(err.Error())
		return false
	}
	
	end, _ := music.EndTime()
	m.musicPlayerData.processBar.UpdateEnd(end)
	_ = m.musicPlayerData.musicName.Set(music.MusicName())
	_ = m.musicPlayerData.singerName.Set(music.SingerName())
	m.curMusic = music
	_ = m.musicPlayerData.PlayStatus.Set(true)
	m.musicPlayerData.processBar.UpdateLyrics(music.Lyrics())
	return true
}

// Lyrics Implementation MusicPlayerFrontendData
func (m *musicPlayer) Lyrics() binding.String {
	return &m.musicPlayerData.processBar.lyrics
}
func (m *musicPlayer) Lyric() binding.String {
	return &m.musicPlayerData.processBar.lyric
}
func (m *musicPlayer) MusicName() binding.String {
	return &m.musicPlayerData.musicName
}
func (m *musicPlayer) SingerName() binding.String {
	return &m.musicPlayerData.singerName
}
func (m *musicPlayer) Volume() binding.Float {
	return &m.musicPlayerData.volume
}
func (m *musicPlayer) ProcessBar() binding.Float {
	return &m.musicPlayerData.processBar.cur
}
func (m *musicPlayer) MusicCurTime() binding.String {
	return &m.musicPlayerData.processBar.curStr
}
func (m *musicPlayer) MusicEndTime() binding.String {
	return &m.musicPlayerData.processBar.end
}
func (m *musicPlayer) MusicTableList() (binding.DataList, binding.Int) {
	return &m.musicPlayerData.tableList.items, &m.musicPlayerData.tableList.index
}
func (m *musicPlayer) MusicList() (binding.DataList, binding.Int, binding.String) {
	return &m.list.items, &m.list.index, &m.list.searchKey
}
func (m *musicPlayer) PictureList() binding.DataList {
	return &m.musicPlayerData.picList.items
}
func (m *musicPlayer) PlayMode() binding.DataItem {
	return &m.musicPlayerData.mode
}
func (m *musicPlayer) PlayStatus() binding.Bool {
	return &m.musicPlayerData.PlayStatus
}
func (m *musicPlayer) StreamMusicList() (binding.DataList, binding.Int, binding.String) {
	return &m.neteaseList.items, &m.neteaseList.index, &m.neteaseList.searchKey
}
func (m *musicPlayer) StreamMusicTableList() (binding.DataList, binding.Int) {
	return &m.musicPlayerData.streamMusicTable.items, &m.musicPlayerData.streamMusicTable.index
}

// AddTable Implementation MusicPlayerBackend
func (m *musicPlayer) AddTable(table model.MusicTable) {
	if err := m.store.SaveMusicTable(table); err != nil {
		m.alert(err.Error())
		return
	}
	if err := m.refreshTable(); err != nil {
		m.alert(err.Error())
	}
	
}
func (m *musicPlayer) DelTable(tableID uint) {
	if err := m.store.DeleteMusicTable(model.MusicTable{
		Model: gorm.Model{
			ID: tableID,
		},
	}); err != nil {
		m.alert(err.Error())
		return
	}
	if err := m.store.DeleteMusicByMusicTableID(tableID); err != nil {
		m.alert(err.Error())
		return
	}
	if err := m.refreshTable(); err != nil {
		m.alert(err.Error())
	}
}
func (m *musicPlayer) ImportMusic(tableID uint, path string) {
	if _, err := m.store.GetMusicTableByID(tableID); err != nil {
		m.alert(err.Error())
		return
	}
	
	items := tool.SearchMusicFileByPath(path)
	for idx := range items {
		items[idx].MusicTableID = tableID
	}
	if err := m.store.SaveMusics(items); err != nil {
		klog.Error(err)
		m.alert(err.Error())
		return
	}
}
func (m *musicPlayer) AddWallpaper(path string) {
	_ = path
}
func (m *musicPlayer) AddMusic(tableID uint, music model.Music) {
	music.MusicTableID = tableID
	music.ID = 0
	if err := m.store.SaveMusic(music); err != nil {
		klog.Error(err)
		if strings.Index(err.Error(), "UNIQUE") >= 0 {
			return
		}
		m.alert(err.Error())
		return
	}
	if err := m.refreshMusic(tableID); err != nil {
		m.alert(err.Error())
	}
}
func (m *musicPlayer) UpdateMusic(tableID uint, music model.Music) {
	music.MusicTableID = tableID
	if err := m.store.UpdateMusic(music); err != nil {
		m.alert(err.Error())
	}
	if m.list.tableId != tableID {
		return
	}
	if err := m.refreshMusic(tableID); err != nil {
		m.alert(err.Error())
		return
	}
	if m.curMusic == nil {
		return
	}
	curMusic, _ := m.curMusic.GetMusic()
	if curMusic.ID != music.ID {
		return
	}
	m.curMusic.Update(music)
	_ = m.musicPlayerData.musicName.Set(m.curMusic.MusicName())
	_ = m.musicPlayerData.singerName.Set(m.curMusic.SingerName())
	m.musicPlayerData.processBar.UpdateLyrics(m.curMusic.Lyrics())
	
}
func (m *musicPlayer) GetPlayedMusic() music.Music {
	return m.curMusic
}

// private
func (m *musicPlayer) refreshTable() error {
	tables, err := m.store.GetMusicTable(0, 5000)
	if err != nil {
		return err
	}
	// 本地列表
	m.musicPlayerData.tableList.items.SetItems(tables)
	return nil
}

func (m *musicPlayer) refreshMusic(tableID uint) error {
	if m.list.tableId != tableID {
		return nil
	}
	for idx, v := range m.musicPlayerData.tableList.items.items {
		if v.ID != tableID {
			continue
		}
		_ = m.musicPlayerData.tableList.index.Set(idx)
	}
	return nil
}

type BindingModel[T any] struct {
	val           T
	DataListeners []binding.DataListener
}

func (b *BindingModel[T]) Get() (T, error) {
	return b.get(), nil
}
func (b *BindingModel[T]) get() T {
	return b.val
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

func (b *BindingDataList[T]) SetItems(items []T) {
	b.items = items
	b.Signal()
}
func (b *BindingDataList[T]) AddListener(listener binding.DataListener) {
	b.DataListeners = append(b.DataListeners, listener)
}
func (b *BindingDataList[T]) RemoveListener(ls binding.DataListener) {
	b.DataListeners = slices.DeleteFunc(b.DataListeners, func(listener binding.DataListener) bool {
		return listener == ls
	})
}
func (b *BindingDataList[T]) AddItem(item T) {
	b.items = append(b.items, item)
	b.Signal()
}
func (b *BindingDataList[T]) GetItem(index int) (binding.DataItem, error) {
	return b.items[index], nil
}
func (b *BindingDataList[T]) Length() int {
	return len(b.items)
}
func (b *BindingDataList[T]) Signal() {
	for _, listener := range b.DataListeners {
		listener.DataChanged()
	}
}

type processBar struct {
	curStr  BindingModel[string]  // 00:00
	cur     BindingModel[float64] // 百分比
	end     BindingModel[string]  // 00:00
	lyrics  BindingModel[string]
	lyric   BindingModel[string]
	lrcLine []LrcLine
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
func (p *processBar) UpdateLyrics(str string) {
	p.lrcLine, _ = decodeLrc(str)
	var lyrics string
	for _, v := range p.lrcLine {
		lyrics += v.Lyrics + "\n"
	}
	_ = p.lyrics.Set(lyrics)
}

type DataListener struct {
	Fn func()
}

func (d *DataListener) DataChanged() {
	d.Fn()
}
