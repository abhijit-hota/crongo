package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Static allocation, baby!
const lessThan60 = `([0-9]|[0-5][0-9])`
const lessThan31 = `([0-2][0-9]|30|31)`
const lessThan24 = `([0-1][0-9]|2[0-4])`
const lessThan12 = `([0-9]|1[0-2])`
const lessThan7 = `([0-6])`

const commonCronPattern = `^(\*|(?P<single>d)|(?P<range>d-d)|(?P<step>\*/d))$`

func makeRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(strings.ReplaceAll(commonCronPattern, "d", pattern))
}

var minuteRegex = makeRegex(lessThan60)
var hourRegex = makeRegex(lessThan24)
var dayRegex = makeRegex(lessThan31)
var monthRegex = makeRegex(lessThan12)
var weekdayRegex = makeRegex(lessThan7)

type thing struct {
	regex  *regexp.Regexp
	values []bool
}

var durations = map[string]thing{
	"minute":  {minuteRegex, make([]bool, 60)},
	"hour":    {hourRegex, make([]bool, 24)},
	"day":     {dayRegex, make([]bool, 31+1)},   // +1 because days are 1-31 and not 0-30
	"month":   {monthRegex, make([]bool, 12+1)}, // same for months
	"weekday": {weekdayRegex, make([]bool, 7)},
}

func parseUint(str string) (uint8, error) {
	u, err := strconv.ParseUint(str, 10, 0)
	return uint8(u), err
}

func makeDurations(singleCronString string, dur string) error {
	values := strings.Split(singleCronString, ",")

	duration := durations[dur]
	re := duration.regex
	single, ranged, step := re.SubexpIndex("single"), re.SubexpIndex("range"), re.SubexpIndex("step")

	for _, value := range values {
		matches := duration.regex.FindAllStringSubmatch(value, -1)
		if matches == nil {
			return fmt.Errorf("invalid minute pattern: %s", value)
		}

		if value == "*" {
			for i := range duration.values {
				duration.values[i] = true
			}
			break
		}

		if matches[0][single] != "" {
			val, err := parseUint(matches[0][single])
			if err != nil {
				return fmt.Errorf("invalid minute value: %s", value)
			}
			duration.values[val] = true
		}

		if matches[0][ranged] != "" {
			// _range+1 because the first match is the whole string
			min, err := parseUint(matches[0][ranged+1])
			if err != nil {
				return fmt.Errorf("invalid minute range: %s", value)
			}

			max, err := parseUint(matches[0][ranged+2])
			if err != nil {
				return fmt.Errorf("invalid minute range: %s", value)
			}

			if min > max {
				return fmt.Errorf("invalid minute range: %s", value)
			}

			for i := min; i <= max; i++ {
				duration.values[i] = true
			}
		}

		if matches[0][step] != "" {
			// step+1 because the first match is the whole string
			stepInt, err := parseUint(matches[0][step+1])
			if err != nil {
				return fmt.Errorf("invalid minute step: %s", value)
			}

			for i := 0; i < len(duration.values); i += int(stepInt) {
				duration.values[i] = true
			}
		}
	}

	return nil
}

var durationOrder = []string{"minute", "hour", "day", "month", "weekday"}

func RunCron(expr string) error {

	fmt.Println("Parsing cron")

	stars := strings.Split(expr, " ")
	if len(stars) != 5 {
		return fmt.Errorf("invalid cron expression: %s", expr)
	}

	for i, star := range stars {
		err := makeDurations(star, durationOrder[i])
		if err != nil {
			return err
		}
	}
	fmt.Println("Waiting to start from next minute!")

	nextMinute := time.Now().Truncate(time.Minute).Add(time.Minute)
	time.Sleep(time.Until(nextMinute))

	fmt.Println("Job started")

	for range time.Tick(time.Minute * 1) {
		now := time.Now()
		if durations["minute"].values[now.Minute()] &&
			durations["hour"].values[now.Hour()] &&
			durations["day"].values[now.Day()] &&
			durations["month"].values[now.Month()] &&
			durations["weekday"].values[now.Weekday()-1] {
			fmt.Println("It's time!", now.String())
		}
	}

	return nil
}

func main() {
	cronExpr := os.Args[1]
	err := RunCron(cronExpr)
	if err != nil {
		panic(err)
	}
}
