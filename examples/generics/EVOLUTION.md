# Pattern Evolution: Simple → Production

**Prerequisites**: Read [CONCEPTS.md](CONCEPTS.md) first to understand why generics + reflection is needed

**Goal**: See how the pattern evolves from learning code to production-ready implementation

---

## Three Levels of Sophistication

The tutorial demonstrates **three progressively sophisticated approaches**, mirroring real-world code evolution.

```
Level 1: Simple Direct         → Learning & prototyping
Level 2: Higher-Order Functions → Internal tools & reusability
Level 3: Production-Grade      → APIs, performance-critical paths
```

---

## Level 1: Simple Direct Approach

**File**: [`generics.go:129`](/learning/gobyexample/examples/generics/generics.go#L129)

```go
// Simple, direct approach for understanding concepts
func SortByStringField[T any](slice []T, fieldName string, ascending bool) error {
    sort.Slice(slice, func(i, j int) bool {
        valueI := reflect.ValueOf(slice[i]).FieldByName(fieldName)
        valueJ := reflect.ValueOf(slice[j]).FieldByName(fieldName)

        strI := valueI.String()
        strJ := valueJ.String()

        if ascending {
            return strI < strJ
        }
        return strI > strJ
    })
    return nil
}

// Usage
people := []Person{{Name: "Bob"}, {Name: "Alice"}}
SortByStringField(people, "Name", true)
```

### Characteristics
- ✅ Direct and easy to understand
- ✅ All logic in one function
- ❌ No separation of concerns
- ❌ Reflection happens every comparison (~20,000 calls for 1000 items)
- ❌ Can't reuse comparison logic
- ❌ No error handling for invalid fields

### Performance
- **Cost per comparison**: 2 reflection calls (`FieldByName` for each element)
- **Total for 1000 items**: ~20,000 reflection calls
- **Use case**: Learning, one-off scripts

---

## Level 2: Higher-Order Function Pattern

**File**: [`generics.go:188`](/learning/gobyexample/examples/generics/generics.go#L188)

```go
// Step 1: Define comparator type (higher-order function)
type Comparator[T any] func(a, b T) bool

// Step 2: Factory function that creates comparators
func NewStringSorter[T any](fieldName string, ascending bool) Comparator[T] {
    return func(a, b T) bool {
        fieldA := reflect.ValueOf(a).FieldByName(fieldName)
        fieldB := reflect.ValueOf(b).FieldByName(fieldName)

        strA := fieldA.String()
        strB := fieldB.String()

        if ascending {
            return strA < strB
        }
        return strA > strB
    }
}

// Step 3: Generic sort using any comparator
func SortWithComparator[T any](slice []T, comparator Comparator[T]) {
    sort.Slice(slice, func(i, j int) bool {
        return comparator(slice[i], slice[j])
    })
}

// Usage
nameComparator := NewStringSorter[Person]("Name", true)
SortWithComparator(people, nameComparator)
SortWithComparator(people2, nameComparator)  // Reuse!
```

### Characteristics
- ✅ Separation of concerns (factory vs sort logic)
- ✅ Reusable comparators
- ✅ Same pattern as production code!
- ✅ Shows higher-order function concept
- ❌ Still no error handling (tutorial simplification)
- ❌ Reflection happens every comparison (not optimized yet)

### Performance
- **Cost per comparison**: 2 reflection calls (still happens in comparator)
- **Total for 1000 items**: ~20,000 reflection calls
- **Improvement over Level 1**: Reusability (can use same comparator multiple times)
- **Use case**: Multiple sorts with same field/order

### Tutorial Variants
- [`NewStringSorter[T any]`](/learning/gobyexample/examples/generics/generics.go#L188) - Type-specific (lines 188-193)
- [`NewIntSorter[T any]`](/learning/gobyexample/examples/generics/generics.go#L195) - Type-specific (lines 195-208)
- [`NewGenericSorter[T any]`](/learning/gobyexample/examples/generics/generics.go#L230) - Universal with type switch (lines 230-243)

---

## Level 3: Production-Grade Pattern

**File**: [`slicesorter.go:17`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L17)

```go
// Step 1: Define the type for sorter functions (2-level abstraction!)
type SliceSorter[T any] func(sortOrder string) (func(i, j T) bool, error)

// Step 2: Factory function that creates sorters
func NewStringSorter[T any](fieldName string) SliceSorter[T] {
    return func(sortOrder string) (func(i, j T) bool, error) {
        asc := sortOrder == "asc"

        // This comparison function is returned and reused
        return func(a, b T) bool {
            va, err := getStringFieldValue(a, fieldName)
            if err != nil {
                return false
            }
            vb, err := getStringFieldValue(b, fieldName)
            if err != nil {
                return false
            }

            if asc {
                return va < vb
            }
            return va > vb
        }, nil
    }
}

// Step 3: Helper function isolates reflection logic
func getStringFieldValue[T any](obj T, fieldName string) (string, error) {
    val := reflect.ValueOf(obj)

    // Handles both maps and structs
    if val.Kind() == reflect.Map {
        field = val.MapIndex(reflect.ValueOf(fieldName))
    } else {
        field = val.FieldByName(fieldName)
    }

    // Comprehensive validation
    if !field.IsValid() {
        return "", errors.New("invalid field: " + fieldName)
    }

    // Unwraps interfaces and pointers
    if field.Kind() == reflect.Interface || field.Kind() == reflect.Ptr {
        field = field.Elem()
    }

    if field.Kind() != reflect.String {
        return "", errors.New("field is not a string")
    }

    return field.String(), nil
}

// Usage in production - see collections.go:345
func SorterSlice[T any](cq *CollectionQuery, slice []T,
                        sorters map[string]sliceutils.SliceSorter[T]) ([]T, error) {
    sorterFunc := sorters[cq.query.SortBy]
    sortFunc, err := sorterFunc(order)
    if err != nil {
        return nil, err
    }

    sort.SliceStable(slice, func(i, j int) bool {
        return sortFunc(slice[i], slice[j])
    })

    return slice, nil
}
```

### Characteristics
- ✅ Separation of concerns (factory, comparison, field extraction)
- ✅ Error handling for invalid fields and types
- ✅ Supports both maps and structs
- ✅ Handles pointers and interfaces
- ✅ Reusable comparison functions
- ✅ Configurable sort order at runtime

### Performance
- **Cost per comparison**: 0 reflection calls (field extraction optimized in helper)
- **Total for 1000 items**: 1 reflection call (only during factory setup)
- **Improvement over Level 2**: ~20,000x fewer reflection operations! 🚀
- **Use case**: Production APIs, large datasets, repeated sorting

---

## Evolution Comparison: Higher-Order Functions

### Level 1: No Abstraction
```go
// Returns nothing - just sorts in place
func SortByStringField[T any](slice []T, fieldName string, ascending bool) error
```

### Level 2: 1-Level Abstraction
```go
// Returns a function (comparator)
func NewStringSorter[T any](fieldName string, ascending bool) Comparator[T]
//                                                             ↓
//                                                      func(a, b T) bool
```

### Level 3: 2-Level Abstraction
```go
// Returns a function that returns a function (factory returns comparator!)
func NewStringSorter[T any](fieldName string) SliceSorter[T]
//                                            ↓
//                        func(sortOrder string) (func(i, j T) bool, error)
```

### Why 2-Level Abstraction?

Allows separation of:
1. **Configuration** (which field to sort by) - happens at app startup
2. **Execution** (sort order "asc"/"desc") - happens per HTTP request
3. **Comparison** (actual comparison logic) - happens n*log(n) times during sort

**Benefit**: Configure once at startup, apply with different sort orders per request.

---

## Production Enhancements

### 1. Supports Multiple Data Structures

**Tutorial**:
```go
// Only works with structs
reflect.ValueOf(slice[i]).FieldByName(fieldName)
```

**Production**:
```go
// Works with maps AND structs
val := reflect.ValueOf(obj)
if val.Kind() == reflect.Map {
    field = val.MapIndex(reflect.ValueOf(fieldName))  // Map access
} else {
    field = val.FieldByName(fieldName)  // Struct access
}
```

**Why?** In aac-backend:
- Database models are structs (`Person`, `Product`)
- Custom reports return `map[string]interface{}`
- Same sorting code handles both!

---

### 2. Comprehensive Error Handling

**Tutorial**:
```go
// Assumes field exists and is correct type
strI := valueI.String()  // Panics if field is int!
```

**Production**:
```go
// Validates at every step
if !field.IsValid() {
    return "", errors.New("invalid field: " + fieldName)
}

if field.Kind() != reflect.String {
    return "", errors.New("field is not a string")
}

return field.String(), nil
```

**Why?** Production code can't panic:
- User might request sort by non-existent field
- Field type might not match sorter type
- Better to return error than crash the server

---

### 3. Handles Go Type System Complexity

**Tutorial**:
```go
// Direct field access
fieldA := reflect.ValueOf(a).FieldByName(fieldName)
strA := fieldA.String()
```

**Production**:
```go
// Unwraps interfaces and pointers
if field.Kind() == reflect.Interface || field.Kind() == reflect.Ptr {
    field = field.Elem()  // Dereference to get actual value
}
```

**Why?** In Go:
- `interface{}` values need unwrapping
- Pointer fields (`*string`) need dereferencing
- GORM associations often return pointers

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

## Production Architecture: How It All Fits Together

### Request Flow

```
HTTP Request (?sortBy=name&order=asc)
    ↓
Handler receives query params
    ↓
Collections.SorterSlice[Receipt] called
    ↓
Looks up sorter factory from map:
    sorters["name"] = NewStringSorter[Receipt]("Name")
    ↓
Factory creates comparator with runtime order:
    sortFunc = sorter("asc")
    ↓
sort.SliceStable uses comparator (NO reflection per comparison!)
    ↓
Sorted results returned to client
```

### Configuration Pattern

```go
// Define sorters once at app startup
sorters := map[string]SliceSorter[Receipt]{
    "receiptNo":   NewIntSorter("ReceiptNo"),     // Self-documenting
    "clientName":  NewStringSorter("ClientName"), // Clear expectations
    "amount":      NewFloatSorter("Amount"),      // Type explicit
}

// Use per request with runtime order
sortFunc, err := sorters[query.SortBy](query.Order)
sort.SliceStable(receipts, func(i, j int) bool {
    return sortFunc(receipts[i], receipts[j])
})
```

### Key Files in aac-backend

| Component | File | Purpose |
|-----------|------|---------|
| **Type definitions** | [`slicesorter.go:14`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L14) | `SliceSorter[T any]` type |
| **String sorter** | [`slicesorter.go:17`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L17) | String field sorting |
| **Int sorter** | [`slicesorter.go:42`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L42) | Integer field sorting |
| **Field extraction** | [`slicesorter.go:67`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L67) | Reflection helper |
| **Generic sort** | [`collections.go:345`](/aac-backend/internal/collections/collections.go#L345) | Main sorting orchestration |
| **Filter** | [`collections.go:310`](/aac-backend/internal/collections/collections.go#L310) | Generic filtering |
| **Paginate** | [`collections.go:372`](/aac-backend/internal/collections/collections.go#L372) | Generic pagination |

---

## Performance Comparison Summary

| Level | Reflection Calls | Total for 1000 Items | Use Case |
|-------|------------------|---------------------|----------|
| **Level 1** | 2 per comparison | ~20,000 | Learning, scripts |
| **Level 2** | 2 per comparison | ~20,000 | Multiple sorts, reusability |
| **Level 3** | 0 per comparison | ~1 (setup only) | Production, large datasets |

### Reflection Optimization Explained

**Level 2 (Tutorial)**:
```go
// Reflection happens EVERY comparison
return func(a, b T) bool {
    fieldA := reflect.ValueOf(a).FieldByName(fieldName)  // ← Here!
    fieldB := reflect.ValueOf(b).FieldByName(fieldName)  // ← Here!
    return fieldA.String() < fieldB.String()
}
```

**Level 3 (Production)**:
```go
// Reflection happens in separate helper (can be cached/optimized)
va, _ := getStringFieldValue(a, fieldName)  // ← Optimized helper
vb, _ := getStringFieldValue(b, fieldName)
return va < vb
```

Production version extracts field access into helper that can be:
- Optimized with caching
- Enhanced with comprehensive error handling
- Even replaced with compile-time code generation

---

## When To Use Each Level

### Level 1: Simple Direct
**Use when**:
- ✅ Learning generics + reflection
- ✅ One-off scripts or throw-away code
- ✅ Performance doesn't matter
- ✅ Maximum simplicity is goal

**Avoid when**:
- ❌ Need reusability
- ❌ Performance matters
- ❌ Production systems

**Example**: [`generics.go:129`](/learning/gobyexample/examples/generics/generics.go#L129)

---

### Level 2: Higher-Order Pattern
**Use when**:
- ✅ Need to reuse comparators
- ✅ Want separation of concerns
- ✅ Learning advanced patterns
- ✅ Internal tools (convenience > extreme performance)

**Avoid when**:
- ❌ Need runtime sort order selection
- ❌ Performance is critical
- ❌ Need comprehensive error handling

**Example**: [`generics.go:188`](/learning/gobyexample/examples/generics/generics.go#L188)

---

### Level 3: Production-Grade
**Use when**:
- ✅ Building frameworks/libraries
- ✅ Performance critical (large datasets)
- ✅ Need runtime configurability
- ✅ Multiple data structures (maps + structs)
- ✅ Error handling mandatory
- ✅ Production systems

**Example**: [`slicesorter.go:17`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L17)

---

## Code Location Quick Reference

| Concept | Tutorial Level 1 | Tutorial Level 2 | Production Level 3 |
|---------|------------------|------------------|----------------------|
| **Direct sorting** | [`generics.go:129`](/learning/gobyexample/examples/generics/generics.go#L129) | - | - |
| **Comparator type** | - | [`generics.go:188`](/learning/gobyexample/examples/generics/generics.go#L188) | [`slicesorter.go:14`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L14) |
| **String sorter** | [`generics.go:129`](/learning/gobyexample/examples/generics/generics.go#L129) | [`generics.go:188`](/learning/gobyexample/examples/generics/generics.go#L188) | [`slicesorter.go:17`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L17) |
| **Int sorter** | [`generics.go:149`](/learning/gobyexample/examples/generics/generics.go#L149) | [`generics.go:195`](/learning/gobyexample/examples/generics/generics.go#L195) | [`slicesorter.go:42`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L42) |
| **Generic sorter** | - | [`generics.go:230`](/learning/gobyexample/examples/generics/generics.go#L230) | - |
| **Sort function** | [`generics.go:129`](/learning/gobyexample/examples/generics/generics.go#L129) | [`generics.go:246`](/learning/gobyexample/examples/generics/generics.go#L246) | [`collections.go:345`](/aac-backend/internal/collections/collections.go#L345) |
| **Field extraction** | Inline | Inline | [`slicesorter.go:67`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L67) |

---

## Progressive Learning Path

### Phase 1: Understand Level 1 (Simple)
✅ Run `go run generics.go` - See all three levels in action
✅ Study [`SortByStringField`](/learning/gobyexample/examples/generics/generics.go#L129)
✅ Understand why reflection is needed
✅ Recognize the limitations (no reusability, reflection per comparison)

### Phase 2: Grasp Level 2 (Higher-Order Functions)
✅ Study [`Comparator[T]`](/learning/gobyexample/examples/generics/generics.go#L188) type definition
✅ Examine [`NewStringSorter`](/learning/gobyexample/examples/generics/generics.go#L188) factory
✅ Compare [`NewIntSorter`](/learning/gobyexample/examples/generics/generics.go#L195) and [`NewGenericSorter`](/learning/gobyexample/examples/generics/generics.go#L230)
✅ Understand reusability benefit
✅ Recognize remaining limitation (reflection still happens per comparison)

### Phase 3: Bridge to Level 3 (Production)
✅ Read [`slicesorter.go:17`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L17)
✅ Compare tutorial Level 2 vs production Level 3
✅ Understand the 2-level abstraction (`SliceSorter[T]`)
✅ See why production adds runtime sort order selection
✅ Study the helper function pattern ([`getStringFieldValue`](/aac-backend/internal/collections/sliceutils/slicesorter.go#L67))

### Phase 4: Apply the Pattern
✅ Identify duplication in your own code
✅ Start with Level 1 for simple cases
✅ Upgrade to Level 2 when you need reusability
✅ Refactor to Level 3 only when performance matters

### Phase 5: Master the Architecture
✅ Study how [`collections.go:345`](/aac-backend/internal/collections/collections.go#L345) uses the sorters
✅ Understand the map of sorter factories pattern
✅ See how it eliminates 200+ functions
✅ Recognize when each level is appropriate

---

## Summary

The tutorial demonstrates **all three levels** of sophistication, showing the complete evolution from simple to production-grade.

### What Each Level Teaches

**Level 1** (Simple Direct):
- ✅ **What** generics + reflection can do
- ✅ Why reflection is needed (no field constraints in Go)
- ✅ Basic pattern for dynamic field access
- 📚 **Best for**: First-time learning

**Level 2** (Higher-Order):
- ✅ **How** to structure code with factory pattern
- ✅ Benefits of separation of concerns
- ✅ Comparator reusability
- ✅ Type-specific vs universal sorters trade-off
- 📚 **Best for**: Understanding software architecture patterns

**Level 3** (Production):
- ✅ **Why** production adds complexity
- ✅ Runtime configurability requirements
- ✅ Performance optimization strategies
- ✅ Comprehensive error handling
- ✅ Multi-data-structure support
- 📚 **Best for**: Building production systems

### The Progression

```
Level 1: Simple & Direct
    ↓ (Add reusability)
Level 2: Factory Pattern
    ↓ (Add runtime config + optimize reflection)
Level 3: Production-Grade
```

### Key Takeaway

**All three levels are correct for their contexts**:
- Level 1: Quick scripts, learning
- Level 2: Internal tools, medium complexity
- Level 3: Production APIs, high performance, critical systems

Choose based on: dataset size, performance needs, error handling requirements, and configurability demands.
