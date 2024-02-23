package gui

import (
	"context"
	"embed"
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/module"
)

type _theme struct {
	fyne.Theme
	font fyne.Resource
}

func (t _theme) Font(fyne.TextStyle) fyne.Resource {
	return t.font
}

type gui struct {
	w          fyne.Window
	op         module.MusicPlayerOperation
	controller *controller
}

func Run(a fyne.App, res *embed.FS) {
	w := a.NewWindow("musicplayer")
	buf, _ := res.ReadFile("resource/font/simkai.ttf")
	a.Settings().SetTheme(_theme{Theme: theme.DefaultTheme(), font: fyne.NewStaticResource("simkai.ttf", buf)})
	var gui gui
	gui.w = w
	gui.op = module.NewController(context.Background(), func(str string) {
		dialog.ShowError(errors.New(str), w)
	})
	gui.w.Resize(fyne.NewSize(1024, 768))
	gui.w.SetMaster()
	gui.View()
	gui.w.ShowAndRun()
}

func (gui *gui) View() {
	w := gui.w
	
	label := widget.NewLabel("本地歌单")
	musicTable := container.NewBorder(label, nil, nil, widget.NewSeparator(), newMusicTable(gui.op).view())
	controller := newController(gui.op)
	musicToolBar := newMusicToolBar(gui.op)
	musicListView := newMusicList(gui.op)
	
	b1 := container.NewBorder(musicToolBar.view(w), nil, nil, nil, musicListView.view(func() {
		if controller.menu.play.Hidden {
			return
		}
		controller.menu.play.Hide()
		controller.menu.pause.Show()
		
	}))
	b2 := container.NewBorder(nil, nil, musicTable, nil, b1)
	
	controllerBox := container.NewBorder(widget.NewSeparator(), nil, nil, nil, controller.View())
	w.SetContent(container.NewBorder(nil, controllerBox, nil, nil, b2))
	w.Show()
}

type musicTable struct {
	op module.MusicPlayerOperation
}

func newMusicTable(op module.MusicPlayerOperation) *musicTable {
	return &musicTable{op: op}
}

func (t musicTable) view() fyne.CanvasObject {
	ml := widget.NewListWithData(
		t.op.GetMusicTableList(),
		func() fyne.CanvasObject {
			bt := widget.NewButton("", func() {})
			return bt
		},
		func(_item binding.DataItem, object fyne.CanvasObject) {
			o := object.(*widget.Button)
			item := _item.(model.MusicTable)
			o.Text = item.Name
			o.OnTapped = func() {
				t.op.SelectTable(item.ID)
			}
			o.Refresh()
		})
	return ml
}

type musicList struct {
	op module.MusicPlayerOperation
}

func newMusicList(op module.MusicPlayerOperation) *musicList {
	return &musicList{op: op}
}

func (t *musicList) view(play func()) fyne.CanvasObject {
	ml := widget.NewListWithData(
		t.op.GetMusicList(),
		func() fyne.CanvasObject {
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
		},
		func(_item binding.DataItem, object fyne.CanvasObject) {
			o := object.(*fyne.Container)
			item := _item.(model.Music)
			
			gridColumns := o
			title := gridColumns.Objects[0].(*widget.Label)
			singerLabel := gridColumns.Objects[1].(*widget.Label)
			album := gridColumns.Objects[2].(*widget.Label)
			length := gridColumns.Objects[3].(*widget.Label)
			button := gridColumns.Objects[4].(*widget.Button)
			button.OnTapped = func() {
				t.op.SelectMusic(item.ID)
				play()
				
			}
			title.Text = item.Name
			singerLabel.Text = item.Singer
			album.Text = item.Album
			length.Text = item.Length
			
			title.Refresh()
			singerLabel.Refresh()
			album.Refresh()
			length.Refresh()
		})
	
	return ml
}

type musicToolBar struct {
	op module.MusicPlayerOperation
}

func newMusicToolBar(op module.MusicPlayerOperation) *musicToolBar {
	return &musicToolBar{op: op}
}
func (t *musicToolBar) view(w fyne.Window) fyne.CanvasObject {
	importButton := widget.NewButton("导入", func() {
		dialog.ShowFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return
			}
			t.op.SaveToTable(reader.Path())
		}, w)
	})
	
	titleLabel := widget.NewLabelWithStyle("歌曲名", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	singerLabel := widget.NewLabel("歌手")
	albumLabel := widget.NewLabel("专辑")
	playLable := widget.NewLabelWithStyle("长度", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	buttonLable := widget.NewLabelWithStyle("播放", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	
	t2 := container.NewGridWithColumns(5, titleLabel, singerLabel, albumLabel, playLable, buttonLable)
	return container.NewBorder(importButton, nil, nil, nil, t2)
}
