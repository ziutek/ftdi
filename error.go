package ftdi

/*
#include <libusb.h>

// Some versions of libusb.h define libusb_strerror as taking a "int", while
// others take a "enum libusb_error". This doesn't matter for C but it matters
// for Go. Create a wrapper with a fixed type, so that Go doesn't complain.
const char* libusb_strerror_wrapper(enum libusb_error e) {
	return libusb_strerror(e);
}


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
	return C.GoString(C.libusb_strerror_wrapper(C.enum_libusb_error(e)))
}
