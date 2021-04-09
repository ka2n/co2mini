package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ka2n/co2mini"
)

var (
	watch  *bool   = flag.Bool("watch", false, "continuous watching value")
	format *string = flag.String("format", "plain", "select format: plain|json")
)

func main() {
	flag.Parse()
	dev := co2mini.FindDevice()
	if dev == nil {
		fmt.Fprintln(os.Stderr, "device not found")
		os.Exit(1)
	}

	var output co2mini.OutputWriter
	switch *format {
	case "json":
		output = co2mini.JSONOutputWriter{}
	case "plain":
		output = co2mini.PlainOutputWriter{}
	default:
		fmt.Fprintln(os.Stderr, "invalid format")
		os.Exit(1)
	}

	if *watch {
		if err := co2mini.Watch(*dev, output); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		if err := co2mini.Oneshot(*dev, output); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
