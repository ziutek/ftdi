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
	"log"
	"time"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const (
	E  = 1 << 4
	RW = 1 << 5
	RS = 1 << 6
)

func wait() {
	time.Sleep(10 * time.Millisecond)
}

func init4bit(d *ftdi.Device) {
	checkErr(d.WriteByte(E | 0x02))
	wait()
	checkErr(d.WriteByte(0x02))
	wait()
}

func send(d *ftdi.Device, rs, b byte) {
	h := b >> 4
	l := b & 0x0f
	checkErr(d.WriteByte(E | rs | h))
	wait()
	checkErr(d.WriteByte(rs | h))
	wait()
	checkErr(d.WriteByte(E | rs | l))
	wait()
	checkErr(d.WriteByte(rs | l))
	wait()
}

func sendCmd(d *ftdi.Device, b byte) {
	send(d, 0, b)
}

func sendData(d *ftdi.Device, b byte) {
	send(d, RS, b)
}

func main() {
	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()

	checkErr(d.SetBitmode(0xff, ftdi.ModeBitbang))

	init4bit(d)
	sendCmd(d, 0x28) // Display ON, Cursor On, Cursor Blinking
	sendCmd(d, 0x0f) // Entry Mode, Increment cursor position, No display shift
	sendCmd(d, 0x06) // Cursor moving direction
	sendCmd(d, 0x01) // Clear screen
	for i := byte(0); i < 80; i++ {
		sendData(d, '0'+i)
	}

}
