package model

import (
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
	
	"fyne.io/fyne/v2/data/binding"
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
	err := db.AutoMigrate(&MusicTable{}, &Music{}, &Picture{})
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
	dummyDataListener `gorm:"-"`
}

type Music struct {
	MusicTableID uint          `gorm:"music_table_id,index:music_table_id"`
	Name         string        `gorm:"name"`
	Singer       string        `gorm:"writer"`
	Album        string        `gorm:"Album"`
	Length       time.Duration `gorm:"length"`
	Path         string        `gorm:"path"`
	Type         MusicType     `gorm:"type"`
	Lyric        string        `gorm:"lyric"`
	Union        string        `gorm:"index:idx_name,unique"`
	
	dummyDataListener `gorm:"-"`
	gorm.Model
}

func (m *Music) SetUnion() {
	m.Union = fmt.Sprintf("%d-%s", m.MusicTableID, m.Path)
}

type Picture struct {
	Path              string `json:"path"`
	dummyDataListener `gorm:"-"`
	gorm.Model
}

type dummyDataListener struct{}

func (t dummyDataListener) AddListener(listener binding.DataListener) {

}
func (t dummyDataListener) RemoveListener(listener binding.DataListener) {}
