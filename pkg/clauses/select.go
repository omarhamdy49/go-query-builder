package clauses

type SelectClause struct {
	Column string
	Alias  string
	Raw    string
}

func NewSelectClause(column string) *SelectClause {
	return &SelectClause{
		Column: column,
	}
}

func NewSelectAsClause(column, alias string) *SelectClause {
	return &SelectClause{
		Column: column,
		Alias:  alias,
	}
}

func NewSelectRawClause(raw string) *SelectClause {
	return &SelectClause{
		Raw: raw,
	}
}

func (s *SelectClause) IsRaw() bool {
	return s.Raw != ""
}

func (s *SelectClause) HasAlias() bool {
	return s.Alias != ""
}

func (s *SelectClause) GetColumn() string {
	return s.Column
}

func (s *SelectClause) GetAlias() string {
	return s.Alias
}

func (s *SelectClause) GetRaw() string {
	return s.Raw
}