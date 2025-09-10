package test

import (
	"context"
	"dubhe/db"
	"dubhe/test/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var testCtx = context.Background()

func setupRepo(t *testing.T) db.IRepo[model.Demo, int64] {
	gdb, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	gdb.AutoMigrate(&model.Demo{})

	repo := db.NewRepo[model.Demo]().WithDB(gdb).WithCtx(&testCtx)
	return repo
}

func prepareData(repo db.IRepo[model.Demo, int64]) {
	for i := 1; i <= 10; i++ {
		repo.Create(&model.Demo{
			Name: "name" + string(rune(i)),
			Age:  i,
		})
	}
}

func TestRepo_Get(t *testing.T) {
	repo := setupRepo(t)
	repo.Create(&model.Demo{Name: "test", Age: 18})
	got, e := repo.Eq("name", "test").Get()
	assert.Nil(t, e)
	if got == nil || got.Name != "test" {
		t.Errorf("Get failed, got: %v", got)
	}
}

func TestRepo_GetOrInit(t *testing.T) {
	repo := setupRepo(t)
	obj, e := repo.Eq("name", "nonexistent").GetOrInit()
	assert.Nil(t, e)
	if obj.Name != "" {
		t.Errorf("Expected new instance, got: %+v", obj)
	}
}

func TestRepo_List(t *testing.T) {
	repo := setupRepo(t)
	prepareData(repo)
	list, e := repo.List()
	assert.Nil(t, e)
	if len(list) != 10 {
		t.Errorf("Expected 10 items, got %d", len(list))
	}
}

func TestRepo_Page(t *testing.T) {
	repo := setupRepo(t)
	prepareData(repo)
	page := &db.Page{
		Page: 1,
		Size: 5,
	}
	paged, e := repo.WithPage(page).Page()
	assert.Nil(t, e)
	if paged.Total != 10 || len(paged.Result) != 5 {
		t.Errorf("Page result incorrect: %+v", paged)
	}
}

func TestRepo_PageT(t *testing.T) {
	repo := setupRepo(t)
	prepareData(repo)
	page := &db.Page{
		Page: 2,
		Size: 4,
	}
	paged, e := repo.WithPage(page).PageT()
	if e != nil {
		t.Errorf("PageT error: %v", e)
	}
	if paged.Total != 10 || len(paged.Result) != 4 {
		t.Errorf("PageT result incorrect: %+v", paged)
	}
}

func TestRepo_Count(t *testing.T) {
	repo := setupRepo(t)
	prepareData(repo)
	count, e := repo.Count()
	if e != nil {
		t.Errorf("Count error: %v", e)
	}
	if count != 10 {
		t.Errorf("Expected count 10, got %d", count)
	}
}

func TestRepo_Scan(t *testing.T) {
	repo := setupRepo(t)
	repo.Create(&model.Demo{Name: "scan-test", Age: 99})
	var out model.Demo
	repo.Raw(`select * from demo where name=?`, "scan-test").Scan(&out)
	if out.Name != "scan-test" {
		t.Errorf("Scan failed, got %+v", out)
	}
}

func TestRepo_Raw(t *testing.T) {
	repo := setupRepo(t)
	m := new(model.Demo).TableName()
	affected, e := repo.Exec("INSERT INTO "+m+" (name, age) VALUES (?, ?)", "raw-user", 33)
	assert.Nil(t, e)
	assert.Equal(t, affected, int64(1))
	var result model.Demo
	repo.Raw("SELECT * FROM "+m+" WHERE name = ?", "raw-user").Scan(&result)
	if result.Name != "raw-user" {
		t.Errorf("Raw query failed: %+v", result)
	}
}
