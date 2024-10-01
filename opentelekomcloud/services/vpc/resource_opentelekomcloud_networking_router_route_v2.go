package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNetworkingRouterRouteV2() *schema.Resource {
	return &schema.Resource{
		CreateContext:      resourceNetworkingRouterRouteV2Create,
		ReadContext:        resourceNetworkingRouterRouteV2Read,
		DeleteContext:      resourceNetworkingRouterRouteV2Delete,
		DeprecationMessage: "use opentelekomcloud_vpc_route_v2 resource instead",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingRouterRouteV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	osMutexKV.Lock(routerID)
	defer osMutexKV.Unlock(routerID)

	n, err := routers.Get(networkingClient, routerID).Extract()
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router: %s", err)
	}

	routes := n.Routes
	dstCIDR := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)
	exists := false

	for _, route := range routes {
		if route.DestinationCIDR == dstCIDR && route.NextHop == nextHop {
			exists = true
			break
		}
	}

	if exists {
		log.Printf("[DEBUG] OpenTelekomCloud Neutron Router %s already has route to %s via %s", routerID, dstCIDR, nextHop)
		return resourceNetworkingRouterRouteV2Read(ctx, d, meta)
	}

	routes = append(routes, routers.Route{
		DestinationCIDR: dstCIDR,
		NextHop:         nextHop,
	})
	updateOpts := routers.UpdateOpts{
		Routes: routes,
	}
	log.Printf("[DEBUG] OpenTelekomCloud Neutron Router %s update options: %#v", routerID, updateOpts)
	_, err = routers.Update(networkingClient, routerID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating OpenTelekomCloud Neutron Router: %s", err)
	}

	d.SetId(fmt.Sprintf("%s-route-%s-%s", routerID, dstCIDR, nextHop))

	return resourceNetworkingRouterRouteV2Read(ctx, d, meta)
}

func resourceNetworkingRouterRouteV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	routerId := d.Get("router_id").(string)

	n, err := routers.Get(networkingClient, routerId).Extract()
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Router %s: %+v", routerId, n)

	destCidr := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
	)

	for _, r := range n.Routes {
		if r.DestinationCIDR == destCidr && r.NextHop == nextHop {
			mErr = multierror.Append(mErr,
				d.Set("destination_cidr", destCidr),
				d.Set("next_hop", nextHop),
			)
			break
		}
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetworkingRouterRouteV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	osMutexKV.Lock(routerID)
	defer osMutexKV.Unlock(routerID)

	n, err := routers.Get(networkingClient, routerID).Extract()
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router: %s", err)
	}

	dstCIDR := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)

	oldRoutes := n.Routes
	var newRoute []routers.Route

	for _, route := range oldRoutes {
		if route.DestinationCIDR != dstCIDR || route.NextHop != nextHop {
			newRoute = append(newRoute, route)
		}
	}

	if len(oldRoutes) == len(newRoute) {
		return diag.Errorf("Can't find route to %s via %s on OpenTelekomCloud Neutron Router %s", dstCIDR, nextHop, routerID)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud Neutron Router %s route to %s via %s", routerID, dstCIDR, nextHop)
	updateOpts := routers.UpdateOpts{
		Routes: newRoute,
	}
	_, err = routers.Update(networkingClient, routerID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating OpenTelekomCloud Neutron Router: %s", err)
	}

	return nil
}
