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
	verbose bool

	// Version is the program version number and is injected from main package.
	VersionFlag bool

	// scanComPorts is used to find com ports.
	scanComPorts bool

	// channel is the output channel.
	channel int

	// voltage is the channel specific voltage.
	voltage string

	// current is the channel specific max current in mA
	current string

	// timeOn is the channel specific ON time in ms.
	timeOn int

	// timeOff is the channel specific OFF time in ms.
	timeOff int
)

func init() {
	flag.BoolVar(&scanComPorts, "s", false, "Scan com ports")
	flag.BoolVar(&VersionFlag, "version", false, "Show version information.")
	flag.BoolVar(&verbose, "v", false, "Show verbose messages")
}

func init() {
	flag.StringVar(&hmp.SerialPortName, "p", "", "Short for hmpPort")
	flag.StringVar(&hmp.SerialPortName, "hmpPort", "", "Use Port for HMP2020 or HMP4040")
	flag.IntVar(&hmp.BaudRate, "hmpBaud", 9600, "Set HMP baud rate.")
	flag.IntVar(&hmp.DataBits, "hmpDataBits", 8, "Set HMP data bit count.")
	flag.StringVar(&hmp.Parity, "hmpParity", "none", "Set parity")
	flag.StringVar(&hmp.StopBits, "hmpStopBits", "1", "Set parity")
	flag.IntVar(&channel, "hmpChannel", 1, "Select HMP channel.")
	flag.IntVar(&channel, "ch", 1, "Short for hmpChannel.")
}

func init() {
	flag.StringVar(&voltage, "V", "1.7", "Set channel output voltage.")
	flag.StringVar(&current, "mA", "10", "Set channel max current.")
	flag.IntVar(&timeOn, "msON", 5000, "Set channel time ON.")
	flag.IntVar(&timeOff, "msOFF", 1000, "Set channel time OFF.")

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
		_, e := com.GetSerialPorts(w, verbose)
		return e
	}
	if hmp.SerialPortName == "" {
		fmt.Println(`no comport name, enter "hmp -help"`)
		return nil
	} else {
		if verbose {
			fmt.Println(`CLI switch '-p' exists, try to connect to HMP...`)
		}

		hmp.Power = hmp.NewDevice(w, hmp.SerialPortName, hmp.BaudRate, hmp.DataBits, hmp.Parity, hmp.StopBits, verbose)
		defer hmp.Power.Close()
		e := hmp.Power.Connect()
		if e != nil {
			log.Fatal(e)
		}
	}

	// example:

	hmp.Power.OutputOFF(-1)
	hmp.Power.SetVoltage(channel, voltage)
	hmp.Power.SetCurrent(channel, current)
	hmp.Power.OutputON(-1)
	for {
		hmp.Power.OutputON(channel)
		time.Sleep(time.Duration(timeOn) * time.Millisecond)
		if verbose {
			fmt.Print("channel:", channel, "V=", hmp.Power.Voltage(channel))
			fmt.Print("channel:", channel, "mA=", hmp.Power.Current(channel))
		}
		hmp.Power.OutputOFF(channel)
		time.Sleep(time.Duration(timeOff) * time.Millisecond)
		hmp.Power.OutputON(channel)
	}
}
