package errs

import (
	"fmt"
	"strings"
	"testing"
)

func FunA() error {
	return Wrap(FuncB(), "wrapping error in FuncA()")
}

func FuncB() error {
	return NewStackError(fmt.Errorf("error from FuncB"))
}

func TestNewStackErrWithWrappedError(t *testing.T) {
	err := FunA()
	fmt.Println(ErrorWithStack(err))
}

func TestErrorWithStack(t *testing.T) {
	err := fmt.Errorf("test error")
	stackErr := NewStackError(err)
	fmt.Println(ErrorWithStack(stackErr))

	errStr := ErrorWithStack(stackErr)
	if !strings.Contains(errStr, "test error") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}
	if !strings.Contains(errStr, "errs/errs_test.go") {
		t.Errorf("ErrorWithStack() failed: %s", errStr)
		return
	}

}
