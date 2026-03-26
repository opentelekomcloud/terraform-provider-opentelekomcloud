package swr

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/domains"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSwrDomainV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrDomainRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deadline": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceSwrDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := domains.GetOpts{
		Namespace:    d.Get("organization").(string),
		Repository:   repository(d.Get("repository").(string)),
		AccessDomain: d.Get("access_domain").(string),
	}
	domain, err := domains.Get(client, opts)
	if err != nil {
		return fmterr.Errorf("error reading domain: %w", err)
	}
	d.SetId(opts.AccessDomain)

	mErr := multierror.Append(
		d.Set("access_domain", strings.ToUpper(domain.AccessDomain)),
		d.Set("repository", domain.Repository),
		d.Set("organization", domain.Namespace),
		d.Set("description", domain.Description),
		d.Set("status", domain.Status),
		d.Set("permission", domain.Permit),
		d.Set("deadline", domain.Deadline),
		d.Set("created", domain.Created),
		d.Set("updated", domain.Updated),
		d.Set("creator_id", domain.CreatorID),
		d.Set("creator_name", domain.CreatorName),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting resource fields: %w", err)
	}

	return nil
}
