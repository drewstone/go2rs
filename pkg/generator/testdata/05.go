package testdata

import tstypes "github.com/drewstone/go2rs/pkg/types"

var (
	// Test05 - 05.ts
	Test05 = map[string]tstypes.Type{
		"github.com/drewstone/go2rs/pkg/parser/testdata.CustomTest": &tstypes.Object{
			Name: "github.com/drewstone/go2rs/pkg/parser/testdata.CustomTest",
			Entries: map[string]tstypes.ObjectEntry{
				"C": {
					RawName:    "C",
					FieldIndex: 0,

					Type: &tstypes.Object{
						Name: "github.com/drewstone/go2rs/pkg/parser/testdata.CustomTestC",
					},
				},
			},
		},
	}
)
