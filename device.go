package ftdi

/*
#include <stdlib.h>
#include <ftdi.h>
#include <libusb.h>

#cgo windows CFLAGS: -Ilibftdi1-1.5/include/libftdi -Ilibftdi1-1.5/include/libusb-1.0
#cgo windows LDFLAGS: ${SRCDIR}/libftdi1-1.5/lib64/libftdi1.a ${SRCDIR}/libftdi1-1.5/lib64/libusb-1.0.a
#cgo linux pkg-config: libftdi1
#cgo darwin pkg-config: libftdi1

// libftdi 1.5 deprecated the purge API.  Use a wrapper to avoid the
// deprecation warnings while still supporting 1.4.
int libftdi_tciflush(struct ftdi_context *ftdi) {
#ifdef SIO_TCIFLUSH
	return ftdi_tciflush(ftdi);
#else
	return ftdi_usb_purge_rx_buffer(ftdi);
#endif
}

int libftdi_tcoflush(struct ftdi_context *ftdi) {
#ifdef SIO_TCIFLUSH
	return ftdi_tcoflush(ftdi);
#else
	return ftdi_usb_purge_tx_buffer(ftdi);
#endif
}

int libftdi_tcioflush(struct ftdi_context *ftdi) {
#ifdef SIO_TCIFLUSH
	return ftdi_tcioflush(ftdi);
#else
	return ftdi_usb_purge_buffers(ftdi);
#endif
}

*/
import "C"

import (
	"errors"
	"runtime"
	"unicode/utf16"
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
	Manufacturer, Description, Serial string

	d *C.struct_libusb_device
}

func (u *USBDev) unref() {
	if u.d == nil {
		panic("USBDev.unref on uninitialized device")
	}
	C.libusb_unref_device(u.d)
	u.d = nil // Help GC.
}

func (u *USBDev) Close() {
	runtime.SetFinalizer(u, nil)
	u.unref()
}

func getLangId(dh *C.libusb_device_handle) (C.uint16_t, error) {
	var buf [128]C.char
	e := C.libusb_get_string_descriptor(
		dh, 0, 0,
		(*C.uchar)(unsafe.Pointer(&buf[0])), C.int(len(buf)),
	)
	if e < 0 {
		return 0, USBError(e)
	}
	if e < 4 {
		return 0, errors.New("not enough data in USB language IDs descriptor")
	}
	return C.uint16_t(uint(buf[2]) | uint(buf[3])<<8), nil
}

func getStringDescriptor(dh *C.libusb_device_handle, id C.uint8_t, langid C.uint16_t) (string, error) {
	var buf [128]C.char
	e := C.libusb_get_string_descriptor(
		dh, id, C.uint16_t(langid),
		(*C.uchar)(unsafe.Pointer(&buf[0])), C.int(len(buf)),
	)
	if e < 0 {
		return "", USBError(e)
	}
	if e < 2 {
		return "", errors.New("not enough data for USB string descriptor")
	}
	l := C.int(buf[0])
	if l > e {
		return "", errors.New("USB string descriptor is too short")
	}
	b := buf[2:l]
	uni16 := make([]uint16, len(b)/2)
	for i := range uni16 {
		uni16[i] = uint16(b[i*2]) | uint16(b[i*2+1])<<8
	}
	return string(utf16.Decode(uni16)), nil
}

// getStrings updates Manufacturer, Description, Serial strings descriptors
// in unicode form. It doesn't use ftdi_usb_get_strings because libftdi
// converts  unicode strings to ASCII.
func (u *USBDev) getStrings(dev *C.libusb_device, ds *C.struct_libusb_device_descriptor) error {
	var (
		err error
		dh  *C.libusb_device_handle
	)
	if e := C.libusb_open(dev, &dh); e != 0 {
		return USBError(e)
	}
	defer C.libusb_close(dh)
	langid, err := getLangId(dh)
	if err != nil {
		return err
	}
	u.Manufacturer, err = getStringDescriptor(dh, ds.iManufacturer, langid)
	if err != nil {
		return err
	}
	u.Description, err = getStringDescriptor(dh, ds.iProduct, langid)
	if err != nil {
		return err
	}
	u.Serial, err = getStringDescriptor(dh, ds.iSerialNumber, langid)
	return err
}

/*func FindAll(vendor, product int) ([]*USBDev, error) {
	ctx := new(C.struct_ftdi_context)
	e := C.ftdi_init(ctx)
	//	defer C.ftdi_deinit(ctx)
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
	for el := dl; el != nil; el = el.next {
		n++
	}
	ret := make([]*USBDev, n)
	i := 0
	for el := dl; el != nil; el = el.next {
		u := new(USBDev)
		u.d = el.dev
		C.libusb_ref_device(el.dev)
		runtime.SetFinalizer(u, (*USBDev).unref)
		if err := u.getStrings(ctx); err != nil {
			return nil, err
		}
		ret[i] = u
		i++
	}
	return ret, nil
}*/

// FindAll search for all USB devices with specified vendor and  product id.
// It returns slice od found devices.
func FindAll(vendor, product int) ([]*USBDev, error) {
	if e := C.libusb_init(nil); e != 0 {
		return nil, USBError(e)
	}
	var dptr **C.struct_libusb_device
	if e := C.libusb_get_device_list(nil, &dptr); e < 0 {
		return nil, USBError(e)
	}
	defer C.libusb_free_device_list(dptr, 1)
	devs := (*[1 << 28]*C.libusb_device)(unsafe.Pointer(dptr))

	var n int
	for i, dev := range devs[:] {
		if dev == nil {
			n = i
			break
		}
	}
	descr := make([]*C.struct_libusb_device_descriptor, n)
	for i, dev := range devs[:n] {
		var ds C.struct_libusb_device_descriptor
		if e := C.libusb_get_device_descriptor(dev, &ds); e < 0 {
			return nil, USBError(e)
		}
		if int(ds.idVendor) == vendor && int(ds.idProduct) == product {
			descr[i] = &ds
			continue
		}
		if vendor == 0 && product == 0 && ds.idVendor == 0x403 {
			switch ds.idProduct {
			case 0x6001, 0x6010, 0x6011, 0x6014, 0x6015:
				descr[i] = &ds
				continue
			}
		}
		n--
	}
	if n == 0 {
		return nil, nil
	}
	ret := make([]*USBDev, n)
	n = 0
	for i, ds := range descr {
		if ds == nil {
			continue
		}
		u := new(USBDev)
		u.d = devs[i]
		C.libusb_ref_device(u.d)
		runtime.SetFinalizer(u, (*USBDev).unref)
		u.getStrings(u.d, ds)
		ret[n] = u
		n++
	}
	return ret, nil
}

// Device represents FTDI device.
type Device struct {
	ctx *C.struct_ftdi_context
}

// Type is numeric type id of FTDI device.
type Type uint32

const (
	TypeAM Type = iota
	TypeBM
	Type2232C
	TypeR
	Type2232H
	Type4232H
	Type232H
	Type230x
)

var types = []string{"AM", "BM", "2232C", "R", "2232H", "4232H", "232H", "230X"}

// String returns text name that describes type id.
func (t Type) String() string {
	if t >= Type(len(types)) {
		return "unknown"
	}
	return types[t]
}

// Type returns type of device d.
func (d *Device) Type() Type {
	return Type(d.ctx._type)
}

func (d *Device) free() {
	if d.ctx == nil {
		panic("Device.free on uninitialized device")
	}
	C.ftdi_free(d.ctx)
	d.ctx = nil
}

func (d *Device) makeError(code C.int) error {
	return makeError(d.ctx, code)
}

func (d *Device) close() error {
	defer d.free()
	if e := C.ftdi_usb_close(d.ctx); e != 0 {
		return d.makeError(e)
	}
	return nil
}

// Close closes device
func (d *Device) Close() error {
	runtime.SetFinalizer(d, nil)
	return d.close()
}

func makeDevice(c Channel) (*Device, error) {
	ctx, err := C.ftdi_new()
	if ctx == nil {
		return nil, err
	}
	d := &Device{ctx}
	if c != ChannelAny {
		if e := C.ftdi_set_interface(d.ctx, uint32(c)); e < 0 {
			defer d.free()
			return nil, d.makeError(e)
		}
	}
	return d, nil
}

// Channel represents channel (interface) of FTDI device. Some devices have more
// than one channel (eg. FT2232 has 2 channels, FT4232 has 4 channels).
type Channel uint32

const (
	ChannelAny Channel = iota
	ChannelA
	ChannelB
	ChannelC
	ChannelD
)

// OpenUSBDev opens channel (interface) c of USB device u.
// u must be FTDI device.
func OpenUSBDev(u *USBDev, c Channel) (*Device, error) {
	d, err := makeDevice(c)
	if err != nil {
		return nil, err
	}
	if e := C.ftdi_usb_open_dev(d.ctx, u.d); e < 0 {
		defer d.free()
		return nil, d.makeError(e)
	}
	runtime.SetFinalizer(d, (*Device).close)
	return d, nil
}

// OpenFirst opens the first device with a given vendor and product ids. Uses
// specified channel c.
func OpenFirst(vendor, product int, c Channel) (*Device, error) {
	d, err := makeDevice(c)
	if err != nil {
		return nil, err
	}
	if e := C.ftdi_usb_open(d.ctx, C.int(vendor), C.int(product)); e < 0 {
		defer d.free()
		return nil, d.makeError(e)
	}
	runtime.SetFinalizer(d, (*Device).close)
	return d, nil
}

// Open opens the index-th device with a given vendor id, product id,
// description and serial. Uses specified channel c.
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
		defer d.free()
		return nil, d.makeError(e)
	}
	runtime.SetFinalizer(d, (*Device).close)
	return d, nil
}

// OpenBusAddr opens the device at a given USB bus and device address. Uses
// specified channel c.
func OpenBusAddr(bus, address int, c Channel) (*Device, error) {
	d, err := makeDevice(c)
	if err != nil {
		return nil, err
	}
	e := C.ftdi_usb_open_bus_addr(
		d.ctx,
		C.uint8_t(bus),
		C.uint8_t(address),
	)
	if e < 0 {
		defer d.free()
		return nil, d.makeError(e)
	}
	runtime.SetFinalizer(d, (*Device).close)
	return d, nil
}

// OpenString opens the ftdi-device described by a description-string. Uses
// specified channel c.
func OpenString(description string, c Channel) (*Device, error) {
	d, err := makeDevice(c)
	if err != nil {
		return nil, err
	}
	e := C.ftdi_usb_open_string(
		d.ctx,
		C.CString(description),
	)
	if e < 0 {
		defer d.free()
		return nil, d.makeError(e)
	}
	runtime.SetFinalizer(d, (*Device).close)
	return d, nil
}

// Mode represents operation mode that FTDI device can work.
type Mode byte

const (
	ModeReset Mode = (1 << iota) >> 1
	ModeBitbang
	ModeMPSSE
	ModeSyncBB
	ModeMCU
	ModeOpto
	ModeCBUS
	ModeSyncFF
	ModeFT1284
)

// MPSSE commands
// See https://www.ftdichip.com/Documents/AppNotes/AN_108_Command_Processor_for_MPSSE_and_MCU_Host_Bus_Emulation_Modes.pdf
// for full details
const (
	MPSSEWriteNeg           byte = C.MPSSE_WRITE_NEG
	MPSSEBitMode            byte = C.MPSSE_BITMODE
	MPSSEReadNeg            byte = C.MPSSE_READ_NEG
	MPSSELSB                byte = C.MPSSE_LSB
	MPSSEDoWrite            byte = C.MPSSE_DO_WRITE
	MPSSEDoRead             byte = C.MPSSE_DO_READ
	MPSSEWriteTMS           byte = C.MPSSE_WRITE_TMS
	MPSSESetBitsLow         byte = C.SET_BITS_LOW
	MPSSESetBitsHigh        byte = C.SET_BITS_HIGH
	MPSSEGetBitsLow         byte = C.GET_BITS_LOW
	MPSSELoopbackStart      byte = C.LOOPBACK_START
	MPSSELoopbackEnd        byte = C.LOOPBACK_END
	MPSSETCKDivisor         byte = C.TCK_DIVISOR
	MPSSEDisableDiv5        byte = C.DIS_DIV_5
	MPSSEEnableDiv5         byte = C.EN_DIV_5
	MPSSEEnable3Phase       byte = C.EN_3_PHASE
	MPSSEDisable3Phase      byte = C.DIS_3_PHASE
	MPSSEClockBits          byte = C.CLK_BITS
	MPSSEClockBytes         byte = C.CLK_BYTES
	MPSSEClockWaitHigh      byte = C.CLK_WAIT_HIGH
	MPSSEClockWaitLow       byte = C.CLK_WAIT_LOW
	MPSSEEnableAdaptive     byte = C.EN_ADAPTIVE
	MPSSEDisableAdaptive    byte = C.DIS_ADAPTIVE
	MPSSEClockBytesOrHigh   byte = C.CLK_BYTES_OR_HIGH
	MPSSEClockBytesOrLow    byte = C.CLK_BYTES_OR_LOW
	MPSSEDriveOpenCollector byte = C.DRIVE_OPEN_COLLECTOR
	MPSSESendImmediate      byte = C.SEND_IMMEDIATE
	MPSSEWaitOnHigh         byte = C.WAIT_ON_HIGH
	MPSSEWaitOnLow          byte = C.WAIT_ON_LOW
	MPSSEReadShort          byte = C.READ_SHORT
	MPSSEReadExtended       byte = C.READ_EXTENDED
	MPSSEWriteShort         byte = C.WRITE_SHORT
	MPSSEWriteExtended      byte = C.WRITE_EXTENDED
)

// MPSSEDivValue calculates the two bytes that are required to be supplied after
// MPSSETCKDivisor to get the desired clock speed (in Hz).
// Set the dvi5 flag if MPSSEEnableDiv5 has been sent, to use a 12MHz base clock,
// instead of 60MHz.
func MPSSEDivValue(rate int, div5 bool) int {
	clk := 60_000_000
	if div5 {
		clk /= 5
	}
	if rate <= 0 || rate > clk {
		return 0
	}
	if (clk/rate)-1 > 0xffff {
		return 0xffff
	}
	return clk/rate - 1
}

// SetBitmode sets operation mode for device d to mode. iomask bitmask
// configures lines corresponding to its bits as input (bit=0) or output (bit=1).
func (d *Device) SetBitmode(iomask byte, mode Mode) error {
	e := C.ftdi_set_bitmode(d.ctx, C.uchar(iomask), C.uchar(mode))
	return d.makeError(e)
}

// Reset resets device d.
func (d *Device) Reset() error {
	return d.makeError(C.ftdi_usb_reset(d.ctx))
}

// PurgeWriteBuffer clears Rx buffer (buffer for data received from USB?).
func (d *Device) PurgeWriteBuffer() error {
	return d.makeError(C.libftdi_tciflush(d.ctx))
}

// PurgeReadBuffer clears Tx buffer (buffer for data that will be sent to USB?).
func (d *Device) PurgeReadBuffer() error {
	return d.makeError(C.libftdi_tcoflush(d.ctx))
}

// PurgeBuffers clears both (Tx and Rx) buffers.
func (d *Device) PurgeBuffers() error {
	return d.makeError(C.libftdi_tcioflush(d.ctx))
}

// ReadChunkSize returns current value of read buffer chunk size.
func (d *Device) ReadChunkSize() (int, error) {
	var cs C.uint
	e := C.ftdi_read_data_get_chunksize(d.ctx, &cs)
	return int(cs), d.makeError(e)
}

// SetReadChunkSize configure read chunk size for device (default is 4096B) and
// size of software buffer dedicated for reading data from device...
func (d *Device) SetReadChunkSize(cs int) error {
	return d.makeError(C.ftdi_read_data_set_chunksize(d.ctx, C.uint(cs)))
}

// WriteChunkSize returns current value of write chunk size.
func (d *Device) WriteChunkSize() (int, error) {
	var cs C.uint
	e := C.ftdi_write_data_get_chunksize(d.ctx, &cs)
	return int(cs), d.makeError(e)
}

// SetWriteChunkSize configure write chunk size (default is 4096). If more than
// cs bytes need to be send to device, they will be split to chunks of size cs.
func (d *Device) SetWriteChunkSize(cs int) error {
	return d.makeError(C.ftdi_write_data_set_chunksize(d.ctx, C.uint(cs)))
}

// LatencyTimer returns latency timer value [ms].
func (d *Device) LatencyTimer() (int, error) {
	var lt C.uchar
	e := C.ftdi_get_latency_timer(d.ctx, &lt)
	return int(lt), d.makeError(e)
}

// SetLatencyTimer sets latency timer to lt (value beetwen 1 and 255). If FTDI
// device has fewer data to completely fill one USB packet (<62 B) it waits for
// lt ms before sending data to USB.
func (d *Device) SetLatencyTimer(lt int) error {
	return d.makeError(C.ftdi_set_latency_timer(d.ctx, C.uchar(lt)))
}

// Read reads data from device to data. It returns number of bytes read.
func (d *Device) Read(data []byte) (int, error) {
	n := C.ftdi_read_data(
		d.ctx,
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(len(data)),
	)
	if n < 0 {
		return 0, d.makeError(n)
	}
	return int(n), nil
}

// Write writes data from buf to device. It retruns number of bytes written.
func (d *Device) Write(data []byte) (int, error) {
	n := C.ftdi_write_data(
		d.ctx,
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(len(data)),
	)
	if n < 0 {
		return 0, d.makeError(n)
	}
	return int(n), nil
}

// WriteString writes bytes from string s to device. It retruns number of bytes written.
func (d *Device) WriteString(s string) (int, error) {
	// BUG: This will cause problems when string implementation changes.
	type stringHeader struct {
		data unsafe.Pointer
		len  int
	}
	n := C.ftdi_write_data(
		d.ctx,
		(*C.uchar)((*stringHeader)(unsafe.Pointer(&s)).data),
		C.int(len(s)),
	)
	if n < 0 {
		return 0, d.makeError(n)
	}
	return int(n), nil
}

// ReadByte reads one byte from device.
func (d *Device) ReadByte() (byte, error) {
	var b byte
	if n := C.ftdi_read_data(d.ctx, (*C.uchar)(&b), 1); n != 1 {
		return 0, d.makeError(n)
	}
	return b, nil
}

// WriteByte writes one byte to device.
func (d *Device) WriteByte(b byte) error {
	if n := C.ftdi_write_data(d.ctx, (*C.uchar)(&b), 1); n != 1 {
		return d.makeError(n)
	}
	return nil
}

// Pins returns current state of pins (circumventing the read buffer).
func (d *Device) Pins() (b byte, err error) {
	if e := C.ftdi_read_pins(d.ctx, (*C.uchar)(&b)); e != 0 {
		err = d.makeError(e)
	}
	return
}

// SetBaudrate sets the rate of data transfer.
//
// For standard USB-UART adapter it sets UART baudrate.
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

// Parity represents the parity mode
type Parity int

const (
	// ParityNone indicates no parity bit is used
	ParityNone Parity = C.NONE
	// ParityOdd indicates an odd parity bit is used
	ParityOdd Parity = C.ODD
	// ParityEven indicates an even parity bit is used
	ParityEven Parity = C.EVEN
	// ParityMark indicates that the parity bit should be a 1
	ParityMark Parity = C.MARK
	// ParitySpace indicates that the parity bit should be a 0
	ParitySpace Parity = C.SPACE
)

// DataBits represents the number of data bits to expect
type DataBits int

const (
	// DataBits7 indicates that 7 data bits are used
	DataBits7 DataBits = C.BITS_7
	// DataBits8 indicates that 8 data bits are used
	DataBits8 DataBits = C.BITS_8
)

// StopBits represents the number of stop bits to expect
type StopBits int

const (
	// StopBits1 indicates only one stop bit is expected
	StopBits1 StopBits = C.STOP_BIT_1
	// StopBits15 indicates one and a half stop bits are expected
	StopBits15 StopBits = C.STOP_BIT_15
	// StopBits2 indicates two stop bits are expected
	StopBits2 StopBits = C.STOP_BIT_2
)

// Break represents the break mode
type Break int

const (
	// BreakOff disables the use of a break signal
	BreakOff Break = C.BREAK_OFF
	// BreakOn enables the use of a break signal
	BreakOn Break = C.BREAK_ON
)

// SetLineProperties sets the uart data bit count, stop bits count, and parity mode
func (d *Device) SetLineProperties(bits DataBits, stopbits StopBits, parity Parity) error {
	e := C.ftdi_set_line_property(
		d.ctx,
		uint32(bits),
		uint32(stopbits),
		uint32(parity),
	)
	return d.makeError(e)
}

// SetLineProperties2 sets the uart data bit count, stop bits count, parity mode,
// and break usage
func (d *Device) SetLineProperties2(bits DataBits, stopbits StopBits, parity Parity, breaks Break) error {
	e := C.ftdi_set_line_property2(
		d.ctx,
		uint32(bits),
		uint32(stopbits),
		uint32(parity),
		uint32(breaks),
	)
	return d.makeError(e)
}

// FlowCtrl represents the flow control mode.
type FlowCtrl int

const (
	// FlowCtrlDisable disables all automatic flow control.
	FlowCtrlDisable FlowCtrl = (1 << iota) >> 1
	// FlowCtrlRTSCTS enables RTS CTS flow control.
	FlowCtrlRTSCTS
	// FlowCtrlDTRDSR enables DTR DSR flow control.
	FlowCtrlDTRDSR
	// FlowCtrlXONXOFF enables XON XOF flow control.
	FlowCtrlXONXOFF
)

// SetFlowControl sets the flow control parameter
func (d *Device) SetFlowControl(flowctrl FlowCtrl) error {
	return d.makeError(C.ftdi_setflowctrl(d.ctx, C.int(flowctrl)))
}

// SetDTRRTS manually sets the DTR and RTS output lines from the
// least significant bit of dtr and rts.
func (d *Device) SetDTRRTS(dtr, rts int) error {
	return d.makeError(C.ftdi_setdtr_rts(d.ctx, C.int(dtr&1), C.int(rts&1)))
}

// SetDTR manually sets the DTR output line from the least significant
// bit of dtr.
func (d *Device) SetDTR(dtr int) error {
	return d.makeError(C.ftdi_setdtr(d.ctx, C.int(dtr&1)))
}

// SetRTS manually sets the RTS output line from the least significant
// bit of rts.
func (d *Device) SetRTS(rts int) error {
	return d.makeError(C.ftdi_setrts(d.ctx, C.int(rts&1)))
}

// ChipID reads FTDI Chip-ID (not all devices support this).
func (d *Device) ChipID() (uint32, error) {
	var id C.uint
	if e := C.ftdi_read_chipid(d.ctx, &id); e < 0 {
		return 0, d.makeError(e)
	}
	return uint32(id), nil
}

// EEPROM returns a handler to the device internal EEPROM subsystem.
func (d *Device) EEPROM() EEPROM {
	return EEPROM{d}
}

type Transfer struct {
	c C.struct_ftdi_transfer_control
}

var errSubmitTransfer = errors.New("libusb_submit_transfer")

func (d *Device) SubmitRead(data []byte) (*Transfer, error) {
	tc, err := C.ftdi_read_data_submit(
		d.ctx,
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(len(data)),
	)
	if tc == nil {
		if err == nil {
			err = errSubmitTransfer
		}
		return nil, err
	}
	return (*Transfer)(unsafe.Pointer(tc)), nil
}

func (d *Device) SubmitWrite(data []byte) (*Transfer, error) {
	tc, err := C.ftdi_write_data_submit(
		d.ctx,
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(len(data)),
	)
	if tc == nil {
		if err == nil {
			err = errSubmitTransfer
		}
		return nil, err
	}
	return (*Transfer)(unsafe.Pointer(tc)), nil
}

func (t *Transfer) Done() (int, error) {
	n := C.ftdi_transfer_data_done(&t.c)
	if n < 0 {
		return 0, makeError(t.c.ftdi, n)
	}
	return int(n), nil
}
