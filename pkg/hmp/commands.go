package hmp

import (
	"strconv"

	"github.com/rokath/HMP2020-Go-Bindings/pkg/com"
)

// OutputOFF disables output on channel n. Valid channel numbers are 1, 2, 3, 4.
// n=0 or -1 switches all channel outputs.
func OutputOFF(p *com.Port, n int) {
	if n <= 0 {
		Command(p, "OUTPUT:GENERAL OFF")
		return
	}
	ch := strconv.Itoa(n)
	Command(p, "INST:SEL OUT"+ch)
	Command(p, "OUTPUT:SELECT OFF")
}

// OutputON enables output on channel n. Valid channel numbers are 1, 2, 3, 4.
// n=0 or -1 switches all channel outputs.
func OutputON(p *com.Port, n int) {
	if n <= 0 {
		Command(p, "OUTPUT:GENERAL ON")
		return
	}
	ch := strconv.Itoa(n)
	Command(p, "INST:SEL OUT"+ch)
	Command(p, "OUTPUT:SELECT ON")
}

// SetVoltage sets voltage on channel n. Valid channel numbers are 1, 2, 3, 4.
func SetVoltage(p *com.Port, n int, v string) {
	ch := strconv.Itoa(n)
	Command(p, "INST:SEL OUT"+ch)
	Command(p, "SOURCE:VOLTAGE:LEVEL "+v)
}

// Voltage returns voltage on channel n. Valid channel numbers are 1, 2, 3, 4.
func Voltage(p *com.Port, n int) (v string) {
	ch := strconv.Itoa(n)
	Command(p, "INST:SEL OUT"+ch)
	return Query(p, "MEASURE:SCALAR:VOLTAGE:DC?")
}

// SetCurrent sets max current on channel n. Valid channel numbers are 1, 2, 3, 4.
func SetCurrent(p *com.Port, n int, v string) {
	ch := strconv.Itoa(n)
	Command(p, "INST:SEL OUT"+ch)
	Command(p, "SOURCE:CURRENT:LEVEL "+v)
}

// Current returns current on channel n. Valid channel numbers are 1, 2, 3, 4.
func Current(p *com.Port, n int) (v string) {
	ch := strconv.Itoa(n)
	Command(p, "INST:SEL OUT"+ch)
	return Query(p, "MEASURE:SCALAR:CURRENT:DC?")
}
