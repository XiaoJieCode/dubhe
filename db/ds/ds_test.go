package ds

import (
	"errors"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open sqlite in-memory db: %v", err)
	}
	return db
}

func TestRegisterAndGetDB(t *testing.T) {
	defer func() { _ = CloseAll() }()

	// 注册数据源
	cfg := DBConfig{
		DSN:             ":memory:",
		Driver:          "sqlite",
		LogLevel:        logger.Silent,
		MaxIdleConns:    5,
		MaxOpenConns:    10,
		ConnMaxLifetime: time.Minute,
		ConnMaxIdleTime: time.Second * 30,
	}
	err := RegisterDataSource(Default, cfg)
	assert.NoError(t, err)

	// 获取数据源
	db, err := GetDB()
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// MustGetDB 不应 panic
	assert.NotPanics(t, func() {
		_ = MustGetDB()
	})
}

func TestRegisterGorm(t *testing.T) {
	defer func() { _ = CloseAll() }()

	db := setupTestDB(t)

	// 使用 RegisterGorm
	err := RegisterGorm("custom", db)
	assert.NoError(t, err)

	got, err := GetDB("custom")
	assert.NoError(t, err)
	assert.Equal(t, db, got)
}

func TestDuplicateRegister(t *testing.T) {
	defer func() { _ = CloseAll() }()

	db := setupTestDB(t)
	err := RegisterGorm("dup", db)
	assert.NoError(t, err)

	// 再次注册同名数据源应报错
	err = RegisterGorm("dup", db)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, err))
}

func TestGetDBErrors(t *testing.T) {
	defer func() { _ = CloseAll() }()

	// 没有数据源时
	_, err := GetDB()
	assert.Error(t, err)

	// 注册两个数据源
	_ = RegisterGorm("db1", setupTestDB(t))
	_ = RegisterGorm("db2", setupTestDB(t))

	// 不指定名字时应报错
	_, err = GetDB()
	assert.Error(t, err)

	// 找不到指定数据源
	_, err = GetDB("not-exist")
	assert.Error(t, err)
}

func TestCloseAll(t *testing.T) {
	_ = RegisterGorm("db1", setupTestDB(t))
	_ = RegisterGorm("db2", setupTestDB(t))

	err := CloseAll()
	assert.NoError(t, err)

	// 确保已经清空
	_, err = GetDB()
	assert.Error(t, err)
}
