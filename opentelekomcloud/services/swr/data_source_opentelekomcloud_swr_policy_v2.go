package swr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSwrPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRead,

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"params": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"tag_selector": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"pattern": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourcePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	namespace := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	id := d.Get("policy_id").(string)
	retPolicy, err := policy.Get(client, namespace, repository, id)
	if err != nil {
		return fmterr.Errorf("error reading retention policy: %w", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("algorithm", retPolicy.Algorithm),
		d.Set("rules", setRules(retPolicy.Rules)),
		d.Set("scope", retPolicy.Scope),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting resource fields: %w", err)
	}

	return nil
}
