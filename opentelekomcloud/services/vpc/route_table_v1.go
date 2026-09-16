package vpc

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/routetables"
)

func selectVpcRouteTableV1(routeTables []routetables.RouteTable, id, name string) (*routetables.RouteTable, error) {
	var selected *routetables.RouteTable
	for i := range routeTables {
		routeTable := &routeTables[i]
		if id != "" && routeTable.ID != id {
			continue
		}
		if name != "" && routeTable.Name != name {
			continue
		}
		if id == "" && name == "" && !routeTable.Default {
			continue
		}
		if selected != nil {
			return nil, fmt.Errorf("query returned more than one VPC route table")
		}
		selected = routeTable
	}
	return selected, nil
}

func setVpcRouteTableV1Fields(d *schema.ResourceData, routeTable *routetables.RouteTable, region string) error {
	return multierror.Append(
		d.Set("region", region),
		d.Set("vpc_id", routeTable.VpcID),
		d.Set("name", routeTable.Name),
		d.Set("description", routeTable.Description),
		d.Set("default", routeTable.Default),
		d.Set("tenant_id", routeTable.TenantID),
		d.Set("route", expandRouteTableRoutes(routeTable.Routes)),
		d.Set("subnets", expandRouteTableSubnets(routeTable.Subnets)),
		d.Set("created_at", routeTable.CreatedAt),
		d.Set("updated_at", routeTable.UpdatedAt),
	).ErrorOrNil()
}

func flattenVpcRouteTablesV1(routeTables []routetables.RouteTable) []map[string]interface{} {
	result := make([]map[string]interface{}, len(routeTables))
	for i, routeTable := range routeTables {
		result[i] = map[string]interface{}{
			"id":          routeTable.ID,
			"name":        routeTable.Name,
			"default":     routeTable.Default,
			"tenant_id":   routeTable.TenantID,
			"vpc_id":      routeTable.VpcID,
			"description": routeTable.Description,
			"subnets":     expandRouteTableSubnets(routeTable.Subnets),
			"routes":      expandRouteTableRoutes(routeTable.Routes),
			"created_at":  routeTable.CreatedAt,
			"updated_at":  routeTable.UpdatedAt,
		}
	}
	return result
}
