// Package types contains structs/interfaces representing Rust types
package rstypes

import "fmt"

// Array represents a fixed-size array type in Rust ([T; N])
type Array struct {
	Common
	Inner Type
	Size  uint64 // Size of the fixed array
}

var _ Type = &Array{}

// UsedAsMapKey returns whether this type can be used as a map key.
// Arrays cannot be used as map keys since they don't implement Hash/Eq.
func (a *Array) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (a *Array) String() string {
	return "[" + a.Inner.String() + "; " + fmt.Sprint(a.Size) + "]"
}
