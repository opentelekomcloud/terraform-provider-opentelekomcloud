package as

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v2/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceASPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASPolicyV2Create,
		ReadContext:   resourceASPolicyV2Read,
		UpdateContext: resourceASPolicyV2Update,
		DeleteContext: resourceASPolicyV2Delete,

		CustomizeDiff: validateAction,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[a-z][a-z0-9._-]+[a-z0-9]+$`),
						"Only lowercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
					),
					validation.StringDoesNotMatch(
						regexp.MustCompile(`_{3,}?|\.{2,}?|-{2,}?`),
						"Periods, underscores, and hyphens cannot be placed next to each other. A maximum of two consecutive underscores are allowed.",
					),
				),
			},
			"scaling_policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALARM", "SCHEDULED", "RECURRENCE",
				}, false),
			},
			"scaling_resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scaling_resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SCALING_GROUP", "BANDWIDTH",
				}, false),
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduled_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"launch_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"recurrence_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Daily", "Weekly", "Monthly",
							}, false),
						},
						"recurrence_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"scaling_policy_action": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ADD", "REMOVE", "REDUCE", "SET",
							}, false),
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 300),
						},
						"percentage": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 20000),
						},
						"limits": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"cool_down_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bandwidth_share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceASPolicyScheduledPolicy(d *schema.ResourceData) policies.SchedulePolicyOpts {
	rawScheduledPolicyList := d.Get("scheduled_policy").(*schema.Set).List()
	var scheduledOpts policies.SchedulePolicyOpts
	if len(rawScheduledPolicyList) > 0 {
		rawScheduledPolicy := rawScheduledPolicyList[0].(map[string]interface{})
		scheduledOpts = policies.SchedulePolicyOpts{
			LaunchTime:      rawScheduledPolicy["launch_time"].(string),
			RecurrenceType:  rawScheduledPolicy["recurrence_type"].(string),
			RecurrenceValue: rawScheduledPolicy["recurrence_value"].(string),
			StartTime:       rawScheduledPolicy["start_time"].(string),
			EndTime:         rawScheduledPolicy["end_time"].(string),
		}
	}
	return scheduledOpts
}

func resourceASPolicyScalingAction(d *schema.ResourceData) policies.ActionOpts {
	rawPolicyActionList := d.Get("scaling_policy_action").(*schema.Set).List()
	var actionOpts policies.ActionOpts
	if len(rawPolicyActionList) > 0 {
		rawPolicyAction := rawPolicyActionList[0].(map[string]interface{})
		actionOpts = policies.ActionOpts{
			Operation:  rawPolicyAction["operation"].(string),
			Size:       rawPolicyAction["size"].(int),
			Percentage: rawPolicyAction["percentage"].(int),
			Limits:     rawPolicyAction["limits"].(int),
		}
	}
	return actionOpts
}

func resourceASPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.AutoscalingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := policies.PolicyOpts{
		PolicyName:     d.Get("scaling_policy_name").(string),
		PolicyType:     d.Get("scaling_policy_type").(string),
		ResourceID:     d.Get("scaling_resource_id").(string),
		ResourceType:   d.Get("scaling_resource_type").(string),
		AlarmID:        d.Get("alarm_id").(string),
		SchedulePolicy: resourceASPolicyScheduledPolicy(d),
		PolicyAction:   resourceASPolicyScalingAction(d),
		CoolDownTime:   d.Get("cool_down_time").(int),
	}
	asPolicyID, err := policies.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating AS Policy: %w", err)
	}
	d.SetId(asPolicyID)

	return resourceASPolicyV2Read(ctx, d, meta)
}

func resourceASPolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.AutoscalingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	asPolicy, err := policies.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS Policy")
	}

	mErr := multierror.Append(
		d.Set("scaling_policy_name", asPolicy.PolicyName),
		d.Set("scaling_policy_type", asPolicy.Type),
		d.Set("scaling_resource_id", asPolicy.ResourceID),
		d.Set("scaling_resource_type", asPolicy.ScalingResourceType),
		d.Set("alarm_id", asPolicy.AlarmID),
		d.Set("cool_down_time", asPolicy.CoolDownTime),
	)

	scheduledPolicy := []map[string]interface{}{
		{
			"launch_time":      asPolicy.SchedulePolicy.LaunchTime,
			"recurrence_type":  asPolicy.SchedulePolicy.RecurrenceType,
			"recurrence_value": asPolicy.SchedulePolicy.RecurrenceValue,
			"start_time":       asPolicy.SchedulePolicy.StartTime,
			"end_time":         asPolicy.SchedulePolicy.EndTime,
		},
	}

	policyAction := []map[string]interface{}{
		{
			"operation":  asPolicy.PolicyAction.Operation,
			"size":       asPolicy.PolicyAction.Size,
			"percentage": asPolicy.PolicyAction.Percentage,
			"limits":     asPolicy.PolicyAction.Limits,
		},
	}

	metadata := []map[string]interface{}{
		{
			"bandwidth_share_type": asPolicy.Metadata.BandwidthShareType,
			"eip_id":               asPolicy.Metadata.EipID,
			"eip_address":          asPolicy.Metadata.EipAddress,
		},
	}
	mErr = multierror.Append(mErr,
		d.Set("scheduled_policy", scheduledPolicy),
		d.Set("scaling_policy_action", policyAction),
		d.Set("metadata", metadata),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting AS PolicyV2 fields: %w", err)
	}

	return nil
}

func resourceASPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.AutoscalingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	updateOpts := policies.PolicyOpts{
		PolicyName:     d.Get("scaling_policy_name").(string),
		PolicyType:     d.Get("scaling_policy_type").(string),
		ResourceID:     d.Get("scaling_resource_id").(string),
		ResourceType:   d.Get("scaling_resource_type").(string),
		AlarmID:        d.Get("alarm_id").(string),
		SchedulePolicy: resourceASPolicyScheduledPolicy(d),
		PolicyAction:   resourceASPolicyScalingAction(d),
		CoolDownTime:   d.Get("cool_down_time").(int),
	}

	asPolicyID, err := policies.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating ASPolicy %q: %w", asPolicyID, err)
	}

	return resourceASPolicyV2Read(ctx, d, meta)
}

func resourceASPolicyV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	if err := policies.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("error deleting AS Policy: %w", err)
	}

	return nil
}

func validateAction(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	rawPolicyActionList := d.Get("scaling_policy_action").(*schema.Set).List()
	if len(rawPolicyActionList) > 0 {
		rawPolicyAction := rawPolicyActionList[0].(map[string]interface{})
		if rawPolicyAction["percentage"].(int) > 0 && rawPolicyAction["size"].(int) > 0 {
			return fmt.Errorf("select one from `percentage` or `size`")
		}
	}
	return nil
}
