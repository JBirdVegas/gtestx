package structfill

import (
	"fmt"
	"reflect"
)

// populate recursively sets default values for the provided reflect.Value.
func populate(v reflect.Value, cfg *config) error {
	if cfg.Debug {
		fmt.Printf("Populating type: %s\n", v.Type().String())
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			newVal := reflect.New(v.Type().Elem())
			v.Set(newVal)
			return populate(newVal.Elem(), cfg)
		}
		return populate(v.Elem(), cfg)
	case reflect.Struct:
		val, has := cfg.CustomTypes[v.Type().String()]
		if has {
			if cfg.Debug {
				fmt.Printf("replacing field with supplied custom type: %s\n", v.Type().String())
			}
			v.Set(reflect.ValueOf(val))
			break
		}

		// Process each field.
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// Only set exported fields.
			if !field.CanSet() {
				continue
			}
			if err := populate(field, cfg); err != nil {
				return err
			}
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := populate(v.Index(i), cfg); err != nil {
				return err
			}
		}
	case reflect.Slice:
		// Create a slice of length 1 if nil.
		if v.IsNil() {
			newSlice := reflect.MakeSlice(v.Type(), 1, 1)
			v.Set(newSlice)
		}
		// Populate each element.
		for i := 0; i < v.Len(); i++ {
			if err := populate(v.Index(i), cfg); err != nil {
				return err
			}
		}
	case reflect.Map:
		// Initialize map if nil.
		if v.IsNil() {
			v.Set(reflect.MakeMap(v.Type()))
		}
		// Try to create a default key/value pair.
		key := reflect.Zero(v.Type().Key())
		switch key.Kind() {
		case reflect.String:
			key = reflect.ValueOf(cfg.StringValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			key = reflect.ValueOf(cfg.Int)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			key = reflect.ValueOf(cfg.Uint)
		case reflect.Float32, reflect.Float64:
			key = reflect.ValueOf(cfg.Float)
		case reflect.Bool:
			key = reflect.ValueOf(cfg.Bool)
		case reflect.Complex64, reflect.Complex128:
			key = reflect.ValueOf(cfg.complex)
		}
		val := reflect.New(v.Type().Elem()).Elem()
		if err := populate(val, cfg); err != nil {
			return err
		}
		v.SetMapIndex(key, val)
	case reflect.Interface:
		// If interface is nil, we cannot determine the concrete type; skip.
		if v.IsNil() {
			return nil
		}
		return populate(v.Elem(), cfg)
	default:
		// For basic types, set defaults.
		setBasicValue(v, cfg)
	}
	return nil
}

// setBasicValue sets default values for supported basic types.
func setBasicValue(v reflect.Value, cfg *config) {
	if cfg.Debug {
		fmt.Printf("Populating type: %s\n", v.Type().String())
	}
	switch v.Kind() {
	case reflect.String:
		v.SetString(cfg.StringValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(cfg.Int)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		v.SetUint(cfg.Uint)
	case reflect.Float32, reflect.Float64:
		v.SetFloat(cfg.Float)
	case reflect.Complex64, reflect.Complex128:
		v.SetComplex(cfg.complex)
	case reflect.Bool:
		v.SetBool(cfg.Bool)
	default:
		if cfg.Debug {
			fmt.Printf("Unhandled type encountered: %s\n", v.Type().String())
		}
		if cfg.PanicOnUnknown {
			panic(fmt.Errorf("unhandled type encountered: %s", v.Type().String()))
		}
	}
}
