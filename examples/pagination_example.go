package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-query-builder/querybuilder"
)

func paginationExample() {
	ctx := context.Background()

	fmt.Println("=== Laravel-Style Pagination Examples ===")
	fmt.Println()

	// Example 1: Basic Pagination
	fmt.Println("1. Basic Pagination")
	fmt.Println("==================")

	result, err := querybuilder.QB().Table("users").Paginate(ctx, 1, 10)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("📄 Page %d of %d\n", result.Meta.CurrentPage, result.Meta.LastPage)
	fmt.Printf("📊 Showing %d-%d of %d users\n", result.Meta.From, result.Meta.To, result.Meta.Total)
	fmt.Printf("🔗 Has more pages: %t\n", result.HasMorePages())
	fmt.Printf("📱 Items per page: %d\n", result.Meta.PerPage)

	// Show pagination helper methods
	fmt.Printf("🏠 On first page: %t\n", result.OnFirstPage())
	fmt.Printf("🔚 On last page: %t\n", result.OnLastPage())
	if next := result.GetNextPageNumber(); next != nil {
		fmt.Printf("➡️  Next page: %d\n", *next)
	}
	if prev := result.GetPreviousPageNumber(); prev != nil {
		fmt.Printf("⬅️  Previous page: %d\n", *prev)
	}

	fmt.Println("\n👥 Users on this page:")
	result.Data.Each(func(user map[string]any) bool {
		fmt.Printf("   - %v: %v (%v)\n", user["id"], user["f_name"], user["email"])
		return true
	})

	// Example 2: Pagination with Filters
	fmt.Println("\n\n2. Pagination with WHERE Filters")
	fmt.Println("===============================")

	filtered, err := querybuilder.QB().Table("users").
		Where("active", false).
		OrderBy("f_name", "asc").
		Paginate(ctx, 1, 5)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("📊 Active adult users: %d total, showing page %d\n",
		filtered.Meta.Total, filtered.Meta.CurrentPage)
	fmt.Printf("📱 %d per page (%d-%d of %d)\n",
		filtered.Meta.PerPage, filtered.Meta.From, filtered.Meta.To, filtered.Meta.Total)

	// Example 3: Pagination with JOINs
	fmt.Println("\n\n3. Pagination with JOIN Queries")
	fmt.Println("==============================")

	joined, err := querybuilder.QB().Table("users").
		Select("users.id", "users.name", "users.email", "COUNT(posts.id) as post_count").
		LeftJoin("posts", "users.id", "posts.author_id").
		Where("users.status", "active").
		GroupBy("users.id", "users.name", "users.email").
		Having("COUNT(posts.id)", ">", 0).
		OrderBy("post_count", "desc").
		Paginate(ctx, 1, 8)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("📊 Users with posts: %d total\n", joined.Meta.Total)
	fmt.Printf("📄 Page %d of %d\n", joined.Meta.CurrentPage, joined.Meta.LastPage)

	fmt.Println("\n👥 Top contributors:")
	joined.Data.Each(func(user map[string]interface{}) bool {
		fmt.Printf("   - %v (%v): %v posts\n",
			user["name"], user["email"], user["post_count"])
		return true
	})

	// Example 4: JSON Response Format
	fmt.Println("\n\n4. JSON Response Format")
	fmt.Println("======================")

	jsonResult, err := querybuilder.QB().Table("posts").
		Select("id", "title", "status", "created_at").
		Where("status", "published").
		OrderBy("created_at", "desc").
		Paginate(ctx, 1, 3)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Convert to JSON (like an API response)
	jsonBytes, err := json.MarshalIndent(jsonResult, "", "  ")
	if err != nil {
		log.Printf("JSON Error: %v", err)
		return
	}

	fmt.Println("📄 API Response JSON:")
	fmt.Println(string(jsonBytes))

	// Example 5: Different Page Sizes
	fmt.Println("\n\n5. Different Page Sizes")
	fmt.Println("======================")

	pageSizes := []int{5, 10, 20}
	for _, pageSize := range pageSizes {
		sizedResult, err := querybuilder.QB().Table("users").
			Paginate(ctx, 1, pageSize)

		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}

		fmt.Printf("📱 Page size %d: %d pages total, showing %d items\n",
			pageSize, sizedResult.Meta.LastPage, sizedResult.Count())
	}

	// Example 6: Navigation Through Pages
	fmt.Println("\n\n6. Page Navigation Example")
	fmt.Println("=========================")

	// Simulate browsing through pages
	for page := 1; page <= 3; page++ {
		navResult, err := querybuilder.QB().Table("posts").
			Where("status", "published").
			OrderBy("id", "asc").
			Paginate(ctx, page, 4)

		if err != nil {
			log.Printf("Error on page %d: %v", page, err)
			continue
		}

		fmt.Printf("\n📄 Page %d/%d:\n", navResult.Meta.CurrentPage, navResult.Meta.LastPage)
		fmt.Printf("   Range: %d-%d of %d total\n",
			navResult.Meta.From, navResult.Meta.To, navResult.Meta.Total)

		// Navigation info
		if navResult.OnFirstPage() {
			fmt.Printf("   🏠 This is the first page\n")
		}
		if navResult.OnLastPage() {
			fmt.Printf("   🔚 This is the last page\n")
		}
		if navResult.HasMorePages() {
			fmt.Printf("   ➡️  More pages available\n")
		}

		// Show items
		fmt.Printf("   📋 Items:\n")
		navResult.Data.Each(func(post map[string]interface{}) bool {
			fmt.Printf("     - Post %v: %v\n", post["id"], post["title"])
			return true
		})

		// Break if no more pages
		if navResult.OnLastPage() {
			break
		}
	}

	// Example 7: Empty Results
	fmt.Println("\n\n7. Empty Results Handling")
	fmt.Println("========================")

	emptyResult, err := querybuilder.QB().Table("users").
		Where("status", "nonexistent_status").
		Paginate(ctx, 1, 10)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if emptyResult.IsEmpty() {
		fmt.Println("📭 No results found")
		fmt.Printf("   Total: %d, Current page: %d\n",
			emptyResult.Meta.Total, emptyResult.Meta.CurrentPage)
		fmt.Printf("   Range: %d-%d\n", emptyResult.Meta.From, emptyResult.Meta.To)
	}

	// Example 8: Advanced Search with Pagination
	fmt.Println("\n\n8. Advanced Search + Pagination")
	fmt.Println("==============================")

	searchTerm := "test"
	searchResult, err := querybuilder.QB().Table("posts").
		Select("id", "title", "content", "created_at").
		Where("title", "LIKE", fmt.Sprintf("%%%s%%", searchTerm)).
		OrWhere("content", "LIKE", fmt.Sprintf("%%%s%%", searchTerm)).
		Where("status", "published").
		OrderBy("created_at", "desc").
		Paginate(ctx, 1, 5)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("🔍 Search for '%s': %d results\n", searchTerm, searchResult.Meta.Total)
	fmt.Printf("📄 Page %d of %d (%d items)\n",
		searchResult.Meta.CurrentPage, searchResult.Meta.LastPage, searchResult.Count())

	if !searchResult.IsEmpty() {
		fmt.Println("\n📋 Search results:")
		searchResult.Data.Each(func(post map[string]interface{}) bool {
			fmt.Printf("   - %v: %v\n", post["id"], post["title"])
			return true
		})
	}

	fmt.Println("\n✨ Pagination examples complete!")
}

// Uncomment to run this example
// func main() {
// 	paginationExample()
// }
