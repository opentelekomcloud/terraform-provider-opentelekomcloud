package swr

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/domains"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSwrDomainsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSwrDomainsRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_domain": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceSwrDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := domains.ListOpts{
		Namespace:  d.Get("organization").(string),
		Repository: repository(d.Get("repository").(string)),
	}
	accessDomains, err := domains.List(client, opts)
	if err != nil {
		return fmterr.Errorf("error reading domain: %w", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("access_domains", setAccessDomains(accessDomains)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting resource fields: %w", err)
	}

	return nil
}

func setAccessDomains(domainsInResp []domains.AccessDomain) []map[string]interface{} {
	var accessDomains []map[string]interface{}
	for _, accessDomainInResp := range domainsInResp {
		accessDomain := map[string]interface{}{
			"access_domain": strings.ToUpper(accessDomainInResp.AccessDomain),
			"description":   accessDomainInResp.Description,
			"permission":    accessDomainInResp.Permit,
			"deadline":      accessDomainInResp.Deadline,
			"created":       accessDomainInResp.Created,
			"updated":       accessDomainInResp.Updated,
			"creator_id":    accessDomainInResp.CreatorID,
			"creator_name":  accessDomainInResp.CreatorName,
			"status":        accessDomainInResp.Status,
		}
		accessDomains = append(accessDomains, accessDomain)
	}
	return accessDomains
}
