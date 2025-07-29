package db

import (
	"dubhe/db/ds"
	"dubhe/db/util/log"
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
			db:           db,
			RepoTemplate: _t.(*RepoTemplate[T]),
		}
	}

	e := db.AutoMigrate(model)
	if e != nil {
		panic(e)
	}
	logger := log.LogFactory.Get(key)
	temp := &RepoTemplate[T]{
		ctx:   nil,
		table: tableName,
		model: &model,
		key:   key,
		cfg:   &cfg,
		log:   logger,
	}
	templates.Store(key, temp)
	return &Repo[T]{
		db:           db,
		RepoTemplate: temp,
	}

}
