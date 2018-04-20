package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/drgrib/ttimer/agent"
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
/// parseArgs
//////////////////////////////////////////////

func parseClock(clock string) (int, int, error) {
	if len(clock) >= 3 {
		// hour
		hourStr := clock[:len(clock)-2]
		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			return 0, 0, fmt.Errorf(
				"Couldn't parse hourStr %#v", hourStr)
		}
		// min
		minStr := clock[len(clock)-2:]
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return 0, 0, fmt.Errorf(
				"Couldn't parse minStr %#v", minStr)

		}
		return hour, min, nil
	}
	hour, err := strconv.Atoi(clock)
	if err != nil {
		return 0, 0, fmt.Errorf(
			"Couldn't parse as hour %#v", clock)
	}
	return hour, 0, nil
}

func parseTime(t string) (time.Duration, string, error) {
	// parameterized location due to not all platforms supporting local detection
	zero := time.Duration(0)
	now := time.Now()
	// track period
	pattern := `(\d+)(a|p)?`
	r := regexp.MustCompile(pattern)
	m := r.FindStringSubmatch(t)
	clock := m[1]
	period := m[2]
	// handle minute case
	if period == "" && len(clock) <= 2 {
		return zero, "", fmt.Errorf("No period, assuming minutes, not Time: %#v", clock)
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
		0, 0, now.Location())
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

func parseArgs(t string) (time.Duration, string) {
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
		d, title, err := parseTime(t)
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
/// main
//////////////////////////////////////////////

func main() {
	// parse
	d, title := parseArgs(args.t)

	// start timer
	t := agent.Timer{Title: title}
	t.Start(d)

	// run UI
	t.CountDown()
}
