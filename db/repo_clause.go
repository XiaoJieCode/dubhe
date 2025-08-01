package db

// Set 赋值字段
func (r *Repo[T, K]) Set(field string, val any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Set(field, val)
	return newR
}

// SetMap 根据Map设置字段
func (r *Repo[T, K]) SetMap(m map[string]any) IRepo[T, K] {
	newR := r.cloneInternal()
	for key, value := range m {
		newR.match.Set(key, value)
	}
	return newR
}

func (r *Repo[T, K]) Where(s string, a ...any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.db = newR.db.Where(s, a...)
	return newR
}

func (r *Repo[T, K]) Desc(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Desc(s)
	return newR
}

func (r *Repo[T, K]) Asc(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Asc(s)
	return newR
}

func (r *Repo[T, K]) Omit(s ...string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.omits = append(newR.omits, s...)
	return newR
}

func (r *Repo[T, K]) Eq(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Eq(s, a)
	return newR
}

func (r *Repo[T, K]) NEq(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.NEq(s, a)
	return newR
}

func (r *Repo[T, K]) In(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.In(s, a)
	return newR
}

func (r *Repo[T, K]) Gte(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Gte(s, a)
	return newR
}

func (r *Repo[T, K]) Gt(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Gt(s, a)
	return newR
}

func (r *Repo[T, K]) Lt(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Lt(s, a)
	return newR
}

func (r *Repo[T, K]) Lte(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Lte(s, a)
	return newR
}

func (r *Repo[T, K]) NotNull(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.NotNull(s)
	return newR
}

func (r *Repo[T, K]) Null(s string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Null(s)
	return newR
}

func (r *Repo[T, K]) Or(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Or(s, a)
	return newR
}

func (r *Repo[T, K]) Like(s string, a any) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.match.Like(s, a)
	return newR
}

func (r *Repo[T, K]) Select(s ...string) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.selects = append(newR.selects, s...)
	return newR
}

func (r *Repo[T, K]) Limit(i int64) IRepo[T, K] {
	newR := r.cloneInternal()
	newR.limit = i
	return newR
}
