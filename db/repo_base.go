package db

import (
	"context"
	"dubhe/db/err"
	clause "dubhe/db/util/clause"
	"dubhe/db/util/log"
	"gorm.io/gorm"
	"slices"
)

type RepoTemplate[T any] struct {
	ctx   *context.Context
	table string
	model *T
	key   string
	cfg   *RepoCfg
	log   log.Log
}

type Repo[T any] struct {
	*RepoTemplate[T]
	db        *gorm.DB
	selects   []string
	match     clause.Match
	page      *Page
	onErrFunc []func(handle IRepoErrHandle[T])
	limit     int64
	err       error
	omits     []string
	data      []*T
	isRaw     bool
}

type ErrHandler[T any] struct {
	repo       IRepo[T]
	err        error
	isContinue bool
	isCancel   bool
	isPanic    bool
}

func newErrHandler[T any](repo IRepo[T], err error) *ErrHandler[T] {
	return &ErrHandler[T]{repo: repo, err: err}
}

func (e ErrHandler[T]) Repo() IRepo[T] {
	return e.repo
}

func (e ErrHandler[T]) Error() error {
	return e.err
}

func (e ErrHandler[T]) Continue() {
	e.isContinue = true
}

func (e ErrHandler[T]) Cancel() {
	e.isCancel = true
}

func (e ErrHandler[T]) Panic() {
	e.isPanic = true
}

func (r *Repo[T]) DB() *gorm.DB {
	return r.db
}

// Tx 返回新的事务 Repo 实例，独立于当前 Repo
func (r *Repo[T]) Tx() IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.db = newRepo.db.Begin().Session(&gorm.Session{NewDB: true})
	return newRepo
}

func (r *Repo[T]) Begin() IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.db = newRepo.db.Begin()
	return newRepo
}

func (r *Repo[T]) Commit() IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.db.Commit()
	return newRepo
}

func (r *Repo[T]) Rollback() IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.db.Rollback()
	return newRepo
}

func (r *Repo[T]) OnErr(f func(IRepoErrHandle[T])) IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.onErrFunc = append(newRepo.onErrFunc, f)
	return newRepo
}

func (r *Repo[T]) Err() *err.Error {
	return err.NewError(r.err)
}

func (r *Repo[T]) cloneInternal() *Repo[T] {
	newTemplate := &RepoTemplate[T]{
		ctx:   r.ctx,
		table: r.table,
		model: r.model,
		key:   r.key,
		cfg:   r.cfg,
	}

	var newPage *Page
	if r.page != nil {
		p := *r.page
		newPage = &p
	}

	return &Repo[T]{
		RepoTemplate: newTemplate,
		db:           r.db,
		selects:      slices.Clone(r.selects),
		match:        *r.match.Clone(),
		page:         newPage,
		err:          r.err,
		onErrFunc:    r.onErrFunc,
		limit:        r.limit,
		omits:        slices.Clone(r.omits),
		data:         nil,
		isRaw:        r.isRaw,
	}
}

// Clone 公开的克隆方法
func (r *Repo[T]) Clone() IRepo[T] {
	return r.cloneInternal()
}

// WithCtx 设置上下文，返回新的 Repo 实例
func (r *Repo[T]) WithCtx(ctx *context.Context) IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.ctx = ctx
	return newRepo
}

// WithDB 设置新的 *gorm.DB，返回当前实例（一般无需克隆）
func (r *Repo[T]) WithDB(db *gorm.DB) IRepo[T] {
	if db == nil {
		panic("db can not be nil")
	}
	r.db = db
	return r
}
