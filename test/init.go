package test

import (
	"dubhe/db/ds"
)

func init() {
	err := ds.RegisterDataSource("Main", ds.DBConfig{
		DSN:    "root:123456@tcp(127.0.0.1:3306)/dubhe?charset=utf8mb4&parseTime=True&loc=Local",
		Driver: "mysql",
	})
	if err != nil {
		panic(err)
	}
}
