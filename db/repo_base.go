package db

import (
	"context"
	"dubhe/db/err"
	clause "dubhe/db/util/condition"
	"gorm.io/gorm"
)

// Repo It provides some methods to do orm operation
type Repo[T any] struct {
	db    *gorm.DB
	ctx   *context.Context
	table string
	model *T
	cfg   *RepoCfg

	match     clause.Match
	page      *Page
	err       error
	onErrFunc func(IRepo[T], *err.Error) bool
}

func (r *Repo[T]) Tx() IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Begin() IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Commit() IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Rollback() IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) OnErr(f func(IRepo[T], *err.Error) bool) IRepo[T] {
	r.onErrFunc = f
	return r
}

func (r *Repo[T]) Clone() IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) WithCtx(ctx *context.Context) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) WithDB(db *gorm.DB) IRepo[T] {
	//TODO implement me
	panic("implement me")
}
