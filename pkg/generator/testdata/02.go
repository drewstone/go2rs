package testdata

import (
	tstypes "github.com/drewstone/go2rs/pkg/types"
)

var (
	// Data02 - 02.ts
	Data02 = map[string]tstypes.Type{
		"github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Data": &tstypes.Object{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Data",
			Entries: map[string]tstypes.ObjectEntry{
				"Hoge": {
					Type: &tstypes.Object{
						Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Hoge",
						Entries: map[string]tstypes.ObjectEntry{
							"Data": {
								Type: &tstypes.Number{},
							},
						},
					},
				},
				"PkgHoge": {
					Type: &tstypes.Object{
						Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict/pkg.Hoge",
						Entries: map[string]tstypes.ObjectEntry{
							"Data": {
								Type: &tstypes.Number{},
							},
						},
					},
				},
			},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Hoge": &tstypes.Object{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict.Hoge",
			Entries: map[string]tstypes.ObjectEntry{
				"Data": {
					Type: &tstypes.Number{},
				},
			},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/conflict/pkg.Hoge": &tstypes.Object{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/conflict/pkg.Hoge",
			Entries: map[string]tstypes.ObjectEntry{
				"Data": {
					Type: &tstypes.Number{},
				},
			},
		},
	}
)
