package main

import (
	"fmt"
	"github.com/ziutek/ftdi"
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

	checkErr(d.SetBitmode(0x00, ftdi.ModeReset))
	checkErr(d.SetBitmode(0x01, ftdi.ModeBitbang))

	checkErr(d.SetBaudrate(256)) // bitbang mode, real B/s = 256*16 = 4096B
	//checkErr(d.SetLatencyTimer(128))
	checkErr(d.SetReadChunkSize(1024)) // 250 ms

	checkErr(d.WriteByte(0))

	buf := make([]byte, 1000)
	for {
		n, err := d.Read(buf)
		checkErr(err)
		ones := 0
		for _, b := range buf[:n] {
			if b&0x02 != 0 {
				ones++
			}
		}
		fmt.Printf("%d/%d\n", ones, n)
	}

	checkErr(d.SetBitmode(0x00, ftdi.ModeReset))
}
