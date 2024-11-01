// Package types contains structs/interfaces representing Rust types
package rstypes

// Any represents a dynamic type in Rust (dyn Any)
type Any struct {
	Common
}

var _ Type = &Any{}

// UsedAsMapKey returns whether this type can be used as the key for map.
// dyn Any cannot be used as a map key since it doesn't implement Hash/Eq.
func (a *Any) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (a *Any) String() string {
	return "dyn Any"
}
