package structfill

import (
	"gopkg.in/yaml.v2"
	"reflect"
)

type config struct {
	CustomTypes    map[string]any `yaml:"custom_types"`
	Int            int64          `yaml:"int"`
	Uint           uint64         `yaml:"uint"`
	Float          float64        `yaml:"float"`
	StringValue    string         `yaml:"string"`
	Bool           bool           `yaml:"bool"`
	complex        complex128     `yaml:"_"` // complex cannot be marshalled
	Debug          bool           `yaml:"debug"`
	PanicOnUnknown bool           `yaml:"panic_on_unknown"`
}

func (cfg *config) String() string {
	wrapper := struct {
		Config config `json:"config"`
	}{Config: *cfg}
	out, err := yaml.Marshal(wrapper)
	if err != nil {
		panic(err)
	}
	return string(out)
}

func makeDefaultConfig() config {
	return config{
		CustomTypes:    make(map[string]any),
		Int:            1,
		Uint:           2,
		Float:          3,
		StringValue:    "string",
		Bool:           true,
		complex:        complex(4, 5),
		Debug:          false,
		PanicOnUnknown: false,
	}
}

type Option func(c *config)
type Options []Option

func WithBool(b bool) Option {
	return func(c *config) {
		c.Bool = b
	}
}

func WithString(s string) Option {
	return func(c *config) {
		c.StringValue = s
	}
}

func WithInt(i int) Option {
	return func(c *config) {
		c.Int = int64(i)
	}
}

func WithUint(u uint) Option {
	return func(c *config) {
		c.Uint = uint64(u)
	}
}

func WithFloat(f float64) Option {
	return func(c *config) {
		c.Float = f
	}
}

func WithComplex(x complex128) Option {
	return func(c *config) {
		c.complex = x
	}
}

func WithCustomType[S any](aStruct S) Option {
	return func(c *config) {
		if c.CustomTypes == nil {
			c.CustomTypes = make(map[string]any)
		}
		c.CustomTypes[reflect.ValueOf(aStruct).Type().String()] = aStruct
	}
}

func WithDebug() Option {
	return func(c *config) {
		c.Debug = true
	}
}

func WithPanicOnUnknown() Option {
	return func(c *config) {
		c.PanicOnUnknown = true
	}
}
