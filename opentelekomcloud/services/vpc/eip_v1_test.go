package vpc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	legacyEips "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/publicips"
)

func TestVpcEIPV1Schemas(t *testing.T) {
	resourceSchema := ResourceVpcEIPV1().Schema
	for _, field := range []string{"enterprise_project_id", "public_border_group", "allow_share_bandwidth_types"} {
		if resourceSchema[field] == nil {
			t.Fatalf("resource schema is missing %s", field)
		}
	}
	publicIPSchema := resourceSchema["publicip"].Elem.(*schema.Resource).Schema
	if publicIPSchema["ip_version"] == nil {
		t.Fatal("resource publicip schema is missing ip_version")
	}

	dataSourceSchema := DataSourceVPCEipV1().Schema
	for _, field := range []string{"enterprise_project_id", "public_border_group", "allow_share_bandwidth_types"} {
		if dataSourceSchema[field] == nil || !dataSourceSchema[field].Computed {
			t.Fatalf("data source %s must be computed", field)
		}
	}
}

func TestFilterPublicIPs(t *testing.T) {
	input := []publicips.PublicIP{
		{
			ID: "first", Status: "DOWN", PrivateIpAddress: "192.168.0.10", PortId: "port-id",
			BandwidthId: "bandwidth-id", PublicIpAddress: "192.0.2.10", Alias: "matching-eip",
		},
		{ID: "second", Status: "ACTIVE", PublicIpAddress: "192.0.2.11", Alias: "other"},
	}

	filtered := filterPublicIPs(input, publicIPFilters{
		ID:               "first",
		Status:           "DOWN",
		PrivateAddress:   "192.168.0.10",
		PortID:           "port-id",
		BandwidthID:      "bandwidth-id",
		PublicAddress:    "192.0.2.10",
		NameRegexPattern: "^matching-",
	})
	if len(filtered) != 1 || filtered[0].ID != "first" {
		t.Fatalf("unexpected filtered EIPs: %#v", filtered)
	}
}

func TestExpandVpcEIPPublicIP(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceVpcEIPV1().Schema, map[string]interface{}{
		"publicip": []interface{}{map[string]interface{}{
			"type":       "5_bgp",
			"name":       "test-eip",
			"ip_address": "192.0.2.10",
			"ip_version": 4,
		}},
	})

	opts := resourcePublicIP(data)
	if opts.Type != "5_bgp" || opts.Alias != "test-eip" ||
		opts.IpAddress != "192.0.2.10" || opts.IPVersion != 4 {
		t.Fatalf("unexpected public IP options: %#v", opts)
	}
}

func TestEIPCreateOptsValueSpecs(t *testing.T) {
	opts := EIPCreateOpts{
		ApplyOpts: legacyEips.ApplyOpts{
			IP:        legacyEips.PublicIpOpts{Type: "5_bgp"},
			Bandwidth: legacyEips.BandwidthOpts{ShareType: "PER"},
		},
		ValueSpecs: map[string]string{"enterprise_project_id": "project-id"},
	}

	body, err := opts.ToPublicIpApplyMap()
	if err != nil {
		t.Fatalf("failed to build legacy create request: %v", err)
	}
	if body["enterprise_project_id"] != "project-id" {
		t.Fatalf("value_specs were not expanded: %#v", body)
	}
	if _, ok := body["value_specs"]; ok {
		t.Fatalf("value_specs wrapper was not removed: %#v", body)
	}
}

func TestResourcePublicIPOmittedVersion(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceVpcEIPV1().Schema, map[string]interface{}{
		"publicip": []interface{}{map[string]interface{}{"type": "5_bgp"}},
	})

	if got := resourcePublicIP(data).IPVersion; got != 0 {
		t.Fatalf("expected omitted IP version to remain zero, got %d", got)
	}
}
