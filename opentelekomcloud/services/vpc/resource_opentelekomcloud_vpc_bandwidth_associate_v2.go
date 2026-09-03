package vpc

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	legacyBandwidths "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/bandwidths"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
	bandwidthsV1 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/bandwidths"
	bandwidthsV2 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const dualStackPublicIPType = "5_dualStack"

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
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	d.SetId(d.Get("bandwidth").(string))

	ips := d.Get("floating_ips").(*schema.Set)
	if err := addIPsToBandwidth(client, d, ips); err != nil {
		return diag.FromErr(err)
	}

	return resourceBandwidthAssociateV2Read(ctx, d, meta)
}

func resourceBandwidthAssociateV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	bandwidth, err := findBandwidthV1(client, d.Id(), allGrantedEnterpriseProjects)
	if err != nil {
		return fmterr.Errorf("error getting bandwidth info: %w", err)
	}
	if bandwidth == nil {
		d.SetId("")
		return nil
	}
	ips := make([]string, len(bandwidth.PublicipInfo))
	for i, ipInfo := range bandwidth.PublicipInfo {
		ips[i] = ipInfo.PublicipId
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
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}
	readClient, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	removed, added := common.GetSetChanges(d, "floating_ips")
	if err := removeIPsFromBandwidth(client, readClient, d, removed); err != nil {
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
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}
	readClient, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	ips := d.Get("floating_ips").(*schema.Set)
	if err := removeIPsFromBandwidth(client, readClient, d, ips); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func addIPsToBandwidth(client *golangsdk.ServiceClient, d *schema.ResourceData, ips *schema.Set) error {
	fips, failedIDs, err := filterExistingFloatingIPs(client, ips)
	if err != nil {
		return err
	}
	if fips.Len() == 0 && len(failedIDs) == 0 {
		return nil
	}

	ipOpts := make([]bandwidthsV2.InsertPublicIPInfo, 0, fips.Len()+len(failedIDs))
	for _, ip := range fips.List() {
		ipOpts = append(ipOpts, bandwidthsV2.InsertPublicIPInfo{
			PublicipId: ip.(string),
		})
	}

	for _, id := range failedIDs {
		if err := insertPortToBandwidth(client, d.Id(), id); err != nil {
			return err
		}
	}

	if len(ipOpts) == 0 {
		return nil
	}

	opts := bandwidthsV2.AddEipOpts{PublicipInfo: ipOpts}

	if _, err := bandwidthsV2.AddEip(client, d.Id(), opts); err != nil {
		return fmt.Errorf("error adding IPs to the bandwidth: %w", err)
	}
	return nil
}

func removeIPsFromBandwidth(
	client, readClient *golangsdk.ServiceClient,
	d *schema.ResourceData,
	ips *schema.Set,
) error {
	bandwidth, err := findBandwidthV1(readClient, d.Id(), allGrantedEnterpriseProjects)
	if err != nil {
		return fmt.Errorf("error reading bandwidth before removing IPs: %w", err)
	}
	if bandwidth == nil {
		return nil
	}

	ipInfo, dualStackPorts := removablePublicIPs(bandwidth.PublicipInfo, ips)
	if len(dualStackPorts) > 0 {
		opts := legacyBandwidths.RemoveOpts{
			ChargeMode:   d.Get("backup_charge_mode").(string),
			Size:         d.Get("backup_size").(int),
			PublicIpInfo: dualStackPorts,
		}
		if err := legacyBandwidths.Remove(client, d.Id(), opts).ExtractErr(); err != nil {
			return fmt.Errorf("error removing IPv6 ports from the bandwidth: %w", err)
		}
	}
	if len(ipInfo) > 0 {
		opts := bandwidthsV2.RemoveEipOpts{
			ChargeMode:   d.Get("backup_charge_mode").(string),
			Size:         d.Get("backup_size").(int),
			PublicipInfo: ipInfo,
		}
		if err := bandwidthsV2.RemoveEip(client, d.Id(), opts); err != nil {
			return fmt.Errorf("error removing IPs from the bandwidth: %w", err)
		}
	}
	return nil
}

func removablePublicIPs(
	attached []bandwidthsV1.PublicIpinfo,
	requested *schema.Set,
) ([]bandwidthsV2.RemovePublicIPInfo, []legacyBandwidths.PublicIpInfoRemove) {
	publicIPs := make([]bandwidthsV2.RemovePublicIPInfo, 0, requested.Len())
	dualStackPorts := make([]legacyBandwidths.PublicIpInfoRemove, 0, requested.Len())
	for _, publicIP := range attached {
		if !requested.Contains(publicIP.PublicipId) {
			continue
		}
		if publicIP.PublicipType == dualStackPublicIPType {
			dualStackPorts = append(dualStackPorts, legacyBandwidths.PublicIpInfoRemove{
				PublicIpID:   publicIP.PublicipId,
				PublicIpType: dualStackPublicIPType,
			})
		} else {
			publicIPs = append(publicIPs, bandwidthsV2.RemovePublicIPInfo{PublicipId: publicIP.PublicipId})
		}
	}
	return publicIPs, dualStackPorts
}

func insertPortToBandwidth(client *golangsdk.ServiceClient, bwID, portID string) error {
	associatedPort, err := ports.Get(client, portID).Extract()
	if err != nil {
		return fmt.Errorf("error fetching port %s: %w", portID, err)
	}
	if associatedPort.Ipv6BandwidthId != "" {
		if associatedPort.Ipv6BandwidthId == bwID {
			return nil
		}

		if err := removePortFromBandwidth(client, associatedPort.Ipv6BandwidthId, portID, "bandwidth", 1); err != nil {
			return err
		}
	}

	insertOpts := bandwidthsV2.AddEipOpts{
		PublicipInfo: []bandwidthsV2.InsertPublicIPInfo{
			{
				PublicipId:   portID,
				PublicipType: dualStackPublicIPType,
			},
		},
	}

	if _, err := bandwidthsV2.AddEip(client, bwID, insertOpts); err != nil {
		return fmt.Errorf("error inserting %s into bandwidth %s: %w", portID, bwID, err)
	}

	return nil
}

func removePortFromBandwidth(client *golangsdk.ServiceClient, bwID, portID, chargeMode string, size int) error {
	if _, err := ports.Get(client, portID).Extract(); err != nil {
		return fmt.Errorf("error fetching port %s: %w", portID, err)
	}

	removeOpts := legacyBandwidths.RemoveOpts{
		ChargeMode: chargeMode,
		Size:       size,
		PublicIpInfo: []legacyBandwidths.PublicIpInfoRemove{
			{
				PublicIpID:   portID,
				PublicIpType: dualStackPublicIPType,
			},
		},
	}

	if err := legacyBandwidths.Remove(client, bwID, removeOpts).ExtractErr(); err != nil {
		return fmt.Errorf("error removing %s from bandwidth %s: %w", portID, bwID, err)
	}

	return nil
}

// filterExistingFloatingIPs returns existing floating IP IDs and unresolved IDs from the given set.
func filterExistingFloatingIPs(clientV2 *golangsdk.ServiceClient, ipIDs *schema.Set) (*schema.Set, []string, error) {
	filtered := schema.NewSet(schema.HashString, []interface{}{})
	unresolved := make(map[string]struct{}, ipIDs.Len())

	for _, raw := range ipIDs.List() {
		unresolved[raw.(string)] = struct{}{}
	}

	// check IPs in v2:
	pages, err := floatingips.List(clientV2, floatingips.ListOpts{}).AllPages()
	if err != nil {
		return nil, nil, fmt.Errorf("error listing floating IPs: %w", err)
	}
	fips, err := floatingips.ExtractFloatingIPs(pages)
	if err != nil {
		return nil, nil, fmt.Errorf("error extracting floating IPs: %w", err)
	}
	for _, ip := range fips {
		if id := ip.ID; ipIDs.Contains(id) {
			filtered.Add(id)
			delete(unresolved, id)
		}
	}

	failedIDs := make([]string, 0, len(unresolved))
	for id := range unresolved {
		failedIDs = append(failedIDs, id)
	}
	sort.Strings(failedIDs)

	return filtered, failedIDs, nil
}
