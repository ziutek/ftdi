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
	"time"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	// Good values for JHD204A (136 FPS on 4x20 display)
	baudrate := 1 << 17
	waitTicks := 6

	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()
	checkErr(d.SetBitmode(0xff, ftdi.ModeBitbang))
	checkErr(d.SetBaudrate(baudrate / 16))

	w := hdc.NewBitbang(d, waitTicks)
	// In conservative mode baudrate can be 1 << 23 on JHD204A
	//w.FastMode(false)

	lcd := hdc.NewDevice(w, 4, 20)
	checkErr(lcd.Init())
	checkErr(lcd.SetDisplay(hdc.DisplayOn | hdc.CursorOn))

	buf1 := make([]byte, 80)
	for i := 0; i < 20; i++ {
		buf1[i] = '0'
		buf1[i+20] = '2'
		buf1[i+40] = '1'
		buf1[i+60] = '3'
	}
	buf2 := make([]byte, 80)
	for i := 0; i < 80; i++ {
		buf2[i] = ' '
	}
	n := 20
	t := time.Now()
	for i := 0; i < n; i++ {
		_, err = lcd.Write(buf2)
		checkErr(err)
		_, err = lcd.Write(buf1)
		checkErr(err)
	}
	fmt.Printf(
		"%.2f FPS\n",
		float64(2*n)*float64(time.Second)/float64(time.Now().Sub(t)),
	)

	for i := 0; i < 4; i++ {
		for i := 0; i < 20; i++ {

		}
	}
}
