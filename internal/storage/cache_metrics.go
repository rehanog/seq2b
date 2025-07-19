package storage

import (
	"sync"
	"sync/atomic"
	"time"
)

// CacheMetrics tracks cache performance and usage
type CacheMetrics struct {
	// Counters (use atomic operations)
	Hits         atomic.Int64
	Misses       atomic.Int64
	Saves        atomic.Int64
	SaveErrors   atomic.Int64
	Gets         atomic.Int64
	GetErrors    atomic.Int64
	Evictions    atomic.Int64
	
	// Sizes
	TotalBytes   atomic.Int64
	EntryCount   atomic.Int64
	
	// Timing
	mu            sync.RWMutex
	SaveDurations []time.Duration
	GetDurations  []time.Duration
	
	// Start time for rate calculations
	StartTime    time.Time
}

// NewCacheMetrics creates a new metrics instance
func NewCacheMetrics() *CacheMetrics {
	return &CacheMetrics{
		StartTime:     time.Now(),
		SaveDurations: make([]time.Duration, 0, 1000),
		GetDurations:  make([]time.Duration, 0, 1000),
	}
}

// RecordSave records a save operation
func (m *CacheMetrics) RecordSave(duration time.Duration, sizeBytes int64, err error) {
	m.Saves.Add(1)
	if err != nil {
		m.SaveErrors.Add(1)
	} else {
		m.TotalBytes.Add(sizeBytes)
		m.EntryCount.Add(1)
	}
	
	m.mu.Lock()
	if len(m.SaveDurations) < cap(m.SaveDurations) {
		m.SaveDurations = append(m.SaveDurations, duration)
	}
	m.mu.Unlock()
}

// RecordGet records a get operation
func (m *CacheMetrics) RecordGet(duration time.Duration, hit bool, err error) {
	m.Gets.Add(1)
	if err != nil {
		m.GetErrors.Add(1)
	} else if hit {
		m.Hits.Add(1)
	} else {
		m.Misses.Add(1)
	}
	
	m.mu.Lock()
	if len(m.GetDurations) < cap(m.GetDurations) {
		m.GetDurations = append(m.GetDurations, duration)
	}
	m.mu.Unlock()
}

// GetStats returns current statistics
func (m *CacheMetrics) GetStats() CacheStats {
	elapsed := time.Since(m.StartTime).Seconds()
	if elapsed == 0 {
		elapsed = 1 // Prevent division by zero
	}
	
	hits := m.Hits.Load()
	misses := m.Misses.Load()
	total := hits + misses
	
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}
	
	m.mu.RLock()
	avgSaveTime := calculateAverage(m.SaveDurations)
	avgGetTime := calculateAverage(m.GetDurations)
	m.mu.RUnlock()
	
	return CacheStats{
		Hits:          hits,
		Misses:        misses,
		HitRate:       hitRate,
		Saves:         m.Saves.Load(),
		SaveErrors:    m.SaveErrors.Load(),
		Gets:          m.Gets.Load(),
		GetErrors:     m.GetErrors.Load(),
		TotalBytes:    m.TotalBytes.Load(),
		EntryCount:    m.EntryCount.Load(),
		SavesPerSec:   float64(m.Saves.Load()) / elapsed,
		GetsPerSec:    float64(m.Gets.Load()) / elapsed,
		AvgSaveTime:   avgSaveTime,
		AvgGetTime:    avgGetTime,
		Uptime:        time.Since(m.StartTime),
	}
}

// Reset resets all metrics
func (m *CacheMetrics) Reset() {
	m.Hits.Store(0)
	m.Misses.Store(0)
	m.Saves.Store(0)
	m.SaveErrors.Store(0)
	m.Gets.Store(0)
	m.GetErrors.Store(0)
	m.Evictions.Store(0)
	m.TotalBytes.Store(0)
	m.EntryCount.Store(0)
	
	m.mu.Lock()
	m.SaveDurations = m.SaveDurations[:0]
	m.GetDurations = m.GetDurations[:0]
	m.mu.Unlock()
	
	m.StartTime = time.Now()
}

// CacheStats represents a snapshot of cache statistics
type CacheStats struct {
	Hits          int64
	Misses        int64
	HitRate       float64
	Saves         int64
	SaveErrors    int64
	Gets          int64
	GetErrors     int64
	TotalBytes    int64
	EntryCount    int64
	SavesPerSec   float64
	GetsPerSec    float64
	AvgSaveTime   time.Duration
	AvgGetTime    time.Duration
	Uptime        time.Duration
}

// calculateAverage calculates average duration
func calculateAverage(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	
	return total / time.Duration(len(durations))
}

// MetricsCacheManager wraps CacheManager with metrics collection
type MetricsCacheManager struct {
	*CacheManager
	metrics *CacheMetrics
}

// NewMetricsCacheManager creates a cache manager with metrics
func NewMetricsCacheManager(libraryPath string) (*MetricsCacheManager, error) {
	cache, err := NewCacheManager(libraryPath)
	if err != nil {
		return nil, err
	}
	
	return &MetricsCacheManager{
		CacheManager: cache,
		metrics:      NewCacheMetrics(),
	}, nil
}

// SavePage saves a page and records metrics
func (m *MetricsCacheManager) SavePage(page interface{}, pageName string, filePath string, dependencies []string) error {
	start := time.Now()
	
	// Estimate size (this is approximate)
	// In production, you'd want more accurate size tracking
	estimatedSize := int64(len(pageName) + len(filePath) + len(dependencies)*20)
	
	err := m.CacheManager.SavePage(page, pageName, filePath, dependencies)
	
	duration := time.Since(start)
	m.metrics.RecordSave(duration, estimatedSize, err)
	
	return err
}

// GetPage retrieves a page and records metrics
func (m *MetricsCacheManager) GetPage(pageName string, filePath string) (interface{}, bool, error) {
	start := time.Now()
	
	page, hit, err := m.CacheManager.GetPage(pageName, filePath)
	
	duration := time.Since(start)
	m.metrics.RecordGet(duration, hit, err)
	
	return page, hit, err
}

// GetMetrics returns current metrics
func (m *MetricsCacheManager) GetMetrics() CacheStats {
	return m.metrics.GetStats()
}

// ResetMetrics resets all metrics
func (m *MetricsCacheManager) ResetMetrics() {
	m.metrics.Reset()
}