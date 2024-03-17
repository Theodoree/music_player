package model

import (
	"errors"
	"fmt"
	
	"gorm.io/gorm"
)

type cacheInterface interface {
	Load(key string) (interface{}, bool)
	Store(key string, value interface{}) bool
	Delete(key string)
}

var (
	NotFoundPrimaryKey = errors.New("primary key not found")
)

const (
	musicListCacheKey = "MUSIC_LIST_%d"
	musicItemCacheKey = "MUSIC_ITEM_%d"
)

func getList[T any](db *gorm.DB, page, limit uint) ([]T, error) {
	var result []T
	err := db.Offset(int(page * limit)).Limit(int(limit)).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
func save[T any](db *gorm.DB, items []T) error {
	return db.CreateInBatches(items, len(items)).Error
}
func update[T any](db *gorm.DB, items T) error {
	return db.Save(items).Error
}

type basicQuery[T any] struct {
	empty T
}

func (q basicQuery[T]) add(db *gorm.DB, item T) error {
	return q.addBatch(db, []T{item})
}
func (q basicQuery[T]) addBatch(db *gorm.DB, items []T) error {
	return save(db, items)
}
func (q basicQuery[T]) delete(db *gorm.DB, ID uint) error {
	if ID == 0 {
		return NotFoundPrimaryKey
	}
	return db.Delete(&q.empty, ID).Error
}
func (q basicQuery[T]) update(db *gorm.DB, u T, ID uint) error {
	if ID == 0 {
		return NotFoundPrimaryKey
	}
	return update(db, u)
}
func (q basicQuery[T]) GetList(db *gorm.DB, page uint, limit uint) ([]T, error) {
	return getList[T](db, page, limit)
}

type MusicTableQuery struct {
	basicQuery[MusicTable]
}

func (q MusicTableQuery) Add(db *gorm.DB, item MusicTable) error {
	return q.basicQuery.add(db, item)
}
func (q MusicTableQuery) AddBatch(db *gorm.DB, items []MusicTable) error {
	return q.basicQuery.addBatch(db, items)
}
func (q MusicTableQuery) CacheKey(id uint) string {
	return fmt.Sprintf(musicListCacheKey, id)
}
func (q MusicTableQuery) Delete(db *gorm.DB, cacheService cacheInterface, item MusicTable) error {
	cacheService.Delete(q.CacheKey(item.ID))
	return q.basicQuery.delete(db, item.ID)
}
func (q MusicTableQuery) Update(db *gorm.DB, cacheService cacheInterface, item MusicTable) error {
	cacheService.Delete(q.CacheKey(item.ID))
	return q.basicQuery.update(db, item, item.ID)
	
}
func (q MusicTableQuery) GetByID(db *gorm.DB, cacheService cacheInterface, id uint) (MusicTable, error) {
	value, _ := cacheService.Load(q.CacheKey(id))
	if value != nil {
		return value.(MusicTable), nil
	}
	var m MusicTable
	err := db.First(&m, id).Error
	if err != nil {
		return m, err
	}
	cacheService.Store(q.CacheKey(id), m)
	return m, nil
}

type MusicQuery struct {
	basicQuery[Music]
}

func (q MusicQuery) Add(db *gorm.DB, item Music) error {
	item.SetUnion()
	return q.basicQuery.add(db, item)
}
func (q MusicQuery) AddBatch(db *gorm.DB, items []Music) error {
	for i := 0; i < len(items); i++ {
		items[i].SetUnion()
	}
	return q.basicQuery.addBatch(db, items)
}
func (q MusicQuery) CacheKey(id uint) string {
	return fmt.Sprintf(musicItemCacheKey, id)
}
func (q MusicQuery) Delete(db *gorm.DB, cacheService cacheInterface, item Music) error {
	cacheService.Delete(q.CacheKey(item.ID))
	return q.basicQuery.delete(db, item.ID)
}
func (q MusicQuery) GetByID(db *gorm.DB, cacheService cacheInterface, id uint) (Music, error) {
	value, _ := cacheService.Load(q.CacheKey(id))
	if value != nil {
		return value.(Music), nil
	}
	var m Music
	err := db.First(&m, id).Error
	if err != nil {
		return m, err
	}
	cacheService.Store(q.CacheKey(id), m)
	return m, nil
}
func (q MusicQuery) DeleteByMusicListID(db *gorm.DB, musicTableTag uint) error {
	return db.Delete(&q.empty, "music_table_id = ?", musicTableTag).Error
}
func (q MusicQuery) GetByMusicListID(db *gorm.DB, musicTableTag uint) ([]Music, error) {
	var items []Music
	return items, db.Find(&items, "music_table_id = ?", musicTableTag).Error
}
func (q MusicQuery) Update(db *gorm.DB, item Music) error {
	return q.basicQuery.update(db, item, item.ID)
}
