package ds

import (
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mu    sync.RWMutex
	dbMap = make(map[string]*gorm.DB) // key为数据源名称，value为对应gorm.DB连接实例
)

// DBConfig 是多数据源注册时传入的数据库配置结构体
type DBConfig struct {
	// 基础连接配置
	DSN    string // 数据源名称（MySQL的连接串）
	Driver string // 数据库驱动，如 "mysql", "postgres", "sqlite" 等

	// GORM日志配置
	LogLevel logger.LogLevel // gorm日志级别

	// 连接池配置
	MaxIdleConns    int           // 连接池最大空闲连接数
	MaxOpenConns    int           // 连接池最大打开连接数
	ConnMaxLifetime time.Duration // 连接最大生命周期
	ConnMaxIdleTime time.Duration // 连接最大空闲时间

	// 其他配置字段可以扩展
}

// RegisterDataSource 注册一个新的数据源连接
func RegisterDataSource(name string, cfg DBConfig) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := dbMap[name]; exists {
		return errors.New("data source already exists: " + name)
	}

	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.DSN)
	case "sqlite":
		dialector = sqlite.Open(cfg.DSN)
	// 其他数据库可扩展
	default:
		return errors.New("unsupported driver: " + cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接池参数
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.ConnMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	}

	dbMap[name] = db
	return nil
}

// GetDB 支持可选参数name：
// - 如果未传name且只有一个数据源，则返回该唯一数据源
// - 如果未传name且多于一个数据源，返回错误
// - 如果传了name，则返回对应数据源或错误
func GetDB(names ...string) (*gorm.DB, error) {
	mu.RLock()
	defer mu.RUnlock()

	if len(names) == 0 {
		if len(dbMap) == 1 {
			for _, db := range dbMap {
				return db, nil
			}
		}
		return nil, errors.New("multiple data sources exist, must specify a name")
	}

	name := names[0]
	db, exists := dbMap[name]
	if !exists {
		return nil, fmt.Errorf("data source not found: %s", name)
	}
	return db, nil
}

// MustGetDB 获取数据源连接，找不到时直接panic，适合初始化阶段使用
func MustGetDB(names ...string) *gorm.DB {
	db, err := GetDB(names...)
	if err != nil {
		panic(err)
	}
	return db
}

// CloseAll 关闭所有数据源连接
func CloseAll() error {
	mu.Lock()
	defer mu.Unlock()

	for name, db := range dbMap {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		err = sqlDB.Close()
		if err != nil {
			return err
		}
		delete(dbMap, name)
	}
	return nil
}
