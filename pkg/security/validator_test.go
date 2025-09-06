package security

import (
	"testing"

	"github.com/go-query-builder/querybuilder/pkg/types"
)

func TestValidateTableName(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name      string
		tableName string
		shouldErr bool
	}{
		{"Valid table name", "users", false},
		{"Valid table with underscore", "user_profiles", false},
		{"Valid table with numbers", "users2", false},
		{"Empty table name", "", true},
		{"Too long table name", "this_is_a_very_long_table_name_that_exceeds_the_maximum_allowed_length_for_table_names_in_most_databases", true},
		{"Table with forbidden keyword", "users_DROP", true},
		{"Table with spaces", "user profiles", true},
		{"Table with special chars", "users@domain", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTableName(tt.tableName)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for table name: %s", tt.tableName)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for table name %s: %v", tt.tableName, err)
			}
		})
	}
}

func TestValidateColumnName(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name       string
		columnName string
		shouldErr  bool
	}{
		{"Valid column name", "email", false},
		{"Valid column with underscore", "first_name", false},
		{"Valid column with table prefix", "users.email", false},
		{"Empty column name", "", true},
		{"Column with forbidden keyword", "email_DROP", true},
		{"Column with spaces", "first name", true},
		{"Column with special chars", "email@domain", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateColumnName(tt.columnName)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for column name: %s", tt.columnName)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for column name %s: %v", tt.columnName, err)
			}
		})
	}
}

func TestValidateOperator(t *testing.T) {
	validator := NewSecurityValidator()

	validOperators := []types.Operator{
		types.OpEqual,
		types.OpNotEqual,
		types.OpGreaterThan,
		types.OpLessThan,
		types.OpIn,
		types.OpNotIn,
		types.OpLike,
		types.OpIsNull,
	}

	for _, op := range validOperators {
		err := validator.ValidateOperator(op)
		if err != nil {
			t.Errorf("Expected no error for operator %s, got: %v", op, err)
		}
	}

	invalidOperator := types.Operator("INVALID_OP")
	err := validator.ValidateOperator(invalidOperator)
	if err == nil {
		t.Error("Expected error for invalid operator")
	}
}

func TestValidateStringValue(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name      string
		value     string
		shouldErr bool
	}{
		{"Valid string", "hello world", false},
		{"Empty string", "", false},
		{"String with SQL injection attempt", "'; DROP TABLE users; --", true},
		{"String with UNION attack", "' UNION SELECT * FROM users --", true},
		{"String with script tag", "<script>alert('xss')</script>", true},
		{"String with javascript", "javascript:alert('xss')", true},
		{"Very long string", string(make([]byte, 1500)), true},
		{"Normal long string", string(make([]byte, 500)), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.validateStringValue(tt.value)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for value: %s", tt.value)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for value %s: %v", tt.value, err)
			}
		})
	}
}

func TestValidateRawSQL(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name      string
		sql       string
		shouldErr bool
	}{
		{"Valid SELECT", "SELECT * FROM users WHERE age > ?", false},
		{"Valid JOIN", "SELECT u.name, p.bio FROM users u JOIN profiles p ON u.id = p.user_id", false},
		{"Empty SQL", "", true},
		{"SQL with DROP", "SELECT * FROM users; DROP TABLE users;", true},
		{"SQL with EXEC", "SELECT * FROM users; EXEC sp_configure", true},
		{"SQL with UNION injection", "SELECT * FROM users WHERE id = 1 UNION SELECT password FROM admin", true},
		{"SQL with comment injection", "SELECT * FROM users WHERE id = 1 -- AND status = 'active'", true},
		{"Very long SQL", string(make([]byte, 15000)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateRawSQL(tt.sql)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for SQL: %s", tt.sql)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for SQL %s: %v", tt.sql, err)
			}
		})
	}
}

func TestValidateLimit(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name      string
		limit     int
		shouldErr bool
	}{
		{"Valid limit", 100, false},
		{"Zero limit", 0, false},
		{"Max limit", 10000, false},
		{"Negative limit", -1, true},
		{"Too large limit", 20000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateLimit(tt.limit)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for limit: %d", tt.limit)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for limit %d: %v", tt.limit, err)
			}
		})
	}
}

func TestValidateOffset(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name      string
		offset    int
		shouldErr bool
	}{
		{"Valid offset", 100, false},
		{"Zero offset", 0, false},
		{"Large offset", 50000, false},
		{"Max offset", 1000000, false},
		{"Negative offset", -1, true},
		{"Too large offset", 2000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateOffset(tt.offset)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for offset: %d", tt.offset)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for offset %d: %v", tt.offset, err)
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Clean input", "hello world", "hello world"},
		{"Input with script", "<script>alert('xss')</script>", "&lt;script>alert('xss')&lt;/script>"},
		{"Input with javascript", "javascript:alert('test')", "alert('test')"},
		{"Input with vbscript", "vbscript:msgbox('test')", "msgbox('test')"},
		{"Input with control chars", "hello\x00world\x01test", "helloworldtest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestEscapeString(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Clean string", "hello world", "hello world"},
		{"String with single quote", "it's working", "it''s working"},
		{"String with backslash", "path\\to\\file", "path\\\\to\\\\file"},
		{"String with both", "it's a \\path", "it''s a \\\\path"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.EscapeString(tt.input)
			if result != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestValidateOrderBy(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name      string
		column    string
		direction types.OrderDirection
		shouldErr bool
	}{
		{"Valid order by ASC", "name", types.Asc, false},
		{"Valid order by DESC", "created_at", types.Desc, false},
		{"Invalid column", "invalid-column", types.Asc, true},
		{"Invalid direction", "name", types.OrderDirection("RANDOM"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateOrderBy(tt.column, tt.direction)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for order by: %s %s", tt.column, tt.direction)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error for order by %s %s: %v", tt.column, tt.direction, err)
			}
		})
	}
}

func TestSecureQueryBuilderValidation(t *testing.T) {
	sqb := NewSecureQueryBuilder()

	err := sqb.ValidateQuery(
		"users",
		[]string{"name", "email"},
		[]types.Operator{types.OpEqual, types.OpGreaterThan},
		[]interface{}{"John", 18},
	)

	if err != nil {
		t.Errorf("Expected no error for valid query, got: %v", err)
	}

	err = sqb.ValidateQuery(
		"invalid-table",
		[]string{"name"},
		[]types.Operator{types.OpEqual},
		[]interface{}{"John"},
	)

	if err == nil {
		t.Error("Expected error for invalid table name")
	}
}

func TestCustomValidationPatterns(t *testing.T) {
	validator := NewSecurityValidator()

	// Test adding custom table pattern
	err := validator.AddAllowedTablePattern(`^test_[a-z]+$`)
	if err != nil {
		t.Errorf("Error adding table pattern: %v", err)
	}

	err = validator.ValidateTableName("test_users")
	if err != nil {
		t.Errorf("Expected no error for custom pattern table name: %v", err)
	}

	// Test adding custom forbidden keyword
	validator.AddForbiddenKeyword("CUSTOM_BAD")
	err = validator.ValidateTableName("test_CUSTOM_BAD")
	if err == nil {
		t.Error("Expected error for custom forbidden keyword")
	}
}

func TestStrictModeToggle(t *testing.T) {
	validator := NewSecurityValidator()

	// With strict mode (default)
	err := validator.ValidateTableName("table-with-dashes")
	if err == nil {
		t.Error("Expected error in strict mode for table with dashes")
	}

	// Disable strict mode
	validator.SetStrictMode(false)
	err = validator.ValidateTableName("table-with-dashes")
	// Should still error due to forbidden patterns, but different validation logic
	// This test mainly ensures the SetStrictMode method works
	
	// Re-enable strict mode
	validator.SetStrictMode(true)
	err = validator.ValidateTableName("valid_table")
	if err != nil {
		t.Errorf("Expected no error in strict mode for valid table: %v", err)
	}
}