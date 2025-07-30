package db

// Set 赋值字段
func (r *Repo[T]) Set(field string, val any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Set(field, val)
	return newR
}

// SetMap 根据Map设置字段
func (r *Repo[T]) SetMap(m map[string]any) IRepo[T] {
	newR := r.cloneInternal()
	for key, value := range m {
		newR.match.Set(key, value)
	}
	return newR
}

func (r *Repo[T]) Where(s string, a ...any) IRepo[T] {
	newR := r.cloneInternal()
	newR.db = newR.db.Where(s, a...)
	return newR
}

func (r *Repo[T]) Desc(s string) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Desc(s)
	return newR
}

func (r *Repo[T]) Asc(s string) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Asc(s)
	return newR
}

func (r *Repo[T]) Omit(s ...string) IRepo[T] {
	newR := r.cloneInternal()
	newR.omits = append(newR.omits, s...)
	return newR
}

func (r *Repo[T]) Eq(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Eq(s, a)
	return newR
}

func (r *Repo[T]) NEq(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.NEq(s, a)
	return newR
}

func (r *Repo[T]) In(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.In(s, a)
	return newR
}

func (r *Repo[T]) Gte(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Gte(s, a)
	return newR
}

func (r *Repo[T]) Gt(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Gt(s, a)
	return newR
}

func (r *Repo[T]) Lt(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Lt(s, a)
	return newR
}

func (r *Repo[T]) Lte(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Lte(s, a)
	return newR
}

func (r *Repo[T]) NotNull(s string) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.NotNull(s)
	return newR
}

func (r *Repo[T]) Null(s string) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Null(s)
	return newR
}

func (r *Repo[T]) Or(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Or(s, a)
	return newR
}

func (r *Repo[T]) Like(s string, a any) IRepo[T] {
	newR := r.cloneInternal()
	newR.match.Like(s, a)
	return newR
}

func (r *Repo[T]) Select(s ...string) IRepo[T] {
	newR := r.cloneInternal()
	newR.selects = append(newR.selects, s...)
	return newR
}

func (r *Repo[T]) Limit(i int64) IRepo[T] {
	newR := r.cloneInternal()
	newR.limit = i
	return newR
}
