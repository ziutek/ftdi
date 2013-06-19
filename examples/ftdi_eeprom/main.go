package main

import (
	"fmt"
	"github.com/ziutek/ftdi"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()

	e := d.EEPROM()
	checkErr(e.Read())
	checkErr(e.Decode())
	fmt.Println(e)
}
