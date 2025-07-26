package testdata

import (
	"dubhe/db"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `gorm:"column:name"`
	Age  int    `gorm:"column:age"`
}

var UserRepo = db.NewRepo[User](db.RepoCfg{
	TableName: "user",
})
