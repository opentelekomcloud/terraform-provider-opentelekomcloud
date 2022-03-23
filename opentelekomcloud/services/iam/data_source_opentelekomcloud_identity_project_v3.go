package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
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
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityProjectV3Read performs the project lookup.
func dataSourceIdentityProjectV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	var domainID string
	rawDomainID, ok := d.GetOk("domain_id")
	if ok {
		domainID = rawDomainID.(string)
	} else {
		domainID = client.DomainID
	}

	listOpts := projects.ListOpts{
		DomainID: domainID,
		Name:     d.Get("name").(string),
		ParentID: d.Get("parent_id").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var project projects.Project
	allPages, err := projects.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to query projects: %s", err)
	}

	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve projects: %s", err)
	}

	var filteredProjects []projects.Project
	var enabled bool
	rawEnabled, eOk := d.GetOk("enabled")
	if eOk {
		enabled = rawEnabled.(bool)
	}
	var isDomain bool
	rawIsDomain, dOk := d.GetOk("is_domain")
	if dOk {
		isDomain = rawIsDomain.(bool)
	}
	for _, v := range allProjects {
		if eOk && v.Enabled != enabled {
			continue
		}
		if dOk && v.IsDomain != isDomain {
			continue
		}
		filteredProjects = append(filteredProjects, v)
	}

	if len(filteredProjects) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(filteredProjects) > 1 {
		projectID := config.HwClient.ProjectID
		for _, p := range filteredProjects {
			if p.ID == projectID {
				filteredProjects = []projects.Project{p}
				break
			}
		}
	}

	project = filteredProjects[0]

	log.Printf("[DEBUG] Single project found: %s", project.ID)

	d.SetId(project.ID)
	mErr := multierror.Append(
		d.Set("is_domain", project.IsDomain),
		d.Set("description", project.Description),
		d.Set("domain_id", project.DomainID),
		d.Set("enabled", project.Enabled),
		d.Set("name", project.Name),
		d.Set("parent_id", project.ParentID),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
