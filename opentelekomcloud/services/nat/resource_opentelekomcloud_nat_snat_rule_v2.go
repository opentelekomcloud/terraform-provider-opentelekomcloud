package nat

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/snatrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNatSnatRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNatSnatRuleV2Create,
		ReadContext:   resourceNatSnatRuleV2Read,
		DeleteContext: resourceNatSnatRuleV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateCIDR,
			},
			"source_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"floating_ip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNatSnatRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	networkID, netOk := d.GetOk("network_id")
	cidr, cidrOk := d.GetOk("cidr")

	if !netOk && !cidrOk {
		return fmterr.Errorf("both network_id and cidr are empty, must specify one of them.")
	}
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	createOpts := &snatrules.CreateOpts{
		NatGatewayID: d.Get("nat_gateway_id").(string),
		NetworkID:    networkID.(string),
		FloatingIPID: d.Get("floating_ip_id").(string),
		SourceType:   d.Get("source_type").(int),
		Cidr:         cidr.(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	snatRule, err := snatrules.Create(NatV2Client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creatting Snat Rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Snat Rule (%s) to become available.", snatRule.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForSnatRuleActive(NatV2Client, snatRule.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Snat Rule: %s", err)
	}

	d.SetId(snatRule.ID)

	return resourceNatSnatRuleV2Read(ctx, d, meta)
}

func resourceNatSnatRuleV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	snatRule, err := snatrules.Get(NatV2Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "Snat Rule"))
	}

	mErr := multierror.Append(
		d.Set("nat_gateway_id", snatRule.NatGatewayID),
		d.Set("network_id", snatRule.NetworkID),
		d.Set("floating_ip_id", snatRule.FloatingIPID),
		d.Set("source_type", snatRule.SourceType),
		d.Set("cidr", snatRule.Cidr),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNatSnatRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	NatV2Client, err := config.NatV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud nat client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSnatRuleDelete(NatV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Snat Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSnatRuleActive(client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := snatrules.Get(client, nId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Snat Rule: %+v", n)
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForSnatRuleDelete(client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Snat Rule %s.\n", nId)

		n, err := snatrules.Get(client, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = snatrules.Delete(client, nId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Snat Rule %s still active.\n", nId)
		return n, "ACTIVE", nil
	}
}
