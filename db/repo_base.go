package db

import (
	"context"
	"dubhe/db/err"
	clause "dubhe/db/util/condition"
	"gorm.io/gorm"
)

type RepoTemplate[T any] struct {
	ctx   *context.Context
	table string
	model *T
	cfg   *RepoCfg
}

// Repo It provides some methods to do orm operation
type Repo[T any] struct {
	*RepoTemplate[T]
	DB        *gorm.DB
	selects   []string
	match     clause.Match
	page      *Page
	err       error
	onErrFunc func(IRepo[T], *err.Error) bool
	limit     int64
	omits     []string
}

func (r *Repo[T]) Tx() IRepo[T] {
	r.DB = r.DB.Begin()
	return r
}

func (r *Repo[T]) Begin() IRepo[T] {
	r.DB = r.DB.Begin()
	return r
}

func (r *Repo[T]) Commit() IRepo[T] {
	r.DB.Commit()
	return r
}

func (r *Repo[T]) Rollback() IRepo[T] {
	r.DB.Rollback()
	return r
}

func (r *Repo[T]) OnErr(f func(IRepo[T], *err.Error) bool) IRepo[T] {
	r.onErrFunc = f
	return r
}
func (r *Repo[T]) Err() *err.Error {
	return err.NewError(r.err)
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
	if db == nil {
		panic("DB can not be nil")
	}
	r.DB = db
	return r
}
