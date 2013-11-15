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
		d, err := ftdi.OpenUSBDev(u, ftdi.ChannelAny)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't open #%d: %s\n", i, err)
			continue
		}

		if d.Type() == ftdi.TypeR {
			cid, err := d.ChipID()
			if err == nil {
				fmt.Printf("%d: 0x%08x\n", i, cid)
			} else {
				fmt.Fprintf(
					os.Stderr,
					"Can't read ChipID for #%d: %s\n",
					i, err,
				)
			}
		}
		d.Close()
	}
}
