package jn

import "fmt"

func JoinByKey[A any, B any, C any, K comparable](
	as []A,
	bs []B,
	keyFnA func(a A) K,
	keyFnB func(b B) K,
	mergeFn func(a A, b B) C,
) []C {
	bMap := make(map[K]B, len(bs))
	for _, b := range bs {
		bMap[keyFnB(b)] = b
	}

	var result []C
	for _, a := range as {
		if b, ok := bMap[keyFnA(a)]; ok {
			result = append(result, mergeFn(a, b))
		}
	}
	return result
}
func main() {
	type User struct {
		UserID int64
		Name   string
	}
	type UserExt struct {
		UserID int64
		Email  string
		Age    int
	}
	type UserResp struct {
		UserID int64
		Name   string
		Email  string
		Age    int
	}
	users := []User{
		{UserID: 1, Name: "a"},
		{UserID: 2, Name: "b"},
		{UserID: 3, Name: "c"},
	}
	usersExt := []UserExt{
		{UserID: 1, Email: "a@gmail.com", Age: 18},
		{UserID: 2, Email: "b@gmail.com", Age: 19},
		{UserID: 3, Email: "c@gmail.com", Age: 20},
	}
	joinByKey := JoinByKey(users, usersExt, func(a User) int64 { return a.UserID }, func(b UserExt) int64 { return b.UserID }, func(a User, b UserExt) UserResp {
		return UserResp{
			UserID: a.UserID,
			Name:   a.Name,
			Email:  b.Email,
			Age:    b.Age,
		}
	})
	fmt.Println(joinByKey)
}
