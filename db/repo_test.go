package db_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"dubhe/db"
)

// ---------- 测试模型 ----------
type User struct {
	db.ModelI64
	Name  string
	Age   int
	Email string
}

func (User) TableName() string { return "users" }
func (u User) RepoDefine() db.RepoCfg {
	return db.RepoCfg{DB: testDB, AutoMigrate: true}
}

// ---------- 全局测试 DB ----------
var testDB *gorm.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = testDB.AutoMigrate(&User{})
	m.Run()
}

// ---------- 单元测试 ----------

func TestRepoCRUD(t *testing.T) {
	repo := db.NewRepo[User, int64]()

	// Create
	u := &User{Name: "Alice", Age: 20, Email: "a@a.com"}
	id, err := repo.Create(u)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if id == 0 {
		t.Fatal("id not returned")
	}

	// GetByID
	u2, err := repo.GetByID(id)
	if err != nil || u2 == nil || u2.Name != "Alice" {
		t.Fatalf("get by id failed: %v, %+v", err, u2)
	}

	// Update
	repo = repo.Set("age", 25).Eq("id", id)
	affected, err := repo.Update()
	if err != nil || affected != 1 {
		t.Fatalf("update failed: %v, affected=%d", err, affected)
	}

	// Save (Update path)
	u2.Age = 30
	_, err = repo.Save(u2)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// List
	list, err := repo.List()
	if err != nil || len(list) != 1 {
		t.Fatalf("list failed: %v, len=%d", err, len(list))
	}

	// Page
	p, err := repo.PageM(10, 1)
	if err != nil || p.Total != 1 {
		t.Fatalf("page failed: %v, %+v", err, p)
	}

	// Del
	affected, err = repo.Eq("id", id).Del()
	if err != nil || affected != 1 {
		t.Fatalf("delete failed: %v, affected=%d", err, affected)
	}
}

func TestRepoTx(t *testing.T) {
	repo := db.NewRepo[User, int64]()

	txRepo := repo.Begin()
	u := &User{Name: "Bob", Age: 22}
	_, err := txRepo.Create(u)
	if err != nil {
		t.Fatalf("create in tx failed: %v", err)
	}

	// rollback
	_ = txRepo.Rollback()
	got, _ := repo.GetByID(u.ID)
	if got != nil {
		t.Fatal("rollback failed, record should not exist")
	}

	// commit
	txRepo = repo.Begin()
	u2 := &User{Name: "Carol", Age: 23}
	_, _ = txRepo.Create(u2)
	_ = txRepo.Commit()

	got, _ = repo.GetByID(u2.ID)
	if got == nil {
		t.Fatal("commit failed, record not found")
	}
}

func TestRepoRawQuery(t *testing.T) {
	repo := db.NewRepo[User, int64]()

	// Insert test data
	repo.Create(&User{Name: "Dave", Age: 40})

	rawRepo := repo.Raw("SELECT * FROM users WHERE age > ?", 30)
	list, err := rawRepo.List()
	if err != nil || len(list) == 0 {
		t.Fatalf("raw list failed: %v", err)
	}
}
func TestCreateBatch(t *testing.T) {
	repo := db.NewRepo[User, int64]()
	batch, err := repo.CreateBatch([]*User{
		&User{Name: "Dave", Age: 40},
		&User{Name: "Bob", Age: 30},
		&User{Name: "Alice", Age: 20},
	})
	if err != nil || batch != 3 {
		t.Fatalf("create batch failed: %v", err)
	}
}
func TestClauseGet(t *testing.T) {
	repo := db.NewRepo[User, int64]()
	batch, err := repo.CreateBatch([]*User{
		{Name: "Dave", Age: 40},
		{Name: "Bob", Age: 30},
		{Name: "Alice", Age: 20},
	})
	if err != nil || batch != 3 {
		t.Fatalf("create batch failed: %v", err)
	}
	get, err := repo.Eq("name", "Dave").Get()
	if err != nil || get == nil || get.Name != "Dave" {
		t.Fatalf("get failed: %v", err)
	}
	u, err := repo.Gt("age", 20).Desc("age").Limit(1).Get()
	if err != nil || u == nil || u.Name != "Dave" {
		t.Fatalf("get failed: %v", err)
	}

}
