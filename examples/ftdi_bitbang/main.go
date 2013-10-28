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

	checkErr(d.SetBaudrate(192))

	log.Print("WriteByte")
	for i := 0; i < 256; i++ {
		checkErr(d.WriteByte(byte(i)))
	}

	log.Print("Ok")
	time.Sleep(time.Second)

	buf := make([]byte, 256)
	for i := range buf {
		buf[i] = byte(i)
	}

	log.Print("Write")
	_, err = d.Write(buf)
	checkErr(err)

	log.Println("Ok")

	checkErr(d.WriteByte(255))
}
