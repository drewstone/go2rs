package testdata

import (
	types "github.com/drewstone/go2rs/pkg/types"
)

var (
	// Data02 - 02.ts
	Data02 = map[string]types.Type{
		"github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Data": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Data",
			Fields: map[string]types.StructField{
				"Hoge": {
					Type: &types.Struct{
						Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Hoge",
						Fields: map[string]types.StructField{
							"Data": {
								Type: &types.Number{},
							},
						},
					},
				},
				"PkgHoge": {
					Type: &types.Struct{
						Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict/pkg.Hoge",
						Fields: map[string]types.StructField{
							"Data": {
								Type: &types.Number{},
							},
						},
					},
				},
			},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Hoge": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Hoge",
			Fields: map[string]types.StructField{
				"Data": {
					Type: &types.Number{},
				},
			},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/conflict/pkg.Hoge": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict/pkg.Hoge",
			Fields: map[string]types.StructField{
				"Data": {
					Type: &types.Number{},
				},
			},
		},
	}
)
