package main

import (
	"flag"
	"fmt"
	ui "github.com/gizak/termui"
	"math"
	"strconv"
	"time"
)

//////////////////////////////////////////////
/// flags
//////////////////////////////////////////////

var args struct {
	t string
	z string
	N bool
}

func init() {
	flag.StringVar(
		&args.t, "t", "1", "time string")
	flag.StringVar(
		&args.z, "z", "-0800", "timezone string")
	flag.BoolVar(
		&args.N, "N", false, "use notifications")
	flag.Parse()
}

//////////////////////////////////////////////
/// util
//////////////////////////////////////////////

func mustBeNil(err error) {
	if err != nil {
		panic(err)
	}
}

//////////////////////////////////////////////
/// argsToDuration
//////////////////////////////////////////////

func argsToDuration(t, z string) time.Duration {
	var d time.Duration
	d = time.Duration(6 * time.Second)
	if len(t) == 1 {
		// simple minute timer
		minutes, err := strconv.Atoi(t)
		mustBeNil(err)
		d = time.Duration(minutes) * time.Minute
	}
	fmt.Println(len(z))
	return d
}

//////////////////////////////////////////////
/// Timer
//////////////////////////////////////////////

type Timer struct {
	title    string
	duration time.Duration
	end      time.Time
	status   string
}

func (t *Timer) start(d time.Duration) {
	t.duration = d
	t.end = time.Now().Add(t.duration)
}

func (t *Timer) update() {
	t.status = "[finished]"
	now := time.Now()
	if !now.After(t.end) {
		left := t.end.Sub(now)
		floorSeconds := math.Floor(left.Seconds())
		rounded := time.Duration(floorSeconds) * time.Second
		t.status = fmt.Sprintf("%v", rounded)
	}
}

func (t *Timer) countDown() {
	// init and close
	err := ui.Init()
	mustBeNil(err)
	defer ui.Close()

	// init cell
	cell := ui.NewPar("")
	cell.TextFgColor = ui.ColorDefault
	cell.Border = false
	cell.X = 2
	cell.Y = 1
	cell.Width = ui.TermWidth()
	cell.Height = ui.TermHeight()

	// draw
	banner := fmt.Sprintf("== %s ==", t.title)
	draw := func(tick int) {
		t.update()
		// render
		cell.Text = fmt.Sprintf("%s\n%v",
			banner,
			t.status)
		ui.Render(cell)
	}

	// handle update
	ms := 50
	timerPath := fmt.Sprintf("/timer/%vms", ms)
	ui.Merge("timer", ui.NewTimerCh(
		time.Duration(ms)*time.Millisecond))
	ui.Handle(timerPath, func(e ui.Event) {
		tick := e.Data.(ui.EvtTimer)
		draw(int(tick.Count))
	})

	// handle resize
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		cell.Width = ui.TermWidth()
		cell.Height = ui.TermHeight()
		ui.Clear()
		ui.Render(ui.Body)
	})

	// handle quit
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	// start loop
	ui.Loop()
}

//////////////////////////////////////////////
/// main
//////////////////////////////////////////////

func main() {
	d := argsToDuration(args.t, args.z)

	// start timer
	timer := Timer{title: "Timer"}
	timer.start(d)

	// run UI
	timer.countDown()
}
