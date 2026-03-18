package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/routetables"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

// ResourceVPCRouteTableRouteV1 manages individual routes within a VPC route table.
func ResourceVPCRouteTableRouteV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcRouteTableRouteCreate,
		ReadContext:   resourceVpcRouteTableRouteRead,
		UpdateContext: resourceVpcRouteTableRouteUpdate,
		DeleteContext: resourceVpcRouteTableRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVpcRouteTableRouteImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateCIDR,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ecs", "eni", "vip", "nat", "peering", "vpn",
					"dc", "egw", "er", "subeni", "local",
				}, false),
			},
			"nexthop": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"route_table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcRouteTableRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	routeTableID := d.Get("route_table_id").(string)
	if routeTableID == "" {
		rtb, err := getDefaultRouteTable(client, d.Get("vpc_id").(string))
		if err != nil {
			return diag.Errorf("error retrieving default route table for VPC %s: %s", d.Get("vpc_id").(string), err)
		}
		routeTableID = rtb.ID
	}

	destination := d.Get("destination").(string)
	desc := d.Get("description").(string)
	routeOpts := routetables.RouteOpts{
		Destination: destination,
		Type:        d.Get("type").(string),
		NextHop:     d.Get("nexthop").(string),
		Description: &desc,
	}

	updateOpts := routetables.UpdateOpts{
		Routes: map[string][]routetables.RouteOpts{
			"add": {routeOpts},
		},
	}

	log.Printf("[DEBUG] OpenTelekomCloud VPC route table route create: rtb=%s, %#v", routeTableID, updateOpts)
	if err := routetables.Update(client, routeTableID, updateOpts); err != nil {
		return diag.Errorf("error creating OpenTelekomCloud VPC route table route: %s", err)
	}

	d.SetId(routeTableID + "/" + destination)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcRouteTableRouteRead(clientCtx, d, meta)
}

func resourceVpcRouteTableRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	routeTableID, destination, err := parseRouteTableRouteID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	routeTable, err := routetables.Get(client, routeTableID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud VPC route table route")
	}

	route := findRouteByDestination(routeTable.Routes, destination)
	if route == nil {
		log.Printf("[WARN] OpenTelekomCloud VPC route %s not found in route table %s, removing from state", destination, routeTableID)
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("vpc_id", routeTable.VpcID),
		d.Set("route_table_id", routeTableID),
		d.Set("route_table_name", routeTable.Name),
		d.Set("destination", route.DestinationCIDR),
		d.Set("type", route.Type),
		d.Set("nexthop", route.NextHop),
		d.Set("description", route.Description),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud VPC route table route: %s", err)
	}

	return nil
}

func resourceVpcRouteTableRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	routeTableID, destination, err := parseRouteTableRouteID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	desc := d.Get("description").(string)
	routeOpts := routetables.RouteOpts{
		Destination: destination,
		Type:        d.Get("type").(string),
		NextHop:     d.Get("nexthop").(string),
		Description: &desc,
	}

	updateOpts := routetables.UpdateOpts{
		Routes: map[string][]routetables.RouteOpts{
			"mod": {routeOpts},
		},
	}

	log.Printf("[DEBUG] OpenTelekomCloud VPC route table route update: rtb=%s, %#v", routeTableID, updateOpts)
	if err := routetables.Update(client, routeTableID, updateOpts); err != nil {
		return diag.Errorf("error updating OpenTelekomCloud VPC route table route: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcRouteTableRouteRead(clientCtx, d, meta)
}

func resourceVpcRouteTableRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	routeTableID, destination, err := parseRouteTableRouteID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	routeOpts := routetables.RouteOpts{
		Destination: destination,
		Type:        d.Get("type").(string),
		NextHop:     d.Get("nexthop").(string),
	}

	updateOpts := routetables.UpdateOpts{
		Routes: map[string][]routetables.RouteOpts{
			"del": {routeOpts},
		},
	}

	log.Printf("[DEBUG] OpenTelekomCloud VPC route table route delete: rtb=%s, %#v", routeTableID, updateOpts)
	if err := routetables.Update(client, routeTableID, updateOpts); err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Route table %s already deleted, removing route from state", routeTableID)
			return nil
		}
		return diag.Errorf("error deleting OpenTelekomCloud VPC route table route: %s", err)
	}

	return nil
}

func resourceVpcRouteTableRouteImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want <route_table_id>/<destination>, got: %s", d.Id())
	}

	d.SetId(parts[0] + "/" + parts[1])
	if err := d.Set("route_table_id", parts[0]); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func getDefaultRouteTable(client *golangsdk.ServiceClient, vpcID string) (*routetables.RouteTable, error) {
	rts, err := routetables.List(client, routetables.ListOpts{VpcID: vpcID})
	if err != nil {
		return nil, fmt.Errorf("error listing route tables for VPC %s: %w", vpcID, err)
	}

	for _, rt := range rts {
		if rt.Default {
			return &rt, nil
		}
	}

	return nil, fmt.Errorf("no default route table found for VPC %s", vpcID)
}

func parseRouteTableRouteID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid route table route ID format: %s", id)
	}
	return parts[0], parts[1], nil
}

func findRouteByDestination(routes []routetables.Route, destination string) *routetables.Route {
	for _, r := range routes {
		if r.DestinationCIDR == destination {
			return &r
		}
	}
	return nil
}
