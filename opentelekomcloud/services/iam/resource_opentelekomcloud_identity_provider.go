package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/federation/providers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityProviderV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProviderV3Create,
		ReadContext:   resourceIdentityProviderV3Read,
		UpdateContext: resourceIdentityProviderV3Update,
		DeleteContext: resourceIdentityProviderV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"remote_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"links": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIdentityProviderV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	opts := providers.CreateOpts{
		ID:          d.Get("name").(string),
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
	}

	p, err := providers.Create(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf(providerError, "creating", err)
	}

	d.SetId(p.ID)

	return resourceIdentityProviderV3Read(ctx, d, meta)
}

func resourceIdentityProviderV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	p, err := providers.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmterr.Errorf(providerError, "reading", err)
	}

	mErr := multierror.Append(
		d.Set("name", p.ID),
		d.Set("description", p.Description),
		d.Set("enabled", p.Enabled),
		d.Set("remote_ids", p.RemoteIDs),
		d.Set("links", p.Links),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting identity provider fields: %w", err)
	}

	return nil
}

func resourceIdentityProviderV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	opts := providers.UpdateOpts{}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		opts.Enabled = &enabled
	}

	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}

	_, err = providers.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf(providerError, "updating", err)
	}

	return resourceIdentityProviderV3Read(ctx, d, meta)
}

func resourceIdentityProviderV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	if err := providers.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf(providerError, "deleting", err)
	}

	return nil
}

const providerError = "error %s identity provider v3: %w"
