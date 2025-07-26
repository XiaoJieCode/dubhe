package db

import (
	"dubhe/db/ds"
	"gorm.io/gorm"
	"sync"
)

var (
	templates = &sync.Map{}
)

type RepoCfg struct {
	DataSource string
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
	if cfg.DataSource == "" {
		db = ds.MustGetDB()
	} else {
		db = ds.MustGetDB(cfg.DataSource)
	}
	key := cfg.DataSource + "_" + tableName
	_t, ok := templates.Load(key)
	if ok {
		return &Repo[T]{
			DB:           db,
			RepoTemplate: _t.(*RepoTemplate[T]),
		}
	}

	e := db.AutoMigrate(model)
	if e != nil {
		panic(e)
	}
	temp := &RepoTemplate[T]{
		ctx:   nil,
		table: tableName,
		model: &model,
		cfg:   &cfg,
	}
	templates.Store(key, temp)
	return &Repo[T]{
		DB:           db,
		RepoTemplate: temp,
	}

}
