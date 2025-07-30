package test

import (
	"dubhe/db"
	"dubhe/test/model"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

// 1. 初始化与设置环境
func TestRepo_WithDB(t *testing.T) {
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(gdb)
	gdb.AutoMigrate(&model.Demo{})
	repo := model.NewDemoRepo()
	for range 1000000 {
		repo.WithDB(gdb).Save(&model.Demo{
			BaseModel: db.BaseModel{},
			Name:      "test",
		})
	}
	demo := len(repo.WithDB(gdb).List())
	fmt.Println(demo)
}

func TestRepo_WithCtx(t *testing.T) {

}
func TestRepo_OnErr(t *testing.T) {

}
func TestRepo_Clone(t *testing.T) {
	repo := model.NewDemoRepo()
	list := repo.List()
	fmt.Println(list)
}

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

// 4. 执行查询
func TestRepo_Get(t *testing.T)       {}
func TestRepo_GetOrInit(t *testing.T) {}
func TestRepo_List(t *testing.T)      {}
func TestRepo_Page(t *testing.T)      {}
func TestRepo_PageT(t *testing.T)     {}
func TestRepo_Count(t *testing.T)     {}
func TestRepo_Scan(t *testing.T)      {}
func TestRepo_Raw(t *testing.T)       {}

// 5. 数据写入操作
func TestRepo_Create(t *testing.T)      {}
func TestRepo_CreateBatch(t *testing.T) {}
func TestRepo_Save(t *testing.T)        {}
func TestRepo_Update(t *testing.T)      {}
func TestRepo_UpdateFull(t *testing.T)  {}
func TestRepo_Set(t *testing.T)         {}
func TestRepo_SetMap(t *testing.T)      {}
func TestRepo_Del(t *testing.T)         {}

// 6. 执行原始 SQL 写操作
func TestRepo_Exec(t *testing.T) {}

// 7. 错误获取
func TestRepo_Err(t *testing.T) {}

// 8. 事务控制
func TestRepo_Begin(t *testing.T)    {}
func TestRepo_Commit(t *testing.T)   {}
func TestRepo_Rollback(t *testing.T) {}
func TestRepo_Tx(t *testing.T)       {}

// 9. 基础访问
func TestRepo_DB(t *testing.T) {}
