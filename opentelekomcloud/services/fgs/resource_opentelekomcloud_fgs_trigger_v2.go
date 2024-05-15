package fgs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/fgs/v2/trigger"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const (
	triggerStatusActive   = "ACTIVE"
	triggerStatusDisabled = "DISABLED"
)

func ResourceFgsTriggerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionTriggerV2Create,
		ReadContext:   resourceFunctionTriggerV2Read,
		UpdateContext: resourceFunctionTriggerV2Update,
		DeleteContext: resourceFunctionTriggerV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceFunctionTriggerImportState,
		},

		Schema: map[string]*schema.Schema{
			"function_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"event_data": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					triggerStatusActive, triggerStatusDisabled,
				}, false),
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

func resourceFunctionTriggerV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpt, err := buildCreateFunctionTriggerBodyParams(d)
	if err != nil {
		return fmterr.Errorf("invalid type of event_data provided (has to be in json): %s", err)
	}

	triggerResp, err := trigger.Create(client, *createOpt)
	if err != nil {
		return diag.Errorf("error creating function trigger: %s", err)
	}

	d.SetId(triggerResp.TriggerId)

	clientCtx := common.CtxWithClient(ctx, client, fgsClientV2)
	return resourceFunctionTriggerV2Read(clientCtx, d, meta)
}

func buildCreateFunctionTriggerBodyParams(d *schema.ResourceData) (*trigger.CreateOpts, error) {
	params := d.Get("event_data").(string)
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &parseResult)
	if err != nil {
		return nil, err
	}
	parseResult = common.RemoveNil(parseResult)

	return &trigger.CreateOpts{
		FuncUrn:         d.Get("function_urn").(string),
		TriggerTypeCode: d.Get("type").(string),
		TriggerStatus:   d.Get("status").(string),
		EventData:       parseResult,
	}, nil
}

func waitForFunctionTriggerStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      functionTriggerStatusRefreshFunc(client, d, []string{"ACTIVE", "DISABLED"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func functionTriggerStatusRefreshFunc(client *golangsdk.ServiceClient, d *schema.ResourceData, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			functionUrn = d.Get("function_urn").(string)
			triggerType = d.Get("type").(string)
			triggerId   = d.Id()
		)
		respBody, err := trigger.Get(client, functionUrn, triggerType, triggerId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "Resource Not Found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		if common.StrSliceContains(targets, respBody.TriggerStatus) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func resourceFunctionTriggerV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		config      = meta.(*cfg.Config)
		region      = config.GetRegion(d)
		functionUrn = d.Get("function_urn").(string)
		triggerType = d.Get("type").(string)
		triggerId   = d.Id()
	)

	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	getTrigger, err := trigger.Get(client, functionUrn, triggerType, triggerId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Function trigger")
	}

	eventData, err := flattenEventData(getTrigger.EventData)
	if err != nil {
		return fmterr.Errorf("Error converting event data response")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("type", getTrigger.TriggerTypeCode),
		d.Set("status", getTrigger.TriggerStatus),
		d.Set("event_data", eventData),
		d.Set("created_at", getTrigger.CreatedTime),
		d.Set("updated_at", getTrigger.LastUpdatedTime),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateFunctionTriggerBodyParams(d *schema.ResourceData) (*trigger.UpdateOpts, error) {
	params := d.Get("event_data").(string)
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &parseResult)
	if err != nil {
		return nil, err
	}
	parseResult = common.RemoveNil(parseResult)

	return &trigger.UpdateOpts{
		TriggerId:       d.Id(),
		FuncUrn:         d.Get("function_urn").(string),
		TriggerTypeCode: d.Get("type").(string),
		TriggerStatus:   d.Get("status").(string),
		EventData:       parseResult,
	}, nil
}

func resourceFunctionTriggerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	updateOpts, err := buildUpdateFunctionTriggerBodyParams(d)
	if err != nil {
		return fmterr.Errorf("invalid type of event_data provided (has to be in json): %s", err)
	}

	_, err = trigger.Update(client, *updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating trigger event: %s", err)
	}

	err = waitForFunctionTriggerStatusCompleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the function trigger (%s) status to become available: %s", d.Id(), err)
	}
	return nil
}

func resourceFunctionTriggerV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		config      = meta.(*cfg.Config)
		functionUrn = d.Get("function_urn").(string)
		triggerType = d.Get("type").(string)
		triggerId   = d.Id()
	)

	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	err = trigger.Delete(client, functionUrn, triggerType, triggerId)
	if err != nil {
		return diag.Errorf("error deleting function trigger: %s", err)
	}

	err = waitForFunctionTriggerDeleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the function trigger (%s) status to become deleted: %s", triggerId, err)
	}
	return nil
}

func waitForFunctionTriggerDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      functionTriggerStatusRefreshFunc(client, d, nil),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceFunctionTriggerImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	var (
		importId = d.Id()
		parts    = strings.Split(importId, "/")
	)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid resource ID format for function trigger, want '<function_urn>/<type>/<id>', but got '%s'", importId)
	}
	d.SetId(parts[2])
	mErr := multierror.Append(
		d.Set("function_urn", parts[0]),
		d.Set("type", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func flattenEventData(resp map[string]interface{}) (string, error) {
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	jsonString := string(jsonBytes)

	return jsonString, err
}
