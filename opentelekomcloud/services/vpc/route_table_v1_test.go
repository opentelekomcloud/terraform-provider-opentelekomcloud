package vpc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/routetables"
)

func TestVpcRouteTableV1SchemasPreserveCompatibility(t *testing.T) {
	resourceSchema := ResourceVPCRouteTableV1().Schema
	for _, field := range []string{"vpc_id", "name", "description", "subnets", "route", "created_at", "updated_at"} {
		if _, ok := resourceSchema[field]; !ok {
			t.Fatalf("resource compatibility field %s is missing", field)
		}
	}
	for _, field := range []string{"default", "tenant_id"} {
		if !resourceSchema[field].Computed {
			t.Fatalf("resource field %s must be computed", field)
		}
	}

	singularSchema := DataSourceVPCRouteTableV1().Schema
	for _, field := range []string{"tenant_id", "created_at", "updated_at"} {
		if !singularSchema[field].Computed {
			t.Fatalf("singular data source field %s must be computed", field)
		}
	}

	pluralItemSchema := DataSourceVpcRouteTablesV1().Schema["routetables"].Elem.(*schema.Resource).Schema
	for _, field := range []string{"created_at", "updated_at"} {
		if !pluralItemSchema[field].Computed {
			t.Fatalf("plural data source item field %s must be computed", field)
		}
	}

	routeSchema := ResourceVPCRouteTableRouteV1().Schema
	for _, field := range []string{
		"vpc_id", "destination", "type", "nexthop", "description", "route_table_id", "route_table_name",
	} {
		if _, ok := routeSchema[field]; !ok {
			t.Fatalf("route resource compatibility field %s is missing", field)
		}
	}

	associationSchema := ResourceVPCRouteTableSubnetAssociateV1().Schema
	for _, field := range []string{"route_table_id", "subnet_id", "vpc_id"} {
		if _, ok := associationSchema[field]; !ok {
			t.Fatalf("subnet association compatibility field %s is missing", field)
		}
	}
}

func TestSelectVpcRouteTableV1(t *testing.T) {
	routeTables := []routetables.RouteTable{
		{ID: "default-id", Name: "default", Default: true},
		{ID: "custom-id", Name: "custom"},
	}

	selected, err := selectVpcRouteTableV1(routeTables, "custom-id", "")
	if err != nil {
		t.Fatalf("unexpected ID selection error: %s", err)
	}
	if selected == nil || selected.ID != "custom-id" {
		t.Fatalf("expected custom route table by ID, got %#v", selected)
	}

	selected, err = selectVpcRouteTableV1(routeTables, "", "")
	if err != nil {
		t.Fatalf("unexpected default selection error: %s", err)
	}
	if selected == nil || selected.ID != "default-id" {
		t.Fatalf("expected default route table, got %#v", selected)
	}

	_, err = selectVpcRouteTableV1([]routetables.RouteTable{
		{ID: "first", Name: "duplicate"},
		{ID: "second", Name: "duplicate"},
	}, "", "duplicate")
	if err == nil {
		t.Fatal("expected ambiguous name selection to fail")
	}
}

func TestVpcRouteTableV1MappingsIncludeAllSdkFields(t *testing.T) {
	routeTable := routetables.RouteTable{
		ID:          "route-table-id",
		Name:        "route-table",
		Default:     false,
		TenantID:    "tenant-id",
		VpcID:       "vpc-id",
		Description: "description",
		Subnets:     []routetables.Subnet{{ID: "subnet-id"}},
		Routes: []routetables.Route{
			{Type: "local", DestinationCIDR: "192.168.0.0/16"},
			{Type: "ecs", DestinationCIDR: "10.0.0.0/24", NextHop: "server-id", Description: "route"},
		},
		CreatedAt: "2026-09-15T08:00:00",
		UpdatedAt: "2026-09-15T08:01:00",
	}

	data := schema.TestResourceDataRaw(t, ResourceVPCRouteTableV1().Schema, map[string]interface{}{
		"vpc_id": "vpc-id",
		"name":   "route-table",
	})
	if err := setVpcRouteTableV1Fields(data, &routeTable, "eu-de"); err != nil {
		t.Fatalf("unexpected mapping error: %s", err)
	}
	for field, expected := range map[string]interface{}{
		"region":     "eu-de",
		"tenant_id":  "tenant-id",
		"created_at": "2026-09-15T08:00:00",
		"updated_at": "2026-09-15T08:01:00",
	} {
		if actual := data.Get(field); actual != expected {
			t.Fatalf("expected %s to be %#v, got %#v", field, expected, actual)
		}
	}
	if routes := data.Get("route").(*schema.Set); routes.Len() != 1 {
		t.Fatalf("expected local route to be excluded, got %d routes", routes.Len())
	}

	flattened := flattenVpcRouteTablesV1([]routetables.RouteTable{routeTable})
	if len(flattened) != 1 {
		t.Fatalf("expected one flattened route table, got %d", len(flattened))
	}
	if flattened[0]["tenant_id"] != "tenant-id" || flattened[0]["created_at"] != "2026-09-15T08:00:00" {
		t.Fatalf("missing SDK response fields: %#v", flattened[0])
	}
	if routes := flattened[0]["routes"].([]map[string]interface{}); len(routes) != 1 {
		t.Fatalf("expected local route to be excluded, got %#v", routes)
	}
}
