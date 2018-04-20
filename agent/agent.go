package agent

import (
	. "fmt"
	"github.com/0xAX/notificator"
	ui "github.com/gizak/termui"
	"math"
	"time"
)

//////////////////////////////////////////////
/// util
//////////////////////////////////////////////

func mustBeNil(err error) {
	if err != nil {
		panic(err)
	}
}

//////////////////////////////////////////////
/// Timer
//////////////////////////////////////////////

type Timer struct {
	Title    string
	Debug    bool
	duration time.Duration
	end      time.Time
	left     time.Duration
	status   string
}

func (t *Timer) Start(d time.Duration) {
	t.duration = d
	if t.Title == "" {
		t.Title = Sprintf("%v Timer", d)
	}
	t.end = time.Now().Add(t.duration)
}

func (t *Timer) update() {
	t.status = "Finished"
	now := time.Now()
	if !now.After(t.end) {
		exactLeft := t.end.Sub(now)
		floorSeconds := math.Floor(exactLeft.Seconds())
		t.left = time.Duration(floorSeconds) * time.Second
		t.status = Sprintf("%v", t.left)
		if t.Debug {
			t.status += "\n"
			t.status += Sprintf("\nt.end: %v", t.end)
		}
	}
}

func (t *Timer) CountDown() {
	// init and close
	err := ui.Init()
	mustBeNil(err)
	defer ui.Close()

	// init notificator
	notify := notificator.New(notificator.Options{
		AppName: t.Title,
	})

	// set and execute pre-notify
	seconds := t.duration.Seconds()
	if seconds > 10 {
		go func() {
			almostSec := math.Floor(seconds * .9)
			almostDur := time.Duration(almostSec) * time.Second
			<-time.After(almostDur)
			message := Sprintf("%v left", t.left)
			notify.Push(
				"", message, "", notificator.UR_CRITICAL)
		}()
	}
	// set and execute notify
	go func() {
		<-time.After(t.duration)
		notify.Push(
			"", "Finished", "", notificator.UR_CRITICAL)
	}()

	// init cell
	cell := ui.NewPar("")
	cell.TextFgColor = ui.ColorDefault
	cell.Border = false
	cell.X = 2
	cell.Y = 1
	cell.Width = ui.TermWidth()
	cell.Height = ui.TermHeight()

	// draw
	banner := Sprintf("== %s ==", t.Title)
	draw := func(tick int) {
		t.update()
		// render
		cell.Text = Sprintf("%s\n%v",
			banner,
			t.status)
		ui.Render(cell)
	}

	// handle update
	ms := 50
	timerPath := Sprintf("/timer/%vms", ms)
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
