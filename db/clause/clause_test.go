package clause

import (
	"reflect"
	"testing"
)

// ---------- 基础构造与链式调用 ----------

func TestNewMatch_Empty(t *testing.T) {
	m := NewMatch()
	if m == nil {
		t.Fatalf("NewMatch() returned nil")
	}
	if len(m.Clauses) != 0 || len(m.Orders) != 0 || len(m.Sets) != 0 {
		t.Fatalf("new match should be empty")
	}
}

func TestClone_DeepCopy(t *testing.T) {
	m := NewMatch().
		Eq("a", 1).
		Gt("b", 2).
		Asc("c").
		Set("x", 100)

	c := m.Clone()

	// 改原对象，不影响克隆体
	m.NEq("d", 3).Desc("e").Set("y", 200)
	if len(c.Clauses) != 2 || len(c.Orders) != 1 || len(c.Sets) != 1 {
		t.Fatalf("clone should be deep-copied and unaffected by original modifications")
	}

	// 改克隆体，不影响原对象
	c.Like("f", "%k%").Asc("g").Set("z", 300)
	if len(m.Clauses) != 3 || len(m.Orders) != 2 || len(m.Sets) != 2 {
		t.Fatalf("original should be unaffected by clone modifications")
	}
}

// ---------- WhereSql 基本操作 ----------

func TestWhereSql_BasicOps(t *testing.T) {
	m := NewMatch().
		Eq("a", 1).
		NEq("b", "x").
		Gt("c", 2).
		Gte("d", 3).
		Lt("e", 4).
		Lte("f", 5).
		Like("g", "%hi%")

	sql, args := m.WhereSql()
	wantSQL := " a = ? AND b <> ? AND c > ? AND d >= ? AND e < ? AND f <= ? AND g LIKE ?"
	wantArgs := []any{1, "x", 2, 3, 4, 5, "%hi%"}
	if sql != wantSQL {
		t.Fatalf("WhereSql mismatch.\n got: %q\nwant: %q", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("WhereSql args mismatch.\n got: %#v\nwant: %#v", args, wantArgs)
	}
}

func TestWhereSql_Null_NotNull(t *testing.T) {
	m := NewMatch().Null("deleted_at").NotNull("created_at")
	sql, args := m.WhereSql()
	wantSQL := " deleted_at IS NULL AND created_at IS NOT NULL"
	if sql != wantSQL {
		t.Fatalf("WhereSql NULL/NOT NULL mismatch.\n got: %q\nwant: %q", sql, wantSQL)
	}
	if len(args) != 0 {
		t.Fatalf("NULL/NOT NULL should not produce args, got: %#v", args)
	}
}

func TestWhereSql_In_Ints(t *testing.T) {
	m := NewMatch().In("id", []int{1, 2, 3})
	sql, args := m.WhereSql()
	wantSQL := " id IN (?,?,?)"
	wantArgs := []any{1, 2, 3}
	if sql != wantSQL {
		t.Fatalf("IN ints SQL mismatch.\n got: %q\nwant: %q", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("IN ints args mismatch.\n got: %#v\nwant: %#v", args, wantArgs)
	}
}

func TestWhereSql_In_Strings(t *testing.T) {
	m := NewMatch().In("code", []string{"A", "B"})
	sql, args := m.WhereSql()
	wantSQL := " code IN (?,?)"
	wantArgs := []any{"A", "B"}
	if sql != wantSQL {
		t.Fatalf("IN strings SQL mismatch.\n got: %q\nwant: %q", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("IN strings args mismatch.\n got: %#v\nwant: %#v", args, wantArgs)
	}
}

func TestWhereSql_In_AnyMixed(t *testing.T) {
	m := NewMatch().In("v", []any{1, "x"})
	sql, args := m.WhereSql()
	wantSQL := " v IN (?,?)"
	wantArgs := []any{1, "x"}
	if sql != wantSQL {
		t.Fatalf("IN []any SQL mismatch.\n got: %q\nwant: %q", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("IN []any args mismatch.\n got: %#v\nwant: %#v", args, wantArgs)
	}
}

// ---------- WhereSql 边界/异常路径 ----------

func TestWhereSql_NoClauses(t *testing.T) {
	m := NewMatch()
	sql, args := m.WhereSql()
	if sql != "" || args != nil {
		t.Fatalf("empty Match should return empty SQL and nil args, got: %q %#v", sql, args)
	}
}

func TestWhereSql_In_EmptyOnly(t *testing.T) {
	// 单独一个空 IN：当前实现会 continue，整体返回空 SQL
	m := NewMatch().In("id", []int{})
	sql, args := m.WhereSql()
	if sql != "" || args != nil {
		t.Fatalf("empty IN only should return empty SQL and nil args, got: %q %#v", sql, args)
	}
}

func TestWhereSql_In_Empty_Middle_BUG(t *testing.T) {
	// 期望行为：跳过空 IN，连接符不应损坏 -> "x = ? AND y = ?"
	m := NewMatch().Eq("x", 1).In("id", []int{}).Eq("y", 2)
	sql, args := m.WhereSql()

	wantSQL := " x = ? AND y = ?"
	wantArgs := []any{1, 2}

	if sql != wantSQL || !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("BUG: empty IN in the middle corrupts connectors.\n got SQL: %q\nwant SQL: %q\n got args: %#v\nwant args: %#v", sql, wantSQL, args, wantArgs)
	}
}

// ---------- OrderSql ----------

func TestOrderSql_Empty(t *testing.T) {
	m := NewMatch()
	if got := m.OrderSql(); got != "" {
		t.Fatalf("empty orders should produce empty string, got: %q", got)
	}
}

func TestOrderSql_Multiple(t *testing.T) {
	m := NewMatch().Asc("name").Desc("age")
	want := "name ASC, age DESC"
	if got := m.OrderSql(); got != want {
		t.Fatalf("OrderSql mismatch.\n got: %q\nwant: %q", got, want)
	}
}

// ---------- SetSql / SetMap ----------

func TestSetSql_And_SetMap(t *testing.T) {
	m := NewMatch().Set("x", 1).Set("y", "abc")
	sql, args := m.SetSql()
	wantSQL := "SET x = ?, y = ?"
	wantArgs := []any{1, "abc"}

	if sql != wantSQL {
		t.Fatalf("SetSql mismatch.\n got: %q\nwant: %q", sql, wantSQL)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("SetSql args mismatch.\n got: %#v\nwant: %#v", args, wantArgs)
	}

	mp := m.SetMap()
	if len(mp) != 2 || mp["x"] != 1 || mp["y"] != "abc" {
		t.Fatalf("SetMap mismatch, got: %#v", mp)
	}
}

func TestSetMap_DuplicateField_LastWins(t *testing.T) {
	m := NewMatch().Set("x", 1).Set("x", 2)
	mp := m.SetMap()
	if len(mp) != 1 || mp["x"] != 2 {
		t.Fatalf("SetMap should take the last value for duplicate fields, got: %#v", mp)
	}

	sql, args := m.SetSql()
	// SetSql 会如实输出两次赋值，这是允许的（由调用方决定是否允许重复字段）
	wantSQL := "SET x = ?, x = ?"
	wantArgs := []any{1, 2}
	if sql != wantSQL || !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("SetSql with duplicate fields mismatch.\n got SQL: %q\nwant SQL: %q\n got args: %#v\nwant args: %#v", sql, wantSQL, args, wantArgs)
	}
}

// ---------- toSlice ----------

func TestToSlice_SupportedTypes(t *testing.T) {
	cases := []struct {
		in   any
		want []any
		ok   bool
	}{
		{[]any{1, "x"}, []any{1, "x"}, true},
		{[]string{"a", "b"}, []any{"a", "b"}, true},
		{[]int{1, 2, 3}, []any{1, 2, 3}, true},
		{[]int64{7, 8}, []any{int64(7), int64(8)}, true},
	}

	for i, c := range cases {
		got, ok := toSlice(c.in)
		if ok != c.ok {
			t.Fatalf("case %d ok mismatch: got %v want %v", i, ok, c.ok)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Fatalf("case %d slice mismatch: got %#v want %#v", i, got, c.want)
		}
	}
}

func TestToSlice_Unsupported(t *testing.T) {
	for i, in := range []any{123, "str", nil, struct{}{}} {
		got, ok := toSlice(in)
		if ok || got != nil {
			t.Fatalf("case %d should be unsupported, got: %#v, ok=%v", i, got, ok)
		}
	}
}
func TestClause_ToSqlStr_Default(t *testing.T) {
	c := Clause{Field: "foo", Value: 42, Op: "UNKNOWN"}
	sql, args := c.ToSqlStr()
	wantSQL := " foo = ?"
	wantArgs := []any{42}
	if sql != wantSQL || !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("default branch mismatch.\n got SQL: %q args=%#v\nwant SQL: %q args=%#v", sql, args, wantSQL, wantArgs)
	}
}
func TestClause_ToSqlStr_OrderOps(t *testing.T) {
	c1 := Clause{Field: "name", Op: OpAsc}
	sql1, args1 := c1.ToSqlStr()
	if sql1 != " name ASC" || args1 != nil {
		t.Fatalf("ASC mismatch, got: %q %#v", sql1, args1)
	}

	c2 := Clause{Field: "age", Op: OpDesc}
	sql2, args2 := c2.ToSqlStr()
	if sql2 != " age DESC" || args2 != nil {
		t.Fatalf("DESC mismatch, got: %q %#v", sql2, args2)
	}
}
func TestToSlice_StructSliceRejected(t *testing.T) {
	type S struct{ X int }
	got, ok := toSlice([]S{{1}, {2}})
	if ok || got != nil {
		t.Fatalf("struct slice should be rejected, got: %#v ok=%v", got, ok)
	}
}
func TestToSlice_AliasTypeSlice(t *testing.T) {
	type MyInt int
	in := []MyInt{10, 20}
	got, ok := toSlice(in)
	want := []any{MyInt(10), MyInt(20)}
	if !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("alias slice mismatch, got: %#v ok=%v", got, ok)
	}
}
