package fgs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const (
	triggerStatusActive   = "ACTIVE"
	triggerStatusDisabled = "DISABLED"
	httpUrl               = "fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}"
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
	var (
		config      = meta.(*cfg.Config)
		httpUrl     = "fgs/triggers/{function_urn}"
		functionUrn = d.Get("function_urn").(string)
	)

	client, err := common.ClientFromCtx(ctx, fgsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.FuncGraphV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{function_urn}", functionUrn)
	createOpt := golangsdk.RequestOpts{
		JSONBody: common.RemoveNil(buildCreateFunctionTriggerBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating function trigger: %s", err)
	}

	respBody, err := common.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId := common.PathSearch("trigger_id", respBody, "")
	d.SetId(resourceId.(string))

	clientCtx := common.CtxWithClient(ctx, client, fgsClientV2)
	return resourceFunctionTriggerV2Read(clientCtx, d, meta)
}

func buildCreateFunctionTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := d.Get("event_data").(string)
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &parseResult)
	if err != nil {
		log.Printf("[ERROR] Invalid type of the params, not json format")
	}
	return map[string]interface{}{
		"trigger_type_code": d.Get("type"),
		"trigger_status":    d.Get("status"),
		"event_data":        parseResult,
	}
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
		respBody, err := GetTriggerById(client, functionUrn, triggerType, triggerId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "Resource Not Found", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		status := common.PathSearch("trigger_status", respBody, "").(string)

		if common.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}
		return respBody, "PENDING", nil
	}
}

func GetTriggerById(client *golangsdk.ServiceClient, functionUrn, triggerType, triggerId string) (interface{}, error) {
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{function_urn}", functionUrn)
	getPath = strings.ReplaceAll(getPath, "{trigger_type_code}", triggerType)
	getPath = strings.ReplaceAll(getPath, "{trigger_id}", triggerId)

	requestResp, err := client.Request("GET", getPath, &golangsdk.RequestOpts{})
	if err != nil {
		return nil, parseTriggerQueryError(err)
	}
	return common.FlattenResponse(requestResp)
}

func parseTriggerQueryError(err error) error {
	var errCode golangsdk.ErrDefault500
	if errors.As(err, &errCode) {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return err
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return err
		}

		// Error code FSS.0500 indicates that the function to which the trigger belongs has been deleted.
		if errorCode == "FSS.0500" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return err
}

func parseEventData(eventData interface{}) interface{} {
	jsonEventData, err := json.Marshal(eventData)
	if err != nil {
		log.Printf("[ERROR] unable to convert the event data of the function trigger, not json format")
		return nil
	}
	return string(jsonEventData)
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

	respBody, err := GetTriggerById(client, functionUrn, triggerType, triggerId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Function trigger")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("type", common.PathSearch("trigger_type_code", respBody, nil)),
		d.Set("status", common.PathSearch("trigger_status", respBody, nil)),
		d.Set("event_data", parseEventData(common.PathSearch("event_data", respBody, nil))),
		d.Set("created_at", common.PathSearch("created_time", respBody, nil)),
		d.Set("updated_at", common.PathSearch("last_updated_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateFunctionTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := d.Get("event_data").(string)
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(params), &parseResult)
	if err != nil {
		log.Printf("[ERROR] Invalid type of the params, not json format")
	}
	return map[string]interface{}{
		"trigger_status": d.Get("status"),
		"event_data":     parseResult,
	}
}

func resourceFunctionTriggerV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", functionUrn)
	updatePath = strings.ReplaceAll(updatePath, "{trigger_type_code}", triggerType)
	updatePath = strings.ReplaceAll(updatePath, "{trigger_id}", triggerId)
	updateOpts := golangsdk.RequestOpts{
		JSONBody: common.RemoveNil(buildUpdateFunctionTriggerBodyParams(d)),
		OkCodes:  []int{200, 201, 202},
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error deleting function trigger: %s", err)
	}

	err = waitForFunctionTriggerStatusCompleted(ctx, client, d)
	if err != nil {
		diag.Errorf("error waiting for the function trigger (%s) status to become available: %s", triggerId, err)
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

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{function_urn}", functionUrn)
	deletePath = strings.ReplaceAll(deletePath, "{trigger_type_code}", triggerType)
	deletePath = strings.ReplaceAll(deletePath, "{trigger_id}", triggerId)
	deleteOpts := golangsdk.RequestOpts{
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
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
