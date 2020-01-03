package orchestrator

import "time"

type clock struct{}

func (c clock) Sleep(d time.Duration) {
	time.Sleep(d)
}
