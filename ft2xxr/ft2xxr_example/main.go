package main

import (
	"fmt"
	"github.com/ziutek/ftdi/ft2xxr"
	"os"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func main() {
	dl, err := ft2xxr.FindDevices(0x0403, 0x6001)
	checkErr(err)
	fmt.Println("Found", len(dl), "devices:")
	for i, d := range dl {
		c, err := d.Connect()
		checkErr(err)
		desc, err := c.Description()
		checkErr(err)
		fmt.Printf("%d: %s\n", i, desc)
	}
}
