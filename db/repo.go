package db

import (
	"context"
	"dubhe/db/err"
	"gorm.io/gorm"
)

// region Base Repository

type GlobalErrHandle func(h IRepoErrHandleBase)

type IRepoErrHandleBase interface {
	Error() error
	Continue()
	Cancel()
	Panic()
}

type IRepoErrHandle[T IModel[K], K ID] interface {
	IRepoErrHandleBase
	Repo() IRepo[T, K]
}

// IRepo Repo Interface T: Model Entity K: Model Primary Key Type
type IRepo[T IModel[K], K ID] interface {
	DB() *gorm.DB
	// Tx 使用现有事务或开启新事务
	Tx() IRepo[T, K]
	// Begin 开启新事务
	Begin() IRepo[T, K]
	// Commit 提交事务
	Commit() IRepo[T, K]
	// Rollback 回滚事务
	Rollback() IRepo[T, K]
	// OnErr 设置错误处理函数
	OnErr(func(IRepoErrHandle[T, K])) IRepo[T, K]
	// Err 返回Err
	Err() *err.Error
	// Clone 克隆当前Repo实例
	Clone() IRepo[T, K]
	// WithCtx 设置上下文
	WithCtx(*context.Context) IRepo[T, K]
	// WithDB 使用自定义DB连接
	WithDB(*gorm.DB) IRepo[T, K]
	// Raw 执行原生SQL
	Raw(string, ...any) IRawQueryRepo[T, K]
	// Exec 执行原生SQL命令
	Exec(string, ...any) int64

	// Get 严格匹配单条数据,记录不存在返回nil
	Get() *T
	// GetByID 根据ID获取记录
	GetByID(K) *T
	// Take 获取单条数据,记录不存在返回nil
	Take() *T
	// GetOrInit 严格匹配单条数据,记录不存在初始化对象
	GetOrInit() *T
	// TakeOrInit 获取单条数据,记录不存在初始化对象
	TakeOrInit() *T
	// List 查询列表数据
	List() []T
	// Page 分页查询（非泛型）
	Page() *Page
	// PageT 泛型分页查询
	PageT() *PageT[T]
	WithPage(page *Page) IRepo[T, K]
	// Count 统计数量
	Count() int64
	// Scan 扫描结果到目标对象
	Scan(dest any)

	// Set 赋值字段
	Set(field string, val any) IRepo[T, K]
	// SetMap 根据Map设置字段
	SetMap(map[string]any) IRepo[T, K]
	// Create 新增单条记录
	Create(*T) int64
	// CreateBatch  批量新增
	CreateBatch([]*T) int64
	// Save 保存（存在即更新，不存在即新增）
	Save(*T) int64
	// Update 根据传入参数更新记录
	Update() int64
	// UpdateFull 全量更新
	UpdateFull(*T) int64
	// Del 删除记录
	Del() int64

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
	// Or 或条件
	Or(string, any) IRepo[T, K]
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
	Get() *T
	// GetOrInit 查询或初始化对象
	GetOrInit() *T
	// List 查询列表数据
	List() []T
	// Scan 扫描结果到目标对象
	Scan(dest any)
}

//endregion Base Repository
