package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityProjectsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProjectsV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_domain": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityProjectsV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	allPages, err := projects.List(client, projects.ListOpts{}).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to query projects: %s", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve projects: %s", err)
	}

	var flattenProjects []map[string]interface{}

	for _, item := range allProjects {
		project := map[string]interface{}{
			"name":        item.Name,
			"description": item.Description,
			"domain_id":   item.DomainID,
			"parent_id":   item.ParentID,
			"enabled":     item.Enabled,
			"is_domain":   item.IsDomain,
			"project_id":  item.ID,
		}
		flattenProjects = append(flattenProjects, project)
	}

	d.SetId(client.DomainID)
	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("projects", flattenProjects),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
