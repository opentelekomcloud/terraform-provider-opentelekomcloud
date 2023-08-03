package as

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceASPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASPolicyCreate,
		ReadContext:   resourceASPolicyRead,
		UpdateContext: resourceASPolicyUpdate,
		DeleteContext: resourceASPolicyDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_policy_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: resourceASPolicyValidateName,
				ForceNew:     false,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: resourceASPolicyValidatePolicyType,
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"scheduled_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"launch_time": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"recurrence_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     false,
							ValidateFunc: resourceASPolicyValidateRecurrenceType,
						},
						"recurrence_value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"start_time": {
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         false,
							Default:          getCurrentUTCwithoutSec(),
							DiffSuppressFunc: common.SuppressDiffAll,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
			"scaling_policy_action": {
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: resourceASPolicyValidateActionOperation,
						},
						"instance_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
			"cool_down_time": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
				Default:  900,
			},
		},
	}
}

func getCurrentUTCwithoutSec() string {
	utcTime := time.Now().UTC().Format(time.RFC3339)
	splits := strings.SplitN(utcTime, ":", 3)
	resultTime := strings.Join(splits[0:2], ":") + "Z"
	return resultTime
}

func validateParameters(d *schema.ResourceData) error {
	log.Printf("[DEBUG] validateParameters for as policy!")
	policyType := d.Get("scaling_policy_type").(string)
	alarmId := d.Get("alarm_id").(string)
	log.Printf("[DEBUG] validateParameters alarmId is :%s", alarmId)
	log.Printf("[DEBUG] validateParameters policyType is :%s", policyType)
	scheduledPolicy := d.Get("scheduled_policy").([]interface{})
	log.Printf("[DEBUG] validateParameters scheduledPolicy is :%#v", scheduledPolicy)
	if policyType == "ALARM" {
		if alarmId == "" {
			return fmt.Errorf("parameter alarm_id should be set if policy type is ALARM")
		}
	}
	if policyType == "SCHEDULED" || policyType == "RECURRENCE" {
		if len(scheduledPolicy) == 0 {
			return fmt.Errorf("parameter scheduled_policy should be set if policy type is RECURRENCE or SCHEDULED")
		}
	}

	if len(scheduledPolicy) == 1 {
		scheduledPolicyMap := scheduledPolicy[0].(map[string]interface{})
		log.Printf("[DEBUG] validateParameters scheduledPolicyMap is :%#v", scheduledPolicyMap)
		recurrenceType := scheduledPolicyMap["recurrence_type"].(string)
		endTime := scheduledPolicyMap["end_time"].(string)
		log.Printf("[DEBUG] validateParameters recurrenceType is :%#v", recurrenceType)
		if policyType == "RECURRENCE" {
			if recurrenceType == "" {
				return fmt.Errorf("parameter recurrence_type should be set if policy type is RECURRENCE")
			}
			if endTime == "" {
				return fmt.Errorf("parameter end_time should be set if policy type is RECURRENCE")
			}
		}
	}

	return nil
}

func getScheduledPolicy(rawScheduledPolicy map[string]interface{}) policies.SchedulePolicyOpts {
	scheduledPolicy := policies.SchedulePolicyOpts{
		LaunchTime:      rawScheduledPolicy["launch_time"].(string),
		RecurrenceType:  rawScheduledPolicy["recurrence_type"].(string),
		RecurrenceValue: rawScheduledPolicy["recurrence_value"].(string),
		StartTime:       rawScheduledPolicy["start_time"].(string),
		EndTime:         rawScheduledPolicy["end_time"].(string),
	}
	return scheduledPolicy
}

func getPolicyAction(rawPolicyAction map[string]interface{}) policies.Action {
	policyAction := policies.Action{
		Operation:   rawPolicyAction["operation"].(string),
		InstanceNum: rawPolicyAction["instance_number"].(int),
	}
	return policyAction
}

func resourceASPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	err = validateParameters(d)
	if err != nil {
		return fmterr.Errorf("error creating ASPolicy: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:         d.Get("scaling_policy_name").(string),
		ID:           d.Get("scaling_group_id").(string),
		Type:         d.Get("scaling_policy_type").(string),
		AlarmID:      d.Get("alarm_id").(string),
		CoolDownTime: d.Get("cool_down_time").(int),
	}
	scheduledPolicyList := d.Get("scheduled_policy").([]interface{})
	if len(scheduledPolicyList) == 1 {
		scheduledPolicyMap := scheduledPolicyList[0].(map[string]interface{})
		scheduledPolicy := getScheduledPolicy(scheduledPolicyMap)
		createOpts.SchedulePolicy = scheduledPolicy
	}
	policyActionList := d.Get("scaling_policy_action").([]interface{})
	if len(policyActionList) == 1 {
		policyActionMap := policyActionList[0].(map[string]interface{})
		policyAction := getPolicyAction(policyActionMap)
		createOpts.Action = policyAction
	}

	log.Printf("[DEBUG] Create AS policy Options: %#v", createOpts)
	asPolicyId, err := policies.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating ASPolicy: %s", err)
	}
	d.SetId(asPolicyId)
	log.Printf("[DEBUG] Create AS Policy %q Success!", asPolicyId)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceASPolicyRead(clientCtx, d, meta)
}

func resourceASPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	asPolicy, err := policies.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS Policy")
	}

	log.Printf("[DEBUG] Retrieved ASPolicy %q: %+v", d.Id(), asPolicy)
	mErr := multierror.Append(
		d.Set("scaling_policy_name", asPolicy.Name),
		d.Set("scaling_policy_type", asPolicy.Type),
		d.Set("alarm_id", asPolicy.AlarmID),
		d.Set("cool_down_time", asPolicy.CoolDownTime),
		d.Set("region", config.GetRegion(d)),
	)

	policyActionInfo := asPolicy.Action
	policyAction := map[string]interface{}{}
	policyAction["operation"] = policyActionInfo.Operation
	policyAction["instance_number"] = policyActionInfo.InstanceNum
	policyActionList := []interface{}{policyAction}
	mErr = multierror.Append(mErr, d.Set("scaling_policy_action", policyActionList))

	scheduledPolicyInfo := asPolicy.SchedulePolicy
	scheduledPolicy := map[string]interface{}{}
	scheduledPolicy["launch_time"] = scheduledPolicyInfo.LaunchTime
	scheduledPolicy["recurrence_type"] = scheduledPolicyInfo.RecurrenceType
	scheduledPolicy["recurrence_value"] = scheduledPolicyInfo.RecurrenceValue
	scheduledPolicy["start_time"] = scheduledPolicyInfo.StartTime
	scheduledPolicy["end_time"] = scheduledPolicyInfo.EndTime
	scheduledPolicies := []interface{}{scheduledPolicy}
	mErr = multierror.Append(mErr, d.Set("scheduled_policy", scheduledPolicies))

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceASPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	err = validateParameters(d)
	if err != nil {
		return fmterr.Errorf("error updating ASPolicy: %s", err)
	}
	updateOpts := policies.UpdateOpts{
		Name:         d.Get("scaling_policy_name").(string),
		Type:         d.Get("scaling_policy_type").(string),
		AlarmID:      d.Get("alarm_id").(string),
		CoolDownTime: d.Get("cool_down_time").(int),
	}
	scheduledPolicyList := d.Get("scheduled_policy").([]interface{})
	if len(scheduledPolicyList) == 1 {
		scheduledPolicyMap := scheduledPolicyList[0].(map[string]interface{})
		scheduledPolicy := getScheduledPolicy(scheduledPolicyMap)
		updateOpts.SchedulePolicy = scheduledPolicy
	}
	policyActionList := d.Get("scaling_policy_action").([]interface{})
	if len(policyActionList) == 1 {
		policyActionMap := policyActionList[0].(map[string]interface{})
		policyAction := getPolicyAction(policyActionMap)
		updateOpts.Action = policyAction
	}
	log.Printf("[DEBUG] Update AS policy Options: %#v", updateOpts)
	asPolicyID, err := policies.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating ASPolicy %q: %s", asPolicyID, err)
	}

	return resourceASPolicyRead(ctx, d, meta)
}

func resourceASPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Begin to delete AS policy %q", d.Id())
	if delErr := policies.Delete(client, d.Id()); delErr != nil {
		return fmterr.Errorf("error deleting AS policy: %s", delErr)
	}

	return nil
}

var RecurrenceTypes = [3]string{"Daily", "Weekly", "Monthly"}

func resourceASPolicyValidateRecurrenceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range RecurrenceTypes {
		if value == RecurrenceTypes[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, RecurrenceTypes))
	return
}

var PolicyActions = [3]string{"ADD", "REMOVE", "SET"}

func resourceASPolicyValidateActionOperation(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range PolicyActions {
		if value == PolicyActions[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, PolicyActions))
	return
}

var PolicyTypes = [3]string{"ALARM", "SCHEDULED", "RECURRENCE"}

func resourceASPolicyValidatePolicyType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range PolicyTypes {
		if value == PolicyTypes[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, PolicyTypes))
	return
}

func resourceASPolicyValidateName(v interface{}, _ string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must contain more than 1 and less than 64 characters", value))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z-_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters, hyphens, and underscores allowed in %q", value))
	}
	return
}
