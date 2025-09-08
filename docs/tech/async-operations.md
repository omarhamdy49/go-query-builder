# Async Operations

Master concurrent queries with goroutines and channels for lightning-fast database operations.

## 🚀 Overview

Go Query Builder provides built-in async functionality that leverages Go's powerful concurrency primitives. Execute multiple queries concurrently, implement non-blocking operations, and build high-performance applications.

## 📋 Table of Contents

- [Basic Async Queries](#basic-async-queries)
- [Concurrent Multiple Queries](#concurrent-multiple-queries)
- [Async Pagination](#async-pagination)
- [Query Racing](#query-racing)  
- [Pipeline Processing](#pipeline-processing)
- [Context Management](#context-management)
- [Error Handling](#error-handling)
- [Performance Patterns](#performance-patterns)

## 🔧 Basic Async Queries

### Single Async Query

```go
import (
    "context"
    "github.com/go-query-builder/querybuilder"
)

ctx := context.Background()

// Start async query
usersChan := querybuilder.QB().Table("users").GetAsync(ctx)

// Do other work while query executes
fmt.Println("Query running in background...")
performOtherWork()

// Get result when ready
result := <-usersChan
if result.Error != nil {
    log.Printf("Error: %v", result.Error)
} else {
    fmt.Printf("Loaded %d users\n", result.Data.Count())
}
```

### Async Count Operation

```go
// Async count
countChan := querybuilder.QB().Table("posts").
    Where("status", "published").
    CountAsync(ctx)

// Non-blocking - continues immediately
fmt.Println("Count query started...")

// Retrieve result
countResult := <-countChan
if countResult.Error == nil {
    fmt.Printf("Total posts: %d\n", countResult.Count)
}
```

### Async Pagination

```go
// Large dataset pagination
paginationChan := querybuilder.QB().Table("users").
    Where("age", ">=", 18).
    OrderBy("created_at", "desc").
    PaginateAsync(ctx, 1, 20)

// Process result asynchronously
paginationResult := <-paginationChan
if paginationResult.Error == nil {
    result := paginationResult.Result
    fmt.Printf("Page %d of %d (%d total)\n", 
        result.Meta.CurrentPage, 
        result.Meta.LastPage, 
        result.Meta.Total)
}
```

## 🎯 Concurrent Multiple Queries

### Parallel Query Execution

```go
import "sync"

var wg sync.WaitGroup
results := make(chan string, 3)

// Launch multiple queries concurrently
queries := []struct {
    name string
    fn   func() <-chan types.AsyncResult
}{
    {"Users", func() <-chan types.AsyncResult { 
        return querybuilder.QB().Table("users").GetAsync(ctx) 
    }},
    {"Posts", func() <-chan types.AsyncResult { 
        return querybuilder.QB().Table("posts").GetAsync(ctx) 
    }},
    {"Comments", func() <-chan types.AsyncResult { 
        return querybuilder.QB().Table("comments").GetAsync(ctx) 
    }},
}

startTime := time.Now()

for _, q := range queries {
    wg.Add(1)
    go func(name string, queryFn func() <-chan types.AsyncResult) {
        defer wg.Done()
        
        resultChan := queryFn()
        result := <-resultChan
        
        if result.Error != nil {
            results <- fmt.Sprintf("%s: Error - %v", name, result.Error)
        } else {
            results <- fmt.Sprintf("%s: %d records", name, result.Data.Count())
        }
    }(q.name, q.fn)
}

// Wait for all queries
go func() {
    wg.Wait()
    close(results)
}()

// Collect all results
for result := range results {
    fmt.Println(result)
}

fmt.Printf("All queries completed in %v\n", time.Since(startTime))
```

### Batch Processing with Worker Pool

```go
// Worker pool for processing large datasets
func processBatchAsync(ctx context.Context, batchSize int, workers int) {
    jobs := make(chan int, 100)
    results := make(chan types.AsyncResult, 100)
    
    // Start workers
    var wg sync.WaitGroup
    for w := 1; w <= workers; w++ {
        wg.Add(1)
        go worker(ctx, jobs, results, &wg)
    }
    
    // Send jobs
    go func() {
        defer close(jobs)
        for i := 0; i < 100; i += batchSize {
            jobs <- i
        }
    }()
    
    // Close results when workers finish
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Process results
    for result := range results {
        if result.Error == nil {
            fmt.Printf("Batch processed: %d records\n", result.Data.Count())
        }
    }
}

func worker(ctx context.Context, jobs <-chan int, results chan<- types.AsyncResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for offset := range jobs {
        resultChan := querybuilder.QB().Table("users").
            Limit(20).
            Offset(offset).
            GetAsync(ctx)
        
        results <- <-resultChan
    }
}
```

## 🏁 Query Racing

### Database Failover Pattern

```go
// Race primary and backup databases
primaryChan := querybuilder.QB().Table("users").GetAsync(ctx)
backupChan := querybuilder.Connection("backup").Table("users").GetAsync(ctx)

select {
case result := <-primaryChan:
    if result.Error == nil {
        fmt.Printf("Primary won: %d users\n", result.Data.Count())
        // Cancel backup query if needed
    } else {
        fmt.Printf("Primary failed: %v, using backup\n", result.Error)
        backupResult := <-backupChan
        if backupResult.Error == nil {
            fmt.Printf("Backup successful: %d users\n", backupResult.Data.Count())
        }
    }
case result := <-backupChan:
    fmt.Printf("Backup won: %d users\n", result.Data.Count())
case <-time.After(5 * time.Second):
    fmt.Println("Both databases timed out")
}
```

### Fastest Response Pattern

```go
// Get data from multiple sources, use fastest
sources := []func() <-chan types.AsyncResult{
    func() <-chan types.AsyncResult { 
        return querybuilder.QB().Table("cache_users").GetAsync(ctx) 
    },
    func() <-chan types.AsyncResult { 
        return querybuilder.QB().Table("users").GetAsync(ctx) 
    },
    func() <-chan types.AsyncResult { 
        return querybuilder.Connection("replica").Table("users").GetAsync(ctx) 
    },
}

// Fan-out to all sources
channels := make([]<-chan types.AsyncResult, len(sources))
for i, source := range sources {
    channels[i] = source()
}

// Use first successful response
for i, ch := range channels {
    select {
    case result := <-ch:
        if result.Error == nil {
            fmt.Printf("Source %d won with %d records\n", i, result.Data.Count())
            return
        }
    default:
        continue
    }
}
```

## 🔄 Pipeline Processing

### Async Data Pipeline

```go
// Stage 1: Fetch users
usersChan := querybuilder.QB().Table("users").
    Where("status", "active").
    Limit(1000).
    GetAsync(ctx)

// Stage 2: Process users as they arrive
go func() {
    result := <-usersChan
    if result.Error != nil {
        log.Printf("Stage 1 error: %v", result.Error)
        return
    }

    fmt.Printf("Stage 1: %d users loaded\n", result.Data.Count())

    // Stage 3: Async processing of each user
    var processingWg sync.WaitGroup
    semaphore := make(chan struct{}, 10) // Limit concurrent processing
    
    result.Data.Each(func(user map[string]any) bool {
        processingWg.Add(1)
        
        go func(userID any) {
            defer processingWg.Done()
            
            // Acquire semaphore
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // Process user (e.g., send email, update records)
            processUser(ctx, userID)
        }(user["id"])
        
        return true
    })

    processingWg.Wait()
    fmt.Println("Stage 2: All users processed")
}()

// Continue with other work...
time.Sleep(time.Second)
```

### Stream Processing

```go
// Continuous stream processing
func streamProcessor(ctx context.Context) {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Process new records asynchronously
            go processNewRecords(ctx)
        case <-ctx.Done():
            fmt.Println("Stream processor stopped")
            return
        }
    }
}

func processNewRecords(ctx context.Context) {
    lastProcessed := getLastProcessedTimestamp()
    
    newRecordsChan := querybuilder.QB().Table("events").
        Where("created_at", ">", lastProcessed).
        OrderBy("created_at", "asc").
        GetAsync(ctx)
    
    result := <-newRecordsChan
    if result.Error == nil && result.Data.Count() > 0 {
        fmt.Printf("Processing %d new records\n", result.Data.Count())
        
        // Process each record
        result.Data.Each(func(record map[string]any) bool {
            processRecord(record)
            return true
        })
        
        updateLastProcessedTimestamp()
    }
}
```

## ⏰ Context Management

### Timeout Control

```go
// Query with timeout
timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()

longQueryChan := querybuilder.QB().Table("users").
    Join("posts", "users.id", "posts.author_id").
    Join("comments", "posts.id", "comments.post_id").
    GetAsync(timeoutCtx)

select {
case result := <-longQueryChan:
    if result.Error == context.DeadlineExceeded {
        fmt.Println("Query timed out")
    } else if result.Error != nil {
        fmt.Printf("Query error: %v\n", result.Error)
    } else {
        fmt.Printf("Query completed: %d results\n", result.Data.Count())
    }
case <-timeoutCtx.Done():
    fmt.Println("Context cancelled")
}
```

### Cancellation Handling

```go
// Cancellable query
cancelCtx, cancel := context.WithCancel(ctx)

queryChan := querybuilder.QB().Table("large_table").GetAsync(cancelCtx)

// Cancel after 1 second
go func() {
    time.Sleep(time.Second)
    cancel()
}()

result := <-queryChan
if result.Error == context.Canceled {
    fmt.Println("Query was cancelled")
}
```

## 🛠️ Error Handling

### Graceful Error Recovery

```go
func resilientAsyncQuery(ctx context.Context, retries int) (*types.Collection, error) {
    for attempt := 0; attempt <= retries; attempt++ {
        resultChan := querybuilder.QB().Table("users").GetAsync(ctx)
        result := <-resultChan
        
        if result.Error == nil {
            return &result.Data, nil
        }
        
        fmt.Printf("Attempt %d failed: %v\n", attempt+1, result.Error)
        
        if attempt < retries {
            backoff := time.Duration(attempt+1) * time.Second
            fmt.Printf("Retrying in %v...\n", backoff)
            time.Sleep(backoff)
        }
    }
    
    return nil, fmt.Errorf("query failed after %d attempts", retries+1)
}

// Usage
users, err := resilientAsyncQuery(ctx, 3)
if err != nil {
    log.Printf("All attempts failed: %v", err)
}
```

### Error Aggregation

```go
// Collect errors from multiple async queries
type QueryResult struct {
    Name   string
    Data   types.Collection
    Error  error
}

func executeMultipleQueries(ctx context.Context) []QueryResult {
    queries := map[string]func() <-chan types.AsyncResult{
        "users":    func() <-chan types.AsyncResult { return querybuilder.QB().Table("users").GetAsync(ctx) },
        "posts":    func() <-chan types.AsyncResult { return querybuilder.QB().Table("posts").GetAsync(ctx) },
        "comments": func() <-chan types.AsyncResult { return querybuilder.QB().Table("comments").GetAsync(ctx) },
    }
    
    resultsChan := make(chan QueryResult, len(queries))
    var wg sync.WaitGroup
    
    // Execute all queries
    for name, queryFn := range queries {
        wg.Add(1)
        go func(queryName string, fn func() <-chan types.AsyncResult) {
            defer wg.Done()
            
            asyncResult := <-fn()
            resultsChan <- QueryResult{
                Name:  queryName,
                Data:  asyncResult.Data,
                Error: asyncResult.Error,
            }
        }(name, queryFn)
    }
    
    // Wait and collect results
    go func() {
        wg.Wait()
        close(resultsChan)
    }()
    
    var results []QueryResult
    for result := range resultsChan {
        results = append(results, result)
    }
    
    return results
}
```

## 🚀 Performance Patterns

### Fan-Out Fan-In

```go
// Fan-out: distribute work
func fanOutFanIn(ctx context.Context, userIDs []int) []map[string]any {
    // Fan-out: create channels for each user
    channels := make([]<-chan types.AsyncResult, len(userIDs))
    
    for i, userID := range userIDs {
        channels[i] = querybuilder.QB().Table("users").
            Select("users.*", "profiles.bio").
            LeftJoin("profiles", "users.id", "profiles.user_id").
            Where("users.id", userID).
            GetAsync(ctx)
    }
    
    // Fan-in: collect results
    var results []map[string]any
    for i, ch := range channels {
        result := <-ch
        if result.Error == nil && result.Data.Count() > 0 {
            userDetail := result.Data.First()
            results = append(results, userDetail)
        } else {
            log.Printf("User %d failed: %v", userIDs[i], result.Error)
        }
    }
    
    return results
}
```

### Async Batch Operations

```go
// Process large datasets in batches
func processBatchesAsync(ctx context.Context, batchSize int) {
    offset := 0
    
    for {
        batchChan := querybuilder.QB().Table("large_table").
            Limit(batchSize).
            Offset(offset).
            GetAsync(ctx)
            
        result := <-batchChan
        if result.Error != nil {
            log.Printf("Batch error: %v", result.Error)
            break
        }
        
        if result.Data.Count() == 0 {
            fmt.Println("No more data to process")
            break
        }
        
        // Process batch asynchronously
        go processBatch(result.Data)
        
        offset += batchSize
        
        // Rate limiting
        time.Sleep(100 * time.Millisecond)
    }
}

func processBatch(data types.Collection) {
    fmt.Printf("Processing batch of %d records\n", data.Count())
    
    data.Each(func(record map[string]any) bool {
        // Process each record
        processRecord(record)
        return true
    })
}
```

## 📊 Monitoring Async Operations

```go
// Monitor async query performance
type AsyncMonitor struct {
    queries    int64
    totalTime  time.Duration
    errors     int64
    mutex      sync.RWMutex
}

func (m *AsyncMonitor) RecordQuery(duration time.Duration, err error) {
    m.mutex.Lock()
    defer m.mutex.Unlock()
    
    m.queries++
    m.totalTime += duration
    if err != nil {
        m.errors++
    }
}

func (m *AsyncMonitor) Stats() (int64, time.Duration, int64) {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    
    avgTime := time.Duration(0)
    if m.queries > 0 {
        avgTime = m.totalTime / time.Duration(m.queries)
    }
    
    return m.queries, avgTime, m.errors
}

// Monitored async query
func monitoredAsyncQuery(ctx context.Context, monitor *AsyncMonitor) {
    start := time.Now()
    
    resultChan := querybuilder.QB().Table("users").GetAsync(ctx)
    result := <-resultChan
    
    monitor.RecordQuery(time.Since(start), result.Error)
}
```

## 🎯 Best Practices

### 1. Always Use Context
```go
// ✅ Good - with context
usersChan := querybuilder.QB().Table("users").GetAsync(ctx)

// ❌ Bad - without context (blocks indefinitely)
usersChan := querybuilder.QB().Table("users").GetAsync(context.Background())
```

### 2. Handle Errors Properly
```go
// ✅ Good - error handling
result := <-usersChan
if result.Error != nil {
    log.Printf("Query failed: %v", result.Error)
    return
}

// ❌ Bad - ignoring errors
result := <-usersChan
users := result.Data // might be nil!
```

### 3. Use Buffered Channels for Known Sizes
```go
// ✅ Good - buffered channel
results := make(chan types.AsyncResult, 10)

// ❌ Bad - unbuffered can cause blocking
results := make(chan types.AsyncResult)
```

### 4. Implement Timeouts
```go
// ✅ Good - with timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// ❌ Bad - no timeout (can hang forever)
ctx := context.Background()
```

### 5. Clean Up Resources
```go
// ✅ Good - cleanup
defer func() {
    if recovery := recover(); recovery != nil {
        log.Printf("Recovered: %v", recovery)
    }
}()
```

## 📈 Performance Tips

1. **Concurrent Queries**: Use async for I/O bound operations
2. **Worker Pools**: Limit concurrent goroutines with semaphores  
3. **Context Timeouts**: Always set reasonable timeouts
4. **Error Handling**: Handle errors gracefully with retries
5. **Resource Management**: Clean up channels and goroutines
6. **Monitoring**: Track performance metrics for optimization

---

**Master async operations to build lightning-fast, concurrent Go applications!** ⚡

[← Query Builder](query-builder.md) | [Performance Optimization →](optimization.md)