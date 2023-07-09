package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	lessThan60 = `([0-9]|[0-5][0-9])`
	lessThan31 = `([0-2][0-9]|30|31)`
	lessThan24 = `([0-1][0-9]|2[0-4])`
	lessThan12 = `([0-9]|1[0-2])`
	lessThan7  = `([0-6])`

	commonCronPattern = `^(\*|(?P<single>d)|(?P<range>d-d)|(?P<step>\*/d))$`
)

var (
	reMinute  = regexp.MustCompile(strings.ReplaceAll(commonCronPattern, "d", lessThan60))
	reHour    = regexp.MustCompile(strings.ReplaceAll(commonCronPattern, "d", lessThan24))
	reDay     = regexp.MustCompile(strings.ReplaceAll(commonCronPattern, "d", lessThan31))
	reMonth   = regexp.MustCompile(strings.ReplaceAll(commonCronPattern, "d", lessThan12))
	reWeekday = regexp.MustCompile(strings.ReplaceAll(commonCronPattern, "d", lessThan7))
)

type duration struct {
	regex   *regexp.Regexp
	allowed []bool
}

var durations = map[string]duration{
	"minute":  {reMinute, make([]bool, 60)},
	"hour":    {reHour, make([]bool, 24)},
	"day":     {reDay, make([]bool, 31+1)},   // +1 because days are 1-31 and not 0-30
	"month":   {reMonth, make([]bool, 12+1)}, // same for months
	"weekday": {reWeekday, make([]bool, 7)},
}

func makeDurations(singleCronString string, dur string) error {
	values := strings.Split(singleCronString, ",")

	allowed := durations[dur].allowed
	re := durations[dur].regex
	single, ranged, step := re.SubexpIndex("single"), re.SubexpIndex("range"), re.SubexpIndex("step")

	for _, value := range values {
		matches := re.FindAllStringSubmatch(value, -1)
		if matches == nil {
			return fmt.Errorf("invalid minute pattern: %s", value)
		}

		if value == "*" {
			for i := range allowed {
				allowed[i] = true
			}
			break
		}

		if matches[0][single] != "" {
			val, err := strconv.ParseUint(matches[0][single], 10, 0)
			if err != nil {
				return fmt.Errorf("invalid minute value: %s", value)
			}
			allowed[uint8(val)] = true
		}

		if matches[0][ranged] != "" {
			// _range+1 because the first match is the whole string
			min, err := strconv.ParseUint(matches[0][ranged+1], 10, 0)
			if err != nil {
				return fmt.Errorf("invalid minute range: %s", value)
			}

			max, err := strconv.ParseUint(matches[0][ranged+2], 10, 0)
			if err != nil {
				return fmt.Errorf("invalid minute range: %s", value)
			}

			if min > max {
				return fmt.Errorf("invalid minute range: %s", value)
			}

			for i := min; i <= max; i++ {
				allowed[uint8(i)] = true
			}
		}

		if matches[0][step] != "" {
			// step+1 because the first match is the whole string
			stepInt, err := strconv.ParseUint(matches[0][step+1], 10, 0)
			if err != nil {
				return fmt.Errorf("invalid minute step: %s", value)
			}

			for i := uint8(0); i < uint8(len(allowed)); i += uint8(stepInt) {
				allowed[i] = true
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
		if durations["minute"].allowed[now.Minute()] &&
			durations["hour"].allowed[now.Hour()] &&
			durations["day"].allowed[now.Day()] &&
			durations["month"].allowed[now.Month()] &&
			durations["weekday"].allowed[now.Weekday()-1] {
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
