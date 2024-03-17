package tool

import (
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"k8s.io/klog"
	"os"
	"path/filepath"
	"time"
	
	"github.com/Theodoree/music_player/internal/model"
	"github.com/dhowden/tag"
)

type File struct {
	path string
	name string
}

func recursiveSearchDirectory(path string, consumer func(entry os.DirEntry, path string)) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return
	}
	
	for _, v := range dirEntry {
		if v.IsDir() {
			recursiveSearchDirectory(filepath.Join(path, v.Name()), consumer)
			continue
		}
		consumer(v, filepath.Join(path, v.Name()))
	}
}

func SearchMusicFileByPath(path string) []model.Music {
	path, _ = filepath.Abs(path)
	var items []model.Music
	consumer := func(entry os.DirEntry, path string) {
		t, ok := model.IsMusicType(entry.Name())
		if !ok {
			return
		}
		ms := model.Music{}
		ms.Name = entry.Name()
		ms.Type = t
		ms.Path = path
		if err := getMetadata(path, &ms); err != nil {
			klog.Error(err)
			return
		}
		items = append(items, ms)
	}
	recursiveSearchDirectory(path, consumer)
	return items
}

func getMetadata(fileName string, music *model.Music) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	m, err := tag.ReadFrom(file)
	if err != nil {
		klog.Error(err)
		return nil
	}
	
	if m.Title() != "" {
		music.Name = m.Title()
	}
	if music.Singer == "" {
		music.Singer = m.AlbumArtist()
	}
	music.Album = m.Album()
	//pic := m.Picture()
	//if pic != nil {
	//	buf, _ := json.Marshal(pic)
	//	music.Pic = string(buf)
	//}
	
	switch m.FileType() {
	case tag.MP3:
		_, _ = file.Seek(0, 0)
		streamer, format, err := mp3.Decode(file)
		if err != nil {
			klog.Error(err)
			_, _, err = flac.Decode(file)
			klog.Error(err)
			return err
		}
		format.SampleRate.D(streamer.Len()).Round(time.Second)
		music.Length = format.SampleRate.D(streamer.Len()).Round(time.Second)
	case tag.FLAC:
		_, _ = file.Seek(0, 0)
		streamer, format, err := flac.Decode(file)
		if err != nil {
			klog.Error(err)
			return err
		}
		format.SampleRate.D(streamer.Len()).Round(time.Second)
		music.Length = format.SampleRate.D(streamer.Len()).Round(time.Second)
	}
	
	return nil
}
