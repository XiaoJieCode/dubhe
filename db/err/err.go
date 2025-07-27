package err

type Error struct {
	msg   string
	code  int
	error error
}

func (e Error) Error() string {
	return e.msg
}

func NewError(err error) *Error {
	if err != nil {
		return &Error{error: err}
	}
	return nil
}
