package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ziutek/ftdi"
)

func clkbits(baudrate int, clk, clkdiv int) (bestbaud int, encdiv uint32) {
	frac_code := []byte{0, 3, 2, 4, 1, 5, 6, 7}
	var divisor, bestdiv int
	if baudrate >= clk/clkdiv {
		encdiv = 0
		bestbaud = clk / clkdiv
	} else if baudrate >= clk/(clkdiv+clkdiv/2) {
		encdiv = 1
		bestbaud = clk / (clkdiv + clkdiv/2)
	} else if baudrate >= clk/(2*clkdiv) {
		encdiv = 2
		bestbaud = clk / (2 * clkdiv)
	} else {
		/* We divide by 16 to have 3 fractional bits and one bit for rounding */
		divisor = clk * 16 / clkdiv / baudrate
		/* Decide if to round up or down*/
		if divisor&1 != 0 {
			bestdiv = divisor/2 + 1
		} else {
			bestdiv = divisor / 2
		}
		if bestdiv > 0x20000 {
			bestdiv = 0x1ffff
		}
		bestbaud = clk * 16 / clkdiv / bestdiv
		/* Decide if to round up or down*/
		if bestbaud&1 != 0 {
			bestbaud = bestbaud/2 + 1
		} else {
			bestbaud = bestbaud / 2
		}
		encdiv = uint32(bestdiv)>>3 | uint32(frac_code[bestdiv&0x7])<<14
	}
	return
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	ft, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)

	checkErr(ft.SetBitmode(0x0f, ftdi.ModeSyncBB))

	// Baudrate for synchronous bitbang mode.
	//
	// FT232R max baudrate is 3 MBaud, USB speed is 12 Mb/s = 1500 kB/s..
	// In best case: 1308 kB/s fdata + 192 kB/s overhead.
	// Theoretical max USB continuous baudrate in one direction: 1308 kBaud
	//
	// FT232R has 256 B USB Tx buffer and 128 B USB Rx buffer. This is much
	// less than max. 8192 B USB packet so in case of sync bit bang this
	// limits max speed.

	const cs = 8192
	const br = 2 * 256 * 1024 / 16
	checkErr(ft.SetReadChunkSize(cs))
	checkErr(ft.SetWriteChunkSize(cs))
	checkErr(ft.SetLatencyTimer(2))
	checkErr(ft.SetBaudrate(br))

	checkErr(ft.PurgeBuffers())

	data := make([]byte, cs)
	const count = 20
	t1 := time.Now()
	for i := 0; i < count; i++ {
		// Schedule read before Write to improve speed.
		go func() {
			_, err = io.ReadFull(ft, data)
			checkErr(err)
		}()
		_, err := ft.Write(data)
		checkErr(err)
	}
	t2 := time.Now()
	realbaud, encdiv := clkbits(br*4, 48e6, 16)
	div := encdiv & 0xffff
	frac := encdiv >> 16
	fmt.Println(
		"br =", br*16,
		"mes =", int64(len(data))*count*int64(time.Second)/int64(t2.Sub(t1)),
		"cfg =", realbaud, div, frac,
	)
	return
}
