package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/bandwidths"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcEIPV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcEIPV1Create,
		ReadContext:   resourceVpcEIPV1Read,
		UpdateContext: resourceVpcEIPV1Update,
		DeleteContext: resourceVpcEIPV1Delete,

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
			"publicip": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"bandwidth": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"unbind_port": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVpcEIPV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := EIPCreateOpts{
		eips.ApplyOpts{
			IP:        resourcePublicIP(d),
			Bandwidth: resourceBandWidth(d),
		},
		common.MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	eip, err := eips.Apply(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error allocating EIP: %w", err)
	}

	log.Printf("[DEBUG] Waiting for EIP %#v to become available.", eip)

	timeout := d.Timeout(schema.TimeoutCreate)
	if err := WaitForEIPActive(ctx, client, eip.ID, timeout); err != nil {
		return fmterr.Errorf("error waiting for EIP (%s) to become ready: %w", eip.ID, err)
	}

	if err := bindToPort(ctx, d, eip.ID, client, timeout); err != nil {
		return fmterr.Errorf("error binding eip: %s to port: %w", eip.ID, err)
	}

	d.SetId(eip.ID)

	if err := addNetworkingTags(d, config, "publicips"); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcEIPV1Read(clientCtx, d, meta)
}

func resourceVpcEIPV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	eip, err := eips.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "eIP")
	}
	bandWidth, err := bandwidths.Get(client, eip.BandwidthID).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching bandwidth: %w", err)
	}

	// Set public ip
	publicIP := []map[string]string{
		{
			"type":       eip.Type,
			"ip_address": eip.PublicAddress,
			"port_id":    eip.PortID,
			"name":       eip.Name,
		},
	}
	if err := d.Set("publicip", publicIP); err != nil {
		return diag.FromErr(err)
	}

	// Set bandwidth
	bw := []map[string]interface{}{
		{
			"name":        bandWidth.Name,
			"size":        eip.BandwidthSize,
			"share_type":  eip.BandwidthShareType,
			"charge_mode": bandWidth.ChargeMode,
		},
	}
	if err := d.Set("bandwidth", bw); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("region", config.GetRegion(d)); err != nil {
		return diag.FromErr(err)
	}

	if err := readNetworkingTags(d, config, "publicips"); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpcEIPV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	// Update bandwidth change
	if d.HasChange("bandwidth") {
		var updateOpts bandwidths.UpdateOpts

		newBWList := d.Get("bandwidth").([]interface{})
		newMap := newBWList[0].(map[string]interface{})
		updateOpts.Size = newMap["size"].(int)
		updateOpts.Name = newMap["name"].(string)

		log.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

		eip, err := eips.Get(client, d.Id()).Extract()
		if err != nil {
			return common.CheckDeletedDiag(d, err, "Error deleting eip")
		}
		_, err = bandwidths.Update(client, eip.BandwidthID, updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating bandwidth: %s", err)
		}
	}

	// Update publicip change
	if d.Get("unbind_port").(bool) {
		timeout := d.Timeout(schema.TimeoutUpdate)
		if err := unbindToPort(ctx, d, d.Id(), client, timeout); err != nil {
			return fmterr.Errorf("error unbinding eip: %s to port: %w", d.Id(), err)
		}
	}

	if d.HasChange("publicip") {
		var updateOpts eips.UpdateOpts
		newIPList := d.Get("publicip").([]interface{})
		newMap := newIPList[0].(map[string]interface{})
		updateOpts.PortID = newMap["port_id"].(string)

		log.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err = eips.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating publicip: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		nwV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf(errCreationV2Client, err)
		}

		if err := common.UpdateResourceTags(nwV2Client, d, "publicips", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcEIPV1Read(clientCtx, d, meta)
}

func resourceVpcEIPV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	if err := unbindToPort(ctx, d, d.Id(), client, timeout); err != nil {
		return fmterr.Errorf("error unbinding eip: %s to port: %w", d.Id(), err)
	}

	if err := eips.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting VPC EIPv1: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    getEIPStatus(client, d.Id()),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting EIP: %w", err)
	}

	d.SetId("")

	return nil
}

func resourcePublicIP(d *schema.ResourceData) eips.PublicIpOpts {
	publicIPRaw := d.Get("publicip").([]interface{})[0].(map[string]interface{})

	publicIpOpts := eips.PublicIpOpts{
		Name:    publicIPRaw["name"].(string),
		Type:    publicIPRaw["type"].(string),
		Address: publicIPRaw["ip_address"].(string),
	}
	return publicIpOpts
}

func resourceBandWidth(d *schema.ResourceData) eips.BandwidthOpts {
	bandwidthRaw := d.Get("bandwidth").([]interface{})[0].(map[string]interface{})

	bandwidthOpts := eips.BandwidthOpts{
		Name:       bandwidthRaw["name"].(string),
		Size:       bandwidthRaw["size"].(int),
		ShareType:  bandwidthRaw["share_type"].(string),
		ChargeMode: bandwidthRaw["charge_mode"].(string),
	}
	return bandwidthOpts
}

func bindToPort(ctx context.Context, d *schema.ResourceData, eipID string, client *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})[0].(map[string]interface{})
	portID, ok := publicIPRaw["port_id"]
	if !ok || portID == "" {
		return nil
	}

	pd := portID.(string)
	log.Printf("[DEBUG] Bind eip: %s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{PortID: pd}
	_, err := eips.Update(client, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}
	return WaitForEIPActive(ctx, client, eipID, timeout)
}

func unbindToPort(ctx context.Context, d *schema.ResourceData, eipID string, client *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})[0].(map[string]interface{})
	portID, ok := publicIPRaw["port_id"]
	if !ok || portID == "" {
		return nil
	}

	pd := portID.(string)
	log.Printf("[DEBUG] Unbind eip: %s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{
		PortID: "",
	}
	_, err := eips.Update(client, eipID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error unbinding port from EIP: %w", err)
	}
	return WaitForEIPActive(ctx, client, eipID, timeout)
}

func WaitForEIPActive(ctx context.Context, client *golangsdk.ServiceClient, eipID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    getEIPStatus(client, eipID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getEIPStatus(client *golangsdk.ServiceClient, eipID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		eip, err := eips.Get(client, eipID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return eip, "DELETED", nil
			}
			return nil, "", err
		}

		if eip.Status == "DOWN" || eip.Status == "ACTIVE" {
			return eip, "ACTIVE", nil
		}

		return eip, "", nil
	}
}
