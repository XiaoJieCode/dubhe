package db

import (
	"errors"
	"gorm.io/gorm"
)

func (r *Repo[T]) Get() *T {
	var model T
	sql, args := r.match.WhereSql()
	db := r.db.Model(&model)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		r.err = err
		r.checkErr(err)
		return nil
	}

	return &model
}

func (r *Repo[T]) GetOrInit() *T {
	model := r.Get()
	if model != nil {
		return model
	}
	// 返回零值指针
	var empty T
	return &empty
}

func (r *Repo[T]) List() []*T {
	var list []*T
	sql, args := r.match.WhereSql()
	db := r.db.Model(new(T))
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Find(&list).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return nil
	}
	return list
}
func (r *Repo[T]) WithPage(page *Page) IRepo[T] {
	r.page = page
	return r
}
func (r *Repo[T]) Page() *Page {
	var list []any // 作为中间容器
	var count int64

	sql, args := r.match.WhereSql()
	db := r.db.Model(new(T))

	// 获取总数
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return &Page{Page: r.page.Page, Size: r.page.Size, Total: 0, Result: nil, Last: 0}
	}

	// 分页参数处理
	offset := (r.page.Page - 1) * r.page.Size
	err = db.Offset(offset).Limit(r.page.Size).Find(&list).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return &Page{Page: r.page.Page, Size: r.page.Size, Total: count, Result: nil, Last: 0}
	}

	// 计算总页数
	last := (count + int64(r.page.Size) - 1) / int64(r.page.Size)

	return &Page{
		Page:   r.page.Page,
		Size:   r.page.Size,
		Total:  count,
		Result: list,
		Last:   last,
	}
}

func (r *Repo[T]) PageT() *PageT[T] {
	var list []*T
	var count int64

	sql, args := r.match.WhereSql()
	db := r.db.Model(new(T))

	// 获取总数
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return &PageT[T]{Page: r.page.Page, Size: r.page.Size, Total: 0, Result: nil, Last: 0}
	}

	// 分页数据
	offset := (r.page.Page - 1) * r.page.Size
	err = db.Offset(offset).Limit(r.page.Size).Find(&list).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return &PageT[T]{Page: r.page.Page, Size: r.page.Size, Total: count, Result: nil, Last: 0}
	}

	last := (count + int64(r.page.Size) - 1) / int64(r.page.Size)

	return &PageT[T]{
		Page:   r.page.Page,
		Size:   r.page.Size,
		Total:  count,
		Result: list,
		Last:   last,
	}
}

func (r *Repo[T]) Count() int64 {
	var count int64
	sql, args := r.match.WhereSql()
	db := r.db.Model(new(T))
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
		return 0
	}
	return count
}

func (r *Repo[T]) Scan(dest any) {
	sql, args := r.match.WhereSql()
	db := r.db.Model(new(T))
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Scan(dest).Error
	if err != nil {
		r.err = err
		r.checkErr(err)
	}
}
