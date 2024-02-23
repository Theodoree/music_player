package db

// Description: This file contains the factory functions for creating new instances of gorm.DB.
// It also contains the factory functions for creating new instances of the cache.
import (
	"os"
	"path/filepath"
	"sync"
	
	"github.com/Theodoree/music_player/internal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SqliteFactory is a factory function that returns a new instance of gorm.DB
func SqliteFactory(path string) internal.Factory {
	// If the path is empty, use the path of the current executable
	if path == "" {
		binaryPath, _ := os.Executable()
		binaryDir := filepath.Dir(binaryPath)
		path = filepath.Join(binaryDir, "storage.db")
	}
	return func() *gorm.DB {
		db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return db
	}
}

func MemoryCacheFactory() internal.CacheFactory {
	return func() internal.CacheInterface {
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
