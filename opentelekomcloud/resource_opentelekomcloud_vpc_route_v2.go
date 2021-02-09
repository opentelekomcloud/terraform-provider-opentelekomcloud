package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/routes"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func resourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcRouteV2Create,
		Read:   resourceVpcRouteV2Read,
		Delete: resourceVpcRouteV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nexthop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCIDR,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVpcRouteV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	createOpts := routes.CreateOpts{
		Type:        d.Get("type").(string),
		NextHop:     d.Get("nexthop").(string),
		Destination: d.Get("destination").(string),
		Tenant_Id:   d.Get("tenant_id").(string),
		VPC_ID:      d.Get("vpc_id").(string),
	}

	route, err := routes.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud VPC route: %s", err)
	}

	log.Printf("[INFO] Vpc Route ID: %s", route.RouteID)
	d.SetId(route.RouteID)

	return resourceVpcRouteV2Read(d, meta)

}

func resourceVpcRouteV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	route, err := routes.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error retrieving OpenTelekomCloud Vpc route: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("type", route.Type),
		d.Set("nexthop", route.NextHop),
		d.Set("destination", route.Destination),
		d.Set("tenant_id", route.Tenant_Id),
		d.Set("vpc_id", route.VPC_ID),
		d.Set("region", GetRegion(d, config)),
	)

	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func resourceVpcRouteV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}

	if err = routes.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("error deleting VPC route: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcRouteDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud VPC route: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcRouteDelete(client *golangsdk.ServiceClient, routeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		route, err := routes.Get(client, routeID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc route %s", routeID)
				return route, "DELETED", nil
			}
			return route, "ACTIVE", err
		}
		return route, "ACTIVE", nil
	}
}
