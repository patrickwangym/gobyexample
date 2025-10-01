# Generics + Reflection: Core Concepts

**Prerequisites**: None - start here to understand the fundamentals

**Next**: Read [EVOLUTION.md](EVOLUTION.md) to see how this pattern evolves from simple to production-grade

---

## The Core Problem

**Question**: How do you sort different struct types by different fields without writing hundreds of duplicate functions?

**Challenge**: Go's generics **cannot** access struct fields directly because there are no field constraints.

```go
// ❌ This WILL NOT compile
func SortByField[T any](slice []T, fieldName string) {
    sort.Slice(slice, func(i, j int) bool {
        return slice[i].fieldName < slice[j].fieldName
        // ERROR: type parameter T has no field or method fieldName
    })
}
```

## Why This Limitation Exists

Go's type system doesn't support "shape constraints" - you can't say "T must have a field called Name".

The only constraints available are:
- `any` - accepts anything
- `comparable` - can use `==` and `!=`
- Interface constraints - must implement specific methods

**There's no way to say**: "T must have a string field called Name"

---

## The Solution: Generics + Reflection

### What is Reflection?

Reflection lets you inspect and manipulate types at **runtime** instead of **compile-time**.

```go
person := Person{Name: "Alice", Age: 30}

// Compile-time: Access fields directly
name := person.Name  // ✅ Compiler knows Person has Name field

// Runtime: Access fields by string name
v := reflect.ValueOf(person)
nameField := v.FieldByName("Name")  // ✅ Works at runtime
nameStr := nameField.String()       // "Alice"
```

**Key Insight**: Reflection trades compile-time safety for runtime flexibility.

### How to Combine Generics with Reflection

```go
func SortByStringField[T any](slice []T, fieldName string, ascending bool) error {
    sort.Slice(slice, func(i, j int) bool {
        // Step 1: Convert to reflection value
        valueI := reflect.ValueOf(slice[i])
        valueJ := reflect.ValueOf(slice[j])

        // Step 2: Access field by name at runtime
        fieldI := valueI.FieldByName(fieldName)
        fieldJ := valueJ.FieldByName(fieldName)

        // Step 3: Extract actual string value
        strI := fieldI.String()
        strJ := fieldJ.String()

        // Step 4: Compare
        if ascending {
            return strI < strJ
        }
        return strI > strJ
    })
    return nil
}
```

**What Each Part Does**:

1. **Generics `[T any]`**: Type-safe container - compiler ensures all elements have same type
2. **Reflection**: Runtime field access - checks if field exists and extracts value
3. **Together**: Type safety at the edges + flexibility in the middle

---

## Design Trade-offs: Generic vs Specialized Sorters

A critical design decision: Should you create **one generic sorter** that handles all types, or **multiple specialized sorters** for each type?

### Approach 1: Universal Generic Sorter

```go
// One function handles string, int, float, etc.
func NewGenericSorter[T any](fieldName string, ascending bool) Comparator[T] {
    return func(a, b T) bool {
        fieldA := reflect.ValueOf(a).FieldByName(fieldName)
        fieldB := reflect.ValueOf(b).FieldByName(fieldName)

        // Type switch evaluated EVERY comparison
        switch fieldA.Kind() {
        case reflect.String:
            return compareStrings(fieldA, fieldB, ascending)
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            return compareInts(fieldA, fieldB, ascending)
        case reflect.Float32, reflect.Float64:
            return compareFloats(fieldA, fieldB, ascending)
        default:
            return false  // ⚠️ Silent failure!
        }
    }
}
```

#### ✅ Pros
1. **DRY** - One function for all primitive types
2. **Convenience** - Don't need to know field type in advance
3. **Extensibility** - Add new types in one place
4. **User-friendly** - Simple API: `NewGenericSorter[Person]("Name", true)`

#### ❌ Cons
1. **Type Safety Loss** ⚠️
   ```go
   // Both compile, but second fails silently at runtime
   NewGenericSorter[Person]("Name", true)    // ✅ Works
   NewGenericSorter[Person]("Invalid", true) // ❌ Returns false, no error!
   ```

2. **Performance Cost**
   - Type switch evaluated **every comparison** (10,000 times for 1000 items)
   - Specialized functions: 0 switches (type known at creation)

3. **Error Handling Problem**
   ```go
   default:
       return false  // Appears as "equal" - breaks sort ordering!
   ```
   - Invalid field types fail silently
   - Hard to debug (sort just doesn't work correctly)

4. **Contract Violation**
   - Comparator should define total ordering
   - Returning `false` for unsupported types breaks mathematical properties
   - Can cause unpredictable sort behavior

5. **Hidden API Complexity**
   ```go
   // User doesn't know which types are supported
   NewGenericSorter[Person]("ComplexField", true)
   // Does this work? Function signature doesn't tell you!
   ```

---

### Approach 2: Specialized Type Sorters

```go
func NewStringSorter[T any](fieldName string, ascending bool) Comparator[T]
func NewIntSorter[T any](fieldName string, ascending bool) Comparator[T]
func NewFloatSorter[T any](fieldName string, ascending bool) Comparator[T]
```

#### ✅ Pros
1. **Explicit Type Expectations** - Function name documents what it does
   ```go
   sorters := map[string]Comparator[Person]{
       "name": NewStringSorter[Person]("Name", true),  // Clear: expects string
       "age":  NewIntSorter[Person]("Age", true),      // Clear: expects int
   }
   ```

2. **Performance** - No type switches during comparisons (10,000x faster)

3. **Better Error Messages**
   ```go
   // Clear error: "field 'Age' is not a string"
   // vs generic: silent failure
   ```

4. **Fail Fast** - Configuration errors caught at factory creation, not buried in sort logic

5. **Self-Documenting API**
   ```go
   NewStringSorter  // I sort string fields
   NewIntSorter     // I sort int fields
   // Function name = contract
   ```

#### ❌ Cons
1. **Code Duplication** - 3+ similar functions
2. **More Functions to Remember** - Need to know which sorter for which type
3. **Boilerplate** - Each type needs separate factory

---

### Approach 3: Hybrid (Best of Both Worlds)

```go
// Public convenience API (with error handling!)
func NewSorter[T any](fieldName string, ascending bool) (Comparator[T], error) {
    // Validate field exists and determine type ONCE at creation
    var zero T
    field := reflect.ValueOf(zero).FieldByName(fieldName)

    if !field.IsValid() {
        return nil, fmt.Errorf("field %s not found", fieldName)
    }

    kind := field.Kind()

    // Type switch ONCE at creation, not every comparison
    switch kind {
    case reflect.String:
        return newStringComparator[T](fieldName, ascending), nil
    case reflect.Int, reflect.Int8, ...:
        return newIntComparator[T](fieldName, ascending), nil
    case reflect.Float32, reflect.Float64:
        return newFloatComparator[T](fieldName, ascending), nil
    default:
        return nil, fmt.Errorf("unsupported field type: %s", kind)
    }
}

// Also expose specialized versions for performance
func NewStringSorter[T any](fieldName string, ascending bool) Comparator[T]
func NewIntSorter[T any](fieldName string, ascending bool) Comparator[T]
```

**Key Improvements**:
1. ✅ **Returns error** - Caller knows immediately if field type is supported
2. ✅ **Type switch once** - Happens at factory creation, not every comparison
3. ✅ **Fails fast** - Invalid configuration caught before sorting starts
4. ✅ **Performance** - Returned comparator has no switch statements
5. ✅ **Flexibility** - Convenience API + performance API

**Usage**:
```go
// Convenience for prototyping
sorter, err := NewSorter[Person]("Name", true)
if err != nil {
    log.Fatal(err)  // Clear error message
}

// Performance for production hot paths
sorter := NewStringSorter[Person]("Name", true)  // Faster, no error check needed
```

---

## The Trade-off Triangle

This is a classic engineering trade-off:

```
        Abstraction
           /\
          /  \
         /    \
        /      \
       /        \
  Type Safety ---- Performance
```

- **More abstraction** → Less type safety, potentially slower
- **More specialization** → More type safety, faster, but more code

**Choose based on**:
1. Dataset size (1,000 vs 100,000 items)
2. Error handling requirements (fail-safe vs fail-fast)
3. API clarity needs (public library vs internal tool)
4. Performance constraints (interactive UI vs batch processing)

---

## Why aac-backend Chose Specialized Approach

The production codebase uses type-specific factories because:

1. **Configuration Clarity**
   - Field types are known at development time
   - Map configuration is self-documenting
   ```go
   sorters := map[string]SliceSorter[Receipt]{
       "receiptNo":   NewIntSorter("ReceiptNo"),     // Obviously an int
       "clientName":  NewStringSorter("ClientName"), // Obviously a string
       "amount":      NewFloatSorter("Amount"),      // Obviously a float
   }
   ```

2. **Production Requirements**
   - Configuration errors need **clear messages** before runtime
   - Performance matters for **large datasets** (100k+ records)
   - **Self-documenting code** > saving 50 lines
   - Field types are static (not determined by user input)

3. **Maintenance Benefits**
   - Each sorter type can evolve independently
   - Easy to add complex types (nullable, nested fields)
   - Clear responsibility boundaries

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
str := field.String()
```

### Pitfall 2: Wrong Type Assertion

```go
// ❌ Will panic or return garbage if field is not a string
str := field.String()  // Field is int → undefined behavior

// ✅ Check type first
if field.Kind() != reflect.String {
    return "", errors.New("field is not a string")
}
str := field.String()
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

**Example**:
```go
type Person struct {
    Name     string   // Direct string
    Nickname *string  // Pointer to string
    Data     any      // Interface wrapping string
}
```

Production code handles all three correctly!

---

## Performance Considerations

**Q: Isn't reflection slow?**

**A: It depends on context**:

- **Reflection overhead**: ~10-100ns per field access
- **Sorting overhead**: O(n log n) comparisons
- **For API responses**: Sorting 1000 items with reflection adds ~0.1ms total

**When to worry**:
- Hot loops executing millions of times per second
- Real-time systems with microsecond latency requirements

**When NOT to worry** (like aac-backend):
- API response formatting (happens once per request)
- Data export (batch operations)
- Admin dashboard queries (human interaction speed)

**The Trade-off**:
- Eliminate 200+ functions and 10,000+ lines of duplicate code
- Accept 0.1ms overhead on API responses that take 50-200ms anyway
- **Result**: Massive win for maintainability at negligible performance cost

---

## Key Takeaways

1. **Generics alone can't solve everything** - Go lacks field constraints

2. **Reflection fills the gap** - Runtime field access enables dynamic operations

3. **Best pattern**: Generics for type safety + Reflection for flexibility

4. **Design trade-offs matter**: Generic convenience vs specialized reliability - choose based on context

5. **When to use**: When code duplication would be massive and performance cost is acceptable

---

## Next Steps

Now that you understand the core concepts, proceed to [EVOLUTION.md](EVOLUTION.md) to see:
- How this pattern evolves from simple to production-grade
- Side-by-side code comparisons
- Performance optimizations
- Real-world architecture in aac-backend
