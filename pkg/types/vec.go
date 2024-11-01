// Package types contains structs/interfaces representing Rust types
package rstypes

// Vec represents a dynamic array type in Rust (Vec<T>)
type Vec struct {
	Common
	Inner Type
}

var _ Type = &Vec{}

// UsedAsMapKey returns whether this type can be used as a map key.
// Vecs cannot be used as map keys since they don't implement Hash/Eq.
func (v *Vec) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (v *Vec) String() string {
	return "Vec<" + v.Inner.String() + ">"
}
