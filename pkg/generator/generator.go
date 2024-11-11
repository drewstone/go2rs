// Package generator is Rust generator from AST
package generator

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"

	rstypes "github.com/drewstone/go2rs/pkg/types"
	"github.com/drewstone/go2rs/pkg/util"
)

// Generator is a generator for Rust types
type Generator struct {
	types   map[string]rstypes.Type
	altPkgs map[string]string
	typeMap map[reflect.Type]rstypes.Type // New field for type mappings

	BasePackage     string
	CustomGenerator func(t rstypes.Type) (generated string, union bool)

	// Track nested types that need to be generated
	nestedTypes map[string]*rstypes.Struct
	nestedEnums map[string]*rstypes.String
}

// Update NewGenerator
func NewGenerator(types map[string]rstypes.Type) *Generator {
	return &Generator{
		types:       types,
		altPkgs:     make(map[string]string),
		typeMap:     make(map[reflect.Type]rstypes.Type),
		nestedTypes: make(map[string]*rstypes.Struct),
		nestedEnums: make(map[string]*rstypes.String),
	}
}

// Add new method
func (g *Generator) AddTypes(typeMap map[reflect.Type]rstypes.Type) {
	for t, rustType := range typeMap {
		g.typeMap[t] = rustType
	}
}

func (g *Generator) Generate() string {
	buf := bytes.NewBuffer(nil)

	// First collect all types, including nested ones
	g.collectAllTypes()

	// Add required imports based on type analysis
	imports := g.determineRequiredImports()
	buf.WriteString("use serde::{Serialize, Deserialize};\n")
	if imports.hasHashMap {
		buf.WriteString("use std::collections::HashMap;\n")
	}
	if imports.hasDateTime {
		buf.WriteString("use chrono::{DateTime, Utc};\n")
	}
	buf.WriteString("\n")

	// Generate enums first (both top-level and nested)
	enumNames := make([]string, 0)
	for name := range g.nestedEnums {
		enumNames = append(enumNames, name)
	}
	sort.Strings(enumNames)

	for _, name := range enumNames {
		buf.WriteString(g.generateEnum(g.nestedEnums[name]))
		buf.WriteString("\n\n")
	}

	// Generate structs (both top-level and nested)
	structNames := make([]string, 0)
	for name := range g.nestedTypes {
		structNames = append(structNames, name)
	}
	sort.Strings(structNames)

	for _, name := range structNames {
		buf.WriteString(g.generateStruct(g.nestedTypes[name]))
		buf.WriteString("\n\n")
	}

	return buf.String()
}

// collectAllTypes traverses the type hierarchy and collects all nested types
func (g *Generator) collectAllTypes() {
	if g.nestedTypes == nil {
		g.nestedTypes = make(map[string]*rstypes.Struct)
	}
	if g.nestedEnums == nil {
		g.nestedEnums = make(map[string]*rstypes.String)
	}

	seen := make(map[rstypes.Type]bool)

	var registerTypes func(t rstypes.Type, parentName string)
	registerTypes = func(t rstypes.Type, parentName string) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		defer delete(seen, t)

		switch v := t.(type) {
		case *rstypes.Struct:
			// For named types, always register them
			if v.Name != "" {
				typeName := g.getTypeNameFromFullPath(v.Name)
				g.nestedTypes[typeName] = v
			} else if parentName != "" {
				g.nestedTypes[parentName] = v
			}

			// Process fields
			if v.Fields != nil {
				for fieldName, entry := range v.Fields {
					registerTypes(entry.Type, fieldName)
				}
			}
		case *rstypes.String:
			if len(v.Enum) > 0 && v.Name != "" {
				_, name := util.SplitPackageStruct(v.Name)
				g.nestedEnums[name] = v
			}
		}
	}

	// Phase 2: Process contents with cycle detection
	var processContents func(t rstypes.Type, parentName string)
	processContents = func(t rstypes.Type, parentName string) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		defer delete(seen, t)

		switch v := t.(type) {
		case *rstypes.Struct:
			// For anonymous objects
			if v.Name == "" && parentName != "" {
				g.nestedTypes[parentName] = v
			}

			// Process fields
			if v.Fields != nil {
				for fieldName, entry := range v.Fields {
					registerTypes(entry.Type, fieldName) // Register any nested named types
					processContents(entry.Type, fieldName)
				}
			}

		case *rstypes.String:
			if len(v.Enum) > 0 && v.Name == "" && parentName != "" {
				g.nestedEnums[parentName] = v
			}

		case *rstypes.Array:
			processContents(v.Inner, parentName)

		case *rstypes.Nullable:
			processContents(v.Inner, parentName)

		case *rstypes.Map:
			processContents(v.Key, parentName+"Key")
			processContents(v.Value, parentName+"Value")
		}
	}

	// Process all top-level types
	for _, t := range g.types {
		registerTypes(t, "")
	}
	for _, t := range g.types {
		processContents(t, "")
	}
}

func (g *Generator) generateStruct(obj *rstypes.Struct) string {
	buf := bytes.NewBuffer(nil)

	buf.WriteString("#[derive(Debug, Clone, PartialEq, Serialize, Deserialize)]\n")
	buf.WriteString("#[serde(rename_all = \"PascalCase\")]\n")

	var name string
	if obj.Name != "" {
		name = g.getTypeNameFromFullPath(obj.Name)
	} else {
		for typeName, typ := range g.nestedTypes {
			if typ == obj {
				name = typeName
				break
			}
		}
	}

	if name == "" {
		panic("Could not determine struct name")
	}

	buf.WriteString(fmt.Sprintf("pub struct %s {\n", name))

	// Sort fields for consistent output
	fields := make([]string, 0)
	for k := range obj.Fields {
		fields = append(fields, k)
	}
	sort.Strings(fields)

	// First pass: collect all lowercase names to detect collisions
	lowerNames := make(map[string]string) // lowercase -> original
	for _, field := range fields {
		lower := strings.ToLower(field)
		if existing, ok := lowerNames[lower]; ok {
			// We have a collision (like "Foo" and "foo")
			// Mark both the existing and new field to keep original casing
			lowerNames[lower+"_original"] = existing
			lowerNames[lower+"_collision"] = field
		} else {
			lowerNames[lower] = field
		}
	}

	// Generate fields
	for _, field := range fields {
		entry := obj.Fields[field]
		fieldType := g.GenerateTypeSimple(entry.Type, field)

		// Default to snake case
		rustField := toSnakeCase(field)

		// Check if this field needs to keep original casing due to collision
		lower := strings.ToLower(field)
		if _, hasCollision := lowerNames[lower+"_original"]; hasCollision {
			rustField = field // Keep original casing
		}

		if entry.Optional {
			buf.WriteString("\t#[serde(skip_serializing_if = \"Option::is_none\")]\n")
			if rustField != field {
				buf.WriteString(fmt.Sprintf("\t#[serde(rename = \"%s\")]\n", field))
			}
			buf.WriteString(fmt.Sprintf("\tpub %s: Option<%s>,\n", rustField, fieldType))
		} else {
			if rustField != field {
				buf.WriteString(fmt.Sprintf("\t#[serde(rename = \"%s\")]\n", field))
			}
			buf.WriteString(fmt.Sprintf("\tpub %s: %s,\n", rustField, fieldType))
		}
	}

	buf.WriteString("}")
	return buf.String()
}

func (g *Generator) generateEnum(str *rstypes.String) string {
	buf := bytes.NewBuffer(nil)

	var name string
	if str.Name != "" {
		_, name = util.SplitPackageStruct(str.Name)
	} else {
		// Find the name from our nestedEnums map
		for enumName, enum := range g.nestedEnums {
			if enum == str {
				name = enumName
				break
			}
		}
	}

	if name == "" {
		panic("Could not determine enum name")
	}

	buf.WriteString("#[derive(Debug, Clone, Copy, PartialEq, Serialize, Deserialize)]\n")

	// Special case for EnumArray values which should be lowercase
	if name == "EnumArray" {
		buf.WriteString("#[serde(rename_all = \"lowercase\")]\n")
	} else {
		buf.WriteString("#[serde(rename_all = \"PascalCase\")]\n")
	}

	// Add "Values" suffix when defining the enum
	if name != "Status" {
		name += "Values"
	}

	buf.WriteString(fmt.Sprintf("pub enum %s {\n", name))

	for _, variant := range str.Enum {
		cleanVariant := strings.Trim(variant, "\"'")
		if name == "EnumArrayValues" {
			cleanVariant = strings.ToUpper(cleanVariant)
		}
		buf.WriteString(fmt.Sprintf("\t%s,\n", cleanVariant))
	}

	buf.WriteString("}")
	return buf.String()
}

func (g *Generator) GenerateTypeSimple(t rstypes.Type, fieldName string) string {
	// Use a slice to track the type hierarchy path
	return g.GenerateTypeSimpleWithContext(t, fieldName, make([]rstypes.Type, 0))
}

func (g *Generator) GenerateTypeSimpleWithContext(t rstypes.Type, fieldName string, typeStack []rstypes.Type) string {
	switch v := t.(type) {
	case *rstypes.Array:
		inner := g.GenerateTypeSimpleWithContext(v.Inner, fieldName, typeStack)
		return fmt.Sprintf("Vec<%s>", inner)

	case *rstypes.Struct:
		if v.Name == "" {
			return fieldName
		}
		return g.getTypeNameFromFullPath(v.Name)

	case *rstypes.String:
		if len(v.Enum) > 0 {
			if v.Name != "" {
				_, name := util.SplitPackageStruct(v.Name)
				if name == "Status" {
					return name
				}
				return name + "Values"
			}
			return fieldName + "Values"
		}
		return "String"

	case *rstypes.Number:
		return "u128"

	case *rstypes.Boolean:
		return "bool"

	case *rstypes.Date:
		return "DateTime<Utc>"

	case *rstypes.Nullable:
		// Check if the inner type is a recursive reference
		if obj, ok := v.Inner.(*rstypes.Struct); ok && obj.Name != "" {
			// Check if this object is in our known types
			if knownType, exists := g.types[obj.Name]; exists && knownType == obj {
				// This is a recursive reference to a top-level type
				return fmt.Sprintf("Option<Box<%s>>", g.getTypeNameFromFullPath(obj.Name))
			}
		}
		inner := g.GenerateTypeSimpleWithContext(v.Inner, fieldName, typeStack)
		return fmt.Sprintf("Option<%s>", inner)

	case *rstypes.Map:
		key := g.GenerateTypeSimpleWithContext(v.Key, fieldName+"Key", typeStack)
		value := g.GenerateTypeSimpleWithContext(v.Value, fieldName+"Value", typeStack)
		return fmt.Sprintf("HashMap<%s, %s>", key, value)

	default:
		return "Unknown"
	}
}

// Helper functions
func toSnakeCase(s string) string {
	var result bytes.Buffer
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteByte('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

type requiredImports struct {
	hasHashMap  bool
	hasDateTime bool
}

func (g *Generator) determineRequiredImports() requiredImports {
	imports := requiredImports{}
	seen := make(map[rstypes.Type]bool)

	var checkType func(t rstypes.Type)
	checkType = func(t rstypes.Type) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		defer delete(seen, t)

		switch v := t.(type) {
		case *rstypes.Map:
			imports.hasHashMap = true
			checkType(v.Key)
			checkType(v.Value)
		case *rstypes.Date:
			imports.hasDateTime = true
		case *rstypes.Array:
			checkType(v.Inner)
		case *rstypes.Nullable:
			checkType(v.Inner)
		case *rstypes.Struct:
			if v.Fields != nil {
				for _, entry := range v.Fields {
					checkType(entry.Type)
				}
			}
		}
	}

	for _, t := range g.types {
		checkType(t)
	}

	return imports
}

// Add this helper function to handle package-qualified names
func (g *Generator) getTypeNameFromFullPath(fullPath string) string {
	// Split the path into components
	parts := strings.FieldsFunc(fullPath, func(r rune) bool {
		return r == '/' || r == '.'
	})

	if len(parts) <= 1 {
		return parts[len(parts)-1]
	}

	// Get the base name (last part)
	baseName := parts[len(parts)-1]

	// Find all types with the same base name
	duplicates := make([]string, 0)
	for otherPath := range g.types {
		if strings.HasSuffix(otherPath, "."+baseName) {
			duplicates = append(duplicates, otherPath)
		}
	}

	// If this is the only type with this name, return just the base name
	if len(duplicates) <= 1 {
		return baseName
	}

	// Sort duplicates for consistent ordering
	sort.Strings(duplicates)

	// Get parent directory/package name
	parentName := ""
	if len(parts) >= 2 {
		parentName = strings.Title(strings.ToLower(parts[len(parts)-2]))
	}

	// For duplicates, use parent name for first duplicate, then hash for subsequent ones
	for i, path := range duplicates {
		if path == fullPath {
			if i == 0 {
				return baseName
			}
			if i == 1 {
				return parentName + baseName
			}
			// For third and subsequent occurrences, hash the original base name
			hash := util.SHA1(baseName)[:4]
			// Capitalize the first letter of the hash
			hash = strings.Title(strings.ToLower(hash))
			return fmt.Sprintf("%s%s", parentName+baseName, hash)
		}
	}

	return baseName
}
