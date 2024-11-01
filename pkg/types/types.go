// Package types contains structs/interfaces representing Rust types
package rstypes

import "go/token"

// Common defines common fields in the all types
type Common struct {
	// PkgName is the package name declared at the beginning of .go files.
	// Currently, only exported types in the root package is available.
	PkgName  string
	Position *token.Position
}

// SetPackageName sets PkgName in Common
func (c *Common) SetPackageName(pkgName string) {
	c.PkgName = pkgName
}

// GetPackageName returns PkgName in Common
func (c *Common) GetPackageName() string {
	return c.PkgName
}

// SetPosition sets Position in Common
func (c *Common) SetPosition(pos *token.Position) {
	c.Position = pos
}

// GetPosition returns Position in Common
func (c *Common) GetPosition() *token.Position {
	return c.Position
}

// Type interface represents all Rust types handled by go-easyparser
type Type interface {
	SetPackageName(pkgName string)
	GetPackageName() string
	UsedAsMapKey() bool
	String() string
	SetPosition(pos *token.Position)
	GetPosition() *token.Position
}

// Enumerable interface represents union types
type Enumerable interface {
	Type

	// AddCandidates adds a candidate for enum
	AddCandidates(key string, v interface{})
}

// NamedType interface represents named types
type NamedType interface {
	Type

	// SetName sets an alternative name
	SetName(name string)
}
