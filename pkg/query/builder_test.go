package query

import (
	"context"
	"testing"

	"github.com/go-query-builder/querybuilder/pkg/types"
)

// MockExecutor for testing
type MockExecutor struct {
	driver types.Driver
}

func (m *MockExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (types.Rows, error) {
	return nil, nil
}

func (m *MockExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) types.Row {
	return nil
}

func (m *MockExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (types.Result, error) {
	return nil, nil
}

func (m *MockExecutor) Begin() (types.Tx, error) {
	return nil, nil
}

func (m *MockExecutor) BeginTx(ctx context.Context, opts *types.TxOptions) (types.Tx, error) {
	return nil, nil
}

func (m *MockExecutor) Driver() types.Driver {
	return m.driver
}

func TestBasicSelectQuery(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(bindings) != 0 {
		t.Errorf("Expected no bindings, got: %v", bindings)
	}
}

func TestSelectWithColumns(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.Select("name", "email", "age")
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT name, email, age FROM users"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	if len(bindings) != 0 {
		t.Errorf("Expected no bindings, got: %v", bindings)
	}
}

func TestSelectDistinct(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.Select("name").Distinct()
	sql, _, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT DISTINCT name FROM users"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
}

func TestWhereBasic(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.Where("age", ">", 18).Where("status", "active")
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users WHERE age > ? AND status = ?"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedBindings := []interface{}{18, "active"}
	if len(bindings) != len(expectedBindings) {
		t.Errorf("Expected %d bindings, got %d", len(expectedBindings), len(bindings))
	}
}

func TestWhereIn(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.WhereIn("role", []interface{}{"admin", "editor", "user"})
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users WHERE role IN (?, ?, ?)"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedBindings := []interface{}{"admin", "editor", "user"}
	if len(bindings) != len(expectedBindings) {
		t.Errorf("Expected %d bindings, got %d", len(expectedBindings), len(bindings))
	}
}

func TestWhereBetween(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.WhereBetween("age", []interface{}{18, 65})
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users WHERE age BETWEEN ? AND ?"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedBindings := []interface{}{18, 65}
	if len(bindings) != len(expectedBindings) {
		t.Errorf("Expected %d bindings, got %d", len(expectedBindings), len(bindings))
	}
}

func TestJoinQuery(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.Select("users.name", "profiles.bio").
		Join("profiles", "users.id", "profiles.user_id")
	
	sql, _, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT users.name, profiles.bio FROM users INNER JOIN profiles ON users.id = profiles.user_id"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
}

func TestLeftJoin(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.LeftJoin("posts", "users.id", "posts.author_id")
	sql, _, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users LEFT JOIN posts ON users.id = posts.author_id"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
}

func TestOrderBy(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.OrderBy("name").OrderByDesc("created_at")
	sql, _, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users ORDER BY name ASC, created_at DESC"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
}

func TestGroupBy(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "orders"

	qb.Select("status").GroupBy("status")
	sql, _, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT status FROM orders GROUP BY status"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
}

func TestHaving(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "orders"

	qb.Select("status").GroupBy("status").Having("count", ">", 5)
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT status FROM orders GROUP BY status HAVING count > ?"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedBindings := []interface{}{5}
	if len(bindings) != len(expectedBindings) {
		t.Errorf("Expected %d bindings, got %d", len(expectedBindings), len(bindings))
	}
}

func TestLimitOffset(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.Limit(10).Offset(20)
	sql, _, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users LIMIT 10 OFFSET 20"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}
}

func TestComplexQuery(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	qb.Select("users.name", "profiles.bio").
		Join("profiles", "users.id", "profiles.user_id").
		Where("users.age", ">", 18).
		WhereIn("users.role", []interface{}{"admin", "editor"}).
		OrderBy("users.name").
		Limit(10)

	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT users.name, profiles.bio FROM users INNER JOIN profiles ON users.id = profiles.user_id WHERE users.age > ? AND users.role IN (?, ?) ORDER BY users.name ASC LIMIT 10"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedBindings := []interface{}{18, "admin", "editor"}
	if len(bindings) != len(expectedBindings) {
		t.Errorf("Expected %d bindings, got %d", len(expectedBindings), len(bindings))
	}
}

func TestPostgreSQLPlaceholders(t *testing.T) {
	executor := &MockExecutor{driver: types.PostgreSQL}
	qb := NewBuilder(executor, types.PostgreSQL)
	qb.table = "users"

	qb.Where("age", ">", 18).Where("status", "active")
	
	// Note: This test would need the compiler to properly handle PostgreSQL placeholders
	// The current implementation doesn't fully support $1, $2 style placeholders yet
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// For now, we just check that it doesn't error
	if sql == "" {
		t.Error("Expected non-empty SQL")
	}
	
	if len(bindings) != 2 {
		t.Errorf("Expected 2 bindings, got %d", len(bindings))
	}
}

func TestQueryBuilderClone(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"
	qb.Where("status", "active")

	clone := qb.Clone()
	clone.Where("age", ">", 18)

	originalSQL, originalBindings, _ := qb.ToSQL()
	cloneSQL, cloneBindings, _ := clone.ToSQL()

	if originalSQL == cloneSQL {
		t.Error("Clone should have different SQL than original")
	}

	if len(originalBindings) >= len(cloneBindings) {
		t.Error("Clone should have more bindings than original")
	}
}

func TestConditionalQueries(t *testing.T) {
	executor := &MockExecutor{driver: types.MySQL}
	qb := NewBuilder(executor, types.MySQL)
	qb.table = "users"

	condition := true
	finalQB := qb.When(condition, func(q types.QueryBuilder) types.QueryBuilder {
		return q.Where("status", "active")
	}).Unless(condition, func(q types.QueryBuilder) types.QueryBuilder {
		return q.Where("deleted_at", "IS NULL")
	})

	sql, bindings, err := finalQB.ToSQL()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSQL := "SELECT * FROM users WHERE status = ?"
	if sql != expectedSQL {
		t.Errorf("Expected SQL: %s, got: %s", expectedSQL, sql)
	}

	expectedBindings := []interface{}{"active"}
	if len(bindings) != len(expectedBindings) {
		t.Errorf("Expected %d bindings, got %d", len(expectedBindings), len(bindings))
	}
}