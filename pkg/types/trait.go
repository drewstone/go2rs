// Package types contains structs/interfaces representing Rust types
package rstypes

// Trait - trait in Rust
type Trait struct {
	Common
	Name string
}

var _ Type = &Trait{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (t *Trait) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (t *Trait) String() string {
	return "trait " + t.Name
}
