package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/mp"
)

type controller struct {
	//mp mp.MusicPlayer
	//menu struct {
	//	play  *widget.Button
	//	pause *widget.Button
	//	prev  *widget.Button
	//	next  *widget.Button
	//}
	//progress struct {
	//	cur         *widget.Label // 进度条
	//	Progressbar *widget.Slider
	//	end         *widget.Label // 进度条
	//}
	
	//volumeProgress *widget.Slider
	
	//playModeSelect *widget.Select
}

func newController() *controller {
	return &controller{}
}

func (b *controller) View(musicPlayer mp.MusicPlayer) fyne.CanvasObject {
	
	// 音乐名和歌手名组件
	musicName := widget.NewLabelWithStyle("          ", fyne.TextAlignCenter, fyne.TextStyle{})
	musicName.Bind(musicPlayer.MusicName())
	playerName := widget.NewLabelWithStyle("          ", fyne.TextAlignCenter, fyne.TextStyle{})
	playerName.Bind(musicPlayer.SingerName())
	left := container.NewVBox(musicName, playerName)
	
	// 播放控制组件
	var playerMenu struct {
		play  *widget.Button
		pause *widget.Button
		prev  *widget.Button
		next  *widget.Button
	}
	playerMenu.play = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), musicPlayer.Play)
	playerMenu.pause = widget.NewButtonWithIcon("", theme.MediaPauseIcon(), musicPlayer.Pause)
	playerMenu.pause.Hide()
	playerMenu.prev = widget.NewButtonWithIcon("", theme.MediaSkipPreviousIcon(), musicPlayer.Prev)
	playerMenu.next = widget.NewButtonWithIcon("", theme.MediaSkipNextIcon(), musicPlayer.Next)
	playStatus := musicPlayer.PlayStatus()
	playStatus.AddListener(&mp.DataListener{Fn: func() {
		ok, _ := playStatus.Get()
		switch ok {
		case true:
			playerMenu.play.Hide()
			playerMenu.pause.Show()
		case false:
			playerMenu.pause.Hide()
			playerMenu.play.Show()
		}
	}})
	midder := container.NewGridWithColumns(3, playerMenu.prev, playerMenu.pause, playerMenu.play, playerMenu.next)
	
	// 音量条
	volumeProgress := widget.NewSlider(0, 100)
	volumeProgress.Bind(musicPlayer.Volume())
	volumeProgress.Refresh()
	val, _ := musicPlayer.Volume().Get()
	volumeProgress.SetValue(val)
	
	// 播放模式
	playModeSelect := widget.NewSelect([]string{"顺序播放", "单曲循环", "随机播放"}, func(s string) {
		mode := musicPlayer.PlayMode().(*mp.BindingModel[mp.PlayMode])
		switch s {
		case "顺序播放":
			_ = mode.Set(mp.PlayModeCycle)
		case "单曲循环":
			_ = mode.Set(mp.PlayModeSingleCycle)
		case "随机播放":
			_ = mode.Set(mp.PlayModeRandom)
		}
	})
	playModeSelect.SetSelectedIndex(0)
	right := container.NewVBox(playModeSelect, volumeProgress)
	topContainer := container.NewGridWithColumns(3, left, midder, right)
	
	var progressWidget struct {
		cur         *widget.Label // 进度条
		Progressbar *widget.Slider
		end         *widget.Label // 进度条
	}
	// 播放条
	progressWidget.cur = widget.NewLabel("00:00")
	progressWidget.cur.Bind(musicPlayer.MusicCurTime())
	progressWidget.Progressbar = widget.NewSlider(0, 100)
	progressWidget.Progressbar.Bind(musicPlayer.ProcessBar())
	progressWidget.end = widget.NewLabel("00:00")
	progressWidget.end.Bind(musicPlayer.MusicEndTime())
	
	bottomContainer := container.NewBorder(nil, nil, progressWidget.cur, progressWidget.end, progressWidget.Progressbar)
	
	return container.NewVBox(topContainer, bottomContainer)
}
