package crontz

import (
	"testing"
	"time"
)

func TestCrontz(t *testing.T) {
	for _, tt := range []struct {
		name      string
		expr      string
		okayTimes []time.Time
		badTimes  []time.Time
	}{{
		name:      "every minute",
		expr:      "* * * * *",
		okayTimes: []time.Time{time.Now()},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			c, err := ParseCron(tt.expr)
			if err != nil {
				t.Fatal(err)
			}
			for _, x := range tt.okayTimes {
				if !c.Matches(x) {
					t.Errorf("expected %s to match %s", tt.expr, x)
				}
			}
			for _, x := range tt.badTimes {
				if c.Matches(x) {
					t.Errorf("expected %s to not match %s", tt.expr, x)
				}
			}
		})
	}
}
