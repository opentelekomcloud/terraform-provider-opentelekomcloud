package iam

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
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
						"condition": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: common.ValidateJsonString,
							StateFunc: func(v interface{}) string {
								jsonString, _ := common.NormalizeJsonString(v)
								return jsonString
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

func buildRolePolicy(d *schema.ResourceData) (policies.CreatePolicy, error) {
	customPolicy := policies.CreatePolicy{
		Version: "1.1",
	}

	statements := d.Get("statement").([]interface{})
	res := make([]policies.CreateStatement, len(statements))

	for i, v := range statements {
		statement := v.(map[string]interface{})
		effect := statement["effect"].(string)
		action := statement["action"].([]interface{})
		var refinedActions []string

		for _, s := range action {
			refinedActions = append(refinedActions, s.(string))
		}

		res[i] = policies.CreateStatement{
			Effect: effect,
			Action: refinedActions,
		}

		resourceCheck := statement["resource"].([]interface{})

		if len(resourceCheck) > 0 {
			var refinedResource []string
			for _, s := range resourceCheck {
				refinedResource = append(refinedResource, s.(string))
			}
			res[i].Resource = refinedResource
		}

		conditionCheck := statement["condition"].(string)

		if len(conditionCheck) > 0 {
			var refinedCondition policies.Condition
			conditionBytes := []byte(conditionCheck)
			if err := json.Unmarshal(conditionBytes, &refinedCondition); err != nil {
				return customPolicy, err
			}
			res[i].Condition = refinedCondition
		}
	}

	customPolicy.Statement = res
	return customPolicy, nil
}

func resourceIdentityRoleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV30, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		fmterr.Errorf(clientV30CreationFail, err)
	}

	roleType, fmtErr := buildRoleType(d)
	if fmtErr != nil {
		return fmtErr
	}

	rolePolicy, err := buildRolePolicy(d)
	if err != nil {
		return fmterr.Errorf("error building policy statements: %s", err)
	}

	opts := policies.CreateOpts{
		Description: d.Get("description").(string),
		DisplayName: d.Get("display_name").(string),
		Type:        roleType.(string),
		Policy:      rolePolicy,
	}

	r, err := policies.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating custom role: %s", err)
	}

	d.SetId(r.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV30)
	return resourceIdentityRoleV3Read(clientCtx, d, meta)
}

func resourceIdentityRoleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV30, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientV30CreationFail, err)
	}

	role, err := policies.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "custom IAM Role")
	}

	statements := make([]interface{}, len(role.Policy.Statement))
	for i, statement := range role.Policy.Statement {
		var condition string
		if len(statement.Condition) > 0 {
			jsonOutput, err := json.Marshal(statement.Condition)
			if err != nil {
				return diag.FromErr(err)
			}
			condition = string(jsonOutput)
		}
		statements[i] = map[string]interface{}{
			"effect":    statement.Effect,
			"action":    statement.Action,
			"resource":  statement.Resource,
			"condition": condition,
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
	client, err := common.ClientFromCtx(ctx, keyClientV30, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientV30CreationFail, err)
	}

	roleType, fmtErr := buildRoleType(d)
	if fmtErr != nil {
		return fmtErr
	}

	rolePolicy, err := buildRolePolicy(d)
	if err != nil {
		return fmterr.Errorf("error building policy statements: %s", err)
	}
	opts := policies.CreateOpts{
		Description: d.Get("description").(string),
		DisplayName: d.Get("display_name").(string),
		Type:        roleType.(string),
		Policy:      rolePolicy,
	}

	_, err = policies.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating (IdentityRoleV3: %v): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV30)
	return resourceIdentityRoleV3Read(clientCtx, d, meta)
}

func resourceIdentityRoleV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV30, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientV30CreationFail, err)
	}

	log.Printf("[DEBUG] Deleting Role %q", d.Id())

	if err := policies.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud IAMv3 role: %s", err)
	}

	d.SetId("")
	return nil
}
