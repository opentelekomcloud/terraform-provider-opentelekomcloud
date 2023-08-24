package iam

import (
	"context"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3.0/users"
	oldusers "github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"email": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: common.SuppressCaseInsensitive,
				ValidateFunc:     common.ValidateEmail,
			},
			"phone": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"country_code"},
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[0-9]{0,32}$"),
					"the phone number must have a maximum of 32 digits"),
			},
			"country_code": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"phone"},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"pwd_reset": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default", "programmatic", "console",
				}, false),
			},
			"send_welcome_email": {
				Type:         schema.TypeBool,
				Optional:     true,
				RequiredWith: []string{"email"},
			},
			"password_strength": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_login": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityUserV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	domainId, err := getDomainID(config, client)
	if err != nil {
		return fmterr.Errorf("domain name or id missing: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	reset := d.Get("pwd_reset").(bool)
	createOpts := users.CreateOpts{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Email:         d.Get("email").(string),
		Phone:         d.Get("phone").(string),
		AreaCode:      d.Get("country_code").(string),
		AccessMode:    d.Get("access_type").(string),
		Enabled:       &enabled,
		PasswordReset: &reset,
		DomainID:      domainId,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	user, err := users.CreateUser(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating IAM user: %s", err)
	}

	d.SetId(user.ID)

	if d.Get("send_welcome_email").(bool) {
		if err := oldusers.SendWelcomeEmail(client, d.Id()).ExtractErr(); err != nil {
			return fmterr.Errorf("error sending a welcome email: %w", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceIdentityUserV3Read(clientCtx, d, meta)
}

func resourceIdentityUserV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	user, err := users.GetUser(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "user")
	}

	log.Printf("[DEBUG] Retrieved IAM user: %#v", user)
	mErr := multierror.Append(nil,
		d.Set("enabled", user.Enabled),
		d.Set("name", user.Name),
		d.Set("description", user.Description),
		d.Set("email", user.Email),
		d.Set("phone", normalizePhoneNumber(user.Phone)),
		d.Set("country_code", user.AreaCode),
		d.Set("access_type", user.AccessMode),
		d.Set("password_strength", user.PasswordStrength),
		d.Set("pwd_reset", user.PasswordStatus),
		d.Set("create_time", user.CreateAt),
		d.Set("last_login", user.LastLogin),
		d.Set("domain_id", user.DomainID),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM user fields: %s", err)
	}
	return nil
}

func resourceIdentityUserV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	var hasChange bool
	var updateOpts users.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("email") {
		updateOpts.Email = d.Get("email").(string)
		hasChange = true
	}

	if d.HasChanges("country_code", "phone") {
		updateOpts.AreaCode = d.Get("country_code").(string)
		updateOpts.Phone = d.Get("phone").(string)
	}

	if d.HasChange("access_type") {
		updateOpts.AccessMode = d.Get("access_type").(string)
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("pwd_reset") {
		reset := d.Get("pwd_reset").(bool)
		updateOpts.PasswordReset = &reset
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	// Add password here so it wouldn't go in the above log entry
	if d.HasChange("password") {
		updateOpts.Password = d.Get("password").(string)
	}

	_, err = users.ModifyUser(client, d.Id(), updateOpts)
	if err != nil {
		return diag.Errorf("error updating IAM user: %s", err)
	}

	if hasChange {
		if d.Get("send_welcome_email").(bool) {
			if err := oldusers.SendWelcomeEmail(client, d.Id()).ExtractErr(); err != nil {
				return fmterr.Errorf("error sending a welcome email: %w", err)
			}
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

	err = oldusers.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud user: %w", err)
	}

	return nil
}

func normalizePhoneNumber(raw string) string {
	phone := raw

	rawList := strings.Split(raw, "-")
	if len(rawList) > 1 {
		phone = rawList[1]
	}

	return phone
}
