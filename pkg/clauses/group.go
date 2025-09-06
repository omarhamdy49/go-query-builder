package clauses

type GroupClause struct {
	Column string
	Raw    string
}

func NewGroupClause(column string) *GroupClause {
	return &GroupClause{
		Column: column,
	}
}

func NewGroupRawClause(raw string) *GroupClause {
	return &GroupClause{
		Raw: raw,
	}
}

func (g *GroupClause) IsRaw() bool {
	return g.Raw != ""
}

func (g *GroupClause) GetColumn() string {
	return g.Column
}

func (g *GroupClause) GetRaw() string {
	return g.Raw
}