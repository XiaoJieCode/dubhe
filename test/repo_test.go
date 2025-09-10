package test

import (
	"dubhe/test/model"
	"fmt"
	"log"

	"gorm.io/gorm"
)

// InitSingleData 插入一条测试数据
func InitSingleData(db *gorm.DB) {
	m := &model.Demo{}
	if err := db.AutoMigrate(m); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	db.Create(&model.Demo{
		Name: "demo",
		Age:  0,
	})
}

// Init100Data 插入100条测试数据
func Init100Data(db *gorm.DB) {
	m := &model.Demo{}
	if err := db.AutoMigrate(m); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}
	var demos []model.Demo
	for i := 1; i <= 100; i++ {
		demos = append(demos, model.Demo{
			Name: fmt.Sprintf("demo-%d", i),
			Age:  i % 30,
		})
	}
	if err := db.Create(&demos).Error; err != nil {
		log.Fatalf("Insert 100 demo data failed: %v", err)
	}
}

// ClearAll 清空所有测试数据
func ClearAll(db *gorm.DB) {
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Demo{}).Error; err != nil {
		log.Fatalf("ClearAll failed: %v", err)
	}
}
