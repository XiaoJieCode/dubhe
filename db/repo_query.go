package db

import (
	"errors"
	"gorm.io/gorm"
)

// Raw 执行原生 SQL 查询
func (r *Repo[T]) Raw(sql string, args ...any) IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.db = newRepo.db.Raw(sql, args...)
	newRepo.isRaw = true
	return newRepo
}

func (r *Repo[T]) Get() *T {
	newRepo := r.cloneInternal()

	var model T
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(&model).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.First(&model).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		newRepo.err = err
		newRepo.checkErr(err)
		return nil
	}

	return &model
}

func (r *Repo[T]) GetOrInit() *T {
	newRepo := r.cloneInternal()
	model := newRepo.Get()
	if model != nil {
		return model
	}
	var empty T
	return &empty
}

func (r *Repo[T]) List() []*T {
	newRepo := r.cloneInternal()
	var list []*T
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}

	var err error
	if newRepo.limit != 0 {
		err = db.Limit(int(newRepo.limit)).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return nil
	}
	return list
}

func (r *Repo[T]) WithPage(page *Page) IRepo[T] {
	newRepo := r.cloneInternal()
	newRepo.page = page
	return newRepo
}

func (r *Repo[T]) Page() *Page {
	newRepo := r.cloneInternal()
	var list []*T
	var count int64

	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)

	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return &Page{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: 0, Result: nil, Last: 0}
	}

	offset := (newRepo.page.Page - 1) * newRepo.page.Size
	err = db.Offset(offset).Limit(newRepo.page.Size).Find(&list).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return &Page{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: count, Result: nil, Last: 0}
	}

	last := (count + int64(newRepo.page.Size) - 1) / int64(newRepo.page.Size)

	return &Page{
		Page:   newRepo.page.Page,
		Size:   newRepo.page.Size,
		Total:  count,
		Result: list,
		Last:   last,
	}
}

func (r *Repo[T]) PageT() *PageT[T] {
	newRepo := r.cloneInternal()
	var list []*T
	var count int64

	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)

	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return &PageT[T]{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: 0, Result: nil, Last: 0}
	}

	offset := (newRepo.page.Page - 1) * newRepo.page.Size
	err = db.Offset(offset).Limit(newRepo.page.Size).Find(&list).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return &PageT[T]{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: count, Result: nil, Last: 0}
	}

	last := (count + int64(newRepo.page.Size) - 1) / int64(newRepo.page.Size)

	return &PageT[T]{
		Page:   newRepo.page.Page,
		Size:   newRepo.page.Size,
		Total:  count,
		Result: list,
		Last:   last,
	}
}

func (r *Repo[T]) Count() int64 {
	newRepo := r.cloneInternal()
	var count int64
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T))
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
		return 0
	}
	return count
}

func (r *Repo[T]) Scan(dest any) {
	newRepo := r.cloneInternal()
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Scan(dest).Error
	if err != nil {
		newRepo.err = err
		newRepo.checkErr(err)
	}
}
