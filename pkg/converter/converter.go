package converter

import (
	"fmt"
	"go/token"
	"reflect"
	"strings"
	"unicode"

	rstypes "github.com/drewstone/go2rs/pkg/types"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// Converter converts golang types to Rust declarations
type Converter struct {
	types         map[reflect.Type]rstypes.Type
	paramNames    map[reflect.Type]string
	ConfigureFunc func(reflect.Type) FuncConfig
}

// FuncConfig configures how functions are converted to Rust
type FuncConfig struct {
	IsAsync         bool
	NoIgnoreContext bool
	IsMethod        bool
	MethodName      string
	ParamNames      []string
}

// NewConverter creates a new converter with primitive type mappings
func NewConverter() *Converter {
	c := &Converter{
		types:      make(map[reflect.Type]rstypes.Type),
		paramNames: make(map[reflect.Type]string),
	}

	// Map Go primitives to Rust types
	primitives := map[reflect.Type]rstypes.Type{
		reflect.TypeOf(false): &rstypes.Number{BitSize: 1}, // bool
		reflect.TypeOf(int(0)): &rstypes.Number{
			BitSize:  32,
			IsSigned: true,
		},
		reflect.TypeOf(int64(0)): &rstypes.Number{
			BitSize:  64,
			IsSigned: true,
		},
		reflect.TypeOf(uint(0)): &rstypes.Number{
			BitSize:  32,
			IsSigned: false,
		},
		reflect.TypeOf(float32(0)): &rstypes.Number{
			BitSize: 32,
			IsFloat: true,
		},
		reflect.TypeOf(float64(0)): &rstypes.Number{
			BitSize: 64,
			IsFloat: true,
		},
		reflect.TypeOf(""): &rstypes.String{},
	}

	for t, rustType := range primitives {
		c.types[t] = rustType
	}

	return c
}

// Convert converts a Go type to a Rust type
func (c *Converter) Convert(t reflect.Type) rstypes.Type {
	// Check if we've already converted this type
	if existing, ok := c.types[t]; ok {
		return existing
	}

	var result rstypes.Type
	switch t.Kind() {
	case reflect.Bool:
		result = &rstypes.Number{BitSize: 1} // bool
	case reflect.String:
		result = &rstypes.String{}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		result = &rstypes.Number{BitSize: 32, IsSigned: true}
	case reflect.Int64:
		result = &rstypes.Number{BitSize: 64, IsSigned: true}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		result = &rstypes.Number{BitSize: 32, IsSigned: false}
	case reflect.Uint64:
		result = &rstypes.Number{BitSize: 64, IsSigned: false}
	case reflect.Float32:
		result = &rstypes.Number{BitSize: 32, IsFloat: true}
	case reflect.Float64:
		result = &rstypes.Number{BitSize: 64, IsFloat: true}
	case reflect.Struct:
		result = c.convertStruct(t)
	case reflect.Slice:
		result = &rstypes.Vec{Inner: c.Convert(t.Elem())}
	case reflect.Array:
		result = &rstypes.Array{
			Inner: c.Convert(t.Elem()),
			Size:  uint64(t.Len()),
		}
	case reflect.Map:
		result = &rstypes.Map{
			Key:   c.Convert(t.Key()),
			Value: c.Convert(t.Elem()),
		}
	case reflect.Chan:
		result = &rstypes.Vec{Inner: c.Convert(t.Elem())} // Changed from Stream to Vec
	case reflect.Ptr:
		result = &rstypes.Nullable{Inner: c.Convert(t.Elem())}
	case reflect.Interface:
		if t.Name() == "error" {
			result = &rstypes.Result{
				Ok:  &rstypes.Unit{},
				Err: &rstypes.String{},
			}
		} else {
			result = &rstypes.Trait{Name: t.Name()}
		}
	case reflect.Func:
		result = c.convertFunc(t)
	default:
		result = &rstypes.Any{}
	}

	// Set package name and position for all types
	if result != nil {
		result.SetPackageName(t.PkgPath())
		// Set a default position since we don't have source position info
		result.SetPosition(&token.Position{})
	}

	c.types[t] = result
	return result
}

// convertStruct converts a Go struct to a Rust struct type
func (c *Converter) convertStruct(t reflect.Type) *rstypes.Struct {
	s := &rstypes.Struct{
		Name:   t.Name(),
		Fields: make(map[string]rstypes.StructField),
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !unicode.IsUpper(rune(field.Name[0])) {
			continue
		}

		fieldType := c.Convert(field.Type)

		// Handle struct tags
		serdeTag := field.Tag.Get("serde")
		rustTag := field.Tag.Get("rust")

		// Combine tags if both exist
		var finalTag string
		if serdeTag != "" {
			if rustTag != "" {
				finalTag = serdeTag + "," + rustTag
			} else {
				finalTag = serdeTag
			}
		} else {
			finalTag = rustTag
		}

		// Check if field is optional
		isOptional := field.Type.Kind() == reflect.Ptr ||
			strings.Contains(rustTag, "optional")

		s.Fields[field.Name] = rstypes.StructField{
			RawName:    field.Name,
			RawTag:     finalTag,
			FieldIndex: i,
			Type:       fieldType,
			Position:   &token.Position{},
			Optional:   isOptional,
		}
	}

	return s
}

// convertFunc converts a Go func to a Rust function type
func (c *Converter) convertFunc(t reflect.Type) *rstypes.Function {
	config := FuncConfig{}
	if c.ConfigureFunc != nil {
		config = c.ConfigureFunc(t)
	}

	fn := &rstypes.Function{
		Name:     t.Name(),
		IsMethod: config.IsMethod,
		IsAsync:  config.IsAsync,
		Params:   make([]rstypes.FunctionParam, 0),
	}

	// Handle parameters
	startIdx := 0
	if config.IsMethod {
		startIdx = 1 // Skip receiver for methods
	}

	for i := startIdx; i < t.NumIn(); i++ {
		paramType := t.In(i)

		// Skip context if configured
		if paramType.Name() == "Context" && !config.NoIgnoreContext {
			continue
		}

		paramName := c.getParamName(paramType, i, config.ParamNames)
		fn.Params = append(fn.Params, rstypes.FunctionParam{
			Name:     paramName,
			Type:     c.Convert(paramType),
			Position: &token.Position{}, // Default position
		})
	}

	// Handle return values
	var returns []rstypes.Type
	hasError := false
	for i := 0; i < t.NumOut(); i++ {
		outType := t.Out(i)
		if outType.Implements(errorType) {
			hasError = true
			continue
		}
		returns = append(returns, c.Convert(outType))
	}

	// Set return type
	if hasError {
		if len(returns) == 0 {
			fn.Returns = &rstypes.Result{
				Ok:  &rstypes.Unit{},
				Err: &rstypes.String{},
			}
		} else {
			fn.Returns = &rstypes.Result{
				Ok:  returns[0],
				Err: &rstypes.String{},
			}
		}
	} else if len(returns) == 0 {
		fn.Returns = &rstypes.Unit{}
	} else if len(returns) == 1 {
		fn.Returns = returns[0]
	} else {
		// Multiple return values become a tuple
		fn.Returns = &rstypes.Tuple{Types: returns}
	}

	return fn
}

// getParamName generates a parameter name for a given type
func (c *Converter) getParamName(t reflect.Type, index int, configNames []string) string {
	// Use configured name if available
	if configNames != nil && index < len(configNames) {
		return configNames[index]
	}

	// Check if we already have a name for this type
	if name, ok := c.paramNames[t]; ok {
		return name
	}

	// Generate name based on type
	switch t.Kind() {
	case reflect.Ptr:
		return c.getParamName(t.Elem(), index, nil)
	case reflect.Struct:
		if t.Name() != "" {
			return strings.ToLower(t.Name())
		}
	}

	// Default to type name or placeholder
	if t.Name() != "" {
		return toSnakeCase(t.Name())
	}
	return fmt.Sprintf("param%d", index)
}

// toSnakeCase converts a string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}
