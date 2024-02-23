package internal

import (
	"gorm.io/gorm"
)

// Factory is a function that returns a new instance of gorm.DB
type Factory func() *gorm.DB

// CacheFactory is a function that returns a new instance of cacheInterface
type CacheFactory func() CacheInterface

type CacheInterface interface {
	Load(key string) (interface{}, bool)
	Store(key string, value interface{}) bool
	Delete(key string)
}

type MusicPlayerOperation interface {
	Play() bool
	Pause() bool
	Prev()
	Next()
}
