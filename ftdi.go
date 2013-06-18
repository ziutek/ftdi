// Go binding for libFTDI library
// http://http://www.intra2net.com/en/developer/libftdi/
package ftdi

/*
#include <ftdi.h>

#cgo pkg-config: libftdi
*/
import "C"

import (
	"unsafe"
)

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

// Device represents some FTDI device
type Device struct {
	ctx C.struct_ftdi_context
}

func makeDevice(i Interface) (*Device, error) {
	d := new(Device)
	e := C.ftdi_init(&d.ctx)
	if e < 0 {
		defer d.deinit()
		return nil, d.makeError(e)
	}
	if i != InterfaceAny {
		e = C.ftdi_set_interface(&d.ctx, C.enum_ftdi_interface(i))
		defer d.deinit()
		return nil, d.makeError(e)
	}
	return d, nil
}

func (d *Device) deinit() {
	C.ftdi_deinit(&d.ctx)
}

func (d *Device) makeError(code C.int) error {
	if code == 0 {
		return nil
	}
	return &Error{
		code: int(code),
		str:  C.GoString(C.ftdi_get_error_string(&d.ctx)),
	}
}

// Close closes device
func (d *Device) Close() error {
	defer d.deinit()
	e := C.ftdi_usb_close(&d.ctx)
	return d.makeError(e)
}

type Interface int

const (
	InterfaceAny Interface = iota
	InterfaceA
	InterfaceB
	InterfaceC
	InterfaceD
)

// OpenFirst opens the first device with a given vendor and product ids. Uses
// specified interface.
func OpenFirst(vendor, product int, i Interface) (*Device, error) {
	d, err := makeDevice(i)
	if err != nil {
		return nil, err
	}
	e := C.ftdi_usb_open(&d.ctx, C.int(vendor), C.int(product))
	if e < 0 {
		defer d.deinit()
		return nil, d.makeError(e)
	}
	return d, nil
}

// Open opens the index-th device with a given vendor id, product id,
// description and serial. Uses specified interface.
func Open(vendor, product int, description, serial string, index uint,
	i Interface) (*Device, error) {

	d, err := makeDevice(i)
	if err != nil {
		return nil, err
	}

	descr := C.CString(description)
	defer C.free(unsafe.Pointer(descr))
	ser := C.CString(serial)
	defer C.free(unsafe.Pointer(ser))

	e := C.ftdi_usb_open_desc_index(
		&d.ctx,
		C.int(vendor), C.int(product),
		descr, ser,
		C.uint(index),
	)
	if e < 0 {
		defer d.deinit()
		return nil, d.makeError(e)
	}
	return d, nil
}

type Bitmode byte

const (
	BitmodeReset Bitmode = iota
	BitmodeBitbang
	BitmodeMPSSE
	BitmodeSyncBB
	BitmodeMCU
	BitmodeOpto
	BitmodeCBUS
	BitmodeSyncFF
	BitmodeFT1284
)

func (d *Device) SetBitmode(iomask byte, mode Bitmode) error {
	e := C.ftdi_set_bitmode(&d.ctx, C.uchar(iomask), C.uchar(mode))
	return d.makeError(e)
}

func (d *Device) Write(buf []byte) (int, error) {
	n := C.ftdi_write_data(
		&d.ctx,
		(*C.uchar)(unsafe.Pointer(&buf[0])),
		C.int(len(buf)),
	)
	if n < 0 {
		return 0, d.makeError(n)
	}
	return int(n), nil
}


func (d *Devic) ReadEEPROM(buf []byte)  {

}
}

