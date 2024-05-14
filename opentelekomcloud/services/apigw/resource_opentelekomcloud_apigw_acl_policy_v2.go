package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	acls "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/acl"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIAclPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIAclPolicyV2Create,
		ReadContext:   resourceAPIAclPolicyV2Read,
		UpdateContext: resourceAPIAclPolicyV2Update,
		DeleteContext: resourceAPIAclPolicyV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIAclPolicyV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAPIAclPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := acls.CreateOpts{
		GatewayID:  d.Get("gateway_id").(string),
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		EntityType: d.Get("entity_type").(string),
		Value:      d.Get("value").(string),
	}
	resp, err := acls.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud APIGW ACL policy: %s", err)
	}
	d.SetId(resp.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIAclPolicyV2Read(clientCtx, d, meta)
}

func resourceAPIAclPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	resp, err := acls.Get(client, d.Get("gateway_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ACL policy")
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("type", resp.Type),
		d.Set("name", resp.Name),
		d.Set("entity_type", resp.EntityType),
		d.Set("value", resp.Value),
		d.Set("updated_at", resp.UpdateTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW ACL policy (%s) fields: %s", d.Id(), err)
	}
	return nil
}

func resourceAPIAclPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := acls.CreateOpts{
		GatewayID:  d.Get("gateway_id").(string),
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		EntityType: d.Get("entity_type").(string),
		Value:      d.Get("value").(string),
	}
	_, err = acls.Update(client, d.Id(), opts)
	if err != nil {
		return diag.Errorf("error updating OpenTelekomCloud APIGW ACL policy: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIAclPolicyV2Read(clientCtx, d, meta)
}

func resourceAPIAclPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	if err = acls.Delete(client, d.Get("gateway_id").(string), d.Id()); err != nil {
		return diag.Errorf("unable to delete the OpenTelekomCloud APIGW ACL policy (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceAPIAclPolicyV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<id>', but '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("gateway_id", parts[0])
}
