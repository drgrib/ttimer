package main

import (
	"fmt"
	"os"

	"github.com/drgrib/ttimer/agent"
	"github.com/drgrib/ttimer/parse"
)

//////////////////////////////////////////////
/// flags
//////////////////////////////////////////////

var args struct {
	t string
	q bool
}

func init() {
	switch len(os.Args) {
	case 3:
		if os.Args[1] == "-q" {
			args.q = true
			args.t = os.Args[2]
		}
		if os.Args[2] == "-q" {
			args.q = true
			args.t = os.Args[1]
		}
	case 2:
		args.t = os.Args[1]
	default:
		args.t = "1m"
	}
}

//////////////////////////////////////////////
/// main
//////////////////////////////////////////////

func main() {
	// parse
	d, title, err := parse.Args(args.t)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("\nPlease refer to https://github.com/drgrib/ttimer for usage instructions.")
		return
	}

	// start timer
	t := agent.Timer{Title: title}
	t.AutoQuit = args.q
	t.Start(d)

	// run UI
	t.CountDown()
}
