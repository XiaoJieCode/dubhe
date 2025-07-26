package condition

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
	OpAsc     = "ASC"
	OpDesc    = "DESC"
	OpSet     = "="
)

type Clause struct {
	Field string
	Value any
	Op    string
}
type Match struct {
	Clauses []Clause
	Orders  []Clause
	Sets    []Clause
}

func NewMatch() *Match {
	return &Match{}
}

func (m *Match) add(field, op string, value any) *Match {
	m.Clauses = append(m.Clauses, Clause{Field: field, Op: op, Value: value})
	return m
}
func (m *Match) order(field, op string) *Match {
	m.Orders = append(m.Orders, Clause{Field: field, Op: op})
	return m
}
func (m *Match) Set(field string, value any) *Match {
	m.Sets = append(m.Sets, Clause{Field: field, Value: value, Op: OpSet})
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

func (m *Match) Asc(field string) *Match {
	m.order(field, OpAsc)
	return m
}
func (m *Match) Desc(field string) *Match {
	m.order(field, OpDesc)
	return m
}
func (m *Match) WhereSql() (string, []any) {
	if len(m.Clauses) == 0 {
		return "", nil
	}

	sql := ""
	args := []any{}
	for i, c := range m.Clauses {
		// 拼接连接符，默认AND，遇到OpOr用OR
		if i > 0 {
			if c.Op == OpOr {
				sql += " OR "
				continue
			} else {
				sql += " AND "
			}
		}

		switch c.Op {
		case OpEq, OpNEq, OpGt, OpGte, OpLt, OpLte, OpLike:
			sql += c.Field + " " + c.Op + " ?"
			args = append(args, c.Value)
		case OpIn:
			// 断言value是切片
			valSlice, ok := toSlice(c.Value)
			if !ok || len(valSlice) == 0 {
				// 这里直接跳过这个条件，不拼接SQL，不加参数
				// 但跳过后连接符处理复杂一点，需要特殊判断。
				// 简单处理：去掉前面的连接符（AND/OR）
				// 方案：因为无法回退，建议先把要拼的条件缓存起来，最后再拼接。
				continue
			}
			sql += c.Field + " IN ("
			for j := range valSlice {
				if j > 0 {
					sql += ","
				}
				sql += "?"
			}
			sql += ")"
			args = append(args, valSlice...)
		case OpNull, OpNotNull:
			sql += c.Field + " " + c.Op
		default:
			// 默认按等号处理
			sql += c.Field + " = ?"
			args = append(args, c.Value)
		}
	}

	return sql, args
}

// OrderSql 构建 ORDER BY 子句，返回类似 "ORDER BY field1 ASC, field2 DESC"
func (m *Match) OrderSql() string {
	if len(m.Orders) == 0 {
		return ""
	}
	sql := "ORDER BY "
	for i, c := range m.Orders {
		if i > 0 {
			sql += ", "
		}
		sql += c.Field + " " + c.Op
	}
	return sql
}

// SetSql 构建 SET 子句和对应参数，通常用于 UPDATE 语句
func (m *Match) SetSql() (string, []any) {
	if len(m.Sets) == 0 {
		return "", nil
	}
	sql := "SET "
	args := make([]any, 0, len(m.Sets))
	for i, c := range m.Sets {
		if i > 0 {
			sql += ", "
		}
		sql += c.Field + " = ?"
		args = append(args, c.Value)
	}
	return sql, args
}

// SetMap 返回一个字段到值的 map，方便直接传给 ORM 的 Updates(map[string]interface{}) 等方法
func (m *Match) SetMap() map[string]any {
	res := make(map[string]any, len(m.Sets))
	for _, c := range m.Sets {
		res[c.Field] = c.Value
	}
	return res
}

// 辅助函数：把any转成[]any切片
func toSlice(value any) ([]any, bool) {
	switch v := value.(type) {
	case []any:
		return v, true
	case []string:
		res := make([]any, len(v))
		for i, s := range v {
			res[i] = s
		}
		return res, true
	case []int:
		res := make([]any, len(v))
		for i, n := range v {
			res[i] = n
		}
		return res, true
	case []int64:
		res := make([]any, len(v))
		for i, n := range v {
			res[i] = n
		}
		return res, true
	// 可以扩展其他切片类型...
	default:
		return nil, false
	}
}
