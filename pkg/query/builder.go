package query

import (
	"context"
	"fmt"

	"github.com/go-query-builder/querybuilder/pkg/clauses"
	"github.com/go-query-builder/querybuilder/pkg/execution"
	"github.com/go-query-builder/querybuilder/pkg/types"
)

type Builder struct {
	executor    types.QueryExecutor
	driver      types.Driver
	table       string
	selects     []*clauses.SelectClause
	wheres      []*clauses.WhereClause
	joins       []*clauses.JoinClause
	orders      []*clauses.OrderClause
	groups      []*clauses.GroupClause
	havings     []*clauses.HavingClause
	unions      []*clauses.UnionClause
	limitValue  *int
	offsetValue *int
	distinct    bool
	lock        *types.LockType
	scopes      []types.ScopeFunc
	bindings    []interface{}
	compiler    *SQLCompiler
	executor_   *execution.QueryExecutor
}

func NewBuilder(executor types.QueryExecutor, driver types.Driver) *Builder {
	qb := &Builder{
		executor:    executor,
		driver:      driver,
		selects:     make([]*clauses.SelectClause, 0),
		wheres:      make([]*clauses.WhereClause, 0),
		joins:       make([]*clauses.JoinClause, 0),
		orders:      make([]*clauses.OrderClause, 0),
		groups:      make([]*clauses.GroupClause, 0),
		havings:     make([]*clauses.HavingClause, 0),
		unions:      make([]*clauses.UnionClause, 0),
		scopes:      make([]types.ScopeFunc, 0),
		bindings:    make([]interface{}, 0),
		compiler:    NewSQLCompiler(driver),
		executor_:   execution.NewQueryExecutor(executor, driver),
	}
	return qb
}

func Table(executor types.QueryExecutor, driver types.Driver, table string) types.QueryBuilder {
	qb := NewBuilder(executor, driver)
	qb.table = table
	return qb
}

func (qb *Builder) Clone() types.QueryBuilder {
	clone := &Builder{
		executor:  qb.executor,
		driver:    qb.driver,
		table:     qb.table,
		selects:   make([]*clauses.SelectClause, len(qb.selects)),
		wheres:    make([]*clauses.WhereClause, len(qb.wheres)),
		joins:     make([]*clauses.JoinClause, len(qb.joins)),
		orders:    make([]*clauses.OrderClause, len(qb.orders)),
		groups:    make([]*clauses.GroupClause, len(qb.groups)),
		havings:   make([]*clauses.HavingClause, len(qb.havings)),
		unions:    make([]*clauses.UnionClause, len(qb.unions)),
		scopes:    make([]types.ScopeFunc, len(qb.scopes)),
		bindings:  make([]interface{}, len(qb.bindings)),
		distinct:  qb.distinct,
		compiler:  NewSQLCompiler(qb.driver),
		executor_: execution.NewQueryExecutor(qb.executor, qb.driver),
	}

	copy(clone.selects, qb.selects)
	copy(clone.wheres, qb.wheres)
	copy(clone.joins, qb.joins)
	copy(clone.orders, qb.orders)
	copy(clone.groups, qb.groups)
	copy(clone.havings, qb.havings)
	copy(clone.unions, qb.unions)
	copy(clone.scopes, qb.scopes)
	copy(clone.bindings, qb.bindings)

	if qb.limitValue != nil {
		limitCopy := *qb.limitValue
		clone.limitValue = &limitCopy
	}
	if qb.offsetValue != nil {
		offsetCopy := *qb.offsetValue
		clone.offsetValue = &offsetCopy
	}
	if qb.lock != nil {
		lockCopy := *qb.lock
		clone.lock = &lockCopy
	}

	return clone
}

func (qb *Builder) From(table string) types.QueryBuilder {
	qb.table = table
	return qb
}

func (qb *Builder) Select(columns ...string) types.QueryBuilder {
	for _, column := range columns {
		qb.selects = append(qb.selects, clauses.NewSelectClause(column))
	}
	return qb
}

func (qb *Builder) SelectRaw(raw string, bindings ...interface{}) types.QueryBuilder {
	qb.selects = append(qb.selects, clauses.NewSelectRawClause(raw))
	qb.bindings = append(qb.bindings, bindings...)
	return qb
}

func (qb *Builder) SelectAs(column, alias string) types.QueryBuilder {
	qb.selects = append(qb.selects, clauses.NewSelectAsClause(column, alias))
	return qb
}

func (qb *Builder) Distinct() types.QueryBuilder {
	qb.distinct = true
	return qb
}

func (qb *Builder) Where(column string, args ...interface{}) types.QueryBuilder {
	clause := qb.parseWhereArgs(column, args...)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) OrWhere(column string, args ...interface{}) types.QueryBuilder {
	clause := qb.parseWhereArgs(column, args...)
	clause.SetBoolean(types.Or)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereNot(column string, args ...interface{}) types.QueryBuilder {
	clause := qb.parseWhereNotArgs(column, args...)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) OrWhereNot(column string, args ...interface{}) types.QueryBuilder {
	clause := qb.parseWhereNotArgs(column, args...)
	clause.SetBoolean(types.Or)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereRaw(raw string, bindings ...interface{}) types.QueryBuilder {
	clause := clauses.NewWhereRawClause(raw)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	qb.bindings = append(qb.bindings, bindings...)
	return qb
}

func (qb *Builder) OrWhereRaw(raw string, bindings ...interface{}) types.QueryBuilder {
	clause := clauses.NewWhereRawClause(raw)
	clause.SetBoolean(types.Or)
	qb.wheres = append(qb.wheres, clause)
	qb.bindings = append(qb.bindings, bindings...)
	return qb
}

func (qb *Builder) WhereBetween(column string, values []interface{}) types.QueryBuilder {
	if len(values) != 2 {
		return qb
	}
	clause := clauses.NewWhereBetweenClause(column, values, false)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereNotBetween(column string, values []interface{}) types.QueryBuilder {
	if len(values) != 2 {
		return qb
	}
	clause := clauses.NewWhereBetweenClause(column, values, true)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereIn(column string, values []interface{}) types.QueryBuilder {
	clause := clauses.NewWhereInClause(column, values, false)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereNotIn(column string, values []interface{}) types.QueryBuilder {
	clause := clauses.NewWhereInClause(column, values, true)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereNull(column string) types.QueryBuilder {
	clause := clauses.NewWhereNullClause(column, false)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereNotNull(column string) types.QueryBuilder {
	clause := clauses.NewWhereNullClause(column, true)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereExists(query types.QueryBuilder) types.QueryBuilder {
	clause := clauses.NewWhereExistsClause(query, false)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereNotExists(query types.QueryBuilder) types.QueryBuilder {
	clause := clauses.NewWhereExistsClause(query, true)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) parseWhereArgs(column string, args ...interface{}) *clauses.WhereClause {
	if len(args) == 0 {
		return clauses.NewWhereClause(column, types.OpEqual, nil)
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpEqual
		value = args[0]
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	return clauses.NewWhereClause(column, operator, value)
}

func (qb *Builder) parseWhereNotArgs(column string, args ...interface{}) *clauses.WhereClause {
	if len(args) == 0 {
		return clauses.NewWhereClause(column, types.OpNotEqual, nil)
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpNotEqual
		value = args[0]
	case 2:
		op := fmt.Sprintf("%v", args[0])
		switch op {
		case "=":
			operator = types.OpNotEqual
		case "!=":
			operator = types.OpEqual
		case ">":
			operator = types.OpLessThanOrEqual
		case ">=":
			operator = types.OpLessThan
		case "<":
			operator = types.OpGreaterThanOrEqual
		case "<=":
			operator = types.OpGreaterThan
		case "LIKE":
			operator = types.OpNotLike
		case "NOT LIKE":
			operator = types.OpLike
		default:
			operator = types.Operator(op)
		}
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	return clauses.NewWhereClause(column, operator, value)
}

func (qb *Builder) Join(table, first string, args ...interface{}) types.QueryBuilder {
	return qb.addJoin(types.InnerJoin, table, first, args...)
}

func (qb *Builder) LeftJoin(table, first string, args ...interface{}) types.QueryBuilder {
	return qb.addJoin(types.LeftJoin, table, first, args...)
}

func (qb *Builder) RightJoin(table, first string, args ...interface{}) types.QueryBuilder {
	return qb.addJoin(types.RightJoin, table, first, args...)
}

func (qb *Builder) CrossJoin(table string) types.QueryBuilder {
	qb.joins = append(qb.joins, clauses.NewCrossJoinClause(table))
	return qb
}

func (qb *Builder) addJoin(joinType types.JoinType, table, first string, args ...interface{}) types.QueryBuilder {
	var operator types.Operator
	var second string

	switch len(args) {
	case 1:
		operator = types.OpEqual
		second = fmt.Sprintf("%v", args[0])
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		second = fmt.Sprintf("%v", args[1])
	default:
		operator = types.OpEqual
		second = fmt.Sprintf("%v", args[0])
	}

	join := clauses.NewJoinClause(joinType, table, first, operator, second)
	qb.joins = append(qb.joins, join)
	return qb
}

func (qb *Builder) OrderBy(column string, direction ...types.OrderDirection) types.QueryBuilder {
	dir := types.Asc
	if len(direction) > 0 {
		dir = direction[0]
	}

	qb.orders = append(qb.orders, clauses.NewOrderClause(column, dir))
	return qb
}

func (qb *Builder) OrderByDesc(column string) types.QueryBuilder {
	return qb.OrderBy(column, types.Desc)
}

func (qb *Builder) OrderByRaw(raw string) types.QueryBuilder {
	qb.orders = append(qb.orders, clauses.NewOrderRawClause(raw))
	return qb
}

func (qb *Builder) GroupBy(columns ...string) types.QueryBuilder {
	for _, column := range columns {
		qb.groups = append(qb.groups, clauses.NewGroupClause(column))
	}
	return qb
}

func (qb *Builder) GroupByRaw(raw string) types.QueryBuilder {
	qb.groups = append(qb.groups, clauses.NewGroupRawClause(raw))
	return qb
}

func (qb *Builder) Having(column string, args ...interface{}) types.QueryBuilder {
	return qb.addHaving(column, types.And, args...)
}

func (qb *Builder) OrHaving(column string, args ...interface{}) types.QueryBuilder {
	return qb.addHaving(column, types.Or, args...)
}

func (qb *Builder) HavingRaw(raw string) types.QueryBuilder {
	clause := clauses.NewHavingRawClause(raw)
	clause.SetBoolean(types.And)
	qb.havings = append(qb.havings, clause)
	return qb
}

func (qb *Builder) OrHavingRaw(raw string) types.QueryBuilder {
	clause := clauses.NewHavingRawClause(raw)
	clause.SetBoolean(types.Or)
	qb.havings = append(qb.havings, clause)
	return qb
}

func (qb *Builder) addHaving(column string, boolean types.BooleanOperator, args ...interface{}) types.QueryBuilder {
	if len(args) == 0 {
		return qb
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpEqual
		value = args[0]
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	clause := clauses.NewHavingClause(column, operator, value)
	clause.SetBoolean(boolean)
	qb.havings = append(qb.havings, clause)
	return qb
}

func (qb *Builder) Limit(limit int) types.QueryBuilder {
	qb.limitValue = &limit
	return qb
}

func (qb *Builder) Offset(offset int) types.QueryBuilder {
	qb.offsetValue = &offset
	return qb
}

func (qb *Builder) Take(limit int) types.QueryBuilder {
	return qb.Limit(limit)
}

func (qb *Builder) Skip(offset int) types.QueryBuilder {
	return qb.Offset(offset)
}

func (qb *Builder) Union(query types.QueryBuilder) types.QueryBuilder {
	qb.unions = append(qb.unions, clauses.NewUnionClause(query))
	return qb
}

func (qb *Builder) UnionAll(query types.QueryBuilder) types.QueryBuilder {
	qb.unions = append(qb.unions, clauses.NewUnionAllClause(query))
	return qb
}

func (qb *Builder) ForUpdate() types.QueryBuilder {
	lock := types.ForUpdate
	qb.lock = &lock
	return qb
}

func (qb *Builder) ForShare() types.QueryBuilder {
	lock := types.ForShare
	qb.lock = &lock
	return qb
}

func (qb *Builder) When(condition bool, callback types.ConditionalFunc) types.QueryBuilder {
	if condition {
		return callback(qb)
	}
	return qb
}

func (qb *Builder) Unless(condition bool, callback types.ConditionalFunc) types.QueryBuilder {
	if !condition {
		return callback(qb)
	}
	return qb
}

func (qb *Builder) Tap(callback types.ConditionalFunc) types.QueryBuilder {
	callback(qb)
	return qb
}

func (qb *Builder) Scope(scopes ...types.ScopeFunc) types.QueryBuilder {
	for _, scope := range scopes {
		qb.scopes = append(qb.scopes, scope)
	}
	return qb
}

func (qb *Builder) Debug() types.QueryBuilder {
	qb.compiler.Debug()
	return qb
}

func (qb *Builder) applyScopes() {
	for _, scope := range qb.scopes {
		scope(qb)
	}
}

func (qb *Builder) ToSQL() (string, []interface{}, error) {
	qb.applyScopes()
	return qb.compiler.CompileSelect(qb)
}

func (qb *Builder) Get(ctx context.Context) (types.Collection, error) {
	return qb.executor_.Get(ctx, qb)
}

func (qb *Builder) First(ctx context.Context) (map[string]interface{}, error) {
	return qb.executor_.First(ctx, qb)
}

func (qb *Builder) Find(ctx context.Context, id interface{}) (map[string]interface{}, error) {
	return qb.executor_.Find(ctx, qb, id)
}

func (qb *Builder) Pluck(ctx context.Context, column string) ([]interface{}, error) {
	return qb.executor_.Pluck(ctx, qb, column)
}

func (qb *Builder) Count(ctx context.Context) (int64, error) {
	return qb.executor_.Count(ctx, qb)
}

func (qb *Builder) Paginate(ctx context.Context, page int, perPage int) (types.PaginationResult, error) {
	// Validate input parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15 // Default Laravel per page
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get total count using a clone to avoid affecting the original query
	countQuery := qb.Clone()
	// Remove limit and offset for count query
	countQuery.(*Builder).limitValue = nil
	countQuery.(*Builder).offsetValue = nil
	// Clear selects for count query to avoid issues with GROUP BY
	if len(qb.groups) == 0 {
		countQuery.(*Builder).selects = make([]*clauses.SelectClause, 0)
	}

	total, err := countQuery.Count(ctx)
	if err != nil {
		return types.PaginationResult{}, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get paginated data
	paginatedQuery := qb.Clone()
	paginatedQuery = paginatedQuery.Limit(perPage).Offset(offset)
	data, err := paginatedQuery.Get(ctx)
	if err != nil {
		return types.PaginationResult{}, fmt.Errorf("failed to get paginated data: %w", err)
	}

	// Calculate pagination metadata
	lastPage := int((total + int64(perPage) - 1) / int64(perPage)) // Ceiling division
	if lastPage == 0 {
		lastPage = 1
	}

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

	// Build pagination result
	result := types.PaginationResult{
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
	}

	return result, nil
}

func (qb *Builder) SimplePaginate(ctx context.Context, page int, perPage int) (types.PaginationResult, error) {
	// Validate input parameters
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 15 // Default Laravel per page
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Get one more item than requested to check if there are more pages
	checkQuery := qb.Clone()
	data, err := checkQuery.Limit(perPage + 1).Offset(offset).Get(ctx)
	if err != nil {
		return types.PaginationResult{}, fmt.Errorf("failed to get paginated data: %w", err)
	}

	// Check if there are more pages
	hasMore := data.Count() > perPage
	if hasMore {
		// Remove the extra item
		items := data.ToSlice()
		data = types.NewCollection(items[:perPage])
	}

	// Calculate metadata (without total count for performance)
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

	// Build simple pagination result (no total, no last_page)
	result := types.PaginationResult{
		Data: data,
		Meta: types.PaginationMeta{
			CurrentPage: page,
			NextPage:    nextPage,
			PerPage:     perPage,
			Total:       -1, // Indicate unknown total
			LastPage:    -1, // Indicate unknown last page
			From:        from,
			To:          to,
		},
	}

	return result, nil
}

// Async query methods
func (qb *Builder) GetAsync(ctx context.Context) <-chan types.AsyncResult {
	resultChan := make(chan types.AsyncResult, 1)
	
	go func() {
		defer close(resultChan)
		data, err := qb.Get(ctx)
		resultChan <- types.AsyncResult{Data: data, Error: err}
	}()
	
	return resultChan
}

func (qb *Builder) CountAsync(ctx context.Context) <-chan types.AsyncCountResult {
	resultChan := make(chan types.AsyncCountResult, 1)
	
	go func() {
		defer close(resultChan)
		count, err := qb.Count(ctx)
		resultChan <- types.AsyncCountResult{Count: count, Error: err}
	}()
	
	return resultChan
}

func (qb *Builder) PaginateAsync(ctx context.Context, page int, perPage int) <-chan types.AsyncPaginationResult {
	resultChan := make(chan types.AsyncPaginationResult, 1)
	
	go func() {
		defer close(resultChan)
		result, err := qb.Paginate(ctx, page, perPage)
		resultChan <- types.AsyncPaginationResult{Result: result, Error: err}
	}()
	
	return resultChan
}

func (qb *Builder) Sum(ctx context.Context, column string) (interface{}, error) {
	return qb.executor_.Sum(ctx, qb, column)
}

func (qb *Builder) Avg(ctx context.Context, column string) (interface{}, error) {
	return qb.executor_.Avg(ctx, qb, column)
}

func (qb *Builder) Min(ctx context.Context, column string) (interface{}, error) {
	return qb.executor_.Min(ctx, qb, column)
}

func (qb *Builder) Max(ctx context.Context, column string) (interface{}, error) {
	return qb.executor_.Max(ctx, qb, column)
}

func (qb *Builder) Insert(ctx context.Context, values map[string]interface{}) error {
	return qb.executor_.Insert(ctx, qb, values)
}

func (qb *Builder) InsertBatch(ctx context.Context, values []map[string]interface{}) error {
	return qb.executor_.InsertBatch(ctx, qb, values)
}

func (qb *Builder) Update(ctx context.Context, values map[string]interface{}) (int64, error) {
	return qb.executor_.Update(ctx, qb, values)
}

func (qb *Builder) Delete(ctx context.Context) (int64, error) {
	return qb.executor_.Delete(ctx, qb)
}

func (qb *Builder) GetTable() string {
	return qb.table
}

func (qb *Builder) GetSelects() []*clauses.SelectClause {
	return qb.selects
}

func (qb *Builder) GetWheres() []*clauses.WhereClause {
	return qb.wheres
}

func (qb *Builder) GetJoins() []*clauses.JoinClause {
	return qb.joins
}

func (qb *Builder) GetOrders() []*clauses.OrderClause {
	return qb.orders
}

func (qb *Builder) GetGroups() []*clauses.GroupClause {
	return qb.groups
}

func (qb *Builder) GetHavings() []*clauses.HavingClause {
	return qb.havings
}

func (qb *Builder) GetUnions() []*clauses.UnionClause {
	return qb.unions
}

func (qb *Builder) GetLimit() *int {
	return qb.limitValue
}

func (qb *Builder) GetOffset() *int {
	return qb.offsetValue
}

func (qb *Builder) IsDistinct() bool {
	return qb.distinct
}

func (qb *Builder) GetLock() *types.LockType {
	return qb.lock
}

func (qb *Builder) GetBindings() []interface{} {
	return qb.bindings
}

func (qb *Builder) GetDriver() types.Driver {
	return qb.driver
}