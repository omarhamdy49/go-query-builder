package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-query-builder/querybuilder"
	"github.com/go-query-builder/querybuilder/pkg/pagination"
	"github.com/go-query-builder/querybuilder/pkg/types"
)

func main() {
	ctx := context.Background()

	config := querybuilder.Config{
		Driver:   querybuilder.PostgreSQL,
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		Username: "postgres",
		Password: "password",
		SSLMode:  "disable",
		Timezone: "UTC",
	}

	db, err := querybuilder.NewConnection(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	dateTimeQueries(ctx, db)
	jsonQueries(ctx, db)
	fullTextSearch(ctx, db)
	upsertExamples(ctx, db)
	chunkingExamples(ctx, db)
	paginationExamples(ctx, db)
	transactionExamples(ctx, db)
	conditionalQueries(ctx, db)
	rawQueries(ctx, db)
}

func dateTimeQueries(ctx context.Context, db querybuilder.DB) {
	fmt.Println("=== Date/Time Query Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "orders")

	today := qb.Clone().WhereToday("created_at")
	sql, bindings, _ := today.ToSQL()
	fmt.Printf("Today's orders SQL: %s\nBindings: %+v\n", sql, bindings)

	thisMonth := qb.Clone().WhereMonth("created_at", time.Now().Month())
	collection, err := thisMonth.Get(ctx)
	if err != nil {
		log.Printf("Error getting this month's orders: %v", err)
	} else {
		fmt.Printf("Found %d orders this month\n", collection.Count())
	}

	pastOrders := qb.Clone().WherePast("shipped_at")
	collection, err = pastOrders.Get(ctx)
	if err != nil {
		log.Printf("Error getting past orders: %v", err)
	} else {
		fmt.Printf("Found %d past orders\n", collection.Count())
	}
}

func jsonQueries(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== JSON Query Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "products")

	jsonContains := qb.Clone().WhereJsonContains("metadata", map[string]interface{}{
		"featured": true,
	})
	collection, err := jsonContains.Get(ctx)
	if err != nil {
		log.Printf("Error with JSON contains: %v", err)
	} else {
		fmt.Printf("Found %d featured products\n", collection.Count())
	}

	jsonPath := qb.Clone().WhereJsonPath("settings", "$.notifications.email", true)
	collection, err = jsonPath.Get(ctx)
	if err != nil {
		log.Printf("Error with JSON path: %v", err)
	} else {
		fmt.Printf("Found %d products with email notifications enabled\n", collection.Count())
	}

	jsonLength := qb.Clone().WhereJsonLength("tags", ">", 3)
	collection, err = jsonLength.Get(ctx)
	if err != nil {
		log.Printf("Error with JSON length: %v", err)
	} else {
		fmt.Printf("Found %d products with more than 3 tags\n", collection.Count())
	}
}

func fullTextSearch(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Full-Text Search Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "articles")

	fullText := qb.Clone().WhereFullText([]string{"title", "content"}, "golang database")
	collection, err := fullText.Get(ctx)
	if err != nil {
		log.Printf("Error with full-text search: %v", err)
	} else {
		fmt.Printf("Found %d articles matching 'golang database'\n", collection.Count())
	}
}

func upsertExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Upsert Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "user_settings")

	values := []map[string]interface{}{
		{
			"user_id": 1,
			"key":     "theme",
			"value":   "dark",
		},
		{
			"user_id": 2,
			"key":     "language",
			"value":   "en",
		},
	}

	options := types.UpsertOptions{
		ConflictTarget: []string{"user_id", "key"},
		UpdateColumns:  []string{"value", "updated_at"},
		ConflictAction: types.DoUpdate,
	}

	// This would work with the execution package
	fmt.Printf("Upsert configuration: %+v\n", options)
	fmt.Println("Upsert values prepared for execution")
}

func chunkingExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Chunking Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "large_table")

	processedCount := 0
	chunkCallback := func(collection querybuilder.Collection) error {
		processedCount += collection.Count()
		fmt.Printf("Processed chunk of %d records (total: %d)\n", collection.Count(), processedCount)
		
		for _, item := range collection.ToSlice() {
			_ = item
		}
		return nil
	}

	fmt.Printf("Chunk callback prepared for processing\n")
	_ = chunkCallback
}

func paginationExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Pagination Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "posts")
	paginator := pagination.NewPaginator(db, db.Driver())

	page1, err := paginator.Paginate(ctx, qb.(interface {
		ToSQL() (string, []interface{}, error)
		Clone() querybuilder.QueryBuilder
		GetTable() string
		Limit(int) querybuilder.QueryBuilder
		Offset(int) querybuilder.QueryBuilder
	}), 1, 10)
	
	if err != nil {
		log.Printf("Error paginating: %v", err)
	} else {
		fmt.Printf("Page 1: %d items of %d total (%d pages)\n", 
			page1.Data.Count(), page1.Total, page1.LastPage)
	}

	simplePage, err := paginator.SimplePaginate(ctx, qb.(interface {
		ToSQL() (string, []interface{}, error)
		Clone() querybuilder.QueryBuilder
		GetTable() string
		Limit(int) querybuilder.QueryBuilder
		Offset(int) querybuilder.QueryBuilder
	}), 1, 10)
	
	if err != nil {
		log.Printf("Error with simple pagination: %v", err)
	} else {
		fmt.Printf("Simple page: %d items, has more: %t\n", 
			simplePage.Data.Count(), simplePage.HasMore)
	}
}

func transactionExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Transaction Examples ===")

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return
	}

	txQB := querybuilder.Table(tx, db.Driver(), "accounts")

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Transaction rolled back due to panic")
		}
	}()

	_, err = txQB.Where("id", 1).Update(ctx, map[string]interface{}{
		"balance": 1000,
	})
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating account: %v", err)
		return
	}

	_, err = txQB.Where("id", 2).Update(ctx, map[string]interface{}{
		"balance": 500,
	})
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating second account: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
	} else {
		fmt.Println("Transaction committed successfully")
	}
}

func conditionalQueries(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Conditional Query Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	includeInactive := false
	filterByAge := true
	minAge := 18

	finalQuery := qb.Clone().
		When(includeInactive, func(q querybuilder.QueryBuilder) querybuilder.QueryBuilder {
			return q.OrWhere("status", "inactive")
		}).
		Unless(!filterByAge, func(q querybuilder.QueryBuilder) querybuilder.QueryBuilder {
			return q.Where("age", ">=", minAge)
		}).
		Tap(func(q querybuilder.QueryBuilder) querybuilder.QueryBuilder {
			fmt.Println("Query building completed")
			return q
		})

	sql, bindings, _ := finalQuery.ToSQL()
	fmt.Printf("Conditional query SQL: %s\nBindings: %+v\n", sql, bindings)
}

func rawQueries(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Raw Query Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "analytics")

	rawSelect := qb.Clone().
		SelectRaw("DATE(created_at) as date, COUNT(*) as count").
		WhereRaw("created_at >= ?", time.Now().AddDate(0, -1, 0)).
		GroupByRaw("DATE(created_at)").
		HavingRaw("COUNT(*) > ?", 10).
		OrderByRaw("date DESC")

	sql, bindings, _ := rawSelect.ToSQL()
	fmt.Printf("Raw query SQL: %s\nBindings: %+v\n", sql, bindings)

	collection, err := rawSelect.Get(ctx)
	if err != nil {
		log.Printf("Error executing raw query: %v", err)
	} else {
		fmt.Printf("Found %d analytics records\n", collection.Count())
	}
}