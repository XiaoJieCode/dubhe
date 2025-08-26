package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// Raw 执行原生 SQL 查询
func (r *Repo[T, K]) Raw(sql string, args ...any) IRawQueryRepo[T, K] {
	newRepo := r.cloneInternal()
	// 绑定模型 T 到原生查询
	newRepo.db = newRepo.db.Model(new(T)).Raw(sql, args...)
	newRepo.isRaw = true
	return newRepo
}

func (r *Repo[T, K]) Get() *T {
	newRepo := r.cloneInternal()
	var models []T

	if r.isRaw {
		r.db.Scan(&models)
		newRepo.handleErr(r.db.Error)
		if len(models) == 0 {
			return nil
		}
		if len(models) > 1 {
			err := errors.New(fmt.Sprintf("%s: raw query found more than one record", newRepo.key))
			newRepo.handleErr(err)
			return nil
		}
		return &models[0]
	}

	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...).Order(newRepo.match.OrderSql())
	}

	err := db.Find(&models).Error
	if err != nil {
		newRepo.handleErr(err)
		return nil
	}

	switch len(models) {
	case 0:
		return nil
	case 1:
		return &models[0]
	default:
		newRepo.err = errors.New(fmt.Sprintf("%s: found more than one record", newRepo.key))
		newRepo.handleErr(newRepo.err)
		return nil
	}
}

func (r *Repo[T, K]) GetByID(id K) *T {
	return r.Eq("id", id).Get()
}

func (r *Repo[T, K]) GetOrInit() *T {
	newRepo := r.cloneInternal()
	var empty T
	model := newRepo.Get()
	if model != nil {
		return model
	}
	return &empty
}

func (r *Repo[T, K]) Take() *T {
	newRepo := r.cloneInternal()
	var model T

	if r.isRaw {
		err := newRepo.db.Limit(1).Scan(&model).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			newRepo.handleErr(newRepo.db.Error)
		}

		return &model
	}

	db := newRepo.db.Model(&model).Omit(newRepo.omits...).Limit(1)
	if len(newRepo.match.Orders) > 0 {
		db = db.Order(newRepo.match.OrderSql())
	}
	sql, args := newRepo.match.WhereSql()
	if sql != "" {
		db = db.Where(sql, args...)
	}

	err := db.First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		newRepo.handleErr(err)
		return nil
	}

	if model.IsNil() {
		return nil
	}

	return &model
}

func (r *Repo[T, K]) TakeOrInit() *T {
	newRepo := r.cloneInternal()
	var empty T
	model := newRepo.Take()
	if model != nil {
		return model
	}
	return &empty
}

func (r *Repo[T, K]) List() []T {
	newRepo := r.cloneInternal()
	var list []T
	if r.isRaw {
		err := r.db.Scan(&list).Error
		if err != nil {
			newRepo.handleErr(err)
			return nil
		}
		return list
	}
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)
	if sql != "" {
		db = db.Where(sql, args...)
	}
	if newRepo.match.OrderSql() != "" {
		db = db.Order(newRepo.match.OrderSql())
	}
	if newRepo.limit != 0 {
		db = db.Limit(int(newRepo.limit))
	}
	err := db.Find(&list).Error
	if err != nil {
		newRepo.handleErr(err)
		return nil
	}
	return list
}

func (r *Repo[T, K]) WithPage(page *Page) IRepo[T, K] {
	newRepo := r.cloneInternal()
	newRepo.page = page
	return newRepo
}

func (r *Repo[T, K]) Page() *Page {
	newRepo := r.cloneInternal()
	if r.isRaw {
		r.log.Error("Page() is not support for raw query")
		return &Page{}
	}
	var list []T
	var count int64

	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)

	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		newRepo.handleErr(err)
		return &Page{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: 0, Result: nil, Last: 0}
	}

	offset := (newRepo.page.Page - 1) * newRepo.page.Size
	err = db.Offset(offset).Limit(newRepo.page.Size).Find(&list).Error
	if err != nil {
		newRepo.handleErr(err)
		return &Page{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: count, Result: nil, Last: 0}
	}

	last := (count + int64(newRepo.page.Size) - 1) / int64(newRepo.page.Size)

	return &Page{
		Page:   newRepo.page.Page,
		Size:   newRepo.page.Size,
		Total:  count,
		Result: ToAnySlice(list),
		Last:   last,
	}
}

func (r *Repo[T, K]) PageT() *PageT[T] {
	if r.isRaw {
		r.log.Error("PageT() not support for raw query")
		return &PageT[T]{}
	}
	newRepo := r.cloneInternal()
	var list []T
	var count int64

	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T)).Omit(newRepo.omits...)
	orders := newRepo.match.OrderSql()
	if sql != "" {
		db = db.Where(sql, args...)
	}
	if orders != "" {
		db = db.Order(orders)
	}
	err := db.Count(&count).Error
	if err != nil {
		newRepo.handleErr(err)
		return &PageT[T]{Page: newRepo.page.Page, Size: newRepo.page.Size, Total: 0, Result: nil, Last: 0}
	}

	offset := (newRepo.page.Page - 1) * newRepo.page.Size
	err = db.Offset(offset).Limit(newRepo.page.Size).Find(&list).Error
	if err != nil {
		newRepo.handleErr(err)
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

func (r *Repo[T, K]) Count() int64 {
	newRepo := r.cloneInternal()
	var count int64
	sql, args := newRepo.match.WhereSql()
	db := newRepo.db.Model(new(T))
	if sql != "" {
		db = db.Where(sql, args...)
	}
	err := db.Count(&count).Error
	if err != nil {
		newRepo.handleErr(err)
		return 0
	}
	return count
}

// Scan 扫描结果到目标对象
func (r *Repo[T, K]) Scan(dest any) {
	newRepo := r.cloneInternal()

	result := newRepo.db.Scan(dest)

	if result.Error != nil {
		newRepo.err = result.Error
		newRepo.handleErr(result.Error)
	}
}
