package ftdi

/*
#include <ftdi.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func cbool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

type EEPROM struct {
	d *Device
}

func (e EEPROM) makeError(code C.int) error {
	if code >= 0 {
		return nil
	}
	return &Error{
		code: int(code),
		str:  C.GoString(C.ftdi_get_error_string(e.d.ctx)),
	}
}

func (e EEPROM) Read() error {
	return e.makeError(C.ftdi_read_eeprom(e.d.ctx))
}

func (e EEPROM) Write() error {
	return e.makeError(C.ftdi_write_eeprom(e.d.ctx))
}

func (e EEPROM) Decode() error {
	return e.makeError(C.ftdi_eeprom_decode(e.d.ctx, 0))
}

func (e EEPROM) Build() error {
	return e.makeError(C.ftdi_eeprom_build(e.d.ctx))
}

func (e EEPROM) VendorId() uint16 {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.VENDOR_ID, &v)
	return uint16(v)
}

func (e EEPROM) SetVendorId(v uint16) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.VENDOR_ID, C.int(v))
}

func (e EEPROM) ProductId() uint16 {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.PRODUCT_ID, &v)
	return uint16(v)
}

func (e EEPROM) SetProductId(v uint16) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.PRODUCT_ID, C.int(v))
}

func (e EEPROM) ReleaseNumber() uint16 {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.RELEASE_NUMBER, &v)
	return uint16(v)
}

func (e EEPROM) SetReleaseNumber(v uint16) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.RELEASE_NUMBER, C.int(v))
}

func (e EEPROM) SelfPowered() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.SELF_POWERED, &v)
	return v != 0
}

func (e EEPROM) SetSelfPowered(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.SELF_POWERED, cbool(v))
}

func (e EEPROM) RemoteWakeup() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.REMOTE_WAKEUP, &v)
	return v != 0
}

func (e EEPROM) SetRemoteWakeup(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.REMOTE_WAKEUP, cbool(v))
}

func (e EEPROM) IsNotPNP() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.IS_NOT_PNP, &v)
	return v != 0
}

func (e EEPROM) SetIsNotPNP(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.IS_NOT_PNP, cbool(v))
}

func (e EEPROM) SuspendDBus7() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.SUSPEND_DBUS7, &v)
	return v != 0
}

func (e EEPROM) SetSuspendDBus7(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.SUSPEND_DBUS7, cbool(v))
}

func (e EEPROM) IsochronousInp() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.IN_IS_ISOCHRONOUS, &v)
	return v != 0
}

func (e EEPROM) SetIsochronousInp(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.IN_IS_ISOCHRONOUS, cbool(v))
}

func (e EEPROM) IsochronousOut() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.OUT_IS_ISOCHRONOUS, &v)
	return v != 0
}

func (e EEPROM) SetIsochronousOut(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.OUT_IS_ISOCHRONOUS, cbool(v))
}

func (e EEPROM) SuspendPullDowns() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.SUSPEND_PULL_DOWNS, &v)
	return v != 0
}

func (e EEPROM) SetSuspendPullDowns(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.SUSPEND_PULL_DOWNS, cbool(v))
}

func (e EEPROM) UseSerial() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.USE_SERIAL, &v)
	return v != 0
}

func (e EEPROM) SetUseSerial(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.USE_SERIAL, cbool(v))
}

func (e EEPROM) USBVersion() uint16 {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.USB_VERSION, &v)
	return uint16(v)
}

func (e EEPROM) SetUSBVersion(v uint16) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.USB_VERSION, C.int(v))
}

func (e EEPROM) UseUSBVersion() bool {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.USE_USB_VERSION, &v)
	return v != 0
}

func (e EEPROM) SetUseUSBVersion(v bool) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.USE_USB_VERSION, cbool(v))
}

// MaxPower returns maximum power consumption (max. current) from USB in mA
func (e EEPROM) MaxPower() int {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.MAX_POWER, &v)
	return int(v)
}

// SetMaxPower sets maximum power consumption (max. current) from USB in mA
func (e EEPROM) SetMaxPower(v int) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.MAX_POWER, C.int(v))
}

func (e EEPROM) channelValue(names []C.enum_ftdi_eeprom_value, c Channel) (
	v C.int) {
	n := int(c - ChannelA)
	if n < 0 || n >= len(names) {
		panic("bad channel")
	}
	C.ftdi_get_eeprom_value(e.d.ctx, uint32(names[n]), &v)
	return
}

func (e EEPROM) setChannelValue(names []C.enum_ftdi_eeprom_value, c Channel,
	v C.int) {

	n := int(c - ChannelA)
	if n < 0 || n >= len(names) {
		panic("bad channel")
	}
	C.ftdi_set_eeprom_value(e.d.ctx, uint32(names[n]), v)
}

type ChannelType byte

const (
	ChannelUART ChannelType = iota
	ChannelFIFO
	ChannelOPTO
	ChannelCPU
	ChannelFT1284
)

var channelTypes = []string{"UART", "FIFO", "OPTO", "CPU", "FT1284"}

func (ct ChannelType) String() string {
	if ct > ChannelFT1284 {
		return "unknown"
	}
	return channelTypes[ct]
}

var channelType = []C.enum_ftdi_eeprom_value{
	C.CHANNEL_A_TYPE,
	C.CHANNEL_B_TYPE,
}

// ChannelType returns type of c channel. c can be: ChannelA, ChannelB
func (e EEPROM) ChannelType(c Channel) ChannelType {
	return ChannelType(e.channelValue(channelType, c))
}

func (e EEPROM) SetChannelType(c Channel, v ChannelType) {
	e.setChannelValue(channelType, c, C.int(v))
}

var channelDriver = []C.enum_ftdi_eeprom_value{
	C.CHANNEL_A_DRIVER,
	C.CHANNEL_B_DRIVER,
	C.CHANNEL_C_DRIVER,
	C.CHANNEL_B_DRIVER,
}

// ChannelDriver returns true if c channel has a driver.
// c can be from range: ChannelA - ChannelD
func (e EEPROM) ChannelDriver(c Channel) bool {
	return e.channelValue(channelDriver, c) != 0
}

func (e EEPROM) SetChannelDriver(c Channel, v bool) {
	e.setChannelValue(channelDriver, c, cbool(v))
}

var channelRS485 = []C.enum_ftdi_eeprom_value{
	C.CHANNEL_A_RS485,
	C.CHANNEL_B_RS485,
	C.CHANNEL_C_RS485,
	C.CHANNEL_D_RS485,
}

// ChannelDriver returns true if c is RS485 channel.
// c can be from range: ChannelA - ChannelD
func (e EEPROM) ChannelRS485(c Channel) bool {
	return e.channelValue(channelRS485, c) != 0
}

func (e EEPROM) SetChannelRS485(c Channel, v bool) {
	e.setChannelValue(channelRS485, c, cbool(v))
}

var highCurrent = []C.enum_ftdi_eeprom_value{
	C.HIGH_CURRENT,
	C.HIGH_CURRENT_A,
	C.HIGH_CURRENT_B,
}

// ChannelDriver returns true if c channel is in high current mode .
// c can be from range: ChannAny - ChannelD (use ChannAny for TypeR device).
func (e EEPROM) HighCurrent(c Channel) bool {
	return e.channelValue(highCurrent, c+ChannelA) != 0
}

func (e EEPROM) SetHighCurrent(c Channel, v bool) {
	var hc C.int
	if v {
		if e.d.Type() == TypeR {
			hc = C.HIGH_CURRENT_DRIVE_R
		} else {
			hc = C.HIGH_CURRENT_DRIVE
		}
	}
	e.setChannelValue(highCurrent, c+ChannelA, hc)
}

type CBusFunction byte

const (
	CBusTxEn CBusFunction = iota
	CBusPwrEn
	CBusRxLED
	CBusTxLED
	CBusTxRxLED
	CBusSleep
	CBusClk48
	CBusClk24
	CBusClk12
	CBusClk6
	CBusIOMode
	CBusBBWR
	CBusBBRD
)

var cbusFunctions = []string{
	"TxD enabled",
	"power enabled",
	"Rx LED",
	"Tx LED",
	"Tx/Rx LED",
	"sleep",
	"clock 48 MHz",
	"clock 24 MHz",
	"clock 12 MHz",
	"clock 6 MHz",
	"bitbang I/O",
	"bitbang write",
	"bitbang read",
}

func (c CBusFunction) String() string {
	if c > CBusBBRD {
		return "unknown"
	}
	return cbusFunctions[c]
}

var cbusFunction = []C.enum_ftdi_eeprom_value{
	C.CBUS_FUNCTION_0,
	C.CBUS_FUNCTION_1,
	C.CBUS_FUNCTION_2,
	C.CBUS_FUNCTION_3,
	C.CBUS_FUNCTION_4,
	C.CBUS_FUNCTION_5,
	C.CBUS_FUNCTION_6,
	C.CBUS_FUNCTION_7,
	C.CBUS_FUNCTION_8,
	C.CBUS_FUNCTION_9,
}

func (e EEPROM) CBusFunction(n int) CBusFunction {
	if n < 0 || n >= len(cbusFunction) {
		panic("bad CBUS number")
	}
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, uint32(cbusFunction[n]), &v)
	return CBusFunction(v)
}

func (e EEPROM) SetCBusFunction(n int, v CBusFunction) {
	if n < 0 || n >= len(cbusFunction) {
		panic("bad CBUS number")
	}
	C.ftdi_set_eeprom_value(e.d.ctx, uint32(cbusFunction[n]), C.int(v))
}

func (e EEPROM) Invert() int {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.INVERT, &v)
	return int(v)
}

func (e EEPROM) SetInvert(v int) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.INVERT, C.int(v))
}

/*
	TODO:
   case GROUP0_DRIVE:
       *value = ftdi->eeprom->group0_drive;
       break;
   case GROUP0_SCHMITT:
       *value = ftdi->eeprom->group0_schmitt;
       break;
   case GROUP0_SLEW:
       *value = ftdi->eeprom->group0_slew;
       break;
   case GROUP1_DRIVE:
       *value = ftdi->eeprom->group1_drive;
       break;
   case GROUP1_SCHMITT:
       *value = ftdi->eeprom->group1_schmitt;
       break;
   case GROUP1_SLEW:
       *value = ftdi->eeprom->group1_slew;
       break;
   case GROUP2_DRIVE:
       *value = ftdi->eeprom->group2_drive;
       break;
   case GROUP2_SCHMITT:
       *value = ftdi->eeprom->group2_schmitt;
       break;
   case GROUP2_SLEW:
       *value = ftdi->eeprom->group2_slew;
       break;
   case GROUP3_DRIVE:
       *value = ftdi->eeprom->group3_drive;
       break;
   case GROUP3_SCHMITT:
       *value = ftdi->eeprom->group3_schmitt;
       break;
   case GROUP3_SLEW:
       *value = ftdi->eeprom->group3_slew;
       break;
   case POWER_SAVE:
       *value = ftdi->eeprom->powersave;
       break;
   case CLOCK_POLARITY:
       *value = ftdi->eeprom->clock_polarity;
       break;
   case DATA_ORDER:
       *value = ftdi->eeprom->data_order;
       break;
   case FLOW_CONTROL:
       *value = ftdi->eeprom->flow_control;
       break;
*/

func (e EEPROM) ChipType() byte {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.CHIP_TYPE, &v)
	return byte(v)
}

func (e EEPROM) SetChipType(v byte) {
	C.ftdi_set_eeprom_value(e.d.ctx, C.CHIP_TYPE, C.int(v))
}

func (e EEPROM) ChipSize() int {
	var v C.int
	C.ftdi_get_eeprom_value(e.d.ctx, C.CHIP_SIZE, &v)
	return int(v)
}

func (e EEPROM) ManufacturerString() string {
	var buf [128]C.char
	C.ftdi_eeprom_get_strings(e.d.ctx, (*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)), nil, 0, nil, 0)
	return C.GoString((*C.char)(unsafe.Pointer(&buf[0])))
}

func (e EEPROM) SetManufacturerString(s string) {
	C.ftdi_eeprom_set_strings(e.d.ctx, C.CString(s), nil, nil)
}

func (e EEPROM) ProductString() string {
	var buf [128]C.char
	C.ftdi_eeprom_get_strings(e.d.ctx, nil, 0, (*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)), nil, 0)
	return C.GoString((*C.char)(unsafe.Pointer(&buf[0])))
}

func (e EEPROM) SetProductString(s string) {
	C.ftdi_eeprom_set_strings(e.d.ctx, nil, C.CString(s), nil)
}

func (e EEPROM) SerialString() string {
	var buf [128]C.char
	C.ftdi_eeprom_get_strings(e.d.ctx, nil, 0, nil, 0, (*C.char)(unsafe.Pointer(&buf[0])), C.int(len(buf)))
	return C.GoString((*C.char)(unsafe.Pointer(&buf[0])))
}

func (e EEPROM) SetSerialString(s string) {
	C.ftdi_eeprom_set_strings(e.d.ctx, nil, nil, C.CString(s))
}

func (e EEPROM) String() string {
	return fmt.Sprintf(""+
		"vendor id:          %04xh\n"+
		"product id:         %04xh\n"+
		"release number:     %04xh\n"+
		"self powered:       %t\n"+
		"remote wakeup:      %t\n"+
		"is not PNP:         %t\n"+
		"isochronous inp:    %t\n"+
		"isochronous out:    %t\n"+
		"suspend pull downs: %t\n"+
		"use serial:         %t\n"+
		"USB version:        %04xh\n"+
		"use USB version:    %t\n"+
		"USB max. current:   %d mA\n"+
		"channel A type:     %s\n"+
		"channel B type:     %s\n"+
		"channel A driver:   %t\n"+
		"channel B driver:   %t\n"+
		"channel C driver:   %t\n"+
		"channel D driver:   %t\n"+
		"channel A RS485:    %t\n"+
		"channel B RS485:    %t\n"+
		"channel C RS485:    %t\n"+
		"channel D RS485:    %t\n"+
		"high current:       %t\n"+
		"high current A:     %t\n"+
		"high current B:     %t\n"+
		"CBUS[0]:            %s\n"+
		"CBUS[1]:            %s\n"+
		"CBUS[2]:            %s\n"+
		"CBUS[3]:            %s\n"+
		"CBUS[4]:            %s\n"+
		"CBUS[5]:            %s\n"+
		"CBUS[6]:            %s\n"+
		"CBUS[7]:            %s\n"+
		"CBUS[8]:            %s\n"+
		"CBUS[9]:            %s\n"+
		"invert:             %010bb\n"+
		"",
		e.VendorId(),
		e.ProductId(),
		e.ReleaseNumber(),
		e.SelfPowered(),
		e.RemoteWakeup(),
		e.IsNotPNP(),
		e.IsochronousInp(),
		e.IsochronousOut(),
		e.SuspendPullDowns(),
		e.UseSerial(),
		e.USBVersion(),
		e.UseUSBVersion(),
		e.MaxPower(),
		e.ChannelType(ChannelA),
		e.ChannelType(ChannelB),
		e.ChannelDriver(ChannelA),
		e.ChannelDriver(ChannelB),
		e.ChannelDriver(ChannelC),
		e.ChannelDriver(ChannelD),
		e.ChannelRS485(ChannelA),
		e.ChannelRS485(ChannelB),
		e.ChannelRS485(ChannelC),
		e.ChannelRS485(ChannelD),
		e.HighCurrent(ChannelAny),
		e.HighCurrent(ChannelA),
		e.HighCurrent(ChannelB),
		e.CBusFunction(0),
		e.CBusFunction(1),
		e.CBusFunction(2),
		e.CBusFunction(3),
		e.CBusFunction(4),
		e.CBusFunction(5),
		e.CBusFunction(6),
		e.CBusFunction(7),
		e.CBusFunction(8),
		e.CBusFunction(9),
		e.Invert(),
	)
}
