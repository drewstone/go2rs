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
