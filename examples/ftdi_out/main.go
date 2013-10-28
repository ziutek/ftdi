package main

import (
	"fmt"
	"github.com/ziutek/ftdi"
	"os"
	"strconv"
)

func errorExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usageExit() {
	fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "BYTE")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		usageExit()
	}
	b, err := strconv.ParseUint(os.Args[1], 0, 8)
	errorExit(err)

	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	errorExit(err)
	defer d.Close()

	errorExit(d.SetBitmode(0xff, ftdi.ModeBitbang))

	errorExit(d.WriteByte(byte(b)))
}
