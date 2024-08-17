package errs

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var ERR_MSG_TITLE string = "Error Messages"
var ERR_STACK_TITLE string = "Stack Trace"
var MAX_STACK_SIZE int = 4096

// StackError is a custom error type that includes a stack trace
type StackError struct {
	err   error
	stack string
	wraps int
}

// Returns the error message without the stack trace
func (obj *StackError) Error() string {
	return obj.err.Error()
}

// Returns the error message with the stack trace
func (obj *StackError) ErrorWithStack() string {
	return obj.Error() + "\n" + string(obj.stack)
}

// Get just the stack trace
func (obj *StackError) Stack() string {
	return string(obj.stack)
}

// Create a new error with a stack trace
func NewStackError(err error) *StackError {
	stack := make([]byte, MAX_STACK_SIZE)
	runtime.Stack(stack, false)
	formattedErr := fmt.Errorf("- [0] %w", err)
	return &StackError{err: formattedErr, stack: cleanedStack(stack)}
}

// Clean the stack trace of the firth 2 lines which contain the
// function call to this package.
func cleanedStack(stack []byte) string {
	stackStr := string(stack)
	// remove the first two lines from the stack trace
	cleanStr := stackStr[strings.Index(stackStr, "\n")+1:]
	cleanStr = cleanStr[strings.Index(cleanStr, "\n")+1:]
	cleanStr = cleanStr[strings.Index(cleanStr, "\n")+1:]
	return cleanStr
}

// This function allows you to pass in an arbirary error and get the
// error message and stack trace if that error wraps another error
// with a stack trace or is a stack trace error itself. It assumes only
// one error contains a stack trace and will return the first error
// that is of the StackError type.
// If the error does not contain a stack trace, it will return the
// error message and a "No Stack" message for the stack trace.
func ErrorWithStack(err error) string {
	var stackErr *StackError
	ok := errors.As(err, &stackErr)
	if ok {
		return fmt.Sprintf("%s\n%s\n%s\n%s", ERR_MSG_TITLE, err.Error(), ERR_STACK_TITLE, stackErr.Stack())
	} else {
		return fmt.Sprintf("%s\n[No Stack]")
	}
}

// wrap allows you to wrap errors and maintain an ordered list of error
// messages with an index in front of them.
// Here is an example:
// Error Messages
// - [0] error from FuncB
// - [1] wrapping error in FuncA()
// Stack Trace
// ...
func Wrap(err error, msg string) error {
	var stackErr *StackError
	if errors.As(err, &stackErr) {
		newErr := fmt.Errorf("%w\n- [%d] %s", err, stackErr.wraps+1, msg)
		stackErr.wraps += 1
		return newErr
	} else {
		return fmt.Errorf("- %w\n- %s", err, msg)
	}
}
