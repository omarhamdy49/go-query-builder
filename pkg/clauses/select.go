package clauses

// SelectClause represents a SELECT clause in an SQL query.
type SelectClause struct {
	Column string
	Alias  string
	Raw    string
}

// NewSelectClause creates a new SELECT clause for the specified column.
func NewSelectClause(column string) *SelectClause {
	return &SelectClause{
		Column: column,
	}
}

// NewSelectAsClause creates a new SELECT clause with an alias.
func NewSelectAsClause(column, alias string) *SelectClause {
	return &SelectClause{
		Column: column,
		Alias:  alias,
	}
}

// NewSelectRawClause creates a new SELECT clause using raw SQL.
func NewSelectRawClause(raw string) *SelectClause {
	return &SelectClause{
		Raw: raw,
	}
}

// IsRaw returns true if this is a raw SELECT clause.
func (s *SelectClause) IsRaw() bool {
	return s.Raw != ""
}

// HasAlias returns true if this SELECT clause has an alias.
func (s *SelectClause) HasAlias() bool {
	return s.Alias != ""
}

// GetColumn returns the column name for this SELECT clause.
func (s *SelectClause) GetColumn() string {
	return s.Column
}

// GetAlias returns the alias for this SELECT clause.
func (s *SelectClause) GetAlias() string {
	return s.Alias
}

// GetRaw returns the raw SQL for this SELECT clause.
func (s *SelectClause) GetRaw() string {
	return s.Raw
}