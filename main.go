package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mewkiz/pkg/pathutil"
	"github.com/mewkiz/pkg/term"
)

var (
	// dbg is a logger which logs debug messages to standard error, prepending
	// the "psxisofix:" prefix.
	dbg = log.New(os.Stderr, term.GreenBold("psxisofix:")+" ", 0)
	// warn is a logger which logs warning messages to standard error, prepending
	// the "psxisofix:" prefix.
	warn = log.New(os.Stderr, term.RedBold("psxisofix:")+" ", 0)
)

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: psxisofix FILE.iso")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	for _, path := range flag.Args() {
		if err := fix(path); err != nil {
			log.Fatalf("%+v", err)
		}
	}
}

// fix converts the given PSX ISO to an ISO file compatible with the mount
// command.
func fix(isoPath string) error {
	dbg.Printf("parsing %q", isoPath)
	input, err := ioutil.ReadFile(isoPath)
	if err != nil {
		return err
	}
	const (
		preSkip = 24
		take    = 0x800
		skip    = 0x130
	)
	input = input[preSkip:]
	data := make([]byte, 0, len(input))
	for len(input) > 0 {
		n := take
		if len(input) < n {
			n = len(input)
		}
		data = append(data, input[:n]...)
		input = input[n:]
		n = skip
		if len(input) < n {
			n = len(input)
		}
		input = input[n:]
	}
	fixPath := pathutil.TrimExt(isoPath) + "_fix.iso"
	dbg.Printf("creating %q", fixPath)
	if err := ioutil.WriteFile(fixPath, data, 0644); err != nil {
		return err
	}
	return nil
}
