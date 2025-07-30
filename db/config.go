package db

import (
	"dubhe/db/ds"
	"dubhe/db/util/log"
	"gorm.io/gorm"
	"sync"
)

// 缓存模板，避免重复构建相同表结构的 RepoTemplate
var templates = &sync.Map{}

// RepoCfg 定义 Repo 的数据库配置
type RepoCfg struct {
	DataSource  string   // 指定数据源名
	DB          *gorm.DB // 指定 DB 实例（优先级高于 DataSource）
	AutoMigrate bool     // 是否自动迁移表结构（默认启用）
}

// RepoDefine 接口用于模型绑定 Repo 配置
type RepoDefine interface {
	RepoDefine() RepoCfg
	TableName() string
}

// NewRepo 创建并初始化一个通用 Repo
func NewRepo[T RepoDefine]() IRepo[T] {
	// 获取模型信息
	model := *new(T)
	tableName := model.TableName()
	cfg := model.RepoDefine()

	// 校验表名
	if tableName == "" {
		panic("table name is empty")
	}

	// 获取数据库实例
	var db *gorm.DB
	switch {
	case cfg.DB != nil:
		db = cfg.DB
	case cfg.DataSource != "":
		db = ds.MustGetDB(cfg.DataSource)
	default:
		db = ds.MustGetDB()
	}

	// 构造缓存键（数据源 + 表名）
	key := tableName
	if cfg.DataSource != "" {
		key = cfg.DataSource + "." + tableName
	} else {
		key = "default." + tableName
	}

	// 尝试从缓存中获取已有模板
	if cached, ok := templates.Load(key); ok {
		return &Repo[T]{
			db:           db,
			RepoTemplate: cached.(*RepoTemplate[T]),
		}
	}

	// 自动迁移表结构（默认开启）
	if cfg.AutoMigrate || cfg.DB == nil {
		if err := db.AutoMigrate(model); err != nil {
			panic(err)
		}
	}

	// 创建并缓存新的 RepoTemplate
	logger := log.LogFactory.Get(key)
	template := &RepoTemplate[T]{
		ctx:   nil,
		table: tableName,
		model: &model,
		key:   key,
		cfg:   &cfg,
		log:   logger,
	}
	templates.Store(key, template)

	return &Repo[T]{
		db:           db,
		RepoTemplate: template,
	}
}
