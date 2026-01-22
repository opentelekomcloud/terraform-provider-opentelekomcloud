package ces

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ces/v2/templates"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func ResourceAlarmTemplateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmTemplateV2Create,
		ReadContext:   resourceAlarmTemplateV2Read,
		UpdateContext: resourceAlarmTemplateV2Update,
		DeleteContext: resourceAlarmTemplateV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z][\w\-().]*$`),
						"The name must start with a letter and can contain letters, digits, underscores (_), "+
							"hyphens (-), parentheses, and periods (.).",
					),
				),
			},
			"type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 2,
				}),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"policies": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 50,
				Set:      resourceTemplatePolicyHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(3, 32),
								validation.StringMatch(
									regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*\.[a-zA-Z][a-zA-Z0-9_]*$`),
									"The namespace must be in the format service.item.",
								),
							),
						},
						"dimension_name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 32),
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 96),
								validation.StringMatch(
									regexp.MustCompile(`^[a-zA-Z][\w]*$`),
									"The metric name must start with a letter.",
								),
							),
						},
						"period": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 1, 300, 1200, 3600, 14400, 86400,
							}),
						},
						"filter": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"average", "max", "min", "sum", "variance",
							}, false),
						},
						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								">", "<", ">=", "<=", "=", "!=",
							}, false),
						},
						"value": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"unit": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 32),
						},
						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 180),
						},
						"alarm_level": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      2,
							ValidateFunc: validation.IntBetween(1, 4),
						},
						"suppress_duration": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 300, 600, 900, 1800, 3600, 10800, 21600, 43200, 86400,
							}),
						},
					},
				},
			},
			"delete_associate_alarm": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTemplatePolicyHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["namespace"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["namespace"].(string)))
	}
	if m["dimension_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["dimension_name"].(string)))
	}
	if m["metric_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_name"].(string)))
	}
	if m["period"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["period"].(int)))
	}
	if m["filter"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["filter"].(string)))
	}
	if m["comparison_operator"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["comparison_operator"].(string)))
	}
	if m["value"] != nil {
		buf.WriteString(fmt.Sprintf("%f-", m["value"].(float64)))
	}
	if m["unit"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["unit"].(string)))
	}
	if m["count"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["count"].(int)))
	}
	if m["alarm_level"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["alarm_level"].(int)))
	}
	if m["suppress_duration"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["suppress_duration"].(int)))
	}

	return hashcode.String(buf.String())
}

func buildPoliciesOpts(d *schema.ResourceData) []templates.Policy {
	rawPolicies := d.Get("policies").(*schema.Set).List()
	if len(rawPolicies) == 0 {
		return nil
	}

	policiesOpts := make([]templates.Policy, len(rawPolicies))
	for i, v := range rawPolicies {
		policy := v.(map[string]interface{})
		policiesOpts[i] = templates.Policy{
			Namespace:          policy["namespace"].(string),
			DimensionName:      policy["dimension_name"].(string),
			MetricName:         policy["metric_name"].(string),
			Period:             policy["period"].(int),
			Filter:             policy["filter"].(string),
			ComparisonOperator: policy["comparison_operator"].(string),
			Value:              policy["value"].(float64),
			Unit:               policy["unit"].(string),
			Count:              policy["count"].(int),
			AlarmLevel:         policy["alarm_level"].(int),
			SuppressDuration:   policy["suppress_duration"].(int),
		}
	}
	return policiesOpts
}

func resourceAlarmTemplateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	createOpts := templates.CreateOpts{
		TemplateName:        d.Get("name").(string),
		TemplateType:        d.Get("type").(int),
		TemplateDescription: d.Get("description").(string),
		Policies:            buildPoliciesOpts(d),
	}

	log.Printf("[DEBUG] CES Alarm Template V2 Create Options: %#v", createOpts)

	templateId, err := templates.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating CES Alarm Template V2: %w", err)
	}

	log.Printf("[DEBUG] CES Alarm Template V2 created with ID: %s", templateId)
	d.SetId(templateId)

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceAlarmTemplateV2Read(clientCtx, d, meta)
}

func resourceAlarmTemplateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	template, err := templates.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CES Alarm Template V2")
	}

	if template == nil {
		log.Printf("[WARN] CES Alarm Template V2 (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	// Convert template_type string to int
	templateType := 0
	if template.TemplateType == "custom_event" {
		templateType = 2
	}

	mErr := multierror.Append(nil,
		d.Set("name", template.TemplateName),
		d.Set("type", templateType),
		d.Set("description", template.TemplateDescription),
		d.Set("template_id", template.TemplateId),
		d.Set("policies", flattenPolicies(template.Policies)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenPolicies(policiesResp []templates.PolicyResp) []map[string]interface{} {
	if len(policiesResp) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(policiesResp))
	for i, policy := range policiesResp {
		result[i] = map[string]interface{}{
			"namespace":           policy.Namespace,
			"dimension_name":      policy.DimensionName,
			"metric_name":         policy.MetricName,
			"period":              policy.Period,
			"filter":              policy.Filter,
			"comparison_operator": policy.ComparisonOperator,
			"value":               policy.Value,
			"unit":                policy.Unit,
			"count":               policy.Count,
			"alarm_level":         policy.AlarmLevel,
			"suppress_duration":   policy.SuppressDuration,
		}
	}
	return result
}

func resourceAlarmTemplateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	if d.HasChanges("name", "description", "policies") {
		updateOpts := templates.UpdateOpts{
			TemplateName:        d.Get("name").(string),
			TemplateType:        d.Get("type").(int),
			TemplateDescription: d.Get("description").(string),
			Policies:            buildPoliciesOpts(d),
		}

		log.Printf("[DEBUG] CES Alarm Template V2 Update Options: %#v", updateOpts)

		err := templates.Update(client, d.Id(), updateOpts)
		if err != nil {
			return fmterr.Errorf("error updating CES Alarm Template V2: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, cesClientV2)
	return resourceAlarmTemplateV2Read(clientCtx, d, meta)
}

func resourceAlarmTemplateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cesClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CesV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	templateId := d.Id()
	log.Printf("[DEBUG] Deleting CES Alarm Template V2: %s", templateId)

	deleteOpts := templates.DeleteOpts{
		TemplateIds:          []string{templateId},
		DeleteAssociateAlarm: d.Get("delete_associate_alarm").(bool),
	}

	_, err = templates.Delete(client, deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CES Alarm Template V2")
	}

	return nil
}
