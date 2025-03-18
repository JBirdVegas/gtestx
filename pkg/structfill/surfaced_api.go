package structfill

import (
	"errors"
	"fmt"
	"reflect"
)

// AutoFill populates a struct (pointed to by aStruct) recursively with default values.
func AutoFill[Ptr *t, t any](aStruct Ptr, options ...Option) error {
	cfg := makeDefaultConfig()
	for i := range options {
		options[i](&cfg)
	}
	rv := reflect.ValueOf(aStruct)
	switch {
	case rv.IsNil():
		return errors.New("aStruct must be a non-nil pointer")
	case rv.Elem().Kind() != reflect.Struct:
		return errors.New("aStruct must be a pointer to a struct")
	case cfg.Debug:
		fmt.Println(cfg.String())
	default:
	}

	return populate(rv.Elem(), &cfg)
}
