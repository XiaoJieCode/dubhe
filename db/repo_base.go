package db

import (
	"context"
	"dubhe/db/err"
	clause "dubhe/db/util/clause"
	"dubhe/db/util/log"
	"gorm.io/gorm"
	"slices"
)

type RepoTemplate[T IModel[K], K ID] struct {
	ctx   *context.Context
	table string
	model *T
	key   string
	cfg   *RepoCfg
	log   log.Log
}

type Repo[T IModel[K], K ID] struct {
	*RepoTemplate[T, K]
	db        *gorm.DB
	selects   []string
	match     clause.Match
	page      *Page
	onErrFunc []func(handle IRepoErrHandle[T, K])
	limit     int64
	err       error
	omits     []string
	isRaw     bool
}

type ErrHandler[T IModel[K], K ID] struct {
	repo       IRepo[T, K]
	err        error
	isContinue bool
	isCancel   bool
	isPanic    bool
}

func newErrHandler[T IModel[K], K ID](repo IRepo[T, K], err error) *ErrHandler[T, K] {
	return &ErrHandler[T, K]{repo: repo, err: err}
}

func (e *ErrHandler[T, K]) Repo() IRepo[T, K] {
	return e.repo
}

func (e *ErrHandler[T, K]) Error() error {
	return e.err
}

func (e *ErrHandler[T, K]) Continue() {
	e.isContinue = true
}

func (e *ErrHandler[T, K]) Cancel() {
	e.isCancel = true
}

func (e *ErrHandler[T, K]) Panic() {
	e.isPanic = true
}

func (r *Repo[T, K]) DB() *gorm.DB {
	return r.db
}

// Tx 返回新的事务 Repo 实例，独立于当前 Repo
func (r *Repo[T, K]) Tx() IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.db = newRepo.db.Begin().Session(&gorm.Session{NewDB: true})
	return newRepo
}

func (r *Repo[T, K]) Begin() IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.db = newRepo.db.Begin()
	return newRepo
}

func (r *Repo[T, K]) Commit() IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.db.Commit()
	return newRepo
}

func (r *Repo[T, K]) Rollback() IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.db.Rollback()
	return newRepo
}

func (r *Repo[T, K]) OnErr(f func(IRepoErrHandle[T, K])) IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.onErrFunc = append(newRepo.onErrFunc, f)
	return newRepo
}

func (r *Repo[T, K]) Err() *err.Error {
	return err.NewError(r.err)
}

func (r *Repo[T, K]) cloneInternal() *Repo[T, K] {

	var newPage *Page
	if r.page != nil {
		p := *r.page
		newPage = &p
	}

	return &Repo[T, K]{
		RepoTemplate: r.RepoTemplate,
		db:           r.db,
		selects:      slices.Clone(r.selects),
		match:        *r.match.Clone(),
		page:         newPage,
		err:          r.err,
		onErrFunc:    r.onErrFunc,
		limit:        r.limit,
		omits:        slices.Clone(r.omits),
		isRaw:        r.isRaw,
	}
}

// Clone 公开的克隆方法
func (r *Repo[T, K]) Clone() IRepo[T, K] {
	return r.cloneInternal()
}

// WithCtx 设置上下文，返回新的 Repo 实例
func (r *Repo[T, K]) WithCtx(ctx *context.Context) IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.ctx = ctx
	return newRepo
}

// WithDB 设置新的 *gorm.DB，返回当前实例（一般无需克隆）
func (r *Repo[T, K]) WithDB(db *gorm.DB) IRepo[T, K] {
	if db == nil {
		panic("db can not be nil")
	}
	r.db = db
	return r
}
