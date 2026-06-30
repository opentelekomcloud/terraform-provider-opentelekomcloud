package cc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cc/v3/policy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCcCentralNetworkPolicyApplyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCcCentralNetworkPolicyApplyV3Create,
		UpdateContext: resourceCcCentralNetworkPolicyApplyV3Create,
		ReadContext:   resourceCcCentralNetworkPolicyApplyV3Read,
		DeleteContext: resourceCcCentralNetworkPolicyApplyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCcCentralNetworkPolicyApplyV3ImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"central_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCcCentralNetworkPolicyApplyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	centralNetworkId := d.Get("central_network_id").(string)
	if err = applyCentralNetworkPolicy(ctx, client, centralNetworkId, d.Get("policy_id").(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(centralNetworkId)

	clientCtx := common.CtxWithClient(ctx, client, ccClientV3)
	return resourceCcCentralNetworkPolicyApplyV3Read(clientCtx, d, meta)
}

func applyCentralNetworkPolicy(ctx context.Context, client *golangsdk.ServiceClient, centralNetworkId, policyId string, timeout time.Duration) error {
	if _, err := policy.Apply(client, centralNetworkId, policyId); err != nil {
		return fmt.Errorf("error applying OpenTelekomCloud CC central network policy: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitForCentralNetworkPolicyApplied(client, centralNetworkId, policyId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for OpenTelekomCloud CC central network policy (%s) to be applied: %s", policyId, err)
	}

	return nil
}

func waitForCentralNetworkPolicyApplied(client *golangsdk.ServiceClient, centralNetworkId, policyId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pol, err := getCentralNetworkPolicy(client, centralNetworkId, policyId)
		if err != nil {
			return nil, "", err
		}
		if pol.IsApplied {
			return pol, "COMPLETED", nil
		}
		return pol, "PENDING", nil
	}
}

func resourceCcCentralNetworkPolicyApplyV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	pol, err := getCentralNetworkPolicy(client, d.Id(), d.Get("policy_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud CC central network policy")
	}
	if !pol.IsApplied {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "OpenTelekomCloud CC central network policy is no longer applied")
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("central_network_id", pol.CentralNetworkId),
		d.Set("policy_id", pol.ID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// resourceCcCentralNetworkPolicyApplyV3Delete reverts the central network to its default policy
// (version 1, with no associated enterprise routers), since a policy cannot simply be "unapplied".
func resourceCcCentralNetworkPolicyApplyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, ccClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.CcV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	centralNetworkId := d.Get("central_network_id").(string)
	resp, err := policy.List(client, policy.ListOpts{CentralNetworkId: centralNetworkId})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud CC central network policies")
	}

	var defaultPolicyId string
	for _, pol := range resp.CentralNetworkPolicies {
		if pol.Version == 1 {
			defaultPolicyId = pol.ID
			break
		}
	}
	if defaultPolicyId == "" {
		return diag.Errorf("error reverting OpenTelekomCloud CC central network policy: no default policy (version 1) found")
	}

	if err = applyCentralNetworkPolicy(ctx, client, centralNetworkId, defaultPolicyId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error reverting OpenTelekomCloud CC central network to default policy: %s", err)
	}

	return nil
}

func resourceCcCentralNetworkPolicyApplyV3ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <central_network_id>/<policy_id>")
	}

	d.SetId(parts[0])
	if err := d.Set("policy_id", parts[1]); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
