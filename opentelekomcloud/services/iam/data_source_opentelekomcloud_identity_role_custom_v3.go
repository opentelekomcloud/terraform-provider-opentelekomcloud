package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceIdentityRoleCustomV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCustomRoleV3Read,

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"id",
				},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ConflictsWith: []string{
					"id",
				},
				ValidateFunc: validation.StringInSlice([]string{
					"domain", "project",
				}, false),
			},

			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ConflictsWith: []string{
					"display_name",
					"type",
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"statement": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"effect": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"condition": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// dataSourceIdentityRoleV3Read performs the role lookup.
func dataSourceIdentityCustomRoleV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	var roleType string
	if v, ok := d.GetOk("type"); ok {
		if v == "project" {
			roleType = "XA"
		} else {
			roleType = "AX"
		}
	}

	listOpts := policies.ListOpts{
		DisplayName: d.Get("display_name").(string),
		Type:        roleType,
		ID:          d.Get("id").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var role policies.Policy
	allPolicies, err := policies.List(identityClient, listOpts)
	if err != nil {
		return fmterr.Errorf("unable to query custome roles: %s", err)
	}

	if len(allPolicies.Roles) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allPolicies.Roles) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allPolicies.Roles)
		return fmterr.Errorf("your query returned more than one result. Please try a more " +
			"specific search criteria.")
	}
	role = allPolicies.Roles[0]

	log.Printf("[DEBUG] Single Role found: %s", role.ID)
	return diag.FromErr(dataSourceIdentityCustomRoleV3Attributes(d, &role))
}

// dataSourceIdentityCustomRoleV3Attributes populates the fields of a custom Role resource.
func dataSourceIdentityCustomRoleV3Attributes(d *schema.ResourceData, role *policies.Policy) error {
	log.Printf("[DEBUG] opentelekomcloud_identity_role_v3 details: %#v", role)

	d.SetId(role.ID)

	displayLayer := role.Type
	if displayLayer == "AX" {
		displayLayer = "domain"
	} else {
		displayLayer = "project"
	}

	statements, err := buildStatementsSet(role)
	if err != nil {
		return err
	}

	mErr := multierror.Append(
		d.Set("name", role.Name),
		d.Set("display_name", role.DisplayName),
		d.Set("domain_id", role.DomainId),
		d.Set("description", role.Description),
		d.Set("type", displayLayer),
		d.Set("statement", statements),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}

	return nil
}
