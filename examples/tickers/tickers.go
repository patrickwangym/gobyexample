// [Timers](timers) are for when you want to do
// something once in the future - _tickers_ are for when
// you want to do something repeatedly at regular
// intervals. Here's an example of a ticker that ticks
// periodically until we stop it.

package main

import (
	"fmt"
	"time"
)

func main() {

	// Tickers use a similar mechanism to timers: a
	// channel that is sent values. Here we'll use the
	// `select` builtin on the channel to await the
	// values as they arrive every 500ms.
	ticker := time.NewTicker(500 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Stopped Ticker successfully")
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			}
		}
	}()

	// Tickers can be stopped like timers. Once a ticker
	// is stopped it won't receive any more values on its
	// channel. We'll stop ours after 1600ms.
	time.Sleep(1600 * time.Millisecond)
	// Stop() does NOT close the channel - it just stops sending values to it.
	// If you did a bare <-ticker.C without any other exit mechanism after stopping, it would block forever.
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
	time.Sleep(500 * time.Millisecond) // Wait to ensure goroutine finishes
}
