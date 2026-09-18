package vpcep

import "testing"

func TestEnterpriseProjectSchemas(t *testing.T) {
	for name, field := range map[string]struct {
		optional bool
		computed bool
		forceNew bool
	}{
		"endpoint resource": {
			ResourceVPCEPEndpointV1().Schema["enterprise_project_id"].Optional,
			ResourceVPCEPEndpointV1().Schema["enterprise_project_id"].Computed,
			ResourceVPCEPEndpointV1().Schema["enterprise_project_id"].ForceNew,
		},
		"service resource": {
			ResourceVPCEPServiceV1().Schema["enterprise_project_id"].Optional,
			ResourceVPCEPServiceV1().Schema["enterprise_project_id"].Computed,
			ResourceVPCEPServiceV1().Schema["enterprise_project_id"].ForceNew,
		},
	} {
		if !field.optional || !field.computed || !field.forceNew {
			t.Fatalf("%s enterprise_project_id must be optional, computed, and force replacement", name)
		}
	}

	if !DataSourceVPCEPServiceV1().Schema["enterprise_project_id"].Computed {
		t.Fatal("VPCEP service data source must export enterprise_project_id")
	}
}
