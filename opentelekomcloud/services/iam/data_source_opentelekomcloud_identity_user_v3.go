package iam

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityUserV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityUserV3Read,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"name": {
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

// dataSourceIdentityUserV3Read performs the user lookup.
func dataSourceIdentityUserV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	listOpts := users.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Enabled:  &enabled,
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var user users.User
	allPages, err := users.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to query users: %s", err)
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve users: %s", err)
	}

	if len(allUsers) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allUsers) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allUsers)
		return fmterr.Errorf("your query returned more than one result")
	}
	user = allUsers[0]

	log.Printf("[DEBUG] Single user found: %s", user.ID)
	return diag.FromErr(dataSourceIdentityUserV3Attributes(d, &user))
}

// dataSourceIdentityUserV3Attributes populates the fields of an User resource.
func dataSourceIdentityUserV3Attributes(d *schema.ResourceData, user *users.User) error {
	log.Printf("[DEBUG] opentelekomcloud_identity_user_v3 details: %#v", user)

	d.SetId(user.ID)
	mErr := multierror.Append(
		d.Set("default_project_id", user.DefaultProjectID),
		d.Set("domain_id", user.DomainID),
		d.Set("enabled", user.Enabled),
		d.Set("name", user.Name),
		d.Set("password_expires_at", user.PasswordExpiresAt.Format(time.RFC3339)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}
	return nil
}
