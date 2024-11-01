package testdata

import types "github.com/drewstone/go2rs/pkg/types"

var (
	// Test05 - 05.ts
	Test05 = map[string]types.Type{
		"github.com/drewstone/go2rs/pkg/parser/testdata.CustomTest": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata.CustomTest",
			Fields: map[string]types.StructField{
				"C": {
					RawName:    "C",
					FieldIndex: 0,

					Type: &types.Struct{
						Name: "github.com/drewstone/go2rs/pkg/parser/testdata.CustomTestC",
					},
				},
			},
		},
	}
)
