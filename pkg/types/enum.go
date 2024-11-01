// Package types contains structs/interfaces representing Rust types
package rstypes

import "go/token"

// EnumVariant represents a variant in an enum
type EnumVariant struct {
	RawName    string
	RawTag     string
	FieldIndex int

	Type     Type
	Position *token.Position
}

// Enum - enum in Rust
type Enum struct {
	Common
	Name string

	Variants map[string]EnumVariant
}

var _ Type = &Enum{}
var _ NamedType = &Enum{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (e *Enum) UsedAsMapKey() bool {
	return false
}

// SetName sets an alternative name
func (e *Enum) SetName(name string) {
	e.Name = name
}

// String returns this type in string representation
func (e *Enum) String() string {
	return "enum " + e.Name
}
