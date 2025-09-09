// Package types provides type definitions, interfaces, and constants used throughout the go-query-builder.
// It includes model structures, query builder interfaces, operators, and other shared types.
package types

import (
	"time"
)

// Config holds database connection configuration.
type Config struct {
	Driver          Driver        `json:"driver"`
	Host            string        `json:"host"`
	Port            int           `json:"port"`
	Database        string        `json:"database"`
	Username        string        `json:"username"`
	Password        string        `json:"password"`
	SSLMode         string        `json:"ssl_mode"`
	Charset         string        `json:"charset"`
	Timezone        string        `json:"timezone"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"conn_max_idle_time"`
}

// PaginationResult represents the result of a paginated query.
type PaginationResult struct {
	Data Collection     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// PaginationMeta contains pagination metadata.
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	NextPage    *int  `json:"next_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}

// HasMorePages returns true if there are more pages available.
func (p PaginationResult) HasMorePages() bool {
	return p.Meta.NextPage != nil
}

// IsEmpty returns true if the pagination result contains no data.
func (p PaginationResult) IsEmpty() bool {
	return p.Meta.Total == 0
}

// Count returns the number of items in the current page.
func (p PaginationResult) Count() int {
	return p.Data.Count()
}

// OnFirstPage returns true if currently on the first page.
func (p PaginationResult) OnFirstPage() bool {
	return p.Meta.CurrentPage == 1
}

// OnLastPage returns true if currently on the last page.
func (p PaginationResult) OnLastPage() bool {
	return p.Meta.CurrentPage == p.Meta.LastPage
}

// GetNextPageNumber returns the next page number if available.
func (p PaginationResult) GetNextPageNumber() *int {
	return p.Meta.NextPage
}

// GetPreviousPageNumber returns the previous page number if available.
func (p PaginationResult) GetPreviousPageNumber() *int {
	if p.Meta.CurrentPage > 1 {
		prev := p.Meta.CurrentPage - 1

		return &prev
	}

	return nil
}

// QueryOptimization contains configuration for query optimization features.
type QueryOptimization struct {
	EnableQueryCache   bool          `json:"enable_query_cache"`
	CacheTTL           time.Duration `json:"cache_ttl"`
	EnablePreparedStmt bool          `json:"enable_prepared_stmt"`
	MaxConcurrency     int           `json:"max_concurrency"`
	EnableQueryLog     bool          `json:"enable_query_log"`
}

// AggregateResult holds the results of aggregate functions.
type AggregateResult struct {
	Count int64       `json:"count"`
	Sum   interface{} `json:"sum"`
	Avg   interface{} `json:"avg"`
	Min   interface{} `json:"min"`
	Max   interface{} `json:"max"`
}

// TimeHelper assists with time-based queries.
type TimeHelper struct {
	Column string
	Value  time.Time
}

// NewTimeHelper creates a new TimeHelper instance.
func NewTimeHelper(column string, value time.Time) TimeHelper {
	return TimeHelper{Column: column, Value: value}
}

// DateHelper assists with date-based queries.
type DateHelper struct {
	Column string
	Value  time.Time
}

// NewDateHelper creates a new DateHelper instance.
func NewDateHelper(column string, value time.Time) DateHelper {
	return DateHelper{Column: column, Value: value}
}

// UpsertOptions configures upsert operations.
type UpsertOptions struct {
	Columns        []string
	UpdateColumns  []string
	ConflictTarget []string
	ConflictAction ConflictAction
}

// BulkInsertOptions configures bulk insert operations.
type BulkInsertOptions struct {
	BatchSize      int
	IgnoreErrors   bool
	OnDuplicateKey string
}

// ChunkOptions configures chunk processing.
type ChunkOptions struct {
	Size    int
	OrderBy string
}

// LazyOptions configures lazy loading.
type LazyOptions struct {
	ChunkSize int
	OrderBy   string
}

// DebugInfo contains debugging information for queries.
type DebugInfo struct {
	SQL      string        `json:"sql"`
	Bindings []interface{} `json:"bindings"`
	Duration time.Duration `json:"duration"`
	Driver   Driver        `json:"driver"`
}

// CollectionImpl implements the Collection interface.
type CollectionImpl struct {
	data []map[string]interface{}
}

// NewCollection creates a new Collection from slice data.
func NewCollection(data []map[string]interface{}) Collection {
	return &CollectionImpl{data: data}
}

// ToSlice returns the underlying data slice.
func (c *CollectionImpl) ToSlice() []map[string]interface{} {
	return c.data
}

// Pluck extracts a column from all rows.
func (c *CollectionImpl) Pluck(column string) []interface{} {
	result := make([]interface{}, len(c.data))
	for i, row := range c.data {
		result[i] = row[column]
	}

	return result
}

// First returns the first item in the collection.
func (c *CollectionImpl) First() map[string]interface{} {
	if len(c.data) == 0 {
		return nil
	}

	return c.data[0]
}

// Count returns the number of items in the collection.
func (c *CollectionImpl) Count() int {
	return len(c.data)
}

// IsEmpty returns true if the collection has no items.
func (c *CollectionImpl) IsEmpty() bool {
	return len(c.data) == 0
}

// Each iterates over each item in the collection.
func (c *CollectionImpl) Each(fn func(map[string]interface{}) bool) {
	for _, item := range c.data {
		if !fn(item) {
			break
		}
	}
}

// Filter returns a new collection with items that match the predicate.
func (c *CollectionImpl) Filter(predicate func(map[string]interface{}) bool) Collection {
	var filtered []map[string]interface{}

	for _, item := range c.data {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}

	return NewCollection(filtered)
}

// Map returns a new collection with each item transformed by the mapper function.
func (c *CollectionImpl) Map(mapper func(map[string]interface{}) map[string]interface{}) Collection {
	mapped := make([]map[string]interface{}, len(c.data))
	for i, item := range c.data {
		mapped[i] = mapper(item)
	}

	return NewCollection(mapped)
}
