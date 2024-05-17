package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3.0/security"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityLoginPolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityLoginPolicyV3Create,
		ReadContext:   resourceIdentityLoginPolicyV3Read,
		UpdateContext: resourceIdentityLoginPolicyV3Update,
		DeleteContext: resourceIdentityLoginPolicyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_validity_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 240),
			},
			"custom_info_for_login": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"lockout_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(15, 30),
			},
			"login_failed_times": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(3, 10),
			},
			"period_with_login_failures": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(15, 60),
			},
			"session_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(15, 1440),
			},
			"show_recent_login_info": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceIdentityLoginPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	domainID, err := getDomainID(config, client)
	if err != nil {
		return fmterr.Errorf("error getting the domain id, err=%s", err)
	}

	loginPolicyOpts := security.UpdateLoginPolicyOpts{
		AccountValidityPeriod:   pointerto.Int(d.Get("account_validity_period").(int)),
		CustomInfoForLogin:      d.Get("custom_info_for_login").(string),
		LockoutDuration:         d.Get("lockout_duration").(int),
		LoginFailedTimes:        d.Get("login_failed_times").(int),
		PeriodWithLoginFailures: d.Get("period_with_login_failures").(int),
		SessionTimeout:          d.Get("session_timeout").(int),
		ShowRecentLoginInfo:     pointerto.Bool(d.Get("show_recent_login_info").(bool)),
	}
	_, err = security.UpdateLoginAuthPolicy(client, domainID, loginPolicyOpts)
	if err != nil {
		return diag.Errorf("error updating the IAM account login policy: %s", err)
	}

	d.SetId(domainID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceIdentityLoginPolicyV3Read(clientCtx, d, meta)
}

func resourceIdentityLoginPolicyV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	loginPolicy, err := security.GetLoginAuthPolicy(client, d.Id())
	if err != nil {
		return diag.Errorf("error fetching the IAM login policy")
	}

	log.Printf("[DEBUG] Retrieved the IAM operation login policy: %#v", loginPolicy)

	mErr := multierror.Append(nil,
		d.Set("login_failed_times", loginPolicy.LoginFailedTimes),
		d.Set("lockout_duration", loginPolicy.LockoutDuration),
		d.Set("custom_info_for_login", loginPolicy.CustomInfoForLogin),
		d.Set("account_validity_period", loginPolicy.AccountValidityPeriod),
		d.Set("period_with_login_failures", loginPolicy.PeriodWithLoginFailures),
		d.Set("session_timeout", loginPolicy.SessionTimeout),
		d.Set("show_recent_login_info", loginPolicy.ShowRecentLoginInfo),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM policy fields: %s", err)
	}
	return nil
}

func resourceIdentityLoginPolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	if d.HasChanges("account_validity_period", "custom_info_for_login", "lockout_duration",
		"login_failed_times", "period_with_login_failures", "session_timeout", "show_recent_login_info") {
		loginPolicyOpts := security.UpdateLoginPolicyOpts{
			AccountValidityPeriod:   pointerto.Int(d.Get("account_validity_period").(int)),
			CustomInfoForLogin:      d.Get("custom_info_for_login").(string),
			LockoutDuration:         d.Get("lockout_duration").(int),
			LoginFailedTimes:        d.Get("login_failed_times").(int),
			PeriodWithLoginFailures: d.Get("period_with_login_failures").(int),
			SessionTimeout:          d.Get("session_timeout").(int),
			ShowRecentLoginInfo:     pointerto.Bool(d.Get("show_recent_login_info").(bool)),
		}
		_, err = security.UpdateLoginAuthPolicy(client, d.Id(), loginPolicyOpts)
		if err != nil {
			return diag.Errorf("error updating the IAM login policy: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceIdentityLoginPolicyV3Read(clientCtx, d, meta)
}

func resourceIdentityLoginPolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	_, err = security.UpdateLoginAuthPolicy(client, d.Id(),
		security.UpdateLoginPolicyOpts{
			AccountValidityPeriod:   pointerto.Int(0),
			CustomInfoForLogin:      "",
			LockoutDuration:         15,
			LoginFailedTimes:        3,
			PeriodWithLoginFailures: 15,
			SessionTimeout:          1395,
			ShowRecentLoginInfo:     pointerto.Bool(false),
		})
	if err != nil {
		return diag.Errorf("error resetting the IAM login policy: %s", err)
	}

	return nil
}
