package mp

import (
	"github.com/Theodoree/music_player/internal/music"
	"math/rand/v2"
	"strings"
)

type PlayMode int

const (
	PlayModeCycle PlayMode = iota
	PlayModeSingleCycle
	PlayModeRandom
)

type list struct {
	tableId uint
	bindingTable[music.Music]
	tmp       []music.Music
	searchKey BindingModel[string]
}

func (t *list) setItems(items []music.Music, tableId uint, cache bool) {
	if cache {
		t.tmp = make([]music.Music, len(items))
		copy(t.tmp, items)
	}
	t.tableId = tableId
	t.items.SetItems(items)
	_ = t.bindingTable.index.Set(-1)
}
func (t *list) prev(mode PlayMode) music.Music {
	index := t.index.get()
	switch mode {
	case PlayModeSingleCycle:
	case PlayModeCycle:
		index -= 1
		if index == -1 {
			index = t.items.Length() - 1
		}
	case PlayModeRandom:
		index = rand.IntN(t.items.Length())
	}
	_ = t.index.Set(index)
	return t.items.items[index]
}
func (t *list) next(mode PlayMode) music.Music {
	index := t.index.get()
	switch mode {
	case PlayModeSingleCycle:
	case PlayModeCycle:
		index += 1
		index %= t.items.Length()
	case PlayModeRandom:
		index = rand.IntN(t.items.Length())
	}
	_ = t.index.Set(index)
	return t.items.items[index]
}
func (t *list) valid() bool {
	return t.items.Length() > 0
}

func (t *list) Search(keyword string) {
	switch {
	case len(keyword) == 0:
		t.setItems(t.tmp, t.tableId, false)
	case len(keyword) > 0:
		var tmp []music.Music
		for i := 0; i < len(t.tmp); i++ {
			m := t.tmp[i]
			if !(strings.Index(m.SingerName(), keyword) >= 0 || strings.Index(m.MusicName(), keyword) >= 0 || strings.Index(m.Album(), keyword) >= 0) {
				continue
			}
			tmp = append(tmp, m)
		}
		t.setItems(tmp, t.tableId, false)
	}
}
