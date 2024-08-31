package errs

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

var ErrB = fmt.Errorf("error from FuncB")
var ErrC = fmt.Errorf("error from FuncC")

func FuncA() error {
	return Wrap(FuncB(), fmt.Errorf("received error in FuncA()"))
}

func FuncB() error {
	return NewStackError(ErrB)
}

func FuncC() error {
	return NewStackError(ErrC)
}

func FuncD() error {
	errB := FuncB()
	errC := FuncC()
	err := Wrap(
		errB,
		fmt.Errorf("received error from FuncB()"),
		fmt.Errorf("while handling error from FuncB() received an error from FuncC()"),
		errC,
	)
	return err
}

func TestStackErrorCanUseErrorsIs(t *testing.T) {

	err := FuncD()
	if !errors.Is(err, ErrB) {
		t.Log(err.(*StackError).wrappedErrs)
		t.Errorf("failed, expected the error to wrap ErrB: %s", err)
		return
	}
	if !errors.Is(err, ErrC) {
		t.Errorf("failed, expected the error to wrap ErrC: %s", err)
		return
	}
	if errors.Is(err, fmt.Errorf("some other error")) {
		t.Errorf("failed: %s", err)
		return
	}

	t.Log("funcs: ", err.(*StackError).funcNames)
	t.Log("lines: ", err.(*StackError).lineNumbers)

}

func TestWrappingMultipleErrors(t *testing.T) {
	sb := strings.Builder{}
	sb.WriteString("Error Messages\n")
	sb.WriteString("- [0] error from FuncB\n")
	sb.WriteString("- [1] received error from FuncB()\n")
	sb.WriteString("- [2] while handling error from FuncB() received an error from FuncC()\n")
	sb.WriteString("- [3] error from FuncC\n")
	sb.WriteString("Primary Stack Trace\n")
	expectedPrefix := sb.String()

	err := FuncD()
	errStr := ErrorWithStack(err)
	fmt.Print(errStr)

	if !strings.HasPrefix(errStr, expectedPrefix) {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}

}

func TestErrorStack(t *testing.T) {

	err := FuncD()
	errStack := ErrorStack(err)

	fmt.Print(errStack)

	if !strings.Contains(errStack, "errs/errs_test.go") {
		t.Errorf("ErrorStack() failed: %s", errStack)
	}
}

func TestErrorMessage(t *testing.T) {
	sb := strings.Builder{}
	sb.WriteString("Error Messages\n")
	sb.WriteString("- [0] error from FuncB\n")
	sb.WriteString("- [1] received error from FuncB()\n")
	sb.WriteString("- [2] while handling error from FuncB() received an error from FuncC()\n")
	sb.WriteString("- [3] error from FuncC")
	expectedStr := sb.String()

	err := FuncD()
	errStr := ErrorMessage(err)
	fmt.Print(errStr)

	if errStr != expectedStr {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}

}

func TestNewStackErrWithWrappedError(t *testing.T) {
	sb := strings.Builder{}
	sb.WriteString("Error Messages\n")
	sb.WriteString("- [0] error from FuncB\n")
	sb.WriteString("- [1] received error in FuncA()\n")
	sb.WriteString("Primary Stack Trace\n")
	expectedPrefix := sb.String()

	err := FuncA()
	formattedErr := ErrorWithStack(err)
	fmt.Print(formattedErr)

	if !strings.HasPrefix(formattedErr, expectedPrefix) {
		t.Errorf("ErrorWithStack() failed: %s", formattedErr)
		return
	}

}

func TestErrorWithStack(t *testing.T) {
	err := fmt.Errorf("test error")
	stackErr := NewStackError(err)
	errStr := ErrorWithStack(stackErr)
	fmt.Print(errStr)

	if !strings.Contains(errStr, "test error") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}
	if !strings.Contains(errStr, "errs/errs_test.go") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}

}

func TestErrorWithStack_noStack(t *testing.T) {
	err := fmt.Errorf("test error")
	errStr := ErrorWithStack(err)
	fmt.Print(errStr)
	if !strings.Contains(errStr, "test error") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}
	if !strings.Contains(errStr, "[No Stack]") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}
}
