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

}

func TestWrappingMultipleErrors(t *testing.T) {

	err := FuncD()
	errStr := ErrorWithStack(err)
	t.Log(errStr)

	if strings.Count(errStr, "\n") < 5 {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}

}

func TestErrorStack(t *testing.T) {

	err := FuncD()
	errStack := ErrorStack(err)

	t.Log(errStack)

	if strings.Count(errStack, "\n") < 5 {
		t.Errorf("ErrorWithStack() failed: %s", errStack)
		return
	}

}

func TestNewStackErrWithWrappedError(t *testing.T) {

	err := FuncA()
	formattedErr := ErrorWithStack(err)
	t.Log(formattedErr)

	if strings.Count(formattedErr, "\n") < 3 {
		t.Errorf("ErrorWithStack() failed: %s", formattedErr)
		return
	}

}

func TestErrorWithStack(t *testing.T) {
	err := fmt.Errorf("test error")
	stackErr := NewStackError(err)
	errStr := ErrorWithStack(stackErr)
	t.Log(errStr)

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
