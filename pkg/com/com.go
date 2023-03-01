// Copyright 2020 Thomas.Hoehenleitner [at] seerose.net
// Use of this source code is governed by a license that can be found in the LICENSE file.

// Package com reads from COM port.
package com

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go.bug.st/serial"
)

var (
	// SerialPortName is the OS specific serial port name used
	SerialPortName string

	// Verbose shows additional information if set true.
	Verbose = false
)

const (
	// Parity is the transmitted bit parity: "even", "odd", "none"
	Parity = "none"

	// BaudRate is the configured baud rate of the serial port. It is set as command line parameter.
	BaudRate = 115200

	// DataBits is the serial port bit count for one "byte".
	DataBits = 8

	// StopBits is the number of stop bits: "1", "1.5", "2"
	StopBits = "1"
)

// COMport is the comport interface type to use different COMports.
type COMport interface {
	Open() bool
	Read(buf []byte) (int, error)
	Write(buf []byte) (int, error)
	Close() error
}

// Port is a serial device trice receiver
type Port struct {
	verbose      bool
	port         string
	serialHandle serial.Port
	serialMode   serial.Mode
	w            io.Writer
}

// NewPort creates an instance of a serial device type trice receiver
func NewPort(w io.Writer, comPortName string, verbose bool) *Port {
	var parity serial.Parity
	switch strings.ToLower(Parity) {
	case "n", "no", "none":
		parity = serial.NoParity
	case "e", "ev", "even":
		parity = serial.EvenParity
	case "o", "odd":
		parity = serial.OddParity
	default:
		log.Fatal("invalid parity value: ", Parity, " Accepting case insensitive: n|no|none|e|even||o|odd.")
	}

	var stopBits serial.StopBits
	switch strings.ToLower(StopBits) {
	case "1", "one":
		stopBits = serial.OneStopBit
	case "1.5":
		stopBits = serial.OnePointFiveStopBits
	case "2", "two":
		stopBits = serial.TwoStopBits
	default:
		log.Fatalf(" Unknown stop bits value \"%s\". Valid are \"1\", \"1.5\". \"2\"\n", StopBits)
	}

	if !(5 <= DataBits && DataBits <= 9) {
		log.Fatalf("Invalid dataBits value %d. Valid are 5-9\n", DataBits)
	}

	r := &Port{
		port: comPortName,
		serialMode: serial.Mode{
			BaudRate: BaudRate,
			DataBits: DataBits,
			Parity:   parity,
			StopBits: stopBits,
		},
	}

	r.w = w
	r.verbose = verbose
	if verbose {
		fmt.Fprintln(w, "New COM port:", r)
	}
	return r
}

// Read blocks until (at least) one byte is received from
// the serial port or an error occurs.
// It stores data received from the serial port into the provided byte array
// buffer. The function returns the number of bytes read.
func (p Port) Read(buf []byte) (int, error) {
	return p.serialHandle.Read(buf)
}

func (p Port) Write(buf []byte) (int, error) {
	return p.serialHandle.Write(buf)
}

// Close releases port.
func (p *Port) Close() error {
	if p.verbose {
		fmt.Fprintln(p.w, "Closing COM port")
	}
	return p.serialHandle.Close()
}

// Open initializes the serial receiver.
//
// It opens a serial port and returns true on success.
func (p *Port) Open() bool {
	var err error
	p.serialHandle, err = serial.Open(p.port, &p.serialMode)
	if err != nil {
		if p.verbose {
			fmt.Fprintln(p.w, err, "try '", os.Args[0], "s' to check for serial ports")
		}
		return false
	}
	return true
}

// GetSerialPorts scans for serial ports.
func GetSerialPorts(w io.Writer) ([]string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		fmt.Fprintln(w, err)
		return ports, err
	}
	if len(ports) == 0 {
		if Verbose {
			fmt.Fprintln(w, "No serial ports found!")
		}
		return ports, err
	}
	for _, port := range ports {
		pS := NewPort(w, port, false)
		if pS.Open() {
			pS.Close()
			fmt.Fprintln(w, "Found port: ", port)
		} else {
			fmt.Fprintln(w, "Found port: ", port, "(used)")
		}
	}
	return ports, err
}

/*
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
*/
