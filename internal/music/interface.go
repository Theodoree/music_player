package music

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/Theodoree/music_player/internal/model"
	"time"
)

type Source interface {
	SearchMusic(tableID uint, keyword string) ([]Music, error)
	List(tableID uint) ([]Music, error)
	Close() error
}

type Callback struct {
	CurTime func(duration time.Duration)
	DoneFn  func(model.Status)
}

// Data 音乐数据接口
type Data interface {
	binding.DataItem
	// TableID 返回表格ID
	TableID() uint
	// Lyrics 返回歌词(全量)
	Lyrics() string
	// MusicName 返回音乐名
	MusicName() string
	// SingerName 返回歌手名
	SingerName() string
	// Album 专辑名
	Album() string
	// AlbumPicture 专辑封面
	AlbumPicture() string
	// CurTime 返回当前播放时间
	CurTime() (time.Duration, error)
	// EndTime 返回音乐结束时间
	EndTime() (time.Duration, error)
}

// Operator 音乐操作接口
type Operator interface {
	// Play 播放音乐
	Play(cb *Callback, volume float64) error
	// Pause 暂停音乐
	Pause() error
	// Stop 停止音乐
	Stop() error
	// SetVolume [0,100]
	SetVolume(f float64)
	
	GetMusic() (model.Music, error)
	
	Seek(f float64)
	
	Update(model model.Music)
}

// Music 音乐接口
type Music interface {
	Operator
	Data
}
