// Package types contains structs/interfaces representing Rust types
package rstypes

// Unit - unit type () in Rust
type Unit struct {
	Common
}

var _ Type = &Unit{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (u *Unit) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (u *Unit) String() string {
	return "()"
}
