package dns

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/resolverrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDNSResolverRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSResolverRuleV2Create,
		ReadContext:   resourceDNSResolverRuleV2Read,
		UpdateContext: resourceDNSResolverRuleV2Update,
		DeleteContext: resourceDNSResolverRuleV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[\w.-]{1,64}$`),
					"must be 1 to 64 characters long and contain only letters, digits, '_', '-' and '.'",
				),
			},
			"domain_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: common.SuppressEqualZoneNames,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_addresses": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsIPAddress,
				},
			},
			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
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

func resourceDNSResolverRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	rule, err := resolverrules.Create(client, resolverrules.CreateOpts{
		Name:        d.Get("name").(string),
		DomainName:  d.Get("domain_name").(string),
		EndpointID:  d.Get("endpoint_id").(string),
		IPAddresses: expandDNSResolverRuleV2IPs(d.Get("ip_addresses").(*schema.Set)),
	})
	if err != nil {
		return fmterr.Errorf("error creating DNS resolver rule: %w", err)
	}
	d.SetId(rule.ID)

	timeout := d.Timeout(schema.TimeoutCreate)
	if err := waitForDNSResolverRuleV2Active(ctx, client, rule.ID, schema.NewSet(schema.HashString, nil), timeout); err != nil {
		return fmterr.Errorf("error waiting for DNS resolver rule %s to become ACTIVE: %w", rule.ID, err)
	}

	if err := updateDNSResolverRuleV2VPCs(ctx, client, d, timeout); err != nil {
		return fmterr.Errorf("error associating VPCs with DNS resolver rule %s: %w", rule.ID, err)
	}

	return resourceDNSResolverRuleV2Read(common.CtxWithClient(ctx, client, keyClientV2), d, meta)
}

func resourceDNSResolverRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	rule, err := GetDNSResolverRuleV2(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error reading DNS resolver rule")
	}

	ips := make([]string, 0, len(rule.IPAddresses))
	for _, ip := range rule.IPAddresses {
		ips = append(ips, ip.IP)
	}
	vpcIDs := make([]string, 0, len(rule.Routers))
	for _, r := range rule.Routers {
		vpcIDs = append(vpcIDs, r.RouterID)
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", rule.Name),
		d.Set("domain_name", rule.DomainName),
		d.Set("endpoint_id", rule.EndpointID),
		d.Set("ip_addresses", ips),
		d.Set("vpc_ids", vpcIDs),
		d.Set("status", rule.Status),
		d.Set("created_at", rule.CreateTime),
		d.Set("updated_at", rule.UpdateTime),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDNSResolverRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	if d.HasChanges("name", "ip_addresses") {
		_, err := resolverrules.Update(client, d.Id(), resolverrules.UpdateOpts{
			Name:        d.Get("name").(string),
			IPAddresses: expandDNSResolverRuleV2IPs(d.Get("ip_addresses").(*schema.Set)),
		})
		if err != nil {
			return fmterr.Errorf("error updating DNS resolver rule %s: %w", d.Id(), err)
		}
		// VPCs are not changed yet, so the old ones are expected.
		oldVPCs, _ := d.GetChange("vpc_ids")
		if err := waitForDNSResolverRuleV2Active(ctx, client, d.Id(), oldVPCs.(*schema.Set), timeout); err != nil {
			return fmterr.Errorf("error waiting for DNS resolver rule %s to become ACTIVE: %w", d.Id(), err)
		}
	}

	if d.HasChange("vpc_ids") {
		if err := updateDNSResolverRuleV2VPCs(ctx, client, d, timeout); err != nil {
			return fmterr.Errorf("error updating VPCs of DNS resolver rule %s: %w", d.Id(), err)
		}
	}

	return resourceDNSResolverRuleV2Read(common.CtxWithClient(ctx, client, keyClientV2), d, meta)
}

func resourceDNSResolverRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if err := resolverrules.Delete(client, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DNS resolver rule")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE", "PENDING_DELETE"},
		Target:       []string{"DELETED"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			rule, err := GetDNSResolverRuleV2(client, d.Id())
			if err != nil {
				var notFound golangsdk.ErrDefault404
				if errors.As(err, &notFound) {
					return true, "DELETED", nil
				}
				return nil, "", err
			}
			return rule, rule.Status, nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for DNS resolver rule %s to be deleted: %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}

// GetDNSResolverRuleV2 returns golangsdk.ErrDefault404 for a rule in DELETED status,
// as a deleted rule may still be returned by the API.
func GetDNSResolverRuleV2(client *golangsdk.ServiceClient, id string) (*resolverrules.ResolverRule, error) {
	rule, err := resolverrules.Get(client, id)
	if err != nil {
		return nil, err
	}
	if rule.Status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}
	return rule, nil
}

func expandDNSResolverRuleV2IPs(set *schema.Set) []resolverrules.IPInfo {
	ips := make([]resolverrules.IPInfo, 0, set.Len())
	for _, ip := range set.List() {
		ips = append(ips, resolverrules.IPInfo{IP: ip.(string)})
	}
	return ips
}

// updateDNSResolverRuleV2VPCs associates added and disassociates removed VPCs.
func updateDNSResolverRuleV2VPCs(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout time.Duration) error {
	oldRaw, newRaw := d.GetChange("vpc_ids")
	oldSet, newSet := oldRaw.(*schema.Set), newRaw.(*schema.Set)

	for _, id := range oldSet.Difference(newSet).List() {
		if _, err := resolverrules.DisassociateVPC(client, d.Id(), resolverrules.RouterOpts{RouterID: id.(string)}); err != nil {
			return fmt.Errorf("error disassociating VPC %s: %w", id, err)
		}
	}
	for _, id := range newSet.Difference(oldSet).List() {
		if _, err := resolverrules.AssociateVPC(client, d.Id(), resolverrules.RouterOpts{RouterID: id.(string)}); err != nil {
			return fmt.Errorf("error associating VPC %s: %w", id, err)
		}
	}

	return waitForDNSResolverRuleV2Active(ctx, client, d.Id(), newSet, timeout)
}

// waitForDNSResolverRuleV2Active waits until the rule is ACTIVE and exactly vpcIDs are associated with it and ACTIVE.
func waitForDNSResolverRuleV2Active(ctx context.Context, client *golangsdk.ServiceClient, id string, vpcIDs *schema.Set, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"ACTIVE"},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			rule, err := GetDNSResolverRuleV2(client, id)
			if err != nil {
				return nil, "", err
			}
			if rule.Status == "ERROR" {
				return rule, "", fmt.Errorf("resolver rule is in ERROR status")
			}
			if rule.Status != "ACTIVE" || len(rule.Routers) != vpcIDs.Len() {
				return rule, "PENDING", nil
			}
			for _, r := range rule.Routers {
				if r.Status == "ERROR" {
					return rule, "", fmt.Errorf("association with VPC %s is in ERROR status", r.RouterID)
				}
				if !vpcIDs.Contains(r.RouterID) || r.Status != "ACTIVE" {
					return rule, "PENDING", nil
				}
			}
			return rule, "ACTIVE", nil
		},
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
