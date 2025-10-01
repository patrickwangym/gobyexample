# Generics + Reflection Tutorial

This directory contains a comprehensive tutorial on how Go generics work with reflection to solve real production problems in the aac-backend codebase.

---

## Quick Start

```bash
cd learning/gobyexample/examples/generics
go run generics.go
```

**What you'll see**:
- Basic generics with type parameters
- Type-specific functions (verbose approach)
- Generics + reflection (concise approach)
- Production pattern (higher-order functions)
- Side-by-side comparison of all approaches

---

## Learning Path

### Step 1: Run the Example
```bash
go run generics.go
```

Observe the output to see how all three levels work.

---

### Step 2: Read the Concepts
📖 **[CONCEPTS.md](CONCEPTS.md)** - Start here

Learn:
- Why Go needs reflection (type system limitations)
- How reflection works
- Generics + reflection combination
- Design trade-offs (generic vs specialized sorters)
- Common pitfalls

---

### Step 3: Study the Evolution
📖 **[EVOLUTION.md](EVOLUTION.md)** - Read after CONCEPTS.md

Understand:
- Level 1 → Level 2 → Level 3 progression
- Side-by-side code comparisons
- Performance optimizations (20,000x improvement!)
- Production architecture in aac-backend
- When to use each level

---

### Step 4: Explore Production Code

**Files to examine** (in order):

1. **[`slicesorter.go:14`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L14)**
   - See the `SliceSorter[T any]` type definition
   - Study `NewStringSorter`, `NewIntSorter`, etc.

2. **[`slicesorter.go:67`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L67)**
   - See how field extraction is isolated
   - Understand error handling approach

3. **[`collections.go:310`](/aac-backend/internal/collections/collections.go#L310)**
   - See `FilterSlice`, `SorterSlice`, `PaginateSlice`
   - Understand how they compose together

4. **[`collections.go:345`](/aac-backend/internal/collections/collections.go#L345)**
   - See real usage of the sorter pattern
   - Understand the map-of-sorters approach

---

## Challenges

Try these hands-on exercises after reading CONCEPTS.md and EVOLUTION.md:

### Challenge 1: Add New Struct Type

```go
type Book struct {
    Title  string
    Author string
    Pages  int
}

books := []Book{
    {Title: "Go Programming", Author: "John", Pages: 300},
    {Title: "Advanced Go", Author: "Alice", Pages: 450},
}

// TODO: Sort by Title, then by Pages
```

### Challenge 2: Handle Nested Fields

```go
type Address struct {
    City    string
    Country string
}

type Person struct {
    Name    string
    Address Address
}

// TODO: Sort by Address.City using reflection
// Hint: Split fieldName by "." and call FieldByName multiple times
```

### Challenge 3: Implement Nullable Field Sorters

```go
type Person struct {
    Name     string
    Nickname *string  // Might be nil
    Age      *int     // Might be nil
}

// TODO: Handle nil values correctly (e.g., sort them last)
// Hint: Check field.IsNil() before dereferencing
```

### Challenge 4: Build Hybrid Approach

Implement a convenience API with performance options:
- Public `NewSorter` function (convenience + error handling)
- Private type-specific helpers (performance)
- Return error for unsupported types
- Type switch only once at creation

---

## Quick Reference

### Generic Type Parameters

```go
func MyFunc[T any](val T) T          // T can be anything
func MyFunc[T comparable](val T)     // T can use == and !=
func MyFunc[T Stringer](val T)       // T must implement Stringer
```

### Reflection Basics

```go
// Get reflection value
v := reflect.ValueOf(myStruct)

// Access field by name
field := v.FieldByName("Name")

// Check if valid
if !field.IsValid() { /* handle error */ }

// Get actual value
str := field.String()   // for string fields
num := field.Int()      // for int fields
```

### Higher-Order Functions

```go
// Function that returns a function
func MakeAdder(x int) func(int) int {
    return func(y int) int {
        return x + y
    }
}

add5 := MakeAdder(5)
result := add5(3)  // 8
```

---

## Common Pitfalls

### Pitfall 1: Assuming Fields Exist

```go
// ❌ Will panic if field doesn't exist
field := v.FieldByName("NonExistent")
str := field.String()  // PANIC!

// ✅ Check validity first
field := v.FieldByName("NonExistent")
if !field.IsValid() {
    return "", errors.New("field not found")
}
```

### Pitfall 2: Wrong Type Assertion

```go
// ❌ Will panic if field is not a string
str := field.String()  // Field is int → returns garbage

// ✅ Check type first
if field.Kind() != reflect.String {
    return "", errors.New("field is not a string")
}
```

### Pitfall 3: Ignoring Pointers/Interfaces

```go
// ❌ Doesn't handle *string or interface{}
str := field.String()

// ✅ Unwrap first
if field.Kind() == reflect.Ptr || field.Kind() == reflect.Interface {
    field = field.Elem()
}
str := field.String()
```

---

## Impact in aac-backend

**Without this pattern, the codebase would need**:
- ~50 models (Client, Receipt, Activity, etc.)
- ×5 sortable fields per model
- ×2 sort orders (asc, desc)
- = **500 functions** just for sorting! 😱

**With this pattern**:
- 10 generic sorter factories (string, int, date, etc.)
- 3 collection functions (filter, sort, paginate)
- = **13 functions total** ✅

**Code reduction**: 97% fewer lines for collection operations!

---

## Additional Resources

### Official Documentation
- [Go Generics Proposal](https://go.dev/blog/intro-generics)
- [Type Parameters Guide](https://go.dev/blog/deconstructing-type-parameters)
- [Reflection Package Docs](https://pkg.go.dev/reflect)

### In aac-backend Codebase
- `/search-concept generics` - Find all generic usage
- `/search-concept reflection` - Find all reflection usage
- `/explain-flow` - Trace how sorting works end-to-end

---

## Key Takeaways

1. **Generics provide type safety** - Compiler ensures type correctness
2. **Reflection provides flexibility** - Runtime field access by name
3. **Together they're powerful** - Type-safe containers + dynamic operations
4. **Higher-order functions** - Separate configuration from execution
5. **Production adds complexity** - For error handling, performance, flexibility
6. **Know when to use** - Balance simplicity vs capability

---

## Questions?

- **Concepts**: Read [CONCEPTS.md](CONCEPTS.md)
- **Evolution**: Read [EVOLUTION.md](EVOLUTION.md)
- **Specific implementation**: Use `/search-concept` to find examples
- **Architecture**: Use `/explain-flow` to trace execution

**Remember**: The tutorial teaches the pattern. The production code shows battle-tested implementation. Both are valuable at different stages of learning!
