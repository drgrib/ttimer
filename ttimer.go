package main

import (
	"fmt"
	ui "github.com/gizak/termui"
	"math"
	"time"
)

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

//////////////////////////////////////////////
/// main
//////////////////////////////////////////////

func main() {

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
	cell.Width = 26
	cell.Height = 2

	// start timer
	timer := Timer{}
	d := time.Duration(6 * time.Second)
	timer.start(d)

	// draw
	banner := "== Time =="
	draw := func(t int) {
		timer.update()
		// render
		cell.Text = fmt.Sprintf("%s\n%v",
			banner,
			timer.status)
		ui.Render(cell)
	}

	// handle update
	ms := 50
	timerPath := fmt.Sprintf("/timer/%vms", ms)
	ui.Merge("timer", ui.NewTimerCh(
		time.Duration(ms)*time.Millisecond))
	ui.Handle(timerPath, func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		draw(int(t.Count))
	})

	// handle quit
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	// start loop
	ui.Loop()
}
