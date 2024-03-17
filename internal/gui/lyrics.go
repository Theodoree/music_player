package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/mp"
	"k8s.io/klog"
)

type lyricsView struct {
	mp mp.MusicPlayer
	w  fyne.Window
}

func newLyricsView(mp mp.MusicPlayer, w fyne.Window) *lyricsView {
	return &lyricsView{mp: mp, w: w}
}
func (t *lyricsView) view() fyne.CanvasObject {
	
	// 创建一个足够长的歌词字符串（这里仅为示例，实际应用中应根据歌词文件生成）
	lyrics := t.mp.Lyrics()
	lyricsLabel := widget.NewLabelWithData(lyrics)
	
	button := widget.NewButton("导入歌词", func() {
		cur := t.mp.GetPlayedMusic()
		if cur == nil {
			return
		}
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				klog.Info(err)
			}
			if closer == nil {
				return
			}
			m, _ := cur.GetMusic()
			m.Lyric = closer.URI().Path()
			t.mp.UpdateMusic(m.MusicTableID, m)
		}, t.w)
	})
	button.Alignment = widget.ButtonAlignCenter
	
	lyrics.AddListener(&mp.DataListener{Fn: func() {
		l, _ := lyrics.Get()
		if len(l) == 0 {
			button.Text = "导入歌词"
			button.Show()
		} else {
			button.Text = "更新歌词"
		}
		button.Refresh()
	}})
	
	// 设置歌词标签的最大行数，以模拟定长（例如，10行）
	lyricsLabel.Alignment = fyne.TextAlignCenter
	lyricsLabel.Wrapping = fyne.TextWrapWord // 使歌词换行
	
	scrollContent := container.New(layout.NewVBoxLayout(), lyricsLabel)
	scroll := container.NewScroll(scrollContent)
	scroll.SetMinSize(fyne.NewSize(400, 0))
	
	return container.NewBorder(button, nil, widget.NewSeparator(), nil, widget.NewCard("", "", scroll))
}
