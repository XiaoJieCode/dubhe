package db

import (
	"dubhe/db/ds"
	"gorm.io/gorm"
)

type RepoCfg struct {
	TableName  string
	Datasource string
}

type RepoDefine interface {
	NewRepo func(cfg RepoCfg) *Repo[]
}
// NewRepo Build Repo
func NewRepo[T any](cfg RepoCfg) *Repo[T] {
	var db *gorm.DB
	if cfg.TableName == "" {
		panic("table name is empty")
	}
	if cfg.Datasource == "" {
		db = ds.MustGetDB()
	} else {
		db = ds.MustGetDB(cfg.Datasource)
	}
	var model T

	return &Repo[T]{
		DB:    db,
		Ctx:   nil,
		Table: cfg.TableName,
		Model: &model,
		Cfg:   &cfg,
	}

}
