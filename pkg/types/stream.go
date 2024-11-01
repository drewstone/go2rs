// Package types contains structs/interfaces representing Rust types
package rstypes

// Stream - stream in Rust
type Stream struct {
	Common
	Name  string
	Inner Type
}

var _ Type = &Stream{}
var _ NamedType = &Stream{}

// UsedAsMapKey returns whether this type can be used as the key for map
func (s *Stream) UsedAsMapKey() bool {
	return false
}

// SetName sets an alternative name
func (s *Stream) SetName(name string) {
	s.Name = name
}

// String returns this type in string representation
func (s *Stream) String() string {
	return "Stream<" + s.Inner.String() + ">"
}
