package model

import (
	"fyne.io/fyne/v2/data/binding"
	"gorm.io/gorm"
	"math"
	"strings"
)

type MusicType int

const (
	MusicTypeMP3 = iota
	MusicTypeFLAC
	MusicTypeWAV
	MusicTypeEnd
	
	MusicTypeUnknown = math.MinInt16
)

func InitModel(db *gorm.DB) {
	err := db.AutoMigrate(&MusicTable{}, &Music{})
	if err != nil {
		panic(err)
	}
}

func IsMusicType(name string) (MusicType, bool) {
	if strings.HasSuffix(name, "mp3") {
		return MusicTypeMP3, true
	}
	if strings.HasSuffix(name, "flac") {
		return MusicTypeFLAC, true
	}
	if strings.HasSuffix(name, "wav") {
		return MusicTypeWAV, true
	}
	
	return MusicTypeUnknown, false
}

type MusicTable struct {
	Name string `gorm:"name"`
	gorm.Model
}

// 该方法用于实现binding.DataItem接口
func (t MusicTable) AddListener(listener binding.DataListener) {

}
func (t MusicTable) RemoveListener(listener binding.DataListener) {

}

type Music struct {
	MusicTableID uint      `gorm:"music_table_id,index"`
	Name         string    `gorm:"name,unique"`
	Singer       string    `gorm:"writer"`
	Album        string    `gorm:"Album"`
	Length       string    `gorm:"length"`
	Pic          string    `gorm:"pic"`
	Path         string    `gorm:"path"`
	Type         MusicType `gorm:"type"`
	gorm.Model
}

// 该方法用于实现binding.DataItem接口
func (t Music) AddListener(listener binding.DataListener) {

}
func (t Music) RemoveListener(listener binding.DataListener) {

}
