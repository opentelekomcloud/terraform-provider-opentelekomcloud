package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/tokens"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityAuthScopeV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityAuthScopeV3Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// computed attributes
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"project_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityAuthScopeV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV3Client("")
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud identity client: %s", err)
	}
	tokenID := config.Token

	d.SetId(d.Get("name").(string))

	result := tokens.Get(identityClient, tokenID)
	if result.Err != nil {
		return diag.FromErr(result.Err)
	}

	user, err := result.ExtractUser()
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("user_name", user.Name),
		d.Set("user_id", user.ID),
		d.Set("user_domain_name", user.Domain.Name),
		d.Set("user_domain_id", user.Domain.ID),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	domain, err := result.ExtractDomain()
	if err != nil {
		return diag.FromErr(err)
	}
	if domain != nil {
		mErr = multierror.Append(mErr,
			d.Set("domain_name", domain.Name),
			d.Set("domain_id", domain.ID),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("domain_name", ""),
			d.Set("domain_id", ""),
		)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	project, err := result.ExtractProject()
	if err != nil {
		return diag.FromErr(err)
	}
	if project != nil {
		mErr = multierror.Append(mErr,
			d.Set("project_name", project.Name),
			d.Set("project_id", project.ID),
			d.Set("project_domain_name", project.Domain.Name),
			d.Set("project_domain_id", project.Domain.ID),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("project_name", ""),
			d.Set("project_id", ""),
			d.Set("project_domain_name", ""),
			d.Set("project_domain_id", ""),
		)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	roles, err := result.ExtractRoles()
	if err != nil {
		return diag.FromErr(err)
	}

	allRoles := flattenIdentityAuthScopeV3Roles(roles)
	if err := d.Set("roles", allRoles); err != nil {
		log.Printf("[DEBUG] Unable to set opentelekomcloud_identity_auth_scope_v3 roles: %s", err)
	}

	_ = d.Set("region", config.GetRegion(d))

	return nil
}

func flattenIdentityAuthScopeV3Roles(roles []tokens.Role) []map[string]string {
	allRoles := make([]map[string]string, len(roles))

	for i, r := range roles {
		allRoles[i] = map[string]string{
			"role_name": r.Name,
			"role_id":   r.ID,
		}
	}

	return allRoles
}
