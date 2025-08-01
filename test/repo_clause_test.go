package test

import "testing"

// 2. 查询条件构造器
func TestRepo_Eq(t *testing.T)      {}
func TestRepo_NEq(t *testing.T)     {}
func TestRepo_In(t *testing.T)      {}
func TestRepo_Gt(t *testing.T)      {}
func TestRepo_Gte(t *testing.T)     {}
func TestRepo_Lt(t *testing.T)      {}
func TestRepo_Lte(t *testing.T)     {}
func TestRepo_Like(t *testing.T)    {}
func TestRepo_Null(t *testing.T)    {}
func TestRepo_NotNull(t *testing.T) {}
func TestRepo_Or(t *testing.T)      {}
func TestRepo_Where(t *testing.T)   {}

// 3. 字段选择、排序、分页设置
func TestRepo_Select(t *testing.T)   {}
func TestRepo_Omit(t *testing.T)     {}
func TestRepo_Asc(t *testing.T)      {}
func TestRepo_Desc(t *testing.T)     {}
func TestRepo_Limit(t *testing.T)    {}
func TestRepo_WithPage(t *testing.T) {}
