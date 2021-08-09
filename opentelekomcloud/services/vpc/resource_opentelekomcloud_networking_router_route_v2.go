package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNetworkingRouterRouteV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterRouteV2Create,
		ReadContext:   resourceNetworkingRouterRouteV2Read,
		DeleteContext: resourceNetworkingRouterRouteV2Delete,

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
	routerId := d.Get("router_id").(string)
	osMutexKV.Lock(routerId)
	defer osMutexKV.Unlock(routerId)

	var destCidr string = d.Get("destination_cidr").(string)
	var nextHop string = d.Get("next_hop").(string)

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	n, err := routers.Get(networkingClient, routerId).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router: %s", err)
	}

	var updateOpts routers.UpdateOpts
	var routeExists bool = false

	var rts []routers.Route = n.Routes
	for _, r := range rts {
		if r.DestinationCIDR == destCidr && r.NextHop == nextHop {
			routeExists = true
			break
		}
	}

	if !routeExists {
		if destCidr != "" && nextHop != "" {
			r := routers.Route{DestinationCIDR: destCidr, NextHop: nextHop}
			log.Printf(
				"[INFO] Adding route %s", r)
			rts = append(rts, r)
		}

		updateOpts.Routes = rts

		log.Printf("[DEBUG] Updating Router %s with options: %+v", routerId, updateOpts)

		_, err = routers.Update(networkingClient, routerId, updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud Neutron Router: %s", err)
		}
		d.SetId(fmt.Sprintf("%s-route-%s-%s", routerId, destCidr, nextHop))
	} else {
		log.Printf("[DEBUG] Router %s has route already", routerId)
	}

	return resourceNetworkingRouterRouteV2Read(ctx, d, meta)
}

func resourceNetworkingRouterRouteV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	routerId := d.Get("router_id").(string)

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	n, err := routers.Get(networkingClient, routerId).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Router %s: %+v", routerId, n)

	var destCidr string = d.Get("destination_cidr").(string)
	var nextHop string = d.Get("next_hop").(string)

	mErr := multierror.Append(
		d.Set("next_hop", ""),
		d.Set("destination_cidr", ""),
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
	routerId := d.Get("router_id").(string)
	osMutexKV.Lock(routerId)
	defer osMutexKV.Unlock(routerId)

	config := meta.(*cfg.Config)

	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	n, err := routers.Get(networkingClient, routerId).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router: %s", err)
	}

	var updateOpts routers.UpdateOpts

	var destCidr string = d.Get("destination_cidr").(string)
	var nextHop string = d.Get("next_hop").(string)

	var oldRts []routers.Route = n.Routes
	var newRts []routers.Route

	for _, r := range oldRts {
		if r.DestinationCIDR != destCidr || r.NextHop != nextHop {
			newRts = append(newRts, r)
		}
	}

	if len(oldRts) != len(newRts) {
		r := routers.Route{DestinationCIDR: destCidr, NextHop: nextHop}
		log.Printf(
			"[INFO] Deleting route %s", r)
		updateOpts.Routes = newRts

		log.Printf("[DEBUG] Updating Router %s with options: %+v", routerId, updateOpts)

		_, err = routers.Update(networkingClient, routerId, updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud Neutron Router: %s", err)
		}
	} else {
		return fmterr.Errorf("route did not exist already")
	}

	return nil
}
