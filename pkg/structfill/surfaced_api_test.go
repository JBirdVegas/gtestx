package structfill

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

type TestStruct struct {
	IntValue     int64
	UintValue    uint64
	FloatValue   float64
	StringValue  string
	BoolValue    bool
	ComplexValue complex128
}

type CustomType struct {
	Field1 string
	Field2 int
}

type OuterType struct {
	CType       CustomType
	StringValue string
}

func TestPopulateStruct(t *testing.T) {
	tests := []struct {
		name     string
		options  []Option
		expected TestStruct
	}{
		{
			name:     "Default Values",
			options:  nil,
			expected: TestStruct{1, 2, 3, "string", true, complex(4, 5)},
		},
		{
			name:     "Custom Integer",
			options:  []Option{WithInt(42)},
			expected: TestStruct{42, 2, 3, "string", true, complex(4, 5)},
		},
		{
			name:     "Custom Uint",
			options:  []Option{WithUint(99)},
			expected: TestStruct{1, 99, 3, "string", true, complex(4, 5)},
		},
		{
			name:     "Custom Float",
			options:  []Option{WithFloat(9.99)},
			expected: TestStruct{1, 2, 9.99, "string", true, complex(4, 5)},
		},
		{
			name:     "Custom StringValue",
			options:  []Option{WithString("hello")},
			expected: TestStruct{1, 2, 3, "hello", true, complex(4, 5)},
		},
		{
			name:     "Custom Bool",
			options:  []Option{WithBool(false)},
			expected: TestStruct{1, 2, 3, "string", false, complex(4, 5)},
		},
		{
			name:     "Custom complex",
			options:  []Option{WithComplex(complex(1, 1))},
			expected: TestStruct{1, 2, 3, "string", true, complex(1, 1)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s TestStruct
			err := AutoFill(&s, tt.options...)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, s, cmp.AllowUnexported(TestStruct{})); len(diff) > 0 {
				t.Errorf("got diff:\n%s", diff)
			}
		})
	}

	t.Run("Custom Type Population", func(t *testing.T) {
		customType := CustomType{Field1: "foo", Field2: 42}
		out := OuterType{}
		err := AutoFill(&out, WithCustomType(customType), WithDebug())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(out.CType, customType) {
			t.Errorf("expected %+v, got %+v", customType, out.CType)
		}
	})
}
