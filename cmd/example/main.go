package main

import (
	"fmt"
	"github.com/alekLukanen/errs"
)

func FuncA() error {
	return errs.Wrap(
		FuncB(),
		fmt.Errorf("received error from FuncB()"),
	)
}

func FuncB() error {
	return errs.NewStackError(fmt.Errorf("error in FuncB"))
}

func main() {
	err := FuncA()
	ErrStack := errs.ErrorWithStack(err)
	fmt.Print(ErrStack)
}
