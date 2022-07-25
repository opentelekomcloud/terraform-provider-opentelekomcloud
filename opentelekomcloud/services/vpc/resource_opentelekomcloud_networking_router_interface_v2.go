package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNetworkingRouterInterfaceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterInterfaceV2Create,
		ReadContext:   resourceNetworkingRouterInterfaceV2Read,
		DeleteContext: resourceNetworkingRouterInterfaceV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingRouterInterfaceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	createOpts := routers.AddInterfaceOpts{
		SubnetID: d.Get("subnet_id").(string),
		PortID:   d.Get("port_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := routers.AddInterface(networkingClient, d.Get("router_id").(string), createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Neutron router interface: %s", err)
	}
	log.Printf("[INFO] Router interface Port ID: %s", n.PortID)

	log.Printf("[DEBUG] Waiting for Router Interface (%s) to become available", n.PortID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    waitForRouterInterfaceActive(networkingClient, n.PortID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for router interface to become active: %w", err)
	}

	d.SetId(n.PortID)

	return resourceNetworkingRouterInterfaceV2Read(ctx, d, meta)
}

func resourceNetworkingRouterInterfaceV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	n, err := ports.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Neutron Router Interface: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Router Interface %s: %+v", d.Id(), n)

	_ = d.Set("region", config.GetRegion(d))

	return nil
}

func resourceNetworkingRouterInterfaceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForRouterInterfaceDelete(networkingClient, d),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Neutron Router Interface: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForRouterInterfaceActive(networkingClient *golangsdk.ServiceClient, rId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := ports.Get(networkingClient, rId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Neutron Router Interface: %+v", r)
		return r, r.Status, nil
	}
}

func waitForRouterInterfaceDelete(networkingClient *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		routerId := d.Get("router_id").(string)
		routerInterfaceId := d.Id()

		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Router Interface %s.", routerInterfaceId)

		removeOpts := routers.RemoveInterfaceOpts{
			SubnetID: d.Get("subnet_id").(string),
			PortID:   d.Get("port_id").(string),
		}

		r, err := ports.Get(networkingClient, routerInterfaceId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Router Interface %s", routerInterfaceId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		_, err = routers.RemoveInterface(networkingClient, routerId, removeOpts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Router Interface %s.", routerInterfaceId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrDefault409); ok {
				if errCode.Actual == 409 {
					log.Printf("[DEBUG] Router Interface %s is still in use.", routerInterfaceId)
					return r, "ACTIVE", nil
				}
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					log.Printf("[DEBUG] Router Interface %s is still in use.", routerInterfaceId)
					return r, "ACTIVE", nil
				}
			}

			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Router Interface %s is still active.", routerInterfaceId)
		return r, "ACTIVE", nil
	}
}
