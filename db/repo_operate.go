package db

func (r Repo[T]) Raw(s string, a ...any) *IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r Repo[T]) Exec(s string, a ...any) *IRepo[T] {
	//TODO implement me
	panic("implement me")
}

func (r Repo[T]) Add(t *T) int64 {
	//TODO implement me
	panic("implement me")
}

func (r Repo[T]) AddBatch(ts []*T) int64 {
	//TODO implement me
	panic("implement me")
}

func (r Repo[T]) Save(t *T) int64 {
	//TODO implement me
	panic("implement me")
}

func (r Repo[T]) Update(t *T) int64 {
	//TODO implement me
	panic("implement me")
}

func (r Repo[T]) Del() int64 {
	//TODO implement me
	panic("implement me")
}
