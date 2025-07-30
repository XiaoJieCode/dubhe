package model

import (
	"dubhe/db"
)

type Demo struct {
	db.BaseModel
	Name string `gorm:"column:name;type:varchar(255);not null"`
}

func (d Demo) RepoDefine() db.RepoCfg {
	return db.RepoCfg{
		DataSource: "Main",
	}
}
func (d Demo) TableName() string {
	return "demo_model"
}
func NewDemoRepo() db.IRepo[Demo] {
	return db.NewRepo[Demo]()
}
