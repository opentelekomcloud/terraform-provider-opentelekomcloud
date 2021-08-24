package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityGroupV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityGroupV3Read,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
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

// dataSourceIdentityGroupV3Read performs the group lookup.
func dataSourceIdentityGroupV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	listOpts := groups.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var group groups.Group
	allPages, err := groups.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to query groups: %s", err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve roles: %s", err)
	}

	if len(allGroups) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allGroups) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allGroups)
		return fmterr.Errorf("your query returned more than one result. Please try a more " +
			"specific search criteria, or set `most_recent` attribute to true.")
	}
	group = allGroups[0]

	log.Printf("[DEBUG] Single group found: %s", group.ID)
	return diag.FromErr(dataSourceIdentityGroupV3Attributes(d, config, &group))
}

// dataSourceIdentityRoleV3Attributes populates the fields of an Role resource.
func dataSourceIdentityGroupV3Attributes(d *schema.ResourceData, config *cfg.Config, group *groups.Group) error {
	log.Printf("[DEBUG] opentelekomcloud_identity_group_v3 details: %#v", group)

	d.SetId(group.ID)
	mErr := multierror.Append(
		d.Set("name", group.Name),
		d.Set("domain_id", group.DomainID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}

	return nil
}
