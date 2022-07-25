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

func ResourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityRoleV3Create,
		ReadContext:   resourceIdentityRoleV3Read,
		UpdateContext: resourceIdentityRoleV3Update,
		DeleteContext: resourceIdentityRoleV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_layer": {
				Type: schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"domain", "project",
				}, false),
				Required: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"statement": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"effect": {
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"Allow", "Deny",
							}, false),
							Required: true,
						},
						"resource": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"catalog": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildRoleType(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	roleType := d.Get("display_layer")

	if roleType == "domain" {
		return "AX", nil
	} else if roleType == "project" {
		return "XA", nil
	}
	return nil, fmterr.Errorf("unknown display layer:%v", roleType)
}

func buildRolePolicy(d *schema.ResourceData) policies.CreatePolicy {
	customPolicy := policies.CreatePolicy{
		Version: "1.1",
	}

	statements := d.Get("statement").([]interface{})
	res := make([]policies.CreateStatement, len(statements))

	for i, v := range statements {
		statement := v.(map[string]interface{})
		effect := statement["effect"].(string)
		action := statement["action"].([]interface{})
		refinedActions := make([]string, 0)

		for _, s := range action {
			refinedActions = append(refinedActions, s.(string))
		}

		res[i] = policies.CreateStatement{
			Effect: effect,
			Action: refinedActions,
		}

		resourceCheck := statement["resource"].([]interface{})

		if len(resourceCheck) > 0 {
			resource := statement["resource"].([]interface{})
			refinedResource := make([]string, 0)
			for _, s := range resource {
				refinedResource = append(refinedResource, s.(string))
			}
			res[i].Resource = refinedResource
		}
	}

	customPolicy.Statement = res
	return customPolicy
}

func resourceIdentityRoleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	roleType, fmtErr := buildRoleType(d)
	if fmtErr != nil {
		return fmtErr
	}

	opts := policies.CreateOpts{
		Description: d.Get("description").(string),
		DisplayName: d.Get("display_name").(string),
		Type:        roleType.(string),
		Policy:      buildRolePolicy(d),
	}

	r, err := policies.Create(client, opts).Extract()

	if err != nil {
		return fmterr.Errorf("error creating custom role: %s", err)
	}

	d.SetId(r.ID)

	return resourceIdentityRoleV3Read(ctx, d, meta)
}

func resourceIdentityRoleV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	role, err := policies.Get(client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error getting role details: %s", err)
	}

	statements := make([]interface{}, len(role.Policy.Statement))
	for i, statement := range role.Policy.Statement {
		statements[i] = map[string]interface{}{
			"effect":   statement.Effect,
			"action":   statement.Action,
			"resource": statement.Resource,
		}
	}

	displayLayer := role.Type
	if displayLayer == "AX" {
		displayLayer = "domain"
	} else {
		displayLayer = "project"
	}

	mErr := multierror.Append(
		d.Set("description", role.Description),
		d.Set("display_name", role.DisplayName),
		d.Set("name", role.Name),
		d.Set("domain_id", role.DomainId),
		d.Set("catalog", role.Catalog),
		d.Set("statement", statements),
		d.Set("display_layer", displayLayer),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting role fields: %s", err)
	}

	return nil
}

func resourceIdentityRoleV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	needsUpdate := false

	if d.HasChange("description") || d.HasChange("display_name") || d.HasChange("display_layer") ||
		d.HasChange("statement") {
		needsUpdate = true
	}

	if needsUpdate {
		roleType, fmtErr := buildRoleType(d)
		if fmtErr != nil {
			return fmtErr
		}

		opts := policies.CreateOpts{
			Description: d.Get("description").(string),
			DisplayName: d.Get("display_name").(string),
			Type:        roleType.(string),
			Policy:      buildRolePolicy(d),
		}

		_, err = policies.Update(client, d.Id(), opts).Extract()
	}
	if err != nil {
		return fmterr.Errorf("error updating (IdentityRoleV3: %v): %s", d.Id(), err)
	}

	return resourceIdentityRoleV3Read(ctx, d, meta)
}

func resourceIdentityRoleV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmterr.Errorf("error creating identity v3.0 client: %s", err)
	}

	log.Printf("[DEBUG] Deleting Role %q", d.Id())

	if err = policies.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud IAMv3 role: %s", err)
	}

	d.SetId("")
	return nil
}
