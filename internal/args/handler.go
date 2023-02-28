package args

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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
)

func init() {
	flag.BoolVar(&Verbose, "v", false, "Show verbose messages")
	flag.BoolVar(&VersionFlag, "version", false, "Show version information.")

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

	if hmp.ComPort != "" {
		if Verbose {
			fmt.Println(`CLI switch '-ph' exists, try to connect to HMP2020.`)
		}
		hmp.Com = hmp.NewCOMPortTarm(os.Stdout)
		defer hmp.Com.Close()
		e := hmp.Connect()
		if e != nil {
			log.Fatal(e)
		}
	}

	// example:
	hmp.OutputGeneralOff()
	hmp.SetVoltageChannel2("24")
	hmp.SetCurrentChannel2("1000")
	hmp.OutputGeneralOn()
	hmp.SetOutputChannel2On()
	fmt.Print(hmp.VoltageChannel2())
	fmt.Print(hmp.CurrentChannel2())

	for {
		time.Sleep(5 * time.Second)
		hmp.SetOutputChannel2Off()
		time.Sleep(1000 * time.Millisecond)
		hmp.SetOutputChannel2On()
	}
}
