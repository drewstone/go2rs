// Place this file in the root directory
package go2rs

import (
	"github.com/drewstone/go2rs/pkg/generator"
	rstypes "github.com/drewstone/go2rs/pkg/types"
)

// Re-export the Generator type and constructor
type Generator = generator.Generator

func NewGenerator(types map[string]rstypes.Type) *Generator {
	return generator.NewGenerator(types)
}

// Re-export all types
type (
	// Core types
	Type       = rstypes.Type
	Enumerable = rstypes.Enumerable
	NamedType  = rstypes.NamedType
	Common     = rstypes.Common

	// Specific types
	Any       = rstypes.Any
	Array     = rstypes.Array
	Boolean   = rstypes.Boolean
	Date      = rstypes.Date
	Enum      = rstypes.Enum
	Function  = rstypes.Function
	Map       = rstypes.Map
	Number    = rstypes.Number
	Nullable  = rstypes.Nullable
	Primitive = rstypes.Primitive
	Result    = rstypes.Result
	Stream    = rstypes.Stream
	String    = rstypes.String
	Struct    = rstypes.Struct
	Trait     = rstypes.Trait
	Tuple     = rstypes.Tuple
	Unit      = rstypes.Unit
	Vec       = rstypes.Vec

	// Field types
	StructField   = rstypes.StructField
	EnumVariant   = rstypes.EnumVariant
	FunctionParam = rstypes.FunctionParam
)

// Constructor functions
func NewAny() *Any                        { return &Any{} }
func NewArray(inner Type) *Array          { return &Array{Inner: inner} }
func NewBoolean() *Boolean                { return &Boolean{} }
func NewDate() *Date                      { return &Date{} }
func NewEnum(name string) *Enum           { return &Enum{Name: name, Variants: make(map[string]EnumVariant)} }
func NewFunction() *Function              { return &Function{} }
func NewMap(key, value Type) *Map         { return &Map{Key: key, Value: value} }
func NewNumber() *Number                  { return &Number{} }
func NewNullable(inner Type) *Nullable    { return &Nullable{Inner: inner} }
func NewPrimitive(name string) *Primitive { return &Primitive{Name: name} }
func NewResult(ok, err Type) *Result      { return &Result{Ok: ok, Err: err} }
func NewStream(inner Type) *Stream        { return &Stream{Inner: inner} }
func NewString() *String                  { return &String{} }
func NewStruct(name string) *Struct {
	return &Struct{
		Name:   name,
		Fields: make(map[string]StructField),
	}
}
func NewTrait(name string) *Trait { return &Trait{Name: name} }
func NewTuple() *Tuple            { return &Tuple{} }
func NewUnit() *Unit              { return &Unit{} }
func NewVec(inner Type) *Vec      { return &Vec{Inner: inner} }
