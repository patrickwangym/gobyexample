// It's possible to use custom types as `error`s by
// implementing the `Error()` method on them. Here's a
// variant on the example above that uses a custom type
// to explicitly represent an argument error.

package main

import (
	"errors"
	"fmt"
)

// A custom error type usually has the suffix "Error".
type argError struct {
	arg     int
	message string
}

// Adding this `Error` method makes `argError` implement
// the `error` interface.
func (e *argError) Error() string {
	return fmt.Sprintf("%d - %s", e.arg, e.message)
}

func f(arg int) (int, error) {
	if arg == 42 {

		// Return our custom error.
		return -1, &argError{arg, "can't work with it"}
	}
	return arg + 3, nil
}

func main() {

	// `errors.As` is a more advanced version of `errors.Is`.
	// It checks that a given error (or any error in its chain)
	// matches a specific error type and converts to a value
	// of that type, returning `true`. If there's no match, it
	// returns `false`.
	_, err := f(42)
	var ae *argError
	if errors.As(err, &ae) {
		fmt.Println(ae.arg)
		fmt.Println(ae.message)
	} else {
		fmt.Println("err doesn't match argError")
	}
}

/*
================================================================================
ERRORS.IS vs ERRORS.AS - Complete Guide
================================================================================

KEY DIFFERENCE:
- errors.Is: Value/identity comparison (checks if errors are THE SAME)
- errors.As: Type comparison (checks if error is OF A TYPE)

================================================================================
PATTERN 1: SENTINEL ERRORS with errors.Is()
================================================================================

What is a Sentinel Error?
A predefined, package-level error variable that serves as a unique identifier.

Example:
    // Define sentinel ONCE at package level
    var ErrNotFound = errors.New("not found")
    var ErrInvalidOTP = &InvalidOTPError{}

    // Use errors.Is() to compare against sentinel
    if errors.Is(err, ErrNotFound) {
        // Handles the specific sentinel error
    }

✅ When to use:
- Error has no fields (simple error)
- Same error instance reused across codebase
- Standard library errors (io.EOF, gorm.ErrRecordNotFound)

✅ Why it works:
- Sentinel is created ONCE → same memory address everywhere
- errors.Is() does pointer/value comparison
- Reliable across all comparisons

❌ NEVER do this (without sentinel):
    if errors.Is(err, &InvalidOTPError{}) {  // Creates NEW pointer each time!
        // This is UNRELIABLE - different pointer addresses
    }

Real-world example:
    // In GORM library
    var ErrRecordNotFound = errors.New("record not found")
    
    // Your code
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        // Handle not found
    }

================================================================================
PATTERN 2: CUSTOM ERROR TYPES with errors.As()
================================================================================

Use errors.As() when:
1. No sentinel error defined
2. Error has fields you might need to access
3. Checking error type (not specific error value)

Example WITHOUT field access:
    // Just type checking, don't care about fields
    if errors.As(err, new(*InvalidOTPError)) {
        return err  // Don't need error details
    }

Example WITH field access:
    // Type checking + accessing fields
    var notFoundErr *RecordNotFoundError
    if errors.As(err, &notFoundErr) {
        fmt.Printf("Model %s not found", notFoundErr.Model)
    }

✅ Why use errors.As():
- Type-safe checking (works with error chain)
- Doesn't require sentinel definition
- Can access error fields after match
- Clear intent: "Is this error TYPE in the chain?"

Real-world example:
    // Stdlib http.MaxBytesError detection
    if errors.As(err, new(*http.MaxBytesError)) {
        return &FileUploadSizeError{}
    }

================================================================================
DECISION TREE
================================================================================

Is there a sentinel error defined?
├─ YES → Use errors.Is()
│         if errors.Is(err, ErrNotFound) { ... }
│
└─ NO → Use errors.As()
          ├─ Need fields? → if errors.As(err, &varErr) { use varErr.Field }
          └─ No fields?   → if errors.As(err, new(*ErrorType)) { ... }

================================================================================
QUICK REFERENCE
================================================================================

errors.Is() - "Is this THE SAME error?"
┌────────────────────────────────────────────────────────────────┐
│ ✅ Sentinel errors (defined once)                              │
│ ✅ Standard library errors (io.EOF, gorm.ErrRecordNotFound)    │
│ ✅ Only need yes/no answer                                     │
│ ❌ NEVER with &ErrorType{} (creates new pointer each time)     │
└────────────────────────────────────────────────────────────────┘

errors.As() - "Is this OF THIS TYPE?"
┌────────────────────────────────────────────────────────────────┐
│ ✅ Custom error types without sentinel                         │
│ ✅ Need to access error fields                                 │
│ ✅ Type checking (not value comparison)                        │
│ ✅ Works with new(*ErrorType) OR var e *ErrorType              │
└────────────────────────────────────────────────────────────────┘

================================================================================
COMMON PATTERNS FROM PRODUCTION CODE
================================================================================

Pattern 1: Sentinel with errors.Is()
    var ErrPersonNotFound = errors.New("person not found")
    
    if errors.Is(err, people.ErrPersonNotFound) {
        // Create new person instead
    }

Pattern 2: Type check without fields
    if errors.As(err, new(*http.MaxBytesError)) {
        return &externalerrors.FileUploadSizeError{}
    }

Pattern 3: Type check with field access
    var validationErr *ValidationError
    if errors.As(err, &validationErr) {
        fmt.Printf("Validation failed: %s", validationErr.Message)
    }

Pattern 4: Standard library sentinel
    if errors.Is(err, gorm.ErrRecordNotFound) {
        // Not found is OK here, continue...
    }

================================================================================
MEMORY & POINTER BEHAVIOR
================================================================================

Sentinel (Same Pointer):
    var ErrTest = &TestError{}  // Address: 0x1234
    errors.Is(err, ErrTest)     // Compares: 0x1234
    errors.Is(err, ErrTest)     // Compares: 0x1234 (SAME!)

Without Sentinel (Different Pointers):
    errors.Is(err, &TestError{})  // Address: 0x1234
    errors.Is(err, &TestError{})  // Address: 0x5678 (DIFFERENT!)
    // ❌ Unreliable comparison!

================================================================================
SUMMARY
================================================================================

errors.Is:  "Is this error THE SAME as sentinel X?" → yes/no
errors.As:  "Is this error OF TYPE X? Give it to me!" → yes/no + error value

Golden Rule:
- Sentinel defined?     → errors.Is()
- No sentinel?          → errors.As()
- Never use errors.Is() with &ErrorType{} (without sentinel)
*/
