package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/policies"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLBPolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBPolicyV3Create,
		ReadContext:   resourceLBPolicyV3Read,
		UpdateContext: resourceLBPoolV3Update,
		DeleteContext: resourceLBPoolV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{},
	}
}

func getRules(d *schema.ResourceData) []policies.Rule {
	rulesRaw := d.Get("rules").([]interface{})
	if len(rulesRaw) == 0 {
		return nil
	}

}

func resourceLBPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := policies.CreateOpts{
		Action:             policies.Action(d.Get("action").(string)),
		Description:        d.Get("description").(string),
		ListenerID:         d.Get("listener_id").(string),
		Name:               d.Get("name").(string),
		Position:           d.Get("position").(int),
		RedirectListenerID: d.Get("redirect_listener_id").(string),
		RedirectPoolID:     d.Get("redirect_pool_id").(string),
		Rules:              getRules(d),
	}

	policy, err := policies.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud pool: %w", err)
	}

	d.SetId(policy.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBPolicyV3Read(clientCtx, d, meta)
}

func resourceLBPolicyV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	policy, err := policies.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "error viewing details of LB Policy v3"))
	}

	mErr := multierror.Append(
		d.Set("name", policy.Name),
		d.Set("description", policy.Description),
		d.Set("lb_algorithm", policy.LBMethod),
		d.Set("project_id", policy.ProjectID),
		d.Set("protocol", policy.Protocol),
		d.Set("session_persistence", expandPersistence(policy.Persistence)),
		d.Set("ip_version", policy.IpVersion),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting LB Policy v3 fields: %w", err)
	}

	return nil
}
