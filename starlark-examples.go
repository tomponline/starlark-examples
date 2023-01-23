package main

import (
	"fmt"

	"go.starlark.net/starlark"
)

// An embedded starlark script.
const fib = `
print("Starlark script started: ", greeting, "\n")
print("inst is: ", inst, "\n")

placementMember("foo")
placementRefuse("not allowed")

def fibonacci(n):
    print("Starlark fibonacci function called")

    res = list(range(n))
    for i in res[2:]:
        res[i] = res[i-2] + res[i-1]
    return res
`

var placementMember = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var memberName string

	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "member", &memberName); err != nil {
		return nil, err
	}

	fmt.Printf("instPlacement: %v\n", memberName)

	return starlark.None, nil
}

var placementRefuse = func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var reason string

	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "reason", &reason); err != nil {
		return nil, err
	}

	fmt.Printf("placementRefuse: %v\n", reason)

	return starlark.None, nil
}

func main() {

	s := starlark.NewDict(0)
	s.SetKey(starlark.String("type"), starlark.String("image"))
	s.SetKey(starlark.String("allow_inconsistent"), starlark.Bool(false))

	d := starlark.NewDict(5)
	d.SetKey(starlark.String("name"), starlark.String("foo"))
	d.SetKey(starlark.String("stateful"), starlark.Bool(false))
	d.SetKey(starlark.String("profiles"), starlark.NewList(nil))
	d.SetKey(starlark.String("config"), starlark.NewDict(0))
	d.SetKey(starlark.String("source"), s)

	// Define the pre-declared global starlark environment the script will be run from.
	predeclared := starlark.StringDict{
		"greeting":        starlark.String("hello"),                                // Global string.
		"placementMember": starlark.NewBuiltin("placementMember", placementMember), // Global function.
		"placementRefuse": starlark.NewBuiltin("placementRefuse", placementRefuse), // Global function.
		"inst":            d,
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
