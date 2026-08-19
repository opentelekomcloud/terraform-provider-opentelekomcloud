package lts

import (
	"context"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cloud_structuring "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/cloud-structuring"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceLtsStructuringCustomTemplatesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLtsStructuringCustomTemplatesV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"demo_log": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"demo_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     ltsStructuringCustomTemplateDemoFieldSchema(),
						},
						"tag_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     ltsStructuringCustomTemplateTagFieldSchema(),
						},
						"rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"param": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"demo_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ltsStructuringCustomTemplateDemoFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_analysis": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"relation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_defined_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ltsStructuringCustomTemplateTagFieldSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_analysis": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"index": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceLtsStructuringCustomTemplatesV3Read(
	_ context.Context,
	d *schema.ResourceData,
	meta interface{},
) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	templates, err := cloud_structuring.List(client, cloud_structuring.ListOpts{
		ID: d.Get("id").(string),
	})
	if err != nil {
		return diag.Errorf("error querying OpenTelekomCloud LTS v3 custom structuring templates: %s", err)
	}

	flattened, ids := flattenLtsStructuringCustomTemplates(templates)
	stateParts := append([]string{config.GetRegion(d), d.Get("id").(string)}, ids...)
	d.SetId(hashcode.Strings(stateParts))

	mErr := multierror.Append(nil,
		d.Set("templates", flattened),
		d.Set("region", config.GetRegion(d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLtsStructuringCustomTemplates(
	templates []cloud_structuring.StructTemplateModel,
) ([]map[string]interface{}, []string) {
	ordered := append([]cloud_structuring.StructTemplateModel(nil), templates...)
	sort.SliceStable(ordered, func(i, j int) bool {
		return ordered[i].ID < ordered[j].ID
	})

	result := make([]map[string]interface{}, len(ordered))
	ids := make([]string, len(ordered))
	for i, template := range ordered {
		result[i] = map[string]interface{}{
			"id":            template.ID,
			"project_id":    template.ProjectId,
			"template_name": template.Name,
			"template_type": template.Type,
			"demo_log":      template.DemoLog,
			"demo_fields":   flattenLtsStructuringCustomTemplateDemoFields(template.DemoFields),
			"tag_fields":    flattenLtsStructuringCustomTemplateTagFields(template.TagFields),
			"rule":          flattenLtsStructuringCustomTemplateRule(template.Rule),
			"demo_label":    template.DemoLabel,
			"created_at":    common.FormatTimeStampRFC3339(template.CreatedAt/1000, false),
		}
		ids[i] = template.ID
	}
	return result, ids
}

func flattenLtsStructuringCustomTemplateDemoFields(
	fields []cloud_structuring.DemoField,
) []map[string]interface{} {
	result := make([]map[string]interface{}, len(fields))
	for i, field := range fields {
		result[i] = map[string]interface{}{
			"field_name":        field.Name,
			"content":           field.Content,
			"type":              field.Type,
			"is_analysis":       field.IsAnalysis,
			"index":             field.Index,
			"relation":          field.Relation,
			"user_defined_name": field.UserDefinedName,
		}
	}
	return result
}

func flattenLtsStructuringCustomTemplateTagFields(
	fields []cloud_structuring.TagFieldNew,
) []map[string]interface{} {
	result := make([]map[string]interface{}, len(fields))
	for i, field := range fields {
		result[i] = map[string]interface{}{
			"field_name":  field.Name,
			"content":     field.Content,
			"type":        field.Type,
			"is_analysis": field.IsAnalysis,
			"index":       field.Index,
		}
	}
	return result
}

func flattenLtsStructuringCustomTemplateRule(rule *cloud_structuring.TemplateRule) []map[string]interface{} {
	if rule == nil {
		return nil
	}
	return []map[string]interface{}{{
		"type":  rule.Type,
		"param": rule.Param,
	}}
}
