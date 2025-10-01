// Starting with version 1.18, Go has added support for
// _generics_, also known as _type parameters_.

package main

import (
	"fmt"
	"reflect"
	"sort"
)

// As an example of a generic function, `SlicesIndex` takes
// a slice of any `comparable` type and an element of that
// type and returns the index of the first occurrence of
// v in s, or -1 if not present. The `comparable` constraint
// means that we can compare values of this type with the
// `==` and `!=` operators. For a more thorough explanation
// of this type signature, see [this blog post](https://go.dev/blog/deconstructing-type-parameters).
// Note that this function exists in the standard library
// as [slices.Index](https://pkg.go.dev/slices#Index).
func SlicesIndex[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// As an example of a generic type, `List` is a
// singly-linked list with values of any type.
type List[T any] struct {
	head, tail *element[T]
}

type element[T any] struct {
	next *element[T]
	val  T
}

// We can define methods on generic types just like we
// do on regular types, but we have to keep the type
// parameters in place. The type is `List[T]`, not `List`.
func (lst *List[T]) Push(v T) {
	if lst.tail == nil {
		lst.head = &element[T]{val: v}
		lst.tail = lst.head
	} else {
		lst.tail.next = &element[T]{val: v}
		lst.tail = lst.tail.next
	}
}

// AllElements returns all the List elements as a slice.
// In the next example we'll see a more idiomatic way
// of iterating over all elements of custom types.
func (lst *List[T]) AllElements() []T {
	var elems []T
	for e := lst.head; e != nil; e = e.next {
		elems = append(elems, e.val)
	}
	return elems
}

// ============================================================
// PART 2: THE PROBLEM - Why We Need Reflection with Generics
// ============================================================

// Let's say we want to sort ANY struct by ANY field name.
// Without reflection, we'd need to write a sorter for each combination:

type Person struct {
	Name string
	Age  int
}

type Product struct {
	Name  string
	Price float64
}

type Book struct {
	Title  string
	Author string
}

// ❌ PROBLEM: This won't compile!
// Generics can't access struct fields directly - no field constraints in Go
/*
func SortByField[T any](slice []T, fieldName string) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].fieldName < slice[j].fieldName  // ERROR: T has no fields!
	})
}
*/

// ============================================================
// SOLUTION 1: Without Generics (Type-Specific, Lots of Duplication)
// ============================================================

func SortPersonByName(people []Person) {
	sort.Slice(people, func(i, j int) bool {
		return people[i].Name < people[j].Name
	})
}

func SortPersonByAge(people []Person) {
	sort.Slice(people, func(i, j int) bool {
		return people[i].Age < people[j].Age
	})
}

func SortProductByName(products []Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].Name < products[j].Name
	})
}

func SortProductByPrice(products []Product) {
	sort.Slice(products, func(i, j int) bool {
		return products[i].Price < products[j].Price
	})
}

// Problem: We need 4 functions for just 2 types with 2 fields each!
// In aac-backend: 50+ types × 5+ fields each = 250+ functions needed! 😱

// ============================================================
// SOLUTION 2: Generics + Reflection (One Function for Everything)
// ============================================================

// Step 1: Generic function that accepts any type
// Step 2: Use reflection to access fields at runtime
func SortByStringField[T any](slice []T, fieldName string, ascending bool) error {
	sort.Slice(slice, func(i, j int) bool {
		// Use reflection to get field values at runtime
		valueI := reflect.ValueOf(slice[i]).FieldByName(fieldName)
		valueJ := reflect.ValueOf(slice[j]).FieldByName(fieldName)

		// Extract the actual string values
		strI := valueI.String()
		strJ := valueJ.String()

		// Compare based on sort order
		if ascending {
			return strI < strJ
		}
		return strI > strJ
	})
	return nil
}

// For numeric fields, we need a separate function (reflection returns different types)
func SortByIntField[T any](slice []T, fieldName string, ascending bool) error {
	sort.Slice(slice, func(i, j int) bool {
		valueI := reflect.ValueOf(slice[i]).FieldByName(fieldName)
		valueJ := reflect.ValueOf(slice[j]).FieldByName(fieldName)

		intI := valueI.Int() // reflection extracts as int64
		intJ := valueJ.Int()

		if ascending {
			return intI < intJ
		}
		return intI > intJ
	})
	return nil
}

func SortByFloatField[T any](slice []T, fieldName string, ascending bool) error {
	sort.Slice(slice, func(i, j int) bool {
		valueI := reflect.ValueOf(slice[i]).FieldByName(fieldName)
		valueJ := reflect.ValueOf(slice[j]).FieldByName(fieldName)

		floatI := valueI.Float() // reflection extracts as float64
		floatJ := valueJ.Float()

		if ascending {
			return floatI < floatJ
		}
		return floatI > floatJ
	})
	return nil
}

// ============================================================
// ADVANCED: Production Pattern from aac-backend
// ============================================================

// This mirrors the exact pattern used in:
// aac-backend/internal/collections/sliceutils/slicesorter.go

// Step 1: Define a type for comparison functions
// This is a "higher-order function" - a function that returns a function
type Comparator[T any] func(a, b T) bool

// Step 2: Factory function that creates type-specific comparators
// This is what makes the pattern so powerful!
func NewStringSorter[T any](fieldName string, ascending bool) Comparator[T] {
	// Return a closure that captures fieldName and ascending
	return func(a, b T) bool {
		// Use reflection to access the field
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

func NewIntSorter[T any](fieldName string, ascending bool) Comparator[T] {
	return func(a, b T) bool {
		fieldA := reflect.ValueOf(a).FieldByName(fieldName)
		fieldB := reflect.ValueOf(b).FieldByName(fieldName)

		intA := fieldA.Int()
		intB := fieldB.Int()

		if ascending {
			return intA < intB
		}
		return intA > intB
	}
}

func NewGenericSorter[T any](fieldName string, ascending bool) Comparator[T] {
	return func(a, b T) bool {
		fieldA := reflect.ValueOf(a).FieldByName(fieldName)
		fieldB := reflect.ValueOf(b).FieldByName(fieldName)

		// Handle different kinds of fields (string, int, float, etc.)
		switch fieldA.Kind() {
		case reflect.String:
			strA := fieldA.String()
			strB := fieldB.String()
			if ascending {
				return strA < strB
			}
			return strA > strB
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intA := fieldA.Int()
			intB := fieldB.Int()
			if ascending {
				return intA < intB
			}
			return intA > intB
		case reflect.Float32, reflect.Float64:
			floatA := fieldA.Float()
			floatB := fieldB.Float()
			if ascending {
				return floatA < floatB
			}
			return floatA > floatB
		default:
			// Unsupported field type
			return false
		}
	}
}

// Generic sort function using the comparator
func SortWithComparator[T any](slice []T, comparator Comparator[T]) {
	sort.Slice(slice, func(i, j int) bool {
		return comparator(slice[i], slice[j])
	})
}

// ============================================================
// DEMONSTRATION: See All Patterns in Action
// ============================================================

func main() {
	fmt.Println("=== BASIC GENERICS ===")
	var s = []string{"foo", "bar", "zoo"}

	// When invoking generic functions, we can often rely
	// on _type inference_. Note that we don't have to
	// specify the types for `S` and `E` when
	// calling `SlicesIndex` - the compiler infers them
	// automatically.
	fmt.Println("index of zoo:", SlicesIndex(s, "zoo"))

	// ... though we could also specify them explicitly.
	_ = SlicesIndex[[]string, string](s, "zoo")

	lst := List[int]{}
	lst.Push(10)
	lst.Push(13)
	lst.Push(23)
	fmt.Println("list:", lst.AllElements())

	fmt.Println("\n=== GENERICS + REFLECTION ===")

	// Create test data
	people := []Person{
		{Name: "Charlie", Age: 35},
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	products := []Product{
		{Name: "Laptop", Price: 999.99},
		{Name: "Mouse", Price: 29.99},
		{Name: "Keyboard", Price: 79.99},
	}

	// Problem: Without generics + reflection, we need separate functions
	fmt.Println("\n--- Old Way (Type-Specific Functions) ---")
	SortPersonByName(people)
	fmt.Println("People sorted by name:", people)

	SortProductByPrice(products)
	fmt.Println("Products sorted by price:", products)

	// Reset data
	people = []Person{
		{Name: "Charlie", Age: 35},
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	// Solution: One generic function works for all types!
	fmt.Println("\n--- New Way (Generics + Reflection) ---")
	SortByStringField(people, "Name", true)
	fmt.Println("People sorted by Name (generic):", people)

	SortByIntField(people, "Age", false)
	fmt.Println("People sorted by Age descending (generic):", people)

	SortByStringField(products, "Name", false)
	fmt.Println("Products sorted by Name descending (generic):", products)

	// Advanced: Production pattern with comparators
	fmt.Println("\n--- Production Pattern (Higher-Order Functions) ---")
	people = []Person{
		{Name: "Charlie", Age: 35},
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	// Create a reusable comparator
	nameComparator := NewStringSorter[Person]("Name", true)
	SortWithComparator(people, nameComparator)
	fmt.Println("People sorted with comparator:", people)

	// The power: Same comparator factory works for ANY type!
	productComparator := NewStringSorter[Product]("Name", false)
	SortWithComparator(products, productComparator)
	fmt.Println("Products sorted with comparator:", products)

	fmt.Println("\n=== KEY INSIGHTS ===")
	fmt.Println("1. Generics alone can't access struct fields (no field constraints)")
	fmt.Println("2. Reflection enables runtime field access by name")
	fmt.Println("3. Generics + Reflection = Type-safe container + Dynamic field access")
	fmt.Println("4. Higher-order functions (Comparator pattern) = Maximum flexibility")
	fmt.Println("\nIn aac-backend:")
	fmt.Println("- This pattern eliminates 200+ duplicate sorting functions")
	fmt.Println("- Single implementation handles all models (Client, Receipt, Activity, etc.)")
	fmt.Println("- Reflection cost is negligible for API response sorting")

	books := []Book{
		{Title: "The Go Programming Language", Author: "Alan A. A. Donovan"},
		{Title: "Introducing Go", Author: "Caleb Doxsey"},
		{Title: "Go in Action", Author: "William Kennedy"},
	}
	SortByStringField(books, "Title", true)
	fmt.Println("Books sorted by Title (generic):", books)

	bookComparatorByTitle := NewStringSorter[Book]("Title", true)
	SortWithComparator(books, bookComparatorByTitle)
	fmt.Println("Books sorted with comparator by Title:", books)

	SortByFloatField(products, "Price", false)
	fmt.Println("Products sorted by Price descending (generic):", products)

	peopleComparatorByAge := NewIntSorter[Person]("Age", true)
	SortWithComparator(people, peopleComparatorByAge)
	fmt.Println("People sorted with comparator by Age:", people)

	bookComparatorByTitle = NewGenericSorter[Book]("Title", false)
	SortWithComparator(books, bookComparatorByTitle)
	fmt.Println("Books sorted with comparator by Title:", books)
}
