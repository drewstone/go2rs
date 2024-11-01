// Package types contains structs/interfaces representing Rust types
package rstypes

import "go/token"

// StructField is a field in structs
type StructField struct {
	RawName    string
	RawTag     string
	FieldIndex int

	Type     Type
	Position *token.Position
	Optional bool
}

// Struct - struct in Rust
type Struct struct {
	Common
	Name string

	Fields map[string]StructField
}

var _ Type = &Struct{}
var _ NamedType = &Struct{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (n *Struct) UsedAsMapKey() bool {
	return false
}

// SetName sets an alternative name
func (n *Struct) SetName(name string) {
	n.Name = name
}

// String returns this type in string representation
func (n *Struct) String() string {
	return n.Name
}
