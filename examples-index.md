# Go by Example - Interactive Learning Index

Click on any topic to jump directly to the Go code example. Perfect for the Week 1 learning plan!

**📝 Progress Tracking**: To mark topics as complete, edit this file and change `- [ ]` to `- [x]` in the source code.

## Basic Concepts (Week 1: Days 1-2)

- [x] 1. [Hello World](./examples/hello-world/hello-world.go)
- [x] 2. [Values](./examples/values/values.go)
- [x] 3. [Variables](./examples/variables/variables.go)
- [x] 4. [Constants](./examples/constants/constants.go)
- [x] 5. [For](./examples/for/for.go)

## Control Flow (Week 1: Day 2-3)

- [x] 6. [If/Else](./examples/if-else/if-else.go)
- [x] 7. [Switch](./examples/switch/switch.go)

## Data Structures (Week 1: Day 3)

- [x] 8. [Arrays](./examples/arrays/arrays.go)
- [x] 9. [Slices](./examples/slices/slices.go)
- [x] 10. [Maps](./examples/maps/maps.go)
      Usage:
      [maps in aac-backend](/maps-usage-in-aac-backend.md)

## Functions (Week 1: Day 3-4)

- [x] 11. [Functions](./examples/functions/functions.go)
- [x] 12. [Multiple Return Values](./examples/multiple-return-values/multiple-return-values.go)
- [x] 13. [Variadic Functions](./examples/variadic-functions/variadic-functions.go)
      Usage:
      [variadic functions and sorter design patterns in aac-backend](/variadic-functions-and-sorter-design-patterns.md)
- [x] 14. [Closures](./examples/closures/closures.go)
      Usage:
      [closures in aac-backend](/closures-in-aac-backend.md)
      Some usage can only be fully grasped when later topics are covered. Will need to revisit this later.
- [x] 15. [Recursion](./examples/recursion/recursion.go)

## Range & Iteration (Week 1: Day 4)

- [x] 16. [Range over Built-in Types](./examples/range-over-built-in-types/range-over-built-in-types.go)

## Advanced Concepts (Week 1: Day 4-5)

- [x] 17. [Pointers](./examples/pointers/pointers.go)
  - [pointers in aac-backend](/pointer-usage-in-aac-backend.md)
  - [pointer usage analysis](/pointer-usage-analysis.md)
  - [pointer patterns analysis](/pointer-patterns-analysis.md)
- [x] 18. [Strings and Runes](./examples/strings-and-runes/strings-and-runes.go)
- [x] 19. [Structs](./examples/structs/structs.go)
      Usage:
      [structs in aac-backend](/struct-usage-in-aac-backend.md)
- [x] 20. [Methods](./examples/methods/methods.go)
- [x] 21. [Interfaces](./examples/interfaces/interfaces.go)
      Usage:
      [interfaces and dependency injection in aac-backend](/interfaces-and-dependency-injection-in-aac-backend.md)

## Modern Go Features

- [x] 22. [Enums](./examples/enums/enums.go)
      usage:
      [enum in aac backend](/enums-in-aac-backend.md)
- [x] 23. [Struct Embedding](./examples/struct-embedding/struct-embedding.go)
  - [struct embedding analysis](/struct-embedding-analysis.md) (callback is a bit hard, revisit later)
  - [nullable types analysis](/nullable-types-analysis.md)
- [x] 24. [Generics](./examples/generics/generics.go)
  - [generics with reflection](/generics-with-reflection.md)
- [x] 25. [Range over Iterators](./examples/range-over-iterators/range-over-iterators.go)
  - [range over iterators analysis](/range-over-iterators-analysis.md) Some code could be refactored to use iterators, revisit later. Need to take note the server must be running Go 1.23+.
  - [range over iterators refactoring opportunity](/range-over-iterators-refactoring-opportunity.md) A detailed analysis of potential refactoring opportunities in the AAC backend using iterators.

## Error Handling (Week 1: Day 5)

- [x] 26. [Errors](./examples/errors/errors.go)
- [x] 27. [Custom Errors](./examples/custom-errors/custom-errors.go)
  - [custom error patterns analysis](/custom-error-patterns-analysis.md)

## Concurrency (Week 2+)

- [ ] 28. [Goroutines](./examples/goroutines/goroutines.go)
- [ ] 29. [Channels](./examples/channels/channels.go)
- [ ] 30. [Channel Buffering](./examples/channel-buffering/channel-buffering.go)
- [ ] 31. [Channel Synchronization](./examples/channel-synchronization/channel-synchronization.go)
- [ ] 32. [Channel Directions](./examples/channel-directions/channel-directions.go)
- [ ] 33. [Select](./examples/select/select.go)
- [ ] 34. [Timeouts](./examples/timeouts/timeouts.go)
- [ ] 35. [Non-Blocking Channel Operations](./examples/non-blocking-channel-operations/non-blocking-channel-operations.go)
- [ ] 36. [Closing Channels](./examples/closing-channels/closing-channels.go)
- [ ] 37. [Range over Channels](./examples/range-over-channels/range-over-channels.go)

## Advanced Concurrency

- [ ] 38. [Timers](./examples/timers/timers.go)
- [ ] 39. [Tickers](./examples/tickers/tickers.go)
- [ ] 40. [Worker Pools](./examples/worker-pools/worker-pools.go)
- [ ] 41. [WaitGroups](./examples/waitgroups/waitgroups.go)
- [ ] 42. [Rate Limiting](./examples/rate-limiting/rate-limiting.go)
- [ ] 43. [Atomic Counters](./examples/atomic-counters/atomic-counters.go)
- [ ] 44. [Mutexes](./examples/mutexes/mutexes.go)
- [ ] 45. [Stateful Goroutines](./examples/stateful-goroutines/stateful-goroutines.go)

## Standard Library

- [ ] 46. [Sorting](./examples/sorting/sorting.go)
- [ ] 47. [Sorting by Functions](./examples/sorting-by-functions/sorting-by-functions.go)

## Error Recovery (Week 1: Day 5)

- [ ] 48. [Panic](./examples/panic/panic.go)
- [ ] 49. [Defer](./examples/defer/defer.go)
- [ ] 50. [Recover](./examples/recover/recover.go)

## String Processing

- [ ] 51. [String Functions](./examples/string-functions/string-functions.go)
- [ ] 52. [String Formatting](./examples/string-formatting/string-formatting.go)
- [ ] 53. [Text Templates](./examples/text-templates/text-templates.go)
- [ ] 54. [Regular Expressions](./examples/regular-expressions/regular-expressions.go)

## Data Formats (Week 2+)

- [ ] 55. [JSON](./examples/json/json.go)
- [ ] 56. [XML](./examples/xml/xml.go)

## Time & Date

- [ ] 57. [Time](./examples/time/time.go)
- [ ] 58. [Epoch](./examples/epoch/epoch.go)
- [ ] 59. [Time Formatting / Parsing](./examples/time-formatting-parsing/time-formatting-parsing.go)

## Utilities

- [ ] 60. [Random Numbers](./examples/random-numbers/random-numbers.go)
- [ ] 61. [Number Parsing](./examples/number-parsing/number-parsing.go)
- [ ] 62. [URL Parsing](./examples/url-parsing/url-parsing.go)
- [ ] 63. [SHA256 Hashes](./examples/sha256-hashes/sha256-hashes.go)
- [ ] 64. [Base64 Encoding](./examples/base64-encoding/base64-encoding.go)

## File Operations

- [ ] 65. [Reading Files](./examples/reading-files/reading-files.go)
- [ ] 66. [Writing Files](./examples/writing-files/writing-files.go)
- [ ] 67. [Line Filters](./examples/line-filters/line-filters.go)
- [ ] 68. [File Paths](./examples/file-paths/file-paths.go)
- [ ] 69. [Directories](./examples/directories/directories.go)
- [ ] 70. [Temporary Files and Directories](./examples/temporary-files-and-directories/temporary-files-and-directories.go)
- [ ] 71. [Embed Directive](./examples/embed-directive/embed-directive.go)

## Testing (Important for Week 1: Day 5)

- [ ] 72. [Testing and Benchmarking](./examples/testing-and-benchmarking/main_test.go)

## Command Line

- [ ] 73. [Command-Line Arguments](./examples/command-line-arguments/command-line-arguments.go)
- [ ] 74. [Command-Line Flags](./examples/command-line-flags/command-line-flags.go)
- [ ] 75. [Command-Line Subcommands](./examples/command-line-subcommands/command-line-subcommands.go)
- [ ] 76. [Environment Variables](./examples/environment-variables/environment-variables.go)

## Logging & HTTP (Critical for AAC Backend)

- [ ] 77. [Logging](./examples/logging/logging.go)
- [ ] 78. [HTTP Client](./examples/http-client/http-client.go)
- [ ] 79. [HTTP Server](./examples/http-server/http-server.go) ⭐ **Key for AAC Backend**

## Advanced Topics

- [ ] 80. [Context](./examples/context/context.go) ⭐ **Important for AAC Backend**
- [ ] 81. [Spawning Processes](./examples/spawning-processes/spawning-processes.go)
- [ ] 82. [Exec'ing Processes](./examples/execing-processes/execing-processes.go)
- [ ] 83. [Signals](./examples/signals/signals.go)
- [ ] 84. [Exit](./examples/exit/exit.go)

---

## Week 1 Learning Path Integration

### Day 1 Focus (Items 1-5)

Start with Hello World through For loops to build basic syntax understanding.

### Day 2 Focus (Items 6-10)

Control flow and data structures - foundation for understanding AAC data models.

### Day 3 Focus (Items 11-21)

Functions and object-oriented concepts - critical for understanding AAC handlers and service patterns.

### Day 4 Focus (Items 17-21 + selected advanced)

Pointers, structs, methods, interfaces - essential for GORM models and service containers.

### Day 5 Focus (Items 26-27, 48-50, 72, 77-79)

Error handling, testing, HTTP servers - directly applicable to AAC backend patterns.

## AAC Backend Connection Points

- **HTTP Server** (Item 79): Direct relevance to `aac-backend` Chi router patterns
- **Context** (Item 80): Used extensively in AAC request handling
- **Structs/Methods** (Items 19-20): Foundation for GORM models
- **Interfaces** (Item 21): Key to service container dependency injection
- **JSON** (Item 55): API response handling in AAC
- **Testing** (Item 72): Essential for AAC development workflow

Click any link to jump directly to the code and start learning!
