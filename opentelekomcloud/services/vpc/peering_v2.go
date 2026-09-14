package vpc

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/peerings"
)

func setVpcPeeringV2Fields(d *schema.ResourceData, peering *peerings.Peering, region string) error {
	return multierror.Append(
		d.Set("region", region),
		d.Set("name", peering.Name),
		d.Set("description", peering.Description),
		d.Set("status", peering.Status),
		d.Set("vpc_id", peering.RequestVpcInfo.VpcID),
		d.Set("vpc_tenant_id", peering.RequestVpcInfo.TenantID),
		d.Set("peer_vpc_id", peering.AcceptVpcInfo.VpcID),
		d.Set("peer_tenant_id", peering.AcceptVpcInfo.TenantID),
		d.Set("created_at", peering.CreatedAt),
		d.Set("updated_at", peering.UpdatedAt),
	).ErrorOrNil()
}

func filterVpcPeerings(
	items []peerings.Peering,
	vpcID, vpcTenantID, peerVpcID, peerTenantID string,
) []peerings.Peering {
	filtered := make([]peerings.Peering, 0, len(items))
	for _, item := range items {
		if vpcID != "" && item.RequestVpcInfo.VpcID != vpcID {
			continue
		}
		if vpcTenantID != "" && item.RequestVpcInfo.TenantID != vpcTenantID {
			continue
		}
		if peerVpcID != "" && item.AcceptVpcInfo.VpcID != peerVpcID {
			continue
		}
		if peerTenantID != "" && item.AcceptVpcInfo.TenantID != peerTenantID {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func flattenPeeringConnections(items []peerings.Peering) []interface{} {
	if items == nil {
		return nil
	}

	result := make([]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":             item.ID,
			"name":           item.Name,
			"description":    item.Description,
			"status":         item.Status,
			"vpc_id":         item.RequestVpcInfo.VpcID,
			"vpc_tenant_id":  item.RequestVpcInfo.TenantID,
			"peer_vpc_id":    item.AcceptVpcInfo.VpcID,
			"peer_tenant_id": item.AcceptVpcInfo.TenantID,
			"created_at":     item.CreatedAt,
			"updated_at":     item.UpdatedAt,
		}
	}
	return result
}
