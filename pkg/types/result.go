// Package types contains structs/interfaces representing Rust types
package rstypes

import "fmt"

// Result - Result<T, E> in Rust
type Result struct {
	Common
	Ok  Type
	Err Type
}

var _ Type = &Result{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (e *Result) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (e *Result) String() string {
	return fmt.Sprintf("Result<%s, %s>", e.Ok.String(), e.Err.String())
}
