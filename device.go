package ftdi

/*
#include <stdlib.h>
#include <ftdi.h>
#include <libusb.h>

#cgo pkg-config: libftdi1
*/
import "C"

import (
	"runtime"
	"unsafe"
)

func makeError(ctx *C.struct_ftdi_context, code C.int) error {
	if code >= 0 {
		return nil
	}
	return &Error{
		code: int(code),
		str:  C.GoString(C.ftdi_get_error_string(ctx)),
	}
}

type USBDev struct {
	d *C.struct_libusb_device
}

func FindAll(vendor, product int) ([]USBDev, error) {
	ctx := new(C.struct_ftdi_context)
	e := C.ftdi_init(ctx)
	defer C.ftdi_deinit(ctx)
	if e < 0 {
		return nil, makeError(ctx, e)
	}
	var dl *C.struct_ftdi_device_list
	e = C.ftdi_usb_find_all(ctx, &dl, C.int(vendor), C.int(product))
	if e < 0 {
		return nil, makeError(ctx, e)
	}
	defer C.ftdi_list_free2(dl)

	n := 0
	for e := dl; e != nil; e = e.next {
		n++
	}
	ret := make([]USBDev, n)
	i := 0
	for e := dl; e != nil; e = e.next {
		ret[i].d = e.dev
		C.libusb_ref_device(e.dev)
		i++
	}
	return ret, nil
}

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
	return makeError(d.ctx, code)
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

func OpenUSBDev(u USBDev, c Channel) (*Device, error) {
	d, err := makeDevice(c)
	if err != nil {
		return nil, err
	}
	e := C.ftdi_usb_open_dev(d.ctx, u.d)
	if e < 0 {
		defer d.deinit()
		return nil, d.makeError(e)
	}
	return d, nil
}

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

// SetBitmode sets i/o mode for device
func (d *Device) SetBitmode(iomask byte, mode Mode) error {
	e := C.ftdi_set_bitmode(d.ctx, C.uchar(iomask), C.uchar(mode))
	return d.makeError(e)
}

func (d *Device) Reset() error {
	return d.makeError(C.ftdi_usb_reset(d.ctx))
}

func (d *Device) PurgeRxBuffer() error {
	return d.makeError(C.ftdi_usb_purge_rx_buffer(d.ctx))
}

func (d *Device) PurgeTxBuffer() error {
	return d.makeError(C.ftdi_usb_purge_tx_buffer(d.ctx))
}

func (d *Device) PurgeBuffers() error {
	return d.makeError(C.ftdi_usb_purge_buffers(d.ctx))
}

func (d *Device) ReadChunkSize() (int, error) {
	var cs C.uint
	e := C.ftdi_read_data_get_chunksize(d.ctx, &cs)
	return int(cs), d.makeError(e)
}

// SetReadChunkSize configure read buffer chunk size (default is 4096).
func (d *Device) SetReadChunkSize(cs int) error {
	return d.makeError(C.ftdi_read_data_set_chunksize(d.ctx, C.uint(cs)))
}

func (d *Device) WriteChunkSize() (int, error) {
	var cs C.uint
	e := C.ftdi_write_data_get_chunksize(d.ctx, &cs)
	return int(cs), d.makeError(e)
}

// SetWriteChunkSize configure write buffer chunk size (default is 4096).
func (d *Device) SetWriteChunkSize(cs int) error {
	return d.makeError(C.ftdi_write_data_set_chunksize(d.ctx, C.uint(cs)))
}

// LatencyTimer returns latency timer value (ms)
func (d *Device) LatencyTimer() (int, error) {
	var lt C.uchar
	e := C.ftdi_get_latency_timer(d.ctx, &lt)
	return int(lt), d.makeError(e)
}

// SetLatencyTimer sets latency timer to lt (value beetwen 1 and 255)
func (d *Device) SetLatencyTimer(lt int) error {
	return d.makeError(C.ftdi_set_latency_timer(d.ctx, C.uchar(lt)))
}

func (d *Device) Read(buf []byte) (int, error) {
	n := C.ftdi_read_data(
		d.ctx,
		(*C.uchar)(unsafe.Pointer(&buf[0])),
		C.int(len(buf)),
	)
	if n < 0 {
		return 0, d.makeError(n)
	}
	return int(n), nil
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

// SetBaudrate sets the rate of data transfer
//
// For standard USB-UART adapter it sets UART boudrate.
//
// For bitbang mode the clock is actually 16 times the br. From the FTDI
// documentation for FT232R bitbang mode:
// "The clock for the Asynchronous Bit Bang mode is actually 16 times the
// Baud rate. A value of 9600 Baud would transfer the data at (9600x16) = 153600
// bytes per second, or 1 every 6.5 μS."
//
// FT232R suports baudrates from 183.1 baud to 3 Mbaud but for real applications
// it should be <= 1 Mbaud: Actual baudrate is set to discrete value that
// satisfies the equation br = 3000000 / (n + x) where n can be an integer
// between 2 and 16384 and x can be a sub-integer of the value 0, 0.125, 0.25,
// 0.375, 0.5, 0.625, 0.75, or 0.875. When n == 1 then x should be 0, i.e.
// baud rate divisors with values between 1 and 2 are not possible.
func (d *Device) SetBaudrate(br int) error {
	return d.makeError(C.ftdi_set_baudrate(d.ctx, C.int(br)))
}

// ChipID reads FTDI Chip-ID (not all devices support this)
func (d *Device) ChipID() (uint32, error) {
	var id C.uint
	e := C.ftdi_read_chipid(d.ctx, &id)
	if e < 0 {
		return 0, d.makeError(e)
	}
	return uint32(id), nil
}

// EEPROM returns a handler to the device internal EEPROM subsystem
func (d *Device) EEPROM() EEPROM {
	return EEPROM{d}
}
