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

	l, err := ftdi.FindAll(0x0403, 0x6001)
	checkErr(err)

	for i, u := range l {
		hexSerial := "0x"
		for i := 0; i < len(u.Serial); i++ {
			hexSerial += fmt.Sprintf("%02x", u.Serial[i])
		}
		fmt.Printf(
			"#%d Manufacturer:'%s' Description:'%s' Serial:'%s'=%s",
			i, u.Manufacturer, u.Description, u.Serial, hexSerial,
		)
		d, err := ftdi.OpenUSBDev(u, ftdi.ChannelAny)
		if err != nil {
			fmt.Printf(" [can't open device: %s]\n", err)
			continue
		}

		if d.Type() == ftdi.TypeR {
			cid, err := d.ChipID()
			if err == nil {
				fmt.Printf(" ChipID:0x%08x\n", cid)
			} else {
				fmt.Printf("[can't read ChipID: %s]\n", err)
			}
		} else {
			fmt.Println()
		}

		d.Close()
	}
}
