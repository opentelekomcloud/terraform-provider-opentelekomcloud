package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityUserV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserV3Create,
		ReadContext:   resourceIdentityUserV3Read,
		UpdateContext: resourceIdentityUserV3Update,
		DeleteContext: resourceIdentityUserV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"default_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"email": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: common.SuppressCaseInsensitive,
			},
			"send_welcome_email": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceIdentityUserV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	enabled := d.Get("enabled").(bool)
	createOpts := users.CreateOpts{
		Name:             d.Get("name").(string),
		DefaultProjectID: d.Get("default_project_id").(string),
		DomainID:         d.Get("domain_id").(string),
		Enabled:          &enabled,
		Description:      d.Get("description").(string),
		Email:            d.Get("email").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Add password here, so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	user, err := users.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud user: %w", err)
	}

	d.SetId(user.ID)

	if d.Get("send_welcome_email").(bool) {
		if err := users.SendWelcomeEmail(client, d.Id()).ExtractErr(); err != nil {
			return fmterr.Errorf("error sending a welcome email: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceIdentityUserV3Read(clientCtx, d, meta)
}

func resourceIdentityUserV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	user, err := users.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "user")
	}

	log.Printf("[DEBUG] Retrieved OpenStack user: %#v", user)

	mErr := multierror.Append(nil,
		d.Set("default_project_id", user.DefaultProjectID),
		d.Set("domain_id", user.DomainID),
		d.Set("enabled", user.Enabled),
		d.Set("name", user.Name),
		d.Set("description", user.Description),
		d.Set("region", config.GetRegion(d)),
		d.Set("email", user.Email),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceIdentityUserV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	var hasChange bool
	var updateOpts users.UpdateOpts

	if d.HasChange("default_project_id") {
		hasChange = true
		updateOpts.DefaultProjectID = d.Get("default_project_id").(string)
	}

	if d.HasChange("domain_id") {
		hasChange = true
		updateOpts.DomainID = d.Get("domain_id").(string)
	}

	if d.HasChange("enabled") {
		hasChange = true
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("email") {
		hasChange = true
		updateOpts.Email = d.Get("email").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	}

	if d.HasChange("password") {
		hasChange = true
		updateOpts.Password = d.Get("password").(string)
	}

	if hasChange {
		_, err := users.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud user: %w", err)
		}
	}

	if d.HasChange("email") && d.Get("send_welcome_email").(bool) {
		if err := users.SendWelcomeEmail(client, d.Id()).ExtractErr(); err != nil {
			return fmterr.Errorf("error sending a welcome email: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceIdentityUserV3Read(clientCtx, d, meta)
}

func resourceIdentityUserV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	err = users.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud user: %w", err)
	}

	return nil
}
