package vpc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/peerings"
)

func TestVpcPeeringV2SchemasPreserveAndExpandCompatibility(t *testing.T) {
	resourceSchema := ResourceVpcPeeringConnectionV2().Schema
	for _, field := range []string{"name", "description", "vpc_id", "peer_vpc_id", "peer_tenant_id"} {
		if _, ok := resourceSchema[field]; !ok {
			t.Fatalf("resource compatibility field %s is missing", field)
		}
	}
	for _, field := range []string{"vpc_tenant_id", "created_at", "updated_at"} {
		if !resourceSchema[field].Computed {
			t.Fatalf("resource field %s must be computed", field)
		}
	}

	accepterSchema := ResourceVpcPeeringConnectionAccepterV2().Schema
	for _, field := range []string{
		"name", "description", "status", "vpc_id", "vpc_tenant_id", "peer_vpc_id",
		"peer_tenant_id", "created_at", "updated_at",
	} {
		if !accepterSchema[field].Computed {
			t.Fatalf("accepter field %s must be computed", field)
		}
	}

	singularSchema := DataSourceVpcPeeringConnectionV2().Schema
	for _, field := range []string{
		"id", "name", "status", "vpc_id", "vpc_tenant_id", "peer_vpc_id", "peer_tenant_id",
	} {
		if !singularSchema[field].Optional {
			t.Fatalf("singular data source filter %s must remain optional", field)
		}
	}
	for _, field := range []string{"created_at", "updated_at"} {
		if !singularSchema[field].Computed {
			t.Fatalf("singular data source field %s must be computed", field)
		}
	}

	pluralSchema := DataSourceVpcPeeringConnectionsV2().Schema
	if !pluralSchema["vpc_tenant_id"].Optional {
		t.Fatal("plural data source vpc_tenant_id must be optional")
	}
	itemSchema := pluralSchema["peering_connections"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"created_at", "updated_at"} {
		if !itemSchema[field].Computed {
			t.Fatalf("plural data source item field %s must be computed", field)
		}
	}
}

func TestFilterVpcPeerings(t *testing.T) {
	items := []peerings.Peering{
		{
			ID:             "first",
			RequestVpcInfo: peerings.VpcInfo{VpcID: "request-vpc", TenantID: "request-tenant"},
			AcceptVpcInfo:  peerings.VpcInfo{VpcID: "peer-vpc", TenantID: "peer-tenant"},
		},
		{
			ID:             "second",
			RequestVpcInfo: peerings.VpcInfo{VpcID: "other-request-vpc", TenantID: "other-request-tenant"},
			AcceptVpcInfo:  peerings.VpcInfo{VpcID: "request-vpc", TenantID: "other-peer-tenant"},
		},
	}

	filtered := filterVpcPeerings(items, "request-vpc", "request-tenant", "peer-vpc", "peer-tenant")
	if len(filtered) != 1 || filtered[0].ID != "first" {
		t.Fatalf("unexpected filtered peerings: %#v", filtered)
	}
	if filtered := filterVpcPeerings(items, "request-vpc", "", "", ""); len(filtered) != 1 || filtered[0].ID != "first" {
		t.Fatalf("expected requester-side VPC match only, got %#v", filtered)
	}
	if filtered := filterVpcPeerings(items, "", "", "", "missing"); len(filtered) != 0 {
		t.Fatalf("expected no matches, got %#v", filtered)
	}
}

func TestFlattenPeeringConnectionsIncludesAllSdkFields(t *testing.T) {
	items := flattenPeeringConnections([]peerings.Peering{{
		ID:             "peering-id",
		Name:           "peering",
		Description:    "description",
		Status:         "ACTIVE",
		RequestVpcInfo: peerings.VpcInfo{VpcID: "request-vpc", TenantID: "request-tenant"},
		AcceptVpcInfo:  peerings.VpcInfo{VpcID: "peer-vpc", TenantID: "peer-tenant"},
		CreatedAt:      "2026-09-11T10:00:00",
		UpdatedAt:      "2026-09-11T10:01:00",
	}})
	if len(items) != 1 {
		t.Fatalf("expected one flattened peering, got %d", len(items))
	}
	item := items[0].(map[string]interface{})
	for field, expected := range map[string]interface{}{
		"vpc_tenant_id":  "request-tenant",
		"peer_vpc_id":    "peer-vpc",
		"peer_tenant_id": "peer-tenant",
		"created_at":     "2026-09-11T10:00:00",
		"updated_at":     "2026-09-11T10:01:00",
	} {
		if item[field] != expected {
			t.Fatalf("expected %s to be %#v, got %#v", field, expected, item[field])
		}
	}
}

func TestSetVpcPeeringV2FieldsIncludesAllSdkFields(t *testing.T) {
	data := schema.TestResourceDataRaw(t, ResourceVpcPeeringConnectionV2().Schema, map[string]interface{}{
		"name":        "peering",
		"vpc_id":      "9daeac7c-a98f-430f-8e38-67f9c044e299",
		"peer_vpc_id": "f583c072-0bb8-4e19-afb2-afb7c1693be5",
	})
	peering := &peerings.Peering{
		Name:           "peering",
		Description:    "description",
		Status:         "ACTIVE",
		RequestVpcInfo: peerings.VpcInfo{VpcID: "request-vpc", TenantID: "request-tenant"},
		AcceptVpcInfo:  peerings.VpcInfo{VpcID: "peer-vpc", TenantID: "peer-tenant"},
		CreatedAt:      "2026-09-11T10:00:00",
		UpdatedAt:      "2026-09-11T10:01:00",
	}
	if err := setVpcPeeringV2Fields(data, peering, "eu-de"); err != nil {
		t.Fatalf("unexpected set error: %s", err)
	}
	for field, expected := range map[string]interface{}{
		"region":         "eu-de",
		"description":    "description",
		"status":         "ACTIVE",
		"vpc_id":         "request-vpc",
		"vpc_tenant_id":  "request-tenant",
		"peer_vpc_id":    "peer-vpc",
		"peer_tenant_id": "peer-tenant",
		"created_at":     "2026-09-11T10:00:00",
		"updated_at":     "2026-09-11T10:01:00",
	} {
		if actual := data.Get(field); actual != expected {
			t.Fatalf("expected %s to be %#v, got %#v", field, expected, actual)
		}
	}
}
