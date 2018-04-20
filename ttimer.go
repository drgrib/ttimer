package main

import (
	"flag"

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
	flag.Parse()
	argList := flag.Args()
	timeSet := (len(argList) > 0)
	if timeSet {
		args.t = argList[0]
	} else {
		args.t = "1m"
	}
}

//////////////////////////////////////////////
/// main
//////////////////////////////////////////////

func main() {
	// parse
	d, title := parse.Args(args.t)

	// start timer
	t := agent.Timer{Title: title}
	t.Start(d)

	// run UI
	t.CountDown()
}
