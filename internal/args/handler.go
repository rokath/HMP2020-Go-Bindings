package args

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

	//  // Verify that a sub-command has been provided: arg[0] is the main command (hmp), arg[1] will be the sub-command.
	//  if len(args) < 2 {
	//  	m := "no args, try: '" + args[0] + " -help'"
	//  	return errors.New(m)
	//  }

	flag.Parse()

	if flag.NFlag() == 0 { // no CLI flags
		fmt.Println(`Enter "hmp -help"`)
		return nil
	}

	hmp.Verbose = Verbose

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

	return nil
}
