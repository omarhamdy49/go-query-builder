package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-query-builder/querybuilder"
)

func main() {
	ctx := context.Background()

	config := querybuilder.Config{
		Driver:   querybuilder.MySQL,
		Host:     "localhost",
		Port:     3306,
		Database: "testdb",
		Username: "root",
		Password: "password",
		Charset:  "utf8mb4",
		Timezone: "UTC",
	}

	db, err := querybuilder.NewConnection(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	basicSelects(ctx, db)
	whereClauseExamples(ctx, db)
	joinExamples(ctx, db)
	aggregateExamples(ctx, db)
	insertExamples(ctx, db)
	updateExamples(ctx, db)
	deleteExamples(ctx, db)
}

func basicSelects(ctx context.Context, db querybuilder.DB) {
	fmt.Println("=== Basic SELECT Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	collection, err := qb.Get(ctx)
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		return
	}
	fmt.Printf("Found %d users\n", collection.Count())

	user, err := qb.Where("id", 1).First(ctx)
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("User: %+v\n", user)
	}

	names, err := qb.Pluck(ctx, "name")
	if err != nil {
		log.Printf("Error plucking names: %v", err)
	} else {
		fmt.Printf("Names: %+v\n", names)
	}
}

func whereClauseExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== WHERE Clause Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	collection, err := qb.Where("age", ">", 18).
		Where("status", "active").
		OrWhere("role", "admin").
		Get(ctx)
	if err != nil {
		log.Printf("Error with where clauses: %v", err)
		return
	}
	fmt.Printf("Found %d active users or admins over 18\n", collection.Count())

	collection, err = qb.Clone().
		WhereBetween("age", []interface{}{25, 65}).
		Get(ctx)
	if err != nil {
		log.Printf("Error with between clause: %v", err)
		return
	}
	fmt.Printf("Found %d users between 25-65\n", collection.Count())

	collection, err = qb.Clone().
		WhereIn("role", []interface{}{"admin", "moderator", "user"}).
		Get(ctx)
	if err != nil {
		log.Printf("Error with in clause: %v", err)
		return
	}
	fmt.Printf("Found %d users with specified roles\n", collection.Count())

	collection, err = qb.Clone().
		WhereNull("deleted_at").
		Get(ctx)
	if err != nil {
		log.Printf("Error with null clause: %v", err)
		return
	}
	fmt.Printf("Found %d non-deleted users\n", collection.Count())
}

func joinExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== JOIN Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	collection, err := qb.Select("users.name", "profiles.bio").
		Join("profiles", "users.id", "profiles.user_id").
		Get(ctx)
	if err != nil {
		log.Printf("Error with join: %v", err)
		return
	}
	fmt.Printf("Found %d users with profiles\n", collection.Count())

	collection, err = qb.Clone().
		Select("users.name", "posts.title").
		LeftJoin("posts", "users.id", "posts.author_id").
		Get(ctx)
	if err != nil {
		log.Printf("Error with left join: %v", err)
		return
	}
	fmt.Printf("Found %d user-post combinations\n", collection.Count())
}

func aggregateExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== Aggregate Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	count, err := qb.Count(ctx)
	if err != nil {
		log.Printf("Error counting: %v", err)
	} else {
		fmt.Printf("Total users: %d\n", count)
	}

	avgAge, err := qb.Avg(ctx, "age")
	if err != nil {
		log.Printf("Error getting average age: %v", err)
	} else {
		fmt.Printf("Average age: %v\n", avgAge)
	}

	maxAge, err := qb.Max(ctx, "age")
	if err != nil {
		log.Printf("Error getting max age: %v", err)
	} else {
		fmt.Printf("Max age: %v\n", maxAge)
	}
}

func insertExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== INSERT Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	err := qb.Insert(ctx, map[string]interface{}{
		"name":       "John Doe",
		"email":      "john@example.com",
		"age":        30,
		"created_at": time.Now(),
	})
	if err != nil {
		log.Printf("Error inserting user: %v", err)
	} else {
		fmt.Println("User inserted successfully")
	}

	users := []map[string]interface{}{
		{
			"name":       "Jane Smith",
			"email":      "jane@example.com",
			"age":        28,
			"created_at": time.Now(),
		},
		{
			"name":       "Bob Johnson",
			"email":      "bob@example.com",
			"age":        35,
			"created_at": time.Now(),
		},
	}

	err = qb.InsertBatch(ctx, users)
	if err != nil {
		log.Printf("Error batch inserting users: %v", err)
	} else {
		fmt.Println("Users batch inserted successfully")
	}
}

func updateExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== UPDATE Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	rowsAffected, err := qb.Where("id", 1).Update(ctx, map[string]interface{}{
		"name":       "John Updated",
		"updated_at": time.Now(),
	})
	if err != nil {
		log.Printf("Error updating user: %v", err)
	} else {
		fmt.Printf("Updated %d rows\n", rowsAffected)
	}

	rowsAffected, err = qb.Clone().Where("age", "<", 18).Update(ctx, map[string]interface{}{
		"status":     "minor",
		"updated_at": time.Now(),
	})
	if err != nil {
		log.Printf("Error updating minors: %v", err)
	} else {
		fmt.Printf("Updated %d minor users\n", rowsAffected)
	}
}

func deleteExamples(ctx context.Context, db querybuilder.DB) {
	fmt.Println("\n=== DELETE Examples ===")

	qb := querybuilder.Table(db, db.Driver(), "users")

	rowsAffected, err := qb.Where("status", "inactive").Delete(ctx)
	if err != nil {
		log.Printf("Error deleting users: %v", err)
	} else {
		fmt.Printf("Deleted %d inactive users\n", rowsAffected)
	}
}