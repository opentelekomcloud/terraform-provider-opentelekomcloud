package nat

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/dnatrules"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNatDnatRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNatDnatRuleCreate,
		ReadContext:   resourceNatDnatRuleRead,
		DeleteContext: resourceNatDnatRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"floating_ip_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"internal_service_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"nat_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"port_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"private_ip"},
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.IsUUID,
			},
			"private_ip": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"port_id"},
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validation.IsIPAddress,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tcp", "udp", "any", "6", "17", "0",
				}, false),
			},
			"external_service_port": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 65535),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"floating_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNatDnatRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	portID, portOk := d.GetOk("port_id")
	privateIp, privateIpOk := d.GetOk("private_ip")
	if !portOk && !privateIpOk {
		return fmterr.Errorf("both `port_id` and `private_ip` are empty, must specify one of them.")
	}

	externalServicePort := d.Get("external_service_port").(int)
	internalServicePort := d.Get("internal_service_port").(int)
	createOpts := dnatrules.CreateOpts{
		NatGatewayID:        d.Get("nat_gateway_id").(string),
		PortID:              portID.(string),
		PrivateIp:           privateIp.(string),
		InternalServicePort: &internalServicePort,
		FloatingIpID:        d.Get("floating_ip_id").(string),
		ExternalServicePort: &externalServicePort,
		Protocol:            d.Get("protocol").(string),
	}

	rule, err := createRuleWithRetry(ctx, client, createOpts, time.Minute)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DNAT Rule: %w", err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForDnatRuleActive(client, rule.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for DNAT rule to become active: %w", err)
	}

	d.SetId(rule.ID)

	return resourceNatDnatRuleRead(ctx, d, meta)
}

var createRulePollInterval = 5 * time.Second

// createRuleWithRetry retries creation of DNAT rule in case err 400 is received (handling DnatRuleInValidPortID erorr)
// time between requests is set by createRulePollInterval
func createRuleWithRetry(ctx context.Context, client *golangsdk.ServiceClient, opts dnatrules.CreateOpts, timeout time.Duration) (*dnatrules.DnatRule, error) {
	var rule *dnatrules.DnatRule
	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		var err error
		rule, err = dnatrules.Create(client, opts)
		if err != nil {
			// we are retrying DnatRuleInValidPortID which is HTTP 400
			if _, ok := err.(golangsdk.ErrDefault400); ok {
				// will wait for some time here, no need to retry just now
				time.Sleep(createRulePollInterval)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rule, nil
}

func resourceNatDnatRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}
	dnatRule, err := dnatrules.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dnat rule")
	}

	mErr := multierror.Append(
		d.Set("floating_ip_id", dnatRule.FloatingIpId),
		d.Set("internal_service_port", dnatRule.InternalServicePort),
		d.Set("nat_gateway_id", dnatRule.NatGatewayId),
		d.Set("port_id", dnatRule.PortId),
		d.Set("private_ip", dnatRule.PrivateIp),
		d.Set("protocol", dnatRule.Protocol),
		d.Set("external_service_port", dnatRule.ExternalServicePort),
		d.Set("floating_ip_address", dnatRule.FloatingIpAddress),
		d.Set("status", dnatRule.Status),
		d.Set("tenant_id", dnatRule.ProjectId),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNatDnatRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForDnatRuleDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud DNAT Rule: %w", err)
	}

	d.SetId("")

	return nil
}

func waitForDnatRuleActive(client *golangsdk.ServiceClient, dnatRuleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := dnatrules.Get(client, dnatRuleID)
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud DNAT Rule: %+v", n)
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForDnatRuleDelete(client *golangsdk.ServiceClient, dnatRuleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud DNAT Rule %s.\n", dnatRuleID)

		n, err := dnatrules.Get(client, dnatRuleID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud DNAT Rule %s", dnatRuleID)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		if err := dnatrules.Delete(client, dnatRuleID); err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud DNAT Rule %s", dnatRuleID)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud DNAT Rule %s still active.\n", dnatRuleID)
		return n, "ACTIVE", nil
	}
}
