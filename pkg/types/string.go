// Package types contains structs/interfaces representing Rust types
package rstypes

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// RawStringEnumCandidate represents a raw candidate for string enum
type RawStringEnumCandidate struct {
	Key   string
	Value string
}

// String represents a Rust String type, which can optionally be an enum
type String struct {
	Common
	Name    string
	Enum    []string                 // Possible enum variants
	RawEnum []RawStringEnumCandidate // Raw enum data with keys
}

var _ Type = &String{}
var _ NamedType = &String{}
var _ Enumerable = &String{}

// UsedAsMapKey returns whether this type can be used as a map key.
// Only non-enum strings can be used as map keys.
func (s *String) UsedAsMapKey() bool {
	return len(s.Enum) == 0
}

// AddCandidates adds a candidate variant to the string enum
func (s *String) AddCandidates(key string, value interface{}) {
	str, ok := value.(string)
	if !ok {
		panic(fmt.Sprintf("expected string for enum variant, got: %s", reflect.TypeOf(value)))
	}

	s.Enum = append(s.Enum, str)
	s.RawEnum = append(s.RawEnum, RawStringEnumCandidate{
		Key:   key,
		Value: str,
	})
}

// SetName sets the type name for this string
func (s *String) SetName(name string) {
	s.Name = name
}

// String returns a string representation of this type
func (s *String) String() string {
	var buf bytes.Buffer
	buf.WriteString("String")

	parts := make([]string, 0, 2)
	if s.Name != "" {
		parts = append(parts, s.Name)
	}
	if len(s.Enum) > 0 {
		parts = append(parts, "["+strings.Join(s.Enum, ",")+"]")
	}

	if len(parts) > 0 {
		buf.WriteString("(")
		buf.WriteString(strings.Join(parts, ", "))
		buf.WriteString(")")
	}

	return buf.String()
}
