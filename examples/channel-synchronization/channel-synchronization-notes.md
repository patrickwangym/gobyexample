# Channel Synchronization in Go

## The Problem We're Solving

When you launch a goroutine with `go functionName()`, it runs concurrently in the background. **The main goroutine doesn't automatically wait for it to finish.** This creates several issues:

### 1. Premature Program Exit

```go
func main() {
    go worker()  // Runs in background
    // main exits immediately - worker might not even start!
}
```

### 2. Race Conditions

Without synchronization, you can't safely:

- Know when a goroutine has completed its work
- Access shared data that a goroutine is modifying
- Ensure operations happen in the correct order

### 3. Resource Dependencies

You might need to:

- Wait for a database operation to complete before closing the connection
- Ensure file writes finish before reading the file
- Complete API calls before processing results

## How Channels Solve This

Channels provide a **synchronization point** where goroutines can coordinate:

```go
done := make(chan bool, 1)
go worker(done)
<-done  // BLOCKS here until worker sends a value
```

The main goroutine **blocks** at `<-done` and can't proceed until the worker goroutine sends a value.

## Buffered vs Unbuffered Channels

### Unbuffered Channel (`make(chan bool)`)

```go
done := make(chan bool)  // No buffer
```

**Behavior:**

- `done <- true` in the worker **blocks** until main executes `<-done`
- Both goroutines must "meet" at the channel at the same time
- It's a **rendezvous** - sender waits for receiver, receiver waits for sender
- **Strictly synchronous handoff**

### Buffered Channel (`make(chan bool, 1)`)

```go
done := make(chan bool, 1)  // Buffer of 1
```

**Behavior:**

- `done <- true` in the worker **does NOT block** (buffer has space)
- Worker can send and **immediately continue/exit** without waiting
- Main still blocks at `<-done` until a value is available
- **Asynchronous send, synchronous receive**

### Key Difference

**Unbuffered = tight coupling:** Worker must wait at send until main receives

**Buffered = looser coupling:** Worker can send and finish; main receives later

## Comparison Table

| Aspect              | Unbuffered           | Buffered (size 1)                    |
| ------------------- | -------------------- | ------------------------------------ |
| **Send operation**  | Blocks until receive | Only blocks if buffer full           |
| **Worker behavior** | Must wait for main   | Can send and exit immediately        |
| **Synchronization** | Strict rendezvous    | Decoupled timing                     |
| **For simple sync** | Works perfectly      | Works perfectly (slight flexibility) |

## Why Use Buffered in This Example?

The buffered channel is a defensive choice that ensures the worker can always send its completion signal without blocking, even if main hasn't reached the receive statement yet.

```go
func worker(done chan bool) {
    fmt.Println("done")
    done <- true  // With buffer: sends immediately, worker can exit
                  // Without buffer: must wait for main to receive
}
```

## Common Pitfall: Multiple Sends with Unbuffered Channel

```go
// Unbuffered - Partial Deadlock/Goroutine Leak
results := make(chan int)  // Unbuffered
go func() {
    results <- 1  // Blocks until main receives
    results <- 2  // Blocks forever - DEADLOCK here!
}()
time.Sleep(2 * time.Second)
fmt.Println(<-results)  // Prints 1 (unblocks first send)
// Goroutine hangs at second send - goroutine leak!
```

**What happens:**

1. ✓ Goroutine sends `1`, blocks
2. ✓ Main wakes up after sleep, receives `1`, prints it
3. ✗ Goroutine tries to send `2`, blocks forever (no corresponding receive)
4. ✗ **Goroutine leak** - the goroutine is stuck waiting

**Key Rule:** Unbuffered channels require a 1:1 pairing of sends and receives. Each send must have a matching receive happening concurrently.

## Real-World Example

```go
func processOrders(orders []Order) {
    done := make(chan bool)

    go func() {
        // Process thousands of orders in background
        for _, order := range orders {
            saveToDatabase(order)
        }
        done <- true  // Signal completion
    }()

    <-done  // Wait here - don't close DB connection yet!
    closeDatabase()  // Safe to close now
}
```

**Without the channel:** The database might close while orders are still being saved, causing errors or data loss.

**With the channel:** We guarantee all orders are saved before cleanup happens.

## Summary

- **Channel synchronization** coordinates goroutines to prevent premature exits and race conditions
- **Unbuffered channels** provide strict synchronization (rendezvous)
- **Buffered channels** allow send to complete without immediate receive
- Both work for simple "done" signals; buffered is slightly more flexible
- Always pair sends with receives to avoid deadlocks and goroutine leaks
