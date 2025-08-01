package test

import (
	"dubhe/db/ds"
)

func init() {
	err := ds.RegisterDataSource("sqlite", ds.DBConfig{
		DSN:    ":memory:",
		Driver: "sqlite",
	})
	if err != nil {
		panic(err)
	}
	err = ds.RegisterDataSource(ds.Default, ds.DBConfig{
		DSN:    "root:123456@tcp(127.0.0.1:3306)/dubhe?charset=utf8mb4&parseTime=True&loc=Local",
		Driver: "mysql",
	})
	if err != nil {
		panic(err)
	}
}
