package db

import (
	"context"
	"dubhe/db/err"
	"gorm.io/gorm"
)

// region Base Repository

// IRepo Repo Interface
type IRepo[T any] interface {
	// Tx 使用现有事务或开启新事务
	Tx() IRepo[T]
	// Begin 开启新事务
	Begin() IRepo[T]
	// Commit 提交事务
	Commit() IRepo[T]
	// Rollback 回滚事务
	Rollback() IRepo[T]
	// OnErr 设置错误处理函数
	OnErr(func(IRepo[T], *err.Error) bool) IRepo[T]
	// Clone 克隆当前Repo实例
	Clone() IRepo[T]
	// WithCtx 设置上下文
	WithCtx(*context.Context) IRepo[T]
	// WithDB 使用自定义DB连接
	WithDB(*gorm.DB) IRepo[T]
	// Raw 执行原生SQL
	Raw(string, ...any) IRepo[T]
	// Exec 执行原生SQL命令
	Exec(string, ...any) IRepo[T]

	// Get 查询单条数据, 查询为nil
	Get() *T
	// GetOrInit 查询或初始化对象
	GetOrInit() *T
	// List 查询列表数据
	List() []*T
	// Page 分页查询（非泛型）
	Page() *Page
	// PageT 泛型分页查询
	PageT() *PageT[T]
	WithPage(page *Page) IRepo[T]
	// Count 统计数量
	Count() int64
	// Scan 扫描结果到目标对象
	Scan(dest any)

	// Set 赋值字段
	Set(field string, val any) IRepo[T]
	// SetMap 根据Map设置字段
	SetMap(map[string]any) IRepo[T]
	// Add 新增单条记录
	Create(*T) int64
	// AddBatch 批量新增
	CreateBatch([]*T) int64
	// Save 保存（存在即更新，不存在即新增）
	Save(*T) int64
	// Update 更新记录
	Update() int64
	UpdateFull(*T) int64
	// Del 删除记录
	Del() int64

	// Desc 降序排序
	Desc(string) IRepo[T]
	// Asc 升序排序
	Asc(string) IRepo[T]
	// Omit 排除字段
	Omit(...string) IRepo[T]

	// Eq 等于
	Eq(string, any) IRepo[T]
	// NEq 不等于
	NEq(string, any) IRepo[T]
	// In 包含
	In(string, any) IRepo[T]
	// Gte 大于等于
	Gte(string, any) IRepo[T]
	// Gt 大于
	Gt(string, any) IRepo[T]
	// Lt 小于
	Lt(string, any) IRepo[T]
	// Lte 小于等于
	Lte(string, any) IRepo[T]
	// NotNull 字段不为空
	NotNull(string) IRepo[T]
	// Null 字段为空
	Null(string) IRepo[T]
	// Or 或条件
	Or(string, ...any) IRepo[T]
	// Like 模糊匹配
	Like(string, any) IRepo[T]
	// Select 指定查询字段
	Select(...string) IRepo[T]
	// Where 自定义条件
	Where(string, ...any) IRepo[T]
	// Limit 限制条数
	Limit(int64) IRepo[T]
}

//endregion Base Repository
