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
	label        string
}

// NewPort creates an instance of a serial device type trice receiver
func NewPort(w io.Writer, label string, comPortName string, baud int, dataBits int, parity string, stopBits string, verbose bool) *Port {
	var parityB serial.Parity
	switch strings.ToLower(parity) {
	case "n", "no", "none":
		parityB = serial.NoParity
	case "e", "ev", "even":
		parityB = serial.EvenParity
	case "o", "odd":
		parityB = serial.OddParity
	default:
		log.Fatal("invalid parity value: ", parity, " Accepting case insensitive: n|no|none|e|even||o|odd.")
	}

	var stopB serial.StopBits
	switch strings.ToLower(stopBits) {
	case "1", "one":
		stopB = serial.OneStopBit
	case "1.5":
		stopB = serial.OnePointFiveStopBits
	case "2", "two":
		stopB = serial.TwoStopBits
	default:
		log.Fatalf(" Unknown stop bits value \"%s\". Valid are \"1\", \"1.5\". \"2\"\n", stopBits)
	}

	if !(5 <= dataBits && dataBits <= 9) {
		log.Fatalf("Invalid dataBits value %d. Valid are 5-9\n", dataBits)
	}

	r := &Port{
		port: comPortName,
		serialMode: serial.Mode{
			BaudRate: baud,
			DataBits: dataBits,
			Parity:   parityB,
			StopBits: stopB,
		},
	}

	r.label = label
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
func (p *Port) Read(buf []byte) (int, error) {
	return p.serialHandle.Read(buf)
}

func (p *Port) Write(buf []byte) (int, error) {
	return p.serialHandle.Write(buf)
}

// Close releases port.
func (p *Port) Close() error {
	if p.verbose {
		fmt.Fprintln(p.w, "Closing "+p.label+" COM port.")
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
func GetSerialPorts(w io.Writer, verbose bool) ([]string, error) {
	ports, err := serial.GetPortsList()

	if err != nil {
		fmt.Fprintln(w, err)
		return ports, err
	}
	if len(ports) == 0 {
		if verbose {
			fmt.Fprintln(w, "No serial ports found!")
		}
		return ports, err
	}
	for _, port := range ports {
		pS := NewPort(w, "try", port, 115200, 8, "N", "1", false)
		if pS.Open() {
			pS.Close()
			fmt.Fprintln(w, "Found port: ", port)
		} else {
			fmt.Fprintln(w, "Found port: ", port, "(used)")
		}
	}
	return ports, err
}
