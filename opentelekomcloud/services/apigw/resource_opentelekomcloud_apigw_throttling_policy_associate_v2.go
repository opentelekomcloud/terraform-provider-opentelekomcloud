package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	throttlingpolicy "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/tr_policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIThrottlingPolicyAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIGWThrottlingPolicyAssociateV2Create,
		ReadContext:   resourceAPIGWThrottlingPolicyAssociateV2Read,
		UpdateContext: resourceAPIGWThrottlingPolicyAssociateV2Update,
		DeleteContext: resourceAPIGWThrottlingPolicyAssociateV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAPIGWThrottlingPolicyAssociateV2ImportState,
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

func resourceAPIGWThrottlingPolicyAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	opts := throttlingpolicy.BindOpts{
		GatewayID:  gatewayId,
		PolicyID:   policyId,
		PublishIds: common.ExpandToStringListBySet(publishIds),
	}
	_, err = throttlingpolicy.BindPolicy(client, opts)
	if err != nil {
		return diag.Errorf("error binding OpenTelekomCloud apigw throttling policy to the API: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", gatewayId, policyId))

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWThrottlingPolicyAssociateV2Read(clientCtx, d, meta)
}

func resourceAPIGWThrottlingPolicyAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	instanceId := d.Get("gateway_id").(string)
	policyId := d.Get("policy_id").(string)
	policies, err := throttlingpolicy.ListAPIBoundPolicy(client, buildListOpts(instanceId, policyId))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting OpenTelekomCloud apigw throttling policies")
	}
	if len(policies) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	mErr := multierror.Append(nil,
		d.Set("publish_ids", flattenApiPublishIds(policies)),
		d.Set("region", config.GetRegion(d)))
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAPIGWThrottlingPolicyAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var (
		gatewayId      = d.Get("gateway_id").(string)
		policyId       = d.Get("policy_id").(string)
		oldRaw, newRaw = d.GetChange("publish_ids")

		addSet = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		rmSet  = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	)

	if rmSet.Len() > 0 {
		err = unbindPolicy(client, gatewayId, policyId, rmSet)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		opts := throttlingpolicy.BindOpts{
			GatewayID:  gatewayId,
			PolicyID:   policyId,
			PublishIds: common.ExpandToStringListBySet(addSet),
		}
		_, err = throttlingpolicy.BindPolicy(client, opts)
		if err != nil {
			return diag.Errorf("error binding OpenTelekomCloud apigw throttling policy to the API: %v", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceAPIGWThrottlingPolicyAssociateV2Read(clientCtx, d, meta)
}

func resourceAPIGWThrottlingPolicyAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var (
		gatewayId  = d.Get("gateway_id").(string)
		policyId   = d.Get("policy_id").(string)
		publishIds = d.Get("publish_ids").(*schema.Set)
	)

	return diag.FromErr(unbindPolicy(client, gatewayId, policyId, publishIds))
}

func resourceAPIGWThrottlingPolicyAssociateV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <gateway_id>/<policy_id>")
	}

	mErr := multierror.Append(nil,
		d.Set("gateway_id", parts[0]),
		d.Set("policy_id", parts[1]),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("error saving OpenTelekomCloud APIGW throttling policy: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}

func buildListOpts(gatewayID, policyId string) throttlingpolicy.ListBoundOpts {
	return throttlingpolicy.ListBoundOpts{
		GatewayID:  gatewayID,
		ThrottleID: policyId,
	}
}

func flattenApiPublishIds(apiList []throttlingpolicy.ApiThrottle) []string {
	if len(apiList) < 1 {
		return nil
	}

	result := make([]string, len(apiList))
	for i, val := range apiList {
		result[i] = val.PublishID
	}
	return result
}

func unbindPolicy(client *golangsdk.ServiceClient, gatewayId, policyId string, unbindSet *schema.Set) error {
	policies, err := throttlingpolicy.ListAPIBoundPolicy(client, buildListOpts(gatewayId, policyId))
	if err != nil {
		return fmt.Errorf("error getting OpenTelekomCloud apigw throttling policies")
	}

	for _, rm := range unbindSet.List() {
		for _, api := range policies {
			if rm == api.PublishID {
				err = throttlingpolicy.UnbindPolicy(client, gatewayId, api.ThrottleApplyID)
				if err != nil {
					return fmt.Errorf("error unbound OpenTelekomCloud apigw throttling policy from the API: %s", err)
				}
				break
			}
		}
	}
	return nil
}
