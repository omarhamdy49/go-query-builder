package main

import (
	"context"
	"fmt"
	"log"
	"time"

	querybuilder "github.com/omarhamdy49/go-query-builder"
)

// Example showing direct table usage without any configuration
func main() {
	ctx := context.Background()

	fmt.Println("=== Direct Table Usage - Zero Configuration ===")
	fmt.Println()

	// Just table name, then start building queries!

	// 1. Simple SELECT
	fmt.Println("1. Simple queries:")
	users, err := querybuilder.QB().Table("users").Get(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Users found: %d\n", users.Count())
	}

	// 2. WHERE conditions
	fmt.Println("\n2. WHERE conditions:")
	activeUsers, err := querybuilder.QB().Table("users").
		Where("status", "active").
		Where("age", ">=", 18).
		Get(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Active adult users: %d\n", activeUsers.Count())
	}

	// 3. JOIN queries
	fmt.Println("\n3. JOIN queries:")
	userPosts, err := querybuilder.QB().Table("users").
		Select("users.name", "posts.title", "posts.created_at").
		Join("posts", "users.id", "posts.author_id").
		Where("posts.status", "published").
		OrderBy("posts.created_at", "desc").
		Limit(5).
		Get(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   User posts: %d\n", userPosts.Count())
		userPosts.Each(func(item map[string]interface{}) bool {
			fmt.Printf("     %s: %s\n", item["name"], item["title"])
			return true
		})
	}

	// 4. Aggregations
	fmt.Println("\n4. Aggregations:")
	totalPosts, err := querybuilder.QB().Table("posts").Count(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Total posts: %d\n", totalPosts)
	}

	avgAge, err := querybuilder.QB().Table("users").Avg(ctx, "age")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Average user age: %v\n", avgAge)
	}

	// 5. GROUP BY with HAVING
	fmt.Println("\n5. GROUP BY with HAVING:")
	authorStats, err := querybuilder.QB().Table("posts").
		Select("author_id", "COUNT(*) as post_count").
		GroupBy("author_id").
		Having("COUNT(*)", ">", 2).
		OrderBy("post_count", "desc").
		Get(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Prolific authors: %d\n", authorStats.Count())
	}

	// 6. INSERT
	fmt.Println("\n6. INSERT operations:")
	newUser := map[string]interface{}{
		"name":       "Alice Smith",
		"email":      "alice@example.com",
		"age":        25,
		"status":     "active",
		"role":       "user",
		"created_at": time.Now(),
	}

	err = querybuilder.QB().Table("users").Insert(ctx, newUser)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Println("   ✓ User inserted")
	}

	// 7. UPDATE
	fmt.Println("\n7. UPDATE operations:")
	affected, err := querybuilder.QB().Table("users").
		Where("email", "alice@example.com").
		Update(ctx, map[string]interface{}{
			"status":     "verified",
			"updated_at": time.Now(),
		})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   ✓ Updated %d records\n", affected)
	}

	// 8. Complex WHERE with subqueries
	fmt.Println("\n8. Complex WHERE conditions:")
	recentActiveUsers, err := querybuilder.QB().Table("users").
		Where("status", "active").
		Where("created_at", ">=", time.Now().AddDate(0, -1, 0)). // Last month
		WhereNotNull("email").
		WhereIn("role", []interface{}{"user", "admin"}).
		OrderBy("created_at", "desc").
		Get(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   Recent active users: %d\n", recentActiveUsers.Count())
	}

	// 9. Laravel-Style Pagination
	fmt.Println("\n9. Laravel-Style Pagination:")
	paginatedPosts, err := querybuilder.QB().Table("posts").
		Where("status", "published").
		OrderBy("created_at", "desc").
		Paginate(ctx, 1, 5) // page 1, 5 per page
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("   📄 Page %d of %d (%d total posts)\n",
			paginatedPosts.Meta.CurrentPage,
			paginatedPosts.Meta.LastPage,
			paginatedPosts.Meta.Total)
		fmt.Printf("   📊 Showing %d-%d of %d posts\n",
			paginatedPosts.Meta.From,
			paginatedPosts.Meta.To,
			paginatedPosts.Meta.Total)
		fmt.Printf("   🔗 Has more pages: %t\n", paginatedPosts.HasMorePages())

		// Show first few posts
		fmt.Println("   📋 Posts on this page:")
		paginatedPosts.Data.Each(func(post map[string]interface{}) bool {
			fmt.Printf("     - %v: %v\n", post["id"], post["title"])
			return true
		})
	}

	// 10. Multiple table operations in sequence
	fmt.Println("\n10. Multiple operations:")

	// Count users before
	beforeCount, _ := querybuilder.QB().Table("users").Count(ctx)

	// Insert batch
	batchUsers := []map[string]interface{}{
		{"name": "User 1", "email": "user1@test.com", "age": 20, "status": "active"},
		{"name": "User 2", "email": "user2@test.com", "age": 25, "status": "active"},
		{"name": "User 3", "email": "user3@test.com", "age": 30, "status": "pending"},
	}

	err = querybuilder.QB().Table("users").InsertBatch(ctx, batchUsers)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		// Count users after
		afterCount, _ := querybuilder.QB().Table("users").Count(ctx)
		fmt.Printf("   Users before: %d, after: %d (added %d)\n",
			beforeCount, afterCount, afterCount-beforeCount)
	}

	fmt.Println("\n✨ All operations complete with zero configuration!")
}
