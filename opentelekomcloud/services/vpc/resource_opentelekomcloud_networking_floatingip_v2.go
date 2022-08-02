package vpc

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

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const (
	PoolID   = "0a2228f2-7f8a-45f1-8e09-9039e1d09975"
	PoolName = "admin_external_net"
)

func ResourceNetworkingFloatingIPV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkFloatingIPV2Create,
		ReadContext:   resourceNetworkFloatingIPV2Read,
		UpdateContext: resourceNetworkFloatingIPV2Update,
		DeleteContext: resourceNetworkFloatingIPV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "admin_external_net",
				ValidateFunc: validation.StringInSlice([]string{
					"admin_external_net",
				}, true),
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkFloatingIPV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	createOpts := FloatingIPCreateOpts{
		floatingips.CreateOpts{
			FloatingNetworkID: PoolID,
			PortID:            d.Get("port_id").(string),
			TenantID:          d.Get("tenant_id").(string),
			FixedIP:           d.Get("fixed_ip").(string),
		},
		common.MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	floatingIP, err := floatingips.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error allocating floating IP: %s", err)
	}

	log.Printf("[DEBUG] Waiting for OpenTelekomCloud Neutron Floating IP (%s) to become available.", floatingIP.ID)

	stateConf := &resource.StateChangeConf{
		Target:       []string{"ACTIVE"},
		Refresh:      waitForFloatingIPActive(client, floatingIP.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for floating IP to become active: %w", err)
	}

	d.SetId(floatingIP.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNetworkFloatingIPV2Read(clientCtx, d, meta)
}

func resourceNetworkFloatingIPV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	floatingIP, err := floatingips.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "floating IP")
	}

	mErr := multierror.Append(
		d.Set("address", floatingIP.FloatingIP),
		d.Set("port_id", floatingIP.PortID),
		d.Set("fixed_ip", floatingIP.FixedIP),
		d.Set("tenant_id", floatingIP.TenantID),
		d.Set("pool", PoolName),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetworkFloatingIPV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	var updateOpts floatingips.UpdateOpts

	if d.HasChange("port_id") {
		portID := d.Get("port_id").(string)
		updateOpts.PortID = &portID
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = floatingips.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating floating IP: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNetworkFloatingIPV2Read(clientCtx, d, meta)
}

func resourceNetworkFloatingIPV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"DELETED"},
		Refresh:      waitForFloatingIPDelete(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		MinTimeout:   3 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Neutron Floating IP: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForFloatingIPActive(client *golangsdk.ServiceClient, fId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		f, err := floatingips.Get(client, fId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Neutron Floating IP: %+v", f)
		if f.Status == "DOWN" || f.Status == "ACTIVE" {
			return f, "ACTIVE", nil
		}

		return f, "", nil
	}
}

func waitForFloatingIPDelete(client *golangsdk.ServiceClient, fId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Floating IP %s.\n", fId)

		f, err := floatingips.Get(client, fId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Floating IP %s", fId)
				return f, "DELETED", nil
			}
			return f, "ACTIVE", err
		}

		err = floatingips.Delete(client, fId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Floating IP %s", fId)
				return f, "DELETED", nil
			}
			return f, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Floating IP %s still active.\n", fId)
		return f, "ACTIVE", nil
	}
}
