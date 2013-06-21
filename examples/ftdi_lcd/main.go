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
	"github.com/ziutek/lcd/hd44780"
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

	lcd := hd44780.New(d, 4, 20)
}
