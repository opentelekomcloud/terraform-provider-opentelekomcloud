package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	throttlingpolicy "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/tr_policy"
	specialpolicy "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/tr_specials"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceAPIThrottlingPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGWThrottlingPolicyV2Create,
		ReadContext:   resourceAPIGWThrottlingPolicyV2Read,
		UpdateContext: resourceAPIGWThrottlingPolicyV2Update,
		DeleteContext: resourceAPIGWThrottlingPolicyV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceThrottlingPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_api_requests": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_app_requests": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_ip_requests": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_user_requests": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(PolicyTypeExclusive),
				ValidateFunc: validation.StringInSlice([]string{
					string(PolicyTypeExclusive),
					string(PolicyTypeShared),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(PeriodUnitMinute),
				ValidateFunc: validation.StringInSlice([]string{
					string(PeriodUnitSecond),
					string(PeriodUnitMinute),
					string(PeriodUnitHour),
					string(PeriodUnitDay),
				}, false),
			},
			"user_throttles": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 30,
				Elem:     specialThrottleSchemaResource(),
			},
			"app_throttles": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 30,
				Elem:     specialThrottleSchemaResource(),
			},
			"created_at": {
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

func specialThrottleSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_api_requests": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"throttling_object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"throttling_object_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildThrottlingPolicyCreateOpts(d *schema.ResourceData) (throttlingpolicy.CreateOpts, error) {
	opt := throttlingpolicy.CreateOpts{
		Name:           d.Get("name").(string),
		TimeInterval:   pointerto.Int(d.Get("period").(int)),
		TimeUnit:       d.Get("period_unit").(string),
		ApiCallLimits:  pointerto.Int(d.Get("max_api_requests").(int)),
		UserCallLimits: pointerto.Int(d.Get("max_user_requests").(int)),
		AppCallLimits:  pointerto.Int(d.Get("max_app_requests").(int)),
		IpCallLimits:   pointerto.Int(d.Get("max_ip_requests").(int)),
		Description:    d.Get("description").(string),
		GatewayID:      d.Get("instance_id").(string),
	}
	pType := d.Get("type").(string)
	var val int
	var ok bool
	if val, ok = policyType[pType]; !ok {
		return opt, fmt.Errorf("invalid throttling policy type: %s", pType)
	}
	opt.Type = pointerto.Int(val)
	return opt, nil
}

func buildThrottlingPolicyUpdateOpts(d *schema.ResourceData) (throttlingpolicy.UpdateOpts, error) {
	opt := throttlingpolicy.UpdateOpts{
		Name:           d.Get("name").(string),
		TimeInterval:   pointerto.Int(d.Get("period").(int)),
		TimeUnit:       d.Get("period_unit").(string),
		ApiCallLimits:  pointerto.Int(d.Get("max_api_requests").(int)),
		UserCallLimits: pointerto.Int(d.Get("max_user_requests").(int)),
		AppCallLimits:  pointerto.Int(d.Get("max_app_requests").(int)),
		IpCallLimits:   pointerto.Int(d.Get("max_ip_requests").(int)),
		Description:    d.Get("description").(string),
		GatewayID:      d.Get("instance_id").(string),
		ThrottleID:     d.Id(),
	}
	pType := d.Get("type").(string)
	var val int
	var ok bool
	if val, ok = policyType[pType]; !ok {
		return opt, fmt.Errorf("invalid throttling policy type: %s", pType)
	}
	opt.Type = pointerto.Int(val)
	return opt, nil
}

func addSpecThrottlingPolicies(client *golangsdk.ServiceClient, policies *schema.Set,
	instanceId, policyId, specType string) error {
	for _, policy := range policies.List() {
		raw := policy.(map[string]interface{})
		specOpts := specialpolicy.CreateOpts{
			ObjectType: specType,
			ObjectID:   raw["throttling_object_id"].(string),
			CallLimits: raw["max_api_requests"].(int),
			GatewayID:  instanceId,
			ThrottleID: policyId,
		}
		_, err := specialpolicy.Create(client, specOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAPIGWThrottlingPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIGW v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

	opts, err := buildThrottlingPolicyCreateOpts(d)
	if err != nil {
		return diag.Errorf("unable to get the create option of the throttling policy: %s", err)
	}
	resp, err := throttlingpolicy.Create(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "throttling policy")
	}
	d.SetId(resp.ID)

	if policies, ok := d.GetOk("user_throttles"); ok {
		err := addSpecThrottlingPolicies(client, policies.(*schema.Set), instanceId, d.Id(), string(PolicyTypeUser))
		if err != nil {
			return diag.Errorf("error creating special user throttling policy: %s", err)
		}
	}
	if policies, ok := d.GetOk("app_throttles"); ok {
		err := addSpecThrottlingPolicies(client, policies.(*schema.Set), instanceId, d.Id(), string(PolicyTypeApplication))
		if err != nil {
			return diag.Errorf("error creating special application throttling policy: %s", err)
		}
	}
	return resourceAPIGWThrottlingPolicyV2Read(ctx, d, meta)
}

func analyseThrottlingPolicyType(pType int) *string {
	for k, v := range policyType {
		if v == pType {
			return &k
		}
	}
	return nil
}

func flattenSpecThrottlingPolicies(specThrottles []specialpolicy.ThrottlingResp) (userThrottles,
	appThrottles []map[string]interface{}, err error) {
	if len(specThrottles) == 0 {
		return nil, nil, nil
	}

	users := make([]map[string]interface{}, 0)
	apps := make([]map[string]interface{}, 0)

	for _, throttle := range specThrottles {
		switch throttle.ObjectType {
		case string(PolicyTypeApplication):
			apps = append(apps, map[string]interface{}{
				"max_api_requests":       throttle.CallLimits,
				"throttling_object_id":   throttle.ObjectID,
				"throttling_object_name": throttle.ObjectName,
				"id":                     throttle.ID,
			})
		case string(PolicyTypeUser):
			users = append(users, map[string]interface{}{
				"max_api_requests":       throttle.CallLimits,
				"throttling_object_id":   throttle.ObjectID,
				"throttling_object_name": throttle.ObjectName,
				"id":                     throttle.ID,
			})
		default:
			return users, apps, fmt.Errorf("invalid policy type, want '%v' or '%v', but '%v'", PolicyTypeApplication,
				PolicyTypeUser, throttle.ObjectType)
		}
	}

	return users, apps, nil
}

func resourceAPIGWThrottlingPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Id()
	)
	resp, err := throttlingpolicy.Get(client, instanceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "throttling policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("type", analyseThrottlingPolicyType(resp.Type)),
		d.Set("name", resp.Name),
		d.Set("period", resp.TimeInterval),
		d.Set("period_unit", resp.TimeUnit),
		d.Set("max_api_requests", resp.ApiCallLimits),
		d.Set("max_user_requests", resp.UserCallLimits),
		d.Set("max_app_requests", resp.AppCallLimits),
		d.Set("max_ip_requests", resp.IpCallLimits),
		d.Set("description", resp.Description),
		// Attributes
		d.Set("created_at", resp.CreateTime),
	)

	if resp.IsIncluSpecialThrottle == includeSpecialThrottle {
		specResp, err := specialpolicy.List(client, specialpolicy.ListOpts{
			GatewayID:  instanceId,
			ThrottleID: d.Id(),
		})
		if err != nil {
			return diag.Errorf("error retrieving special throttle: %s", err)
		}
		userThrottles, appThrottles, err := flattenSpecThrottlingPolicies(specResp)
		if err != nil {
			return diag.Errorf("error retrieving special throttle: %s", err)
		}
		mErr = multierror.Append(mErr,
			d.Set("user_throttles", userThrottles),
			d.Set("app_throttles", appThrottles),
		)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving throttling policy (%s) fields: %s", policyId, err)
	}
	return nil
}

func updateSpecThrottlingPolicyCallLimit(client *golangsdk.ServiceClient, instanceId, policyId, strategyId string,
	limit int) error {
	opts := specialpolicy.UpdateOpts{
		CallLimits:      limit,
		GatewayID:       instanceId,
		ThrottleID:      policyId,
		SpecialPolicyID: strategyId,
	}
	_, err := specialpolicy.Update(client, opts)
	if err != nil {
		return err
	}
	return nil
}

func removeSpecThrottlingPolicies(client *golangsdk.ServiceClient, policies *schema.Set,
	instanceId, policyId string) error {
	for _, policy := range policies.List() {
		raw := policy.(map[string]interface{})
		err := specialpolicy.Delete(client, instanceId, policyId, raw["id"].(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func updateSpecThrottlingPolicies(d *schema.ResourceData, client *golangsdk.ServiceClient,
	paramName, specType string) error {
	oldRaws, newRaws := d.GetChange(paramName)
	addRaws := newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
	removeRaws := oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	instanceId := d.Get("instance_id").(string)

	for _, rm := range removeRaws.List() {
		rmPolicy := rm.(map[string]interface{})
		rmObject := rmPolicy["throttling_object_id"].(string)
		for _, add := range addRaws.List() {
			addPolicy := add.(map[string]interface{})
			if rmObject == addPolicy["throttling_object_id"].(string) {
				strategyId := rmPolicy["id"].(string)
				limit := addPolicy["max_api_requests"].(int)
				err := updateSpecThrottlingPolicyCallLimit(client, instanceId, d.Id(), strategyId, limit)
				if err != nil {
					return err
				}
				removeRaws.Remove(rm)
				addRaws.Remove(add)
			}
		}
	}
	err := removeSpecThrottlingPolicies(client, removeRaws, instanceId, d.Id())
	if err != nil {
		return err
	}
	err = addSpecThrottlingPolicies(client, addRaws, instanceId, d.Id(), specType)
	if err != nil {
		return err
	}
	return nil
}

func resourceAPIGWThrottlingPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	if d.HasChangesExcept("user_throttles", "app_throttles") {
		opts, err := buildThrottlingPolicyUpdateOpts(d)
		if err != nil {
			return diag.Errorf("unable to get the update option of the throttling policy: %s", err)
		}

		_, err = throttlingpolicy.Update(client, opts)
		if err != nil {
			return diag.Errorf("error updating throttling policy: %s", err)
		}
	}
	if d.HasChange("user_throttles") {
		err = updateSpecThrottlingPolicies(d, client, "user_throttles", string(PolicyTypeUser))
		if err != nil {
			return diag.Errorf("error updating special user throttles: %s", err)
		}
	}
	if d.HasChange("app_throttles") {
		err = updateSpecThrottlingPolicies(d, client, "app_throttles", string(PolicyTypeApplication))
		if err != nil {
			return diag.Errorf("error updating special app throttles: %s", err)
		}
	}

	return resourceAPIGWThrottlingPolicyV2Read(ctx, d, meta)
}

func resourceAPIGWThrottlingPolicyV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Id()
	)
	if err = throttlingpolicy.Delete(client, instanceId, policyId); err != nil {
		return diag.Errorf("unable to delete the throttling policy (%s): %s", policyId, err)
	}

	return nil
}

func resourceThrottlingPolicyImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIGW v2 client: %s", err)
	}

	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}

	instanceId := parts[0]
	name := parts[1]
	opt := throttlingpolicy.ListOpts{
		PolicyName: name,
		GatewayID:  instanceId,
	}
	resp, err := throttlingpolicy.List(client, opt)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving throttling policies: %s", err)
	}

	d.SetId(resp[0].ID)

	return []*schema.ResourceData{d}, d.Set("instance_id", instanceId)
}
