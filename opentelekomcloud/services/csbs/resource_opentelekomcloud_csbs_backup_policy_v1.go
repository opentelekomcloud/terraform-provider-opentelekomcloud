package csbs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/csbs/v1/policies"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const errorSaveMsg = "[DEBUG] Error saving %s to state for OpenTelekomCloud CSBS backup policy (%s): %s"

func ResourceCSBSBackupPolicyV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCSBSBackupPolicyCreate,
		ReadContext:   resourceCSBSBackupPolicyRead,
		UpdateContext: resourceCSBSBackupPolicyUpdate,
		DeleteContext: resourceCSBSBackupPolicyDelete,

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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "fc4d5750-22e7-4798-8a46-f48f62c4c1da",
				ForceNew: true,
			},
			"common": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_operation": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
						"max_backups": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"retention_duration_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"day_backups": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"week_backups": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"month_backups": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"year_backups": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"permanent": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"trigger_pattern": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operation_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"resource": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: common.ValidateTags,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceCSBSBackupPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	policyClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CSBSv1 client: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ProviderId:  d.Get("provider_id").(string),
		Parameters: policies.PolicyParam{
			Common: resourceCSBSCommonParamsV1(d),
		},

		ScheduledOperations: resourceCSBSScheduleV1(d),

		Resources: resourceCSBSResourceV1(d),
		Tags:      resourceCSBSTagsV1(d),
	}

	backupPolicy, err := policies.Create(policyClient, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Backup Policy : %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"suspended"},
		Refresh:    waitForCSBSBackupPolicyActive(policyClient, backupPolicy.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for Backup Policy (%s) to become available: %s", backupPolicy.ID, err)
	}

	d.SetId(backupPolicy.ID)
	return resourceCSBSBackupPolicyRead(ctx, d, meta)
}

func resourceCSBSBackupPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	policyClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CSBSv1 client: %s", err)
	}

	backupPolicy, err := policies.Get(policyClient, d.Id())
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[WARN] Removing backup policy %s as it's already gone", d.Id())
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving backup policy: %s", err)
	}

	if err := d.Set("resource", flattenCSBSPolicyResources(*backupPolicy)); err != nil {
		return fmterr.Errorf(errorSaveMsg, "resource", d.Id(), err)
	}

	scheduledOperations := flattenCSBSScheduledOperations(*backupPolicy)
	scheduledOperationsRaw := d.Get("scheduled_operation").(*schema.Set).List()
	if len(scheduledOperationsRaw) == 1 {
		operationDefinitionZero := scheduledOperationsRaw[0].(map[string]interface{})
		scheduledOperations[0]["day_backups"] = operationDefinitionZero["day_backups"]
		scheduledOperations[0]["week_backups"] = operationDefinitionZero["week_backups"]
		scheduledOperations[0]["month_backups"] = operationDefinitionZero["month_backups"]
		scheduledOperations[0]["year_backups"] = operationDefinitionZero["year_backups"]
		scheduledOperations[0]["timezone"] = operationDefinitionZero["timezone"]
	}

	if err := d.Set("scheduled_operation", scheduledOperations); err != nil {
		return fmterr.Errorf(errorSaveMsg, "scheduler_operation", d.Id(), err)
	}

	tagsMap := make(map[string]string)
	for _, tag := range backupPolicy.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	me := multierror.Append(nil,
		d.Set("name", backupPolicy.Name),
		d.Set("common", backupPolicy.Parameters.Common),
		d.Set("status", backupPolicy.Status),
		d.Set("description", backupPolicy.Description),
		d.Set("provider_id", backupPolicy.ProviderId),
		d.Set("created_at", backupPolicy.CreatedAt.Format(time.RFC3339)),
		d.Set("region", config.GetRegion(d)),
		d.Set("tags", tagsMap),
	)

	return diag.FromErr(me.ErrorOrNil())
}

func resourceCSBSBackupPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	policyClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CSBSv1 client: %s", err)
	}
	var updateOpts policies.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("resource") {
		updateOpts.Resources = resourceCSBSResourceV1(d)
	}
	if d.HasChange("scheduled_operation") {
		updateOpts.ScheduledOperations = resourceCSBSScheduleUpdateV1(d)
	}

	_, err = policies.Update(policyClient, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating Backup Policy: %s", err)
	}

	return resourceCSBSBackupPolicyRead(ctx, d, meta)
}

func resourceCSBSBackupPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	policyClient, err := config.CsbsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating CSBS client: %s", err)
	}

	err = policies.Delete(policyClient, d.Id())
	if err != nil {
		return fmterr.Errorf("error delete CSBSv1 policy; %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available"},
		Target:     []string{"deleted"},
		Refresh:    waitForCSBSPolicyDelete(policyClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for delete backup Policy: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCSBSBackupPolicyActive(policyClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := policies.Get(policyClient, policyID)
		if err != nil {
			return nil, "", err
		}

		if policy.Status == "error" {
			return policy, policy.Status, nil
		}
		return policy, policy.Status, nil
	}
}

func waitForCSBSPolicyDelete(policyClient *golangsdk.ServiceClient, policyID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := policies.Get(policyClient, policyID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Backup Policy %s", policyID)
				return policy, "deleted", nil
			}
			return policy, "available", err
		}

		return policy, policy.Status, nil
	}
}

func resourceCSBSScheduleV1(d *schema.ResourceData) []policies.ScheduledOperation {
	scheduledOperations := d.Get("scheduled_operation").(*schema.Set).List()
	scheduledOperation := make([]policies.ScheduledOperation, len(scheduledOperations))
	for i, raw := range scheduledOperations {
		rawMap := raw.(map[string]interface{})
		scheduledOperation[i] = policies.ScheduledOperation{
			Name:          rawMap["name"].(string),
			Description:   rawMap["description"].(string),
			Enabled:       rawMap["enabled"].(bool),
			OperationType: rawMap["operation_type"].(string),
			Trigger: policies.Trigger{
				Properties: policies.TriggerProperties{
					Pattern: rawMap["trigger_pattern"].(string),
				},
			},

			OperationDefinition: policies.OperationDefinition{
				MaxBackups:            pointerto.Int(rawMap["max_backups"].(int)),
				RetentionDurationDays: rawMap["retention_duration_days"].(int),
				Permanent:             rawMap["permanent"].(bool),
				DayBackups:            rawMap["day_backups"].(int),
				WeekBackups:           rawMap["week_backups"].(int),
				MonthBackups:          rawMap["month_backups"].(int),
				YearBackups:           rawMap["year_backups"].(int),
				TimeZone:              rawMap["timezone"].(string),
			},
		}
	}

	return scheduledOperation
}

func resourceCSBSResourceV1(d *schema.ResourceData) []policies.Resource {
	resources := d.Get("resource").(*schema.Set).List()
	res := make([]policies.Resource, len(resources))
	for i, raw := range resources {
		rawMap := raw.(map[string]interface{})
		res[i] = policies.Resource{
			Name: rawMap["name"].(string),
			Id:   rawMap["id"].(string),
			Type: rawMap["type"].(string),
		}
	}
	return res
}

func resourceCSBSScheduleUpdateV1(d *schema.ResourceData) []policies.ScheduledOperationToUpdate {
	oldSORaw, newSORaw := d.GetChange("scheduled_operation")
	oldSOList := oldSORaw.(*schema.Set).List()
	newSOSetList := newSORaw.(*schema.Set).List()

	// scheduledOperations := d.Get("scheduled_operation").(*schema.Set).List()
	schedule := make([]policies.ScheduledOperationToUpdate, len(newSOSetList))
	for i, raw := range newSOSetList {
		rawNewMap := raw.(map[string]interface{})
		rawOldMap := oldSOList[i].(map[string]interface{})
		schedule[i] = policies.ScheduledOperationToUpdate{
			Id:          rawOldMap["id"].(string),
			Name:        rawNewMap["name"].(string),
			Description: rawNewMap["description"].(string),
			Enabled:     rawNewMap["enabled"].(bool),
			Trigger: policies.Trigger{
				Properties: policies.TriggerProperties{
					Pattern: rawNewMap["trigger_pattern"].(string),
				},
			},
			OperationDefinition: policies.OperationDefinition{
				MaxBackups:            pointerto.Int(rawNewMap["max_backups"].(int)),
				RetentionDurationDays: rawNewMap["retention_duration_days"].(int),
				Permanent:             rawNewMap["permanent"].(bool),
				DayBackups:            rawNewMap["day_backups"].(int),
				WeekBackups:           rawNewMap["week_backups"].(int),
				MonthBackups:          rawNewMap["month_backups"].(int),
				YearBackups:           rawNewMap["year_backups"].(int),
				TimeZone:              rawNewMap["timezone"].(string),
			},
		}
	}

	return schedule
}

func resourceCSBSCommonParamsV1(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("common").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func flattenCSBSScheduledOperations(backupPolicy policies.BackupPolicy) []map[string]interface{} {
	var scheduledOperationList []map[string]interface{}
	for _, schedule := range backupPolicy.ScheduledOperations {
		mapping := map[string]interface{}{
			"enabled":                 schedule.Enabled,
			"trigger_id":              schedule.TriggerID,
			"name":                    schedule.Name,
			"description":             schedule.Description,
			"operation_type":          schedule.OperationType,
			"max_backups":             schedule.OperationDefinition.MaxBackups,
			"retention_duration_days": schedule.OperationDefinition.RetentionDurationDays,
			"permanent":               schedule.OperationDefinition.Permanent,
			"trigger_name":            schedule.Trigger.Name,
			"trigger_type":            schedule.Trigger.Type,
			"trigger_pattern":         schedule.Trigger.Properties.Pattern,
			"id":                      schedule.ID,
		}
		scheduledOperationList = append(scheduledOperationList, mapping)
	}

	return scheduledOperationList
}

func flattenCSBSPolicyResources(backupPolicy policies.BackupPolicy) []map[string]interface{} {
	var resourceList []map[string]interface{}
	for _, resources := range backupPolicy.Resources {
		mapping := map[string]interface{}{
			"id":   resources.Id,
			"type": resources.Type,
			"name": resources.Name,
		}
		resourceList = append(resourceList, mapping)
	}

	return resourceList
}
