// Package security provides security validation utilities for SQL queries and inputs.
package security

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// Validator provides security validation for SQL queries and database identifiers.
type Validator struct {
	strictMode            bool
	allowedTablePatterns  []*regexp.Regexp
	allowedColumnPatterns []*regexp.Regexp
	forbiddenKeywords     []string
	maxQueryLength        int
}

// NewValidator creates a new security validator with default configuration.
func NewValidator() *Validator {
	return &Validator{
		strictMode: true,
		allowedTablePatterns: []*regexp.Regexp{
			regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`),
		},
		allowedColumnPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`),
		},
		forbiddenKeywords: []string{
			"DROP", "ALTER", "CREATE", "TRUNCATE", "EXEC", "EXECUTE",
			"xp_", "sp_", "INFORMATION_SCHEMA", "SYSTEM",
		},
		maxQueryLength: 10000,
	}
}

// SetStrictMode enables or disables strict validation mode.
func (v *Validator) SetStrictMode(strict bool) *Validator {
	v.strictMode = strict
	return v
}

// AddAllowedTablePattern adds a regex pattern to the list of allowed table name patterns.
func (v *Validator) AddAllowedTablePattern(pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid table pattern: %w", err)
	}
	v.allowedTablePatterns = append(v.allowedTablePatterns, regex)
	return nil
}

// AddAllowedColumnPattern adds a regex pattern to the list of allowed column name patterns.
func (v *Validator) AddAllowedColumnPattern(pattern string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid column pattern: %w", err)
	}
	v.allowedColumnPatterns = append(v.allowedColumnPatterns, regex)
	return nil
}

// AddForbiddenKeyword adds a keyword to the list of forbidden terms.
func (v *Validator) AddForbiddenKeyword(keyword string) *Validator {
	v.forbiddenKeywords = append(v.forbiddenKeywords, strings.ToUpper(keyword))
	return v
}

// SetMaxQueryLength sets the maximum allowed length for SQL queries.
func (v *Validator) SetMaxQueryLength(length int) *Validator {
	v.maxQueryLength = length
	return v
}

// ValidateTableName validates a table name against security rules.
func (v *Validator) ValidateTableName(table string) error {
	if table == "" {
		return fmt.Errorf("table name cannot be empty")
	}

	if len(table) > 64 {
		return fmt.Errorf("table name too long: %d characters (max 64)", len(table))
	}

	for _, keyword := range v.forbiddenKeywords {
		if strings.Contains(strings.ToUpper(table), keyword) {
			return fmt.Errorf("table name contains forbidden keyword: %s", keyword)
		}
	}

	if v.strictMode {
		valid := false
		for _, pattern := range v.allowedTablePatterns {
			if pattern.MatchString(table) {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("table name does not match allowed patterns: %s", table)
		}
	}

	return nil
}

// ValidateColumnName validates a column name against security rules.
func (v *Validator) ValidateColumnName(column string) error {
	if column == "" {
		return fmt.Errorf("column name cannot be empty")
	}

	if len(column) > 64 {
		return fmt.Errorf("column name too long: %d characters (max 64)", len(column))
	}

	for _, keyword := range v.forbiddenKeywords {
		if strings.Contains(strings.ToUpper(column), keyword) {
			return fmt.Errorf("column name contains forbidden keyword: %s", keyword)
		}
	}

	if v.strictMode {
		valid := false
		for _, pattern := range v.allowedColumnPatterns {
			if pattern.MatchString(column) {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("column name does not match allowed patterns: %s", column)
		}
	}

	return nil
}

// ValidateOperator validates that an SQL operator is allowed.
func (v *Validator) ValidateOperator(operator types.Operator) error {
	allowedOperators := map[types.Operator]bool{
		types.OpEqual:              true,
		types.OpNotEqual:           true,
		types.OpGreaterThan:        true,
		types.OpGreaterThanOrEqual: true,
		types.OpLessThan:           true,
		types.OpLessThanOrEqual:    true,
		types.OpLike:               true,
		types.OpNotLike:            true,
		types.OpILike:              true,
		types.OpNotILike:           true,
		types.OpIn:                 true,
		types.OpNotIn:              true,
		types.OpBetween:            true,
		types.OpNotBetween:         true,
		types.OpIsNull:             true,
		types.OpIsNotNull:          true,
		types.OpExists:             true,
		types.OpNotExists:          true,
		types.OpJSONContains:       true,
		types.OpJSONExtract:        true,
		types.OpFullText:           true,
	}

	if !allowedOperators[operator] {
		return fmt.Errorf("operator not allowed: %s", operator)
	}

	return nil
}

// ValidateValue validates a query parameter value for security risks.
func (v *Validator) ValidateValue(value interface{}) error {
	if value == nil {
		return nil
	}

	switch val := value.(type) {
	case string:
		return v.validateStringValue(val)
	case []interface{}:
		for _, item := range val {
			if err := v.ValidateValue(item); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Validator) validateStringValue(value string) error {
	if len(value) > 1000 {
		return fmt.Errorf("string value too long: %d characters (max 1000)", len(value))
	}

	suspicious := []string{
		"<script", "javascript:", "vbscript:", "onload=", "onerror=",
		"UNION", "SELECT", "INSERT", "UPDATE", "DELETE", "DROP",
		"--", "/*", "*/", "xp_", "sp_",
	}

	upperValue := strings.ToUpper(value)
	for _, pattern := range suspicious {
		if strings.Contains(upperValue, strings.ToUpper(pattern)) {
			return fmt.Errorf("string value contains suspicious content: %s", pattern)
		}
	}

	return nil
}

// ValidateRawSQL validates raw SQL strings for injection attacks and forbidden patterns.
func (v *Validator) ValidateRawSQL(sql string) error {
	if sql == "" {
		return fmt.Errorf("raw SQL cannot be empty")
	}

	if len(sql) > v.maxQueryLength {
		return fmt.Errorf("raw SQL too long: %d characters (max %d)", len(sql), v.maxQueryLength)
	}

	for _, keyword := range v.forbiddenKeywords {
		if strings.Contains(strings.ToUpper(sql), keyword) {
			return fmt.Errorf("raw SQL contains forbidden keyword: %s", keyword)
		}
	}

	if err := v.checkForSQLInjectionPatterns(sql); err != nil {
		return err
	}

	return nil
}

func (v *Validator) checkForSQLInjectionPatterns(sql string) error {
	dangerous := []string{
		`(?i)(union\s+select)`,
		`(?i)(or\s+1\s*=\s*1)`,
		`(?i)(and\s+1\s*=\s*1)`,
		`(?i)('|\"|;|--|\*/)`,
		`(?i)(exec\s*\()`,
		`(?i)(drop\s+table)`,
		`(?i)(alter\s+table)`,
		`(?i)(create\s+table)`,
		`(?i)(truncate\s+table)`,
		`(?i)(information_schema)`,
	}

	for _, pattern := range dangerous {
		if matched, _ := regexp.MatchString(pattern, sql); matched {
			return fmt.Errorf("raw SQL contains potentially dangerous pattern")
		}
	}

	return nil
}

// SanitizeInput removes potentially dangerous characters from user input.
func (v *Validator) SanitizeInput(input string) string {
	if input == "" {
		return input
	}

	result := strings.Builder{}
	for _, r := range input {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			result.WriteRune(r)
		}
	}

	sanitized := result.String()

	sanitized = strings.ReplaceAll(sanitized, "<script", "&lt;script")
	sanitized = strings.ReplaceAll(sanitized, "javascript:", "")
	sanitized = strings.ReplaceAll(sanitized, "vbscript:", "")

	return sanitized
}

// EscapeString escapes special characters in a string to prevent SQL injection.
func (v *Validator) EscapeString(input string) string {
	if input == "" {
		return input
	}

	escaped := strings.ReplaceAll(input, "'", "''")
	escaped = strings.ReplaceAll(escaped, "\\", "\\\\")

	return escaped
}

// ValidateLimit validates that a LIMIT value is within acceptable bounds.
func (v *Validator) ValidateLimit(limit int) error {
	if limit < 0 {
		return fmt.Errorf("limit cannot be negative: %d", limit)
	}
	if limit > 10000 {
		return fmt.Errorf("limit too large: %d (max 10000)", limit)
	}
	return nil
}

// ValidateOffset validates that an OFFSET value is within acceptable bounds.
func (v *Validator) ValidateOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset cannot be negative: %d", offset)
	}
	if offset > 1000000 {
		return fmt.Errorf("offset too large: %d (max 1000000)", offset)
	}
	return nil
}

// ValidateOrderBy validates an ORDER BY clause for security issues.
func (v *Validator) ValidateOrderBy(column string, direction types.OrderDirection) error {
	if err := v.ValidateColumnName(column); err != nil {
		return fmt.Errorf("invalid order by column: %w", err)
	}

	if direction != types.Asc && direction != types.Desc {
		return fmt.Errorf("invalid order direction: %s", direction)
	}

	return nil
}

// ValidateGroupBy validates a GROUP BY column for security issues.
func (v *Validator) ValidateGroupBy(column string) error {
	return v.ValidateColumnName(column)
}

// SecureQueryBuilder provides a secure wrapper for query building with validation.
type SecureQueryBuilder struct {
	validator *Validator
}

// NewSecureQueryBuilder creates a new secure query builder with default validation settings.
func NewSecureQueryBuilder() *SecureQueryBuilder {
	return &SecureQueryBuilder{
		validator: NewValidator(),
	}
}

// GetValidator returns the underlying security validator.
func (sqb *SecureQueryBuilder) GetValidator() *Validator {
	return sqb.validator
}

// ValidateQuery performs comprehensive validation on all components of a database query.
func (sqb *SecureQueryBuilder) ValidateQuery(table string, columns []string, operators []types.Operator, values []interface{}) error {
	if err := sqb.validator.ValidateTableName(table); err != nil {
		return fmt.Errorf("table validation failed: %w", err)
	}

	for _, column := range columns {
		if err := sqb.validator.ValidateColumnName(column); err != nil {
			return fmt.Errorf("column validation failed: %w", err)
		}
	}

	for _, operator := range operators {
		if err := sqb.validator.ValidateOperator(operator); err != nil {
			return fmt.Errorf("operator validation failed: %w", err)
		}
	}

	for _, value := range values {
		if err := sqb.validator.ValidateValue(value); err != nil {
			return fmt.Errorf("value validation failed: %w", err)
		}
	}

	return nil
}
