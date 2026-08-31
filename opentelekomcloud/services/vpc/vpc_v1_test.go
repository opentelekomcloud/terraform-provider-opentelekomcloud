package vpc

import (
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/vpcs"
)

func TestVpcV1Schemas(t *testing.T) {
	resourceSchema := ResourceVirtualPrivateCloudV1().Schema
	if !resourceSchema["enterprise_project_id"].Optional ||
		!resourceSchema["enterprise_project_id"].Computed ||
		!resourceSchema["enterprise_project_id"].ForceNew {
		t.Fatal("resource enterprise_project_id must be optional, computed, and force replacement")
	}
	for _, field := range []string{"routes", "tenant_id", "created_at", "updated_at"} {
		if !resourceSchema[field].Computed || resourceSchema[field].Optional {
			t.Fatalf("resource %s must be computed only", field)
		}
	}

	dataSourceSchema := DataSourceVirtualPrivateCloudVpcV1().Schema
	for _, field := range []string{"description", "routes", "tenant_id", "created_at", "updated_at"} {
		if !dataSourceSchema[field].Computed || dataSourceSchema[field].Optional {
			t.Fatalf("data source %s must be computed only", field)
		}
	}
	if !dataSourceSchema["enterprise_project_id"].Optional || !dataSourceSchema["enterprise_project_id"].Computed {
		t.Fatal("data source enterprise_project_id must be optional and computed")
	}
}

func TestFilterVpcs(t *testing.T) {
	shared := false
	input := []vpcs.Vpc{
		{ID: "first", Name: "match", CIDR: "192.168.0.0/16", Status: "OK", EnableSharedSnat: false},
		{ID: "second", Name: "match", CIDR: "192.168.0.0/16", Status: "OK", EnableSharedSnat: true},
		{ID: "third", Name: "other", CIDR: "10.0.0.0/8", Status: "CREATING", EnableSharedSnat: false},
	}

	filtered := filterVpcs(input, vpcFilters{
		Name:   "match",
		CIDR:   "192.168.0.0/16",
		Status: "OK",
		Shared: &shared,
	})
	if len(filtered) != 1 || filtered[0].ID != "first" {
		t.Fatalf("unexpected filtered VPCs: %#v", filtered)
	}
}
