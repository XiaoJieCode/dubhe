package clause

const (
	OpEq      = "="
	OpNEq     = "<>"
	OpGt      = ">"
	OpGte     = ">="
	OpLt      = "<"
	OpLte     = "<="
	OpIn      = "IN"
	OpLike    = "LIKE"
	OpNull    = "IS NULL"
	OpNotNull = "IS NOT NULL"
	OpOr      = "OR"
)

type Clause struct {
	Field string
	Value any
	Op    string
}
type Match struct {
	Clauses []Clause
}

func NewMatch() *Match {
	return &Match{}
}

func (m *Match) add(field, op string, value any) *Match {
	m.Clauses = append(m.Clauses, Clause{Field: field, Op: op, Value: value})
	return m
}

// Eq 相等条件
func (m *Match) Eq(field string, value any) *Match {
	return m.add(field, OpEq, value)
}

// NEq 不等条件
func (m *Match) NEq(field string, value any) *Match {
	return m.add(field, OpNEq, value)
}

// Gt 大于条件
func (m *Match) Gt(field string, value any) *Match {
	return m.add(field, OpGt, value)
}

// Gte 大于等于条件
func (m *Match) Gte(field string, value any) *Match {
	return m.add(field, OpGte, value)
}

// Lt 小于条件
func (m *Match) Lt(field string, value any) *Match {
	return m.add(field, OpLt, value)
}

// Lte 小于等于条件
func (m *Match) Lte(field string, value any) *Match {
	return m.add(field, OpLte, value)
}

// In 包含条件
func (m *Match) In(field string, value any) *Match {
	return m.add(field, OpIn, value)
}

// Like 模糊匹配
func (m *Match) Like(field string, value any) *Match {
	return m.add(field, OpLike, value)
}

// Null 字段为 NULL
func (m *Match) Null(field string) *Match {
	// Null操作没有value
	return m.add(field, OpNull, nil)
}

// NotNull 字段不为 NULL
func (m *Match) NotNull(field string) *Match {
	// NotNull操作没有value
	return m.add(field, OpNotNull, nil)
}

// Or 或条件
func (m *Match) Or(field string, value any) *Match {
	return m.add(field, OpOr, value)
}
