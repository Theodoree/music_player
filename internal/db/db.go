package db

import (
	"github.com/Theodoree/music_player/internal"
	"github.com/Theodoree/music_player/internal/model"
	"gorm.io/gorm"
	"k8s.io/klog"
)

type MusicStore interface {
	GetMusicTable(page uint, limit uint) ([]model.MusicTable, error)
	GetMusicTableByID(id uint) (model.MusicTable, error)
	SaveMusicTable(tag model.MusicTable) error
	
	GetMusicByMusicTableID(muscleListID uint) ([]model.Music, error)
	DeleteMusicByMusicTableID(muscleListID uint)
	SaveMusics(item []model.Music) error
	SaveMusic(item model.Music) error
}

const DefaultTableID = 1

type db struct {
	*gorm.DB
	cache internal.CacheInterface
}

func New(factory internal.Factory, cache internal.CacheFactory) MusicStore {
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
			Name: "默认表格",
		})
	}
	
}
func (db *db) GetMusicTable(page uint, limit uint) ([]model.MusicTable, error) {
	return model.MusicTableQuery{}.GetList(db.DB, page, limit)
}
func (db *db) GetMusicTableByID(id uint) (model.MusicTable, error) {
	return model.MusicTableQuery{}.GetByID(db.DB, db.cache, id)
}
func (db *db) SaveMusicTable(list model.MusicTable) error {
	return model.MusicTableQuery{}.Add(db.DB, list)
}
func (db *db) DeleteMusicTable(item model.MusicTable) error {
	if item.ID == DefaultTableID {
		return nil
	}
	return model.MusicTableQuery{}.Delete(db.DB, db.cache, item)
}

func (db *db) GetMusicByMusicTableID(musicTableTagID uint) ([]model.Music, error) {
	return model.MusicQuery{}.GetByMusicListID(db.DB, musicTableTagID)
}
func (db *db) DeleteMusicByMusicTableID(muscleListID uint) {
	model.MusicQuery{}.DeleteByMusicListID(db.DB, muscleListID)
}
func (db *db) SaveMusics(item []model.Music) error {
	return model.MusicQuery{}.AddBatch(db.DB, item)
}
func (db *db) SaveMusic(item model.Music) error {
	return model.MusicQuery{}.Add(db.DB, item)
}
