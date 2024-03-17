package netease

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"fyne.io/fyne/v2/data/binding"
	"github.com/Theodoree/music_player/internal/decode"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/music"
)

var _ music.Source = (*neteaseSource)(nil)
var server = "39.101.203.25:3000"
var NoDecodeError = errors.New("no decode")

type neteaseSource struct {
	ctx      context.Context
	cancel   context.CancelFunc
	savePath string
}

func Source(ctx context.Context, savePath string) music.Source {
	var n neteaseSource
	n.ctx, n.cancel = context.WithCancel(ctx)
	n.savePath = savePath
	return &n
}

func (api *neteaseSource) SearchMusic(_ uint, keyword string) ([]music.Music, error) {
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/cloudsearch?limit=%d&keywords=%s", server, 30, url.QueryEscape(keyword)), nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	
	var result neteaseSearchInfo
	err = json.Unmarshal(buf, &result)
	if err != nil {
		return nil, err
	}
	
	var items []music.Music
	for _, song := range result.Result.Songs {
		// 仅试听
		if len(song.Privilege.ChargeInfoList) == 0 || song.Privilege.ChargeInfoList[0].ChargeType == 1 {
			continue
		}
		
		items = append(items, newNeteaseMusic(song, api.savePath))
	}
	
	return items, nil
	
}
func (api *neteaseSource) List(_ uint) ([]music.Music, error) {
	return api.SearchMusic(0, "")
	
}
func (api *neteaseSource) Close() error {
	api.cancel()
	return nil
}

type neteaseSearchInfo struct {
	Result struct {
		Songs []song `json:"songs"`
	} `json:"result"`
	Code int `json:"code"`
}

type song struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Pst  int    `json:"pst"`
	T    int    `json:"t"`
	Ar   []struct {
		Id    int           `json:"id"`
		Name  string        `json:"name"`
		Tns   []interface{} `json:"tns"`
		Alias []string      `json:"alias"`
		Alia  []string      `json:"alia,omitempty"`
	} `json:"ar"`
	Alia []interface{} `json:"alia"`
	Pop  int           `json:"pop"`
	St   int           `json:"st"`
	Rt   *string       `json:"rt"`
	Fee  int           `json:"fee"`
	V    int           `json:"v"`
	Crbt interface{}   `json:"crbt"`
	Cf   string        `json:"cf"`
	Al   struct {
		Id     int           `json:"id"`
		Name   string        `json:"name"`
		PicUrl string        `json:"picUrl"`
		Tns    []interface{} `json:"tns"`
		PicStr string        `json:"pic_str,omitempty"`
		Pic    int64         `json:"pic"`
	} `json:"al"`
	Dt int `json:"dt"`
	H  struct {
		Br   int `json:"br"`
		Fid  int `json:"fid"`
		Size int `json:"size"`
		Vd   int `json:"vd"`
		Sr   int `json:"sr"`
	} `json:"h"`
	M struct {
		Br   int `json:"br"`
		Fid  int `json:"fid"`
		Size int `json:"size"`
		Vd   int `json:"vd"`
		Sr   int `json:"sr"`
	} `json:"m"`
	L struct {
		Br   int `json:"br"`
		Fid  int `json:"fid"`
		Size int `json:"size"`
		Vd   int `json:"vd"`
		Sr   int `json:"sr"`
	} `json:"l"`
	Sq *struct {
		Br   int `json:"br"`
		Fid  int `json:"fid"`
		Size int `json:"size"`
		Vd   int `json:"vd"`
		Sr   int `json:"sr"`
	} `json:"sq"`
	Hr                   interface{}   `json:"hr"`
	A                    interface{}   `json:"a"`
	Cd                   string        `json:"cd"`
	No                   int           `json:"no"`
	RtUrl                interface{}   `json:"rtUrl"`
	Ftype                int           `json:"ftype"`
	RtUrls               []interface{} `json:"rtUrls"`
	DjId                 int           `json:"djId"`
	Copyright            int           `json:"copyright"`
	SId                  int           `json:"s_id"`
	Mark                 int64         `json:"mark"`
	OriginCoverType      int           `json:"originCoverType"`
	OriginSongSimpleData *struct {
		SongId  int    `json:"songId"`
		Name    string `json:"name"`
		Artists []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"artists"`
		AlbumMeta struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"albumMeta"`
	} `json:"originSongSimpleData"`
	TagPicList        interface{} `json:"tagPicList"`
	ResourceState     bool        `json:"resourceState"`
	Version           int         `json:"version"`
	SongJumpInfo      interface{} `json:"songJumpInfo"`
	EntertainmentTags interface{} `json:"entertainmentTags"`
	Single            int         `json:"single"`
	NoCopyrightRcmd   interface{} `json:"noCopyrightRcmd"`
	Mst               int         `json:"mst"`
	Cp                int         `json:"cp"`
	Rtype             int         `json:"rtype"`
	Rurl              interface{} `json:"rurl"`
	Mv                int         `json:"mv"`
	PublishTime       int64       `json:"publishTime"`
	Privilege         struct {
		Id                 int         `json:"id"`
		Fee                int         `json:"fee"`
		Payed              int         `json:"payed"`
		St                 int         `json:"st"`
		Pl                 int         `json:"pl"`
		Dl                 int         `json:"dl"`
		Sp                 int         `json:"sp"`
		Cp                 int         `json:"cp"`
		Subp               int         `json:"subp"`
		Cs                 bool        `json:"cs"`
		Maxbr              int         `json:"maxbr"`
		Fl                 int         `json:"fl"`
		Toast              bool        `json:"toast"`
		Flag               int         `json:"flag"`
		PreSell            bool        `json:"preSell"`
		PlayMaxbr          int         `json:"playMaxbr"`
		DownloadMaxbr      int         `json:"downloadMaxbr"`
		MaxBrLevel         string      `json:"maxBrLevel"`
		PlayMaxBrLevel     string      `json:"playMaxBrLevel"`
		DownloadMaxBrLevel string      `json:"downloadMaxBrLevel"`
		PlLevel            string      `json:"plLevel"`
		DlLevel            string      `json:"dlLevel"`
		FlLevel            string      `json:"flLevel"`
		Rscl               interface{} `json:"rscl"`
		FreeTrialPrivilege struct {
			ResConsumable      bool        `json:"resConsumable"`
			UserConsumable     bool        `json:"userConsumable"`
			ListenType         interface{} `json:"listenType"`
			CannotListenReason interface{} `json:"cannotListenReason"`
		} `json:"freeTrialPrivilege"`
		RightSource    int `json:"rightSource"`
		ChargeInfoList []struct {
			Rate          int         `json:"rate"`
			ChargeUrl     interface{} `json:"chargeUrl"`
			ChargeMessage interface{} `json:"chargeMessage"`
			ChargeType    int         `json:"chargeType"`
		} `json:"chargeInfoList"`
	} `json:"privilege"`
}

type neteaseAudioInfo struct {
	Data []struct {
		ID                 int         `json:"id"`
		URL                string      `json:"url"`
		Br                 int         `json:"br"`
		Size               int         `json:"size"`
		Md5                string      `json:"md5"`
		Code               int         `json:"code"`
		Expi               int         `json:"expi"`
		Type               string      `json:"type"`
		Gain               float64     `json:"gain"`
		Fee                int         `json:"fee"`
		Uf                 interface{} `json:"uf"`
		Payed              int         `json:"payed"`
		Flag               int         `json:"flag"`
		CanExtend          bool        `json:"canExtend"`
		FreeTrialInfo      interface{} `json:"freeTrialInfo"`
		Level              string      `json:"level"`
		EncodeType         string      `json:"encodeType"`
		FreeTrialPrivilege struct {
			ResConsumable  bool        `json:"resConsumable"`
			UserConsumable bool        `json:"userConsumable"`
			ListenType     interface{} `json:"listenType"`
		} `json:"freeTrialPrivilege"`
		FreeTimeTrialPrivilege struct {
			ResConsumable  bool `json:"resConsumable"`
			UserConsumable bool `json:"userConsumable"`
			Type           int  `json:"type"`
			RemainTime     int  `json:"remainTime"`
		} `json:"freeTimeTrialPrivilege"`
		URLSource   int         `json:"urlSource"`
		RightSource int         `json:"rightSource"`
		PodcastCtrp interface{} `json:"podcastCtrp"`
		EffectTypes interface{} `json:"effectTypes"`
		Time        int         `json:"time"`
	} `json:"data"`
	Code int `json:"code"`
}

var _ music.Music = (*neteaseMusic)(nil)

type neteaseMusic struct {
	id       int
	title    string
	singer   string
	album    string
	albumPic string
	savePath string
	
	decode   decode.Decoder
	audio    string
	audioURL string
	endTime  time.Duration
}

func strJoin[T any](items []T, fn func(T) string, sep string) string {
	var strs []string
	for _, item := range items {
		strs = append(strs, fn(item))
	}
	return strings.Join(strs, sep)
}

func newNeteaseMusic(s song, savePath string) music.Music {
	var m neteaseMusic
	m.id = s.Id
	m.title = s.Name
	m.singer = strJoin(s.Ar, func(struct {
		Id    int           `json:"id"`
		Name  string        `json:"name"`
		Tns   []interface{} `json:"tns"`
		Alias []string      `json:"alias"`
		Alia  []string      `json:"alia,omitempty"`
	}) string(func(s struct {
		Id    int
		Name  string
		Tns   []interface{}
		Alias []string
		Alia  []string
	}) string {
		return s.Name
	}), "/")
	m.album = s.Al.Name
	m.albumPic = s.Al.PicUrl
	m.savePath = savePath
	
	return &m
}

func (n *neteaseMusic) getFileName() string {
	return fmt.Sprintf("%s-%s-%d.mp3", strings.Replace(n.singer, "/", " ", -1), strings.Replace(n.title, "/", " ", -1), n.id)
}
func (n *neteaseMusic) lyricFileName() string {
	return fmt.Sprintf("%s-%s-%d.lrc", strings.Replace(n.singer, "/", " ", -1), strings.Replace(n.title, "/", " ", -1), n.id)
}
func (n *neteaseMusic) getReader() (io.ReadSeekCloser, error) {
	file, err := os.Open(filepath.Join(n.savePath, n.getFileName()))
	if err == nil {
		return file, nil
	}
	if err := n.getAudio(); err != nil {
		return nil, err
	}
	m, err := n.GetMusic()
	if err != nil {
		return nil, err
	}
	return os.Open(m.Path)
}
func (n *neteaseMusic) getAudio() error {
	if n.audio != "" {
		return nil
	}
	_, err := os.Stat(filepath.Join(n.savePath, n.getFileName()))
	if err == nil {
		n.audio = filepath.Join(n.savePath, n.getFileName())
		return nil
	}
	
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/song/url?id=%d", server, n.id), nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	buf, _ := io.ReadAll(res.Body)
	
	var info neteaseAudioInfo
	err = json.Unmarshal(buf, &info)
	if err != nil {
		return err
	}
	if len(info.Data) == 0 {
		return errors.New("no audio")
	}
	item := info.Data[0]
	
	duration := time.Duration(item.Time) * time.Millisecond / time.Second * time.Second
	if duration == 0 || duration == time.Second*30 || item.URL == "" {
		return errors.New("no audio")
	}
	n.audioURL = item.URL
	n.endTime = duration
	
	return nil
}

// TableID Implementation music.Data
func (n *neteaseMusic) TableID() uint {
	return 0
}
func (n *neteaseMusic) Lyrics() string {
	path := filepath.Join(n.savePath, n.lyricFileName())
	file, err := os.Open(path)
	if err == nil {
		buf, _ := io.ReadAll(file)
		return string(buf)
	}
	lrcics := getLyricByID(n.id)
	err = os.WriteFile(path, []byte(lrcics), 0644)
	return lrcics
	
}
func (n *neteaseMusic) MusicName() string {
	return n.title
}
func (n *neteaseMusic) SingerName() string {
	return n.singer
}
func (n *neteaseMusic) Album() string {
	return n.album
}
func (n *neteaseMusic) AlbumPicture() string {
	return n.albumPic
}
func (n *neteaseMusic) CurTime() (time.Duration, error) {
	if n.decode == nil {
		return 0, NoDecodeError
	}
	return n.decode.CurTime(), nil
}
func (n *neteaseMusic) EndTime() (time.Duration, error) {
	if n.decode != nil {
		return n.decode.EndTime(), nil
	}
	
	if n.endTime == 0 {
		_ = n.getAudio()
	}
	
	return n.endTime, nil
	
}

// Play Implementation  music.Operator
func (n *neteaseMusic) Play(cb *music.Callback, volume float64) error {
	if n.decode != nil {
		n.decode.Play()
		return nil
	}
	reader, err := n.getReader()
	if err != nil {
		return err
	}
	// tryGetCache
	d, err := decode.NewDecoder(context.TODO(), model.MusicTypeMP3, reader, volume, cb)
	if err != nil {
		_ = reader.Close()
		return err
	}
	n.decode = d
	n.decode.Play()
	return nil
	
}
func (n *neteaseMusic) Pause() error {
	if n.decode == nil {
		return NoDecodeError
	}
	n.decode.Pause()
	return nil
}
func (n *neteaseMusic) Stop() error {
	if n.decode == nil {
		return NoDecodeError
	}
	n.decode.Stop()
	return nil
}
func (n *neteaseMusic) SetVolume(f float64) {
	if n.decode == nil {
		return
	}
	n.decode.SetVolume(f)
}
func (n *neteaseMusic) DownloadMusic() error {
	_, err := os.Stat(filepath.Join(n.savePath, n.getFileName()))
	if err == nil {
		return nil
	}
	return download(n.audioURL, filepath.Join(n.savePath, n.getFileName()))
}
func (n *neteaseMusic) GetMusic() (model.Music, error) {
	if err := n.DownloadMusic(); err != nil {
		return model.Music{}, err
	}
	n.Lyrics()
	e, _ := n.EndTime()
	return model.Music{
		Name:   n.title,
		Singer: n.singer,
		Album:  n.album,
		Length: e,
		Type:   model.MusicTypeMP3,
		Path:   filepath.Join(n.savePath, n.getFileName()),
		Lyric:  filepath.Join(n.savePath, n.lyricFileName()),
	}, nil
}
func (n *neteaseMusic) Seek(f float64) {
	if n.decode == nil {
		return
	}
	n.decode.Seek(f)
}
func (n *neteaseMusic) Update(model model.Music) {}

// AddListener Implementation  binding.DataItem
func (n *neteaseMusic) AddListener(_ binding.DataListener) {

}
func (n *neteaseMusic) RemoveListener(_ binding.DataListener) {

}
func download(url, path string) error {
	r, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http.Get:%s", err.Error())
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("r.StatusCode:%d", r.StatusCode)
	}
	defer func() {
		_ = r.Body.Close()
	}()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll:%s", err.Error())
	}
	err = os.WriteFile(path, b, 0755)
	if err != nil {
		return fmt.Errorf("os.WriteFile:%s", err.Error())
	}
	return nil
}

func getLyricByID(id int) string {
	uuu := fmt.Sprintf("http://%s/lyric?id=%d", server, id)
	r, err := http.Get(uuu)
	if err != nil {
		log.Println("歌词获取失败：", err)
		return ""
	}
	defer func() {
		_ = r.Body.Close()
	}()
	bbb, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("歌词获取失败：", err)
		return ""
	}
	var v struct {
		Sgc       bool `json:"sgc"`
		Sfy       bool `json:"sfy"`
		Qfy       bool `json:"qfy"`
		LyricUser struct {
			Id       int    `json:"id"`
			Status   int    `json:"status"`
			Demand   int    `json:"demand"`
			Userid   int    `json:"userid"`
			Nickname string `json:"nickname"`
			Uptime   int64  `json:"uptime"`
		} `json:"lyricUser"`
		Lrc struct {
			Version int    `json:"version"`
			Lyric   string `json:"lyric"`
		} `json:"lrc"`
		Klyric struct {
			Version int    `json:"version"`
			Lyric   string `json:"lyric"`
		} `json:"klyric"`
		Tlyric struct {
			Version int    `json:"version"`
			Lyric   string `json:"lyric"`
		} `json:"tlyric"`
		Romalrc struct {
			Version int    `json:"version"`
			Lyric   string `json:"lyric"`
		} `json:"romalrc"`
		Code int `json:"code"`
	}
	err = json.Unmarshal(bbb, &v)
	if err != nil {
		log.Println("歌词获取失败：", err)
		return ""
	}
	return v.Lrc.Lyric
}
