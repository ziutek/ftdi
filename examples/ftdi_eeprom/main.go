package main

import (
	"flag"
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

var (
	set     = flag.Bool("set", false, "Set EEPROM variables")
	vendor  = flag.Int("vendor", 0x0403, "PCI vendor id")
	product = flag.Int("product", 0x6001, "PCI product id")
	invert  = flag.Int(
		"invert",
		0,
		"Set invert flags (use 0x15 to set all lines down)",
	)
	cbusFunction = flag.Int(
		"cbus",
		int(ftdi.CBusIOMode),
		"Function id for all CBUS lines",
	)
	maxCurrent  = flag.Int("maxI", 200, "Maximum USB current (mA)")
	highCurrent = flag.Bool("highDrive", false, "Set high current drive flag")
)

func main() {
	flag.Parse()

	d, err := ftdi.OpenFirst(*vendor, *product, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()

	e := d.EEPROM()
	checkErr(e.Read())
	checkErr(e.Decode())
	fmt.Println(e)

	if !*set {
		return
	}

	modified := false

	if e.Invert() != *invert {
		e.SetInvert(*invert)
		modified = true
	}
	if e.MaxPower() != *maxCurrent {
		e.SetMaxPower(*maxCurrent)
		modified = true
	}
	cbusf := ftdi.CBusFunction(*cbusFunction)
	for n := 0; n < 4; n++ {
		if e.CBusFunction(n) != cbusf {
			e.SetCBusFunction(n, cbusf)
			modified = true
		}
	}
	if e.HighCurrent(ftdi.ChannelAny) != *highCurrent {
		e.SetHighCurrent(ftdi.ChannelAny, *highCurrent)
		modified = true
	}

	if modified {
		checkErr(e.Build())
		checkErr(e.Write())

		checkErr(e.Read())
		checkErr(e.Decode())
		fmt.Println(e)
	}
}
