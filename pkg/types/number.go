// Package types contains structs/interfaces representing Rust types
package rstypes

import (
	"fmt"
	"go/types"
	"reflect"
)

// RawNumberEnumCandidate represents a raw candidate for number enum
type RawNumberEnumCandidate struct {
	Key   string
	Value interface{}
}

// Number - numeric types in Rust
type Number struct {
	Common
	Name    string
	RawType types.BasicKind

	// For enum variants
	Enum    []int64
	RawEnum []RawNumberEnumCandidate

	// Rust numeric type
	IsFloat   bool
	IsSigned  bool
	BitSize   int
	IsUnsized bool // For isize/usize
}

var _ Type = &Number{}
var _ NamedType = &Number{}
var _ Enumerable = &Number{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (e *Number) UsedAsMapKey() bool {
	return len(e.Enum) == 0
}

// AddCandidates adds an candidate for enum
func (e *Number) AddCandidates(key string, v interface{}) {
	val := reflect.ValueOf(v)
	var int64Val int64

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		int64Val = val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		int64Val = int64(val.Uint())
	case reflect.Float32, reflect.Float64:
		int64Val = int64(val.Float())
	default:
		panic(fmt.Sprintf("expected numeric type for enum variant, got: %s", val.Type()))
	}

	e.Enum = append(e.Enum, int64Val)
	e.RawEnum = append(e.RawEnum, RawNumberEnumCandidate{
		Key:   key,
		Value: v,
	})
}

// SetName sets an alternative name
func (e *Number) SetName(name string) {
	e.Name = name
}

// String returns this type in string representation
func (e *Number) String() string {
	if e.IsFloat {
		return fmt.Sprintf("f%d", e.BitSize)
	}

	if e.IsUnsized {
		if e.IsSigned {
			return "isize"
		}
		return "usize"
	}

	if e.IsSigned {
		return fmt.Sprintf("i%d", e.BitSize)
	}
	return fmt.Sprintf("u%d", e.BitSize)
}
