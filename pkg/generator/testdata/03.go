package testdata

import (
	tstypes "github.com/drewstone/go2rs/pkg/types"
)

var (
	// Data03 - 03.ts
	Data03 = map[string]tstypes.Type{
		"github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive": &tstypes.Object{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive",
			Entries: map[string]tstypes.ObjectEntry{
				"Re": {}, // Overwritten by init()
				"Children": {
					Type: &tstypes.Array{
						Inner: &tstypes.Object{
							Name: "github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive",
						},
					},
				},
			},
		},
		"github.com/drewstone/go2rs/pkg/parser/testdata/recursive.RecursiveMap": &tstypes.Object{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata/recursive.RecursiveMap",
			Entries: map[string]tstypes.ObjectEntry{
				"Map": {
					Type: &tstypes.Map{
						Key: &tstypes.String{},
						Value: &tstypes.Object{
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
	re := Data03["github.com/drewstone/go2rs/pkg/parser/testdata/recursive.Recursive"].(*tstypes.Object)

	re.Entries["Re"] = tstypes.ObjectEntry{
		Type: &tstypes.Nullable{
			Inner: re,
		},
	}

}
