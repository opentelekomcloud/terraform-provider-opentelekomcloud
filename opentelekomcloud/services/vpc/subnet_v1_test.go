package vpc

import (
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/subnets"
)

func TestSubnetV1Schemas(t *testing.T) {
	resourceSchema := ResourceVpcSubnetV1().Schema
	for _, field := range []string{"scope", "tenant_id", "created_at", "updated_at"} {
		if !resourceSchema[field].Computed || resourceSchema[field].Optional {
			t.Fatalf("resource %s must be computed only", field)
		}
	}

	dataSourceSchema := DataSourceVpcSubnetV1().Schema
	for _, field := range []string{"tenant_id", "created_at", "updated_at"} {
		if !dataSourceSchema[field].Computed || dataSourceSchema[field].Optional {
			t.Fatalf("data source %s must be computed only", field)
		}
	}
	for _, field := range []string{"description", "ntp_addresses", "scope"} {
		if !dataSourceSchema[field].Computed || !dataSourceSchema[field].Optional {
			t.Fatalf("data source %s must be optional and computed", field)
		}
	}
}

func TestFilterSubnets(t *testing.T) {
	input := []subnets.Subnet{
		{
			ID: "first", Name: "match", Description: "description", CIDR: "192.168.0.0/24",
			Status: "ACTIVE", GatewayIP: "192.168.0.1", PrimaryDNS: "100.125.4.25",
			SecondaryDNS: "100.125.129.199", AvailabilityZone: "eu-de-01", Scope: "center",
			ExtraDHCPOpts: []subnets.ExtraDHCPOpt{{OptName: "ntp", OptValue: "10.100.0.33"}},
		},
		{ID: "second", Name: "other", CIDR: "10.0.0.0/24", Status: "ACTIVE", VpcID: "other-vpc"},
	}

	filtered := filterSubnets(input, subnetFilters{
		ID:               "first",
		Name:             "match",
		Description:      "description",
		CIDR:             "192.168.0.0/24",
		Status:           "ACTIVE",
		GatewayIP:        "192.168.0.1",
		PrimaryDNS:       "100.125.4.25",
		SecondaryDNS:     "100.125.129.199",
		AvailabilityZone: "eu-de-01",
		NtpAddresses:     "10.100.0.33",
		Scope:            "center",
	})
	if len(filtered) != 1 || filtered[0].ID != "first" {
		t.Fatalf("unexpected filtered subnets: %#v", filtered)
	}
}

func TestSubnetNtpAddresses(t *testing.T) {
	opts := []subnets.ExtraDHCPOpt{
		{OptName: "other", OptValue: "ignored"},
		{OptName: "ntp", OptValue: "10.100.0.33,10.100.0.34"},
	}
	if actual := subnetNtpAddresses(opts); actual != "10.100.0.33,10.100.0.34" {
		t.Fatalf("unexpected NTP addresses: %q", actual)
	}
	if actual := subnetNtpAddresses(nil); actual != "" {
		t.Fatalf("expected empty NTP addresses, got %q", actual)
	}
}
