package v2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/lbaas_v2/listeners"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceL7RuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7RuleV2Create,
		ReadContext:   resourceL7RuleV2Read,
		UpdateContext: resourceL7RuleV2Update,
		DeleteContext: resourceL7RuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceL7RuleV2Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"HOST_NAME", "PATH",
				}, true),
			},

			"compare_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"STARTS_WITH", "EQUAL_TO", "REGEX",
				}, true),
			},

			"l7policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if len(v.(string)) == 0 {
						errors = append(errors, fmt.Errorf("'value' field should not be empty"))
					}
					return
				},
			},

			"key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"admin_state_up": {
				Type:         schema.TypeBool,
				Default:      true,
				Optional:     true,
				ValidateFunc: common.ValidateTrueOnly,
			},
		},
	}
}

func resourceL7RuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	// Assign some required variables for use in creation.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := ""
	ruleType := d.Get("type").(string)
	key := d.Get("key").(string)
	compareType := d.Get("compare_type").(string)
	adminStateUp := d.Get("admin_state_up").(bool)

	// Ensure the right combination of options have been specified.
	if err := checkL7RuleType(ruleType, key); err != nil {
		return fmterr.Errorf("unable to create L7 Rule: %s", err)
	}

	createOpts := l7policies.CreateRuleOpts{
		TenantID:     d.Get("tenant_id").(string),
		RuleType:     l7policies.RuleType(ruleType),
		CompareType:  l7policies.CompareType(compareType),
		Value:        d.Get("value").(string),
		Key:          key,
		AdminStateUp: &adminStateUp,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(client, l7policyID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err = getListenerIDForL7Policy(client, l7policyID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(client, listenerID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent L7 Policy to become active before continuing
	if err := waitForLBV2L7Policy(ctx, client, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create L7 Rule")
	var l7Rule *l7policies.Rule
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		l7Rule, err = l7policies.CreateRule(client, l7policyID, createOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("error creating L7 Rule: %s", err)
	}

	// Wait for L7 Rule to become active before continuing
	if err := waitForLBV2L7Rule(ctx, client, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Rule.ID)
	_ = d.Set("listener_id", listenerID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceL7RuleV2Read(clientCtx, d, meta)
}

func resourceL7RuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	l7policyID := d.Get("l7policy_id").(string)

	l7Rule, err := l7policies.GetRule(client, l7policyID, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "L7 Rule")
	}

	log.Printf("[DEBUG] Retrieved L7 Rule %s: %#v", d.Id(), l7Rule)

	mErr := multierror.Append(
		d.Set("l7policy_id", l7policyID),
		d.Set("type", l7Rule.RuleType),
		d.Set("compare_type", l7Rule.CompareType),
		d.Set("tenant_id", l7Rule.TenantID),
		d.Set("value", l7Rule.Value),
		d.Set("key", l7Rule.Key),
		d.Set("admin_state_up", l7Rule.AdminStateUp),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7RuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	// Assign some required variables for use in updating.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)
	ruleType := d.Get("type").(string)
	key := d.Get("key").(string)

	// Key should always be set
	updateOpts := l7policies.UpdateRuleOpts{
		Key: &key,
	}

	if d.HasChange("compare_type") {
		updateOpts.CompareType = l7policies.CompareType(d.Get("compare_type").(string))
	}
	if d.HasChange("value") {
		updateOpts.Value = d.Get("value").(string)
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	// Ensure the right combination of options have been specified.
	if err := checkL7RuleType(ruleType, key); err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(client, listenerID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(client, l7policyID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to get parent L7 Policy: %s", err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(client, l7policyID, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("unable to get L7 Rule: %s", err)
	}

	// Wait for parent L7 Policy to become active before continuing
	if err := waitForLBV2L7Policy(ctx, client, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout); err != nil {
		return diag.FromErr(err)
	}

	// Wait for L7 Rule to become active before continuing
	if err := waitForLBV2L7Rule(ctx, client, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating L7 Rule %s with options: %#v", d.Id(), updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err := l7policies.UpdateRule(client, l7policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmterr.Errorf("unable to update L7 Rule %s: %s", d.Id(), err)
	}

	// Wait for L7 Rule to become active before continuing
	if err = waitForLBV2L7Rule(ctx, client, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceL7RuleV2Read(clientCtx, d, meta)
}

func resourceL7RuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreationV2Client, err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(client, listenerID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve parent listener (%s) for the L7 Rule: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(client, l7policyID).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve parent L7 Policy (%s) for the L7 Rule: %s", l7policyID, err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(client, l7policyID, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Unable to retrieve L7 Rule")
	}

	// Wait for parent L7 Policy to become active before continuing
	if err := waitForLBV2L7Policy(ctx, client, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to delete L7 Rule %s", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = l7policies.DeleteRule(client, l7policyID, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error deleting L7 Rule")
	}

	if err := waitForLBV2L7Rule(ctx, client, parentListener, parentL7Policy, l7Rule, "DELETED", lbPendingDeleteStatuses, timeout); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// TODO: DROP IT
func resourceL7RuleV2Import(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for L7 Rule. Format must be <policy id>/<rule id>")
		return nil, err
	}

	config := meta.(*cfg.Config)
	client, err := config.ElbV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf(ErrCreationV2Client, err)
	}

	listenerID := ""
	l7policyID := parts[0]
	l7ruleID := parts[1]

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(client, l7policyID).Extract()
	if err != nil {
		return nil, fmt.Errorf("unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err = getListenerIDForL7Policy(client, l7policyID)
		if err != nil {
			return nil, err
		}
	}

	d.SetId(l7ruleID)
	mErr := multierror.Append(
		d.Set("l7policy_id", l7policyID),
		d.Set("listener_id", listenerID),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func checkL7RuleType(ruleType, key string) error {
	keyRequired := []string{"COOKIE", "HEADER"}
	if common.StrSliceContains(keyRequired, ruleType) && key == "" {
		return fmt.Errorf("key attribute is required, when the L7 Rule type is %s", strings.Join(keyRequired, " or "))
	} else if !common.StrSliceContains(keyRequired, ruleType) && key != "" {
		return fmt.Errorf("key attribute must not be used, when the L7 Rule type is not %s", strings.Join(keyRequired, " or "))
	}
	return nil
}
