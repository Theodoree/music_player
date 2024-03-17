package gui

import (
	"context"
	"embed"
	"errors"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Theodoree/music_player/internal/mp"
)

type _theme struct {
	fyne.Theme
	font fyne.Resource
}

func (t _theme) Font(fyne.TextStyle) fyne.Resource {
	return t.font
}

type gui struct {
	ctx context.Context
	//mp  mp.MusicPlayer
	
	//musicListView  *musicListView
	//musicTableView *musicTableView
	//controller     *controller
}

func Run(a fyne.App, res *embed.FS) {
	w := a.NewWindow("musicplayer")
	buf, err := res.ReadFile("resource/font/simkai.ttf")
	a.Settings().SetTheme(_theme{Theme: theme.DefaultTheme(), font: fyne.NewStaticResource("simkai.ttf", buf)})
	var gui gui
	gui.ctx = context.Background()
	musicPlayer, err := mp.NewMusicPlayer(gui.ctx, func(str string) {
		dialog.ShowError(errors.New(str), w)
	})
	if err != nil {
		panic(err)
	}
	w.Resize(fyne.NewSize(1024, 768))
	w.SetMaster()
	gui.InitMenu(w, musicPlayer)
	gui.View(w, musicPlayer)
	w.ShowAndRun()
}

func (app *gui) InitMenu(window fyne.Window, musicPlayer mp.MusicPlayer) {
	// 添加子菜单项到“File”菜单下
	pauseMenuItem := fyne.NewMenuItem("暂停", musicPlayer.Play)
	prevMenuItem := fyne.NewMenuItem("上一首", musicPlayer.Prev)
	nextMenuItem := fyne.NewMenuItem("下一首", musicPlayer.Next)
	// 创建主菜单并添加菜单项
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("播放控制", pauseMenuItem, prevMenuItem, nextMenuItem),
	)
	// 设置窗口的菜单栏
	window.SetMainMenu(mainMenu)
}

func (app *gui) View(window fyne.Window, musicPlayer mp.MusicPlayer) {
	
	entry := newSelectEntry(musicPlayer)
	// 表格列表视图
	musicTableView := newMusicTableView(musicPlayer, window, entry)
	// 音乐列表视图
	musicListView := newMusicListView(musicPlayer, window, entry)
	// 歌词视图
	lyrics := newLyricsView(musicPlayer, window)
	// 控制器视图(歌手、音乐名、播放进度条、音量、播放控制按钮)
	controller := newController()
	
	// 组装容器
	// left: 表格列表视图  middle: 音乐列表视图 right: 歌词视图
	topContainer := container.NewBorder(nil, nil, musicTableView.view(musicListView.SwapTable), lyrics.view(), musicListView.view())
	bottomContainer := container.NewBorder(widget.NewSeparator(), nil, nil, nil, controller.View(musicPlayer))
	view := container.NewBorder(nil, bottomContainer, nil, nil, topContainer)
	
	window.SetContent(view)
	window.Show()
}
