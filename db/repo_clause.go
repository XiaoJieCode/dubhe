package db

func (r *Repo[T]) Where(s string, a ...any) IRepo[T] {
	//TODO implement me
	panic("implement me")
}
func (r *Repo[T]) Desc(s string) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Asc(s string) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Omit(s ...string) IRepo[T] {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) NotNull(s string) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Null(s string) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Or(s string, a ...any) IRepo[T] {
	//TODO implement me
	panic("implement me")
}
func (r *Repo[T]) Like(s string, a any) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Select(s ...string) IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r *Repo[T]) Limit(i int64) IRepo[T] {
	//TODO implement me
	panic("implement me")
}
