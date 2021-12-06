package main

import (
	"dispswitch/pkg/dispswitcher"
	"flag"
)

func main() {
	var pos string
	flag.StringVar(&pos, "side-of-primary", "right", "Location of the Second monitor relative to the Primary [right, left, same, above, bellow]")
	flag.Parse()

	outs := dispswitcher.OutputList{}
	outs.Scan()
	switch pos {
	case "right":
		fallthrough
	case "left":
		pos += "-of" // + outs.Primary().Name
	case "same":
		pos += "-as" // + outs.Primary().Name
	case "above":
	case "bellow":
	default:
		pos = ""
	}

	if outs.Connected().Count() == 1 {
		return
	}

	if outs.Enabled().Count() == 1 {
		outs.Connected().NotPrimary().First().Enable("--"+pos, outs.Primary().Name)
	} else {
		outs.Connected().NotPrimary().First().Disable()
	}
}
