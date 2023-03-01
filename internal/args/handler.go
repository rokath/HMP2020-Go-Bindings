package args

import (
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

	// c is the hmp com port.
	//p *com.Port

	p *hmp.Device

	// scanComPorts is used to find com ports.
	scanComPorts bool
)

func init() {
	flag.BoolVar(&scanComPorts, "s", false, "Scan com ports")
	flag.BoolVar(&VersionFlag, "version", false, "Show version information.")
	flag.BoolVar(&Verbose, "v", false, "Show verbose messages")
}

func init() {
	flag.StringVar(&hmp.SerialPortName, "hmpPort", "", "Use Port for HMP2020 or HMP4040")
	flag.StringVar(&hmp.SerialPortName, "p", "", "Short for hmpPort")
	flag.IntVar(&hmp.BaudRate, "hmpBaud", 115200, "Set HMP baud rate.")
	flag.IntVar(&hmp.DataBits, "hmpDataBits", 8, "Set HMP data bit count.")
	flag.StringVar(&hmp.Parity, "hmpParity", "none", "Set parity")
	flag.StringVar(&hmp.StopBits, "hmpStopBits", "1", "Set parity")
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
	if hmp.SerialPortName == "" {
		fmt.Println(`no comport name, enter "hmp -help"`)
		return nil
	} else {
		if Verbose {
			fmt.Println(`CLI switch '-p' exists, try to connect to HMP...`)
		}

		p = hmp.NewDevice(w, hmp.SerialPortName, hmp.BaudRate, hmp.DataBits, hmp.Parity, hmp.StopBits, Verbose)

		e := p.Connect()
		if e != nil {
			log.Fatal(e)
		}

		defer p.Close()
	}

	// example:

	p.OutputOFF(-1)
	p.SetVoltage(2, "2")
	p.SetCurrent(2, "10")
	p.OutputON(-1)
	p.OutputON(2)
	fmt.Print(p.Voltage(2))
	fmt.Print(p.Current(2))
	for {
		time.Sleep(3 * time.Second)
		p.OutputOFF(2)
		time.Sleep(1000 * time.Millisecond)
		p.OutputON(2)
	}
}
