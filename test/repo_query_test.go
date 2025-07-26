package test

import (
	"dubhe/db"
	"dubhe/db/err"
	"dubhe/test/model"
	"github.com/spf13/cast"
	"testing"
)

func Panic(db.IRepo[model.User], *err.Error) bool {
	return true
}
func TestLimit(t *testing.T) {
	list := model.UserRepo().Limit(1).List()
	println(len(list))
}
func TestGet(t *testing.T) {
	model.UserRepo().OnErr(func(i db.IRepo[model.User], e *err.Error) bool {

		return true
	}).Eq("id", -1).Get()

	u1 := model.UserRepo().Eq("id", 2).
		OnErr(Panic).
		Get()
	println(u1)

	u2 := model.UserRepo().Eq("id", -1).GetOrInit()
	println(u2)

	users := model.UserRepo().List()
	println(users)

	pageUsers := model.UserRepo().WithPage(&db.Page{
		Page: 1,
		Size: 10,
	}).Page()
	println(pageUsers)
	pageTUsers := model.UserRepo().WithPage(&db.Page{
		Page: 1,
		Size: 10,
	}).PageT()
	println(pageTUsers)
	counts := model.UserRepo().Gt("id", 2).Count()
	println(counts)
	type User struct {
		ID int64 `json:"id"`
	}

	dest := new(User)
	model.UserRepo().Eq("id", 2).Scan(dest)
	println(dest)

}
func TestTx(t *testing.T) {
	repo := model.UserRepo().Tx()
	affect := int64(0)
	before := len(repo.List())
	defer func() {
		println("row affect: " + cast.ToString(affect))
		p := model.UserRepo()
		after := p.List()
		println("before: " + cast.ToString(before))
		println("after: " + cast.ToString(len(after)))
	}()
	defer repo.Rollback()
	affect += repo.Create(&model.User{
		BaseModel: db.BaseModel{},
		Name:      "",
		Age:       0,
	})
	affect += repo.Create(&model.User{
		BaseModel: db.BaseModel{},
		Name:      "",
		Age:       0,
	})
	affect += repo.Create(&model.User{
		BaseModel: db.BaseModel{},
		Name:      "",
		Age:       0,
	})
}
