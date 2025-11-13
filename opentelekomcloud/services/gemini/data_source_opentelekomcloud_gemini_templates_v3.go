package gemini

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceGeminiDBV3Templates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBTemplatesRead,

		Schema: map[string]*schema.Schema{
			"templates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_version_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_name": {
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
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_defined": {
							Type:     schema.TypeBool,
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

func dataSourceGeminiDBTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	opts := template.ListOpts{
		Offset: 0,
		Limit:  100,
	}

	allTemplates := make([]template.Configuration, 0)
	for {
		templates, err := template.List(client, opts)
		if err != nil {
			return diag.Errorf("error listing GeminiDB templates: %s", err)
		}

		if len(templates) == 0 {
			break
		}

		allTemplates = append(allTemplates, templates...)

		if len(templates) < opts.Limit {
			break
		}
		opts.Offset += opts.Limit
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("templates", flattenTemplates(allTemplates)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting GeminiDB templates attributes: %s", err)
	}

	return nil
}

func flattenTemplates(templates []template.Configuration) []map[string]interface{} {
	result := make([]map[string]interface{}, len(templates))
	for i, tmpl := range templates {
		result[i] = map[string]interface{}{
			"id":                     tmpl.Id,
			"name":                   tmpl.Name,
			"description":            tmpl.Description,
			"datastore_version_name": tmpl.DataStoreVersionName,
			"datastore_name":         tmpl.DataStoreName,
			"created_at":             tmpl.Created,
			"updated_at":             tmpl.Updated,
			"mode":                   tmpl.Mode,
			"user_defined":           tmpl.UserDefined,
		}
	}
	return result
}
