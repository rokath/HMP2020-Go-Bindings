package hmp

import (
	"strconv"
)

// OutputOFF disables output on channel n. Valid channel numbers are 1, 2, 3, 4.
// n=0 or -1 switches all channel outputs.
func (p *Device) OutputOFF(n int) {
	if n <= 0 {
		p.Command("OUTPUT:GENERAL OFF")
		return
	}
	ch := strconv.Itoa(n)
	p.Command("INST:SEL OUT" + ch)
	p.Command("OUTPUT:SELECT OFF")
}

// OutputON enables output on channel n. Valid channel numbers are 1, 2, 3, 4.
// n=0 or -1 switches all channel outputs.
func (p *Device) OutputON(n int) {
	if n <= 0 {
		p.Command("OUTPUT:GENERAL ON")
		return
	}
	ch := strconv.Itoa(n)
	p.Command("INST:SEL OUT" + ch)
	p.Command("OUTPUT:SELECT ON")
}

// SetVoltage sets voltage on channel n. Valid channel numbers are 1, 2, 3, 4.
func (p *Device) SetVoltage(n int, v string) {
	ch := strconv.Itoa(n)
	p.Command("INST:SEL OUT" + ch)
	p.Command("SOURCE:VOLTAGE:LEVEL " + v)
}

// Voltage returns voltage on channel n. Valid channel numbers are 1, 2, 3, 4.
func (p *Device) Voltage(n int) (v string) {
	ch := strconv.Itoa(n)
	p.Command("INST:SEL OUT" + ch)
	return p.Query("MEASURE:SCALAR:VOLTAGE:DC?")
}

// SetCurrent sets max current on channel n. Valid channel numbers are 1, 2, 3, 4.
func (p *Device) SetCurrent(n int, v string) {
	ch := strconv.Itoa(n)
	p.Command("INST:SEL OUT" + ch)
	p.Command("SOURCE:CURRENT:LEVEL " + v)
}

// Current returns current on channel n. Valid channel numbers are 1, 2, 3, 4.
func (p *Device) Current(n int) (v string) {
	ch := strconv.Itoa(n)
	p.Command("INST:SEL OUT" + ch)
	return p.Query("MEASURE:SCALAR:CURRENT:DC?")
}
