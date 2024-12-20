// Package types contains structs/interfaces representing Rust types
package rstypes

// Date - RFC3399 string in Rust
type Date struct {
	Common
}

var _ Type = &Date{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (*Date) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (*Date) String() string {
	return "Date"
}
