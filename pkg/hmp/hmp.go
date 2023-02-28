// Package hmp provides remote access to Rohde&Schwarz HMP2020 power supply
package hmp

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	tarm "github.com/tarm/serial"
)

var (

	// ComPort holds the COMPort name the T-10A is connected to.
	ComPort string

	// Com holds the serial port values.
	Com *Tarm

	// Verbose, if true gives additional information for issue fixing.
	Verbose bool
)

func init() {
	flag.StringVar(&ComPort, "portHMP2020", "", "Use Port for HMP2020")
	flag.StringVar(&ComPort, "ph", "", "Short for portHMP2020")

}

// Tarm is a serial device type.
type Tarm struct {
	config  tarm.Config
	stream  *tarm.Port
	w       io.Writer
	Verbose bool
}

// NewCOMPortTarm creates an instance of a serial device type trice receiver.
func NewCOMPortTarm(w io.Writer) *Tarm {
	var p = new(Tarm)
	p.w = w
	p.Verbose = Verbose
	p.config.Name = ComPort
	p.config.Baud = 115200 // fixed baud rate
	p.config.ReadTimeout = 1000 * time.Millisecond
	p.config.Size = 8 // "character length"
	p.config.StopBits = 1
	p.config.Parity = tarm.ParityNone // no parity
	var err error
	p.stream, err = tarm.OpenPort(&p.config)
	if err != nil {
		if p.Verbose {
			fmt.Fprintln(w, ComPort, "not found")
			fmt.Fprintln(w, "try 'trice scan' to find com ports as one possibility")
		}
		log.Fatal(err)
	}
	if p.Verbose {
		fmt.Fprintln(w, "NewCOMPortTarm:", p.config)
	}
	return p
}

// Close returns an error in case of failure.
func (p *Tarm) Close() error {
	if p.Verbose {
		fmt.Fprintln(p.w, "Closing Tarm COM port")
	}
	return p.stream.Close()
}

// Read blocks until (at least) one byte is received from
// the serial port or an error occurs.
// It stores data received from the serial port into the provided byte array
// buffer. The function returns the number of bytes read.
func (p *Tarm) Read(buf []byte) (int, error) {
	return p.stream.Read(buf)
}

// Write ...
func (p *Tarm) Write(buf []byte) (int, error) {
	return p.stream.Write(buf)
}

// Connect tries to get contact to HMP2020.
func Connect() error {
	response, e := SendAndReceive("*IDN?\n", 100) // 50 is the fractal border.
	if e != nil {
		return e
	}
	exp := "ROHDE&SCHWARZ,HMP2020"
	if response[:len(exp)] != exp {
		return fmt.Errorf("HMP2020 not connected")
	}
	if Verbose {
		fmt.Println(strings.TrimSuffix(response, "\n"), "connected")
	}
	return nil
}

// Send transmits cmd and waits ms afterwards before returning.
func Send(cmd string, ms time.Duration) {
	b := []byte(cmd)
	m, e := Com.Write(b)
	if m != len(b) || e != nil {
		fmt.Printf("Wrote only %d bytes and not %d bytes, error %v", m, len(b), e)
	}
	time.Sleep(ms * time.Millisecond)
}

// SendAndReceive transmits cmd, waits ms milliseconds and returns n read bytes in b. Set ms = 0 if no answer is expected.
// The returned string is "as is" from HMP2020.
func SendAndReceive(cmd string, ms time.Duration) (string, error) {
	Send(cmd, ms)         // needs a line end
	b := make([]byte, 64) // 32 needed
	n, e := Com.Read(b)
	if e != nil {
		fmt.Println(e)
	}
	if n == 0 {
		e = errors.New("no answer from HMP2020")
		return "", e
	}
	b = b[:n]
	r := string(b)
	return r, e
}

func Query(cmd string) (response string) {
	if Verbose {
		fmt.Println("query:", cmd)
	}
	response, e := SendAndReceive(cmd+"\r\n", 500) // 50 is the fractal border.
	if e != nil {
		fmt.Println(e) // log.Fatal(e)
	}
	return
}

func Command(cmd string) {
	if Verbose {
		fmt.Println("command:", cmd)
	}
	Send(cmd+"\r\n", 100) // 50 is the fractal border.
}

func OutputGeneralOff() {
	Command("OUTPUT:GENERAL OFF")
}

func OutputGeneralOn() {
	Command("OUTPUT:GENERAL ON")
}

func SetOutputChannel1Off() {
	Command("INST:SEL OUT1")
	Command("OUTPUT:SELECT OFF")
}

func SetOutputChannel1On() {
	Command("INST:SEL OUT1")
	Command("OUTPUT:SELECT ON")
}

func SetOutputChannel2Off() {
	Command("INST:SEL OUT2")
	Command("OUTPUT:SELECT OFF")
}

func SetOutputChannel2On() {
	Command("INST:SEL OUT2")
	Command("OUTPUT:SELECT ON")
}

func SetVoltageChannel1(v string) {
	Command("INST:SEL OUT1")
	Command("SOURCE:VOLTAGE:LEVEL " + v)
}

func SetVoltageChannel2(v string) {
	Command("INST:SEL OUT2")
	Command("SOURCE:VOLTAGE:LEVEL " + v)
}

func SetCurrentChannel1(v string) {
	Command("INST:SEL OUT1")
	Command("SOURCE:CURRENT:LEVEL:IMM " + v)
}

func SetCurrentChannel2(v string) {
	Command("INST:SEL OUT2")
	Command("SOURCE:CURRENT:LEVEL:IMM " + v)
}

func VoltageChannel1() (v string) {
	Command("INST:SEL OUT1")
	return Query("MEASURE:SCALAR:VOLTAGE:DC?")
}

func VoltageChannel2() (v string) {
	Command("INST:SEL OUT2")
	return Query("MEASURE:SCALAR:VOLTAGE:DC?")
}

func CurrentChannel1() (v string) {
	Command("INST:SEL OUT1")
	return Query("MEASURE:SCALAR:CURRENT:DC?")
}

func CurrentChannel2() (v string) {
	Command("INST:SEL OUT2")
	return Query("MEASURE:SCALAR:CURRENT:DC?")
}
