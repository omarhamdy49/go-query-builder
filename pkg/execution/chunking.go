package execution

import (
	"context"
	"fmt"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

func (e *QueryExecutor) Chunk(ctx context.Context, qb QueryBuilderInterface, size int, callback types.ChunkFunc) error {
	if size <= 0 {
		return fmt.Errorf("chunk size must be positive")
	}

	offset := 0
	orderBy := "id"

	for {
		clone := qb.Clone()
		chunkQB := clone.OrderBy(orderBy).Limit(size).Offset(offset)
		
		collection, err := e.Get(ctx, chunkQB.(QueryBuilderInterface))
		if err != nil {
			return fmt.Errorf("failed to get chunk: %w", err)
		}

		if collection.IsEmpty() {
			break
		}

		if err := callback(collection); err != nil {
			return fmt.Errorf("chunk callback error: %w", err)
		}

		if collection.Count() < size {
			break
		}

		offset += size
	}

	return nil
}

func (e *QueryExecutor) ChunkById(ctx context.Context, qb QueryBuilderInterface, size int, callback types.ChunkFunc, column ...string) error {
	if size <= 0 {
		return fmt.Errorf("chunk size must be positive")
	}

	idColumn := "id"
	if len(column) > 0 && column[0] != "" {
		idColumn = column[0]
	}

	var lastId interface{}

	for {
		clone := qb.Clone()
		chunkQB := clone.OrderBy(idColumn).Limit(size)
		
		if lastId != nil {
			chunkQB = chunkQB.Where(idColumn, ">", lastId)
		}

		collection, err := e.Get(ctx, chunkQB.(QueryBuilderInterface))
		if err != nil {
			return fmt.Errorf("failed to get chunk by id: %w", err)
		}

		if collection.IsEmpty() {
			break
		}

		if err := callback(collection); err != nil {
			return fmt.Errorf("chunk callback error: %w", err)
		}

		if collection.Count() < size {
			break
		}

		lastItem := collection.ToSlice()[collection.Count()-1]
		lastId = lastItem[idColumn]
	}

	return nil
}

func (e *QueryExecutor) Each(ctx context.Context, qb QueryBuilderInterface, callback types.LazyFunc, chunkSize ...int) error {
	size := 1000
	if len(chunkSize) > 0 && chunkSize[0] > 0 {
		size = chunkSize[0]
	}

	return e.Chunk(ctx, qb, size, func(collection types.Collection) error {
		for _, item := range collection.ToSlice() {
			if err := callback(item); err != nil {
				return err
			}
		}
		return nil
	})
}

func (e *QueryExecutor) EachById(ctx context.Context, qb QueryBuilderInterface, callback types.LazyFunc, chunkSize ...int) error {
	size := 1000
	if len(chunkSize) > 0 && chunkSize[0] > 0 {
		size = chunkSize[0]
	}

	return e.ChunkById(ctx, qb, size, func(collection types.Collection) error {
		for _, item := range collection.ToSlice() {
			if err := callback(item); err != nil {
				return err
			}
		}
		return nil
	})
}

func (e *QueryExecutor) Lazy(ctx context.Context, qb QueryBuilderInterface, chunkSize ...int) (*LazyCollection, error) {
	size := 1000
	if len(chunkSize) > 0 && chunkSize[0] > 0 {
		size = chunkSize[0]
	}

	return &LazyCollection{
		executor: e,
		qb:       qb,
		ctx:      ctx,
		size:     size,
		offset:   0,
		orderBy:  "id",
	}, nil
}

func (e *QueryExecutor) LazyById(ctx context.Context, qb QueryBuilderInterface, column string, chunkSize ...int) (*LazyCollection, error) {
	size := 1000
	if len(chunkSize) > 0 && chunkSize[0] > 0 {
		size = chunkSize[0]
	}

	idColumn := "id"
	if column != "" {
		idColumn = column
	}

	return &LazyCollection{
		executor:     e,
		qb:           qb,
		ctx:          ctx,
		size:         size,
		orderBy:      idColumn,
		useIdCursor:  true,
		idColumn:     idColumn,
	}, nil
}

type LazyCollection struct {
	executor     *QueryExecutor
	qb           QueryBuilderInterface
	ctx          context.Context
	size         int
	offset       int
	orderBy      string
	useIdCursor  bool
	idColumn     string
	lastId       interface{}
	currentBatch types.Collection
	batchIndex   int
	finished     bool
}

func (lc *LazyCollection) Next() bool {
	if lc.finished {
		return false
	}

	if lc.currentBatch == nil || lc.batchIndex >= lc.currentBatch.Count() {
		if err := lc.loadNextBatch(); err != nil {
			lc.finished = true
			return false
		}
		
		if lc.currentBatch.IsEmpty() {
			lc.finished = true
			return false
		}
		
		lc.batchIndex = 0
	}

	return !lc.finished
}

func (lc *LazyCollection) Value() map[string]interface{} {
	if lc.currentBatch == nil || lc.batchIndex >= lc.currentBatch.Count() {
		return nil
	}

	item := lc.currentBatch.ToSlice()[lc.batchIndex]
	lc.batchIndex++
	return item
}

func (lc *LazyCollection) Each(callback types.LazyFunc) error {
	for lc.Next() {
		if err := callback(lc.Value()); err != nil {
			return err
		}
	}
	return nil
}

func (lc *LazyCollection) loadNextBatch() error {
	clone := lc.qb.Clone()
	
	if lc.useIdCursor {
		batchQB := clone.OrderBy(lc.idColumn).Limit(lc.size)
		if lc.lastId != nil {
			batchQB = batchQB.Where(lc.idColumn, ">", lc.lastId)
		}
		
		collection, err := lc.executor.Get(lc.ctx, batchQB.(QueryBuilderInterface))
		if err != nil {
			return err
		}
		
		lc.currentBatch = collection
		
		if !collection.IsEmpty() {
			lastItem := collection.ToSlice()[collection.Count()-1]
			lc.lastId = lastItem[lc.idColumn]
		}
	} else {
		batchQB := clone.OrderBy(lc.orderBy).Limit(lc.size).Offset(lc.offset)
		
		collection, err := lc.executor.Get(lc.ctx, batchQB.(QueryBuilderInterface))
		if err != nil {
			return err
		}
		
		lc.currentBatch = collection
		lc.offset += lc.size
	}

	return nil
}

func (lc *LazyCollection) ToSlice() ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	
	for lc.Next() {
		result = append(result, lc.Value())
	}
	
	return result, nil
}

func (lc *LazyCollection) Filter(predicate func(map[string]interface{}) bool) *LazyCollection {
	// Create a new lazy collection that applies the filter
	return &LazyCollection{
		executor:     lc.executor,
		qb:           lc.qb,
		ctx:          lc.ctx,
		size:         lc.size,
		orderBy:      lc.orderBy,
		useIdCursor:  lc.useIdCursor,
		idColumn:     lc.idColumn,
		// Note: In a full implementation, you'd need to handle filtering properly
		// For now, we'll return the original collection
	}
}

func (lc *LazyCollection) Map(mapper func(map[string]interface{}) map[string]interface{}) *LazyCollection {
	// Create a new lazy collection that applies the mapper
	return &LazyCollection{
		executor:     lc.executor,
		qb:           lc.qb,
		ctx:          lc.ctx,
		size:         lc.size,
		orderBy:      lc.orderBy,
		useIdCursor:  lc.useIdCursor,
		idColumn:     lc.idColumn,
		// Note: In a full implementation, you'd need to handle mapping properly
		// For now, we'll return the original collection
	}
}

func (e *QueryExecutor) Cursor(ctx context.Context, qb QueryBuilderInterface) (*Cursor, error) {
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := e.executor.QueryContext(ctx, sql, bindings...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	columns, err := rows.Columns()
	if err != nil {
		rows.Close()
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	return &Cursor{
		rows:    rows,
		columns: columns,
	}, nil
}

type Cursor struct {
	rows    types.Rows
	columns []string
	closed  bool
}

func (c *Cursor) Next() bool {
	if c.closed {
		return false
	}
	return c.rows.Next()
}

func (c *Cursor) Scan(dest ...interface{}) error {
	return c.rows.Scan(dest...)
}

func (c *Cursor) ScanStruct(dest interface{}) error {
	return fmt.Errorf("struct scanning not implemented")
}

func (c *Cursor) ScanMap() (map[string]interface{}, error) {
	values := make([]interface{}, len(c.columns))
	valuePtrs := make([]interface{}, len(c.columns))
	
	for i := range c.columns {
		valuePtrs[i] = &values[i]
	}

	if err := c.rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, col := range c.columns {
		val := values[i]
		if b, ok := val.([]byte); ok {
			result[col] = string(b)
		} else {
			result[col] = val
		}
	}

	return result, nil
}

func (c *Cursor) Close() error {
	if !c.closed {
		c.closed = true
		return c.rows.Close()
	}
	return nil
}

func (c *Cursor) Err() error {
	return c.rows.Err()
}

func (c *Cursor) Columns() []string {
	return c.columns
}