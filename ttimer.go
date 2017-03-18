package main

import (
	"flag"
	"fmt"
	ui "github.com/gizak/termui"
	"math"
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
		&args.t, "t", "100", "time string")
	flag.StringVar(
		&args.z, "z", "-0800", "timezone string")
	flag.BoolVar(
		&args.N, "N", true, "use notifications")
	flag.Parse()
}

//////////////////////////////////////////////
/// functions
//////////////////////////////////////////////

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
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
	panicErr(err)
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

	// start timer
	timer := Timer{title: "Timer"}
	d := time.Duration(6 * time.Second)
	timer.start(d)

	// run UI
	timer.countDown()
}
