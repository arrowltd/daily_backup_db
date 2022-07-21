package utils

import (
	"fmt"
	"time"

	"github.com/arrowltd/daily_backup_db/date"
)

func IntervalEverydayAt(label string, h, m int, fn func()) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Error on IntervalEverydayAt | %s : %s", label, err)
		}
	}()
	now := date.Now()
	firstRun := true
	runAt := date.ChangeTo(now, h, m, 0)
	if runAt.Before(now) {
		runAt = runAt.Add(24 * time.Hour)
	}
	for {
		duration := 24 * time.Hour
		if firstRun {
			duration = runAt.Sub(now)
			if duration < 0 {
				duration = 0
			}
			firstRun = false
		} else {
			now = date.Now()
		}
		fmt.Printf("* RUNNING %s in %v at %v\n", label, duration, now.Add(duration))
		<-time.After(duration)
		go fn()
	}
}
