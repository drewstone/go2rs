// Package types contains structs/interfaces representing Rust types
package rstypes

// Primitive represents primitive types in Rust
type Primitive struct {
	Common
	Name string
}

var _ Type = &Primitive{}

func (p *Primitive) UsedAsMapKey() bool {
	return true
}

func (p *Primitive) String() string {
	return p.Name
}
