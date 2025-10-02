// Starting with version 1.23, Go has added support for
// [iterators](https://go.dev/blog/range-functions),
// which lets us range over pretty much anything!

package main

import (
	"fmt"
	"iter"
	"slices"
)

/*
Real-World Analogy
Think of it like a buffet line:
- Iterator: The server behind the counter
- Yield function: Handing you each dish one at a time
- Return value of yield: You saying "yes, more please" (true) or "I'm full, stop" (false)
- Range loop: You, the customer, processing each dish as it's handed to you
The server doesn't prepare all dishes at once and pile them on your plate - they hand them one at a time, and stop when you say you're done.
*/

// Let's look at the `List` type from the
// [previous example](generics) again. In that example
// we had an `AllElements` method that returned a slice
// of all elements in the list. With Go iterators, we
// can do it better - as shown below.
type List[T any] struct {
	head, tail *element[T]
}

type element[T any] struct {
	next *element[T]
	val  T
}

func (lst *List[T]) Push(v T) {
	if lst.tail == nil {
		lst.head = &element[T]{val: v}
		lst.tail = lst.head
	} else {
		lst.tail.next = &element[T]{val: v}
		lst.tail = lst.tail.next
	}
}

// All returns an _iterator_, which in Go is a function
// with a [special signature](https://pkg.go.dev/iter#Seq).
func (lst *List[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		// The iterator function takes another function as
		// a parameter, called `yield` by convention (but
		// the name can be arbitrary). It will call `yield` for
		// every element we want to iterate over, and note `yield`'s
		// return value for a potential early termination.
		for e := lst.head; e != nil; e = e.next {
			if !yield(e.val) {
				return
			}
		}
	}
}

// Iteration doesn't require an underlying data structure,
// and doesn't even have to be finite! Here's a function
// returning an iterator over Fibonacci numbers: it keeps
// running as long as `yield` keeps returning `true`.
func genFib() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 1, 1

		for {
			if !yield(a) {
				return // Early termination
			}
			a, b = b, a+b
		}
	}
}

func main() {
	lst := List[int]{}
	lst.Push(10)
	lst.Push(13)
	lst.Push(23)

	// Since `List.All` returns an iterator, we can use it
	// in a regular `range` loop.
	// Generate values on-demand (memory efficient)
	for e := range lst.All() {
		fmt.Println(e)
	}
	/*
		// Inside lst.All() iterator
		func(yield func(T) bool) {
			for e := lst.head; e != nil; e = e.next {

				if !yield(e.val) {  // ← Calls yield
					//      ↓
					//      └─── Execution jumps here ───┐
					return                            //  │
				}                                     //  │
			}                                         //  │
		}                                             //  │
													  //  │
		// Range loop                                 //  │
		for e := range lst.All() {  // ←───────────────┘
			fmt.Println(e)  // This code runs when yield is called
			// After this executes, control returns to the iterator
		}

		// Think of yield as a remote control for the range loop body:
		yield(value) {
			// 1. Set loop variable to value
			// 2. Press "play" on the range body
			// 3. Wait for range body to finish
			// 4. Check if loop wants to continue (break? or continue?)
			// 5. Return true (if continue) / false (if break) to iterator
		}
	*/

	// Packages like [slices](https://pkg.go.dev/slices) have
	// a number of useful functions to work with iterators.
	// For example, `Collect` takes any iterator and collects
	// all its values into a slice.
	all := slices.Collect(lst.All())
	fmt.Println("all:", all)

	for n := range genFib() {

		// Once the loop hits `break` or an early return, the `yield` function
		// passed to the iterator will return `false`.
		if n >= 10 {
			break
		}
		fmt.Println(n)
	}
}

/*
The Japanese sushi buffet with the conveyor belt (kaiten-zushi/回転寿司) is a perfect analogy for Go iterators:

Kaiten-Zushi vs Go Iterators
Traditional Buffet (Old Way - AllElements())
// Load everything into memory at once
func (lst *List[T]) AllElements() []T {
    var all []T
    for e := lst.head; e != nil; e = e.next {
        all = append(all, e.val)  // Pile everything on one big plate
    }
    return all
}
- Chef prepares ALL sushi at once
- Everything sits on your table
- Takes up lots of space
- Some might go to waste if you're full

Kaiten-Zushi (Iterator Way - All())
func (lst *List[T]) All() iter.Seq[T] {
    return func(yield func(T) bool) {
        for e := lst.head; e != nil; e = e.next {
            if !yield(e.val) {  // Offer next plate
                return  // Customer says "I'm full, stop sending!"
            }
        }
    }
}
- Chef sends one plate at a time on the belt
- You pick it up (yield delivers it)
- Eat it (execute range body)
- Empty plate goes back (yield returns true = "send more!")
- If you're full, press the stop button (break = yield returns false)
- Chef stops making sushi immediately

The Flow Visualized
// The Sushi Chef (Iterator)
func() {
    for each_sushi_type {
        if !yield(sushi) {  // Put plate on belt
            return          // Customer pressed stop button!
        }
        // Wait for customer to finish eating...
        // Customer signals "ready for next plate"
    }
}

// You at the Counter (Range Loop)
for sushi := range chef.All() {
    eat(sushi)              // Process current plate
    // Implicitly signal: "Send next plate!" (yield returns true)

    if (feeling_full) {
        break               // Press stop button (yield returns false)
    }
}

Why This Is Better
1. Memory Efficiency
// Traditional: All sushi on table at once (2GB of sushi!)
all := lst.AllElements()  // 2 million plates on your table!

// Kaiten-zushi: One plate at a time (just a few KB)
for sushi := range lst.All() {  // Only current plate on table
    eat(sushi)
}
2. Early Termination
for n := range genFib() {
    if n >= 10 {
        break  // "I'm full, stop the belt!"
    }
    fmt.Println(n)
}
// Chef immediately stops making more sushi
3. Infinite Sequences
func genFib() iter.Seq[int] {
    return func(yield func(int) bool) {
        a, b := 1, 1
        for {  // Infinite sushi generation!
            if !yield(a) {
                return  // Customer stops you
            }
            a, b = b, a+b
        }
    }
}
// Chef can keep making sushi forever,
// but stops when you say stop

The Key Innovation
In a traditional buffet, the chef must prepare everything upfront. In kaiten-zushi:
- On-demand production: Sushi is made as it's consumed
- No waste: Stop when customer is satisfied
- Infinite menu: Can keep making new types forever
- Space efficient: Only one plate on the belt at a time
That's exactly how Go iterators work with yield! 🎯

*/
