package vpc

import (
	"reflect"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/quotas"
)

func TestDataSourceVpcQuotasV1Schema(t *testing.T) {
	dataSource := DataSourceVpcQuotasV1()

	if !dataSource.Schema["type"].Optional {
		t.Fatal("type must be an optional filter")
	}
	if !dataSource.Schema["region"].Computed {
		t.Fatal("region must be computed")
	}
	if !dataSource.Schema["quotas"].Computed {
		t.Fatal("quotas must be computed")
	}
}

func TestFlattenVpcQuotas(t *testing.T) {
	input := []quotas.Quota{
		{Type: "publicIp", Used: 2, Quota: 10, Min: 0},
		{Type: "vpc", Used: 4, Quota: -1, Min: 0},
	}
	expected := []map[string]interface{}{
		{"type": "publicIp", "used": 2, "quota": 10, "min": 0},
		{"type": "vpc", "used": 4, "quota": -1, "min": 0},
	}

	if actual := flattenVpcQuotas(input); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("unexpected flattened quotas: %#v", actual)
	}
	if actual := flattenVpcQuotas(nil); actual != nil {
		t.Fatalf("expected nil for empty quotas, got %#v", actual)
	}
}
