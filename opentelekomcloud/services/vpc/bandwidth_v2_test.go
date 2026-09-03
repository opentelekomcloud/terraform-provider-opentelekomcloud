package vpc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/bandwidths"
)

func TestBandwidthV2Schemas(t *testing.T) {
	resourceSchema := ResourceBandwidthV2().Schema
	for _, field := range []string{
		"share_type", "bandwidth_type", "charge_mode", "billing_info", "tenant_id",
		"enterprise_project_id", "public_border_group", "created_at", "updated_at", "publicip_info",
	} {
		if resourceSchema[field] == nil {
			t.Fatalf("resource schema is missing %s", field)
		}
	}

	dataSourceSchema := DataSourceBandWidthV2().Schema
	for _, field := range []string{
		"id", "billing_info", "tenant_id", "enterprise_project_id",
		"public_border_group", "created_at", "updated_at", "publicip_info",
	} {
		if dataSourceSchema[field] == nil {
			t.Fatalf("data source schema is missing %s", field)
		}
	}

	publicIPSchema := resourceSchema["publicip_info"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"id", "address", "ipv6_address", "ip_version", "type"} {
		if publicIPSchema[field] == nil {
			t.Fatalf("publicip_info schema is missing %s", field)
		}
	}
}

func TestFilterBandwidthsV1(t *testing.T) {
	input := []bandwidths.BandWidth{
		{ID: "first", Name: "matching", Size: 10, ShareType: "WHOLE", PublicBorderGroup: "center"},
		{ID: "second", Name: "other", Size: 20, ShareType: "PER", PublicBorderGroup: "edge"},
	}

	filtered := filterBandwidthsV1(input, bandwidthFilters{
		ID:                "first",
		Name:              "matching",
		Size:              10,
		ShareType:         "WHOLE",
		PublicBorderGroup: "center",
	})
	if len(filtered) != 1 || filtered[0].ID != "first" {
		t.Fatalf("unexpected filtered bandwidths: %#v", filtered)
	}
}

func TestFlattenBandwidthPublicIPs(t *testing.T) {
	flattened := flattenBandwidthPublicIPs([]bandwidths.PublicIpinfo{{
		PublicipId:        "public-ip-id",
		PublicipAddress:   "192.0.2.10",
		Publicipv6Address: "2001:db8::10",
		IPVersion:         6,
		PublicipType:      "5_dualStack",
	}})
	if len(flattened) != 1 ||
		flattened[0]["id"] != "public-ip-id" ||
		flattened[0]["ipv6_address"] != "2001:db8::10" ||
		flattened[0]["ip_version"] != 6 {
		t.Fatalf("unexpected flattened public IPs: %#v", flattened)
	}
}

func TestRemovablePublicIPsSkipsDetachedIDs(t *testing.T) {
	attached := []bandwidths.PublicIpinfo{
		{PublicipId: "still-attached"},
		{PublicipId: "unmanaged"},
		{PublicipId: "ipv6-port", PublicipType: dualStackPublicIPType},
	}
	requested := schema.NewSet(schema.HashString, []interface{}{"already-detached", "still-attached", "ipv6-port"})

	publicIPs, ports := removablePublicIPs(attached, requested)
	if len(publicIPs) != 1 || publicIPs[0].PublicipId != "still-attached" {
		t.Fatalf("unexpected removable public IPs: %#v", publicIPs)
	}
	if len(ports) != 1 ||
		ports[0].PublicIpID != "ipv6-port" ||
		ports[0].PublicIpType != dualStackPublicIPType {
		t.Fatalf("unexpected removable IPv6 ports: %#v", ports)
	}
}
