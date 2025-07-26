package err

type Error struct {
	msg  string
	code int
}

func (e Error) Error() string {
	return e.msg
}
