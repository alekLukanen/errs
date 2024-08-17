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
	fmt.Printf(errStr)

	if !strings.HasPrefix(errStr, expectedPrefix) {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}

}

func TestErrorStack(t *testing.T) {

	err := FuncD()
	errStack := ErrorStack(err)

	fmt.Printf(errStack)

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
	fmt.Printf(errStr)

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
