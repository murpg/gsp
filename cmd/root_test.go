package cmd

import (
	"reflect"
	"sort"
	"testing"
)

var testFixtures = []struct {
	name string
	in   []string
	out  []string
}{
	{
		"sliceWithDuplicates",
		[]string{
			"cfcs/database/fileExchange/fe.cfc",
			"pmses/common/fileExchange/feRemoveItem.cfm",
			"pmses/common/fileExchange/feRemoveUser.cfm",
			"pmses/common/fileExchange/js/feCommon.js",
			"pmses/common/fileExchange/js/feCommon.js",
		},
		[]string{
			"cfcs/database/fileExchange/fe.cfc",
			"pmses/common/fileExchange/feRemoveItem.cfm",
			"pmses/common/fileExchange/feRemoveUser.cfm",
			"pmses/common/fileExchange/js/feCommon.js",
		},
	},
	{
		"sliceWithoutDuplicates",
		[]string{
			"cfcs/database/fileExchange/fe.cfc",
			"pmses/common/fileExchange/feRemoveItem.cfm",
		},
		[]string{
			"cfcs/database/fileExchange/fe.cfc",
			"pmses/common/fileExchange/feRemoveItem.cfm",
		},
	},
}

func TestRemoveDuplicates(t *testing.T) {
	for _, tt := range testFixtures {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicates(tt.in)
			if len(got) != len(tt.out) {
				t.Errorf("removeDuplicates(%s) got %v, want %v", tt.in, got, tt.out)
			}

			gotCopy := make([]string, len(got))
			copy(gotCopy, got)
			sort.Strings(gotCopy)

			ttCopy := make([]string, len(tt.out))
			copy(ttCopy, tt.out)
			sort.Strings(ttCopy)

			if !reflect.DeepEqual(gotCopy, ttCopy) {
				t.Errorf("got %v, want %v", gotCopy, ttCopy)
			}
		})
	}
}
