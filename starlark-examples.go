package main

import (
	"fmt"
	"strings"

	"go.starlark.net/starlark"
)

// An embedded starlark script.
const fib = `
print("Starlark script started: ", greeting, "\n")

print(repeat("foo1", 3), "\n")
print(repeat(s="foo2", n=2), "\n")

def fibonacci(n):
    print("Starlark fibonacci function called")

    res = list(range(n))
    for i in res[2:]:
        res[i] = res[i-2] + res[i-1]
    return res
`

// Defines a Go function that can be called from starlark called repeat(str, n=1).
var repeat = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var s string
	var n int = 1

	fmt.Printf("Go repeat function called: Position args: %+v, Keyword args: %+v\n", args, kwargs)

	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "s", &s, "n?", &n); err != nil {
		return nil, err
	}

	return starlark.String(strings.Repeat(s, n)), nil
}

func main() {
	// Define the pre-declared global starlark environment the script will be run from.
	predeclared := starlark.StringDict{
		"greeting": starlark.String("hello"),              // Global string.
		"repeat":   starlark.NewBuiltin("repeat", repeat), // Global function.
	}

	// Execute starlark script.
	thread := &starlark.Thread{Name: "my thread"}
	globals, err := starlark.ExecFile(thread, "somename.star", fib, predeclared)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Retrieve a global variable from starlark environment.
	fibonacci := globals["fibonacci"]

	// Call starlark function from Go.
	v, err := starlark.Call(thread, fibonacci, starlark.Tuple{starlark.MakeInt(10)}, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("fibonacci(10) result in Go is: %v\n", v) // fibonacci(10) = [0, 1, 1, 2, 3, 5, 8, 13, 21, 34]
}
