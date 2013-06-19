package ftdi

/*
#include <stdlib.h>
#include <ftdi.h>

#cgo pkg-config: libftdi1
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// Device represents some FTDI device
type Device struct {
	ctx *C.struct_ftdi_context
}

type Type uint32

const (
	TypeAM Type = iota
	TypeBM
	Type2232C
	TypeR
	Type2232H
	Type4232H
	Type232H
)

var types = []string{"AM", "BM", "2232C", "R", "2232H", "4232H", "232H"}

func (t Type) String() string {
	if t >= Type(len(types)) {
		return "unknown"
	}
	return types[t]
}

func (d *Device) Type() Type {
	return Type(d.ctx._type)
}

func makeDevice(c Channel) (*Device, error) {
	d := new(Device)
	d.ctx = new(C.struct_ftdi_context)
	e := C.ftdi_init(d.ctx)
	if e < 0 {
		defer d.deinit()
		return nil, d.makeError(e)
	}
	if c != ChannelAny {
		e = C.ftdi_set_interface(d.ctx, C.enum_ftdi_interface(c))
		if e < 0 {
			defer d.deinit()
			return nil, d.makeError(e)
		}
	}
	runtime.SetFinalizer(d, (*Device).Close)
	return d, nil
}

func (d *Device) deinit() {
	C.ftdi_deinit(d.ctx)
}

func (d *Device) makeError(code C.int) error {
	if code == 0 {
		return nil
	}
	return &Error{
		code: int(code),
		str:  C.GoString(C.ftdi_get_error_string(d.ctx)),
	}
}

// Close closes device
func (d *Device) Close() error {
	defer d.deinit()
	e := C.ftdi_usb_close(d.ctx)
	runtime.SetFinalizer(d, nil)
	return d.makeError(e)
}

type Channel uint32

const (
	ChannelAny Channel = iota
	ChannelA
	ChannelB
	ChannelC
	ChannelD
)

// OpenFirst opens the first device with a given vendor and product ids. Uses
// specified interface.
func OpenFirst(vendor, product int, c Channel) (*Device, error) {
	d, err := makeDevice(c)
	if err != nil {
		return nil, err
	}
	e := C.ftdi_usb_open(d.ctx, C.int(vendor), C.int(product))
	if e < 0 {
		defer d.deinit()
		return nil, d.makeError(e)
	}
	return d, nil
}

// Open opens the index-th device with a given vendor id, product id,
// description and serial. Uses specified interface.
func Open(vendor, product int, description, serial string, index uint,
	c Channel) (*Device, error) {

	d, err := makeDevice(c)
	if err != nil {
		return nil, err
	}

	descr := C.CString(description)
	defer C.free(unsafe.Pointer(descr))
	ser := C.CString(serial)
	defer C.free(unsafe.Pointer(ser))

	e := C.ftdi_usb_open_desc_index(
		d.ctx,
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

type Mode byte

const (
	ModeReset Mode = iota
	ModeBitbang
	ModeMPSSE
	ModeSyncBB
	ModeMCU
	ModeOpto
	ModeCBUS
	ModeSyncFF
	ModeFT1284
)

func (d *Device) SetBitmode(iomask byte, mode Mode) error {
	e := C.ftdi_set_bitmode(d.ctx, C.uchar(iomask), C.uchar(mode))
	return d.makeError(e)
}

func (d *Device) Write(buf []byte) (int, error) {
	n := C.ftdi_write_data(
		d.ctx,
		(*C.uchar)(unsafe.Pointer(&buf[0])),
		C.int(len(buf)),
	)
	if n < 0 {
		return 0, d.makeError(n)
	}
	return int(n), nil
}

func (d *Device) WriteByte(b byte) error {
	n := C.ftdi_write_data(d.ctx, (*C.uchar)(&b), 1)
	if n != 1 {
		return d.makeError(n)
	}
	return nil
}

func (d *Device) SetBaudrate(br int) error {
	return d.makeError(C.ftdi_set_baudrate(d.ctx, C.int(br)))
}

// EEPROM returns a handler to the device internal EEPROM subsystem
func (d *Device) EEPROM() EEPROM {
	return EEPROM{d}
}
