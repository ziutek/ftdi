package main

import (
	"log"

	"github.com/ziutek/ftdi"
)

func main() {
	d, err := ftdi.OpenFirst(0x0403, 0x6011, ftdi.ChannelA)
	if err != nil {
		log.Fatalf("Unable to open FTDI device: %s", err)
	}
	defer d.Close()

	// Channel A is the SPI bus
	if err := d.SetBitmode(0xff, ftdi.ModeMPSSE); err != nil {
		log.Fatalf("Unable to set Bitmode: %s", err)
	}

	clk := byte(1 << 0)
	mosi := byte(1 << 1)
	miso := byte(1 << 2)
	cs := byte(1 << 3)

	outputs := clk | mosi | miso | cs

	// We want a 3MHz clock, and we're not dividing down the 60MHz clock
	speed := ftdi.MPSSEDivValue(3_000_000, false)

	// Set up the SPI Bus - all the pins will be idle high
	mpsse_commands := []byte{
		ftdi.MPSSEDisableDiv5, // Disable /5 divisor to use the 60MHz clock
		ftdi.MPSSETCKDivisor,  // set the clock divisor
		byte(speed & 0xff),    // low byte of clock rate
		byte(speed >> 8),      // high byte of clock rate
		ftdi.MPSSESetBitsLow,  // set low-bit values
		outputs,               // What values to set (all 1)
		outputs,               // Which pins to apply the above values to
	}
	if _, err := d.Write(mpsse_commands); err != nil {
		log.Fatalf("Unable to write MPSSE commands: %s", err)
	}
	tx_data := []byte{1, 2, 3, 4} // Data bytes to appears on the SPI bus

	xfer := []byte{
		ftdi.MPSSESetBitsLow,
		^(cs | clk), // Set all outputs except the chipselect & clock to be high
		outputs,
		ftdi.MPSSEDoWrite | ftdi.MPSSEWriteNeg,
		byte(len(tx_data) & 0xff),
		byte(len(tx_data) >> 8),
	}
	// Add in the actual data we want to send
	xfer = append(xfer, tx_data...)

	// After the transfer, put the pins high, except the clock
	xfer = append(xfer, ftdi.MPSSESetBitsLow)
	xfer = append(xfer, outputs&^clk)
	xfer = append(xfer, outputs)
	if _, err := d.Write(xfer); err != nil {
		log.Fatalf("Unable to write SPI transfer: %s", err)
	}
}
