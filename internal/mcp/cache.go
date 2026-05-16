package mcp

import (
	"sync"
	"time"

	"github.com/Swif7ify/Obelisk-CLI/internal/ai"
	"github.com/Swif7ify/Obelisk-CLI/internal/engine"
)

// ResultCache stores the latest scan results for resource access
type ResultCache struct {
	mu          sync.RWMutex
	lastResult  *engine.Result
	lastScanAt  time.Time
	projectPath string
}

// NewResultCache creates a new result cache
func NewResultCache() *ResultCache {
	return &ResultCache{}
}

// Set stores a scan result in the cache
func (c *ResultCache) Set(result *engine.Result, projectPath string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.lastResult = result
	c.lastScanAt = time.Now()
	c.projectPath = projectPath
}

// Get retrieves the latest scan result
func (c *ResultCache) Get() (*engine.Result, time.Time, string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if c.lastResult == nil {
		return nil, time.Time{}, "", false
	}
	
	return c.lastResult, c.lastScanAt, c.projectPath, true
}

// GetHealthScore retrieves just the health score
func (c *ResultCache) GetHealthScore() (*ai.HealthReport, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	if c.lastResult == nil || c.lastResult.Report == nil {
		return nil, false
	}
	
	return c.lastResult.Report, true
}

// Clear removes all cached data
func (c *ResultCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.lastResult = nil
	c.lastScanAt = time.Time{}
	c.projectPath = ""
}

// HasData checks if cache has data
func (c *ResultCache) HasData() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	return c.lastResult != nil
}

// Made with Bob
