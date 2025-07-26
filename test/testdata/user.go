package testdata

import (
	"dubhe/db"
)

type User struct {
	db.BaseModel
	Name string `gorm:"column:name"`
	Age  int    `gorm:"column:age"`
}

func (u User) RepoDefine() db.RepoCfg {
	return db.RepoCfg{}
}

func (u User) TableName() string {
	return "user"
}
