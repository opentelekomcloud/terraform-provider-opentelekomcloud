package iam

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/projects"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityProjectV3Read,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

// dataSourceIdentityProjectV3Read performs the project lookup.
func dataSourceIdentityProjectV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	listOpts := projects.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Name:     d.Get("name").(string),
		ParentID: d.Get("parent_id").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var project projects.Project
	allPages, err := projects.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("Unable to query projects: %s", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return fmterr.Errorf("Unable to retrieve projects: %s", err)
	}

	if len(allProjects) < 1 {
		return fmterr.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allProjects) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allProjects)
		return fmterr.Errorf("Your query returned more than one result")
	}
	project = allProjects[0]

	log.Printf("[DEBUG] Single project found: %s", project.ID)
	return diag.FromErr(dataSourceIdentityProjectV3Attributes(d, &project))
}

// dataSourceIdentityProjectV3Attributes populates the fields of an Project resource.
func dataSourceIdentityProjectV3Attributes(d *schema.ResourceData, project *projects.Project) error {
	log.Printf("[DEBUG] opentelekomcloud_identity_project_v3 details: %#v", project)

	d.SetId(project.ID)
	d.Set("is_domain", project.IsDomain)
	d.Set("description", project.Description)
	d.Set("domain_id", project.DomainID)
	d.Set("enabled", project.Enabled)
	d.Set("name", project.Name)
	d.Set("parent_id", project.ParentID)

	return nil
}
