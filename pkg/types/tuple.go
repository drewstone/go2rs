// Package types contains structs/interfaces representing Rust types
package rstypes

import (
	"fmt"
	"strings"
)

// Tuple - tuple in Rust
type Tuple struct {
	Common
	Name  string
	Types []Type
}

var _ Type = &Tuple{}
var _ NamedType = &Tuple{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (t *Tuple) UsedAsMapKey() bool {
	return false
}

// SetName sets an alternative name
func (t *Tuple) SetName(name string) {
	t.Name = name
}

// String returns this type in string representation
func (t *Tuple) String() string {
	if len(t.Types) == 0 {
		return "()"
	}

	types := make([]string, len(t.Types))
	for i, typ := range t.Types {
		types[i] = typ.String()
	}
	return fmt.Sprintf("(%s)", strings.Join(types, ", "))
}
