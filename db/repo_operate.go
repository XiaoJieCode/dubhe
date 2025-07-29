package db

import (
	"fmt"
)

// Exec 执行原生 SQL 写操作（Insert/Update/Delete）
func (r *Repo[T]) Exec(sql string, args ...any) int64 {
	newRepo := r.cloneInternal()
	tx := newRepo.db.Exec(sql, args...)
	if tx.Error != nil {
		newRepo.err = tx.Error
		newRepo.checkErr(tx.Error)
		return 0
	}
	return tx.RowsAffected
}

// Create 插入单条数据
func (r *Repo[T]) Create(t *T) int64 {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(t).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Save(t).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return 0
	}
	return db.RowsAffected
}

// CreateBatch 批量插入
func (r *Repo[T]) CreateBatch(ts []*T) int64 {
	newRepo := r.cloneInternal()
	var t T
	db := newRepo.db.Model(&t).Omit(newRepo.omits...)
	sql, args := newRepo.match.WhereSql()
	if sql != "" {
		db = db.Where(sql, args...)
	}

	err := db.CreateInBatches(ts, 1000).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return 0
	}

	return int64(len(ts))
}

// Save 根据 ID 存在与否执行 Create 或 Update
func (r *Repo[T]) Save(t *T) int64 {
	newRepo := r.cloneInternal()
	model, ok := any(t).(IBaseModel)
	if !ok {
		newRepo.err = fmt.Errorf("type does not implement IBaseModel")
		newRepo.checkErr(newRepo.err)
		return 0
	}
	if model.GetID() == 0 {
		return newRepo.Create(t)
	}
	db := newRepo.db.Omit(newRepo.omits...)
	db.Statement.Dest = t
	newRepo.db = db
	return newRepo.Update()
}

// Update 部分字段更新
func (r *Repo[T]) Update() int64 {
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
		newRepo.checkErr(result.Error)
		return 0
	}
	return result.RowsAffected
}

// UpdateFull 用结构体全字段更新
func (r *Repo[T]) UpdateFull(t *T) int64 {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(t).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	result := db.Updates(t)
	if result.Error != nil {
		newRepo.err = result.Error
		newRepo.checkErr(result.Error)
		return 0
	}
	return result.RowsAffected
}

// Del 删除
func (r *Repo[T]) Del() int64 {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	if sql == "" {
		newRepo.err = fmt.Errorf("delete operation requires a condition")
		newRepo.checkErr(newRepo.err)
		return 0
	}
	var t T
	db := newRepo.db.Where(sql, args...)
	result := db.Delete(&t)
	if result.Error != nil {
		newRepo.err = result.Error
		newRepo.checkErr(result.Error)
		return 0
	}
	return result.RowsAffected
}
