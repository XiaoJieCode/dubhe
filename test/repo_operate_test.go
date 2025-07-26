package test

import (
	"dubhe/test/model"
	"testing"
)

func TestUpdateFull(t *testing.T) {
	before := model.UserRepo().Eq("id", 1).Get()
	affect := model.UserRepo().Eq("id", 1).
		UpdateFull(&model.User{
			Name: "test",
			Age:  1,
		})
	after := model.UserRepo().Eq("id", 1).Get()

	println("before", before)
	println("after", after)
	println("affect", affect)
}

func TestUpdate(t *testing.T) {
	before := model.UserRepo().Eq("id", 1).Get()
	affect := model.UserRepo().
		Eq("id", 1).
		Set("age", 100).
		Update()

	after := model.UserRepo().Eq("id", 1).Get()
	println("before", before)
	println("after", after)
	println("affect", affect)

}
