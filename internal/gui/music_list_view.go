package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/mp"
	"github.com/Theodoree/music_player/internal/music"
	"k8s.io/klog"
	"time"
)

type listType int

const (
	ListType listType = iota
	StreamListType
)

type musicListView struct {
	mp   mp.MusicPlayer
	w    fyne.Window
	list struct {
		list        *widget.List
		allSelected int // 0:init 1:select 2:unselect
	}
	listContainer       *fyne.Container
	streamListContainer *fyne.Container
	
	selectEntry struct {
		*selectEntry
		
		curTableID    uint
		musicSelected []bool
	}
}

func newMusicListView(mp mp.MusicPlayer, w fyne.Window, selectEntry *selectEntry) *musicListView {
	var mlv = &musicListView{mp: mp, w: w}
	mlv.selectEntry.selectEntry = selectEntry
	return mlv
}
func (m *musicListView) view() fyne.CanvasObject {
	m.listContainer = container.NewBorder(m.toolBar(), nil, nil, nil, m.musicList())
	m.streamListContainer = container.NewBorder(m.streamToolBar(), nil, nil, nil, m.streamMusicList())
	m.streamListContainer.Hide()
	return container.NewBorder(nil, nil, nil, nil, m.listContainer, m.streamListContainer)
}
func (m *musicListView) toolBar() fyne.CanvasObject {
	_, _, keyword := m.mp.MusicList()
	
	importButton := widget.NewButton("导入", func() {
		dialog.ShowFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, m.w)
				return
			}
			if reader == nil {
				return
			}
			ml, index := m.mp.MusicTableList()
			idx, _ := index.Get()
			item, _ := ml.GetItem(idx)
			table := item.(model.MusicTable)
			m.mp.ImportMusic(table.ID, reader.Path())
			_ = index.Set(idx)
		}, m.w)
	})
	
	from := dialog.NewForm("添加到", "确认", "取消", []*widget.FormItem{widget.NewFormItem("列表名称", m.selectEntry.Select)}, func(ok bool) {
		defer func() {
			m.selectEntry.Refresh(false)
		}()
		if !ok {
			return
		}
		
		tableID := m.selectEntry.GetTableID()
		items, _, _ := m.mp.MusicList()
		for idx, selected := range m.selectEntry.musicSelected {
			if !selected {
				continue
			}
			_item, err := items.GetItem(idx)
			if err != nil {
				klog.Error(err)
				continue
			}
			item := _item.(music.Music)
			
			_music, err := item.GetMusic()
			if err != nil {
				klog.Error(err)
				continue
			}
			m.mp.AddMusic(tableID, _music)
		}
	}, m.w)
	addMusicButton := widget.NewButton("添加到", func() {
		from.Show()
	})
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索")
	searchEntry.OnSubmitted = func(k string) {
		_ = keyword.Set(k)
	}
	toolBox := container.NewBorder(nil, widget.NewSeparator(), container.NewHBox(importButton, addMusicButton), nil, searchEntry)
	
	selectLabel := widget.NewCheck("", func(b bool) {
		if b {
			m.list.allSelected = 1
		} else {
			m.list.allSelected = 2
		}
		m.list.list.Refresh()
	})
	titleLabel := widget.NewLabelWithStyle("歌曲名", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	singerLabel := widget.NewLabel("歌手")
	albumLabel := widget.NewLabel("专辑")
	playLabel := widget.NewLabelWithStyle("长度", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	buttonLabel := widget.NewLabelWithStyle("播放", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	tableHeader := container.NewGridWithColumns(6, selectLabel, titleLabel, singerLabel, albumLabel, playLabel, buttonLabel)
	return container.NewBorder(toolBox, nil, nil, nil, tableHeader)
}
func (m *musicListView) musicList() fyne.CanvasObject {
	items, index, _ := m.mp.MusicList()
	
	ml := widget.NewList(items.Length, func() fyne.CanvasObject {
		checkBox := widget.NewCheck("", nil)
		titleLabel := widget.NewLabelWithStyle("歌曲名", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		titleLabel.Truncation = fyne.TextTruncateEllipsis
		singerLabel := widget.NewLabel("歌手")
		singerLabel.Truncation = fyne.TextTruncateEllipsis
		album := widget.NewLabel("专辑")
		album.Truncation = fyne.TextTruncateEllipsis
		length := widget.NewLabel("长度")
		length.Truncation = fyne.TextTruncateEllipsis
		button := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
		return container.NewGridWithColumns(6, checkBox, titleLabel, singerLabel, album, length, button)
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		_item, _ := items.GetItem(id)
		o := object.(*fyne.Container)
		item := _item.(music.Music)
		
		gridColumns := o
		check := gridColumns.Objects[0].(*widget.Check)
		title := gridColumns.Objects[1].(*widget.Label)
		singerLabel := gridColumns.Objects[2].(*widget.Label)
		album := gridColumns.Objects[3].(*widget.Label)
		length := gridColumns.Objects[4].(*widget.Label)
		button := gridColumns.Objects[5].(*widget.Button)
		switch m.list.allSelected {
		case 1:
			check.SetChecked(true)
		case 2:
			check.SetChecked(false)
		default:
			check.SetChecked(false)
		}
		if m.selectEntry.curTableID != item.TableID() {
			m.selectEntry.musicSelected = make([]bool, items.Length())
			m.selectEntry.curTableID = item.TableID()
		}
		check.OnChanged = func(b bool) {
			if id >= len(m.selectEntry.musicSelected) {
				return
			}
			m.selectEntry.musicSelected[id] = b
		}
		button.OnTapped = func() {
			_ = index.Set(id - 1)
			m.mp.Next()
		}
		title.Text = item.MusicName()
		singerLabel.Text = item.SingerName()
		album.Text = item.Album()
		end, _ := item.EndTime()
		length.Text = fmt.Sprintf("%01d:%02d", end/time.Second/60, end/time.Second%60)
		
		check.Refresh()
		title.Refresh()
		singerLabel.Refresh()
		album.Refresh()
		length.Refresh()
		o.Show()
	})
	
	items.AddListener(&mp.DataListener{Fn: func() {
		ml.Refresh()
	}})
	m.list.list = ml
	return ml
}

func (m *musicListView) streamToolBar() fyne.CanvasObject {
	_, _, searchKey := m.mp.StreamMusicList()
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("搜索")
	searchEntry.OnSubmitted = func(s string) {
		_ = searchKey.Set(s)
	}
	
	importButton := widget.NewButton("导入", func() {})
	importButton.Disable()
	
	toolBox := container.NewBorder(nil, widget.NewSeparator(), importButton, nil, searchEntry)
	
	titleLabel := widget.NewLabelWithStyle("歌曲名", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	singerLabel := widget.NewLabel("歌手")
	albumLabel := widget.NewLabel("专辑")
	playLabel := widget.NewLabelWithStyle("长度", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	buttonLabel := widget.NewLabelWithStyle("播放", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	tableHeader := container.NewGridWithColumns(5, titleLabel, singerLabel, albumLabel, playLabel, buttonLabel)
	
	return container.NewBorder(toolBox, nil, nil, nil, tableHeader)
}
func (m *musicListView) streamMusicList() fyne.CanvasObject {
	items, index, _ := m.mp.StreamMusicList()
	ml := widget.NewList(items.Length, func() fyne.CanvasObject {
		titleLabel := widget.NewLabelWithStyle("歌曲名", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		titleLabel.Truncation = fyne.TextTruncateEllipsis
		singerLabel := widget.NewLabel("歌手")
		singerLabel.Truncation = fyne.TextTruncateEllipsis
		album := widget.NewLabel("专辑")
		album.Truncation = fyne.TextTruncateEllipsis
		length := widget.NewLabel("长度")
		length.Truncation = fyne.TextTruncateEllipsis
		button := widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {})
		return container.NewGridWithColumns(5, titleLabel, singerLabel, album, length, button)
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		_item, _ := items.GetItem(id)
		o := object.(*fyne.Container)
		item := _item.(music.Music)
		
		gridColumns := o
		title := gridColumns.Objects[0].(*widget.Label)
		singerLabel := gridColumns.Objects[1].(*widget.Label)
		album := gridColumns.Objects[2].(*widget.Label)
		length := gridColumns.Objects[3].(*widget.Label)
		button := gridColumns.Objects[4].(*widget.Button)
		button.OnTapped = func() {
			_ = index.Set(id - 1)
			m.mp.Next()
		}
		title.Text = item.MusicName()
		singerLabel.Text = item.SingerName()
		album.Text = item.Album()
		end, _ := item.EndTime()
		length.Text = fmt.Sprintf("%01d:%02d", end/time.Second/60, end/time.Second%60)
		title.Refresh()
		singerLabel.Refresh()
		album.Refresh()
		length.Refresh()
	})
	items.AddListener(&mp.DataListener{Fn: func() {
		ml.Refresh()
	}})
	
	return ml
}

func (m *musicListView) SwapTable(t listType) {
	switch t {
	case ListType:
		if m.listContainer.Visible() {
			return
		}
		m.listContainer.Show()
		m.streamListContainer.Hide()
	case StreamListType:
		if m.streamListContainer.Visible() {
			return
		}
		m.streamListContainer.Show()
		m.listContainer.Hide()
		
	}
}

func (m *musicListView) Refresh() {
	if m.listContainer.Visible() {
		m.listContainer.Refresh()
	}
	if m.streamListContainer.Visible() {
		m.streamListContainer.Refresh()
	}
}
