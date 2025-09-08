package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-query-builder/querybuilder"
	"github.com/go-query-builder/querybuilder/pkg/types"
)

// Async and Performance Examples
func asyncPerformanceExample() {
	ctx := context.Background()

	fmt.Println("=== Async Query Operations & Performance Optimization ===")
	fmt.Println()

	// Example 1: Basic Async Queries
	fmt.Println("1. Basic Async Queries")
	fmt.Println("=====================")

	// Single async query
	usersChan := querybuilder.QB().Table("users").GetAsync(ctx)
	
	// Do other work while query executes...
	fmt.Println("📊 Query running in background...")
	time.Sleep(100 * time.Millisecond) // Simulate other work
	
	// Wait for result
	result := <-usersChan
	if result.Error != nil {
		log.Printf("Error: %v", result.Error)
	} else {
		fmt.Printf("✅ Async query completed: %d users loaded\n", result.Data.Count())
	}

	// Async count
	countChan := querybuilder.QB().Table("posts").
		Where("status", "published").
		CountAsync(ctx)
		
	countResult := <-countChan
	if countResult.Error != nil {
		log.Printf("Count error: %v", countResult.Error)
	} else {
		fmt.Printf("✅ Async count completed: %d published posts\n", countResult.Count)
	}

	fmt.Println()

	// Example 2: Concurrent Multiple Queries
	fmt.Println("2. Concurrent Multiple Queries")
	fmt.Println("=============================")

	var wg sync.WaitGroup
	results := make(chan string, 4)

	// Launch multiple concurrent queries
	queries := []struct {
		name  string
		query func() <-chan types.AsyncResult
	}{
		{"Users", func() <-chan types.AsyncResult { 
			return querybuilder.QB().Table("users").Where("status", "active").GetAsync(ctx) 
		}},
		{"Posts", func() <-chan types.AsyncResult { 
			return querybuilder.QB().Table("posts").Where("status", "published").GetAsync(ctx) 
		}},
		{"Comments", func() <-chan types.AsyncResult { 
			return querybuilder.QB().Table("comments").Where("approved", true).GetAsync(ctx) 
		}},
		{"Categories", func() <-chan types.AsyncResult { 
			return querybuilder.QB().Table("categories").GetAsync(ctx) 
		}},
	}

	startTime := time.Now()

	for _, q := range queries {
		wg.Add(1)
		go func(name string, queryFunc func() <-chan types.AsyncResult) {
			defer wg.Done()
			resultChan := queryFunc()
			result := <-resultChan
			
			if result.Error != nil {
				results <- fmt.Sprintf("❌ %s: Error - %v", name, result.Error)
			} else {
				results <- fmt.Sprintf("✅ %s: %d records", name, result.Data.Count())
			}
		}(q.name, q.query)
	}

	// Wait for all queries to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	for result := range results {
		fmt.Printf("   %s\n", result)
	}

	duration := time.Since(startTime)
	fmt.Printf("🚀 All concurrent queries completed in %v\n", duration)

	fmt.Println()

	// Example 3: Async Pagination
	fmt.Println("3. Async Pagination")
	fmt.Println("==================")

	paginationChan := querybuilder.QB().Table("users").
		Where("age", ">=", 18).
		OrderBy("created_at", "desc").
		PaginateAsync(ctx, 1, 10)

	fmt.Println("📄 Pagination running asynchronously...")
	
	paginationResult := <-paginationChan
	if paginationResult.Error != nil {
		log.Printf("Pagination error: %v", paginationResult.Error)
	} else {
		p := paginationResult.Result
		fmt.Printf("✅ Async pagination completed:\n")
		fmt.Printf("   📊 Page %d of %d (%d total users)\n", 
			p.Meta.CurrentPage, p.Meta.LastPage, p.Meta.Total)
		fmt.Printf("   📋 Loaded %d users on this page\n", p.Count())
		
		if p.HasMorePages() {
			fmt.Printf("   ➡️  Next page available: %d\n", *p.GetNextPageNumber())
		}
	}

	fmt.Println()

	// Example 4: Async Query Racing
	fmt.Println("4. Async Query Racing (First Result Wins)")
	fmt.Println("=========================================")

	// Race multiple data sources for the fastest response
	primaryChan := querybuilder.QB().Table("users").GetAsync(ctx)
	backupChan := querybuilder.QB().Connection("backup").Table("users").GetAsync(ctx)
	
	select {
	case result := <-primaryChan:
		if result.Error == nil {
			fmt.Printf("🏆 Primary database won: %d users\n", result.Data.Count())
		} else {
			fmt.Printf("❌ Primary failed: %v, waiting for backup...\n", result.Error)
			backupResult := <-backupChan
			if backupResult.Error == nil {
				fmt.Printf("🏆 Backup database responded: %d users\n", backupResult.Data.Count())
			} else {
				fmt.Printf("❌ Both databases failed\n")
			}
		}
	case result := <-backupChan:
		if result.Error == nil {
			fmt.Printf("🏆 Backup database won: %d users\n", result.Data.Count())
		} else {
			fmt.Printf("❌ Backup failed: %v, waiting for primary...\n", result.Error)
			primaryResult := <-primaryChan
			if primaryResult.Error == nil {
				fmt.Printf("🏆 Primary database responded: %d users\n", primaryResult.Data.Count())
			} else {
				fmt.Printf("❌ Both databases failed\n")
			}
		}
	case <-time.After(5 * time.Second):
		fmt.Printf("⏰ Timeout: Both databases took too long\n")
	}

	fmt.Println()

	// Example 5: Pipeline Processing with Async
	fmt.Println("5. Async Pipeline Processing")
	fmt.Println("============================")

	// Stage 1: Fetch users
	usersChan2 := querybuilder.QB().Table("users").
		Where("status", "active").
		Limit(100).
		GetAsync(ctx)

	// Stage 2: Process users as they arrive
	go func() {
		result := <-usersChan2
		if result.Error != nil {
			log.Printf("Pipeline stage 1 error: %v", result.Error)
			return
		}

		fmt.Printf("📦 Stage 1 completed: %d users loaded\n", result.Data.Count())

		// Stage 3: Process each user asynchronously
		var processingWg sync.WaitGroup
		processedCount := 0
		
		result.Data.Each(func(user map[string]any) bool {
			processingWg.Add(1)
			go func(userID any) {
				defer processingWg.Done()
				
				// Simulate async processing (e.g., sending emails, updating records)
				time.Sleep(50 * time.Millisecond)
				processedCount++
			}(user["id"])
			
			return true
		})

		processingWg.Wait()
		fmt.Printf("📦 Stage 2 completed: %d users processed\n", processedCount)
	}()

	// Continue with other work...
	time.Sleep(2 * time.Second)

	fmt.Println()

	// Example 6: Async Batch Operations
	fmt.Println("6. Async Batch Operations")
	fmt.Println("========================")

	// Prepare batch data
	batchData := []map[string]any{
		{"name": "User 1", "email": "user1@example.com", "status": "active"},
		{"name": "User 2", "email": "user2@example.com", "status": "active"},
		{"name": "User 3", "email": "user3@example.com", "status": "pending"},
	}

	// Async batch insert
	insertChan := make(chan error, 1)
	go func() {
		defer close(insertChan)
		err := querybuilder.QB().Table("users").InsertBatch(ctx, batchData)
		insertChan <- err
	}()

	// Do other work while insert runs
	fmt.Println("📦 Batch insert running in background...")
	time.Sleep(100 * time.Millisecond)

	// Check result
	if err := <-insertChan; err != nil {
		log.Printf("Batch insert error: %v", err)
	} else {
		fmt.Printf("✅ Async batch insert completed: %d records inserted\n", len(batchData))
	}

	fmt.Println()

	// Example 7: Context Cancellation
	fmt.Println("7. Context Cancellation and Timeouts")
	fmt.Println("====================================")

	// Create context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Start a long-running query
	longQueryChan := querybuilder.QB().Table("users").
		Join("posts", "users.id", "posts.author_id").
		Join("comments", "posts.id", "comments.post_id").
		GetAsync(timeoutCtx)

	select {
	case result := <-longQueryChan:
		if result.Error != nil {
			if result.Error == context.DeadlineExceeded {
				fmt.Printf("⏰ Query timed out after 2 seconds\n")
			} else {
				fmt.Printf("❌ Query error: %v\n", result.Error)
			}
		} else {
			fmt.Printf("✅ Long query completed: %d results\n", result.Data.Count())
		}
	case <-timeoutCtx.Done():
		fmt.Printf("⏰ Context cancelled or timed out\n")
	}

	fmt.Println()

	// Example 8: Fan-out Fan-in Pattern
	fmt.Println("8. Fan-out Fan-in Pattern")
	fmt.Println("========================")

	// Fan-out: Split work across multiple goroutines
	userIDs := []int{1, 2, 3, 4, 5}
	userDetailChans := make([]<-chan types.AsyncResult, len(userIDs))

	for i, userID := range userIDs {
		userDetailChans[i] = querybuilder.QB().Table("users").
			Select("users.*", "profiles.bio", "COUNT(posts.id) as post_count").
			LeftJoin("profiles", "users.id", "profiles.user_id").
			LeftJoin("posts", "users.id", "posts.author_id").
			Where("users.id", userID).
			GroupBy("users.id", "profiles.bio").
			GetAsync(ctx)
	}

	// Fan-in: Collect all results
	allUserDetails := make([]map[string]any, 0)
	for i, ch := range userDetailChans {
		result := <-ch
		if result.Error != nil {
			fmt.Printf("❌ User %d details error: %v\n", userIDs[i], result.Error)
		} else if result.Data.Count() > 0 {
			userDetail := result.Data.First()
			allUserDetails = append(allUserDetails, userDetail)
			fmt.Printf("✅ User %d details loaded: %s\n", 
				userIDs[i], userDetail["name"])
		}
	}

	fmt.Printf("📊 Fan-out Fan-in completed: %d user details collected\n", len(allUserDetails))

	fmt.Println()
	fmt.Println("✨ Async and Performance examples complete!")
}

// Example 9: Performance Monitoring
func performanceMonitoringExample() {
	fmt.Println("9. Performance Monitoring")
	fmt.Println("========================")

	ctx := context.Background()

	// Simulate multiple queries for monitoring
	queries := []func(){
		func() { querybuilder.QB().Table("users").Get(ctx) },
		func() { querybuilder.QB().Table("posts").Where("status", "published").Get(ctx) },
		func() { querybuilder.QB().Table("comments").Count(ctx) },
		func() { querybuilder.QB().Table("users").Where("age", ">", 18).Paginate(ctx, 1, 20) },
	}

	start := time.Now()
	for _, query := range queries {
		query()
	}
	duration := time.Since(start)

	fmt.Printf("📊 Executed %d queries in %v\n", len(queries), duration)
	fmt.Printf("📈 Average query time: %v\n", duration/time.Duration(len(queries)))

	// In a real app, you would collect these metrics and send them to monitoring systems
	fmt.Printf("🔍 Performance metrics collected for monitoring dashboard\n")

	fmt.Println()
}

// Uncomment to run these examples
// func main() {
// 	asyncPerformanceExample()
// 	performanceMonitoringExample()
// }