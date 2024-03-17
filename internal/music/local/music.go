package local

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"time"
	
	"github.com/Theodoree/music_player/internal/db"
	"github.com/Theodoree/music_player/internal/decode"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/music"
)

var NoDecodeError = errors.New("no decode")

type localSource struct {
	ctx    context.Context
	cancel context.CancelFunc
	db     db.MusicStore
}

func Source(ctx context.Context, db db.MusicStore) music.Source {
	var n localSource
	n.ctx, n.cancel = context.WithCancel(ctx)
	n.db = db
	return &n
}

func (api *localSource) SearchMusic(tableID uint, keyWord string) ([]music.Music, error) {
	musics, err := api.db.GetMusicByMusicTableID(tableID)
	if err != nil {
		return nil, err
	}
	
	var items []music.Music
	for _, m := range musics {
		if len(keyWord) > 0 && strings.Index(keyWord, m.Name) == -1 && strings.Index(keyWord, m.Singer) == -1 {
			continue
		}
		items = append(items, newMusic(m))
	}
	
	return items, nil
}
func (api *localSource) List(tableID uint) ([]music.Music, error) {
	return api.SearchMusic(tableID, "")
}
func (api *localSource) Close() error {
	api.cancel()
	return nil
}

type _music struct {
	model.Music
	decode decode.Decoder
}

func newMusic(m model.Music) music.Music {
	return &_music{Music: m}
}

func (n *_music) getReader() (io.ReadSeekCloser, error) {
	file, err := os.Open(n.Music.Path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// TableID Implementation music.Data
func (n *_music) TableID() uint {
	return n.MusicTableID
}
func (n *_music) Lyrics() string {
	buf, _ := os.ReadFile(n.Music.Lyric)
	return string(buf)
	
}
func (n *_music) MusicName() string {
	return n.Music.Name
}
func (n *_music) SingerName() string {
	return n.Music.Singer
}
func (n *_music) Album() string {
	return n.Music.Album
}
func (n *_music) AlbumPicture() string {
	//return n.Music.Pic
	return ""
}
func (n *_music) CurTime() (time.Duration, error) {
	if n.decode == nil {
		return 0, NoDecodeError
	}
	if n.decode == nil {
		return 0, NoDecodeError
	}
	return n.decode.CurTime(), nil
}
func (n *_music) EndTime() (time.Duration, error) {
	return n.Music.Length, nil
}

// Play Implementation music.Operator
func (n *_music) Play(cb *music.Callback, volume float64) error {
	if n.decode != nil {
		n.decode.Play()
		return nil
	}
	reader, err := n.getReader()
	if err != nil {
		return err
	}
	// tryGetCache
	decoder, err := decode.NewDecoder(context.TODO(), n.Type, reader, volume, cb)
	if err != nil {
		_ = reader.Close()
		return err
	}
	n.decode = decoder
	n.decode.Play()
	return nil
	
}
func (n *_music) Pause() error {
	if n.decode == nil {
		return NoDecodeError
	}
	n.decode.Pause()
	return nil
}
func (n *_music) Stop() error {
	if n.decode == nil {
		return NoDecodeError
	}
	n.decode.Stop()
	n.decode = nil
	return nil
}
func (n *_music) SetVolume(f float64) {
	if n.decode == nil {
		return
	}
	n.decode.SetVolume(f)
}
func (n *_music) GetMusic() (model.Music, error) {
	return n.Music, nil
}
func (n *_music) Seek(f float64) {
	if n.decode == nil {
		return
	}
	n.decode.Seek(f)
}
func (n *_music) Update(model model.Music) {
	n.Music = model
}
