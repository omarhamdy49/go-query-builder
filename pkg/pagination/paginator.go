// Package pagination provides utilities for paginating database query results.
// It includes support for Laravel-style pagination with data and meta structures,
// simple pagination for performance optimization, and cursor-based pagination
// for large datasets.
package pagination

import (
	"context"
	"math"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

type Paginator struct {
	executor types.QueryExecutor
	driver   types.Driver
}

func NewPaginator(executor types.QueryExecutor, driver types.Driver) *Paginator {
	return &Paginator{
		executor: executor,
		driver:   driver,
	}
}

type QueryBuilderInterface interface {
	ToSQL() (string, []any, error)
	Clone() types.QueryBuilder
	GetTable() string
	Limit(int) types.QueryBuilder
	Offset(int) types.QueryBuilder
}

func (p *Paginator) Paginate(ctx context.Context, qb QueryBuilderInterface, page, perPage int) (*types.PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	countQB := qb.Clone()
	total, err := p.getTotal(ctx, countQB.(QueryBuilderInterface))
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * perPage
	dataQB := qb.Clone().Limit(perPage).Offset(offset)
	
	data, err := p.getData(ctx, dataQB.(QueryBuilderInterface))
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	from := offset + 1
	to := offset + data.Count()
	
	if total == 0 {
		from = 0
		to = 0
	}

	var nextPage *int
	if page < lastPage {
		next := page + 1
		nextPage = &next
	}

	return &types.PaginationResult{
		Data: data,
		Meta: types.PaginationMeta{
			CurrentPage: page,
			NextPage:    nextPage,
			PerPage:     perPage,
			Total:       total,
			LastPage:    lastPage,
			From:        from,
			To:          to,
		},
	}, nil
}

func (p *Paginator) SimplePaginate(ctx context.Context, qb QueryBuilderInterface, page, perPage int) (*types.PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15
	}

	offset := (page - 1) * perPage
	dataQB := qb.Clone().Limit(perPage + 1).Offset(offset)
	
	data, err := p.getData(ctx, dataQB.(QueryBuilderInterface))
	if err != nil {
		return nil, err
	}

	hasMore := data.Count() > perPage
	if hasMore {
		items := data.ToSlice()
		data = types.NewCollection(items[:perPage])
	}

	from := offset + 1
	to := offset + data.Count()
	
	if data.Count() == 0 {
		from = 0
		to = 0
	}

	var nextPage *int
	if hasMore {
		next := page + 1
		nextPage = &next
	}

	return &types.PaginationResult{
		Data: data,
		Meta: types.PaginationMeta{
			CurrentPage: page,
			NextPage:    nextPage,
			PerPage:     perPage,
			Total:       -1, // Unknown for simple pagination
			LastPage:    -1, // Unknown for simple pagination
			From:        from,
			To:          to,
		},
	}, nil
}

func (p *Paginator) CursorPaginate(ctx context.Context, qb QueryBuilderInterface, cursor string, perPage int, cursorColumn ...string) (*CursorPaginationResult, error) {
	if perPage < 1 {
		perPage = 15
	}

	column := "id"
	if len(cursorColumn) > 0 && cursorColumn[0] != "" {
		column = cursorColumn[0]
	}

	dataQB := qb.Clone().OrderBy(column).Limit(perPage + 1)
	
	if cursor != "" {
		dataQB = dataQB.Where(column, ">", cursor)
	}

	data, err := p.getData(ctx, dataQB.(QueryBuilderInterface))
	if err != nil {
		return nil, err
	}

	hasMore := data.Count() > perPage
	var nextCursor string
	
	if hasMore {
		items := data.ToSlice()
		data = types.NewCollection(items[:perPage])
		nextCursor = items[perPage-1][column].(string)
	} else if data.Count() > 0 {
		items := data.ToSlice()
		nextCursor = items[data.Count()-1][column].(string)
	}

	return &CursorPaginationResult{
		Data:       data,
		PerPage:    perPage,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (p *Paginator) getTotal(ctx context.Context, qb QueryBuilderInterface) (int64, error) {
	countQB := qb.Clone()
	
	sql, bindings, err := countQB.ToSQL()
	if err != nil {
		return 0, err
	}

	countSQL := p.wrapCountQuery(sql)
	
	row := p.executor.QueryRowContext(ctx, countSQL, bindings...)
	
	var total int64
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (p *Paginator) getData(ctx context.Context, qb QueryBuilderInterface) (types.Collection, error) {
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := p.executor.QueryContext(ctx, sql, bindings...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return p.scanRows(rows)
}

func (p *Paginator) scanRows(rows types.Rows) (types.Collection, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]any
	
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return types.NewCollection(results), nil
}

func (p *Paginator) wrapCountQuery(sql string) string {
	return "SELECT COUNT(*) FROM (" + sql + ") AS count_query"
}

type SimplePaginationResult struct {
	Data        types.Collection `json:"data"`
	PerPage     int              `json:"per_page"`
	CurrentPage int              `json:"current_page"`
	From        int              `json:"from"`
	To          int              `json:"to"`
	HasMore     bool             `json:"has_more"`
}

type CursorPaginationResult struct {
	Data       types.Collection `json:"data"`
	PerPage    int              `json:"per_page"`
	NextCursor string           `json:"next_cursor"`
	HasMore    bool             `json:"has_more"`
}

func (r *SimplePaginationResult) HasPages() bool {
	return r.HasMore
}

func (r *SimplePaginationResult) IsEmpty() bool {
	return r.Data.IsEmpty()
}

func (r *CursorPaginationResult) IsEmpty() bool {
	return r.Data.IsEmpty()
}