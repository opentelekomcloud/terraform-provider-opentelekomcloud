package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	acls "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/acl"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAclPolicyAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAclPolicyAssociateV2Create,
		ReadContext:   resourceAclPolicyAssociateV2Read,
		UpdateContext: resourceAclPolicyAssociateV2Update,
		DeleteContext: resourceAclPolicyAssociateV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAclPolicyAssociateV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"publish_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAclPolicyAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)
	policyId := d.Get("policy_id").(string)
	publishIds := d.Get("publish_ids").(*schema.Set)

	opts := acls.BindOpts{
		GatewayID:  gatewayId,
		PolicyID:   policyId,
		PublishIds: common.ExpandToStringListBySet(publishIds),
	}
	_, err = acls.BindPolicy(client, opts)
	if err != nil {
		return diag.Errorf("error binding policy to the OpenTelekomCloud APIGW API: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", gatewayId, policyId))

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAclPolicyAssociateV2Read(clientCtx, d, meta)
}

func resourceAclPolicyAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	gatewayId := d.Get("gateway_id").(string)
	policyId := d.Get("policy_id").(string)

	resp, err := acls.ListAPIBoundPolicy(client, buildAclPolicyListOpts(gatewayId, policyId))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting OpenTelekomCloud apigw ACL policy association")
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	return diag.FromErr(d.Set("publish_ids", flattenApiPublishIdsForAclPolicy(resp)))
}

func resourceAclPolicyAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)
	policyId := d.Get("policy_id").(string)
	oldRaw, newRaw := d.GetChange("publish_ids")

	addSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
	rmSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))

	if rmSet.Len() > 0 {
		err = unbindAclPolicy(client, gatewayId, policyId, rmSet)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		opt := acls.BindOpts{
			GatewayID:  gatewayId,
			PolicyID:   policyId,
			PublishIds: common.ExpandToStringListBySet(addSet),
		}
		_, err = acls.BindPolicy(client, opt)
		if err != nil {
			return diag.Errorf("error binding published APIs to the OpenTelekomCloud APIGW  ACL policy (%s): %s", policyId, err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAclPolicyAssociateV2Read(clientCtx, d, meta)
}

func resourceAclPolicyAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	gatewayId := d.Get("gateway_id").(string)
	policyId := d.Get("policy_id").(string)
	publishIds := d.Get("publish_ids").(*schema.Set)

	return diag.FromErr(unbindAclPolicy(client, gatewayId, policyId, publishIds))
}

func resourceAclPolicyAssociateV2ImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<policy_id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("gateway_id", parts[0]),
		d.Set("policy_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

func unbindAclPolicy(client *golangsdk.ServiceClient, gatewayId, policyId string, unbindSet *schema.Set) error {
	policies, err := acls.ListAPIBoundPolicy(client, buildAclPolicyListOpts(gatewayId, policyId))
	if err != nil {
		return fmt.Errorf("error getting binding APIs based on OpenTelekomCloud APIGW ACL policy (%s): %s", policyId, err)
	}
	for _, rm := range unbindSet.List() {
		for _, api := range policies {
			if rm == api.PublishId {
				err = acls.UnbindPolicy(client, gatewayId, api.BindingId)
				if err != nil {
					return fmt.Errorf("error unbound OpenTelekomCloud APIGW ACL policy from the API: %s", err)
				}
				break
			}
		}
	}
	return nil
}

func buildAclPolicyListOpts(gatewayId, policyId string) acls.ListBoundOpts {
	return acls.ListBoundOpts{
		GatewayID: gatewayId,
		ID:        policyId,
		Limit:     500,
	}
}

func flattenApiPublishIdsForAclPolicy(apiList []acls.ApiAcl) []string {
	if len(apiList) < 1 {
		return nil
	}

	result := make([]string, len(apiList))
	for i, val := range apiList {
		result[i] = val.PublishId
	}
	return result
}
