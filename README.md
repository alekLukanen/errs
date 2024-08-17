# errs
A Go package which allows you to attach stack traces to errors.

## Example

Simply create the initial error with the `NewStackError()` function and then use the `Wrap()` function each time you want to add a new error to the error on the way back up the stack. Finally when you are ready to log the error you can use the `ErrorWithStack()` function to get the full error message with the stack trace.
```go
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
	fmt.Printf(ErrStack)
}
```

Running the above example:

```shell
$ go run -trimpath cmd/example/main.go
Error Messages
- [0] error in FuncB
- [1] received error from FuncB()
Primary Stack Trace
main.FuncB()
	./main.go:16 +0x2b
main.FuncA()
	./main.go:10 +0x13
main.main()
	./main.go:20 +0x13```
```

Because this package uses the `runtime` stack function you should always run your program with the `-trimpath` build option set so that the file paths are removed from the stack trace.
