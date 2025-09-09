// Package optimization provides query optimization features including caching,
// connection pooling optimization, prepared statement management, and query analysis.
package optimization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// QueryCache implements a thread-safe in-memory query cache
type QueryCache struct {
	cache map[string]*CacheEntry
	mutex sync.RWMutex
	ttl   time.Duration
}

type CacheEntry struct {
	Data      types.Collection
	Count     int64
	CreatedAt time.Time
	Hits      int64
}

// NewQueryCache creates a new query cache with specified TTL
func NewQueryCache(ttl time.Duration) *QueryCache {
	cache := &QueryCache{
		cache: make(map[string]*CacheEntry),
		ttl:   ttl,
	}
	
	// Start cleanup goroutine
	go cache.cleanupExpired()
	
	return cache
}

// Get retrieves cached query result
func (qc *QueryCache) Get(key string) (types.Collection, int64, bool) {
	qc.mutex.RLock()
	defer qc.mutex.RUnlock()
	
	entry, exists := qc.cache[key]
	if !exists {
		return nil, 0, false
	}
	
	// Check if expired
	if time.Since(entry.CreatedAt) > qc.ttl {
		return nil, 0, false
	}
	
	entry.Hits++
	return entry.Data, entry.Count, true
}

// Set stores query result in cache
func (qc *QueryCache) Set(key string, data types.Collection, count int64) {
	qc.mutex.Lock()
	defer qc.mutex.Unlock()
	
	qc.cache[key] = &CacheEntry{
		Data:      data,
		Count:     count,
		CreatedAt: time.Now(),
		Hits:      0,
	}
}

// Clear removes all cached entries
func (qc *QueryCache) Clear() {
	qc.mutex.Lock()
	defer qc.mutex.Unlock()
	
	qc.cache = make(map[string]*CacheEntry)
}

// Stats returns cache statistics
func (qc *QueryCache) Stats() CacheStats {
	qc.mutex.RLock()
	defer qc.mutex.RUnlock()
	
	var totalHits int64
	expired := 0
	
	for _, entry := range qc.cache {
		totalHits += entry.Hits
		if time.Since(entry.CreatedAt) > qc.ttl {
			expired++
		}
	}
	
	return CacheStats{
		TotalEntries: len(qc.cache),
		TotalHits:    totalHits,
		ExpiredCount: expired,
	}
}

// cleanupExpired removes expired cache entries
func (qc *QueryCache) cleanupExpired() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		qc.mutex.Lock()
		for key, entry := range qc.cache {
			if time.Since(entry.CreatedAt) > qc.ttl {
				delete(qc.cache, key)
			}
		}
		qc.mutex.Unlock()
	}
}

type CacheStats struct {
	TotalEntries int   `json:"total_entries"`
	TotalHits    int64 `json:"total_hits"`
	ExpiredCount int   `json:"expired_count"`
}

// QueryOptimizer provides query optimization capabilities
type QueryOptimizer struct {
	cache          *QueryCache
	config         types.QueryOptimization
	preparedStmts  map[string]*PreparedStatement
	stmtMutex     sync.RWMutex
	queryLog      []QueryLogEntry
	logMutex      sync.RWMutex
}

type PreparedStatement struct {
	SQL      string
	Hash     string
	UsageCount int64
	CreatedAt time.Time
}

type QueryLogEntry struct {
	SQL       string    `json:"sql"`
	Bindings  []any     `json:"bindings"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
}

// NewQueryOptimizer creates a new query optimizer
func NewQueryOptimizer(config types.QueryOptimization) *QueryOptimizer {
	var cache *QueryCache
	if config.EnableQueryCache {
		cache = NewQueryCache(config.CacheTTL)
	}
	
	return &QueryOptimizer{
		cache:         cache,
		config:        config,
		preparedStmts: make(map[string]*PreparedStatement),
		queryLog:      make([]QueryLogEntry, 0),
	}
}

// GenerateCacheKey creates a cache key from SQL and bindings
func (qo *QueryOptimizer) GenerateCacheKey(sql string, bindings []any) string {
	hasher := sha256.New()
	hasher.Write([]byte(sql))
	for _, binding := range bindings {
		_, _ = fmt.Fprintf(hasher, "%v", binding)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetCachedResult retrieves cached query result
func (qo *QueryOptimizer) GetCachedResult(key string) (types.Collection, int64, bool) {
	if qo.cache == nil {
		return nil, 0, false
	}
	return qo.cache.Get(key)
}

// CacheResult stores query result in cache
func (qo *QueryOptimizer) CacheResult(key string, data types.Collection, count int64) {
	if qo.cache != nil {
		qo.cache.Set(key, data, count)
	}
}

// RegisterPreparedStatement tracks prepared statement usage
func (qo *QueryOptimizer) RegisterPreparedStatement(sql string) string {
	if !qo.config.EnablePreparedStmt {
		return ""
	}
	
	hash := qo.generateSQLHash(sql)
	
	qo.stmtMutex.Lock()
	defer qo.stmtMutex.Unlock()
	
	if stmt, exists := qo.preparedStmts[hash]; exists {
		stmt.UsageCount++
		return hash
	}
	
	qo.preparedStmts[hash] = &PreparedStatement{
		SQL:        sql,
		Hash:       hash,
		UsageCount: 1,
		CreatedAt:  time.Now(),
	}
	
	return hash
}

// LogQuery records query execution for analysis
func (qo *QueryOptimizer) LogQuery(sql string, bindings []any, duration time.Duration, err error) {
	if !qo.config.EnableQueryLog {
		return
	}
	
	qo.logMutex.Lock()
	defer qo.logMutex.Unlock()
	
	entry := QueryLogEntry{
		SQL:       sql,
		Bindings:  bindings,
		Duration:  duration,
		Timestamp: time.Now(),
	}
	
	if err != nil {
		entry.Error = err.Error()
	}
	
	qo.queryLog = append(qo.queryLog, entry)
	
	// Keep only last 1000 entries
	if len(qo.queryLog) > 1000 {
		qo.queryLog = qo.queryLog[len(qo.queryLog)-1000:]
	}
}

// GetQueryStats returns query performance statistics
func (qo *QueryOptimizer) GetQueryStats() QueryStats {
	qo.logMutex.RLock()
	defer qo.logMutex.RUnlock()
	
	if len(qo.queryLog) == 0 {
		return QueryStats{}
	}
	
	var totalDuration time.Duration
	slowQueries := 0
	errorCount := 0
	
	for _, entry := range qo.queryLog {
		totalDuration += entry.Duration
		if entry.Duration > 1*time.Second {
			slowQueries++
		}
		if entry.Error != "" {
			errorCount++
		}
	}
	
	avgDuration := totalDuration / time.Duration(len(qo.queryLog))
	
	stats := QueryStats{
		TotalQueries:    len(qo.queryLog),
		AverageDuration: avgDuration,
		SlowQueries:     slowQueries,
		ErrorCount:      errorCount,
	}
	
	if qo.cache != nil {
		cacheStats := qo.cache.Stats()
		stats.CacheStats = &cacheStats
	}
	
	return stats
}

// GetSlowQueries returns queries that took longer than threshold
func (qo *QueryOptimizer) GetSlowQueries(threshold time.Duration) []QueryLogEntry {
	qo.logMutex.RLock()
	defer qo.logMutex.RUnlock()
	
	var slowQueries []QueryLogEntry
	for _, entry := range qo.queryLog {
		if entry.Duration > threshold {
			slowQueries = append(slowQueries, entry)
		}
	}
	
	return slowQueries
}

// ClearStats clears all optimization statistics
func (qo *QueryOptimizer) ClearStats() {
	qo.logMutex.Lock()
	qo.queryLog = make([]QueryLogEntry, 0)
	qo.logMutex.Unlock()
	
	if qo.cache != nil {
		qo.cache.Clear()
	}
	
	qo.stmtMutex.Lock()
	qo.preparedStmts = make(map[string]*PreparedStatement)
	qo.stmtMutex.Unlock()
}

func (qo *QueryOptimizer) generateSQLHash(sql string) string {
	hasher := sha256.New()
	hasher.Write([]byte(sql))
	return hex.EncodeToString(hasher.Sum(nil))
}

type QueryStats struct {
	TotalQueries    int            `json:"total_queries"`
	AverageDuration time.Duration  `json:"average_duration"`
	SlowQueries     int            `json:"slow_queries"`
	ErrorCount      int            `json:"error_count"`
	CacheStats      *CacheStats    `json:"cache_stats,omitempty"`
}

// ConcurrencyManager manages concurrent query execution
type ConcurrencyManager struct {
	semaphore chan struct{}
}

// NewConcurrencyManager creates a concurrency manager with max concurrent queries
func NewConcurrencyManager(maxConcurrency int) *ConcurrencyManager {
	return &ConcurrencyManager{
		semaphore: make(chan struct{}, maxConcurrency),
	}
}

// Acquire acquires a slot for concurrent execution
func (cm *ConcurrencyManager) Acquire(ctx context.Context) error {
	select {
	case cm.semaphore <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release releases a concurrency slot
func (cm *ConcurrencyManager) Release() {
	<-cm.semaphore
}

// ExecuteWithConcurrencyLimit executes function with concurrency control
func (cm *ConcurrencyManager) ExecuteWithConcurrencyLimit(ctx context.Context, fn func() error) error {
	if err := cm.Acquire(ctx); err != nil {
		return err
	}
	defer cm.Release()
	
	return fn()
}