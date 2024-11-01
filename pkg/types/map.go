// Package types contains structs/interfaces representing Rust types
package rstypes

import "fmt"

// Map represents a HashMap<K,V> type in Rust
type Map struct {
	Common
	Key   Type
	Value Type
}

var _ Type = &Map{}

// UsedAsMapKey returns whether this type can be used as a map key.
// Maps cannot be used as map keys since they don't implement Hash/Eq.
func (m *Map) UsedAsMapKey() bool {
	return false
}

// String returns this type in string representation
func (m *Map) String() string {
	return fmt.Sprintf("HashMap<%s, %s>", m.Key.String(), m.Value.String())
}
