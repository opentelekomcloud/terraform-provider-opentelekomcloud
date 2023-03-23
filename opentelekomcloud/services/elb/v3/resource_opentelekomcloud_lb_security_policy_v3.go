package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/security_policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLBSecurityPolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLBSecurityPolicyV3Create,
		ReadContext:   resourceLBSecurityPolicyV3Read,
		UpdateContext: resourceLBSecurityPolicyV3Update,
		DeleteContext: resourceLBSecurityPolicyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocols": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ciphers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getAllElems(d *schema.ResourceData, parameter string) []string {
	rawElems := d.Get(parameter).([]interface{})
	elems := make([]string, len(rawElems))
	for i, raw := range rawElems {
		elems[i] = raw.(string)
	}

	return elems
}

func flattenListeners(policy security_policy.SecurityPolicy) []string {
	var listenerIDs []string

	for _, elem := range policy.SecurityPolicy.Listeners {
		listenerIDs = append(listenerIDs, elem.ID)
	}

	return listenerIDs
}

func resourceLBSecurityPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := security_policy.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Protocols:   getAllElems(d, "protocols"),
		Ciphers:     getAllElems(d, "ciphers"),
	}

	policy, err := security_policy.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud security policy: %w", err)
	}

	d.SetId(policy.SecurityPolicy.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBSecurityPolicyV3Read(clientCtx, d, meta)
}

func resourceLBSecurityPolicyV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	policy, err := security_policy.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error viewing details of LB Security Policy v3")
	}

	mErr := multierror.Append(
		d.Set("name", policy.SecurityPolicy.Name),
		d.Set("description", policy.SecurityPolicy.Description),
		d.Set("created_at", policy.SecurityPolicy.CreatedAt),
		d.Set("updated_at", policy.SecurityPolicy.UpdatedAt),
		d.Set("project_id", policy.SecurityPolicy.ProjectId),
		d.Set("protocols", policy.SecurityPolicy.Protocols),
		d.Set("ciphers", policy.SecurityPolicy.Ciphers),
		d.Set("listener_ids", flattenListeners(*policy)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting LB Security Policy v3 fields: %w", err)
	}

	return nil
}

func resourceLBSecurityPolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := security_policy.UpdateOpts{}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = description
	}
	if d.HasChange("protocols") {
		opts.Protocols = getAllElems(d, "protocols")
	}
	if d.HasChange("ciphers") {
		opts.Ciphers = getAllElems(d, "ciphers")
	}

	_, err = security_policy.Update(client, opts, d.Id())
	if err != nil {
		return fmterr.Errorf("error updating LB Security Policy v3: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLBSecurityPolicyV3Read(clientCtx, d, meta)
}

func resourceLBSecurityPolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if err := security_policy.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("error deleting LB Policy v3: %w", err)
	}

	return nil
}
