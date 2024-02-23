package tool

import (
	"encoding/json"
	"os"
	"path/filepath"
	
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
		getMetadata(path, &ms)
		items = append(items, ms)
	}
	recursiveSearchDirectory(path, consumer)
	return items
}

func getMetadata(fileName string, music *model.Music) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	m, err := tag.ReadFrom(file)
	if err != nil {
		return
	}
	
	music.Name = m.Title()
	music.Singer = m.Artist()
	if music.Singer == "" {
		music.Singer = m.AlbumArtist()
	}
	music.Album = m.Album()
	pic := m.Picture()
	if pic != nil {
		buf, _ := json.Marshal(pic)
		music.Pic = string(buf)
	}
	
}
