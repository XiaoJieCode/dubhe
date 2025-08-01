package db

import (
	"fmt"
)

// Exec 执行原生 SQL 写操作（Insert/Update/Delete）
func (r *Repo[T, K]) Exec(sql string, args ...any) int64 {
	newRepo := r.cloneInternal()
	tx := newRepo.db.Exec(sql, args...)
	if tx.Error != nil {
		newRepo.handleErr(tx.Error)
		return 0
	}
	return tx.RowsAffected
}

// Create 插入单条数据
func (r *Repo[T, K]) Create(t *T) int64 {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(t).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Create(t).Error
	if err != nil {
		newRepo.handleErr(err)
		return 0
	}
	return db.RowsAffected
}

// CreateBatch 批量插入
func (r *Repo[T, K]) CreateBatch(ts []*T) int64 {
	newRepo := r.cloneInternal()
	var t T
	db := newRepo.db.Model(&t).Omit(newRepo.omits...)
	sql, args := newRepo.match.WhereSql()
	if sql != "" {
		db = db.Where(sql, args...)
	}

	err := db.CreateInBatches(ts, 1000).Error
	if err != nil {
		newRepo.handleErr(err)
		return 0
	}

	return int64(len(ts))
}

// Save 根据 ID 存在与否执行 Create 或 Update
func (r *Repo[T, K]) Save(t *T) int64 {
	newRepo := r.cloneInternal()
	model, ok := any(t).(IModel[K])
	if !ok {
		newRepo.handleErr(fmt.Errorf("type does not implement IModel"))
		return 0
	}
	if model.IsNil() {
		return newRepo.Create(t)
	}
	db := newRepo.db.Omit(newRepo.omits...)
	db.Statement.Dest = t
	newRepo.db = db
	return newRepo.UpdateFull(t)
}

// Update 部分字段更新
func (r *Repo[T, K]) Update() int64 {
	newRepo := r.cloneInternal()
	var t T
	sql, args := newRepo.match.WhereSql()
	updateMap := newRepo.match.SetMap()
	db := newRepo.db.Model(t).Omit(newRepo.omits...)

	if sql != "" {
		db = db.Where(sql, args...)
	}

	result := db.Updates(updateMap)
	if result.Error != nil {
		newRepo.err = result.Error
		newRepo.handleErr(result.Error)
		return 0
	}
	return result.RowsAffected
}

// UpdateFull 用结构体全字段更新
func (r *Repo[T, K]) UpdateFull(t *T) int64 {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(t).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	result := db.Updates(t)
	if result.Error != nil {
		newRepo.handleErr(result.Error)
		return 0
	}
	return result.RowsAffected
}

// Del 删除
func (r *Repo[T, K]) Del() int64 {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	if sql == "" {
		newRepo.err = fmt.Errorf("delete operation requires a condition")
		newRepo.handleErr(newRepo.err)
		return 0
	}
	var t T
	db := newRepo.db.Where(sql, args...)
	result := db.Delete(&t)
	if result.Error != nil {
		newRepo.handleErr(result.Error)
		return 0
	}
	return result.RowsAffected
}
