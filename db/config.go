package db

import (
	"dubhe/db/ds"
	"gorm.io/gorm"
)

type RepoCfg struct {
	UseDataSource string
}

type RepoDefine interface {
	RepoDefine() RepoCfg
	TableName() string
}

// NewRepo Build Repo
func NewRepo[T RepoDefine]() IRepo[T] {
	model := *new(T)
	tableName := model.TableName()
	cfg := model.RepoDefine()
	var db *gorm.DB
	if tableName == "" {
		panic("table name is empty")
	}
	if cfg.UseDataSource == "" {
		db = ds.MustGetDB()
	} else {
		db = ds.MustGetDB(cfg.UseDataSource)
	}
	e := db.AutoMigrate(model) // 自动创建或更新 user 表
	if e != nil {
		panic(e)
	}
	return &Repo[T]{
		db:    db,
		ctx:   nil,
		table: tableName,
		model: &model,
		cfg:   &cfg,
	}

}
