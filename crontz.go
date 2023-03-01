package crontz

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Crontab is a simplified implementation that only supports matching against
// times.
type Crontab [5]uint64

// bounds for minute, hour, day, month and weekday.
var (
	loBounds = [5]int{0, 0, 1, 1, 0}
	hiBounds = [5]int{59, 23, 31, 12, 6}
)

// ParseCron to a set of accepted time values, only parses numeric values.
func ParseCron(exp string) (ct Crontab, err error) {
	fields := strings.Split(exp, " ")
	if l := len(fields); l != 5 {
		return ct, fmt.Errorf("invalid cron %s length %d", exp, l)
	}

	for i, field := range fields {
		for _, part := range strings.Split(field, ",") {
			loN, hiN, stepN := loBounds[i], hiBounds[i], 1

			switch step := strings.Split(part, "/"); {
			case len(step) == 2:
				stepN, err = strconv.Atoi(step[1])
				if err != nil {
					return ct, fmt.Errorf("invalid cron %s part %s", exp, part)
				}
				part = step[0]
			case len(step) > 2:
				return ct, fmt.Errorf("invalid cron %s part %s", exp, part)
			}

			switch split := strings.Split(part, "-"); {
			case len(split) == 2:
				loN, err = strconv.Atoi(split[0])
				if err != nil {
					return ct, fmt.Errorf("invalid cron %s err %v", exp, err)
				}
				hiN, err = strconv.Atoi(split[1])
				if err != nil {
					return ct, fmt.Errorf("invalid cron %s err %v", exp, err)
				}
			case len(split) == 1:
				if part != "*" {
					loN, err = strconv.Atoi(part)
					if err != nil {
						return ct, fmt.Errorf("invalid cron %s err %v", exp, err)
					}
					hiN = loN
				}
			default:
				return ct, fmt.Errorf("invalid cron %s part %s", exp, part)
			}

			if loN < loBounds[i] || hiN > hiBounds[i] {
				return ct, fmt.Errorf("invalid cron %s part %s", exp, part)
			}
			for n := loN; n <= hiN; n += stepN {
				ct[i] |= 1 << uint(n)
			}
		}
	}
	return ct, nil
}

// Matches checks the time is within the cronjob time. Matches checks for
// timezone jumping and will target the upcomming timezone.
func (ct Crontab) Matches(t time.Time) bool {
	if t.Location() != time.UTC {

		current, after := t.Hour(), t.Add(1*time.Hour).Hour()
		if after < current {
			after += 24
		}

		switch after - current {
		case 0: // clock duplication
			return false
		case 2: // clock skip
			if ct.Matches(time.Date(
				t.Year(),
				t.Month(),
				t.Day(),
				t.Hour()+1,
				t.Minute(),
				0, 0,
				time.UTC,
			)) {
				return true
			}
		}
	}
	return ct[0]&(1<<uint(t.Minute())) != 0 && // 0-59
		ct[1]&(1<<uint(t.Hour())) != 0 && // 0-23
		ct[2]&(1<<uint(t.Day())) != 0 && // 1-31
		ct[3]&(1<<uint(t.Month())) != 0 && // 1-12
		ct[4]&(1<<uint(t.Weekday())) != 0 // 0-6
}
