## gtestx

Help for testing golang apps

- [structfill](./pkg/structfill/README.md) populates a struct according to options supplied by the caller. This is
  particularly useful for testing functions that translate one object type to another.

```go
package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/jbirdvegas/gtestx/pkg/structfill"
)

type inner struct {
	One string
}

type outer struct {
	One   string
	Two   int
	Three inner
	Four  string
}

func main() {
	opts := structfill.Options{
		structfill.WithString("foobar"),
		structfill.WithInt(123),
		structfill.WithCustomType(inner{One: "Hello World!"}),
	}
	got := outer{}
	if err := structfill.AutoFill(&got, opts...); err != nil {
		panic(err)
	}
	want := outer{
		One: "foobar",
		Two: 123,
		Three: inner{
			One: "Hello World!",
		},
		Four: "foobar",
	}
	if diff := cmp.Diff(want, got); len(diff) > 0 {
		panic(fmt.Sprintf("AutoFill() mismatch (-want +got):\n%s", diff))
	}
	println("Done")
}

```