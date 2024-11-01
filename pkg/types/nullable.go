// Package types contains structs/interfaces representing Rust types
package rstypes

// Nullable - ... i.e. Option in Rust
type Nullable struct {
	Common
	Inner Type
}

var _ Type = &Nullable{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (e *Nullable) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (e *Nullable) String() string {
	return "Option<" + e.Inner.String() + ">"
}
