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

	modified := false

	invert := 0x15
	if e.Invert() != invert {
		e.SetInvert(invert)
		modified = true
	}
	maxCurrent := 200 //mA
	if e.MaxPower() != maxCurrent {
		e.SetMaxPower(maxCurrent)
		modified = true
	}
	cbusFunction := ftdi.CBusIOMode
	for n := 0; n < 4; n++ {
		if e.CBusFunction(n) != cbusFunction {
			e.SetCBusFunction(n, cbusFunction)
			modified = true
		}
	}

	if modified {
		checkErr(e.Build())
		checkErr(e.Write())

		checkErr(e.Read())
		checkErr(e.Decode())
		fmt.Println(e)
	}
}
