package db

import "fmt"

func (r *Repo[T]) Raw(s string, a ...any) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Exec(s string, a ...any) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Create(t *T) int64 {
	sql, args := r.match.WhereSql()
	db := r.DB.Model(t).Omit(r.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Save(t).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return 0
	}
	return db.RowsAffected
}

func (r *Repo[T]) CreateBatch(ts []*T) int64 {
	var t T
	db := r.DB.Model(&t).Omit(r.omits...)
	sql, args := r.match.WhereSql()
	if sql != "" {
		db = db.Where(sql, args...)
	}

	err := db.CreateInBatches(ts, 1000).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return 0
	}

	return int64(len(ts))
}

func (r *Repo[T]) Save(t *T) int64 {
	model, ok := any(t).(IBaseModel)
	if !ok {
		r.err = fmt.Errorf("type does not implement IBaseModel")
		r.checkErr(r.err)
		return 0
	}

	if model.GetID() == 0 {
		return r.Create(t)
	}

	// 正确传参给 Update
	db := r.DB.Omit(r.omits...)
	db.Statement.Dest = t
	r.DB = db
	return r.Update()
}

func (r *Repo[T]) Update() int64 {
	var t T
	sql, args := r.match.WhereSql()
	updateMap := r.match.SetMap()
	db := r.DB.Model(t).Omit(r.omits...)

	if sql != "" {
		db = db.Where(sql, args...)
	}

	result := db.Updates(updateMap)
	if result.Error != nil {
		r.err = result.Error
		r.checkErr(result.Error)
		return 0
	}
	return result.RowsAffected
}
func (r *Repo[T]) UpdateFull(t *T) int64 {
	sql, args := r.match.WhereSql()
	db := r.DB.Model(t).Omit(r.omits...)

	if sql != "" {
		db = db.Where(sql, args...)
	}

	result := db.Updates(t)
	if result.Error != nil {
		r.err = result.Error
		r.checkErr(result.Error)
		return 0
	}
	return result.RowsAffected
}

func (r *Repo[T]) Del() int64 {
	sql, args := r.match.WhereSql()
	if sql == "" {
		r.err = fmt.Errorf("delete operation requires a condition")
		r.checkErr(r.err)
		return 0
	}
	var t T
	db := r.DB.Where(sql, args...)
	result := db.Delete(&t)
	if result.Error != nil {
		r.err = result.Error
		r.checkErr(result.Error)
		return 0
	}
	return result.RowsAffected
}
