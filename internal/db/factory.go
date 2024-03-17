package db

// Description: This file contains the factory functions for creating new instances of gorm.DB.
// It also contains the factory functions for creating new instances of the cache.
import (
	"k8s.io/klog"
	"os"
	"path/filepath"
	"sync"
	
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CacheInterface interface {
	Load(key string) (interface{}, bool)
	Store(key string, value interface{}) bool
	Delete(key string)
}

// Factory is a function that returns a new instance of gorm.DB
type Factory func() *gorm.DB

// CacheFactory is a function that returns a new instance of cacheInterface
type CacheFactory func() CacheInterface

// SqliteFactory is a factory function that returns a new instance of gorm.DB
func SqliteFactory(path string) Factory {
	// If the path is empty, use the path of the current executable
	if path == "" {
		binaryPath, _ := os.Executable()
		binaryDir := filepath.Dir(binaryPath)
		path = filepath.Join(binaryDir, "storage.db")
	} else {
		path = filepath.Join(path, "storage.db")
	}
	klog.Info(path)
	return func() *gorm.DB {
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return db
	}
}

func MemoryCacheFactory() CacheFactory {
	return func() CacheInterface {
		return &memoryCache{}
	}
}

type memoryCache struct{ sync.Map }

func (m *memoryCache) Load(key string) (value any, ok bool) {
	return m.Map.Load(key)
}

func (m *memoryCache) Store(key string, value interface{}) bool {
	m.Map.Store(key, value)
	return true
}

func (m *memoryCache) Delete(key string) {
	m.Map.Delete(key)
}
