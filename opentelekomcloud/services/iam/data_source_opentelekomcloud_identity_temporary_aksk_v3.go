package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityTempAKSKV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityTempAKSKV3Read,

		Schema: map[string]*schema.Schema{
			"duration_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"security_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityTempAKSKV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client()
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	opts := credentials.CreateTemporaryOpts{
		Methods: []string{"token"},
		Token:   client.Token(),
	}

	if v, ok := d.GetOk("duration_seconds"); ok {
		opts.Duration = v.(int)
	}

	if v, ok := d.GetOk("agency_name"); ok {
		opts.AgencyName = v.(string)
	}

	credential, err := credentials.CreateTemporary(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating temporary AK/SK: %s", err)
	}

	d.SetId(credential.AccessKey)

	mErr := multierror.Append(nil,
		d.Set("access", credential.AccessKey),
		d.Set("secret", credential.SecretKey),
		d.Set("security_token", credential.SecurityToken),
		d.Set("expires_at", credential.ExpiresAt),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting temporary credentials attributes: %s", err)
	}

	return nil
}
