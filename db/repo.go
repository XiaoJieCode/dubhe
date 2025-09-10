package db

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/xiaojiecode/dubhe/db/clause"
	"gorm.io/gorm"
)

// region Base Model Define

type ID interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

// IModel ModelT Interface
type IModel[T ID] interface {
	GetID() T
	IsNil() bool
}

// ModelT Define

type ModelT[T ID] struct {
	ID        T              `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetID Get ID
func (m ModelT[T]) GetID() T {
	return m.ID
}

func (m ModelT[T]) IsNil() bool {
	var t T
	if m.ID == t {
		return true
	}
	return false
}

type ModelI64 struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m ModelI64) GetID() int64 {
	return m.ID
}

func (m ModelI64) IsNil() bool {
	return m.ID == 0
}

// endregion Base Model Define

// region Page Define

type Page struct {
	Page   int64 `json:"page" form:"page"`
	Size   int64 `json:"size" form:"size"`
	Total  int64 `json:"total" form:"total"`
	Result []any `json:"result" form:"result"`
}

type PageT[T any] struct {
	Page   int64 `json:"page" form:"page"`
	Size   int64 `json:"size" form:"size"`
	Total  int64 `json:"total" form:"total"`
	Result []T   `json:"result" form:"result"`
}

// ToPageT 将非泛型分页转换成泛型分页，适用于泛型处理
func (p *Page) ToPageT() *PageT[any] {
	var res []any
	for _, item := range p.Result {
		res = append(res, item)
	}
	return &PageT[any]{
		Page:   p.Page,
		Size:   p.Size,
		Total:  p.Total,
		Result: res,
	}
}

// ToPage 将泛型分页转换为非泛型分页，方便统一JSON输出或兼容旧接口
func (p PageT[T]) ToPage() *Page {
	return &Page{
		Page:   p.Page,
		Size:   p.Size,
		Total:  p.Total,
		Result: ToAnySlice(p.Result),
	}
}

// ToAnySlice 辅助函数，将任意类型切片转换成 []any，方便赋值给Page.Result
func ToAnySlice[T any](src []T) []any {
	res := make([]any, len(src))
	for i, v := range src {
		res[i] = v
	}
	return res
}

// endregion Page Define

// region Base Repository

// IRepo Repo Interface T: Model Entity K: Model Primary Key Type
type IRepo[T IModel[K], K ID] interface {
	DB() *gorm.DB
	// Tx 使用现有事务或开启新事务
	Tx() IRepo[T, K]
	// Begin 开启新事务
	Begin() IRepo[T, K]
	// Commit 提交事务
	Commit() error
	// Rollback 回滚事务
	Rollback() error
	// Clone 克隆当前Repo实例
	Clone() IRepo[T, K]
	// WithCtx 设置上下文
	WithCtx(*context.Context) IRepo[T, K]
	// WithDB 使用自定义DB连接
	WithDB(*gorm.DB) IRepo[T, K]
	// Raw 执行原生SQL
	Raw(string, ...any) IRawQueryRepo[T, K]
	// Exec 执行原生SQL命令
	Exec(string, ...any) (int64, error)

	// Get 匹配获取新纪录, 不存在返回nil
	Get() (*T, error)
	// GetByID 根据ID获取记录, 不存在返回nil
	GetByID(K) (*T, error)
	// GetOrInit 获取单挑记录, 不存在返回空记录
	GetOrInit() (*T, error)
	// List 查询列表数据
	List() ([]T, error)
	// Page 分页查询（非泛型）
	Page() (*Page, error)
	PageM(size int64, page int64) (*Page, error)
	// PageT 泛型分页查询
	PageT() (*PageT[T], error)
	PageMT(size int64, page int64) (*PageT[T], error)
	WithPage(page *Page) IRepo[T, K]
	// Count 统计数量
	Count() (int64, error)
	// Scan 扫描结果到目标对象
	Scan(dest any) error

	// Set 赋值字段
	Set(field string, val any) IRepo[T, K]
	// SetMap 根据Map设置字段
	SetMap(map[string]any) IRepo[T, K]
	// Create 新增单条记录
	Create(*T) (K, error)
	// CreateBatch  批量新增, 返回插入数量
	CreateBatch([]*T) (int64, error)
	// Save 保存: 存在即更新，不存在即新增, 返回id
	Save(*T) (K, error)
	// Update 根据传入参数更新记录, 返回更新数量
	Update() (int64, error)
	// UpdateFull 全量更新, 返回更新数量
	UpdateFull(*T) (int64, error)
	// Del 删除记录, 返回更新数量
	Del() (int64, error)

	// Desc 降序排序
	Desc(string) IRepo[T, K]
	// Asc 升序排序
	Asc(string) IRepo[T, K]
	// Omit 排除字段
	Omit(...string) IRepo[T, K]

	// Eq 等于
	Eq(string, any) IRepo[T, K]
	// NEq 不等于
	NEq(string, any) IRepo[T, K]
	// In 包含
	In(string, any) IRepo[T, K]
	// Gte 大于等于
	Gte(string, any) IRepo[T, K]
	// Gt 大于
	Gt(string, any) IRepo[T, K]
	// Lt 小于
	Lt(string, any) IRepo[T, K]
	// Lte 小于等于
	Lte(string, any) IRepo[T, K]
	// NotNull 字段不为空
	NotNull(string) IRepo[T, K]
	// Null 字段为空
	Null(string) IRepo[T, K]
	// Like 模糊匹配
	Like(string, any) IRepo[T, K]
	// Select 指定查询字段
	Select(...string) IRepo[T, K]
	// Where 自定义条件
	Where(string, ...any) IRepo[T, K]
	// Limit 限制条数
	Limit(int64) IRepo[T, K]
}

type IRawQueryRepo[T IModel[K], K ID] interface {
	// Get 查询单条数据，查询不到返回nil
	Get() (*T, error)
	// GetOrInit 查询或初始化对象
	GetOrInit() (*T, error)
	// List 查询列表数据
	List() ([]T, error)
	// Scan 扫描结果到目标对象
	Scan(dest any) error
}

//endregion Base Repository

// region IRepo Bases Impl

type RepoTemplate[T IModel[K], K ID] struct {
	ctx   *context.Context
	table string
	model *T
	key   string
	cfg   *RepoCfg
}

type Repo[T IModel[K], K ID] struct {
	*RepoTemplate[T, K]
	db      *gorm.DB
	selects []string
	omits   []string
	match   clause.Match
	page    *Page
	limit   int64
	isRaw   bool
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

func (r *Repo[T, K]) Commit() error {
	newRepo := r.cloneInternal()
	db := newRepo.db.Commit()
	return db.Error
}

func (r *Repo[T, K]) Rollback() error {
	newRepo := r.cloneInternal()
	return newRepo.db.Rollback().Error
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

// WithDB 设置新的 *gorm.DB，返回当前实例
func (r *Repo[T, K]) WithDB(db *gorm.DB) IRepo[T, K] {
	if db == nil {
		panic("db can not be nil")
	}
	r.db = db
	return r
}

// endregion IRepo Bases Impl

// region IRepo Clauses Impl

// Set 赋值字段
func (r *Repo[T, K]) Set(field string, val any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Set(field, val)
	return newR
}

// SetMap 根据Map设置字段
func (r *Repo[T, K]) SetMap(m map[string]any) IRepo[T, K] {
	newR := r.cloneInternal()
	for key, value := range m {
		newR.match.Set(key, value)
	}
	return newR
}

func (r *Repo[T, K]) Where(s string, a ...any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.db = newR.db.Where(s, a...)
	newR.isRaw = true
	return newR
}

func (r *Repo[T, K]) Desc(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Desc(s)
	return newR
}

func (r *Repo[T, K]) Asc(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Asc(s)
	return newR
}

func (r *Repo[T, K]) Omit(s ...string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.omits = append(newR.omits, s...)
	return newR
}

func (r *Repo[T, K]) Eq(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Eq(s, a)
	return newR
}

func (r *Repo[T, K]) NEq(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.NEq(s, a)
	return newR
}

func (r *Repo[T, K]) In(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.In(s, a)
	return newR
}

func (r *Repo[T, K]) Gte(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Gte(s, a)
	return newR
}

func (r *Repo[T, K]) Gt(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Gt(s, a)
	return newR
}

func (r *Repo[T, K]) Lt(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Lt(s, a)
	return newR
}

func (r *Repo[T, K]) Lte(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Lte(s, a)
	return newR
}

func (r *Repo[T, K]) NotNull(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.NotNull(s)
	return newR
}

func (r *Repo[T, K]) Null(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Null(s)
	return newR
}

func (r *Repo[T, K]) Like(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Like(s, a)
	return newR
}

func (r *Repo[T, K]) Select(s ...string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.selects = append(newR.selects, s...)
	return newR
}

func (r *Repo[T, K]) Limit(i int64) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.limit = i
	return newR
}

// endregion IRepo Clauses Impl

// region IRepo Operators Impl

// Exec 执行原生 SQL 写操作（Insert/Update/Delete）
func (r *Repo[T, K]) Exec(sql string, args ...any) (int64, error) {
	newRepo := r.cloneInternal()
	tx := newRepo.db.Exec(sql, args...)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// Create 插入单条数据
func (r *Repo[T, K]) Create(t *T) (K, error) {
	var k K
	if t == nil {
		return k, fmt.Errorf("t is nil")
	}
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(t).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Create(t).Error
	if err != nil {
		return k, err
	}
	return (*t).GetID(), err
}

// CreateBatch 批量插入
func (r *Repo[T, K]) CreateBatch(ts []*T) (int64, error) {
	newRepo := r.cloneInternal()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)
	sql, args := newRepo.match.WhereSql()
	if sql != "" {
		db = db.Where(sql, args...)
	}

	err := db.CreateInBatches(ts, 1000).Error
	if err != nil {
		return 0, err
	}

	return db.RowsAffected, nil
}

// Save 根据 ID 存在与否执行 Create 或 Updated
func (r *Repo[T, K]) Save(t *T) (K, error) {
	var zeroK K
	if t == nil {
		return zeroK, fmt.Errorf("save param t cannot be nil")
	}
	if (*t).IsNil() {
		return r.Create(t)
	} else {
		_, err := r.UpdateFull(t)
		if err != nil {
			return zeroK, err
		}
		return (*t).GetID(), nil
	}
}

// Update 部分字段更新
func (r *Repo[T, K]) Update() (int64, error) {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	updateMap := newRepo.match.SetMap()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)

	if sql != "" {
		db = db.Where(sql, args...)
	}

	result := db.Updates(updateMap)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateFull 用结构体全字段更新
func (r *Repo[T, K]) UpdateFull(t *T) (int64, error) {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(t).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	result := db.Updates(t)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// Del 删除
func (r *Repo[T, K]) Del() (int64, error) {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	if sql == "" {
		return 0, fmt.Errorf("delete operation requires a condition")
	}
	db := newRepo.db.Model(new(T)).Where(sql, args...)
	result := db.Delete(new(T))
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// endregion IRepo Operators Impl

// region IRepo Query Impl
func (r *Repo[T, K]) supportQuery() *Repo[T, K] {
	c := r.cloneInternal()
	c.db = c.db.Model(new(T))
	if c.isRaw {
		return c
	}
	db := c.db.Select(c.selects).Omit(c.omits...)
	sql, args := c.match.WhereSql()
	if sql != "" {
		db = db.Where(sql, args...)
	}
	orders := c.match.OrderSql()
	if orders != "" {
		db = db.Order(orders)
	}
	if c.limit > 0 {
		db = db.Limit(int(c.limit))
	}
	c.db = db
	return c
}

// Raw 执行原生 SQL 查询
func (r *Repo[T, K]) Raw(sql string, args ...any) IRawQueryRepo[T, K] {
	newRepo := r.cloneInternal()
	// 绑定模型 T 到原生查询
	newRepo.db = newRepo.db.Model(new(T)).Raw(sql, args...)
	newRepo.isRaw = true
	return newRepo
}

func (r *Repo[T, K]) Get() (*T, error) {
	c := r.supportQuery()
	var models []T
	if c.isRaw {
		err := c.db.Scan(&models).Error
		if err != nil {
			return nil, err
		} else if len(models) == 0 {
			return nil, nil
		} else if len(models) > 1 {
			return nil, errors.New(fmt.Sprintf("%s: raw query found more than one record", c.key))
		}
		return &models[0], nil
	}
	err := c.db.Find(&models).Error

	if err != nil {
		return nil, err
	} else if len(models) == 0 {
		return nil, nil
	} else if len(models) == 1 {
		return &models[0], nil
	}
	return nil, errors.New(fmt.Sprintf("%s: raw query found more than one record", c.key))
}

func (r *Repo[T, K]) GetByID(id K) (*T, error) {
	return r.Eq("id", id).Get()
}

func (r *Repo[T, K]) GetOrInit() (*T, error) {
	newRepo := r.cloneInternal()
	var empty T
	model, err := newRepo.Get()
	if err != nil {
		return nil, err
	}
	if model != nil {
		return model, nil
	}
	return &empty, nil
}

func (r *Repo[T, K]) List() ([]T, error) {
	var list []T
	if r.isRaw {
		err := r.db.Scan(&list).Error
		if err != nil {
			return nil, err
		}
		return list, nil
	}
	newRepo := r.supportQuery()
	err := newRepo.db.Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func (r *Repo[T, K]) WithPage(page *Page) IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.page = page
	return newRepo
}

func (r *Repo[T, K]) Page() (*Page, error) {
	newRepo := r.supportQuery()
	if newRepo.page == nil {
		newRepo.page = &Page{Page: 1, Size: 10}
	}
	offset := (newRepo.page.Page - 1) * newRepo.page.Size
	var list []T
	var count int64
	err := newRepo.db.Count(&count).Error
	if err != nil {
		return &Page{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: 0, Result: nil}, err
	}

	err = newRepo.db.Offset(int(offset)).Limit(int(newRepo.page.Size)).Find(&list).Error
	if err != nil {
		return &Page{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: count, Result: nil}, nil
	}
	return &Page{
		Page:   newRepo.page.Page,
		Size:   newRepo.page.Size,
		Total:  count,
		Result: ToAnySlice(list),
	}, nil
}
func (r *Repo[T, K]) PageM(size int64, page int64) (*Page, error) {
	return r.WithPage(&Page{Size: size, Page: page}).Page()
}

func (r *Repo[T, K]) PageT() (*PageT[T], error) {
	newRepo := r.supportQuery()
	if newRepo.page == nil {
		newRepo.page = &Page{Page: 1, Size: 10}
	}
	var list []T
	var count int64

	err := newRepo.db.Count(&count).Error
	if err != nil {
		return &PageT[T]{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: 0, Result: nil}, err
	}

	offset := (newRepo.page.Page - 1) * newRepo.page.Size
	err = newRepo.db.Offset(int(offset)).Limit(int(newRepo.page.Size)).Find(&list).Error
	if err != nil {
		return &PageT[T]{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: count, Result: nil}, err
	}

	return &PageT[T]{
		Page:   newRepo.page.Page,
		Size:   newRepo.page.Size,
		Total:  count,
		Result: list,
	}, nil
}
func (r *Repo[T, K]) PageMT(size int64, page int64) (*PageT[T], error) {
	return r.WithPage(&Page{Size: size, Page: page}).PageT()
}

func (r *Repo[T, K]) Count() (int64, error) {
	newRepo := r.supportQuery()
	var count int64
	err := newRepo.DB().Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

// Scan 扫描结果到目标对象
func (r *Repo[T, K]) Scan(dest any) error {
	newRepo := r.cloneInternal()
	err := newRepo.db.Scan(dest).Error
	if err != nil {
		return err
	}
	return nil
}

// endregion IRepo Query Impl
