# Examples Guide - Practical Usage Patterns

## Table of Contents

1. [Basic Examples](#basic-examples)
2. [CRUD Operations](#crud-operations)
3. [Advanced Queries](#advanced-queries)
4. [Pagination Examples](#pagination-examples)
5. [Async Operations](#async-operations)
6. [Security Examples](#security-examples)
7. [Performance Examples](#performance-examples)
8. [Real-World Use Cases](#real-world-use-cases)
9. [Error Handling](#error-handling)
10. [Best Practices](#best-practices)

---

## Basic Examples

### Getting Started

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/go-query-builder/querybuilder"
)

func main() {
    // Create context
    ctx := context.Background()
    
    // Get query builder instance (auto-configured from .env)
    qb := querybuilder.QB()
    
    // Simple query
    users, err := qb.Table("users").
        Where("active", true).
        Get(ctx)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d active users\n", users.Count())
    
    // Iterate through results
    users.Each(func(user map[string]any) bool {
        fmt.Printf("User ID: %v, Name: %v, Email: %v\n", 
            user["id"], user["name"], user["email"])
        return true // Continue iteration
    })
}
```

### Environment Configuration

Create `.env` file in your project root:

```env
# Database Configuration
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSL_MODE=disable

# Connection Pool Settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_MAX_LIFETIME=5m
DB_MAX_IDLE_TIME=2m
```

### Simple Queries

```go
// Count records
count, err := querybuilder.QB().Table("users").Count(ctx)

// Get first record
user, err := querybuilder.QB().Table("users").
    Where("email", "user@example.com").
    First(ctx)

// Find by ID
user, err := querybuilder.QB().Table("users").Find(ctx, 123)

// Get single value
name, err := querybuilder.QB().Table("users").
    Where("id", 123).
    Value(ctx, "name")

// Get column values
names, err := querybuilder.QB().Table("users").
    Where("active", true).
    Pluck(ctx, "name")
```

---

## CRUD Operations

### Create (INSERT)

#### Single Insert
```go
// Insert new user
err := querybuilder.QB().Table("users").Insert(ctx, map[string]interface{}{
    "name":       "John Doe",
    "email":      "john@example.com",
    "active":     true,
    "age":        30,
    "created_at": time.Now(),
})

if err != nil {
    log.Printf("Insert failed: %v", err)
    return
}

fmt.Println("User created successfully")
```

#### Batch Insert
```go
// Insert multiple users at once
users := []map[string]interface{}{
    {
        "name":       "Alice Smith",
        "email":      "alice@example.com",
        "active":     true,
        "age":        25,
        "created_at": time.Now(),
    },
    {
        "name":       "Bob Johnson",
        "email":      "bob@example.com", 
        "active":     false,
        "age":        35,
        "created_at": time.Now(),
    },
    {
        "name":       "Charlie Brown",
        "email":      "charlie@example.com",
        "active":     true,
        "age":        28,
        "created_at": time.Now(),
    },
}

err := querybuilder.QB().Table("users").InsertBatch(ctx, users)
if err != nil {
    log.Printf("Batch insert failed: %v", err)
    return
}

fmt.Printf("Successfully inserted %d users\n", len(users))
```

### Read (SELECT)

#### Basic Queries
```go
// Get all active users
activeUsers, err := querybuilder.QB().Table("users").
    Where("active", true).
    OrderBy("name", "asc").
    Get(ctx)

// Get users with specific columns
users, err := querybuilder.QB().Table("users").
    Select("id", "name", "email").
    Where("active", true).
    Get(ctx)

// Get users with conditions
adultUsers, err := querybuilder.QB().Table("users").
    Where("age", ">=", 18).
    Where("active", true).
    OrderBy("created_at", "desc").
    Limit(10).
    Get(ctx)
```

#### Complex WHERE Conditions
```go
// Multiple conditions with different operators
users, err := querybuilder.QB().Table("users").
    Where("age", ">=", 18).
    Where("age", "<=", 65).
    Where("country", "US").
    Where("active", true).
    Get(ctx)

// IN conditions
targetIDs := []interface{}{1, 2, 3, 4, 5}
users, err := querybuilder.QB().Table("users").
    WhereIn("id", targetIDs).
    Get(ctx)

// LIKE conditions for search
searchTerm := "john"
users, err := querybuilder.QB().Table("users").
    Where("name", "LIKE", "%"+searchTerm+"%").
    OrWhere("email", "LIKE", "%"+searchTerm+"%").
    Get(ctx)

// Date range queries
lastMonth := time.Now().AddDate(0, -1, 0)
recentUsers, err := querybuilder.QB().Table("users").
    Where("created_at", ">=", lastMonth).
    OrderBy("created_at", "desc").
    Get(ctx)

// NULL checks
usersWithoutProfiles, err := querybuilder.QB().Table("users").
    WhereNull("profile_completed_at").
    Get(ctx)

usersWithProfiles, err := querybuilder.QB().Table("users").
    WhereNotNull("profile_completed_at").
    Get(ctx)
```

#### JOIN Queries
```go
// INNER JOIN
usersWithProfiles, err := querybuilder.QB().Table("users").
    Join("profiles", "users.id", "profiles.user_id").
    Select("users.id", "users.name", "users.email", "profiles.bio", "profiles.avatar").
    Where("users.active", true).
    Get(ctx)

// LEFT JOIN
allUsersWithOptionalProfiles, err := querybuilder.QB().Table("users").
    LeftJoin("profiles", "users.id", "profiles.user_id").
    Select("users.*", "profiles.bio", "profiles.avatar").
    Where("users.active", true).
    Get(ctx)

// Multiple JOINs
userOrderSummary, err := querybuilder.QB().Table("users").
    LeftJoin("orders", "users.id", "orders.user_id").
    LeftJoin("order_items", "orders.id", "order_items.order_id").
    Select(
        "users.id",
        "users.name", 
        "users.email",
        "COUNT(DISTINCT orders.id) as order_count",
        "SUM(order_items.quantity * order_items.price) as total_spent",
    ).
    GroupBy("users.id", "users.name", "users.email").
    Having("COUNT(DISTINCT orders.id)", ">", 0).
    OrderBy("total_spent", "desc").
    Get(ctx)
```

### Update (UPDATE)

#### Single Record Update
```go
// Update specific user
affected, err := querybuilder.QB().Table("users").
    Where("id", 123).
    Update(ctx, map[string]interface{}{
        "name":       "John Updated",
        "active":     true,
        "updated_at": time.Now(),
    })

if err != nil {
    log.Printf("Update failed: %v", err)
    return
}

fmt.Printf("Updated %d user(s)\n", affected)
```

#### Conditional Updates
```go
// Update inactive users who haven't logged in recently
lastLoginThreshold := time.Now().AddDate(0, 0, -90) // 90 days ago

affected, err := querybuilder.QB().Table("users").
    Where("active", true).
    Where("last_login", "<", lastLoginThreshold).
    Update(ctx, map[string]interface{}{
        "active":     false,
        "status":     "dormant",
        "updated_at": time.Now(),
    })

fmt.Printf("Marked %d users as dormant\n", affected)
```

#### Update with Safety Check
```go
// Always check how many records will be affected first
count, err := querybuilder.QB().Table("users").
    Where("status", "pending").
    Count(ctx)

if err != nil {
    return err
}

fmt.Printf("About to update %d users. Continue? (y/n): ", count)
// Get user confirmation...

if count > 100 {
    return fmt.Errorf("refusing to update %d users - too many records", count)
}

// Proceed with update
affected, err := querybuilder.QB().Table("users").
    Where("status", "pending").
    Update(ctx, map[string]interface{}{
        "status":     "active",
        "activated_at": time.Now(),
        "updated_at": time.Now(),
    })
```

### Delete (DELETE)

#### Safe Deletion with Conditions
```go
// Delete inactive users older than 2 years
twoYearsAgo := time.Now().AddDate(-2, 0, 0)

deleted, err := querybuilder.QB().Table("users").
    Where("active", false).
    Where("created_at", "<", twoYearsAgo).
    Where("last_login", "<", twoYearsAgo).
    Delete(ctx)

if err != nil {
    log.Printf("Deletion failed: %v", err)
    return
}

fmt.Printf("Deleted %d old inactive users\n", deleted)
```

#### Soft Delete Pattern
```go
// Instead of actual deletion, mark as deleted
affected, err := querybuilder.QB().Table("users").
    Where("id", userID).
    Where("active", true).  // Only soft-delete active users
    Update(ctx, map[string]interface{}{
        "deleted_at": time.Now(),
        "active":     false,
        "status":     "deleted",
    })

if affected == 0 {
    return fmt.Errorf("user not found or already deleted")
}
```

---

## Advanced Queries

### Aggregation Functions

```go
// User statistics
stats := make(map[string]interface{})

// Total users
totalUsers, err := querybuilder.QB().Table("users").Count(ctx)
stats["total_users"] = totalUsers

// Active users
activeUsers, err := querybuilder.QB().Table("users").
    Where("active", true).Count(ctx)
stats["active_users"] = activeUsers

// Average age
avgAge, err := querybuilder.QB().Table("users").
    Where("age", ">", 0).Avg(ctx, "age")
stats["average_age"] = avgAge

// Oldest and youngest users
oldestAge, err := querybuilder.QB().Table("users").Max(ctx, "age")
youngestAge, err := querybuilder.QB().Table("users").Min(ctx, "age")
stats["oldest_age"] = oldestAge
stats["youngest_age"] = youngestAge

// Total orders value
totalOrderValue, err := querybuilder.QB().Table("orders").
    Where("status", "completed").Sum(ctx, "total_amount")
stats["total_order_value"] = totalOrderValue
```

### Subqueries and Complex Queries

```go
// Users who have placed orders
usersWithOrders, err := querybuilder.QB().Table("users").
    WhereExists(func(sub querybuilder.QueryBuilder) querybuilder.QueryBuilder {
        return sub.Table("orders").
            WhereRaw("orders.user_id = users.id").
            Where("status", "completed")
    }).
    Get(ctx)

// Users who haven't placed any orders
usersWithoutOrders, err := querybuilder.QB().Table("users").
    WhereNotExists(func(sub querybuilder.QueryBuilder) querybuilder.QueryBuilder {
        return sub.Table("orders").
            WhereRaw("orders.user_id = users.id")
    }).
    Get(ctx)

// Top spending customers
topCustomers, err := querybuilder.QB().Table("users").
    Join("orders", "users.id", "orders.user_id").
    Select(
        "users.id",
        "users.name",
        "users.email", 
        "COUNT(orders.id) as order_count",
        "SUM(orders.total_amount) as total_spent",
    ).
    Where("orders.status", "completed").
    GroupBy("users.id", "users.name", "users.email").
    Having("COUNT(orders.id)", ">=", 5).
    OrderBy("total_spent", "desc").
    Limit(10).
    Get(ctx)
```

### Conditional Query Building

```go
func searchUsers(ctx context.Context, filters UserFilters) (Collection, error) {
    qb := querybuilder.QB().Table("users").
        Select("id", "name", "email", "active", "created_at")
    
    // Add conditions based on filters
    if filters.Active != nil {
        qb = qb.Where("active", *filters.Active)
    }
    
    if filters.MinAge > 0 {
        qb = qb.Where("age", ">=", filters.MinAge)
    }
    
    if filters.MaxAge > 0 {
        qb = qb.Where("age", "<=", filters.MaxAge)
    }
    
    if filters.Country != "" {
        qb = qb.Where("country", filters.Country)
    }
    
    if filters.SearchTerm != "" {
        qb = qb.Where(func(q querybuilder.QueryBuilder) querybuilder.QueryBuilder {
            return q.Where("name", "LIKE", "%"+filters.SearchTerm+"%").
                OrWhere("email", "LIKE", "%"+filters.SearchTerm+"%")
        })
    }
    
    if filters.CreatedAfter != nil {
        qb = qb.Where("created_at", ">=", *filters.CreatedAfter)
    }
    
    // Apply sorting
    if filters.SortBy != "" {
        direction := "asc"
        if filters.SortDesc {
            direction = "desc"
        }
        qb = qb.OrderBy(filters.SortBy, direction)
    } else {
        qb = qb.OrderBy("created_at", "desc")
    }
    
    return qb.Get(ctx)
}

type UserFilters struct {
    Active       *bool
    MinAge       int
    MaxAge       int
    Country      string
    SearchTerm   string
    CreatedAfter *time.Time
    SortBy       string
    SortDesc     bool
}
```

---

## Pagination Examples

### Basic Pagination

```go
func getUsersPage(ctx context.Context, page int, perPage int) (*PaginationResult, error) {
    result, err := querybuilder.QB().Table("users").
        Where("active", true).
        OrderBy("created_at", "desc").
        Paginate(ctx, page, perPage)
    
    if err != nil {
        return nil, err
    }
    
    return &result, nil
}

// Usage
result, err := getUsersPage(ctx, 1, 20) // Page 1, 20 items
if err != nil {
    return err
}

fmt.Printf("Page %d of %d\n", result.Meta.CurrentPage, result.Meta.LastPage)
fmt.Printf("Showing %d-%d of %d users\n", 
    result.Meta.From, result.Meta.To, result.Meta.Total)

// Check if there are more pages
if result.HasMorePages() {
    fmt.Printf("Next page available: %d\n", *result.GetNextPageNumber())
}

// Access the data
result.Data.Each(func(user map[string]any) bool {
    fmt.Printf("User: %v (%v)\n", user["name"], user["email"])
    return true
})
```

### Pagination with Filters

```go
func searchUsersWithPagination(ctx context.Context, search UserSearchRequest) (*UserSearchResponse, error) {
    qb := querybuilder.QB().Table("users").
        Select("id", "name", "email", "active", "created_at")
    
    // Apply search filters
    if search.Query != "" {
        qb = qb.Where("name", "LIKE", "%"+search.Query+"%").
            OrWhere("email", "LIKE", "%"+search.Query+"%")
    }
    
    if search.Active != nil {
        qb = qb.Where("active", *search.Active)
    }
    
    if search.Country != "" {
        qb = qb.Where("country", search.Country)
    }
    
    // Apply sorting
    if search.SortBy != "" {
        qb = qb.OrderBy(search.SortBy, search.SortDirection)
    } else {
        qb = qb.OrderBy("created_at", "desc")
    }
    
    // Execute paginated query
    result, err := qb.Paginate(ctx, search.Page, search.PerPage)
    if err != nil {
        return nil, err
    }
    
    // Convert to response format
    users := make([]User, 0, result.Data.Count())
    result.Data.Each(func(row map[string]any) bool {
        users = append(users, User{
            ID:        row["id"].(int64),
            Name:      row["name"].(string),
            Email:     row["email"].(string),
            Active:    row["active"].(bool),
            CreatedAt: row["created_at"].(time.Time),
        })
        return true
    })
    
    return &UserSearchResponse{
        Users:      users,
        Pagination: result.Meta,
    }, nil
}

type UserSearchRequest struct {
    Query         string
    Active        *bool
    Country       string
    SortBy        string
    SortDirection string
    Page          int
    PerPage       int
}

type UserSearchResponse struct {
    Users      []User          `json:"users"`
    Pagination PaginationMeta  `json:"pagination"`
}
```

### Cursor-Based Pagination (High Performance)

```go
// For very large datasets, cursor-based pagination is more efficient
func getUsersCursor(ctx context.Context, cursor string, limit int) (*CursorResult, error) {
    qb := querybuilder.QB().Table("users").
        Select("id", "name", "email", "created_at").
        Where("active", true).
        OrderBy("id", "asc").
        Limit(limit + 1) // Get one extra to check if there are more
    
    // Apply cursor condition
    if cursor != "" {
        cursorID, err := strconv.ParseInt(cursor, 10, 64)
        if err != nil {
            return nil, fmt.Errorf("invalid cursor: %w", err)
        }
        qb = qb.Where("id", ">", cursorID)
    }
    
    users, err := qb.Get(ctx)
    if err != nil {
        return nil, err
    }
    
    // Check if there are more results
    hasMore := users.Count() > limit
    if hasMore {
        // Remove the extra item
        users = users.Take(limit)
    }
    
    // Get next cursor from last item
    var nextCursor string
    if hasMore && users.Count() > 0 {
        lastUser := users.Last()
        nextCursor = fmt.Sprintf("%v", lastUser["id"])
    }
    
    return &CursorResult{
        Data:       users,
        NextCursor: nextCursor,
        HasMore:    hasMore,
    }, nil
}

type CursorResult struct {
    Data       Collection `json:"data"`
    NextCursor string     `json:"next_cursor,omitempty"`
    HasMore    bool       `json:"has_more"`
}
```

---

## Async Operations

### Basic Async Queries

```go
// Execute multiple queries concurrently
func loadDashboardData(ctx context.Context) (*DashboardData, error) {
    // Start all queries concurrently
    userCountChan := querybuilder.QB().Table("users").CountAsync(ctx)
    activeUsersChan := querybuilder.QB().Table("users").
        Where("active", true).CountAsync(ctx)
    recentOrdersChan := querybuilder.QB().Table("orders").
        Where("created_at", ">=", time.Now().AddDate(0, 0, -7)).
        GetAsync(ctx)
    topProductsChan := querybuilder.QB().Table("products").
        OrderBy("sales_count", "desc").
        Limit(5).
        GetAsync(ctx)
    
    // Wait for results
    userCountResult := <-userCountChan
    activeUsersResult := <-activeUsersChan
    recentOrdersResult := <-recentOrdersChan
    topProductsResult := <-topProductsChan
    
    // Check for errors
    if userCountResult.Error != nil {
        return nil, fmt.Errorf("failed to get user count: %w", userCountResult.Error)
    }
    
    if activeUsersResult.Error != nil {
        return nil, fmt.Errorf("failed to get active users: %w", activeUsersResult.Error)
    }
    
    if recentOrdersResult.Error != nil {
        return nil, fmt.Errorf("failed to get recent orders: %w", recentOrdersResult.Error)
    }
    
    if topProductsResult.Error != nil {
        return nil, fmt.Errorf("failed to get top products: %w", topProductsResult.Error)
    }
    
    return &DashboardData{
        TotalUsers:    userCountResult.Count,
        ActiveUsers:   activeUsersResult.Count,
        RecentOrders:  recentOrdersResult.Data,
        TopProducts:   topProductsResult.Data,
    }, nil
}

type DashboardData struct {
    TotalUsers   int64      `json:"total_users"`
    ActiveUsers  int64      `json:"active_users"`
    RecentOrders Collection `json:"recent_orders"`
    TopProducts  Collection `json:"top_products"`
}
```

### Async with Timeout

```go
func loadDataWithTimeout(ctx context.Context, timeout time.Duration) (*Data, error) {
    // Create context with timeout
    timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    
    // Start async operations
    dataChan := querybuilder.QB().Table("large_table").
        Where("processed", false).
        GetAsync(timeoutCtx)
    
    // Wait for result or timeout
    select {
    case result := <-dataChan:
        if result.Error != nil {
            return nil, result.Error
        }
        return &Data{Items: result.Data}, nil
        
    case <-timeoutCtx.Done():
        return nil, fmt.Errorf("operation timed out after %v", timeout)
    }
}
```

### Parallel Processing with Worker Pool

```go
func processUsersInParallel(ctx context.Context) error {
    // Get all user IDs to process
    userIDs, err := querybuilder.QB().Table("users").
        Where("needs_processing", true).
        Pluck(ctx, "id")
    
    if err != nil {
        return err
    }
    
    // Create worker pool
    numWorkers := 10
    userIDChan := make(chan interface{}, len(userIDs))
    errorsChan := make(chan error, numWorkers)
    
    // Start workers
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for userID := range userIDChan {
                if err := processUser(ctx, userID); err != nil {
                    errorsChan <- err
                    return
                }
            }
        }()
    }
    
    // Send user IDs to workers
    go func() {
        defer close(userIDChan)
        for _, userID := range userIDs {
            userIDChan <- userID
        }
    }()
    
    // Wait for completion
    wg.Wait()
    close(errorsChan)
    
    // Check for errors
    for err := range errorsChan {
        return fmt.Errorf("worker error: %w", err)
    }
    
    return nil
}

func processUser(ctx context.Context, userID interface{}) error {
    // Fetch user data
    user, err := querybuilder.QB().Table("users").Find(ctx, userID)
    if err != nil {
        return err
    }
    
    // Process user (example: send email)
    if err := sendWelcomeEmail(user); err != nil {
        return err
    }
    
    // Update user status
    _, err = querybuilder.QB().Table("users").
        Where("id", userID).
        Update(ctx, map[string]interface{}{
            "needs_processing": false,
            "processed_at":     time.Now(),
        })
    
    return err
}
```

---

## Security Examples

### Input Validation

```go
func createUserSecurely(ctx context.Context, userData UserCreateRequest) error {
    // Validate input
    if err := validateUserInput(userData); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Check if email already exists
    existingCount, err := querybuilder.QB().Table("users").
        Where("email", userData.Email).
        Count(ctx)
    
    if err != nil {
        return err
    }
    
    if existingCount > 0 {
        return fmt.Errorf("email already exists")
    }
    
    // Hash password securely
    hashedPassword, err := hashPassword(userData.Password)
    if err != nil {
        return err
    }
    
    // Insert user with parameterized query (automatically secure)
    err = querybuilder.QB().Table("users").Insert(ctx, map[string]interface{}{
        "name":            userData.Name,
        "email":           userData.Email,
        "password_hash":   hashedPassword,  // Never store plain passwords
        "active":          true,
        "email_verified":  false,
        "created_at":      time.Now(),
        "updated_at":      time.Now(),
    })
    
    return err
}

func validateUserInput(userData UserCreateRequest) error {
    // Name validation
    if len(userData.Name) < 2 || len(userData.Name) > 100 {
        return errors.New("name must be between 2 and 100 characters")
    }
    
    // Email validation
    if !isValidEmail(userData.Email) {
        return errors.New("invalid email format")
    }
    
    // Password strength validation
    if !isStrongPassword(userData.Password) {
        return errors.New("password does not meet security requirements")
    }
    
    // Check for potential injection attacks
    if containsMaliciousPatterns(userData.Name) {
        return errors.New("name contains invalid characters")
    }
    
    return nil
}

func isValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    matched, _ := regexp.MatchString(pattern, email)
    return matched && len(email) <= 320 // RFC 5321 limit
}

func isStrongPassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    
    hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
    hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
    hasNumber := regexp.MustCompile(`\d`).MatchString(password)
    hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
    
    return hasUpper && hasLower && hasNumber && hasSpecial
}
```

### Safe Update Operations

```go
func updateUserSafely(ctx context.Context, userID int64, updates map[string]interface{}) error {
    // Whitelist allowed fields
    allowedFields := map[string]bool{
        "name":       true,
        "email":      true,
        "active":     true,
        "updated_at": true,
    }
    
    // Filter to only allowed fields
    safeUpdates := make(map[string]interface{})
    for field, value := range updates {
        if allowedFields[field] {
            safeUpdates[field] = value
        } else {
            log.Printf("Warning: Attempted to update disallowed field: %s", field)
        }
    }
    
    if len(safeUpdates) == 0 {
        return errors.New("no valid fields to update")
    }
    
    // Add automatic timestamp
    safeUpdates["updated_at"] = time.Now()
    
    // Validate the user exists and get current data
    currentUser, err := querybuilder.QB().Table("users").
        Where("id", userID).
        First(ctx)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return fmt.Errorf("user not found")
        }
        return err
    }
    
    // Additional security checks
    if emailUpdate, hasEmail := safeUpdates["email"]; hasEmail {
        // Check if new email is already taken
        count, err := querybuilder.QB().Table("users").
            Where("email", emailUpdate).
            Where("id", "!=", userID).
            Count(ctx)
        
        if err != nil {
            return err
        }
        
        if count > 0 {
            return fmt.Errorf("email already in use")
        }
    }
    
    // Perform the update with specific WHERE condition
    affected, err := querybuilder.QB().Table("users").
        Where("id", userID).           // Specific user
        Where("active", true).         // Only update active users
        Update(ctx, safeUpdates)
    
    if err != nil {
        return err
    }
    
    if affected == 0 {
        return fmt.Errorf("user not found or not active")
    }
    
    log.Printf("Successfully updated user %d with fields: %v", userID, getMapKeys(safeUpdates))
    return nil
}
```

---

## Performance Examples

### Efficient Batch Operations

```go
func importUsersEfficiently(ctx context.Context, users []UserImportData) error {
    batchSize := 1000
    
    for i := 0; i < len(users); i += batchSize {
        end := i + batchSize
        if end > len(users) {
            end = len(users)
        }
        
        batch := users[i:end]
        
        // Convert to map format
        batchMaps := make([]map[string]interface{}, len(batch))
        for j, user := range batch {
            batchMaps[j] = map[string]interface{}{
                "name":       user.Name,
                "email":      user.Email,
                "active":     true,
                "created_at": time.Now(),
                "updated_at": time.Now(),
            }
        }
        
        // Insert batch
        err := querybuilder.QB().Table("users").InsertBatch(ctx, batchMaps)
        if err != nil {
            return fmt.Errorf("failed to insert batch starting at %d: %w", i, err)
        }
        
        log.Printf("Imported batch %d-%d (%d users)", i+1, end, len(batch))
    }
    
    return nil
}
```

### Query Optimization Examples

```go
// ✅ OPTIMIZED: Use indexes effectively
func getUsersOptimized(ctx context.Context, filters UserFilters) (Collection, error) {
    return querybuilder.QB().Table("users").
        Select("id", "name", "email", "created_at").    // Only needed columns
        Where("active", true).                          // Indexed column first
        Where("created_at", ">=", filters.StartDate).   // Use indexed date column
        OrderBy("id", "desc").                          // Use primary key for sorting
        Limit(100).                                     // Limit results
        Get(ctx)
}

// ❌ UNOPTIMIZED: Inefficient query patterns
func getUsersUnoptimized(ctx context.Context) (Collection, error) {
    return querybuilder.QB().Table("users").
        Select("*").                                    // All columns (wasteful)
        Where("UPPER(name)", "LIKE", "%JOHN%").        // Function prevents index use
        OrderBy("RAND()").                             // Random order is expensive
        Get(ctx)                                       // No limit (could return millions)
}
```

### Caching Examples

```go
func getCachedUserStats(ctx context.Context) (*UserStats, error) {
    // Try cache first
    cacheKey := "user_stats"
    if cached := getFromCache(cacheKey); cached != nil {
        return cached.(*UserStats), nil
    }
    
    // Query database
    totalUsers, err := querybuilder.QB().Table("users").Count(ctx)
    if err != nil {
        return nil, err
    }
    
    activeUsers, err := querybuilder.QB().Table("users").
        Where("active", true).Count(ctx)
    if err != nil {
        return nil, err
    }
    
    newUsersToday, err := querybuilder.QB().Table("users").
        WhereDate("created_at", time.Now().Format("2006-01-02")).
        Count(ctx)
    if err != nil {
        return nil, err
    }
    
    stats := &UserStats{
        Total:        totalUsers,
        Active:       activeUsers,
        NewToday:     newUsersToday,
        CachedAt:     time.Now(),
    }
    
    // Cache for 5 minutes
    setCache(cacheKey, stats, 5*time.Minute)
    
    return stats, nil
}
```

---

## Real-World Use Cases

### E-commerce Order System

```go
type OrderService struct {
    db querybuilder.QueryBuilder
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    // Validate inventory
    for _, item := range req.Items {
        product, err := s.db.Table("products").Find(ctx, item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product %d not found", item.ProductID)
        }
        
        stock := product["stock"].(int64)
        if stock < int64(item.Quantity) {
            return nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
        }
    }
    
    // Calculate total
    var total float64
    for _, item := range req.Items {
        product, _ := s.db.Table("products").Find(ctx, item.ProductID)
        price := product["price"].(float64)
        total += price * float64(item.Quantity)
    }
    
    // Create order
    err := s.db.Table("orders").Insert(ctx, map[string]interface{}{
        "user_id":      req.UserID,
        "status":       "pending",
        "total_amount": total,
        "created_at":   time.Now(),
        "updated_at":   time.Now(),
    })
    
    if err != nil {
        return nil, err
    }
    
    // Get created order ID (implementation depends on your setup)
    order, err := s.db.Table("orders").
        Where("user_id", req.UserID).
        OrderBy("created_at", "desc").
        First(ctx)
    
    if err != nil {
        return nil, err
    }
    
    orderID := order["id"].(int64)
    
    // Create order items
    orderItems := make([]map[string]interface{}, len(req.Items))
    for i, item := range req.Items {
        product, _ := s.db.Table("products").Find(ctx, item.ProductID)
        
        orderItems[i] = map[string]interface{}{
            "order_id":   orderID,
            "product_id": item.ProductID,
            "quantity":   item.Quantity,
            "price":      product["price"],
            "created_at": time.Now(),
        }
    }
    
    err = s.db.Table("order_items").InsertBatch(ctx, orderItems)
    if err != nil {
        return nil, err
    }
    
    // Update inventory
    for _, item := range req.Items {
        _, err = s.db.Table("products").
            Where("id", item.ProductID).
            Update(ctx, map[string]interface{}{
                "stock": gorm.Expr("stock - ?", item.Quantity),
            })
        
        if err != nil {
            return nil, err
        }
    }
    
    return s.GetOrder(ctx, orderID)
}

func (s *OrderService) GetOrderHistory(ctx context.Context, userID int64, page int) (*PaginationResult, error) {
    return s.db.Table("orders").
        Select("orders.*", "COUNT(order_items.id) as item_count").
        LeftJoin("order_items", "orders.id", "order_items.order_id").
        Where("orders.user_id", userID).
        GroupBy("orders.id").
        OrderBy("orders.created_at", "desc").
        Paginate(ctx, page, 20)
}
```

### User Authentication System

```go
type AuthService struct {
    db querybuilder.QueryBuilder
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*User, error) {
    // Validate input
    if err := s.validateRegistration(req); err != nil {
        return nil, err
    }
    
    // Check if email exists
    count, err := s.db.Table("users").
        Where("email", req.Email).Count(ctx)
    if err != nil {
        return nil, err
    }
    if count > 0 {
        return nil, errors.New("email already registered")
    }
    
    // Hash password
    hashedPassword, err := hashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    // Create user
    err = s.db.Table("users").Insert(ctx, map[string]interface{}{
        "name":              req.Name,
        "email":             req.Email,
        "password_hash":     hashedPassword,
        "email_verified":    false,
        "active":            true,
        "registration_ip":   getClientIP(ctx),
        "created_at":        time.Now(),
        "updated_at":        time.Now(),
    })
    
    if err != nil {
        return nil, err
    }
    
    // Get created user
    user, err := s.db.Table("users").
        Where("email", req.Email).First(ctx)
    if err != nil {
        return nil, err
    }
    
    // Send verification email (async)
    go s.sendVerificationEmail(user)
    
    return s.mapToUser(user), nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
    // Get user
    user, err := s.db.Table("users").
        Where("email", email).
        Where("active", true).
        First(ctx)
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("invalid credentials")
        }
        return nil, err
    }
    
    // Verify password
    if !verifyPassword(password, user["password_hash"].(string)) {
        // Log failed attempt
        s.logLoginAttempt(ctx, email, false)
        return nil, errors.New("invalid credentials")
    }
    
    // Update last login
    _, err = s.db.Table("users").
        Where("id", user["id"]).
        Update(ctx, map[string]interface{}{
            "last_login":    time.Now(),
            "last_login_ip": getClientIP(ctx),
            "updated_at":    time.Now(),
        })
    
    if err != nil {
        log.Printf("Failed to update last login: %v", err)
    }
    
    // Log successful attempt
    s.logLoginAttempt(ctx, email, true)
    
    // Generate token
    token, err := generateJWTToken(user["id"].(int64))
    if err != nil {
        return nil, err
    }
    
    return &LoginResponse{
        Token: token,
        User:  s.mapToUser(user),
    }, nil
}

func (s *AuthService) logLoginAttempt(ctx context.Context, email string, success bool) {
    s.db.Table("login_attempts").Insert(ctx, map[string]interface{}{
        "email":      email,
        "success":    success,
        "ip_address": getClientIP(ctx),
        "user_agent": getUserAgent(ctx),
        "created_at": time.Now(),
    })
}
```

---

## Error Handling

### Comprehensive Error Handling

```go
func handleUserOperation(ctx context.Context, userID int64) error {
    user, err := querybuilder.QB().Table("users").Find(ctx, userID)
    if err != nil {
        switch {
        case err == sql.ErrNoRows:
            return &NotFoundError{
                Resource: "user",
                ID:       userID,
            }
        case isConnectionError(err):
            return &DatabaseConnectionError{
                Operation: "find_user",
                Cause:     err,
            }
        case isTimeoutError(err):
            return &TimeoutError{
                Operation: "find_user",
                Duration:  30 * time.Second,
                Cause:     err,
            }
        default:
            return &DatabaseError{
                Operation: "find_user",
                Cause:     err,
            }
        }
    }
    
    // Process user...
    return nil
}

// Custom error types
type NotFoundError struct {
    Resource string
    ID       int64
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

type DatabaseConnectionError struct {
    Operation string
    Cause     error
}

func (e *DatabaseConnectionError) Error() string {
    return fmt.Sprintf("database connection failed during %s: %v", e.Operation, e.Cause)
}

type TimeoutError struct {
    Operation string
    Duration  time.Duration
    Cause     error
}

func (e *TimeoutError) Error() string {
    return fmt.Sprintf("operation %s timed out after %v: %v", e.Operation, e.Duration, e.Cause)
}
```

### Retry Logic

```go
func executeWithRetry(ctx context.Context, operation func() error, maxRetries int) error {
    var lastErr error
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil // Success
        }
        
        lastErr = err
        
        // Don't retry on certain errors
        if !isRetryableError(err) {
            return err
        }
        
        if attempt < maxRetries {
            // Exponential backoff
            backoff := time.Duration(attempt*attempt) * time.Second
            log.Printf("Attempt %d failed: %v. Retrying in %v...", attempt+1, err, backoff)
            
            select {
            case <-time.After(backoff):
                continue
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", maxRetries+1, lastErr)
}

func isRetryableError(err error) bool {
    // Retry on connection errors, timeouts, temporary failures
    if isConnectionError(err) || isTimeoutError(err) {
        return true
    }
    
    // Don't retry on validation errors, not found, etc.
    if isValidationError(err) || isNotFoundError(err) {
        return false
    }
    
    return true
}

// Usage
err := executeWithRetry(ctx, func() error {
    return querybuilder.QB().Table("users").Insert(ctx, userData)
}, 3)
```

---

Your Go Query Builder provides comprehensive functionality with **Laravel-compatible syntax**, **enterprise-grade security**, and **high-performance operations** for all your database needs.

🚀 **Ready for production use with real-world examples and best practices!** 🚀