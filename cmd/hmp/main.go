package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	arg "github.com/rokath/HMP2020-Go-Bindings/internal/args"
	"github.com/spf13/afero"
)

var (
	// do not initialize, goreleaser will handle that
	version string

	// do not initialize, goreleaser will handle that
	commit string

	// do not initialize, goreleaser will handle that
	date string
)

// main is the entry point.
func main() {
	fSys := &afero.Afero{Fs: afero.NewOsFs()} // os.DirFS("")
	doit(os.Stdout, fSys, os.Args)
}

// doit is the action.
func doit(w io.Writer, fSys *afero.Afero, args []string) {

	// inject values
	arg.Version = version
	arg.Commit = commit
	arg.Date = date

	rand.Seed(time.Now().UnixNano())

	e := arg.Handler(w, fSys, args)
	if e != nil {
		fmt.Fprintln(w, error.Error(e))
	}
}
