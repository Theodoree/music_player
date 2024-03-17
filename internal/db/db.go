package db

import (
	"gorm.io/gorm"
	"k8s.io/klog"
	
	"github.com/Theodoree/music_player/internal/model"
)

type tableOperator interface {
	GetMusicTable(page uint, limit uint) ([]model.MusicTable, error)
	GetMusicTableByID(id uint) (model.MusicTable, error)
	SaveMusicTable(tag model.MusicTable) error
	DeleteMusicTable(item model.MusicTable) error
}
type musicOperator interface {
	GetMusicByMusicTableID(muscleListID uint) ([]model.Music, error)
	DeleteMusicByMusicTableID(muscleListID uint) error
	SaveMusics(item []model.Music) error
	SaveMusic(item model.Music) error
	UpdateMusic(item model.Music) error
}

type MusicStore interface {
	tableOperator
	musicOperator
}

const DefaultTableID = 1

type db struct {
	*gorm.DB
	cache CacheInterface
}

func New(factory Factory, cache CacheFactory) MusicStore {
	var store db
	store.DB = factory()
	store.cache = cache()
	model.InitModel(store.DB)
	store.Init()
	return &store
}
func (db *db) Init() {
	_, err := db.GetMusicTableByID(DefaultTableID)
	if err != nil {
		klog.Error(err)
		_ = db.SaveMusicTable(model.MusicTable{
			Name: "本地列表",
		})
	}
	
}

// implementation tableOperator

func (db *db) GetMusicTable(page uint, limit uint) ([]model.MusicTable, error) {
	return model.MusicTableQuery{}.GetList(db.DB, page, limit)
}
func (db *db) GetMusicTableByID(id uint) (model.MusicTable, error) {
	return model.MusicTableQuery{}.GetByID(db.DB, db.cache, id)
}
func (db *db) SaveMusicTable(item model.MusicTable) error {
	return model.MusicTableQuery{}.Add(db.DB, item)
}
func (db *db) DeleteMusicTable(item model.MusicTable) error {
	if item.ID == DefaultTableID {
		return nil
	}
	return model.MusicTableQuery{}.Delete(db.DB, db.cache, item)
}

// implementation musicOperator

func (db *db) GetMusicByMusicTableID(musicTableID uint) ([]model.Music, error) {
	return model.MusicQuery{}.GetByMusicListID(db.DB, musicTableID)
}
func (db *db) DeleteMusicByMusicTableID(musicTableID uint) error {
	return model.MusicQuery{}.DeleteByMusicListID(db.DB, musicTableID)
}
func (db *db) SaveMusics(item []model.Music) error {
	return model.MusicQuery{}.AddBatch(db.DB, item)
}
func (db *db) SaveMusic(item model.Music) error {
	return model.MusicQuery{}.Add(db.DB, item)
}
func (db *db) UpdateMusic(item model.Music) error {
	return model.MusicQuery{}.Update(db.DB, item)
}
