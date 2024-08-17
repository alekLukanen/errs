package errs

import (
	"fmt"
	"strings"
	"testing"
)

func FuncA() error {
	return Wrap(FuncB(), fmt.Errorf("received error in FuncA()"))
}

func FuncB() error {
	return NewStackError(fmt.Errorf("error from FuncB"))
}

func FuncC() error {
	return NewStackError(fmt.Errorf("error from FuncC"))
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

func TestWrappingMultipleErrors(t *testing.T) {
	err := FuncD()
	errStr := ErrorWithStack(err)
	fmt.Printf(errStr)
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
	fmt.Printf(formattedErr)

	if !strings.HasPrefix(formattedErr, expectedPrefix) {
		t.Errorf("ErrorWithStack() failed: %s", formattedErr)
		return
	}

}

func TestErrorWithStack(t *testing.T) {
	err := fmt.Errorf("test error")
	stackErr := NewStackError(err)
	errStr := ErrorWithStack(stackErr)
	fmt.Printf(errStr)

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
	fmt.Printf(errStr)
	if !strings.Contains(errStr, "test error") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}
	if !strings.Contains(errStr, "[No Stack]") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}
}
