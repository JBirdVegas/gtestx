## structfill

```go
package main

import "reflect"
import "github.com/jbirdvegas/gtestx/pkg/structfill"

type foo struct {
	One string
	Two int
}

func main() {
	f := foo{}
	if err := structfill.AutoFill(&f); err != nil {
		panic(err)
	}
	reflect.DeepEqual(foo{One: "string", Two: 1}, f)
}

```