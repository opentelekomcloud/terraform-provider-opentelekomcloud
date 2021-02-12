package vpc

import (
	"fmt"
	"log"
	"time"

	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/bandwidths"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceVpcEIPV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcEIPV1Create,
		Read:   resourceVpcEIPV1Read,
		Update: resourceVpcEIPV1Update,
		Delete: resourceVpcEIPV1Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		},
	}
}

func resourceVpcEIPV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating NetworkingV1 client: %s", err)
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
		return fmt.Errorf("error allocating EIP: %s", err)
	}

	log.Printf("[DEBUG] Waiting for EIP %#v to become available.", eip)

	timeout := d.Timeout(schema.TimeoutCreate)
	err = WaitForEIPActive(client, eip.ID, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for EIP (%s) to become ready: %s", eip.ID, err)
	}

	err = bindToPort(d, eip.ID, client, timeout)
	if err != nil {
		return fmt.Errorf("error binding eip:%s to port:%s", eip.ID, err)
	}

	d.SetId(eip.ID)

	if err := addNetworkingTags(d, config, "publicips"); err != nil {
		return err
	}

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating NetworkingV1 client: %s", err)
	}

	eip, err := eips.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "eIP")
	}
	bandWidth, err := bandwidths.Get(client, eip.BandwidthID).Extract()
	if err != nil {
		return fmt.Errorf("error fetching bandwidth: %s", err)
	}

	// Set public ip
	publicIP := []map[string]string{
		{
			"type":       eip.Type,
			"ip_address": eip.PublicAddress,
			"port_id":    eip.PortID,
		},
	}
	if err := d.Set("publicip", publicIP); err != nil {
		return err
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
		return err
	}
	if err := d.Set("region", config.GetRegion(d)); err != nil {
		return err
	}

	if err := readNetworkingTags(d, config, "publicips"); err != nil {
		return err
	}

	return nil
}

func resourceVpcEIPV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating NetworkingV1 client: %s", err)
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
			return common.CheckDeleted(d, err, "Error deleting eip")
		}
		_, err = bandwidths.Update(client, eip.BandwidthID, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating bandwidth: %s", err)
		}

	}

	// Update publicip change
	if d.HasChange("publicip") {
		var updateOpts eips.UpdateOpts

		newIPList := d.Get("publicip").([]interface{})
		newMap := newIPList[0].(map[string]interface{})
		updateOpts.PortID = newMap["port_id"].(string)

		log.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err = eips.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating publicip: %s", err)
		}

	}

	// update tags
	if d.HasChange("tags") {
		NetworkingV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
		}

		if err := common.UpdateResourceTags(NetworkingV2Client, d, "publicips", d.Id()); err != nil {
			return fmt.Errorf("error updating tags: %s", err)
		}
	}

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating NetworkingV1 client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = unbindToPort(d, d.Id(), client, timeout)
	if err != nil {
		return fmt.Errorf("error unbinding eip:%s to port: %s", d.Id(), err)
	}

	if err = eips.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("error deleting VPC EIPv1: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    getEIPStatus(client, d.Id()),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error deleting EIP: %s", err)
	}

	d.SetId("")

	return nil
}

func resourcePublicIP(d *schema.ResourceData) eips.PublicIpOpts {
	publicIPRaw := d.Get("publicip").([]interface{})[0].(map[string]interface{})

	publicIpOpts := eips.PublicIpOpts{
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

func bindToPort(d *schema.ResourceData, eipID string, client *golangsdk.ServiceClient, timeout time.Duration) error {
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
	return WaitForEIPActive(client, eipID, timeout)
}

func unbindToPort(d *schema.ResourceData, eipID string, client *golangsdk.ServiceClient, timeout time.Duration) error {
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
		return fmt.Errorf("error unbinding port from EIP: %s", err)
	}
	return WaitForEIPActive(client, eipID, timeout)
}

func WaitForEIPActive(networkingClient *golangsdk.ServiceClient, eipID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    getEIPStatus(networkingClient, eipID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
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
