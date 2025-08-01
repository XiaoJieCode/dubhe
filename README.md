# GORM CRUD 工具封装文档

## 概述

本工具封装了 GORM 框架，提供了一套类型安全、链式调用的 CRUD 操作接口。通过泛型支持，实现了更简洁的数据库操作方式，并包含事务管理、错误处理等高级特性。

## 核心接口

### `IRepo[T any]`

泛型仓库接口，`T` 表示模型类型

### `IRepoErrHandle[T any]`

错误处理接口

## 方法分类说明

### 1. 数据库连接与事务管理

| 方法                          | 参数    | 返回值        | 说明           |
|-----------------------------|-------|------------|--------------|
| `DB()`                      | -     | `*gorm.DB` | 获取底层 GORM 连接 |
| `Tx()`                      | -     | `IRepo[T]` | 使用现有事务或开启新事务 |
| `Begin()`                   | -     | `IRepo[T]` | 显式开启新事务      |
| `Commit()`                  | -     | `IRepo[T]` | 提交事务         |
| `Rollback()`                | -     | `IRepo[T]` | 回滚事务         |
| `WithDB(*gorm.DB)`          | 数据库连接 | `IRepo[T]` | 使用自定义 DB 连接  |
| `WithCtx(*context.Context)` | 上下文   | `IRepo[T]` | 设置上下文        |
| `Clone()`                   | -     | `IRepo[T]` | 克隆当前 Repo 实例 |

### 2. 错误处理

| 方法                               | 参数     | 返回值          | 说明         |
|----------------------------------|--------|--------------|------------|
| `OnErr(func(IRepoErrHandle[T]))` | 错误处理函数 | `IRepo[T]`   | 设置自定义错误处理器 |
| `Err()`                          | -      | `*err.Error` | 获取操作错误     |

### 3. 查询操作

| 方法                     | 参数   | 返回值         | 说明                  |
|------------------------|------|-------------|---------------------|
| `Get()`                | -    | `*T`        | 查询单条记录（未找到返回nil）    |
| `GetOrInit()`          | -    | `*T`        | 查询或初始化对象（未找到返回空结构体） |
| `List()`               | -    | `[]T`       | 查询列表数据              |
| `Count()`              | -    | `int64`     | 统计数量                |
| `Page()`               | -    | `*Page`     | 分页查询（返回通用分页对象）      |
| `PageT()`              | -    | `*PageT[T]` | 泛型分页查询（包含类型化数据）     |
| `Scan(dest any)`       | 目标对象 | -           | 扫描结果到指定对象           |
| `WithPage(page *Page)` | 分页对象 | `IRepo[T]`  | 设置分页参数              |

### 4. 写入操作

| 方法                       | 参数      | 返回值        | 说明                   |
|--------------------------|---------|------------|----------------------|
| `Create(*T)`             | 模型指针    | `int64`    | 新增单条记录（返回影响行数）       |
| `CreateBatch([]*T)`      | 模型切片    | `int64`    | 批量新增记录               |
| `Save(*T)`               | 模型指针    | `int64`    | 保存记录（存在更新，不存在新增）     |
| `Update()`               | -       | `int64`    | 执行更新（需先调用Set/SetMap） |
| `UpdateFull(*T)`         | 模型指针    | `int64`    | 完整更新模型               |
| `Del()`                  | -       | `int64`    | 删除记录（需先调用条件方法）       |
| `Set(string, any)`       | 字段名, 值  | `IRepo[T]` | 设置单个更新字段             |
| `SetMap(map[string]any)` | 字段映射    | `IRepo[T]` | 批量设置更新字段             |
| `Exec(string, ...any)`   | SQL, 参数 | `int64`    | 执行原生SQL命令            |
| `Raw(string, ...any)`    | SQL, 参数 | `IRepo[T]` | 设置原生SQL查询            |

### 5. 条件构建

| 方法                      | 参数     | 返回值        | 说明         |
|-------------------------|--------|------------|------------|
| `Eq(string, any)`       | 字段, 值  | `IRepo[T]` | 等于条件       |
| `NEq(string, any)`      | 字段, 值  | `IRepo[T]` | 不等于条件      |
| `In(string, any)`       | 字段, 值  | `IRepo[T]` | IN 条件      |
| `Gt(string, any)`       | 字段, 值  | `IRepo[T]` | 大于条件       |
| `Gte(string, any)`      | 字段, 值  | `IRepo[T]` | 大于等于条件     |
| `Lt(string, any)`       | 字段, 值  | `IRepo[T]` | 小于条件       |
| `Lte(string, any)`      | 字段, 值  | `IRepo[T]` | 小于等于条件     |
| `Like(string, any)`     | 字段, 值  | `IRepo[T]` | 模糊匹配       |
| `NotNull(string)`       | 字段     | `IRepo[T]` | 非NULL条件    |
| `Null(string)`          | 字段     | `IRepo[T]` | NULL条件     |
| `Or(string, any)`       | 字段, 值  | `IRepo[T]` | OR条件       |
| `Where(string, ...any)` | 条件, 参数 | `IRepo[T]` | 自定义WHERE条件 |
| `Select(...string)`     | 字段列表   | `IRepo[T]` | 指定查询字段     |
| `Omit(...string)`       | 字段列表   | `IRepo[T]` | 排除字段       |
| `Desc(string)`          | 字段     | `IRepo[T]` | 降序排序       |
| `Asc(string)`           | 字段     | `IRepo[T]` | 升序排序       |
| `Limit(int64)`          | 数量     | `IRepo[T]` | 限制结果数量     |

## 使用示例

### 基础CRUD操作

```go
// 初始化Repo
userRepo := NewRepo[User]()

// 创建记录
newUser := &User{Name: "John", Age: 30}
userRepo.Create(newUser)

// 查询单条
user := userRepo.Eq("id", 123).Get()

// 更新记录
userRepo.Eq("id", 123).Set("age", 31).Update()

// 删除记录
userRepo.Eq("id", 456).Del()
```

### 事务操作

```go


err := userRepo.Tx(func (txRepo IRepo[User]) error {
// 转出方扣款
if affected := txRepo.Eq("id", 1).Set("balance", gorm.Expr("balance - ?", 100)).Update(); affected == 0 {
return errors.New("扣款失败")
}

// 接收方加款
if affected := txRepo.Eq("id", 2).Set("balance", gorm.Expr("balance + ?", 100)).Update(); affected == 0 {
    return errors.New("加款失败")
}
    
    return nil
})

if err != nil {
// 处理事务错误
}

```

## 设计特点

1. **链式调用**：所有方法返回 `IRepo[T]` 接口，支持链式调用
2. **类型安全**：通过泛型确保模型类型一致性
3. **错误处理**：提供灵活的错误处理机制
4. **事务支持**：简化事务管理流程
5. **条件构建**：提供丰富的条件构建方法

> 注意：实际使用前需初始化全局数据库连接，并注册模型对应的错误处理器