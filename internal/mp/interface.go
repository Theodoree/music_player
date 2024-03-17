package mp

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/Theodoree/music_player/internal/model"
	"github.com/Theodoree/music_player/internal/music"
)

type MusicPlayerPlaybackModule interface {
	// Play 播放音乐
	Play()
	// Pause 暂停音乐
	Pause()
	// Stop 停止音乐
	Stop()
	// Prev 上一首
	Prev()
	// Next 下一首
	Next()
}

type MusicPlayerDataModule interface {
	// Lyrics 返回一个动态绑定的歌词(全量)
	Lyrics() binding.String
	// Lyric 返回一个动态绑定的歌词(一行)
	Lyric() binding.String
	// MusicName 返回一个动态绑定的音乐名
	MusicName() binding.String
	// SingerName 返回一个动态绑定的歌手名
	SingerName() binding.String
	// Volume 返回一个动态绑定的float64[0,100]
	Volume() binding.Float
	// ProcessBar 返回一个动态绑定的进度条值[0,100]
	ProcessBar() binding.Float
	// MusicCurTime 返回一个动态绑定的当前播放时间
	MusicCurTime() binding.String
	// MusicEndTime 返回一个动态绑定的音乐结束时间
	MusicEndTime() binding.String
	// MusicTableList 自定义表格列表和索引
	MusicTableList() (binding.DataList, binding.Int)
	// MusicList 音乐列表和索引
	MusicList() (binding.DataList, binding.Int, binding.String)
	// PictureList 返回一个墙纸列表
	PictureList() binding.DataList
	// PlayMode 播放模式
	PlayMode() binding.DataItem
	// PlayStatus 返回一个动态绑定的播放状态
	PlayStatus() binding.Bool
	
	// StreamMusicList 流媒体列表和索引
	StreamMusicList() (binding.DataList, binding.Int, binding.String)
	// StreamMusicTableList 流媒体表格列表和索引
	StreamMusicTableList() (binding.DataList, binding.Int)
}

type MusicPlayerOperationModule interface {
	// AddTable 新增表格
	AddTable(table model.MusicTable)
	// DelTable  删除表格
	DelTable(tableID uint)
	// ImportMusic 导入音乐
	ImportMusic(tableID uint, path string)
	// AddWallpaper 新增墙纸
	AddWallpaper(path string)
	// AddMusic 新增音乐项至指定表格
	AddMusic(tableID uint, music model.Music)
	// UpdateMusic 更新音乐
	UpdateMusic(tableID uint, music model.Music)
	// GetPlayedMusic 获取当前音乐
	GetPlayedMusic() music.Music
}

type MusicPlayer interface {
	MusicPlayerPlaybackModule
	MusicPlayerDataModule
	MusicPlayerOperationModule
}
