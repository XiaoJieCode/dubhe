package db

// Set 赋值字段
func (r *Repo[T]) Set(field string, val any) IRepo[T] {
	r.match.Set(field, val)
	return r
}

// SetMap 根据Map设置字段
func (r *Repo[T]) SetMap(m map[string]any) IRepo[T] {
	for key, value := range m {
		r.match.Set(key, value)
	}
	return r
}
func (r *Repo[T]) Where(s string, a ...any) IRepo[T] {
	r.DB = r.DB.Where(s, a)
	return r
}
func (r *Repo[T]) Desc(s string) IRepo[T] {
	r.match.Desc(s)
	return r
}

func (r *Repo[T]) Asc(s string) IRepo[T] {
	r.match.Asc(s)
	return r
}

func (r *Repo[T]) Omit(s ...string) IRepo[T] {
	r.omits = append(r.omits, s...)
	return r
}

func (r *Repo[T]) Eq(s string, a any) IRepo[T] {
	r.match.Eq(s, a)
	return r
}

func (r *Repo[T]) NEq(s string, a any) IRepo[T] {
	r.match.NEq(s, a)
	return r

}

func (r *Repo[T]) In(s string, a any) IRepo[T] {
	r.match.In(s, a)
	return r
}

func (r *Repo[T]) Gte(s string, a any) IRepo[T] {
	r.match.Gte(s, a)
	return r

}

func (r *Repo[T]) Gt(s string, a any) IRepo[T] {
	r.match.Gt(s, a)
	return r
}

func (r *Repo[T]) Lt(s string, a any) IRepo[T] {
	r.match.Lt(s, a)
	return r
}

func (r *Repo[T]) Lte(s string, a any) IRepo[T] {
	r.match.Lte(s, s)
	return r
}

func (r *Repo[T]) NotNull(s string) IRepo[T] {
	r.match.NotNull(s)
	return r
}

func (r *Repo[T]) Null(s string) IRepo[T] {
	r.match.Null(s)
	return r
}

func (r *Repo[T]) Or(s string, a ...any) IRepo[T] {
	r.match.Or(s, a)
	return r
}
func (r *Repo[T]) Like(s string, a any) IRepo[T] {
	r.match.Like(s, a)
	return r
}

func (r *Repo[T]) Select(s ...string) IRepo[T] {
	r.selects = append(r.selects, s...)
	return r
}

func (r *Repo[T]) Limit(i int64) IRepo[T] {
	r.limit = i
	return r
}
