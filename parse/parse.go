package parse

import (
	. "fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

//////////////////////////////////////////////
/// parseArgs
//////////////////////////////////////////////

func parseClock(clock string) (int, int, error) {
	if len(clock) >= 3 {
		// hour
		hourStr := clock[:len(clock)-2]
		hour, err := strconv.Atoi(hourStr)
		if err != nil {
			return 0, 0, Errorf(
				"Couldn't parse hourStr %#v", hourStr)
		}
		// min
		minStr := clock[len(clock)-2:]
		min, err := strconv.Atoi(minStr)
		if err != nil {
			return 0, 0, Errorf(
				"Couldn't parse minStr %#v", minStr)

		}
		return hour, min, nil
	}
	hour, err := strconv.Atoi(clock)
	if err != nil {
		return 0, 0, Errorf(
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
	if len(m) < 3 {
		return zero, "", Errorf("Could not parse as Time: %#v", t)
	}
	clock := m[1]
	period := m[2]
	// handle minute case
	if period == "" && len(clock) <= 2 {
		return zero, "", Errorf("No period, assuming minutes, not Time: %#v", clock)
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
	title := Sprintf("%v Timer", formatted)
	return d, title, nil
}

func Args(t string) (time.Duration, string, error) {
	f, err := strconv.ParseFloat(t, 64)
	switch {
	case err == nil:
		floatMinutes := math.Floor(f)
		seconds := int64(math.Floor((f - floatMinutes) * 60))
		minutes := int64(floatMinutes)

		if seconds == 0 && len(t) > 1 {
			// parse as time
			d, title, err := parseTime(t)
			if err == nil {
				return d, title, nil
			}
		}

		d := time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
		title := Sprintf("%vm Timer", f)
		return d, title, nil
	case len(t) == 1:
		// simple minute timer
		minutes, err := strconv.Atoi(t)
		if err != nil {
			return 0, "", err
		}

		d := time.Duration(minutes) * time.Minute
		title := Sprintf("%vm Timer", t)
		return d, title, nil
	default:
		// parse as duration
		d, err := time.ParseDuration(t)
		if err == nil {
			title := Sprintf("%v Timer", t)
			return d, title, nil
		}
		// parse as time
		d, title, err := parseTime(t)
		if err == nil {
			return d, title, nil
		}
		// if not time, parse as minute
		minutes, err := strconv.Atoi(t)
		if err != nil {
			return 0, "", err
		}
		d = time.Duration(minutes) * time.Minute
		title = Sprintf("%vm Timer", t)
		return d, title, nil
	}
}
