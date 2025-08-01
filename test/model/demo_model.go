package model

import (
	"dubhe/db"
	"dubhe/db/ds"
)

type Demo struct {
	ID   int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name;type:varchar(255);not null"`
	Age  int    `gorm:"column:age;type:int;not null"`
}

func (d Demo) RepoDefine() db.RepoCfg {
	return db.RepoCfg{
		DataSource: ds.Default,
	}
}

func (d Demo) TableName() string {
	return "demo"
}

func (d Demo) GetID() int64 {
	return d.ID
}

func (d Demo) IsNil() bool {
	return d.ID == 0
}

func NewDemoRepo() db.IRepo[Demo, int64] {
	return db.NewRepo[Demo, int64]()
}
