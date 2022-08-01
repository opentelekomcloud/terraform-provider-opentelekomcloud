package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceNetworkingFloatingIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingFloatingIPAssociateV2Create,
		ReadContext:   resourceNetworkingFloatingIPAssociateV2Read,
		DeleteContext: resourceNetworkingFloatingIPAssociateV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"floating_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingFloatingIPAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	floatingIP := d.Get("floating_ip").(string)
	portID := d.Get("port_id").(string)

	floatingIPID, err := resourceNetworkingFloatingIPAssociateV2IP2ID(client, floatingIP)
	if err != nil {
		return fmterr.Errorf("unable to get ID of floating IP: %s", err)
	}

	updateOpts := floatingips.UpdateOpts{
		PortID: &portID,
	}

	log.Printf("[DEBUG] Floating IP Associate Create Options: %#v", updateOpts)

	_, err = floatingips.Update(client, floatingIPID, updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error associating floating IP %s to port %s: %s", floatingIPID, portID, err)
	}

	d.SetId(floatingIPID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceNetworkingFloatingIPAssociateV2Read(clientCtx, d, meta)
}

func resourceNetworkingFloatingIPAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		d.Set("floating_ip", floatingIP.FloatingIP),
		d.Set("port_id", floatingIP.PortID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNetworkingFloatingIPAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	portID := d.Get("port_id").(string)
	updateOpts := floatingips.UpdateOpts{
		PortID: nil,
	}

	log.Printf("[DEBUG] Floating IP Delete Options: %#v", updateOpts)

	_, err = floatingips.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error disassociating floating IP %s from port %s: %s", d.Id(), portID, err)
	}

	return nil
}

func resourceNetworkingFloatingIPAssociateV2IP2ID(client *golangsdk.ServiceClient, floatingIP string) (string, error) {
	listOpts := floatingips.ListOpts{
		FloatingIP: floatingIP,
	}

	allPages, err := floatingips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return "", err
	}

	if len(allFloatingIPs) != 1 {
		return "", fmt.Errorf("unable to determine the ID of %s", floatingIP)
	}

	return allFloatingIPs[0].ID, nil
}
