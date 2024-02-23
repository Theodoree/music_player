package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/module"
)

type controller struct {
	op module.MusicPlayerOperation
	// 左侧显示框
	stat struct {
		//picture
		musicName  *widget.Label
		playerName *widget.Label
	}
	
	// menu
	menu struct {
		play  *widget.Button
		pause *widget.Button
		prev  *widget.Button
		next  *widget.Button
	}
	progress struct {
		cur         *widget.Label // 进度条
		Progressbar *widget.Slider
		end         *widget.Label // 进度条
	}
}

func newController(op module.MusicPlayerOperation) *controller {
	return &controller{op: op}
}

func (b *controller) View() fyne.CanvasObject {
	b.stat.musicName = widget.NewLabel("perfect")
	b.stat.musicName.Bind(b.op.GetMusicName())
	b.stat.playerName = widget.NewLabel("Ed Sheeran")
	b.stat.playerName.Bind(b.op.GetSingerName())
	
	statView := container.NewVBox(container.NewCenter(container.NewVBox(b.stat.musicName, b.stat.playerName)))
	
	// menu
	b.menu.play = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), func() {
		if b.op.Play() {
			b.menu.play.Hide()
			b.menu.pause.Show()
		}
	})
	b.menu.pause = widget.NewButtonWithIcon("", theme.MediaPauseIcon(), func() {
		if b.op.Pause() {
			b.menu.pause.Hide()
			b.menu.play.Show()
		}
	})
	b.menu.pause.Hide()
	b.menu.prev = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), b.op.Prev)
	b.menu.next = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), b.op.Next)
	
	menu := container.NewGridWithColumns(4, b.menu.prev, b.menu.pause, b.menu.play, b.menu.next)
	
	// 播放条
	b.progress.cur = widget.NewLabel("00:00")
	b.progress.cur.Bind(b.op.GetMusicCurTime())
	b.progress.Progressbar = widget.NewSlider(0, 100)
	b.progress.Progressbar.Bind(b.op.GetProcessBar())
	b.progress.end = widget.NewLabel("00:00")
	b.progress.end.Bind(b.op.GetMusicEndTime())
	
	l1 := container.NewGridWithColumns(3, statView, menu)
	line1 := container.NewBorder(nil, nil, nil, nil, l1)
	
	t2 := container.NewBorder(nil, nil, b.progress.cur, b.progress.end, b.progress.Progressbar)
	
	return container.NewVBox(line1, t2)
}
