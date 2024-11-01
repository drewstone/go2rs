package testdata

import (
	types "github.com/drewstone/go2rs/pkg/types"
)

var (
	// Data01 - 01.ts
	Data01 = map[string]types.Type{
		"github.com/drewstone/go2rs/pkg/parser/testdata/success.Embedded": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/success.Embedded",
			Fields: map[string]types.StructField{
				"foo": {
					Optional: true,
					Type:     &types.Number{},
				},
			},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/success.Status": &types.String{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/success.Status",
			Enum: []string{"Failure", "OK"},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/success.Data": &types.Struct{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/success.Data",
			Fields: map[string]types.StructField{
				"Time": {
					Type: &types.Date{},
				},
				"Package": {
					Type: &types.Nullable{
						Inner: &types.Struct{
							Fields: map[string]types.StructField{
								"data": {
									Type: &types.Number{},
								},
							},
						},
					},
				},
				"foo": {
					Optional: true,
					Type:     &types.Number{},
				},
				"A": {
					Type: &types.Number{},
				},
				"b": {
					Optional: true,
					Type:     &types.Number{},
				},
				"C": {
					Type: &types.String{},
				},
				"D": {
					Type: &types.Nullable{
						Inner: &types.Number{},
					},
				},
				"EnumArray": {
					Type: &types.Array{
						Inner: &types.String{
							Enum: []string{"a", "b", "c"},
						},
					},
				},
				"Array": {
					Type: &types.Nullable{
						Inner: &types.Array{
							Inner: &types.Number{},
						},
					},
				},
				"Map": {
					Type: &types.Map{
						Key: &types.String{},
						Value: &types.String{
							Name: "github.com/drewstone/go2rs/pkg/parser/testdata/success.Status",
							Enum: []string{"Failure", "OK"},
						},
					},
				},
				"OptionalArray": {
					Type: &types.Array{
						Inner: &types.Nullable{
							Inner: &types.String{},
						},
					},
				},
				"Status": {
					Type: &types.String{
						Name: "github.com/drewstone/go2rs/pkg/parser/testdata/success.Status",
						Enum: []string{"Failure", "OK"},
					},
				},
				"Foo": {
					Optional: true,
					Type: &types.Struct{
						Fields: map[string]types.StructField{
							"V": {
								Type: &types.Number{},
							},
						},
					},
				},
				"U": {
					Type: &types.Struct{
						Fields: map[string]types.StructField{
							"Data": {
								Type: &types.Number{},
							},
						},
					},
				},
			},
		},
	}
)
