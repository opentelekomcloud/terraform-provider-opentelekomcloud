package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/agency"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityAgencyV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityAgencyV3Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trust_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"trust_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIdentityAgencyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := agencyClient(d, config)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud identity client: %s", err)
	}

	opts := agency.ListOpts{
		Name:          d.Get("name").(string),
		DomainID:      client.DomainID,
		TrustDomainID: d.Get("trust_domain_id").(string),
	}
	pages, err := agency.List(client, opts).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing agencies: %w", err)
	}
	agencies, err := agency.ExtractAgencies(pages)
	if err != nil {
		return fmterr.Errorf("error extracting agencies: %w", err)
	}
	if len(agencies) < 1 {
		return common.DataSourceTooFewDiag
	}
	if len(agencies) > 1 {
		return common.DataSourceTooManyDiag
	}

	result := agencies[0]

	d.SetId(result.ID)
	mErr := multierror.Append(
		d.Set("name", result.Name),
		d.Set("trust_domain_id", result.DelegatedDomainID),
		d.Set("trust_domain_name", result.DelegatedDomainName),
		d.Set("description", result.Description),
		d.Set("duration", result.Duration),
		d.Set("expire_time", result.ExpireTime),
		d.Set("create_time", result.CreateTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting agency fields: %w", err)
	}

	return nil
}
