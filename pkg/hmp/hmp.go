package hmp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/rokath/HMP2020-Go-Bindings/pkg/com"
)

var (

	// SerialPortName is the OS specific serial port name used for HMP device.
	SerialPortName string

	// Parity is the transmitted bit parity: "even", "odd", "none"
	Parity string

	// BaudRate is the configured baud rate of the serial port. It is set as command line parameter.
	BaudRate int

	// DataBits is the serial port bit count for one "byte".
	DataBits int

	// StopBits is the number of stop bits: "1", "1.5", "2"
	StopBits = "1"

	// Power is the power device to be controlled.
	Power *Device
)

type Device struct {

	// Ports is the assigned communictation port.
	Port *com.Port

	verbose bool

	w io.Writer
}

func NewDevice(w io.Writer, serialPortName string, baudRate int, dataBits int, parity string, stopBits string, verbose bool) (p *Device) {
	p = new(Device)
	p.verbose = verbose
	p.w = w
	p.Port = com.NewPort(w, "HMP", SerialPortName, baudRate, dataBits, parity, stopBits, verbose)

	if !p.Port.Open() {
		log.Fatal(errors.New(serialPortName + " port failure, try with -v for more information"))
	}
	return
}

// Read blocks until (at least) one byte is received from
// the serial port or an error occurs.
// It stores data received from the serial port into the provided byte array
// buffer. The function returns the number of bytes read.
func (p *Device) Read(buf []byte) (int, error) {
	return p.Port.Read(buf)
}

func (p *Device) Write(buf []byte) (int, error) {
	return p.Port.Write(buf)
}

// Close releases port.
func (p *Device) Close() error {
	if p.verbose {
		fmt.Fprintln(p.w, "Closing COM port")
	}
	return p.Port.Close()
}

// Connect tries to get contact to HMP2020.
func (p *Device) Connect() error {
	response, e := p.SendAndReceive("*IDN?\n", 100) // 50 is the fractal border.
	if e != nil {
		return e
	}
	exp := "ROHDE&SCHWARZ,HMP2020"
	if response[:len(exp)] != exp {
		return fmt.Errorf("HMP2020 not connected")
	}
	if p.verbose {
		fmt.Println(strings.TrimSuffix(response, "\n"), "connected")
	}
	return nil
}

// Send transmits cmd and waits ms afterwards before returning.
func (p *Device) Send(cmd string, ms time.Duration) {
	b := []byte(cmd)
	m, e := p.Write(b)
	if m != len(b) || e != nil {
		fmt.Printf("Wrote only %d bytes and not %d bytes, error %v", m, len(b), e)
	}
	time.Sleep(ms * time.Millisecond)
}

// SendAndReceive transmits cmd, waits ms milliseconds and returns response in r.
// The returned string is "as is" from HMP2020.
func (p *Device) SendAndReceive(cmd string, ms time.Duration) (r string, e error) {
	p.Send(cmd, ms)       // needs a line end
	b := make([]byte, 64) // 32 needed
	var n int
	n, e = p.Read(b)
	if e != nil {
		fmt.Println(e)
	}
	if n == 0 {
		e = errors.New("no answer from HMP2020")
		return "", e
	}
	b = b[:n]
	r = string(b)
	return r, e
}

func (p *Device) Query(cmd string) (response string) {
	if p.verbose {
		fmt.Println("query:", cmd)
	}
	response, e := p.SendAndReceive(cmd+"\r\n", 500) // 50 is the fractal border.
	if e != nil {
		fmt.Println(e) // log.Fatal(e)
	}
	return
}

func (p *Device) Command(cmd string) {
	if p.verbose {
		fmt.Println("command:", cmd)
	}
	p.Send(cmd+"\r\n", 100) // 50 is the fractal border.
}
