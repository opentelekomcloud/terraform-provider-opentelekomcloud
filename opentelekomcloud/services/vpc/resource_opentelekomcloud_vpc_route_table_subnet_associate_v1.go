package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/routetables"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

// ResourceVPCRouteTableSubnetAssociateV1 manages the association between a subnet and a route table.
func ResourceVPCRouteTableSubnetAssociateV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcRouteTableSubnetAssociateCreate,
		ReadContext:   resourceVpcRouteTableSubnetAssociateRead,
		DeleteContext: resourceVpcRouteTableSubnetAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVpcRouteTableSubnetAssociateImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcRouteTableSubnetAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	rtID := d.Get("route_table_id").(string)
	subnetID := d.Get("subnet_id").(string)

	actionOpts := routetables.ActionOpts{
		Subnets: routetables.ActionSubnetsOpts{
			Associate: []string{subnetID},
		},
	}

	log.Printf("[DEBUG] OpenTelekomCloud VPC route table subnet associate: rtb=%s subnet=%s", rtID, subnetID)
	_, err = routetables.Action(client, rtID, actionOpts)
	if err != nil {
		return diag.Errorf("error associating subnet %s with route table %s: %s", subnetID, rtID, err)
	}

	d.SetId(rtID + "/" + subnetID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcRouteTableSubnetAssociateRead(clientCtx, d, meta)
}

func resourceVpcRouteTableSubnetAssociateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	rtID, subnetID, err := parseRouteTableSubnetAssociateID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	routeTable, err := routetables.Get(client, rtID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud VPC route table subnet association")
	}

	if !subnetInRouteTable(routeTable.Subnets, subnetID) {
		log.Printf("[WARN] Subnet %s not found in route table %s, removing association from state", subnetID, rtID)
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("route_table_id", rtID),
		d.Set("subnet_id", subnetID),
		d.Set("vpc_id", routeTable.VpcID),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud VPC route table subnet association: %s", err)
	}

	return nil
}

func resourceVpcRouteTableSubnetAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	rtID, subnetID, err := parseRouteTableSubnetAssociateID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	actionOpts := routetables.ActionOpts{
		Subnets: routetables.ActionSubnetsOpts{
			Disassociate: []string{subnetID},
		},
	}

	log.Printf("[DEBUG] OpenTelekomCloud VPC route table subnet disassociate: rtb=%s subnet=%s", rtID, subnetID)
	_, err = routetables.Action(client, rtID, actionOpts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Route table %s already deleted, removing association from state", rtID)
			return nil
		}
		return diag.Errorf("error disassociating subnet %s from route table %s: %s", subnetID, rtID, err)
	}

	return nil
}

func resourceVpcRouteTableSubnetAssociateImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want <route_table_id>/<subnet_id>, got: %s", d.Id())
	}

	d.SetId(parts[0] + "/" + parts[1])

	mErr := multierror.Append(nil,
		d.Set("route_table_id", parts[0]),
		d.Set("subnet_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func parseRouteTableSubnetAssociateID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid route table subnet association ID format: %s", id)
	}
	return parts[0], parts[1], nil
}

func subnetInRouteTable(subnets []routetables.Subnet, subnetID string) bool {
	for _, s := range subnets {
		if s.ID == subnetID {
			return true
		}
	}
	return false
}
