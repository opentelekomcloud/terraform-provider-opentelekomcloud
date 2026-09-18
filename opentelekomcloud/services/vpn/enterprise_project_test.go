package vpn

import "testing"

func TestEnterpriseProjectSchemas(t *testing.T) {
	for name, field := range map[string]struct {
		optional bool
		computed bool
		forceNew bool
	}{
		"connection resource": {
			ResourceEnterpriseConnection().Schema["enterprise_project_id"].Optional,
			ResourceEnterpriseConnection().Schema["enterprise_project_id"].Computed,
			ResourceEnterpriseConnection().Schema["enterprise_project_id"].ForceNew,
		},
		"gateway resource": {
			ResourceEnterpriseVpnGateway().Schema["enterprise_project_id"].Optional,
			ResourceEnterpriseVpnGateway().Schema["enterprise_project_id"].Computed,
			ResourceEnterpriseVpnGateway().Schema["enterprise_project_id"].ForceNew,
		},
	} {
		if !field.optional || !field.computed || !field.forceNew {
			t.Fatalf("%s enterprise_project_id must be optional, computed, and force replacement", name)
		}
	}

	if !DataSourceEnterpriseConnection().Schema["enterprise_project_id"].Computed {
		t.Fatal("VPN connection data source must export enterprise_project_id")
	}
	if !DataSourceEnterpriseVpnGateway().Schema["enterprise_project_id"].Computed {
		t.Fatal("VPN gateway data source must export enterprise_project_id")
	}
}
