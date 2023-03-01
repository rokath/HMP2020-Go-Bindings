package hmp

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/rokath/HMP2020-Go-Bindings/pkg/com"
)

var (
	// Verbose, if true gives additional information for issue fixing.
	Verbose bool

	Port com.Port
)

func init() {
	flag.StringVar(&com.SerialPortName, "portHMP", "", "Use Port for HMP2020 or HMP4040")
	flag.StringVar(&com.SerialPortName, "ph", "", "Short for portHMP")

}

// Connect tries to get contact to HMP2020.
func Connect(p *com.Port) error {
	response, e := SendAndReceive(p, "*IDN?\n", 100) // 50 is the fractal border.
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
func Send(p *com.Port, cmd string, ms time.Duration) {
	b := []byte(cmd)
	m, e := p.Write(b)
	if m != len(b) || e != nil {
		fmt.Printf("Wrote only %d bytes and not %d bytes, error %v", m, len(b), e)
	}
	time.Sleep(ms * time.Millisecond)
}

// SendAndReceive transmits cmd, waits ms milliseconds and returns response in r.
// The returned string is "as is" from HMP2020.
func SendAndReceive(p *com.Port, cmd string, ms time.Duration) (r string, e error) {
	Send(p, cmd, ms)      // needs a line end
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

func Query(p *com.Port, cmd string) (response string) {
	if Verbose {
		fmt.Println("query:", cmd)
	}
	response, e := SendAndReceive(p, cmd+"\r\n", 500) // 50 is the fractal border.
	if e != nil {
		fmt.Println(e) // log.Fatal(e)
	}
	return
}

func Command(p *com.Port, cmd string) {
	if Verbose {
		fmt.Println("command:", cmd)
	}
	Send(p, cmd+"\r\n", 100) // 50 is the fractal border.
}
