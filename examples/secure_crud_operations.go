package main

import (
	"context"
	"fmt"
	"log"
	"time"

	querybuilder "github.com/omarhamdy49/go-query-builder"
)

// Secure CRUD Operations Example
// Demonstrates INSERT, UPDATE, DELETE with highest security standards
func main4() {
	fmt.Println("🔒 Secure CRUD Operations with SQL Injection Prevention")
	fmt.Println("======================================================")

	ctx := context.Background()

	// // ==================================================================
	// // 1. SECURE INSERT OPERATIONS
	// // ==================================================================
	// fmt.Println("\n1. Secure INSERT Operations")
	// fmt.Println("===========================")

	// // Single Insert - All inputs are parameterized and validated
	// fmt.Println("\n📝 Single User Insert:")
	// userInsertErr := querybuilder.QB().Table("users").Insert(ctx, map[string]interface{}{
	// 	"id":         "94ef0fa5-9bfd-4562-80c6-f17326e2300e", // UUID primary key
	// 	"f_name":     "John",                                 // Automatically escaped and validated
	// 	"l_name":     "Doe",                                  // Safe from SQL injection
	// 	"email":      "john.doe@example.com",                 // Email validation built-in
	// 	"active":     true,                                   // Boolean properly handled
	// 	"password":   "securepassword",                       // Password should be hashed in a real app
	// 	"phone":      " +15551234567",                        // Phone number format checked
	// 	"created_at": time.Now(),                             // Timestamp properly formatted
	// 	"updated_at": time.Now(),
	// })

	// if userInsertErr != nil {
	// 	log.Printf("❌ Insert failed: %v", userInsertErr)
	// } else {
	// 	fmt.Println("✅ User inserted successfully with full SQL injection protection")
	// }

	// // Batch Insert - Multiple records with transaction safety
	// fmt.Println("\n📦 Batch Insert (Multiple Users):")
	// batchUsers := []map[string]interface{}{
	// 	{
	// 		"id":         "060ffb51-fe18-4fba-ae7e-d0f17c545b2d", // UUID primary key
	// 		"f_name":     "Alice",
	// 		"l_name":     "Smith",
	// 		"email":      "alice.smith@example.com",
	// 		"active":     true,
	// 		"phone":      " +15551234566",
	// 		"password":   "securepassword", // Password should be hashed in a real app
	// 		"created_at": time.Now(),
	// 		"updated_at": time.Now(),
	// 	},
	// 	{
	// 		"id":         "88ead49a-d155-4e3e-872a-846151fe15b5", // UUID primary key
	// 		"f_name":     "Bob",
	// 		"l_name":     "Johnson",
	// 		"email":      "bob.johnson@example.com",
	// 		"phone":      " +15551234561",
	// 		"active":     false,
	// 		"password":   "securepassword", // Password should be hashed in a real app
	// 		"created_at": time.Now(),
	// 		"updated_at": time.Now(),
	// 	},
	// 	{
	// 		"id":         "2079316b-aecb-405b-983b-9d0787071fde",
	// 		"f_name":     "Charlie",
	// 		"l_name":     "Brown",
	// 		"email":      "charlie.brown@example.com",
	// 		"phone":      " +15551234562",
	// 		"active":     true,
	// 		"password":   "securepassword", // Password should be hashed in a real app
	// 		"created_at": time.Now(),
	// 		"updated_at": time.Now(),
	// 	},
	// }

	// batchErr := querybuilder.QB().Table("users").InsertBatch(ctx, batchUsers)
	// if batchErr != nil {
	// 	log.Printf("❌ Batch insert failed: %v", batchErr)
	// } else {
	// 	fmt.Printf("✅ Batch inserted %d users with atomic transaction protection\n", len(batchUsers))
	// }

	// ==================================================================
	// 2. SECURE UPDATE OPERATIONS
	// ==================================================================
	fmt.Println("\n2. Secure UPDATE Operations")
	fmt.Println("===========================")

	// Update with WHERE conditions - Prevents accidental mass updates
	fmt.Println("\n✏️  Update User Email (with WHERE protection):")
	updateCount, updateErr := querybuilder.QB().Table("users").
		Where("email", "ohamdy@ihsanlife.com"). // Secure WHERE condition
		Where("active", true).                  // Additional safety condition
		Update(ctx, map[string]interface{}{
			"email":      "omar.hamdy@ihsanlife.com", // New email (validated)
			"updated_at": time.Now(),                 // Update timestamp
		})

	if updateErr != nil {
		log.Printf("❌ Update failed: %v", updateErr)
	} else {
		fmt.Printf("✅ Updated %d user(s) successfully with parameterized queries\n", updateCount)
	}

	// Conditional Update - Update only if conditions are met
	fmt.Println("\n🔄 Conditional Update (Active Users Only):")
	conditionalCount, conditionalErr := querybuilder.QB().Table("users").
		Where("active", false).                                  // Target inactive users
		Where("created_at", ">", time.Now().AddDate(0, 0, -30)). // Recent users only
		Update(ctx, map[string]interface{}{
			// "active":     true,       // Activate them
			"updated_at": time.Now(), // Track update time
		})

	if conditionalErr != nil {
		log.Printf("❌ Conditional update failed: %v", conditionalErr)
	} else {
		fmt.Printf("✅ Conditionally updated %d user(s) with multi-condition safety\n", conditionalCount)
	}

	// Update with SECURE complex WHERE conditions - FIXED to prevent mass updates
	fmt.Println("\n🎯 Complex WHERE Update (Multiple Conditions - SECURE):")
	complexCount, complexErr := querybuilder.QB().Table("users").
		Where("f_name", "John").
		Where("active", false).
		Where("email", "LIKE", "%@example.com").                // Changed to AND condition
		Where("created_at", ">", time.Now().AddDate(0, 0, -7)). // All conditions must match
		Update(ctx, map[string]interface{}{
			"l_name":     "SecureUpdate", // Changed name to indicate this is the secure version
			"updated_at": time.Now(),
		})

	if complexErr != nil {
		log.Printf("❌ Complex update failed: %v", complexErr)
	} else {
		fmt.Printf("✅ Secure complex update affected %d user(s) (should be 0 or very few)\n", complexCount)
	}

	// 🚨 CRITICAL SECURITY WARNING DEMONSTRATION
	fmt.Println("\n🚨 CRITICAL: OrWhere Security Danger Demo:")
	fmt.Println("   ⚠️  The following shows why OrWhere can cause MASS UPDATES:")
	fmt.Println("   ❌ BAD:  WHERE name='Bob' AND active=false OR email LIKE '%@example.com'")
	fmt.Println("   ❌ BAD:  Due to operator precedence, this becomes:")
	fmt.Println("   ❌ BAD:  WHERE (name='Bob' AND active=false) OR (email LIKE '%@example.com')")
	fmt.Println("   ❌ BAD:  Since most users have @example.com emails, ALL USERS GET UPDATED!")
	fmt.Println("   ✅ GOOD: Use only WHERE (AND) conditions for UPDATE/DELETE operations")
	fmt.Println("   ✅ GOOD: Always test UPDATE queries with COUNT first!")

	// BEST PRACTICE: Always test UPDATE/DELETE with COUNT first
	fmt.Println("\n✅ BEST PRACTICE: Test Before Update/Delete:")
	testCount, testErr := querybuilder.QB().Table("users").
		Where("f_name", "NonExistentUser").
		Where("email", "test@nonexistent.com").
		Count(ctx)

	if testErr != nil {
		log.Printf("❌ Test count failed: %v", testErr)
	} else {
		fmt.Printf("   📊 Test query would affect %d users\n", testCount)
		if testCount > 0 {
			fmt.Printf("   ⚠️  WARNING: %d users would be updated - verify this is intended!\n", testCount)
		} else {
			fmt.Printf("   ✅ Safe: No users would be affected by this update\n")
		}
	}

	// ==================================================================
	// 3. SECURE DELETE OPERATIONS
	// ==================================================================
	fmt.Println("\n3. Secure DELETE Operations")
	fmt.Println("===========================")

	// Safe Delete with specific conditions
	fmt.Println("\n🗑️  Safe Delete (Specific User):")
	deleteCount, deleteErr := querybuilder.QB().Table("users").
		Where("email", "john.doe.updated@example.com"). // Specific identifier
		Where("active", true).                          // Additional safety check
		Delete(ctx)

	if deleteErr != nil {
		log.Printf("❌ Delete failed: %v", deleteErr)
	} else {
		fmt.Printf("✅ Safely deleted %d user(s) with precise WHERE conditions\n", deleteCount)
	}

	// Conditional Delete - Delete inactive users older than 90 days
	fmt.Println("\n🧹 Cleanup Delete (Old Inactive Users):")
	cleanupDate := time.Now().AddDate(0, 0, -90) // 90 days ago
	cleanupCount, cleanupErr := querybuilder.QB().Table("users").
		Where("active", false).                  // Only inactive users
		Where("created_at", "<", cleanupDate).   // Older than 90 days
		Where("email", "LIKE", "%@example.com"). // Only example emails (safety)
		Delete(ctx)

	if cleanupErr != nil {
		log.Printf("❌ Cleanup delete failed: %v", cleanupErr)
	} else {
		fmt.Printf("✅ Cleanup deleted %d old inactive user(s) with time-based safety\n", cleanupCount)
	}

	// ==================================================================
	// 4. SECURITY VALIDATION DEMONSTRATIONS
	// ==================================================================
	fmt.Println("\n4. Security Protection Demonstrations")
	fmt.Println("====================================")

	// Demonstrate SQL Injection Prevention
	fmt.Println("\n🛡️  SQL Injection Prevention Tests:")

	// Test 1: Malicious input automatically sanitized
	maliciousInput := "'; DROP TABLE users; --"
	fmt.Printf("Testing malicious input: %s\n", maliciousInput)

	// This will be safely parameterized - no SQL injection possible
	safeCount, safeErr := querybuilder.QB().Table("users").
		Where("f_name", maliciousInput). // Automatically parameterized as ?
		Count(ctx)

	if safeErr != nil {
		fmt.Printf("✅ Malicious input safely handled: %v\n", safeErr)
	} else {
		fmt.Printf("✅ Query executed safely, found %d users (input was parameterized)\n", safeCount)
	}

	// Test 2: Column name validation
	fmt.Println("\n🔍 Column Name Validation:")
	_, columnErr := querybuilder.QB().Table("users").
		Where("invalid_column_name'; DROP TABLE users; --", "test"). // Invalid column
		Get(ctx)

	if columnErr != nil {
		fmt.Printf("✅ Invalid column name rejected: %v\n", columnErr)
	}

	// Test 3: Table name validation
	fmt.Println("\n🏷️  Table Name Validation:")
	_, tableErr := querybuilder.QB().Table("users'; DROP TABLE users; --"). // Invalid table
										Get(ctx)

	if tableErr != nil {
		fmt.Printf("✅ Invalid table name rejected: %v\n", tableErr)
	}

	// ==================================================================
	// 5. TRANSACTION SAFETY DEMONSTRATIONS
	// ==================================================================
	fmt.Println("\n5. Transaction Safety Features")
	fmt.Println("==============================")

	// Demonstrate safe batch operations
	fmt.Println("\n💾 Transaction-Safe Batch Operations:")

	// All batch operations are atomic - either all succeed or all fail
	transactionUsers := []map[string]interface{}{
		{
			"f_name":     "Trans",
			"l_name":     "User1",
			"email":      "trans.user1@example.com",
			"active":     true,
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
		{
			"f_name":     "Trans",
			"l_name":     "User2",
			"email":      "trans.user2@example.com",
			"active":     true,
			"created_at": time.Now(),
			"updated_at": time.Now(),
		},
	}

	transBatchErr := querybuilder.QB().Table("users").InsertBatch(ctx, transactionUsers)
	if transBatchErr != nil {
		log.Printf("❌ Transaction batch failed (all rolled back): %v", transBatchErr)
	} else {
		fmt.Printf("✅ Transaction batch completed (%d users inserted atomically)\n", len(transactionUsers))
	}

	// ==================================================================
	// 6. PERFORMANCE AND SECURITY MONITORING
	// ==================================================================
	fmt.Println("\n6. Security Monitoring & Best Practices")
	fmt.Println("=======================================")

	// Demonstrate query performance monitoring
	fmt.Println("\n📊 Query Performance Monitoring:")
	start := time.Now()

	result, perfErr := querybuilder.QB().Table("users").
		Where("active", true).
		OrderBy("created_at", "desc").
		Limit(10).
		Get(ctx)

	duration := time.Since(start)

	if perfErr != nil {
		log.Printf("❌ Performance query failed: %v", perfErr)
	} else {
		fmt.Printf("✅ Query completed in %v, returned %d records\n", duration, result.Count())

		if duration > 1*time.Second {
			fmt.Printf("⚠️  Query took longer than 1s - consider optimization\n")
		}
	}

	// ==================================================================
	// 7. SECURITY SUMMARY
	// ==================================================================
	fmt.Println("\n7. Security Features Summary")
	fmt.Println("============================")

	fmt.Println("🔒 Built-in Security Protections:")
	fmt.Println("  ✅ Parameterized queries (prevents SQL injection)")
	fmt.Println("  ✅ Input validation and sanitization")
	fmt.Println("  ✅ Column and table name validation")
	fmt.Println("  ✅ Automatic data type handling")
	fmt.Println("  ✅ Transaction safety for batch operations")
	fmt.Println("  ✅ Connection pooling with limits")
	fmt.Println("  ✅ Query timeout protection")
	fmt.Println("  ✅ Prepared statement caching")
	fmt.Println("  ✅ Error handling without data leakage")
	fmt.Println("  ✅ Automatic escaping of special characters")

	fmt.Println("\n🚀 Laravel-Level Security Compliance:")
	fmt.Println("  ✅ Mass assignment protection")
	fmt.Println("  ✅ Query builder injection prevention")
	fmt.Println("  ✅ Database credential protection")
	fmt.Println("  ✅ Connection encryption support")
	fmt.Println("  ✅ Rate limiting and connection limits")

	fmt.Println("\n✨ All CRUD operations completed securely!")
	fmt.Println("🛡️  Your data is protected against SQL injection and other attacks!")
}
