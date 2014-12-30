package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ziutek/ftdi"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func read(r io.Reader, in []byte) {
	os.Stdout.Write([]byte{'r'})
	_, err := io.ReadFull(r, in)
	checkErr(err)
}

func main() {
	ft, err := ftdi.OpenFirst(0x0403, 0x6001, ftdi.ChannelAny)
	checkErr(err)

	checkErr(ft.SetBitmode(0x0f, ftdi.ModeSyncBB))

	// Baudrate for synchronous bitbang mode.
	//
	// FT232R max baudrate is 3 MBaud but in bitbang mode it is additionally
	// limited by USB speed.
	//
	// USB full speed is:
	//     19*64 B / 1ms = 1216000 B/s.
	// Theoretical max. continuous sync bitbang baudrate is:
	//     (1216000 - 2*overhead) Baud / 2 = (608000 - overhead) Baud
	//
	// TODO: Calculate overhead.
	//
	// FT232R has 256 B Tx buffer (for sending to USB host) and 128 B Rx buffer
	// (for receiving from USB host). So if the long term baudrate isn't exceed,
	// the short term (burst) baudrate (for no more than 256 symbols) can be
	// up to (1216000 - overhead) Baud.

	const cs = 64 * 1024
	const br = 65580 * 16
	checkErr(ft.SetReadChunkSize(cs))
	checkErr(ft.SetWriteChunkSize(cs))
	checkErr(ft.SetLatencyTimer(2))
	checkErr(ft.SetBaudrate(br / 16))

	checkErr(ft.PurgeBuffers())

	in := make([]byte, cs)
	out := make([]byte, cs)
	const count = 10

	t1 := time.Now()
	for i := 0; i < count; i++ {
		go read(ft, in)
		// Uncoment following line to see speed for reverse order.
		//time.Sleep(time.Microsecond)
		os.Stdout.Write([]byte{'w'})
		_, err := ft.Write(out)
		checkErr(err)
	}
	t2 := time.Now()

	fmt.Println(
		"\nbr =", br,
		"mes =", cs*count*int64(time.Second)/int64(t2.Sub(t1)),
	)
	return
}
