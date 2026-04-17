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
	fl "github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/flow-logs"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceErFlowLogV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceErFlowLogCreateV3,
		UpdateContext: resourceErFlowLogUpdateV3,
		ReadContext:   resourceErFlowLogReadV3,
		DeleteContext: resourceErFlowLogDeleteV3,

		Importer: &schema.ResourceImporter{
			StateContext: resourceErFlowLogV3ImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_store_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
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
		},
	}
}

func resourceErFlowLogCreateV3(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}
	opts := fl.CreateOpts{
		RouterID: d.Get("instance_id").(string),
		FlowLog: &fl.FlowLog{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			ResourceType: d.Get("resource_type").(string),
			ResourceId:   d.Get("resource_id").(string),
			LogGroupId:   d.Get("log_group_id").(string),
			LogStreamId:  d.Get("log_stream_id").(string),
			LogStoreType: d.Get("log_store_type").(string),
		},
	}
	log, err := fl.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud ER V3 Flow Log: %s", err)
	}

	d.SetId(log.ID)

	err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the create of OpenTelekomCloud ER V3 flow log (%s) to complete: %s", d.Id(), err)
	}

	enabled := d.Get("enabled").(bool)
	if !enabled {
		_, err = fl.Disable(client, d.Get("instance_id").(string), d.Id())
		if err != nil {
			return diag.FromErr(err)
		}

		err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of OpenTelekomCloud ER V3 flow log (%s) to complete: %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceErFlowLogReadV3(clientCtx, d, meta)
}

func getFlowLogInfo(ctx context.Context, d *schema.ResourceData, meta interface{}) (*fl.FlowLogResponse, error) {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return nil, err
	}
	resp, err := fl.Get(client, d.Get("instance_id").(string), d.Id())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func resourceErFlowLogReadV3(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	resp, err := fl.Get(client, d.Get("instance_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "ER flow log")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("log_store_type", resp.LogStoreType),
		d.Set("log_group_id", resp.LogGroupId),
		d.Set("log_stream_id", resp.LogStreamId),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("resource_type", resp.ResourceType),
		d.Set("resource_id", resp.ResourceId),
		d.Set("state", resp.Status),
		d.Set("created_at", common.FormatTimeStampRFC3339(common.ConvertTimeStrToNanoTimestamp(resp.CreatedAt)/1000, false)),
		d.Set("updated_at", common.FormatTimeStampRFC3339(common.ConvertTimeStrToNanoTimestamp(resp.UpdatedAt)/1000, false)),
		d.Set("enabled", resp.Enabled),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceErFlowLogUpdateV3(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if d.HasChanges("name", "description") {
		_, err = fl.Update(client, d.Id(), fl.UpdateOpts{
			RouterID:    d.Get("instance_id").(string),
			Name:        d.Get("name").(string),
			Description: pointerto.String(d.Get("description").(string)),
		})
		if err != nil {
			return diag.Errorf("error updating OpenTelekomCloud ER V3 flow log: %s", err)
		}

		err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of OpenTelekomCloud ER V3 flow log (%s) to complete: %s", d.Id(), err)
		}
	}

	if d.HasChanges("enabled") {
		enabled := d.Get("enabled").(bool)
		if enabled {
			_, err = fl.Enable(client, d.Get("instance_id").(string), d.Id())
		} else {
			_, err = fl.Disable(client, d.Get("instance_id").(string), d.Id())
		}
		if err != nil {
			return diag.FromErr(err)
		}
		err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of OpenTelekomCloud ER V3 flow log (%s) to complete: %s", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, erClientV3)
	return resourceErFlowLogReadV3(clientCtx, d, meta)
}

func resourceErFlowLogDeleteV3(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, erClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.ErV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = fl.Delete(client, d.Get("instance_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("error deleting Flow Log: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      flowLogStatusRefreshFunc(ctx, d, meta, true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flowLogStatusRefreshFunc(ctx context.Context, d *schema.ResourceData, meta interface{}, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getFlowLogInfo(ctx, d, meta)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				return "Resource Not Found", "DELETED", nil
			}

			return nil, "ERROR", err
		}

		if common.StrSliceContains([]string{"fail"}, resp.Status) {
			return resp, "", fmt.Errorf("unexpected status: '%s'", resp.Status)
		}

		if common.StrSliceContains([]string{"available"}, resp.Status) {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func flowLogWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      flowLogStatusRefreshFunc(ctx, d, meta, false),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceErFlowLogV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
