package args

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/rokath/HMP2020-Go-Bindings/pkg/com"
	"github.com/rokath/HMP2020-Go-Bindings/pkg/hmp"
	"github.com/spf13/afero"
)

var (
	// Version is the program version number and is injected from main package.
	Version string

	// Commit is the program checksum and injected from main package.
	Commit string

	// Date is the compile time and injected from main package.
	Date string

	// Verbose, if true gives additional information for issue fixing.
	Verbose bool

	// Version is the program version number and is injected from main package.
	VersionFlag bool

	// p is the hmp com port.
	p *com.Port

	// scanComPorts is used to find com ports.
	scanComPorts bool
)

func init() {
	flag.BoolVar(&scanComPorts, "s", false, "Scan com ports")
	flag.BoolVar(&VersionFlag, "version", false, "Show version information.")
	flag.BoolVar(&Verbose, "v", false, "Show verbose messages")

}

// Handler is called in main, evaluates args and calls the appropriate functions.
// It returns for program exit.
// All output is written to w.
// fSys is used aas file system.
func Handler(w io.Writer, fSys *afero.Afero, args []string) error {

	if Date == "" { // goreleaser will set Date, otherwise use file info.
		fi, err := fSys.Stat(args[0])
		if nil == err { // On running main tests file-info is invalid, so do not use in that case.
			Date = fi.ModTime().String()
		}
	}

	flag.Parse()

	if flag.NFlag() == 0 { // no CLI flags
		fmt.Println(`Enter "hmp -help"`)
		return nil
	}

	com.Verbose = Verbose
	hmp.Verbose = Verbose

	if VersionFlag {
		if Version != "" {
			fmt.Print("version=", Version)
		}
		if Commit != "" {
			fmt.Print("commit=", Commit)
		}
		fmt.Println("date=", Date)
		return nil
	}

	if scanComPorts {
		_, e := com.GetSerialPorts(w)
		return e
	}
	if com.SerialPortName == "" {
		fmt.Println(`no comport name, enter "hmp -help"`)
		return nil
	} else {
		if Verbose {
			fmt.Println(`CLI switch '-ph' exists, try to connect to HMP...`)
		}
		p = com.NewPort(w, com.SerialPortName, Verbose)
		if !p.Open() {
			return errors.New(com.SerialPortName + " port failure, try with -v for more information")
		}
		defer p.Close()
		e := hmp.Connect(p)
		if e != nil {
			log.Fatal(e)
		}
	}

	// example:

	hmp.OutputOFF(p, -1)
	hmp.SetVoltage(p, 2, "2")
	hmp.SetCurrent(p, 2, "10")
	hmp.OutputON(p, -1)
	hmp.OutputON(p, 2)
	fmt.Print(hmp.Voltage(p, 2))
	fmt.Print(hmp.Current(p, 2))
	for {
		time.Sleep(5 * time.Second)
		hmp.OutputOFF(p, 2)
		time.Sleep(1000 * time.Millisecond)
		hmp.OutputON(p, 2)
	}
}
