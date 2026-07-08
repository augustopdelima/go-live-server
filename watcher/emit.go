package watcher

import (
	"time"
)

func Debounce(input <- chan Event, interval time.Duration) <- chan Event {

	output := make(chan Event)

	go runDebounce(
		input,
		output,
		interval,
	)

	return output
}

func runDebounce(
	input <- chan Event,
	output chan<-Event,
	interval time.Duration,
) {
	defer close(output)

	var (
		timer *time.Timer
		pending Event
	)

	for {
		var timerChannel <- chan time.Time

		if timer != nil {
			timerChannel = timer.C
		}

		select {
			case event, ok := <-input:
				if !ok {
					if timer != nil {
						output <- pending
					}
					return
				}

				pending = event

				timer = resetTimer(
					timer,
					interval,
				)

			case <- timerChannel:
				output <- pending
				timer = nil
		}
	}
}


func resetTimer(
	timer *time.Timer,
	duration time.Duration,
) *time.Timer {
	if timer == nil {
		return time.NewTimer(duration)
	}

	if !timer.Stop() {
		select {
			case <-timer.C:
			default:
		}
	}

	timer.Reset(duration)

	return timer
}
