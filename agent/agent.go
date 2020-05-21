package agent

import (
	. "fmt"
	"math"
	"time"

	"github.com/0xAX/notificator"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	termX = 1
	termY = 0
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
}

func (t *Timer) update() {
	t.status = "Finished\n\n[r]estart\n[q]uit"
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

	p := widgets.NewParagraph()
	termWidth, termHeight := ui.TerminalDimensions()
	p.SetRect(termX, termY, termWidth, termHeight)
	p.TextStyle.Fg = ui.ColorClear
	p.Border = false

	// draw
	banner := Sprintf("== %s ==", t.Title)
	draw := func(tick int) {
		t.update()
		// render
		p.Text = Sprintf("%s\n%v",
			banner,
			t.status)
		ui.Render(p)
	}

	tickerCount := 1
	draw(tickerCount)
	tickerCount++
	ticker := time.NewTicker(100 * time.Millisecond).C

	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>", "<Escape>":
				return
			case "r":
				if time.Now().After(t.end) {
					t.Start(t.duration)
				}
			case "<Resize>":
				resize := e.Payload.(ui.Resize)
				p.SetRect(termX, termY, resize.Width, resize.Height)
			}
		case <-ticker:
			draw(tickerCount)
			tickerCount++
		}
	}
}
