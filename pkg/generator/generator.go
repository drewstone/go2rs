// Package generator is Rust generator from AST
package generator

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"unicode"

	tstypes "github.com/drewstone/go2rs/pkg/types"
	"github.com/drewstone/go2rs/pkg/util"
)

// Generator is a generator for Rust types
type Generator struct {
	types   map[string]tstypes.Type
	altPkgs map[string]string

	BasePackage     string
	CustomGenerator func(t tstypes.Type) (generated string, union bool)

	// Track nested types that need to be generated
	nestedTypes map[string]*tstypes.Object
	nestedEnums map[string]*tstypes.String
}

// NewGenerator returns a new Generator
func NewGenerator(types map[string]tstypes.Type) *Generator {
	return &Generator{
		types:       types,
		altPkgs:     make(map[string]string),
		nestedTypes: make(map[string]*tstypes.Object),
		nestedEnums: make(map[string]*tstypes.String),
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
		g.nestedTypes = make(map[string]*tstypes.Object)
	}
	if g.nestedEnums == nil {
		g.nestedEnums = make(map[string]*tstypes.String)
	}

	seen := make(map[tstypes.Type]bool)

	var registerTypes func(t tstypes.Type, parentName string)
	registerTypes = func(t tstypes.Type, parentName string) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		defer delete(seen, t)

		switch v := t.(type) {
		case *tstypes.Object:
			// For named types, always register them
			if v.Name != "" {
				typeName := g.getTypeNameFromFullPath(v.Name)
				g.nestedTypes[typeName] = v
			} else if parentName != "" {
				g.nestedTypes[parentName] = v
			}

			// Process fields
			if v.Entries != nil {
				for fieldName, entry := range v.Entries {
					registerTypes(entry.Type, fieldName)
				}
			}
		case *tstypes.String:
			if len(v.Enum) > 0 && v.Name != "" {
				_, name := util.SplitPackageStruct(v.Name)
				g.nestedEnums[name] = v
			}
		}
	}

	// Phase 2: Process contents with cycle detection
	var processContents func(t tstypes.Type, parentName string)
	processContents = func(t tstypes.Type, parentName string) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		defer delete(seen, t)

		switch v := t.(type) {
		case *tstypes.Object:
			// For anonymous objects
			if v.Name == "" && parentName != "" {
				g.nestedTypes[parentName] = v
			}

			// Process fields
			if v.Entries != nil {
				for fieldName, entry := range v.Entries {
					registerTypes(entry.Type, fieldName) // Register any nested named types
					processContents(entry.Type, fieldName)
				}
			}

		case *tstypes.String:
			if len(v.Enum) > 0 && v.Name == "" && parentName != "" {
				g.nestedEnums[parentName] = v
			}

		case *tstypes.Array:
			processContents(v.Inner, parentName)

		case *tstypes.Nullable:
			processContents(v.Inner, parentName)

		case *tstypes.Map:
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

func (g *Generator) generateStruct(obj *tstypes.Object) string {
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
	for k := range obj.Entries {
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
		entry := obj.Entries[field]
		fieldType := g.generateTypeSimple(entry.Type, field)

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

func (g *Generator) generateEnum(str *tstypes.String) string {
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

func (g *Generator) generateTypeSimple(t tstypes.Type, fieldName string) string {
	// Use a map to track types being processed to prevent infinite recursion
	return g.generateTypeSimpleWithContext(t, fieldName, make(map[tstypes.Type]bool))
}

func (g *Generator) generateTypeSimpleWithContext(t tstypes.Type, fieldName string, inProcess map[tstypes.Type]bool) string {
	// If we've seen this type before in the current chain, break the recursion
	if inProcess[t] {
		// For recursive types, use the field name as the type
		if obj, ok := t.(*tstypes.Object); ok && obj.Name != "" {
			_, name := util.SplitPackageStruct(obj.Name)
			return name
		}
		return fieldName
	}

	inProcess[t] = true
	defer delete(inProcess, t)

	switch v := t.(type) {
	case *tstypes.Array:
		return fmt.Sprintf("Vec<%s>", g.generateTypeSimpleWithContext(v.Inner, fieldName, inProcess))
	case *tstypes.Object:
		if v.Name == "" {
			return fieldName
		}
		return g.getTypeNameFromFullPath(v.Name)
	case *tstypes.String:
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
	case *tstypes.Number:
		return "u128"
	case *tstypes.Boolean:
		return "bool"
	case *tstypes.Date:
		return "DateTime<Utc>"
	case *tstypes.Nullable:
		return fmt.Sprintf("Option<%s>", g.generateTypeSimpleWithContext(v.Inner, fieldName, inProcess))
	case *tstypes.Map:
		return fmt.Sprintf("HashMap<%s, %s>",
			g.generateTypeSimpleWithContext(v.Key, fieldName+"Key", inProcess),
			g.generateTypeSimpleWithContext(v.Value, fieldName+"Value", inProcess))
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
	seen := make(map[tstypes.Type]bool)

	var checkType func(t tstypes.Type)
	checkType = func(t tstypes.Type) {
		if t == nil || seen[t] {
			return
		}
		seen[t] = true
		defer delete(seen, t)

		switch v := t.(type) {
		case *tstypes.Map:
			imports.hasHashMap = true
			checkType(v.Key)
			checkType(v.Value)
		case *tstypes.Date:
			imports.hasDateTime = true
		case *tstypes.Array:
			checkType(v.Inner)
		case *tstypes.Nullable:
			checkType(v.Inner)
		case *tstypes.Object:
			if v.Entries != nil {
				for _, entry := range v.Entries {
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
			return fmt.Sprintf("%s_%s", parentName+baseName, hash)
		}
	}

	return baseName
}
