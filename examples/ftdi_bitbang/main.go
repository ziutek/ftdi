package main

import (
	"github.com/ziutek/ftdi"
	"log"
	"time"
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

	checkErr(d.SetBitmode(0xff, ftdi.ModeBitbang))

	Bps := 115200

	log.Print("WriteByte")
	for i := 0; i < 256; i++ {
		checkErr(d.WriteByte(byte(i)))
		time.Sleep(time.Second / time.Duration(Bps))
	}

	log.Print("Ok")
	time.Sleep(time.Second)

	buf := make([]byte, 80 * 2 * 3)
	for i := range buf {
		buf[i] = byte(i)
	}

	checkErr(d.SetBaudrate(Bps / 16)) // bitbang mode so real Bps / 16

	log.Print("Write")
	t := time.Now()
	n := 80
	for i := 0; i < n; i++ {
		_, err := d.Write(buf)
		checkErr(err)
	}
	dt := time.Now().Sub(t)

	n *= len(buf)
	log.Printf(
		"%d bytes written in %s (%d B/s)",
		n, dt, time.Duration(n)*time.Second/dt,
	)
}
