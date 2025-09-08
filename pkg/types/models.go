package types

import (
	"time"
)

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

type PaginationResult struct {
	Data Collection    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	NextPage    *int  `json:"next_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}

// Helper methods for PaginationResult
func (p PaginationResult) HasMorePages() bool {
	return p.Meta.NextPage != nil
}

func (p PaginationResult) IsEmpty() bool {
	return p.Meta.Total == 0
}

func (p PaginationResult) Count() int {
	return p.Data.Count()
}

func (p PaginationResult) OnFirstPage() bool {
	return p.Meta.CurrentPage == 1
}

func (p PaginationResult) OnLastPage() bool {
	return p.Meta.CurrentPage == p.Meta.LastPage
}

func (p PaginationResult) GetNextPageNumber() *int {
	return p.Meta.NextPage
}

func (p PaginationResult) GetPreviousPageNumber() *int {
	if p.Meta.CurrentPage > 1 {
		prev := p.Meta.CurrentPage - 1
		return &prev
	}
	return nil
}

// Query optimization configuration
type QueryOptimization struct {
	EnableQueryCache    bool          `json:"enable_query_cache"`
	CacheTTL           time.Duration `json:"cache_ttl"`
	EnablePreparedStmt bool          `json:"enable_prepared_stmt"`
	MaxConcurrency     int           `json:"max_concurrency"`
	EnableQueryLog     bool          `json:"enable_query_log"`
}

type AggregateResult struct {
	Count int64       `json:"count"`
	Sum   interface{} `json:"sum"`
	Avg   interface{} `json:"avg"`
	Min   interface{} `json:"min"`
	Max   interface{} `json:"max"`
}

type TimeHelper struct {
	Column string
	Value  time.Time
}

func NewTimeHelper(column string, value time.Time) TimeHelper {
	return TimeHelper{Column: column, Value: value}
}

type DateHelper struct {
	Column string
	Value  time.Time
}

func NewDateHelper(column string, value time.Time) DateHelper {
	return DateHelper{Column: column, Value: value}
}

type UpsertOptions struct {
	Columns        []string
	UpdateColumns  []string
	ConflictTarget []string
	ConflictAction ConflictAction
}

type BulkInsertOptions struct {
	BatchSize      int
	IgnoreErrors   bool
	OnDuplicateKey string
}

type ChunkOptions struct {
	Size    int
	OrderBy string
}

type LazyOptions struct {
	ChunkSize int
	OrderBy   string
}

type DebugInfo struct {
	SQL      string        `json:"sql"`
	Bindings []interface{} `json:"bindings"`
	Duration time.Duration `json:"duration"`
	Driver   Driver        `json:"driver"`
}

type CollectionImpl struct {
	data []map[string]interface{}
}

func NewCollection(data []map[string]interface{}) Collection {
	return &CollectionImpl{data: data}
}

func (c *CollectionImpl) ToSlice() []map[string]interface{} {
	return c.data
}

func (c *CollectionImpl) Pluck(column string) []interface{} {
	result := make([]interface{}, len(c.data))
	for i, row := range c.data {
		result[i] = row[column]
	}
	return result
}

func (c *CollectionImpl) First() map[string]interface{} {
	if len(c.data) == 0 {
		return nil
	}
	return c.data[0]
}

func (c *CollectionImpl) Count() int {
	return len(c.data)
}

func (c *CollectionImpl) IsEmpty() bool {
	return len(c.data) == 0
}

func (c *CollectionImpl) Each(fn func(map[string]interface{}) bool) {
	for _, item := range c.data {
		if !fn(item) {
			break
		}
	}
}

func (c *CollectionImpl) Filter(predicate func(map[string]interface{}) bool) Collection {
	var filtered []map[string]interface{}
	for _, item := range c.data {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return NewCollection(filtered)
}

func (c *CollectionImpl) Map(mapper func(map[string]interface{}) map[string]interface{}) Collection {
	mapped := make([]map[string]interface{}, len(c.data))
	for i, item := range c.data {
		mapped[i] = mapper(item)
	}
	return NewCollection(mapped)
}