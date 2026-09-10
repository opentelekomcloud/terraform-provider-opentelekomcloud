package lts

import (
	"reflect"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/groups"
)

func TestDataSourceLtsGroupsV2Filters(t *testing.T) {
	dataSource := DataSourceLtsGroupsV2()
	for _, field := range []string{"group_id", "name"} {
		if !dataSource.Schema[field].Optional {
			t.Fatalf("%s must be an optional filter", field)
		}
	}
}

func TestFilterLtsGroups(t *testing.T) {
	allGroups := []groups.LogGroup{
		{LogGroupId: "first-id", LogGroupName: "first"},
		{LogGroupId: "second-id", LogGroupName: "second"},
	}

	tests := []struct {
		name     string
		id       string
		group    string
		expected []groups.LogGroup
	}{
		{name: "no filters", expected: allGroups},
		{name: "by ID", id: "first-id", expected: allGroups[:1]},
		{name: "by name", group: "second", expected: allGroups[1:]},
		{name: "combined match", id: "second-id", group: "second", expected: allGroups[1:]},
		{name: "combined mismatch", id: "first-id", group: "second", expected: []groups.LogGroup{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := filterLtsGroups(allGroups, test.id, test.group)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Fatalf("unexpected groups: %#v", actual)
			}
		})
	}
}
