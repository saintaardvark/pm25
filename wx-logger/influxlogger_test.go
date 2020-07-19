package main

import (
	"testing"
)

var TestTable = []struct {
	name string
	want string
}{
	{
		"PrcpMtr",
		"precipitation_meter",
	},
}

func TestGetAbbrev(t *testing.T) {
	for _, test := range TestTable {
		abbrev := GetAbbrev(test.name)
		if abbrev != test.want {
			t.Errorf("GetAbbrev(%v) returned (%v), want (%v)",
				test.name, abbrev, test.want)
		}
	}
}
