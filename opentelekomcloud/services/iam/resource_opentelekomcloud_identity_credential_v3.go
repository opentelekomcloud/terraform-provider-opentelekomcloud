package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/credentials"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityCredentialV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCredentialV3Create,
		ReadContext:   resourceIdentityCredentialV3Read,
		UpdateContext: resourceIdentityCredentialV3Update,
		DeleteContext: resourceIdentityCredentialV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_use_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCredentialV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client()
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}

	userID, ok := d.GetOk("user_id")
	if !ok {
		userID = client.UserID
	}

	if userID == "" {
		return fmterr.Errorf("error defining current user ID, please either provide " +
			"`user_id` or authenticate with token auth (not using AK/SK)")
	}

	credential, err := credentials.Create(client, credentials.CreateOpts{
		UserID:      userID.(string),
		Description: d.Get("description").(string),
	}).Extract()
	if err != nil {
		return fmterr.Errorf("error creating AK/SK: %s", err)
	}

	d.SetId(credential.AccessKey)
	_ = d.Set("secret", credential.SecretKey) // secret key returned only once

	return resourceIdentityCredentialV3Read(ctx, d, meta)
}

func resourceIdentityCredentialV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client()
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}
	credential, err := credentials.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IAM credentials")
	}
	mErr := multierror.Append(nil,
		d.Set("user_id", credential.UserID),
		d.Set("access", credential.AccessKey),
		d.Set("status", credential.Status),
		d.Set("create_time", credential.CreateTime),
		d.Set("last_use_time", credential.LastUseTime),
		d.Set("description", credential.Description),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting AK/SK attributes: %s", err)
	}
	return nil
}

func resourceIdentityCredentialV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client()
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}
	opts := credentials.UpdateOpts{}
	if d.HasChange("status") {
		opts.Status = d.Get("status").(string)
	}
	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}
	_, err = credentials.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating AK/SK: %s", err)
	}
	return resourceIdentityCredentialV3Read(ctx, d, meta)
}

func resourceIdentityCredentialV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client()
	if err != nil {
		return fmterr.Errorf("error creating OpenStack identity client: %s", err)
	}
	err = credentials.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting AK/SK: %s", err)
	}
	d.SetId("")
	return nil
}
