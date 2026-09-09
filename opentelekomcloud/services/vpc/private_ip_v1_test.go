package vpc

import (
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/privateips"
)

func TestVpcSubnetPrivateIPV1Schemas(t *testing.T) {
	resourceSchema := ResourceVpcSubnetPrivateIPV1().Schema
	if !resourceSchema["subnet_id"].Required || !resourceSchema["subnet_id"].ForceNew {
		t.Fatal("resource subnet_id must be required and ForceNew")
	}
	if !resourceSchema["ip_address"].Optional || !resourceSchema["ip_address"].Computed ||
		!resourceSchema["ip_address"].ForceNew {
		t.Fatal("resource ip_address must be optional, computed, and ForceNew")
	}
	for _, field := range []string{"status", "tenant_id", "device_owner"} {
		if !resourceSchema[field].Computed {
			t.Fatalf("resource %s must be computed", field)
		}
	}

	dataSourceSchema := DataSourceVpcSubnetPrivateIPV1().Schema
	if !dataSourceSchema["id"].Required {
		t.Fatal("data source id must be required")
	}
	for _, field := range []string{"subnet_id", "ip_address", "status", "tenant_id", "device_owner"} {
		if !dataSourceSchema[field].Computed {
			t.Fatalf("data source %s must be computed", field)
		}
	}
}

func TestValidateVpcSubnetPrivateIPCreateResponse(t *testing.T) {
	if err := validateVpcSubnetPrivateIPCreateResponse([]privateips.PrivateIP{{ID: "private-ip-id"}}); err != nil {
		t.Fatalf("unexpected valid response error: %s", err)
	}
	for _, response := range [][]privateips.PrivateIP{
		nil,
		{{ID: "first"}, {ID: "second"}},
	} {
		if err := validateVpcSubnetPrivateIPCreateResponse(response); err == nil {
			t.Fatalf("expected response length %d to be rejected", len(response))
		}
	}
}
