package gui

import (
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/db"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/mp"
)

type selectEntry struct {
	Select   *widget.Select
	refresh  func()
	tableIds []uint
}

func newSelectEntry(mp mp.MusicPlayer) *selectEntry {
	fn := func() ([]string, []uint) {
		items, _ := mp.MusicTableList()
		var options []string
		var idx []uint
		for i := 0; i < items.Length(); i++ {
			_item, _ := items.GetItem(i)
			item := _item.(model.MusicTable)
			if item.ID == db.DefaultTableID {
				continue
			}
			options = append(options, item.Name)
			idx = append(idx, item.ID)
		}
		return options, idx
	}
	option, idx := fn()
	var s selectEntry
	s.Select = widget.NewSelect(option, func(s string) {})
	s.tableIds = idx
	s.refresh = func() {
		option, idx := fn()
		s.Select.SetOptions(option)
		s.tableIds = idx
	}
	return &s
}

func (s *selectEntry) GetTableID() uint {
	return s.tableIds[s.Select.SelectedIndex()]
}
func (s *selectEntry) Refresh(optionRefresh bool) {
	if optionRefresh {
		s.refresh()
	}
	s.Select.Selected = ""
	s.Select.Refresh()
}
