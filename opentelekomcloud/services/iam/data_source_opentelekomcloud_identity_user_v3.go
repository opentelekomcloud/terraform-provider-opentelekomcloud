package iam

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3.0/security"
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
			"password_expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mfa_device": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityUserV3Read performs the user lookup.
func dataSourceIdentityUserV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	listOpts := users.ListOpts{
		DomainID: d.Get("domain_id").(string),
		Enabled:  &enabled,
		Name:     d.Get("name").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	allPages, err := users.List(client, listOpts).AllPages()
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

	user := allUsers[0]

	log.Printf("[DEBUG] Single user found: %s", user.ID)

	d.SetId(user.ID)
	mErr := multierror.Append(
		d.Set("domain_id", user.DomainID),
		d.Set("enabled", user.Enabled),
		d.Set("name", user.Name),
		d.Set("password_expires_at", user.PasswordExpiresAt.Format(time.RFC3339)),
	)

	client.Endpoint = strings.Replace(client.Endpoint, "v3", "v3.0", 1)

	mfa, err := security.GetUserMfaDevice(client, user.ID)
	switch err.(type) {
	case golangsdk.ErrDefault403:
		log.Printf("[DEBUG] Security administrator permissions needed to set MFA")
	case nil, golangsdk.ErrDefault404:
		if mfa == nil {
			mErr = multierror.Append(mErr, d.Set("mfa_device", ""))
		} else {
			mErr = multierror.Append(mErr, d.Set("mfa_device", mfa.SerialNumber))
		}
	default:
		log.Printf("[DEBUG] Error getting MFA info: %v", err.Error())
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting user fields: %s", err)
	}
	return nil
}
