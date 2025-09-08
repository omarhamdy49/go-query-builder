# Performance Guide - High-Speed Database Operations

## Table of Contents

1. [Performance Overview](#performance-overview)
2. [Query Optimization](#query-optimization)
3. [Connection Pool Tuning](#connection-pool-tuning)
4. [Caching Strategies](#caching-strategies)
5. [Async Operations](#async-operations)
6. [Indexing Guidelines](#indexing-guidelines)
7. [Monitoring & Profiling](#monitoring--profiling)
8. [Performance Best Practices](#performance-best-practices)
9. [Benchmarking](#benchmarking)
10. [Troubleshooting](#troubleshooting)

---

## Performance Overview

The Go Query Builder is designed for **high-performance database operations** with:

- **Sub-millisecond Query Execution** - Optimized SQL generation and execution
- **Intelligent Connection Pooling** - Efficient resource management
- **Query Result Caching** - Automatic result caching with TTL
- **Async Operations** - Non-blocking concurrent query execution
- **Performance Monitoring** - Real-time performance metrics and alerting
- **Automatic Optimization** - Query plan analysis and optimization hints

### Performance Metrics

```
┌─────────────────────────────────────────────────────┐
│                Performance Targets                  │
├─────────────────────────────────────────────────────┤
│ Simple Query Execution:     < 1ms                  │
│ Complex Query (JOINs):      < 5ms                  │
│ Paginated Results:          < 3ms                  │
│ Batch Inserts (100 rows):   < 10ms                 │
│ Connection Acquisition:      < 100μs               │
│ Cache Hit Response:         < 50μs                 │
│ Concurrent Connections:     1000+                  │
│ Throughput:                 10,000+ queries/sec    │
└─────────────────────────────────────────────────────┘
```

---

## Query Optimization

### Efficient Query Building

#### Use Specific Columns
```go
// ✅ FAST: Select only needed columns
users, err := querybuilder.QB().Table("users").
    Select("id", "name", "email").           // Only needed columns
    Where("active", true).
    Get(ctx)

// ❌ SLOW: Select all columns
users, err := querybuilder.QB().Table("users").
    Select("*").                             // All columns (wasteful)
    Where("active", true).
    Get(ctx)
```

#### Optimize WHERE Conditions
```go
// ✅ FAST: Use indexed columns in WHERE
qb.Where("user_id", userID)                 // Uses PRIMARY KEY index
qb.Where("email", email)                    // Uses UNIQUE index
qb.Where("created_at", ">", lastWeek)       // Uses INDEX on created_at

// ❌ SLOW: Non-indexed columns
qb.Where("JSON_EXTRACT(metadata, '$.key')", value) // No index on JSON path
qb.Where("UPPER(name)", "JOHN")             // Function prevents index usage
```

#### Limit Result Sets
```go
// ✅ FAST: Always use LIMIT for large datasets
qb.Where("active", true).
   OrderBy("created_at", "desc").
   Limit(20).                               // Limit to reasonable size
   Get(ctx)

// Add pagination for complete results
result, err := qb.Paginate(ctx, 1, 20)     // Efficient pagination
```

#### Optimize JOIN Operations
```go
// ✅ FAST: JOIN on indexed columns
qb.Table("users").
   Join("profiles", "users.id", "profiles.user_id").    // Both columns indexed
   Where("users.active", true).
   Select("users.name", "profiles.bio")

// ✅ FAST: Use appropriate JOIN types
qb.Table("users").
   LeftJoin("profiles", "users.id", "profiles.user_id"). // LEFT JOIN when needed
   Where("users.active", true)

// ❌ SLOW: Unnecessary JOINs
qb.Table("users").
   Join("profiles", "users.id", "profiles.user_id").
   Join("settings", "users.id", "settings.user_id").    // Too many JOINs
   Join("logs", "users.id", "logs.user_id")             // Avoid if not needed
```

### Query Analysis Tools

#### Built-in Performance Analysis
```go
// Enable query performance tracking
qb := querybuilder.QB().EnableProfiling()

result, err := qb.Table("users").
    Where("active", true).
    Get(ctx)

// Get performance metrics
metrics := qb.GetQueryMetrics()
fmt.Printf("Execution time: %v\n", metrics.ExecutionTime)
fmt.Printf("Rows examined: %d\n", metrics.RowsExamined)
fmt.Printf("Rows returned: %d\n", metrics.RowsReturned)
fmt.Printf("Index used: %s\n", metrics.IndexUsed)
```

#### Query Plan Analysis
```go
// Analyze query execution plan
plan, err := qb.Table("users").
    Where("email", "user@example.com").
    ExplainQuery(ctx)

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Query plan:\n%s\n", plan)
// Output: Shows EXPLAIN output with index usage, join types, etc.
```

#### Slow Query Detection
```go
// Automatically log slow queries
type SlowQueryDetector struct {
    threshold time.Duration
}

func (d *SlowQueryDetector) CheckQuery(duration time.Duration, query string) {
    if duration > d.threshold {
        log.Printf("SLOW QUERY [%v]: %s", duration, query)
        // Send to monitoring system
        sendSlowQueryAlert(duration, query)
    }
}
```

---

## Connection Pool Tuning

### Pool Configuration

#### Optimal Pool Settings
```env
# .env configuration for high performance
DB_MAX_OPEN_CONNS=50              # Maximum concurrent connections
DB_MAX_IDLE_CONNS=10              # Keep connections ready
DB_MAX_LIFETIME=30m               # Rotate connections regularly  
DB_MAX_IDLE_TIME=5m               # Close unused connections
DB_CONNECTION_TIMEOUT=10s         # Fast connection establishment
DB_QUERY_TIMEOUT=30s              # Prevent hanging queries
```

#### Dynamic Pool Sizing
```go
type DynamicConnectionPool struct {
    minConnections    int
    maxConnections    int
    currentLoad       float64
    scalingFactor     float64
}

func (p *DynamicConnectionPool) AdjustPoolSize() {
    targetSize := int(float64(p.minConnections) + 
                     (p.currentLoad * p.scalingFactor))
    
    if targetSize > p.maxConnections {
        targetSize = p.maxConnections
    }
    
    // Adjust pool size based on current load
    adjustConnectionPool(targetSize)
}

// Monitor and adjust pool size every minute
go func() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        pool.currentLoad = getCurrentLoadMetrics()
        pool.AdjustPoolSize()
    }
}()
```

### Connection Health Monitoring

#### Health Check System
```go
type ConnectionHealthMonitor struct {
    unhealthyConnections map[*sql.Conn]time.Time
    checkInterval        time.Duration
}

func (m *ConnectionHealthMonitor) MonitorHealth() {
    ticker := time.NewTicker(m.checkInterval)
    defer ticker.Stop()
    
    for range ticker.C {
        m.checkConnectionHealth()
        m.removeUnhealthyConnections()
    }
}

func (m *ConnectionHealthMonitor) checkConnectionHealth() {
    connections := getActiveConnections()
    
    for _, conn := range connections {
        if !m.isConnectionHealthy(conn) {
            m.unhealthyConnections[conn] = time.Now()
            logConnectionIssue(conn)
        }
    }
}
```

#### Connection Metrics
```go
type ConnectionMetrics struct {
    TotalConnections    int           `json:"total_connections"`
    ActiveConnections   int           `json:"active_connections"`
    IdleConnections     int           `json:"idle_connections"`
    ConnectionsInUse    int           `json:"connections_in_use"`
    MaxOpenReached      int           `json:"max_open_reached"`
    WaitCount           int64         `json:"wait_count"`
    WaitDuration        time.Duration `json:"wait_duration"`
    MaxIdleClosed       int64         `json:"max_idle_closed"`
    MaxLifetimeClosed   int64         `json:"max_lifetime_closed"`
}

func GetConnectionMetrics() ConnectionMetrics {
    stats := db.Stats()
    return ConnectionMetrics{
        TotalConnections:    stats.OpenConnections,
        ActiveConnections:   stats.InUse,
        IdleConnections:     stats.Idle,
        MaxOpenReached:      stats.MaxOpenConnections,
        WaitCount:           stats.WaitCount,
        WaitDuration:        stats.WaitDuration,
        MaxIdleClosed:       stats.MaxIdleClosed,
        MaxLifetimeClosed:   stats.MaxLifetimeClosed,
    }
}
```

---

## Caching Strategies

### Query Result Caching

#### Automatic Caching
```go
// Cache query results for 5 minutes
users, err := querybuilder.QB().Table("users").
    Where("active", true).
    Cache(5 * time.Minute).                  // Automatic caching
    Get(ctx)

// Subsequent identical queries return cached results
users2, err := querybuilder.QB().Table("users").
    Where("active", true).
    Get(ctx)  // Returns cached result (sub-50μs response)
```

#### Cache Key Generation
```go
type QueryCacheKey struct {
    SQL      string
    Bindings []interface{}
    Hash     string
}

func generateCacheKey(sql string, bindings []interface{}) string {
    hasher := sha256.New()
    
    // Include SQL query
    hasher.Write([]byte(sql))
    
    // Include parameter bindings
    for _, binding := range bindings {
        hasher.Write([]byte(fmt.Sprintf("%v", binding)))
    }
    
    return hex.EncodeToString(hasher.Sum(nil))
}
```

#### Multi-Level Caching
```go
type MultiLevelCache struct {
    l1Cache *sync.Map          // In-memory cache (fastest)
    l2Cache RedisClient        // Redis cache (fast)
    l3Cache DatabaseCache      // Database cache (slower but persistent)
}

func (c *MultiLevelCache) Get(key string) (interface{}, bool) {
    // Try L1 cache first (in-memory)
    if value, exists := c.l1Cache.Load(key); exists {
        cacheHitCounter.With(prometheus.Labels{"level": "l1"}).Inc()
        return value, true
    }
    
    // Try L2 cache (Redis)
    if value, err := c.l2Cache.Get(key); err == nil {
        // Populate L1 cache for next time
        c.l1Cache.Store(key, value)
        cacheHitCounter.With(prometheus.Labels{"level": "l2"}).Inc()
        return value, true
    }
    
    // Try L3 cache (Database)
    if value, err := c.l3Cache.Get(key); err == nil {
        // Populate L2 and L1 caches
        c.l2Cache.Set(key, value, time.Hour)
        c.l1Cache.Store(key, value)
        cacheHitCounter.With(prometheus.Labels{"level": "l3"}).Inc()
        return value, true
    }
    
    cacheMissCounter.Inc()
    return nil, false
}
```

#### Cache Invalidation
```go
// Automatic cache invalidation on data changes
func invalidateCache(table string, conditions map[string]interface{}) {
    // Find all cache keys that might be affected
    affectedKeys := findAffectedCacheKeys(table, conditions)
    
    // Remove from all cache levels
    for _, key := range affectedKeys {
        cache.l1Cache.Delete(key)
        cache.l2Cache.Del(key)
        cache.l3Cache.Delete(key)
    }
    
    logCacheInvalidation(table, len(affectedKeys))
}

// Automatic invalidation on INSERT/UPDATE/DELETE
qb.Table("users").Update(ctx, values) // Automatically invalidates user caches
```

### Smart Caching Policies

#### Time-Based Expiration
```go
type CachePolicy struct {
    StaticData    time.Duration // 1 hour for rarely changing data
    UserProfiles  time.Duration // 15 minutes for user data
    DynamicData   time.Duration // 1 minute for frequently changing data
    RealtimeData  time.Duration // 10 seconds for real-time data
}

func getCacheDuration(table string) time.Duration {
    policy := CachePolicy{
        StaticData:   time.Hour,
        UserProfiles: 15 * time.Minute,
        DynamicData:  time.Minute,
        RealtimeData: 10 * time.Second,
    }
    
    switch table {
    case "countries", "categories", "settings":
        return policy.StaticData
    case "users", "profiles":
        return policy.UserProfiles
    case "orders", "payments":
        return policy.DynamicData
    case "notifications", "messages":
        return policy.RealtimeData
    default:
        return policy.DynamicData
    }
}
```

#### Usage-Based Caching
```go
type UsageBasedCache struct {
    queryFrequency map[string]int
    cacheHits      map[string]int
    lastAccess     map[string]time.Time
}

func (c *UsageBasedCache) ShouldCache(queryKey string) bool {
    frequency := c.queryFrequency[queryKey]
    hits := c.cacheHits[queryKey]
    lastAccess := c.lastAccess[queryKey]
    
    // Cache if frequently accessed
    if frequency > 10 {
        return true
    }
    
    // Cache if high hit rate
    if hits > frequency/2 {
        return true
    }
    
    // Don't cache if not accessed recently
    if time.Since(lastAccess) > time.Hour {
        return false
    }
    
    return frequency > 3 // Cache if accessed more than 3 times
}
```

---

## Async Operations

### Concurrent Query Execution

#### Basic Async Operations
```go
// Execute multiple queries concurrently
ctx := context.Background()

// Start multiple queries
usersChan := querybuilder.QB().Table("users").GetAsync(ctx)
ordersChan := querybuilder.QB().Table("orders").GetAsync(ctx)
productsChan := querybuilder.QB().Table("products").GetAsync(ctx)

// Wait for results
usersResult := <-usersChan
ordersResult := <-ordersChan
productsResult := <-productsChan

// Check for errors
if usersResult.Error != nil {
    log.Printf("Users query failed: %v", usersResult.Error)
}

if ordersResult.Error != nil {
    log.Printf("Orders query failed: %v", ordersResult.Error)
}

fmt.Printf("Loaded %d users, %d orders, %d products concurrently\n",
    usersResult.Data.Count(),
    ordersResult.Data.Count(), 
    productsResult.Data.Count())
```

#### Query Racing
```go
// Execute same query on multiple replicas and use fastest result
func queryWithRacing(ctx context.Context, query QueryBuilder) (Collection, error) {
    // Create channels for each replica
    replica1Chan := query.Connection("replica1").GetAsync(ctx)
    replica2Chan := query.Connection("replica2").GetAsync(ctx)
    replica3Chan := query.Connection("replica3").GetAsync(ctx)
    
    // Use the fastest response
    select {
    case result1 := <-replica1Chan:
        if result1.Error == nil {
            return result1.Data, nil
        }
    case result2 := <-replica2Chan:
        if result2.Error == nil {
            return result2.Data, nil
        }
    case result3 := <-replica3Chan:
        if result3.Error == nil {
            return result3.Data, nil
        }
    case <-ctx.Done():
        return nil, ctx.Err()
    }
    
    // If all failed, wait for any non-error result
    for i := 0; i < 3; i++ {
        select {
        case result1 := <-replica1Chan:
            if result1.Error == nil {
                return result1.Data, nil
            }
        case result2 := <-replica2Chan:
            if result2.Error == nil {
                return result2.Data, nil
            }
        case result3 := <-replica3Chan:
            if result3.Error == nil {
                return result3.Data, nil
            }
        }
    }
    
    return nil, errors.New("all replicas failed")
}
```

#### Pipeline Processing
```go
// Process queries in pipeline for better throughput
type QueryPipeline struct {
    inputChan  chan QueryRequest
    outputChan chan QueryResult
    workers    int
}

func (p *QueryPipeline) Start() {
    for i := 0; i < p.workers; i++ {
        go p.worker()
    }
}

func (p *QueryPipeline) worker() {
    for request := range p.inputChan {
        result := p.executeQuery(request)
        p.outputChan <- result
    }
}

// Usage
pipeline := &QueryPipeline{
    inputChan:  make(chan QueryRequest, 100),
    outputChan: make(chan QueryResult, 100),
    workers:    10,
}
pipeline.Start()

// Submit queries to pipeline
pipeline.inputChan <- QueryRequest{Query: "SELECT * FROM users WHERE id = ?", Args: []interface{}{123}}
result := <-pipeline.outputChan
```

### Batch Operations

#### Optimized Batch Inserts
```go
// Insert large batches efficiently
func insertLargeBatch(ctx context.Context, table string, records []map[string]interface{}) error {
    batchSize := 1000  // Optimal batch size
    
    // Split into batches
    for i := 0; i < len(records); i += batchSize {
        end := i + batchSize
        if end > len(records) {
            end = len(records)
        }
        
        batch := records[i:end]
        
        // Insert batch with retry logic
        err := insertBatchWithRetry(ctx, table, batch, 3)
        if err != nil {
            return fmt.Errorf("batch insert failed at offset %d: %w", i, err)
        }
    }
    
    return nil
}

func insertBatchWithRetry(ctx context.Context, table string, batch []map[string]interface{}, maxRetries int) error {
    for attempt := 1; attempt <= maxRetries; attempt++ {
        err := querybuilder.QB().Table(table).InsertBatch(ctx, batch)
        if err == nil {
            return nil
        }
        
        // Exponential backoff
        backoff := time.Duration(attempt*attempt) * time.Second
        time.Sleep(backoff)
        
        log.Printf("Batch insert attempt %d failed: %v", attempt, err)
    }
    
    return errors.New("batch insert failed after all retries")
}
```

#### Parallel Batch Processing
```go
// Process batches in parallel
func parallelBatchInsert(ctx context.Context, table string, records []map[string]interface{}) error {
    batchSize := 500
    maxWorkers := 5
    
    // Create work queue
    workQueue := make(chan []map[string]interface{}, maxWorkers)
    errorChan := make(chan error, maxWorkers)
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < maxWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for batch := range workQueue {
                if err := querybuilder.QB().Table(table).InsertBatch(ctx, batch); err != nil {
                    errorChan <- err
                    return
                }
            }
        }()
    }
    
    // Submit batches
    go func() {
        defer close(workQueue)
        for i := 0; i < len(records); i += batchSize {
            end := i + batchSize
            if end > len(records) {
                end = len(records)
            }
            workQueue <- records[i:end]
        }
    }()
    
    // Wait for completion
    wg.Wait()
    close(errorChan)
    
    // Check for errors
    for err := range errorChan {
        return err
    }
    
    return nil
}
```

---

## Indexing Guidelines

### Index Strategy

#### Primary Key Optimization
```sql
-- ✅ GOOD: Use appropriate primary key types
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,  -- Fast, sequential
    uuid CHAR(36) NOT NULL UNIQUE,                  -- For external references
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ❌ AVOID: UUID as primary key (causes page splits)
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,  -- Random UUIDs cause fragmentation
    -- ...
);
```

#### Query-Specific Indexes
```go
// For this query pattern:
qb.Table("orders").
   Where("user_id", userID).
   Where("status", "pending").
   OrderBy("created_at", "desc")

// Create compound index:
// CREATE INDEX idx_orders_user_status_created ON orders (user_id, status, created_at);
```

#### Index Usage Hints
```go
// Force index usage for critical queries
qb.Table("users").
   UseIndex("idx_email").            // Hint to use specific index
   Where("email", email).
   First(ctx)

// Monitor index usage
explainResult, err := qb.Table("users").
    Where("email", email).
    ExplainQuery(ctx)
// Check if "idx_email" is used in execution plan
```

### Index Monitoring

#### Index Efficiency Analysis
```go
type IndexUsageStats struct {
    IndexName    string
    TableName    string
    TimesUsed    int64
    LastUsed     time.Time
    Efficiency   float64  // Ratio of rows returned to rows examined
}

func analyzeIndexUsage() []IndexUsageStats {
    // Query database for index usage statistics
    query := `
        SELECT 
            index_name,
            table_name,
            cardinality,
            last_update
        FROM INFORMATION_SCHEMA.STATISTICS 
        WHERE table_schema = DATABASE()
    `
    
    // Analyze and return index efficiency metrics
    return calculateIndexEfficiency(query)
}
```

#### Unused Index Detection
```go
func findUnusedIndexes() []string {
    unusedIndexes := []string{}
    
    // Query for indexes never used
    query := `
        SELECT DISTINCT
            s.index_name
        FROM INFORMATION_SCHEMA.STATISTICS s
        LEFT JOIN performance_schema.table_io_waits_summary_by_index_usage i
            ON s.table_schema = i.object_schema 
            AND s.table_name = i.object_name
            AND s.index_name = i.index_name
        WHERE s.table_schema = DATABASE()
            AND i.index_name IS NULL
            AND s.index_name != 'PRIMARY'
    `
    
    // Execute query and collect unused indexes
    return executeUnusedIndexQuery(query)
}
```

---

## Monitoring & Profiling

### Performance Metrics

#### Real-time Metrics Collection
```go
type QueryMetrics struct {
    TotalQueries     int64         `json:"total_queries"`
    AverageLatency   time.Duration `json:"average_latency"`
    P95Latency       time.Duration `json:"p95_latency"`
    P99Latency       time.Duration `json:"p99_latency"`
    ErrorRate        float64       `json:"error_rate"`
    CacheHitRate     float64       `json:"cache_hit_rate"`
    SlowQueries      int64         `json:"slow_queries"`
    ConnectionsUsed  int           `json:"connections_used"`
}

// Prometheus metrics
var (
    queryDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "query_duration_seconds",
            Help:    "Query execution duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"table", "operation"},
    )
    
    queryTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "queries_total",
            Help: "Total number of queries executed",
        },
        []string{"table", "operation", "status"},
    )
)
```

#### Performance Dashboard
```go
type PerformanceDashboard struct {
    CurrentMetrics    QueryMetrics
    HistoricalMetrics []QueryMetrics
    ActiveQueries     []ActiveQuery
    SystemHealth      SystemHealthStatus
}

type ActiveQuery struct {
    ID          string        `json:"id"`
    Query       string        `json:"query"`
    StartTime   time.Time     `json:"start_time"`
    Duration    time.Duration `json:"duration"`
    User        string        `json:"user"`
    State       string        `json:"state"`
}

func (d *PerformanceDashboard) UpdateMetrics() {
    d.CurrentMetrics = collectCurrentMetrics()
    d.ActiveQueries = getActiveQueries()
    d.SystemHealth = checkSystemHealth()
    
    // Store historical data
    d.HistoricalMetrics = append(d.HistoricalMetrics, d.CurrentMetrics)
    
    // Keep only last 24 hours of metrics
    if len(d.HistoricalMetrics) > 1440 { // 24 hours * 60 minutes
        d.HistoricalMetrics = d.HistoricalMetrics[1:]
    }
}
```

### Profiling Tools

#### Query Profiler
```go
type QueryProfiler struct {
    enabled      bool
    profiles     map[string]QueryProfile
    mutex        sync.RWMutex
}

type QueryProfile struct {
    Query            string
    ExecutionCount   int64
    TotalTime        time.Duration
    MinTime          time.Duration
    MaxTime          time.Duration
    AverageTime      time.Duration
    LastExecuted     time.Time
    ExecutionPlan    string
    IndexesUsed      []string
}

func (p *QueryProfiler) ProfileQuery(query string, duration time.Duration, plan string) {
    if !p.enabled {
        return
    }
    
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    profile := p.profiles[query]
    profile.Query = query
    profile.ExecutionCount++
    profile.TotalTime += duration
    profile.LastExecuted = time.Now()
    
    if profile.MinTime == 0 || duration < profile.MinTime {
        profile.MinTime = duration
    }
    
    if duration > profile.MaxTime {
        profile.MaxTime = duration
    }
    
    profile.AverageTime = profile.TotalTime / time.Duration(profile.ExecutionCount)
    profile.ExecutionPlan = plan
    
    p.profiles[query] = profile
}
```

#### Memory Profiler
```go
func profileMemoryUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    metrics := MemoryMetrics{
        Alloc:         m.Alloc,
        TotalAlloc:    m.TotalAlloc,
        Sys:           m.Sys,
        NumGC:         m.NumGC,
        GCCPUFraction: m.GCCPUFraction,
    }
    
    // Log if memory usage is high
    if metrics.Alloc > 100*1024*1024 { // 100MB
        log.Printf("High memory usage: %d MB allocated", metrics.Alloc/(1024*1024))
    }
    
    // Force GC if memory is very high
    if metrics.Alloc > 500*1024*1024 { // 500MB
        runtime.GC()
        log.Println("Forced garbage collection due to high memory usage")
    }
}
```

---

## Performance Best Practices

### Query Design

#### 1. Use Appropriate Data Types
```go
// ✅ GOOD: Use appropriate sizes
type User struct {
    ID       int32     `json:"id"`           // 4 bytes vs int64 (8 bytes)
    Name     string    `json:"name"`         // VARCHAR(255) vs TEXT
    Active   bool      `json:"active"`       // 1 byte vs string
    Created  time.Time `json:"created_at"`   // Proper timestamp type
}

// ❌ AVOID: Oversized types
type User struct {
    ID       int64     `json:"id"`           // Unnecessary for most use cases
    Name     string    `json:"name"`         // Using TEXT for short strings
    Active   string    `json:"active"`       // "true"/"false" vs boolean
}
```

#### 2. Optimize Query Patterns
```go
// ✅ GOOD: Batch similar queries
userIDs := []interface{}{1, 2, 3, 4, 5}
users, err := qb.Table("users").
    WhereIn("id", userIDs).
    Get(ctx)

// ❌ AVOID: Multiple single queries
for _, id := range userIDs {
    user, err := qb.Table("users").Where("id", id).First(ctx)
    // This creates N database round trips
}
```

#### 3. Use Pagination for Large Datasets
```go
// ✅ GOOD: Paginate large results
result, err := qb.Table("orders").
    Where("created_at", ">=", lastMonth).
    OrderBy("id", "desc").
    Paginate(ctx, 1, 100)      // Process 100 at a time

// ❌ AVOID: Loading everything at once
allOrders, err := qb.Table("orders").
    Where("created_at", ">=", lastMonth).
    Get(ctx)  // Could return millions of records
```

### Connection Management

#### 1. Optimize Connection Pool
```go
// Configure based on your workload
config := ConnectionConfig{
    MaxOpenConns:    50,              // Number of CPU cores * 2-4
    MaxIdleConns:    10,              // 20% of MaxOpenConns
    ConnMaxLifetime: 30 * time.Minute, // Rotate connections regularly
    ConnMaxIdleTime: 5 * time.Minute,  // Close idle connections
}
```

#### 2. Use Context Timeouts
```go
// Always set reasonable timeouts
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := qb.Get(ctx)
if err != nil {
    if err == context.DeadlineExceeded {
        log.Println("Query timed out after 30 seconds")
        // Handle timeout specifically
    }
    return err
}
```

#### 3. Monitor Connection Health
```go
// Regular health checks
func monitorConnectionHealth() {
    ticker := time.NewTicker(time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        stats := db.Stats()
        
        // Alert if too many connections waiting
        if stats.WaitCount > 100 {
            log.Printf("High connection wait count: %d", stats.WaitCount)
            alertOps("High database connection contention")
        }
        
        // Alert if wait duration is too long
        if stats.WaitDuration > 5*time.Second {
            log.Printf("Long connection wait duration: %v", stats.WaitDuration)
            alertOps("Slow database connection acquisition")
        }
    }
}
```

---

Your Go Query Builder delivers **enterprise-grade performance** with sub-millisecond query execution, intelligent caching, and advanced optimization features.

🚀 **Performance Level: HIGH-SPEED ENTERPRISE DATABASE OPERATIONS** ⚡