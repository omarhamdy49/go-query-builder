# Query Optimization & Performance

Supercharge your database queries with advanced optimization techniques, caching, and performance monitoring.

## 🚀 Overview

Go Query Builder includes comprehensive optimization features to ensure your applications run at peak performance. From intelligent caching to prepared statement management, every query can be optimized for speed and efficiency.

## 📋 Table of Contents

- [Query Caching](#query-caching)
- [Prepared Statements](#prepared-statements)
- [Connection Pooling](#connection-pooling)
- [Concurrency Management](#concurrency-management)
- [Query Analysis](#query-analysis)
- [Performance Monitoring](#performance-monitoring)
- [Best Practices](#best-practices)

## 💾 Query Caching

### Automatic Query Caching

```go
import (
    "time"
    "github.com/go-query-builder/querybuilder"
    "github.com/go-query-builder/querybuilder/pkg/types"
    "github.com/go-query-builder/querybuilder/pkg/optimization"
)

// Enable query caching globally
config := types.QueryOptimization{
    EnableQueryCache:   true,
    CacheTTL:          5 * time.Minute,
    EnablePreparedStmt: true,
    MaxConcurrency:    50,
    EnableQueryLog:    true,
}

optimizer := optimization.NewQueryOptimizer(config)

// Queries are automatically cached
users, err := querybuilder.QB().Table("users").
    Where("status", "active").
    Get(ctx) // First call - hits database

users2, err := querybuilder.QB().Table("users").
    Where("status", "active").  
    Get(ctx) // Second call - served from cache!
```

### Manual Cache Management

```go
// Create cache with custom TTL
cache := optimization.NewQueryCache(10 * time.Minute)

// Generate cache key
sql := "SELECT * FROM users WHERE status = ?"
bindings := []any{"active"}
cacheKey := optimizer.GenerateCacheKey(sql, bindings)

// Check cache
if data, count, found := cache.Get(cacheKey); found {
    fmt.Printf("Cache hit! %d records\n", data.Count())
} else {
    // Execute query and cache result
    users, err := querybuilder.QB().Table("users").Where("status", "active").Get(ctx)
    if err == nil {
        cache.Set(cacheKey, users, int64(users.Count()))
    }
}

// Cache statistics
stats := cache.Stats()
fmt.Printf("Cache entries: %d, Hits: %d\n", stats.TotalEntries, stats.TotalHits)
```

### Smart Cache Invalidation

```go
// Cache with automatic invalidation
type SmartCache struct {
    cache     *optimization.QueryCache
    tableTags map[string][]string // Track which tables each cache key uses
    mutex     sync.RWMutex
}

func (sc *SmartCache) InvalidateTable(tableName string) {
    sc.mutex.Lock()
    defer sc.mutex.Unlock()
    
    // Find and invalidate all queries that use this table
    for cacheKey, tables := range sc.tableTags {
        for _, table := range tables {
            if table == tableName {
                sc.cache.Delete(cacheKey)
                delete(sc.tableTags, cacheKey)
                break
            }
        }
    }
}

// Invalidate cache when data changes
func updateUser(ctx context.Context, userID int, data map[string]any) error {
    err := querybuilder.QB().Table("users").Where("id", userID).Update(ctx, data)
    if err == nil {
        smartCache.InvalidateTable("users") // Clear related cache entries
    }
    return err
}
```

## ⚡ Prepared Statements

### Automatic Prepared Statement Management

```go
// Prepared statements are managed automatically
optimizer := optimization.NewQueryOptimizer(types.QueryOptimization{
    EnablePreparedStmt: true,
})

// First execution - creates prepared statement
users1, _ := querybuilder.QB().Table("users").Where("id", 1).Get(ctx)

// Second execution - reuses prepared statement  
users2, _ := querybuilder.QB().Table("users").Where("id", 2).Get(ctx)

// Third execution - reuses same prepared statement
users3, _ := querybuilder.QB().Table("users").Where("id", 3).Get(ctx)
```

### Prepared Statement Statistics

```go
// Track prepared statement usage
sql := "SELECT * FROM users WHERE id = ?"
stmtHash := optimizer.RegisterPreparedStatement(sql)

fmt.Printf("Statement hash: %s\n", stmtHash)

// Get usage statistics
stats := optimizer.GetQueryStats()
fmt.Printf("Total queries: %d\n", stats.TotalQueries)
```

### Manual Prepared Statement Control

```go
type PreparedStatementManager struct {
    statements map[string]*sql.Stmt
    db         *sql.DB
    mutex      sync.RWMutex
}

func (psm *PreparedStatementManager) GetOrPrepare(sql string) (*sql.Stmt, error) {
    psm.mutex.RLock()
    if stmt, exists := psm.statements[sql]; exists {
        psm.mutex.RUnlock()
        return stmt, nil
    }
    psm.mutex.RUnlock()
    
    psm.mutex.Lock()
    defer psm.mutex.Unlock()
    
    // Double-check after acquiring write lock
    if stmt, exists := psm.statements[sql]; exists {
        return stmt, nil
    }
    
    stmt, err := psm.db.Prepare(sql)
    if err != nil {
        return nil, err
    }
    
    psm.statements[sql] = stmt
    return stmt, nil
}

func (psm *PreparedStatementManager) CloseAll() {
    psm.mutex.Lock()
    defer psm.mutex.Unlock()
    
    for _, stmt := range psm.statements {
        stmt.Close()
    }
    psm.statements = make(map[string]*sql.Stmt)
}
```

## 🏊 Connection Pooling

### Optimal Pool Configuration

```go
// Configure connection pool for performance
config := types.Config{
    Driver:          types.MySQL,
    Host:            "localhost",
    Port:            3306,
    Database:        "myapp",
    Username:        "user",
    Password:        "password",
    MaxOpenConns:    25,  // Maximum open connections
    MaxIdleConns:    5,   // Maximum idle connections
    ConnMaxLifetime: 5 * time.Minute,  // Connection lifetime
    ConnMaxIdleTime: 2 * time.Minute,  // Connection idle time
}

db, err := querybuilder.NewConnection(config)
if err != nil {
    log.Fatal(err)
}

// Monitor connection pool
stats := db.Stats()
fmt.Printf("Open connections: %d\n", stats.OpenConnections)
fmt.Printf("Idle connections: %d\n", stats.Idle) 
fmt.Printf("Connections in use: %d\n", stats.InUse)
```

### Dynamic Pool Scaling

```go
type DynamicPool struct {
    db          types.DB
    monitor     *PoolMonitor
    minConns    int
    maxConns    int
    scaleUpAt   float64 // Usage percentage to scale up
    scaleDownAt float64 // Usage percentage to scale down
}

func (dp *DynamicPool) AutoScale() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := dp.db.Stats()
        usage := float64(stats.InUse) / float64(stats.OpenConnections)
        
        if usage > dp.scaleUpAt && stats.OpenConnections < dp.maxConns {
            // Scale up
            newMax := min(stats.OpenConnections+5, dp.maxConns)
            dp.updateMaxConnections(newMax)
            log.Printf("Scaled up to %d connections", newMax)
        } else if usage < dp.scaleDownAt && stats.OpenConnections > dp.minConns {
            // Scale down
            newMax := max(stats.OpenConnections-2, dp.minConns)
            dp.updateMaxConnections(newMax)
            log.Printf("Scaled down to %d connections", newMax)
        }
    }
}
```

## 🚦 Concurrency Management

### Concurrency Limiting

```go
// Limit concurrent queries to prevent database overload
concurrencyManager := optimization.NewConcurrencyManager(10) // Max 10 concurrent queries

// Execute with concurrency control
err := concurrencyManager.ExecuteWithConcurrencyLimit(ctx, func() error {
    users, err := querybuilder.QB().Table("users").Get(ctx)
    if err != nil {
        return err
    }
    processUsers(users)
    return nil
})

if err != nil {
    log.Printf("Query failed: %v", err)
}
```

### Query Queue Management

```go
type QueryQueue struct {
    queue     chan QueryTask
    workers   int
    semaphore chan struct{}
}

type QueryTask struct {
    Query    func() error
    Priority int
    Timeout  time.Duration
}

func NewQueryQueue(workers, queueSize int) *QueryQueue {
    qq := &QueryQueue{
        queue:     make(chan QueryTask, queueSize),
        workers:   workers,
        semaphore: make(chan struct{}, workers),
    }
    
    // Start workers
    for i := 0; i < workers; i++ {
        go qq.worker()
    }
    
    return qq
}

func (qq *QueryQueue) worker() {
    for task := range qq.queue {
        qq.semaphore <- struct{}{} // Acquire
        
        go func(t QueryTask) {
            defer func() { <-qq.semaphore }() // Release
            
            ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
            defer cancel()
            
            done := make(chan error, 1)
            go func() {
                done <- t.Query()
            }()
            
            select {
            case err := <-done:
                if err != nil {
                    log.Printf("Query error: %v", err)
                }
            case <-ctx.Done():
                log.Printf("Query timeout")
            }
        }(task)
    }
}

func (qq *QueryQueue) Submit(task QueryTask) {
    select {
    case qq.queue <- task:
        // Task queued successfully
    default:
        log.Printf("Query queue full, dropping task")
    }
}
```

### Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    maxFailures   int
    resetTimeout  time.Duration
    failureCount  int
    lastFailTime  time.Time
    state         string // "closed", "open", "half-open"
    mutex         sync.RWMutex
}

func (cb *CircuitBreaker) Execute(query func() error) error {
    cb.mutex.RLock()
    state := cb.state
    cb.mutex.RUnlock()
    
    if state == "open" {
        if time.Since(cb.lastFailTime) > cb.resetTimeout {
            cb.setState("half-open")
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    }
    
    err := query()
    
    if err != nil {
        cb.recordFailure()
    } else {
        cb.recordSuccess()
    }
    
    return err
}

func (cb *CircuitBreaker) recordFailure() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    cb.failureCount++
    cb.lastFailTime = time.Now()
    
    if cb.failureCount >= cb.maxFailures {
        cb.state = "open"
    }
}

func (cb *CircuitBreaker) recordSuccess() {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    cb.failureCount = 0
    cb.state = "closed"
}
```

## 📈 Query Analysis

### Slow Query Detection

```go
// Detect and log slow queries
func (optimizer *QueryOptimizer) ExecuteWithLogging(ctx context.Context, query func() error, sql string, bindings []any) error {
    start := time.Now()
    err := query()
    duration := time.Since(start)
    
    // Log the query
    optimizer.LogQuery(sql, bindings, duration, err)
    
    // Alert on slow queries
    if duration > 2*time.Second {
        log.Printf("SLOW QUERY (%v): %s", duration, sql)
        
        // Send to monitoring system
        sendSlowQueryAlert(sql, duration, err)
    }
    
    return err
}

// Get slow queries for analysis
slowQueries := optimizer.GetSlowQueries(1 * time.Second)
for _, query := range slowQueries {
    fmt.Printf("Slow query: %s (took %v)\n", query.SQL, query.Duration)
}
```

### Query Pattern Analysis

```go
type QueryAnalyzer struct {
    patterns map[string]*QueryPattern
    mutex    sync.RWMutex
}

type QueryPattern struct {
    SQL           string
    Count         int64
    TotalDuration time.Duration
    MaxDuration   time.Duration
    MinDuration   time.Duration
    ErrorCount    int64
}

func (qa *QueryAnalyzer) AnalyzeQuery(sql string, duration time.Duration, err error) {
    qa.mutex.Lock()
    defer qa.mutex.Unlock()
    
    pattern, exists := qa.patterns[sql]
    if !exists {
        pattern = &QueryPattern{
            SQL:         sql,
            MinDuration: duration,
        }
        qa.patterns[sql] = pattern
    }
    
    pattern.Count++
    pattern.TotalDuration += duration
    
    if duration > pattern.MaxDuration {
        pattern.MaxDuration = duration
    }
    if duration < pattern.MinDuration {
        pattern.MinDuration = duration
    }
    
    if err != nil {
        pattern.ErrorCount++
    }
}

func (qa *QueryAnalyzer) GetTopQueries(limit int) []QueryPattern {
    qa.mutex.RLock()
    defer qa.mutex.RUnlock()
    
    var patterns []QueryPattern
    for _, pattern := range qa.patterns {
        patterns = append(patterns, *pattern)
    }
    
    // Sort by total time
    sort.Slice(patterns, func(i, j int) bool {
        return patterns[i].TotalDuration > patterns[j].TotalDuration
    })
    
    if len(patterns) > limit {
        patterns = patterns[:limit]
    }
    
    return patterns
}
```

## 📊 Performance Monitoring

### Real-time Performance Dashboard

```go
type PerformanceDashboard struct {
    metrics *QueryMetrics
    server  *http.Server
}

type QueryMetrics struct {
    TotalQueries    int64         `json:"total_queries"`
    QueriesPerSec   float64       `json:"queries_per_sec"`
    AvgDuration     time.Duration `json:"avg_duration"`
    SlowQueries     int64         `json:"slow_queries"`
    CacheHitRatio   float64       `json:"cache_hit_ratio"`
    ActiveConns     int           `json:"active_connections"`
    ErrorRate       float64       `json:"error_rate"`
    mutex           sync.RWMutex
}

func (pd *PerformanceDashboard) StartServer(port int) {
    mux := http.NewServeMux()
    
    // Metrics endpoint
    mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        pd.metrics.mutex.RLock()
        defer pd.metrics.mutex.RUnlock()
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(pd.metrics)
    })
    
    // Health check
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    pd.server = &http.Server{
        Addr:    fmt.Sprintf(":%d", port),
        Handler: mux,
    }
    
    log.Printf("Performance dashboard started on port %d", port)
    log.Fatal(pd.server.ListenAndServe())
}

func (pd *PerformanceDashboard) RecordQuery(duration time.Duration, cached bool, err error) {
    pd.metrics.mutex.Lock()
    defer pd.metrics.mutex.Unlock()
    
    pd.metrics.TotalQueries++
    
    // Update moving averages
    pd.updateMovingAverages(duration, cached, err)
}
```

### Alerting System

```go
type AlertingSystem struct {
    thresholds AlertThresholds
    notifiers  []Notifier
}

type AlertThresholds struct {
    SlowQueryThreshold    time.Duration
    HighErrorRate        float64
    LowCacheHitRate      float64
    HighConnectionUsage  float64
}

type Notifier interface {
    SendAlert(alert Alert)
}

type Alert struct {
    Type        string    `json:"type"`
    Message     string    `json:"message"`
    Severity    string    `json:"severity"`
    Timestamp   time.Time `json:"timestamp"`
    Metadata    map[string]any `json:"metadata,omitempty"`
}

func (as *AlertingSystem) CheckThresholds(metrics QueryMetrics) {
    // Check slow queries
    if metrics.AvgDuration > as.thresholds.SlowQueryThreshold {
        alert := Alert{
            Type:      "slow_queries",
            Message:   fmt.Sprintf("Average query duration is %v", metrics.AvgDuration),
            Severity:  "warning",
            Timestamp: time.Now(),
            Metadata: map[string]any{
                "avg_duration": metrics.AvgDuration,
                "threshold":    as.thresholds.SlowQueryThreshold,
            },
        }
        as.sendAlert(alert)
    }
    
    // Check error rate
    if metrics.ErrorRate > as.thresholds.HighErrorRate {
        alert := Alert{
            Type:      "high_error_rate",
            Message:   fmt.Sprintf("Error rate is %.2f%%", metrics.ErrorRate*100),
            Severity:  "critical",
            Timestamp: time.Now(),
            Metadata: map[string]any{
                "error_rate": metrics.ErrorRate,
                "threshold":  as.thresholds.HighErrorRate,
            },
        }
        as.sendAlert(alert)
    }
    
    // Check cache hit rate
    if metrics.CacheHitRatio < as.thresholds.LowCacheHitRate {
        alert := Alert{
            Type:      "low_cache_hit_rate", 
            Message:   fmt.Sprintf("Cache hit rate is %.2f%%", metrics.CacheHitRatio*100),
            Severity:  "info",
            Timestamp: time.Now(),
        }
        as.sendAlert(alert)
    }
}

func (as *AlertingSystem) sendAlert(alert Alert) {
    for _, notifier := range as.notifiers {
        go notifier.SendAlert(alert)
    }
}
```

## 🎯 Best Practices

### 1. Enable Appropriate Caching

```go
// ✅ Good - cache stable data
config := types.QueryOptimization{
    EnableQueryCache: true,
    CacheTTL:        5 * time.Minute, // Reasonable TTL
}

// ❌ Bad - caching frequently changing data
config := types.QueryOptimization{
    EnableQueryCache: true,
    CacheTTL:        24 * time.Hour, // Too long for dynamic data
}
```

### 2. Use Connection Pooling Wisely

```go
// ✅ Good - appropriate pool size
config.MaxOpenConns = 25  // Based on expected load
config.MaxIdleConns = 5   // Reasonable for idle connections

// ❌ Bad - unlimited connections
config.MaxOpenConns = 0   // Can overwhelm database
```

### 3. Monitor Performance Continuously

```go
// ✅ Good - comprehensive monitoring
optimizer := optimization.NewQueryOptimizer(config)
go func() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        stats := optimizer.GetQueryStats()
        logMetrics(stats)
        checkAlerts(stats)
    }
}()

// ❌ Bad - no monitoring
// optimizer := optimization.NewQueryOptimizer(config)
// // No monitoring setup
```

### 4. Handle Circuit Breaker States

```go
// ✅ Good - graceful degradation
cb := NewCircuitBreaker()
err := cb.Execute(func() error {
    return querybuilder.QB().Table("users").Get(ctx)
})

if err != nil {
    // Fallback to cache or alternative
    return getCachedUsers()
}

// ❌ Bad - no fallback mechanism
cb.Execute(queryFunction) // Fails without fallback
```

### 5. Optimize Query Patterns

```go
// ✅ Good - efficient query
users := querybuilder.QB().Table("users").
    Select("id", "name", "email").        // Only needed columns
    Where("status", "active").            // Indexed column
    Where("created_at", ">", yesterday).  // Limit date range
    Limit(100).                          // Reasonable limit
    Get(ctx)

// ❌ Bad - inefficient query  
users := querybuilder.QB().Table("users").  // SELECT *
    Where("LOWER(name)", "LIKE", "%john%"). // Function on column
    Get(ctx)                               // No limit
```

## 📋 Performance Checklist

- [ ] **Query Caching** - Enable for stable, frequently accessed data
- [ ] **Prepared Statements** - Automatic preparation for repeated queries
- [ ] **Connection Pooling** - Properly sized pools for your workload
- [ ] **Concurrency Control** - Limit concurrent queries to prevent overload
- [ ] **Slow Query Monitoring** - Track and optimize slow queries
- [ ] **Error Rate Monitoring** - Alert on high error rates
- [ ] **Cache Hit Ratio** - Monitor cache effectiveness
- [ ] **Circuit Breaker** - Implement for external service calls
- [ ] **Index Optimization** - Ensure proper database indexing
- [ ] **Query Analysis** - Regular analysis of query patterns

## 🔧 Troubleshooting Performance Issues

### High CPU Usage
```go
// Check for N+1 queries
// ❌ Bad - N+1 query pattern
for _, user := range users {
    posts, _ := querybuilder.QB().Table("posts").Where("user_id", user["id"]).Get(ctx)
}

// ✅ Good - single query with JOIN
userPosts, _ := querybuilder.QB().Table("users").
    Select("users.*", "posts.title").
    Join("posts", "users.id", "posts.user_id").
    Get(ctx)
```

### High Memory Usage
```go
// ✅ Good - process in chunks
err := querybuilder.QB().Table("large_table").
    Chunk(ctx, 1000, func(records types.Collection) error {
        processRecords(records)
        return nil
    })

// ❌ Bad - load all data at once
allRecords, _ := querybuilder.QB().Table("large_table").Get(ctx) // OOM risk
```

### Slow Queries
```go
// ✅ Good - optimized query
users := querybuilder.QB().Table("users").
    Select("id", "name").              // Only needed columns
    Where("status", "active").         // Use indexed column
    Where("id", ">", lastProcessedID). // Range query on primary key
    OrderBy("id", "asc").             // Use index order
    Limit(1000).                      // Reasonable batch size
    Get(ctx)
```

---

**Optimize your queries for lightning-fast performance and scalable applications!** ⚡

[← Async Operations](async-operations.md) | [Security Guide →](security.md)