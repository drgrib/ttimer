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
		&args.z, "z", "America/Los_Angeles", "timezone")
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
		// min
		minStr := clock[len(clock)-2:]
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return 0, 0, errors.New(
				fmt.Sprintf(
					"Couldn't parse minStr %#v", minStr))

		}
		return hour, min, nil
	}
	hour, err := strconv.Atoi(clock)
	if err != nil {
		return 0, 0, errors.New(
			fmt.Sprintf(
				"Couldn't parse as hour %#v", clock))
	}
	return hour, 0, nil
}

func parseTime(t string, z string) (time.Duration, string, error) {
	// parameterized location due to not all platforms supporting local detection
	loc, err := time.LoadLocation(z)
	zero := time.Duration(0)
	if err != nil {
		return zero, "", err
	}
	now := time.Now().In(loc)
	// track period
	pattern := `(\d+)(a|p)?`
	r := regexp.MustCompile(pattern)
	m := r.FindStringSubmatch(t)
	clock := m[1]
	period := m[2]
	// handle minute case
	if period == "" && len(clock) <= 2 {
		return zero, "", errors.New(
			fmt.Sprintf("No period, assuming minutes, not Time: %#v", clock))
	}
	// handle clock
	hour, min, err := parseClock(clock)
	if err != nil {
		return zero, "", err
	}
	// estimate endTime
	endTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		hour, min,
		0, 0,
		loc)
	// increment by 12 hours until after now
	for endTime.Before(now) {
		endTime = endTime.Add(12 * time.Hour)
	}
	// final increment if wrong period
	if period == "a" && endTime.Hour() >= 12 {
		endTime = endTime.Add(12 * time.Hour)
	}
	if period == "p" && endTime.Hour() < 12 {
		endTime = endTime.Add(12 * time.Hour)
	}
	// calculate the duration
	d := endTime.Sub(now)
	// format the title
	layout := "304pm"
	if min == 0 {
		// further truncate
		layout = "3pm"
	}
	formatted := endTime.Format(layout)
	// truncate period
	formatted = formatted[:len(formatted)-1]
	title := fmt.Sprintf("%v Timer", formatted)
	return d, title, nil
}

func parseArgs(t string, z string) (time.Duration, string) {
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
		d, title, err := parseTime(t, z)
		if err == nil {
			return d, title
		}
		// if not time, parse as minute
		minutes, err := strconv.Atoi(t)
		if err != nil {
			break
		}
		d = time.Duration(minutes) * time.Minute
		title = fmt.Sprintf("%vm Timer", t)
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

	// start timer
	timer := Timer{title: title}
	timer.start(d)

	// run UI
	timer.countDown()
}
