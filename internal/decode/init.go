package decode

import (
	"context"
	"errors"
	"github.com/Theodoree/music_player/internal/music"
	"io"
	"time"
	
	"github.com/Theodoree/music_player/internal/model"
)

type Metadata struct {
	// 音乐名
	Title string
	// 歌手名
	Artist string
	// 专辑名称
	Album string
	// 专辑页面(BASE64)
	AlbumPicture string
	// 专辑页面类型
	MimeType string
	// 结束时间
	EndTime time.Duration
}

type Decoder interface {
	// Play 播放音乐
	Play()
	// Pause 暂停音乐
	Pause()
	// Stop 停止音乐
	Stop()
	// SetVolume [0.01,1]
	SetVolume(f float64)
	// Metadata 获取音乐元数据
	Metadata() (Metadata, error)
	// CurTime 获取当前播放时间
	CurTime() time.Duration
	// EndTime 获取音乐结束时间
	EndTime() time.Duration
	// Seek 跳转到指定百分比[0,100]
	Seek(f float64)
}

// FLAC、MP3、WAV decoder power by @github.com/faiface/beep
func NewDecoder(ctx context.Context, MusicType model.MusicType, reader io.ReadSeekCloser, volume float64, cb *music.Callback) (Decoder, error) {
	switch MusicType {
	case model.MusicTypeFLAC:
		return newFLACDecoder(ctx, reader, volume, cb)
	case model.MusicTypeMP3:
		return newMP3Decoder(ctx, reader, volume, cb)
	case model.MusicTypeWAV:
		return newWAVDecoder(ctx, reader, volume, cb)
	default:
		return nil, errors.New("unhandled default case")
	}
}

func GetMetadata(filepath string, MusicType model.MusicType) Metadata {
	return Metadata{}
}
