package ftdi

type Error struct {
	code int
	str  string
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Error() string {
	return e.str
}
