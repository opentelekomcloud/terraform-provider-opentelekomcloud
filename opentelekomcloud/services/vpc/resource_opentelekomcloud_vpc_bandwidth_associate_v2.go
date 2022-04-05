package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/bandwidths"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceBandwidthAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandwidthAssociateV2Create,
		ReadContext:   resourceBandwidthAssociateV2Read,
		UpdateContext: resourceBandwidthAssociateV2Update,
		DeleteContext: resourceBandwidthAssociateV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"floating_ips": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"backup_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "bandwidth",
			},
			"backup_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
		},
	}
}

func resourceBandwidthAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	d.SetId(d.Get("bandwidth").(string))

	ips := d.Get("floating_ips").(*schema.Set)
	if err := addIPsToBandwidth(client, d, ips); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceBandwidthAssociateV2Read(clientCtx, d, meta)
}

func resourceBandwidthAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	bandwidth, err := bandwidths.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting bandwidth info")
	}
	ips := make([]string, len(bandwidth.PublicIpInfo))
	for i, ipInfo := range bandwidth.PublicIpInfo {
		ips[i] = ipInfo.ID
	}
	mErr := multierror.Append(
		d.Set("bandwidth", d.Id()),
		d.Set("floating_ips", ips),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting bandwidth associate fields: %w", err)
	}

	return nil
}

func resourceBandwidthAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	removed, added := common.GetSetChanges(d, "floating_ips")
	if err := removeIPsFromBandwidth(client, d, removed); err != nil {
		return diag.FromErr(err)
	}
	if err := addIPsToBandwidth(client, d, added); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceBandwidthAssociateV2Read(clientCtx, d, meta)
}

func resourceBandwidthAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	ips := d.Get("floating_ips").(*schema.Set)
	if err := removeIPsFromBandwidth(client, d, ips); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func addIPsToBandwidth(client *golangsdk.ServiceClient, d *schema.ResourceData, ips *schema.Set) error {
	ips, err := filterExistingFloatingIPs(client, ips)
	if err != nil {
		return err
	}
	if ips.Len() == 0 {
		return nil
	}

	ipOpts := make([]bandwidths.PublicIpInfoInsertOpts, ips.Len())
	for i, ip := range ips.List() {
		ipOpts[i] = bandwidths.PublicIpInfoInsertOpts{
			PublicIpID: ip.(string),
		}
	}
	opts := bandwidths.InsertOpts{PublicIpInfo: ipOpts}

	if _, err := bandwidths.Insert(client, d.Id(), opts).Extract(); err != nil {
		return fmt.Errorf("error adding IPs to the bandwidth: %w", err)
	}
	return nil
}

func removeIPsFromBandwidth(client *golangsdk.ServiceClient, d *schema.ResourceData, ips *schema.Set) error {
	ips, err := filterExistingFloatingIPs(client, ips)
	if err != nil {
		return err
	}
	if ips.Len() == 0 {
		return nil
	}

	ipInfo := make([]bandwidths.PublicIpInfoID, ips.Len())
	for i, ip := range ips.List() {
		ipInfo[i] = bandwidths.PublicIpInfoID{
			PublicIpID: ip.(string),
		}
	}

	opts := bandwidths.RemoveOpts{
		ChargeMode:   d.Get("backup_charge_mode").(string),
		Size:         d.Get("backup_size").(int),
		PublicIpInfo: ipInfo,
	}

	if err := bandwidths.Remove(client, d.Id(), opts).ExtractErr(); err != nil {
		return fmt.Errorf("error removing IPs from the bandwidth: %w", err)
	}
	return nil
}

// filterExistingFloatingIPs returns only existing IPs from given slice
func filterExistingFloatingIPs(clientV2 *golangsdk.ServiceClient, ipIDs *schema.Set) (*schema.Set, error) {
	filtered := schema.NewSet(schema.HashString, []interface{}{})

	// check IPs in v2:
	pages, err := floatingips.List(clientV2, floatingips.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error listing floating IPs: %w", err)
	}
	fips, err := floatingips.ExtractFloatingIPs(pages)
	if err != nil {
		return nil, fmt.Errorf("error extracting floating IPs: %w", err)
	}
	for _, ip := range fips {
		if id := ip.ID; ipIDs.Contains(id) {
			filtered.Add(id)
		}
	}
	return filtered, nil
}
