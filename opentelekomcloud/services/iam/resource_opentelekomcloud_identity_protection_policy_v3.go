package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3.0/security"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityProtectionPolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProtectionPolicyV3Create,
		ReadContext:   resourceIdentityProtectionPolicyV3Read,
		UpdateContext: resourceIdentityProtectionPolicyV3Update,
		DeleteContext: resourceIdentityProtectionPolicyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"enable_operation_protection_policy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceIdentityProtectionPolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	domainID := client.DomainID

	enable := d.Get("enable_operation_protection_policy").(bool)
	opPolicyOpts := security.UpdateProtectionPolicyOpts{
		OperationProtection: pointerto.Bool(enable),
	}
	_, err = security.UpdateOperationProtectionPolicy(client, domainID, opPolicyOpts)
	if err != nil {
		return diag.Errorf("error updating the IAM operation protection policy: %s", err)
	}

	d.SetId(domainID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceIdentityProtectionPolicyV3Read(clientCtx, d, meta)
}

func resourceIdentityProtectionPolicyV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	opPolicy, err := security.GetOperationProtectionPolicy(client, d.Id())
	if err != nil {
		return diag.Errorf("error fetching the IAM operation protection policy")
	}

	log.Printf("[DEBUG] Retrieved the IAM operation protection policy: %#v", opPolicy)

	mErr := multierror.Append(nil,
		d.Set("enable_operation_protection_policy", opPolicy.OperationProtection),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM policy fields: %s", err)
	}
	return nil
}

func resourceIdentityProtectionPolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	if d.HasChange("enable_operation_protection_policy") {
		passPolicyOpts := security.UpdateProtectionPolicyOpts{
			OperationProtection: pointerto.Bool(d.Get("enable_operation_protection_policy").(bool)),
		}
		_, err = security.UpdateOperationProtectionPolicy(client, d.Id(), passPolicyOpts)
		if err != nil {
			return diag.Errorf("error updating the IAM operation protection policy: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceIdentityProtectionPolicyV3Read(clientCtx, d, meta)
}

func resourceIdentityProtectionPolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV30Client()
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	_, err = security.UpdateOperationProtectionPolicy(client, d.Id(),
		security.UpdateProtectionPolicyOpts{OperationProtection: pointerto.Bool(false)},
	)
	if err != nil {
		return diag.Errorf("error resetting the IAM protection policy: %s", err)
	}

	return nil
}
