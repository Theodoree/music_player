package module

import (
	"github.com/Theodoree/music_player/internal/model"
	"math/rand/v2"
	"slices"
)

type PlaybackMode int

const (
	PlaybackModeCycle PlaybackMode = iota
	PlaybackModeSingleCycle
	PlaybackModeRandom
)

type list struct {
	BindingDataList[model.Music]
	idx    int
	mode   PlaybackMode
	signal BindingModel[struct{}]
}

func (t *list) SearchByID(musicID uint) int {
	return slices.IndexFunc(t.items, func(music model.Music) bool {
		return music.ID == musicID
	})
}

func (t *list) setItems(items []model.Music) {
	t.items = items
	t.idx = 0
	t.signal.Signal()
}
func (t *list) setMode(mode PlaybackMode) {
	t.mode = mode
}
func (t *list) setIdx(idx int) {
	if idx >= len(t.items) {
		panic("out of range")
	}
	t.idx = idx
}
func (t *list) prev() model.Music {
	switch t.mode {
	case PlaybackModeSingleCycle:
	case PlaybackModeCycle:
		t.idx -= 1
		if t.idx == -1 {
			t.idx = len(t.items) - 1
		}
	case PlaybackModeRandom:
		t.idx = rand.IntN(len(t.items))
	}
	return t.items[t.idx]
}
func (t *list) next() model.Music {
	switch t.mode {
	case PlaybackModeSingleCycle:
	case PlaybackModeCycle:
		t.idx += 1
		t.idx %= len(t.items)
	case PlaybackModeRandom:
		t.idx = rand.IntN(len(t.items))
	}
	return t.items[t.idx]
}
func (t *list) valid() bool {
	return len(t.items) > 0
}
