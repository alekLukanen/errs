package errs

import (
	"fmt"
	"runtime"
	"strings"
)

var ERR_MSG_TITLE string = "Error Messages"
var ERR_STACK_TITLE string = "Primary Stack Trace"
var MAX_STACK_SIZE int = 4096

// StackError is a custom error type that includes a stack trace
type StackError struct {
	err         error
	stack       string
	wrappedErrs []error
}

// Returns the error message without the stack trace
func (obj *StackError) Error() string {

	var sb strings.Builder
	sb.WriteString(ERR_MSG_TITLE)
	sb.WriteString("\n")
	sb.WriteString("- [0] ")
	sb.WriteString(obj.err.Error())

	if len(obj.wrappedErrs) == 0 {
		return sb.String()
	}

	sb.WriteString("\n")
	for i, err := range obj.wrappedErrs {
		sErr, ok := err.(*StackError)
		if ok {
			sb.WriteString(fmt.Sprintf("- [%d] %s", i+1, sErr.err.Error()))
		} else {
			sb.WriteString(fmt.Sprintf("- [%d] %s", i+1, err.Error()))
		}
		if i != len(obj.wrappedErrs)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (obj *StackError) Unwrap() []error {
	return obj.wrappedErrs
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
	return &StackError{err: err, stack: cleanedStack(stack)}
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
// error message and stack trace if that error is a StackError.
// If the error is not a StackError, it will return the
// error message and a "No Stack" message for the stack trace.
func ErrorWithStack(err error) string {
	stackErr, ok := err.(*StackError)
	if ok {
		return fmt.Sprintf("%s\n%s\n%s", stackErr.Error(), ERR_STACK_TITLE, stackErr.Stack())
	} else {
		return fmt.Sprintf("%s\n%s\n%s\n[No Stack]\n", ERR_MSG_TITLE, err.Error(), ERR_STACK_TITLE)
	}
}

// Get the error's string message
func ErrorMessage(err error) string {
	stackErr, ok := err.(*StackError)
	if ok {
		return stackErr.Error()
	} else {
		return fmt.Sprintf("%s\n", err.Error())
	}
}

// Get the error's stack trace if it is a StackError
// else return a "No Stack" message.
func ErrorStack(err error) string {
	stackErr, ok := err.(*StackError)
	if ok {
		return stackErr.Stack()
	} else {
		return "[No Stack]\n"
	}
}

// wrap allows you to wrap errors and maintain an ordered list of error
// messages with an index in front of them. The errors themselves are all
// stored in the top level StackError struct and are stored in the order
// provided to this function.
// Here is an example:
// Error Messages
// - [0] error from FuncB
// - [1] wrapping error in FuncA()
// Primary Stack Trace
// ...
func Wrap(primaryErr error, newErrs ...error) error {
	if len(newErrs) == 0 {
		return primaryErr
	}

	stackErr, ok := primaryErr.(*StackError)
	if ok {
		for _, newErr := range newErrs {
			stackErr.wrappedErrs = append(stackErr.wrappedErrs, newErr)
		}
		return stackErr
	} else {
		var retErr error
		for _, newErr := range newErrs {
			retErr = fmt.Errorf("- %w\n- %w", primaryErr, newErr)
		}
		return retErr
	}
}
