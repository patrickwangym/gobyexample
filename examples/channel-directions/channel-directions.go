// When using channels as function parameters, you can
// specify if a channel is meant to only send or receive
// values. This specificity increases the type-safety of
// the program.

package main

import "fmt"

// This `ping` function only accepts a channel for sending
// values. It would be a compile-time error to try to
// receive on this channel.
func ping(pings chan<- string, msg string) {
	pings <- msg
}

// Think of the channel as a box [chan] and the arrow shows data movement:
// <-[chan] = data exits the box (you receive it)
// Arrow points OUT OF the channel → Data flows from channel to you → Receive only
// [chan]<- = data enters the box (you send it)
// Arrow points INTO the channel → Data flows from you into channel → Send only
// The `pong` function accepts one channel for receives
// (`pings`) and a second for sends (`pongs`).
func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings // <-pings (receive) matches <-chan (receive-only)
	pongs <- msg   // pongs<- (send) matches chan<- (send-only)
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
