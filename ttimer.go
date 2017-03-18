package ttimer

import (
	"fmt"
	ui "github.com/gizak/termui"
	"math"
	"time"
)

func main() {

	// init and close
	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()

	// init cell
	cell := ui.NewPar("")
	cell.TextFgColor = ui.ColorDefault
	cell.Border = false
	cell.X = 2
	cell.Y = 1
	cell.Width = 26
	cell.Height = 2

	// calculate duration
	duration := time.Duration(6 * time.Second)
	endTime := time.Now().Add(duration)

	// draw
	banner := "== Time =="
	draw := func(t int) {
		value := "[finished]"
		// handle time subtraction
		now := time.Now()
		if !now.After(endTime) {
			left := endTime.Sub(now)
			floorSeconds := math.Floor(left.Seconds())
			rounded := time.Duration(floorSeconds) * time.Second
			value = fmt.Sprintf("%v", rounded)
		}
		// render
		cell.Text = fmt.Sprintf("%s\n%v",
			banner,
			value)
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
