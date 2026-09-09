package iam

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/tokens"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

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

	d.SetId(d.Get("name").(string))
	if err := d.Set("region", config.GetRegion(d)); err != nil {
		return diag.FromErr(err)
	}

	if tokenID := authScopeTokenID(config); tokenID != "" {
		return authScopeFromToken(d, identityClient, tokenID)
	}
	return authScopeFromAKSK(d, config, identityClient)
}

func authScopeTokenID(config *cfg.Config) string {
	if config.Token != "" {
		return config.Token
	}
	if config.HwClient != nil && config.HwClient.Token() != "" {
		return config.HwClient.Token()
	}
	if config.DomainClient != nil {
		return config.DomainClient.Token()
	}
	return ""
}

func authScopeFromToken(d *schema.ResourceData, client *golangsdk.ServiceClient, tokenID string) diag.Diagnostics {
	result := tokens.Get(client, tokenID)
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

	return nil
}

// authScopeFromAKSK resolves the auth scope for AK/SK authentication, which
// provides no token to introspect. The owning user is looked up by the access
// key, and the remaining scope is taken from the authenticated clients.
func authScopeFromAKSK(d *schema.ResourceData, config *cfg.Config, client *golangsdk.ServiceClient) diag.Diagnostics {
	// The identity client is derived from the domain client, so the domain client's
	// options describe the credentials the lookup below is actually signed with.
	// They differ from the provider configuration when an agency is assumed.
	akskOpts := golangsdk.AKSKAuthOptions{}
	if config.DomainClient != nil {
		akskOpts = config.DomainClient.AKSKOptions()
	}

	if config.AgencyName != "" && config.AgencyDomainName != "" {
		return fmterr.Errorf("unable to determine the auth scope: an agency assumed with AK/SK authentication is " +
			"backed by temporary credentials which can't be resolved to a user, use username/password or " +
			"token authentication")
	}
	if config.SecurityToken != "" || akskOpts.SecurityToken != "" {
		return fmterr.Errorf("unable to determine the auth scope: temporary AK/SK credentials are not queryable, " +
			"use permanent AK/SK, username/password or token authentication")
	}
	accessKey := akskOpts.AccessKey
	if accessKey == "" {
		accessKey = config.AccessKey
	}
	if accessKey == "" {
		return fmterr.Errorf("unable to determine the auth scope: neither a token nor an access key is available")
	}

	credential, err := credentials.Get(client, accessKey).Extract()
	if err != nil {
		return fmterr.Errorf("error retrieving the user of the configured access key: %s", err)
	}

	projectOpts := golangsdk.AKSKAuthOptions{}
	if config.HwClient != nil {
		projectOpts = config.HwClient.AKSKOptions()
	}
	domainName := config.DomainName
	if domainName == "" {
		domainName = akskOpts.Domain
	}

	domainID := config.GetDomainID()

	var diags diag.Diagnostics
	userName := ""
	user, err := users.Get(client, credential.UserID).Extract()
	if err != nil {
		log.Printf("[WARN] Unable to retrieve details of user %s: %s", credential.UserID, err)
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Unable to retrieve the name of the current user",
			Detail: fmt.Sprintf("The IAM user owning the configured access key can't be read, so `user_name` is "+
				"left empty. Querying an IAM user requires the corresponding IAM permissions: %s", err),
		})
	} else {
		userName = user.Name
		if user.DomainID != "" {
			domainID = user.DomainID
		}
	}

	mErr := multierror.Append(nil,
		d.Set("user_id", credential.UserID),
		d.Set("user_name", userName),
		d.Set("user_domain_id", domainID),
		d.Set("user_domain_name", domainName),
		// AK/SK is always project-scoped, so the domain scope stays empty.
		d.Set("domain_id", ""),
		d.Set("domain_name", ""),
		d.Set("project_id", projectOpts.ProjectId),
		d.Set("project_name", projectOpts.ProjectName),
		d.Set("project_domain_id", domainID),
		d.Set("project_domain_name", domainName),
		// Role assignments are only exposed by token introspection.
		d.Set("roles", []map[string]string{}),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return diags
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
