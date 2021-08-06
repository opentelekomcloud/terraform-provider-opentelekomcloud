package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/roles"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityRoleV3Read,

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

// dataSourceIdentityRoleV3Read performs the role lookup.
func dataSourceIdentityRoleV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	listOpts := roles.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var role roles.Role
	allPages, err := roles.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to query roles: %s", err)
	}

	allRoles, err := roles.ExtractRoles(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve roles: %s", err)
	}

	if len(allRoles) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allRoles) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allRoles)
		return fmterr.Errorf("your query returned more than one result. Please try a more " +
			"specific search criteria, or set `most_recent` attribute to true.")
	}
	role = allRoles[0]

	log.Printf("[DEBUG] Single Role found: %s", role.ID)
	return diag.FromErr(dataSourceIdentityRoleV3Attributes(d, config, &role))
}

// dataSourceIdentityRoleV3Attributes populates the fields of an Role resource.
func dataSourceIdentityRoleV3Attributes(d *schema.ResourceData, config *cfg.Config, role *roles.Role) error {
	log.Printf("[DEBUG] opentelekomcloud_identity_role_v3 details: %#v", role)

	d.SetId(role.ID)
	mErr := multierror.Append(
		d.Set("name", role.Name),
		d.Set("domain_id", role.DomainID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}

	return nil
}
