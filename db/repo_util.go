package db

import "dubhe/db/err"

func (r *Repo[T]) checkErr(e error) {
	if r.onErrFunc != nil {
		needPanic := r.onErrFunc(r, err.NewError(e))
		if needPanic {
			panic(r.err)
		}
	}
}
