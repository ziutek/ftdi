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
	"github.com/ziutek/ftdi"
	"github.com/ziutek/lcd/hdc"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()
	checkErr(d.SetBitmode(0xff, ftdi.ModeBitbang))
	// Set output speed in bytes per second. For communication with
	// HD44780 in 4-bit mode there are two bytes send for one 4-bit
	// nibble (first with E bit set, second with E bit unset).
	boudrate := 1024 // bytes/s
	checkErr(d.SetBaudrate(boudrate / 16))

	lcd := hdc.NewDevice(hdc.NewBitbang(d), 4, 20)
	checkErr(lcd.Init())
	checkErr(lcd.SetDisplay(hdc.DisplayOn | hdc.CursorOn))

	checkErr(lcd.Write([]byte{
	'1','1','1','1','1','1','1'i,'1','1','1','1','1','1','1'
	}))
	}
}
