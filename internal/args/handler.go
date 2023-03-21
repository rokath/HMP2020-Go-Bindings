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
)

var (
	verbose        bool   // Verbose, if true gives additional information for issue fixing.
	VersionFlag    bool   // Version is the program version number and is injected from main package.
	scanComPorts   bool   // scanComPorts is used to find com ports.
	channel        int    // channel is the output channel.
	voltage        string // voltage is the channel specific voltage.
	current        string // current is the channel specific max current in A.
	msON           int    // timeOn is the channel specific ON time in ms.
	msOFF          int    // timeOff is the channel specific OFF time in ms.
	stepDirection  string // The StepDirection is "UP" or "DOWN".
	stepSize       string // StepSize has Volt as unit like "0.1" for 100 mV.
	stepMsDuration int    // StepMsDuration is the duration of one step in milliseconds.
	count          int    // count is the count of executed steps or switch loops.
	output         bool   // Enable channel output.
	beep           bool
)

func init() {
	flag.BoolVar(&scanComPorts, "s", false, "Scan com ports")
	flag.BoolVar(&VersionFlag, "version", false, "Show version information.")
	flag.BoolVar(&verbose, "v", false, "Show verbose messages.")
	flag.BoolVar(&beep, "beep", false, "Beeps after command execution.")
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
	flag.StringVar(&current, "A", "0.010", "Set channel max current.")

	flag.BoolVar(&output, "output", false, "Enable output.")

	flag.IntVar(&msON, "msON", 5000, "Set channel time ON.")
	flag.IntVar(&msOFF, "msOFF", 1000, "Set channel time OFF.")

	flag.StringVar(&stepDirection, "stepDirection", "DOWN", "Set step direction to UP or DOWN.")
	flag.StringVar(&stepSize, "stepSize", "0.1", "Set step size in Volt.")
	flag.IntVar(&stepMsDuration, "stepMs", 1000, "Set step duration in milliseconds.")

	flag.IntVar(&count, "count", 10, "Set step/switch count.")
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
	}
	if verbose {
		fmt.Println(`CLI switch '-p' exists, try to connect to HMP...`)
	}

	hmp.Power = hmp.NewDevice(w, hmp.SerialPortName, hmp.BaudRate, hmp.DataBits, hmp.Parity, hmp.StopBits, verbose)
	defer hmp.Power.Close()
	e := hmp.Power.Connect()
	if e != nil {
		log.Fatal(e)
	}

	if isFlagPassed("V") {
		if verbose {
			fmt.Printf("SetVoltage(channel %d, %s V)\n", channel, voltage)
		}
		hmp.Power.SetVoltage(channel, voltage)
	}
	if isFlagPassed("A") {
		if verbose {
			fmt.Printf("SetCurrent(channel %d, %s A)\n", channel, current)
		}
		hmp.Power.SetCurrent(channel, current)
	}
	if isFlagPassed("output") && !output {
		if verbose {
			fmt.Printf("OutputOFF(channel %d)\n", channel)
		}
		hmp.Power.OutputOFF(channel)
	}
	if isFlagPassed("output") && output {
		if verbose {
			fmt.Printf("OutputON(channel %d)\n", channel)
		}
		hmp.Power.OutputON(channel)
	}
	if verbose {
		fmt.Print("channel:", channel, "V=", hmp.Power.Voltage(channel))
		fmt.Print("channel:", channel, "A=", hmp.Power.Current(channel))
	}

	if isFlagPassed("stepDirection") || isFlagPassed("stepSize") || isFlagPassed("stepMs") {
		if verbose {
			fmt.Printf("VoltageRamp(channel %d, %s, %s deltaV, %d ms, %d steps))\n", channel, stepDirection, stepSize, stepMsDuration, count)
		}
		dur := time.Duration(stepMsDuration) * time.Millisecond
		hmp.Power.VoltageRamp(channel, stepDirection, stepSize, dur, count)
		if verbose {
			fmt.Print("channel:", channel, "V=", hmp.Power.Voltage(channel))
			fmt.Print("channel:", channel, "A=", hmp.Power.Current(channel))
		}
	}

	if isFlagPassed("msON") || isFlagPassed("msOFF") {
		if verbose {
			fmt.Printf("channel %d: OutputON(%d ms), OutputOFF(%d ms)\n", channel, msON, msOFF)
		}
		for count > 0 {
			count--
			hmp.Power.OutputON(channel)
			time.Sleep(time.Duration(msON) * time.Millisecond)
			if verbose {
				fmt.Print("channel:", channel, "V=", hmp.Power.Voltage(channel))
				fmt.Print("channel:", channel, "A=", hmp.Power.Current(channel))
			}
			hmp.Power.OutputOFF(channel)
			time.Sleep(time.Duration(msOFF) * time.Millisecond)
		}
	}
	if beep {
		hmp.Power.Command("SYST:BEEP")
	}
	return nil
	// examples:
	//  hmp.Power.OutputOFF(-1) // all channels OFF
	//  hmp.Power.OutputON(-1)  // all channels ON
}

// https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
