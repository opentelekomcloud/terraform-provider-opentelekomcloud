package er

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/route"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceErStaticRouteV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStaticRouteV3Create,
		UpdateContext: resourceStaticRouteV3Update,
		ReadContext:   resourceStaticRouteV3Read,
		DeleteContext: resourceStaticRouteV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceStaticRouteV3ImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_blackhole": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildStaticRouteCreateOpts(d *schema.ResourceData) route.CreateOpts {
	return route.CreateOpts{
		RouteTableId: d.Get("route_table_id").(string),
		Destination:  d.Get("destination").(string),
		AttachmentId: d.Get("attachment_id").(string),
		IsBlackhole:  pointerto.Bool(d.Get("is_blackhole").(bool)),
	}
}

func staticRouteStatusRefreshFunc(client *golangsdk.ServiceClient, routeTableId, staticRouteId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := route.Get(client, routeTableId, staticRouteId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "NOT_FOUND", "COMPLETED", nil
			}
			return "AN_ERROR_OCCURRED", "ERROR", err
		}

		if common.IsStrContainsSliceElement(resp.State, targets, false, true) {
			return resp, "available", nil
		}
		return resp, "pending", nil
	}
}

func resourceStaticRouteV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	opts := buildStaticRouteCreateOpts(d)

	resp, err := route.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating static route: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      staticRouteStatusRefreshFunc(client, d.Get("route_table_id").(string), d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the create operation completed: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceStaticRouteV3Read(clientCtx, d, meta)
}

func resourceStaticRouteV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		routeTableId  = d.Get("route_table_id").(string)
		staticRouteId = d.Id()
	)

	resp, err := route.Get(client, routeTableId, staticRouteId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER static route")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("destination", resp.Destination),
		d.Set("is_blackhole", resp.IsBlackhole),
		d.Set("type", resp.Type),
		d.Set("status", resp.State),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)

	if len(resp.Attachments) > 0 && resp.Attachments[0].AttachmentId != "" {
		mErr = multierror.Append(mErr, d.Set("attachment_id", resp.Attachments[0].AttachmentId))
	} else {
		mErr = multierror.Append(mErr, d.Set("attachment_id", nil))
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving static route (%s) fields: %s", staticRouteId, mErr)
	}
	return nil
}

func resourceStaticRouteV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	opts := route.UpdateOpts{
		RouteTableId: d.Get("route_table_id").(string),
		RouteId:      d.Id(),
		AttachmentId: d.Get("attachment_id").(string),
		IsBlackhole:  pointerto.Bool(d.Get("is_blackhole").(bool)),
	}

	_, err = route.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating static route (%s): %s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"available"},
		Refresh:      staticRouteStatusRefreshFunc(client, d.Get("route_table_id").(string), d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the update operation completed: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceStaticRouteV3Read(clientCtx, d, meta)
}

func resourceStaticRouteV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		routeTableId  = d.Get("route_table_id").(string)
		staticRouteId = d.Id()
	)
	err = route.Delete(client, routeTableId, staticRouteId)
	if err != nil {
		return diag.Errorf("error deleting static route (%s): %s", staticRouteId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"COMPLETED"},
		Refresh:      staticRouteStatusRefreshFunc(client, routeTableId, d.Id(), nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the delete operation completed: %s", err)
	}
	return nil
}

func resourceStaticRouteV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid format for import ID, want '<route_table_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	if err := d.Set("route_table_id", parts[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}
