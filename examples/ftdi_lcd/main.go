// Pin connection between FT232R and LCD (HD44780 compatible):
// TxD  (DBUS0) <--> D4
// RxD  (DBUS1) <--> D5
// RTS# (DBUS2) <--> D6
// CTS# (DBUS3) <--> D7
// DTR# (DBUS4) <--> E
// DSR# (DBUS5) <--> R/W#
// DCD# (DBUS6) <--> RS
package main

import (
	"fmt"
	"github.com/ziutek/ftdi"
	"github.com/ziutek/lcd/hdc"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()
	checkErr(d.SetBitmode(0xff, ftdi.ModeBitbang))

	// Set output bitbang speed in bytes per second.
	//
	// hdc.Bitbang sends 3 bytes for one 4-bit nibble:
	// 1. with E bit unset, need >= 140 ns
	// 2. with E bit set,   need >= 450 ns
	// 3. with E bit unset, need >= 10 ns
	// Full E cycle time >= 1000 ns
	// We can specify fixed baudrate so we use longest time satisfy all time
	// constrains: 495 ns.
	baudrate := int(1e9) / 495 // bytes/s
	baudrate = 1024
	fmt.Println("Setting baudrate to %d B/s", baudrate)
	checkErr(d.SetBaudrate(baudrate / 16))

	lcd := hdc.NewDevice(hdc.NewBitbang(d), 4, 20)
	checkErr(lcd.Init())
	checkErr(lcd.SetDisplay(hdc.DisplayOn | hdc.CursorOn))

	buf := make([]byte, 80)
	for i := 0; i < 20; i++ {
		buf[i] = 0
		buf[i+20] = 1
		buf[i+40] = 2
		buf[i+60] = 3
	}
	_, err = lcd.Write(buf)
	checkErr(err)
}
