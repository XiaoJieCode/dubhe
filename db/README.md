# 📘 Repo 工具类使用文档

> 本文档介绍 `Repo` 工具类的完整用法，涵盖初始化、CRUD、分页、条件查询、事务控制、原生 SQL 等场景。  
> 基于 `GORM` 封装，支持泛型与统一的接口抽象。

---

## 1️⃣ 基础概念

### 📌 模型接口定义
所有模型必须实现 `IModel` 接口（提供主键获取与空对象判断）：

````go
type IModel[T ID] interface {
    GetID() T
    IsNil() bool
}

````

可使用内置的基础模型：

```go
type User struct {
    db.ModelI64          // 提供 ID、CreatedAt、UpdatedAt、DeletedAt
    Name string
    Age  int
}
```

---

## 2️⃣ Repo 初始化

```go
// 创建 Repo 实例（泛型）
repo := NewRepo[User, int64]()
```

- `User`：模型类型
- `int64`：主键类型

---

## 3️⃣ 基础方法

### 获取 DB 实例
```go
db := repo.DB()
```

### 克隆 Repo
```go
newRepo := repo.Clone()
```

### 设置上下文
```go
ctx := context.Background()
repo = repo.WithCtx(&ctx)
```

### 切换数据库连接
```go
repo = repo.WithDB(customDB)
```

---

## 4️⃣ 事务控制

```go
// 开启事务
txRepo := repo.Begin()

// 提交事务
err := txRepo.Commit()

// 回滚事务
err := txRepo.Rollback()

// 或直接使用 Tx() 创建新的事务性 Repo
txRepo := repo.Tx()
```

---

## 5️⃣ CRUD 操作

### Create 新增
```go
user := &User{Name: "Tom", Age: 20}
id, err := repo.Create(user)
```

### CreateBatch 批量新增
```go
users := []*User{
    {Name: "Tom", Age: 20},
    {Name: "Jack", Age: 25},
}
count, err := repo.CreateBatch(users)
```

### Save 新增或更新
```go
user := &User{Name: "Tom", Age: 20}
id, err := repo.Save(user)
```

### Update 部分字段更新
```go
rows, err := repo.Eq("id", 1).Set("name", "Updated").Update()
```

### UpdateFull 全量更新
```go
user := &User{ID: 1, Name: "FullUpdate", Age: 30}
rows, err := repo.UpdateFull(&user)
```

### Del 删除
```go
rows, err := repo.Eq("id", 1).Del()
```

---

## 6️⃣ 查询操作

### Get 获取单条
```go
user, err := repo.Eq("id", 1).Get()
```

### GetByID
```go
user, err := repo.GetByID(1)
```

### GetOrInit
```go
user, err := repo.Eq("name", "Tom").GetOrInit()
```

### List 获取多条
```go
users, err := repo.List()
```

### Count 统计数量
```go
count, err := repo.Eq("age", 20).Count()
```

### Scan 扫描到自定义对象
```go
var results []map[string]any
err := repo.Scan(&results)
```

---

## 7️⃣ 分页查询

### Page（非泛型）
```go
page, err := repo.WithPage(&Page{Page: 1, Size: 10}).Page()
```

### PageM 简便方式
```go
page, err := repo.PageM(10, 1)
```

### PageT（泛型分页）
```go
pageT, err := repo.PageT()
```

### PageMT 简便泛型分页
```go
pageT, err := repo.PageMT(10, 1)
```

---

## 8️⃣ 条件构造

所有条件方法均返回新的 Repo，可链式调用。

```go
repo = repo.
    Eq("age", 20).
    Gte("id", 10).
    Like("name", "%Tom%").
    In("status", []int{1, 2, 3}).
    NotNull("email").
    Desc("created_at").
    Limit(5)
```

支持方法：
- `Eq` / `NEq` / `In`
- `Gte` / `Gt` / `Lte` / `Lt`
- `NotNull` / `Null`
- `Like`
- `Asc` / `Desc`
- `Limit`
- `Select` / `Omit`
- `Where`（自定义条件）

---

## 9️⃣ 原生 SQL

### Exec 执行 DML
```go
rows, err := repo.Exec("UPDATE users SET age = ? WHERE id = ?", 30, 1)
```

### Raw 查询单条
```go
user, err := repo.Raw("SELECT * FROM users WHERE id = ?", 1).Get()
```

### Raw 查询列表
```go
users, err := repo.Raw("SELECT * FROM users WHERE age > ?", 20).List()
```

### Raw 分页查询
```go
page, err := repo.Raw("SELECT * FROM users").WithPage(&Page{Page: 1, Size: 5}).Page()
```

### Raw 扫描到结构体
```go
var data []User
err := repo.Raw("SELECT * FROM users WHERE age > ?", 20).Scan(&data)
```

---

## 🔟 错误与异常处理

- `Create(t *T)`：若传入 `nil`，返回 `t is nil` 错误
- `Save(t *T)`：若 `nil`，返回 `save param t cannot be nil`
- `Del()`：若没有条件，返回 `delete operation requires a condition`
- `Get()`：若查询返回多条记录，报错 `"raw query found more than one record"`

---

## 1️⃣1️⃣ 高级用法

### 设置多个字段
```go
repo.Eq("id", 1).SetMap(map[string]any{
"name": "Updated",
"age":  30,
}).Update()
```

### 指定查询字段
```go
users, _ := repo.Select("id", "name").List()
```

### 排除字段
```go
users, _ := repo.Omit("password").List()
```

### 组合条件查询
```go
users, _ := repo.Where("age > ? AND name LIKE ?", 18, "%Tom%").List()
```

### 使用事务提交多个操作
```go
tx := repo.Begin()
defer tx.Rollback()

id, _ := tx.Create(&User{Name: "Alice"})
rows, _ := tx.Eq("id", id).Set("age", 30).Update()

tx.Commit()
```

---

## ✅ 最佳实践

1. **强制条件删除**
    - `Del()` 必须有条件，避免误删全表。

2. **事务优先**
    - 复杂写操作使用 `repo.Begin()` + `Commit()/Rollback()` 保证一致性。

3. **链式调用统一风格**  
   ```go
   repo.Eq("status", 1).Desc("created_at").Limit(10).List()
   ```

4. **分页统一**
    - 封装统一的分页响应结构（`Page` / `PageT`），方便前端消费。  
