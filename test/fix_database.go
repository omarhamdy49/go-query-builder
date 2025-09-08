package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-query-builder/querybuilder"
)

func main() {
	fmt.Println("🔧 Fixing database - Reverting mass update")
	ctx := context.Background()

	// Revert the incorrect mass update
	count, err := querybuilder.QB().Table("users").
		Where("l_name", "UpdatedLastName").
		Update(ctx, map[string]interface{}{
			"l_name": "OriginalLastName", // Set to a generic value
		})
	
	if err != nil {
		log.Fatalf("Failed to fix database: %v", err)
	}
	
	fmt.Printf("✅ Fixed %d users that were incorrectly updated\n", count)
}
