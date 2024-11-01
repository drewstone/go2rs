// Package types contains structs/interfaces representing Rust types
package rstypes

import (
	"fmt"
	"go/token"
	"strings"
)

// FunctionParam represents a parameter in a function
type FunctionParam struct {
	Name     string
	Type     Type
	Position *token.Position
}

// Function represents a function or method in Rust
type Function struct {
	Common
	Name     string
	IsMethod bool
	IsAsync  bool
	Params   []FunctionParam // Changed from map to slice to preserve order
	Returns  Type            // Changed from ReturnVal for consistency
	Receiver Type            // For methods, represents self type
}

var _ Type = &Function{}
var _ NamedType = &Function{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (f *Function) UsedAsMapKey() bool {
	return false
}

// SetName sets an alternative name
func (f *Function) SetName(name string) {
	f.Name = name
}

func (f *Function) String() string {
	var parts []string

	if f.IsAsync {
		parts = append(parts, "async")
	}

	parts = append(parts, "fn")

	if f.Name != "" {
		parts = append(parts, f.Name)
	}

	// Format parameters
	paramStrs := make([]string, len(f.Params))
	for i, param := range f.Params {
		paramStrs[i] = fmt.Sprintf("%s: %s", param.Name, param.Type.String())
	}

	parts = append(parts, "("+strings.Join(paramStrs, ", ")+")")

	// Add return type if not unit
	if f.Returns != nil && f.Returns.String() != "()" {
		parts = append(parts, "->", f.Returns.String())
	}

	return strings.Join(parts, " ")
}
