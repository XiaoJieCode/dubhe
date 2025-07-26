package test

import (
	"dubhe/db"
	"dubhe/db/err"
	"dubhe/test/testdata"
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	UserRepo := db.NewRepo[testdata.User]()

	user := UserRepo.Eq("id", 1).OnErr(func(i db.IRepo[testdata.User], e *err.Error) bool {
		return true
	}).Get()
	fmt.Println(user)

}
