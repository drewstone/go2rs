package testdata

import (
	types "github.com/drewstone/go2rs/pkg/types"
)

var (
	// Data03 - 03.ts
	Data03 = map[string]types.Type{
		"github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive",
			Fields: map[string]types.StructField{
				"Re": {}, // Overwritten by init()
				"Children": {
					Type: &types.Array{
						Inner: &types.Struct{
							Name: "github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive",
						},
					},
				},
			},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/recursive.RecursiveMap": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/recursive.RecursiveMap",
			Fields: map[string]types.StructField{
				"Map": {
					Type: &types.Map{
						Key: &types.String{},
						Value: &types.Struct{
							Name: "github.com/drewstone/go2rs/pkg/parser/testdata/recursive.RecursiveMap",
						},
					},
				},
			},
		},
	}
)

func init() {
	//nolint
	re := Data03["github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive"].(*types.Struct)

	re.Fields["Re"] = types.StructField{
		Type: &types.Nullable{
			Inner: re,
		},
	}

}
