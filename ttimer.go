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
}

func init() {
	if len(os.Args) > 1 {
		args.t = os.Args[1]
	} else {
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
	t.Start(d)

	// run UI
	t.CountDown()
}
