package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/routes"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVPCRouteV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcRouteV2Create,
		ReadContext:   resourceVpcRouteV2Read,
		DeleteContext: resourceVpcRouteV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: common.ValidateCIDR,
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

func resourceVpcRouteV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
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
		return fmterr.Errorf("error creating OpenTelekomCloud VPC route: %s", err)
	}

	log.Printf("[INFO] Vpc Route ID: %s", route.RouteID)
	d.SetId(route.RouteID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceVpcRouteV2Read(clientCtx, d, meta)
}

func resourceVpcRouteV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	route, err := routes.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "route")
	}

	mErr := multierror.Append(nil,
		d.Set("type", route.Type),
		d.Set("nexthop", route.NextHop),
		d.Set("destination", route.Destination),
		d.Set("tenant_id", route.Tenant_Id),
		d.Set("vpc_id", route.VPC_ID),
		d.Set("region", config.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceVpcRouteV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	if err := routes.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting VPC route: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"DELETED"},
		Refresh:      waitForVpcRouteDelete(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud VPC route: %s", err)
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
