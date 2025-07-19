package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
)

// CacheManager handles persistent storage of parsed data
type CacheManager struct {
	db           *badger.DB
	libraryPath  string
}

// CachedPage represents a page with metadata for caching
type CachedPage struct {
	Page         json.RawMessage   `json:"page"`
	FileModTime  time.Time         `json:"file_mod_time"`
	Dependencies []string          `json:"dependencies"` // Pages this page links to
}

// CacheMetadata stores overall cache information
type CacheMetadata struct {
	Version      string    `json:"version"`
	LastUpdated  time.Time `json:"last_updated"`
	LibraryPath  string    `json:"library_path"`
}

const (
	cacheVersion = "1.0"
	metadataKey  = "cache_metadata"
	pagePrefix   = "page:"
	backlinksPrefix = "backlinks:"
)

// NewCacheManager creates a new cache manager
func NewCacheManager(libraryPath string) (*CacheManager, error) {
	// Determine cache directory based on library path
	cacheDir, err := getCacheDir(libraryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	// Open BadgerDB
	opts := badger.DefaultOptions(cacheDir)
	opts.Logger = nil // Disable verbose logging
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %w", err)
	}

	return &CacheManager{
		db:          db,
		libraryPath: libraryPath,
	}, nil
}

// Close closes the cache database
func (cm *CacheManager) Close() error {
	return cm.db.Close()
}

// SavePage saves a parsed page to the cache
func (cm *CacheManager) SavePage(page interface{}, pageName string, filePath string, dependencies []string) error {
	// Get file modification time
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	// Marshal the page separately
	pageData, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("failed to marshal page: %w", err)
	}

	cached := CachedPage{
		Page:         pageData,
		FileModTime:  fileInfo.ModTime(),
		Dependencies: dependencies,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cached page: %w", err)
	}

	err = cm.db.Update(func(txn *badger.Txn) error {
		key := []byte(pagePrefix + pageName)
		return txn.Set(key, data)
	})

	if err != nil {
		return fmt.Errorf("failed to save page to cache: %w", err)
	}

	return nil
}

// GetPage retrieves a cached page if it's still valid
func (cm *CacheManager) GetPage(pageName string, filePath string) (interface{}, bool, error) {
	// Check if file has been modified
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, false, nil // File doesn't exist, return cache miss
	}

	var cached CachedPage
	err = cm.db.View(func(txn *badger.Txn) error {
		key := []byte(pagePrefix + pageName)
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &cached)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, false, nil // Cache miss
	}

	if err != nil {
		return nil, false, fmt.Errorf("failed to read from cache: %w", err)
	}

	// Check if cached data is still valid
	if fileInfo.ModTime().After(cached.FileModTime) {
		return nil, false, nil // File has been modified, cache miss
	}

	// Return the raw JSON for the caller to unmarshal
	return cached.Page, true, nil
}

// SaveBacklinks saves the backlinks for a page
func (cm *CacheManager) SaveBacklinks(pageName string, backlinks interface{}) error {
	data, err := json.Marshal(backlinks)
	if err != nil {
		return fmt.Errorf("failed to marshal backlinks: %w", err)
	}

	err = cm.db.Update(func(txn *badger.Txn) error {
		key := []byte(backlinksPrefix + pageName)
		return txn.Set(key, data)
	})

	if err != nil {
		return fmt.Errorf("failed to save backlinks to cache: %w", err)
	}

	return nil
}

// GetBacklinks retrieves cached backlinks for a page
func (cm *CacheManager) GetBacklinks(pageName string) (interface{}, bool, error) {
	var backlinks interface{}

	err := cm.db.View(func(txn *badger.Txn) error {
		key := []byte(backlinksPrefix + pageName)
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &backlinks)
		})
	})

	if err == badger.ErrKeyNotFound {
		return nil, false, nil // Cache miss
	}

	if err != nil {
		return nil, false, fmt.Errorf("failed to read backlinks from cache: %w", err)
	}

	return backlinks, true, nil
}

// SaveMetadata saves cache metadata
func (cm *CacheManager) SaveMetadata() error {
	metadata := CacheMetadata{
		Version:     cacheVersion,
		LastUpdated: time.Now(),
		LibraryPath: cm.libraryPath,
	}

	data, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	err = cm.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(metadataKey), data)
	})

	if err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}

	return nil
}

// ValidateCache checks if the cache is valid for the current library
func (cm *CacheManager) ValidateCache() (bool, error) {
	var metadata CacheMetadata

	err := cm.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(metadataKey))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &metadata)
		})
	})

	if err == badger.ErrKeyNotFound {
		return false, nil // No metadata, invalid cache
	}

	if err != nil {
		return false, fmt.Errorf("failed to read metadata: %w", err)
	}

	// Check if it's the same library and version
	if metadata.LibraryPath != cm.libraryPath || metadata.Version != cacheVersion {
		return false, nil
	}

	return true, nil
}

// Clear removes all cached data
func (cm *CacheManager) Clear() error {
	return cm.db.DropAll()
}

// getCacheDir returns the appropriate cache directory for the platform
func getCacheDir(libraryPath string) (string, error) {
	var cacheDir string
	
	// If running tests, use a temp directory within the repo
	if os.Getenv("SEQ2B_TEST_MODE") == "true" || libraryPath == "" {
		// Find the seq2b root directory by looking for go.mod
		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		
		rootDir := currentDir
		for {
			if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
				// Found the root
				cacheDir = filepath.Join(rootDir, "tmp", "cache")
				break
			}
			parent := filepath.Dir(rootDir)
			if parent == rootDir {
				// Reached filesystem root without finding go.mod
				// Fall back to system temp
				cacheDir = filepath.Join(os.TempDir(), "seq2b-test-cache")
				break
			}
			rootDir = parent
		}
	} else {
		// For normal operation, store cache in cache subdirectory of the library
		cacheDir = filepath.Join(libraryPath, "cache")
	}
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}