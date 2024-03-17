package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/db"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/mp"
)

type musicTableView struct {
	fn          mp.MusicPlayerOperationModule
	data        mp.MusicPlayerDataModule
	w           fyne.Window
	selectEntry *selectEntry
}

func newMusicTableView(mp mp.MusicPlayer, w fyne.Window, selectEntry *selectEntry) *musicTableView {
	return &musicTableView{fn: mp, data: mp, w: w, selectEntry: selectEntry}
}

func (m *musicTableView) view(swap func(t listType)) fyne.CanvasObject {
	label0 := widget.NewLabelWithStyle("云音乐", fyne.TextAlignCenter, fyne.TextStyle{})
	streamMusicTable := m.streamTable(swap)
	musicTable := m.table(swap)
	
	label := container.NewBorder(container.NewVBox(label0, streamMusicTable, widget.NewSeparator(), m.addTableButton(), m.delButton(), widget.NewSeparator()), nil, nil, nil, musicTable)
	return container.NewBorder(nil, nil, nil, widget.NewSeparator(), label)
}
func (m *musicTableView) streamTable(swap func(t listType)) fyne.CanvasObject {
	items, index := m.data.StreamMusicTableList()
	ml := widget.NewList(
		items.Length,
		func() fyne.CanvasObject {
			bt := widget.NewButton("", func() {})
			return bt
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			_item, _ := items.GetItem(id)
			o := object.(*widget.Button)
			item := _item.(*mp.BindingModel[string])
			o.SetIcon(theme.MediaMusicIcon())
			o.Text, _ = item.Get()
			o.OnTapped = func() {
				_ = index.Set(id)
				swap(StreamListType)
			}
			o.Refresh()
		})
	items.AddListener(&mp.DataListener{Fn: func() {
		ml.Refresh()
	}})
	return ml
}
func (m *musicTableView) table(swap func(t listType)) fyne.CanvasObject {
	items, index := m.data.MusicTableList()
	ml := widget.NewList(
		items.Length,
		func() fyne.CanvasObject {
			bt := widget.NewButton("", func() {})
			return bt
		},
		func(id widget.ListItemID, object fyne.CanvasObject) {
			o := object.(*widget.Button)
			_item, _ := items.GetItem(id)
			item := _item.(model.MusicTable)
			if item.ID == db.DefaultTableID {
				o.SetIcon(theme.HomeIcon())
			}
			o.Text = item.Name
			o.OnTapped = func() {
				_ = index.Set(id)
				swap(ListType)
			}
			o.Refresh()
		})
	items.AddListener(&mp.DataListener{Fn: func() {
		ml.Refresh()
	}})
	return ml
}
func (m *musicTableView) addTableButton() fyne.CanvasObject {
	nameEntry := widget.NewEntry()
	from := dialog.NewForm("新增列表", "确认", "取消", []*widget.FormItem{widget.NewFormItem("列表名称", nameEntry)}, func(ok bool) {
		defer func() {
			nameEntry.Text = ""
			m.selectEntry.Refresh(ok)
		}()
		if !ok {
			return
		}
		m.fn.AddTable(model.MusicTable{
			Name: nameEntry.Text,
		})
	}, m.w)
	
	from.Resize(fyne.NewSize(400, 0))
	return widget.NewButtonWithIcon("新增列表", theme.ContentAddIcon(), func() {
		from.Show()
	})
}

func (m *musicTableView) delButton() fyne.CanvasObject {
	
	from := dialog.NewForm("删除列表", "确认", "取消", []*widget.FormItem{widget.NewFormItem("列表名称", m.selectEntry.Select)}, func(ok bool) {
		defer func() {
			m.selectEntry.Refresh(ok)
		}()
		if !ok {
			return
		}
		m.fn.DelTable(m.selectEntry.GetTableID())
	}, m.w)
	
	from.Resize(fyne.NewSize(400, 0))
	return widget.NewButtonWithIcon("删除列表", theme.ContentRemoveIcon(), func() {
		from.Show()
	})
}
