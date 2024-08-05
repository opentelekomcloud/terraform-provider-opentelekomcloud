package er

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/propagation"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceErPropagationV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePropagationV3Create,
		ReadContext:   resourcePropagationV3Read,
		DeleteContext: resourcePropagationV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePropagationImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attachment_type": {
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

func resourcePropagationV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		instanceId   = d.Get("instance_id").(string)
		routeTableId = d.Get("route_table_id").(string)

		opts = propagation.CreateOpts{
			RouterID:     instanceId,
			RouteTableID: routeTableId,
			AttachmentID: d.Get("attachment_id").(string),
		}
	)

	resp, err := propagation.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating the propagation to the route table: %s", err)
	}
	d.SetId(resp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      propagationStatusRefreshFunc(client, instanceId, routeTableId, d.Id(), []string{"available"}),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePropagationV3Read(ctx, d, meta)
}

func QueryPropagationById(client *golangsdk.ServiceClient, instanceId, routeTableId,
	propagationId string) (*propagation.Propagation, error) {
	resp, err := propagation.List(client, propagation.ListOpts{
		RouterId:     instanceId,
		RouteTableId: routeTableId,
	})
	if err != nil {
		return nil, err
	}

	filter := map[string]interface{}{
		"ID": propagationId,
	}
	result, err := common.FilterSliceWithField(resp.Propagations, filter)
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte(fmt.Sprintf("the propagation (%s) does not exist", propagationId)),
			},
		}
	}

	log.Printf("[DEBUG] The result filtered by resource ID (%s) is: %#v", propagationId, result)
	association, ok := result[0].(propagation.Propagation)
	if !ok {
		return nil, fmt.Errorf("the element type of filter result is incorrect, want 'propagations.Propagation', but got '%T'", result[0])
	}

	return &association, nil
}

func propagationStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, routeTableId, propagationId string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := QueryPropagationById(client, instanceId, routeTableId, propagationId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return resp, "COMPLETED", nil
			}

			return nil, "", err
		}

		if common.StrSliceContains([]string{"failed"}, resp.State) {
			return resp, "", fmt.Errorf("unexpected status '%s'", resp.State)
		}
		if common.StrSliceContains(targets, resp.State) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourcePropagationV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		propagationId = d.Id()
	)

	resp, err := QueryPropagationById(client, instanceId, routeTableId, propagationId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER propagation")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("route_table_id", resp.RouteTableID),
		d.Set("attachment_id", resp.AttachmentID),
		d.Set("attachment_type", resp.ResourceType),
		d.Set("status", resp.State),
		d.Set("created_at", resp.CreatedAt),
		d.Set("updated_at", resp.UpdatedAt),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving propagation (%s) fields: %s", propagationId, mErr)
	}
	return nil
}

func resourcePropagationV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var (
		instanceId    = d.Get("instance_id").(string)
		routeTableId  = d.Get("route_table_id").(string)
		propagationId = d.Id()

		opts = propagation.DeleteOpts{
			RouterID:     instanceId,
			RouteTableID: routeTableId,
			AttachmentID: d.Get("attachment_id").(string),
		}
	)
	err = propagation.Delete(client, opts)
	if err != nil {
		return diag.Errorf("error deleting propagation (%s): %s", propagationId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      propagationStatusRefreshFunc(client, instanceId, routeTableId, propagationId, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePropagationImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<route_table_id>/<propagation_id>', but got '%s'", d.Id())
	}

	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("route_table_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
