package main

import (
	"errors"
	"flag"
	"fmt"
	ui "github.com/gizak/termui"
	"math"
	"regexp"
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
/// parseArgs
//////////////////////////////////////////////

func parseClock(clock string) (int, int, error) {
	if len(clock) >= 3 {
		// hour
		hourStr := clock[:len(clock)-2]
		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			return 0, 0, errors.New(
				fmt.Sprintf(
					"Couldn't parse hourStr %#v", hourStr))
		}
		// minute
		minStr := clock[len(clock)-2:]
		minute, err := strconv.Atoi(minStr)
		if err != nil {
			return 0, 0, errors.New(
				fmt.Sprintf(
					"Couldn't parse minStr %#v", minStr))

		}
		return hour, minute, nil
	}
	hour, err := strconv.Atoi(clock)
	if err != nil {
		return 0, 0, errors.New(
			fmt.Sprintf(
				"Couldn't parse as hour %#v", clock))
	}
	return hour, 0, nil
}

func parseAsTime(t, z string) (time.Time, error) {
	failTime := time.Now()
	failErr := errors.New(fmt.Sprintf("Couldn't parse as time %#v", t))
	// track period
	pattern := `(\d+)(a|p)?`
	r := regexp.MustCompile(pattern)
	m := r.FindStringSubmatch(t)
	period := m[2]
	clock := m[1]
	hour, minute, err := parseClock(clock)
	fmt.Println(hour, minute)
	if err != nil {
		return failTime, failErr
	}
	if period == "" {
		// try no period string
	} else {
		// try period string
	}
	// fail whale
	return failTime, failErr
}

func parseArgs(t, z string) (time.Duration, string) {
	switch {
	case len(t) == 1:
		// simple minute timer
		minutes, err := strconv.Atoi(t)
		if err != nil {
			break
		}
		d := time.Duration(minutes) * time.Minute
		title := fmt.Sprintf("%vm Timer", t)
		return d, title
	default:
		// parse as duration
		d, err := time.ParseDuration(t)
		if err == nil {
			title := fmt.Sprintf("%v Timer", t)
			return d, title
		}
		// parse as time
		// timeVal, err := parseAsTime(t, z)
		_, err = parseAsTime(t, z)
		if err == nil {
			return d, "Made It"
		}
		// if not time, parse as minute
		minutes, err := strconv.Atoi(t)
		if err != nil {
			break
		}
		d = time.Duration(minutes) * time.Minute
		title := fmt.Sprintf("%vm Timer", t)
		return d, title
	}
	fmt.Printf(
		"%#v couldn't be parsed, starting 1m timer\n", t)
	d := time.Duration(1 * time.Minute)
	title := "1m Timer"
	return d, title

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
	// parse
	d, title := parseArgs(args.t, args.z)
	return
	// start timer
	timer := Timer{title: title}
	timer.start(d)

	// run UI
	timer.countDown()
}
