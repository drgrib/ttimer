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
/// AfterWallClock
//////////////////////////////////////////////

func AfterWallClock(d time.Duration) <-chan time.Time {
	c := make(chan time.Time, 1)
	go func() {
		end := time.Now().Add(d)
		// clear monotonic clock
		end = end.Round(0)
		for !time.Now().After(end) {
			time.Sleep(100 * time.Millisecond)
		}
		c <- time.Now()
	}()
	return c
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
	// strip monotonic time to account for system changes
	t.end = time.Now().Add(t.duration).Round(0)
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
			t.status += Sprintf("\nnow: %v", now)
			t.status += Sprintf("\nexactLeft: %v", exactLeft)
			t.status += Sprintf("\nt.end: %v", t.end)
			t.status += Sprintf("\nt.end.Sub(now): %v", t.end.Sub(now))
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
			<-AfterWallClock(almostDur)
			message := Sprintf("%v left", t.left)
			notify.Push(
				"", message, "", notificator.UR_CRITICAL)
		}()
	}
	// set and execute notify
	go func() {
		<-AfterWallClock(t.duration)
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
	ui.Handle("/sys/kbd/C-c", func(ui.Event) {
		ui.StopLoop()
	})

	// start loop
	ui.Loop()
}
