package ftdi

/*
#include <libusb.h>
*/
import "C"

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

type USBError int

func (e USBError) Error() string {
	return C.GoString(C.libusb_strerror(C.int(e)))
}
