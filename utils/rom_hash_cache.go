package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

type RomHashCache struct {
	Cache       map[string]map[string]int `json:"cache"` // slug -> sha1 -> romID
	LastUpdated time.Time                 `json:"last_updated"`
}

func getCacheFilePath() string {
	return filepath.Join(".cache", "rom_hash_cache.json")
}

func loadRomHashCache() (*RomHashCache, error) {
	logger := gaba.GetLogger()
	cachePath := getCacheFilePath()

	// If file doesn't exist, return empty cache (not an error)
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		logger.Debug("No cache file found, starting with empty cache")
		return &RomHashCache{
			Cache:       make(map[string]map[string]int),
			LastUpdated: time.Now(),
		}, nil
	}

	// Read and unmarshal
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var cache RomHashCache
	if err := json.Unmarshal(data, &cache); err != nil {
		// If corrupted, return empty cache and log warning
		logger.Warn("Cache file corrupted, starting fresh", "error", err)
		return &RomHashCache{
			Cache:       make(map[string]map[string]int),
			LastUpdated: time.Now(),
		}, nil
	}

	logger.Debug("Loaded cache", "platforms", len(cache.Cache), "lastUpdated", cache.LastUpdated)
	return &cache, nil
}

func (c *RomHashCache) save() error {
	logger := gaba.GetLogger()
	c.LastUpdated = time.Now()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	cachePath := getCacheFilePath()
	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	logger.Debug("Saved cache", "path", cachePath)
	return nil
}

func (c *RomHashCache) lookup(platformSlug, sha1 string) (int, bool) {
	if platformCache, ok := c.Cache[platformSlug]; ok {
		if romID, found := platformCache[sha1]; found {
			return romID, true
		}
	}
	return 0, false
}

func (c *RomHashCache) set(platformSlug, sha1 string, romID int) {
	if c.Cache[platformSlug] == nil {
		c.Cache[platformSlug] = make(map[string]int)
	}
	c.Cache[platformSlug][sha1] = romID
}
