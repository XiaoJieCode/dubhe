package db

import (
	"dubhe/db/err"
	"fmt"
)

func (r *Repo[T]) checkErr(e error) {
	if r.onErrFunc != nil {
		needPanic := r.onErrFunc(r, err.NewError(e))
		if needPanic {
			panic(r.err)
		}
	} else {
		fmt.Println("error not handle: ", e)
	}
}
