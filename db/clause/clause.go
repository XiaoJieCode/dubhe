package clause

import (
	"reflect"
	"slices"
	"strings"
)

// 定义一组 SQL 操作符常量，用于构建不同的条件/排序/更新表达式
const (
	OpEq      = "="           // 等于
	OpNEq     = "<>"          // 不等于
	OpGt      = ">"           // 大于
	OpGte     = ">="          // 大于等于
	OpLt      = "<"           // 小于
	OpLte     = "<="          // 小于等于
	OpIn      = "IN"          // 包含
	OpLike    = "LIKE"        // 模糊匹配
	OpNull    = "IS NULL"     // 为 NULL
	OpNotNull = "IS NOT NULL" // 不为 NULL
	OpOr      = "OR"          // OR 逻辑连接符
	OpAnd     = "AND"         // AND 逻辑连接符
	OpAsc     = "ASC"         // 升序排序
	OpDesc    = "DESC"        // 降序排序
	OpSet     = "="           // 用于 UPDATE SET 的赋值
)

// Clause 表示一个 SQL 子句（条件、排序、更新字段等）
// - Field: 字段名，例如 "name"、"id"
// - Value: 对应的值，例如 "Tom"、123（部分操作如 NULL/NOT NULL 不需要值）
// - Op: 操作符，例如 "="、">"、"<"、"LIKE"、"IN"、"DESC"（排序时用）
type Clause struct {
	Field string // 字段名
	Value any    // 字段值
	Op    string // 操作符
}

// Match 用来描述一组 SQL 子句集合，包括查询条件、排序条件、更新字段
// - Clauses: 存放 WHERE 条件，例如 age > 18、status = 'active'
// - Orders: 存放 ORDER BY 排序条件，例如 created_at DESC、id ASC
// - Sets:   存放 UPDATE SET 语句的赋值，例如 name = 'Tom'、count = 100
type Match struct {
	Clauses []Clause // WHERE 条件子句集合
	Orders  []Clause // ORDER BY 子句集合
	Sets    []Clause // UPDATE SET 子句集合
}

// Clone 深拷贝当前 Match，避免引用同一底层 slice 导致的副作用
func (m *Match) Clone() *Match {
	if m == nil {
		return NewMatch()
	}
	return &Match{
		Clauses: slices.Clone(m.Clauses),
		Orders:  slices.Clone(m.Orders),
		Sets:    slices.Clone(m.Sets),
	}
}

// NewMatch 创建一个新的空 Match
func NewMatch() *Match {
	return &Match{}
}

// add 内部方法，向 Clauses 添加一个条件子句
func (m *Match) add(field, op string, value any) *Match {
	m.Clauses = append(m.Clauses, Clause{Field: field, Op: op, Value: value})
	return m
}

// order 内部方法，向 Orders 添加一个排序子句
func (m *Match) order(field, op string) *Match {
	m.Orders = append(m.Orders, Clause{Field: field, Op: op})
	return m
}

// Set 添加一个 SET 子句（UPDATE 用）
func (m *Match) Set(field string, value any) *Match {
	m.Sets = append(m.Sets, Clause{Field: field, Value: value, Op: OpSet})
	return m
}

// ====== 以下为条件构造器 (WHERE 子句) ======

// Eq 相等条件，例如 age = 18
func (m *Match) Eq(field string, value any) *Match {
	return m.add(field, OpEq, value)
}

// NEq 不等条件，例如 age <> 18
func (m *Match) NEq(field string, value any) *Match {
	return m.add(field, OpNEq, value)
}

// Gt 大于条件，例如 age > 18
func (m *Match) Gt(field string, value any) *Match {
	return m.add(field, OpGt, value)
}

// Gte 大于等于条件，例如 age >= 18
func (m *Match) Gte(field string, value any) *Match {
	return m.add(field, OpGte, value)
}

// Lt 小于条件，例如 age < 18
func (m *Match) Lt(field string, value any) *Match {
	return m.add(field, OpLt, value)
}

// Lte 小于等于条件，例如 age <= 18
func (m *Match) Lte(field string, value any) *Match {
	return m.add(field, OpLte, value)
}

// In 包含条件，例如 id IN (1,2,3)
func (m *Match) In(field string, value any) *Match {
	return m.add(field, OpIn, value)
}

// Like 模糊匹配，例如 name LIKE '%Tom%'
func (m *Match) Like(field string, value any) *Match {
	return m.add(field, OpLike, value)
}

// Null 字段为 NULL，例如 deleted_at IS NULL
func (m *Match) Null(field string) *Match {
	return m.add(field, OpNull, nil)
}

// NotNull 字段不为 NULL，例如 updated_at IS NOT NULL
func (m *Match) NotNull(field string) *Match {
	return m.add(field, OpNotNull, nil)
}

// ====== 排序构造器 (ORDER BY 子句) ======

// Asc 升序排序，例如 id ASC
func (m *Match) Asc(field string) *Match {
	m.order(field, OpAsc)
	return m
}

// Desc 降序排序，例如 created_at DESC
func (m *Match) Desc(field string) *Match {
	m.order(field, OpDesc)
	return m
}

// WhereSql 生成 WHERE 子句及其参数，例如 "age > ? AND status = ?"  [18, "active"]
func (m *Match) WhereSql() (string, []any) {
	if len(m.Clauses) == 0 {
		return "", nil
	}

	sql := ""
	var args []any
	for i, c := range m.Clauses {
		sqlPart, arg := c.ToSqlStr()
		if sqlPart == "" {
			continue
		}
		if i > 0 {
			if c.Op != OpOr {
				sql += " " + OpAnd
			}
		}
		sql += sqlPart
		if arg != nil {
			args = append(args, arg...)
		}
	}

	return sql, args
}

// OrderSql 生成 ORDER BY 子句，例如 "ORDER BY created_at DESC, id ASC"
func (m *Match) OrderSql() string {
	if len(m.Orders) == 0 {
		return ""
	}
	sql := ""
	for i, c := range m.Orders {
		if i > 0 {
			sql += ", "
		}
		sql += c.Field + " " + c.Op
	}
	return sql
}

// SetSql 生成 SET 子句和参数（用于 UPDATE）
// 例如： "SET name = ?, age = ?"  [ "Tom", 20 ]
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

// SetMap 返回一个字段到值的映射，方便 ORM 执行 Updates(map[string]interface{})
func (m *Match) SetMap() map[string]any {
	res := make(map[string]any, len(m.Sets))
	for _, c := range m.Sets {
		res[c.Field] = c.Value
	}
	return res
}

// toSlice 辅助函数：尝试将任意类型转换为 []any

func toSlice(input any) ([]any, bool) {
	if input == nil {
		return nil, false
	}

	val := reflect.ValueOf(input)
	typ := val.Type()

	// 必须是切片
	if typ.Kind() != reflect.Slice {
		return nil, false
	}

	elemType := typ.Elem()
	// 如果是 struct，不支持
	if elemType.Kind() == reflect.Struct {
		return nil, false
	}

	length := val.Len()
	res := make([]any, length)
	for i := 0; i < length; i++ {
		res[i] = val.Index(i).Interface()
	}

	return res, true
}
func (c Clause) ToSqlStr() (string, []any) {
	switch c.Op {
	case OpEq, OpNEq, OpGt, OpGte, OpLt, OpLte, OpLike:
		return " " + c.Field + " " + c.Op + " ?", []any{c.Value}
	case OpNull, OpNotNull:
		return " " + c.Field + " " + c.Op, nil
	case OpIn:
		valSlice, ok := toSlice(c.Value)
		if !ok || len(valSlice) == 0 {
			return "", nil
		}
		return " " + c.Field + " IN (" + strings.Repeat("?,", len(valSlice)-1) + "?)", valSlice
	case OpAsc, OpDesc:
		return " " + c.Field + " " + c.Op, nil

	default:
		return " " + c.Field + " = ?", []any{c.Value}
	}
}
