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

const (
	baudrate  = 192
	Bps       = baudrate * 32 // for FT323R
	chunkSize = 2 * 64
	bufLen    = 2 * 62 // 62 because of 2 status bytes in every 64B USB packet
	dt        = time.Second / Bps
)

func find1(buf []byte) int {
	for i, b := range buf {
		if b&0x02 != 0 {
			return i
		}
	}
	return -1
}

func find0(buf []byte) int {
	for i, b := range buf {
		if b&0x02 == 0 {
			return i
		}
	}
	return -1
}

func main() {
	d, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)
	defer d.Close()

	checkErr(d.SetBitmode(0x00, ftdi.ModeReset))
	checkErr(d.SetBitmode(0x01, ftdi.ModeBitbang))

	checkErr(d.SetBaudrate(baudrate))

	//checkErr(d.SetLatencyTimer(4))
	lat, err := d.LatencyTimer()
	checkErr(err)
	log.Println("latency:", lat)

	checkErr(d.SetReadChunkSize(chunkSize))

	checkErr(d.WriteByte(0))

	buf := make([]byte, bufLen)
	pulseLen := 0
	for {
		n, err := d.Read(buf)
		checkErr(err)
		data := buf
		if n != len(buf) {
			log.Printf("Partial buffer: %d/%d", n, len(buf))
			data = buf[:n]
		}

	analysis:
		if pulseLen > 0 {
			if i := find0(data); i != -1 {
				pulseLen += i

				ms := time.Duration(pulseLen) * dt / time.Millisecond
				c := '?'
				if ms >= 40 && ms <= 130 {
					c = '0'
				} else if ms >= 140 && ms <= 250 {
					c = '1'
				}
				log.Printf("%c (%s)", c, time.Duration(pulseLen)*dt)

				data = data[i+1:]
				pulseLen = 0
				goto analysis
			} else {
				pulseLen += len(data)
			}
		} else {
			if i := find1(data); i != -1 {
				data = data[i+1:]
				pulseLen = 1
				goto analysis
			}
		}
	}

	checkErr(d.SetBitmode(0x00, ftdi.ModeReset))
}
