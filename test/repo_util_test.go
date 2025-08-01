package test

import (
	"dubhe/db"
	"dubhe/test/model"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	_ = gdb.AutoMigrate(&model.Demo{})
	repo := model.NewDemoRepo()
	for range 1000000 {
		repo.WithDB(gdb).Save(&model.Demo{
			Name: "test",
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

// 7. 错误获取
func TestRepo_Err(t *testing.T) {
	affected := model.NewDemoRepo().OnErr(func(handle db.IRepoErrHandle[model.Demo, int64]) {
		handle.Panic()
	}).Create(&model.Demo{
		ID:   1,
		Name: "test",
		Age:  0,
	})
	fmt.Println(affected)
}

func TestRepo_OnSQLErr(t *testing.T) {
	list := model.NewDemoRepo().Where("a<d").List()
	assert.Equal(t, 0, len(list))
}

// 8. 事务控制
func TestRepo_Begin(t *testing.T)    {}
func TestRepo_Commit(t *testing.T)   {}
func TestRepo_Rollback(t *testing.T) {}
func TestRepo_Tx(t *testing.T)       {}

// 9. 基础访问
func TestRepo_DB(t *testing.T) {}
