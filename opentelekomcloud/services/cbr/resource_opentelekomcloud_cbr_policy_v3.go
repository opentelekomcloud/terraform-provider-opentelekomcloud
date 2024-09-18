package cbr

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cbr/v3/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCBRPolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCBRPolicyV3Create,
		ReadContext:   resourceCBRPolicyV3Read,
		UpdateContext: resourceCBRPolicyV3Update,
		DeleteContext: resourceCBRPolicyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateName,
			},
			"operation_definition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day_backups": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"max_backups": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(-1, 99999),
						},
						"month_backups": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"retention_duration_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(-1, 99999),
						},
						"timezone": {
							Type:     schema.TypeString,
							Required: true,
						},
						"week_backups": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
						"year_backups": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(0, 100),
						},
					},
				},
			},
			"operation_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"backup", "replication",
				}, false),
			},
			"trigger_pattern": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"destination_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"destination_region"},
			},
		},
	}
}

func resourceCBRPolicyV3OpDefinition(d *schema.ResourceData) *policies.PolicyODCreate {
	opDefinitionRaw := d.Get("operation_definition").([]interface{})
	if len(opDefinitionRaw) == 1 {
		policyODCreate := policies.PolicyODCreate{}
		opDefinition := opDefinitionRaw[0].(map[string]interface{})
		if destinationProjectID, ok := d.GetOk("destination_project_id"); ok {
			policyODCreate.DestinationProjectId = destinationProjectID.(string)
			policyODCreate.DestinationRegion = d.Get("destination_region").(string)
		}
		policyODCreate.DailyBackups = opDefinition["day_backups"].(int)
		policyODCreate.WeekBackups = opDefinition["week_backups"].(int)
		policyODCreate.YearBackups = opDefinition["year_backups"].(int)
		policyODCreate.MonthBackups = opDefinition["month_backups"].(int)
		policyODCreate.MaxBackups = opDefinition["max_backups"].(int)
		policyODCreate.RetentionDurationDays = opDefinition["retention_duration_days"].(int)
		policyODCreate.Timezone = opDefinition["timezone"].(string)
		return &policyODCreate
	}
	return &policies.PolicyODCreate{
		Timezone: "UTC+00:00",
	}
}

func resourceCBRPolicyV3TriggerPattern(d *schema.ResourceData) []string {
	triggerPatternRaw := d.Get("trigger_pattern").([]interface{})
	patterns := make([]string, 0)
	for _, v := range triggerPatternRaw {
		patterns = append(patterns, v.(string))
	}
	return patterns
}

func resourceCBRPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	enabled := d.Get("enabled").(bool)

	createOpts := policies.CreateOpts{
		Name:                d.Get("name").(string),
		OperationDefinition: resourceCBRPolicyV3OpDefinition(d),
		Enabled:             &enabled,
		OperationType:       policies.OperationType(d.Get("operation_type").(string)),
		Trigger: &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: resourceCBRPolicyV3TriggerPattern(d),
			},
		},
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	cbrPolicy, err := policies.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 policy: %s", err)
	}

	// Store the ID
	d.SetId(cbrPolicy.ID)

	return resourceCBRPolicyV3Read(ctx, d, meta)
}

func resourceCBRPolicyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	cbrPolicy, err := policies.Get(client, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing CBR policy %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}
		return fmterr.Errorf("error retrieving CBRv3 policy: %s", err)
	}

	log.Printf("[DEBUG] Retrieved policy %s: %+v", d.Id(), cbrPolicy)
	mErr := multierror.Append(nil,
		d.Set("enabled", cbrPolicy.Enabled),
		d.Set("name", cbrPolicy.Name),
		d.Set("operation_type", cbrPolicy.OperationType),
		d.Set("trigger_pattern", cbrPolicy.Trigger.Properties.Pattern),
		d.Set("region", config.GetRegion(d)),
	)

	var opDefinitionList []map[string]interface{}
	opDefinition := make(map[string]interface{})
	cbrPolicyOD := cbrPolicy.OperationDefinition
	opDefinition["day_backups"] = cbrPolicyOD.DailyBackups
	opDefinition["max_backups"] = cbrPolicyOD.MaxBackups
	opDefinition["month_backups"] = cbrPolicyOD.MonthBackups
	opDefinition["retention_duration_days"] = cbrPolicyOD.RetentionDurationDays
	opDefinition["timezone"] = cbrPolicyOD.Timezone
	opDefinition["week_backups"] = cbrPolicyOD.WeekBackups
	opDefinition["year_backups"] = cbrPolicyOD.YearBackups
	opDefinitionList = append(opDefinitionList, opDefinition)
	if err := d.Set("operation_definition", opDefinitionList); err != nil {
		return fmterr.Errorf("error setting operetion_definition: %s", err)
	}
	mErr = multierror.Append(mErr,
		d.Set("destination_project_id", cbrPolicyOD.DestinationProjectId),
		d.Set("destination_region", cbrPolicyOD.DestinationRegion),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceCBRPolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	var updateOpts policies.UpdateOpts

	if d.HasChange("name") {
		newName := d.Get("name")
		updateOpts.Name = newName.(string)
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("trigger_pattern") {
		pattern := resourceCBRPolicyV3TriggerPattern(d)
		updateOpts.Trigger = &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: pattern,
			},
		}
	}

	if d.HasChange("operation_definition") {
		opDefinition := resourceCBRPolicyV3OpDefinition(d)
		updateOpts.OperationDefinition = opDefinition
	}

	_, err = policies.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud CBRv3 policy: %s", err)
	}

	return resourceCBRPolicyV3Read(ctx, d, meta)
}

func resourceCBRPolicyV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud CBRv3 client: %s", err)
	}

	log.Printf("[DEBUG] Deleting CBRv3 policy %s", d.Id())

	if err = policies.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("eror deleting OpenTelekomCloud CBRv3 policy: %s", err)
	}

	d.SetId("")
	return nil
}
