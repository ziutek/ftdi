package main

import (
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

	checkErr(d.SetBaudrate(256)) // bitbang mode, actually: 256*32 = 8192 B/s
	//checkErr(d.SetLatencyTimer(4))
	lat, err := d.LatencyTimer()
	checkErr(err)
	log.Println("latency:", lat)
	checkErr(d.SetReadChunkSize(512))

	checkErr(d.WriteByte(0))

	buf := make([]byte, 512)
	for {
		n, err := d.Read(buf)
		checkErr(err)
		ones := 0
		for _, b := range buf[:n] {
			if b&0x02 != 0 {
				ones++
			}
		}
		if n == len(buf) {
			log.Printf("%d", ones)
		} else {
			log.Printf("! %d/%d/%d", ones, n, len(buf))
		}
	}

	checkErr(d.SetBitmode(0x00, ftdi.ModeReset))
}
