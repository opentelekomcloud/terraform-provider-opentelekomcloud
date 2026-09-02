package vpc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/bandwidths"
	legacyEips "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/publicips"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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
						"ip_version": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntInSlice([]int{4, 6}),
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
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ExactlyOneOf: []string{"bandwidth.0.name"},
						},
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							RequiredWith: []string{"bandwidth.0.size"},
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							RequiredWith: []string{"bandwidth.0.name"},
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_share_bandwidth_types": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	bandwidth, err := resourceBandwidth(d, config)
	if err != nil {
		return diag.FromErr(err)
	}

	var eipID string
	valueSpecs := common.MapValueSpecs(d)
	if len(valueSpecs) > 0 {
		if _, ok := valueSpecs["enterprise_project_id"]; !ok {
			valueSpecs["enterprise_project_id"] = config.GetEnterpriseProjectID(d, "0")
		}
		legacyClient, err := config.NetworkingV1Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf(errCreationV1Client, err)
		}
		publicIPOpts := resourcePublicIP(d)
		var ipVersion string
		if publicIPOpts.IPVersion != 0 {
			ipVersion = strconv.Itoa(publicIPOpts.IPVersion)
		}
		createOpts := EIPCreateOpts{
			ApplyOpts: legacyEips.ApplyOpts{
				IP: legacyEips.PublicIpOpts{
					Name:    publicIPOpts.Alias,
					Type:    publicIPOpts.Type,
					Address: publicIPOpts.IpAddress,
					Version: ipVersion,
				},
				Bandwidth: legacyEips.BandwidthOpts{
					Id:         bandwidth.ID,
					Name:       bandwidth.Name,
					Size:       bandwidth.Size,
					ShareType:  bandwidth.ShareType,
					ChargeMode: bandwidth.ChargeMode,
				},
			},
			ValueSpecs: valueSpecs,
		}
		log.Printf("[DEBUG] Legacy EIP create options for value_specs compatibility: %#v", createOpts)
		eip, err := legacyEips.Apply(legacyClient, createOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error allocating EIP: %w", err)
		}
		eipID = eip.ID
	} else {
		createOpts := publicips.CreateOpts{
			Publicip:            resourcePublicIP(d),
			Bandwidth:           bandwidth,
			EnterpriseProjectId: config.GetEnterpriseProjectID(d, "0"),
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		eip, err := publicips.Create(client, createOpts)
		if err != nil {
			return fmterr.Errorf("error allocating EIP: %w", err)
		}
		eipID = eip.ID
	}

	d.SetId(eipID)

	log.Printf("[DEBUG] Waiting for EIP %s to become available.", eipID)

	timeout := d.Timeout(schema.TimeoutCreate)
	if err := WaitForEIPActive(ctx, client, eipID, timeout); err != nil {
		return fmterr.Errorf("error waiting for EIP (%s) to become ready: %w", eipID, err)
	}

	if err := bindToPort(ctx, d, eipID, client, timeout); err != nil {
		return fmterr.Errorf("error binding eip: %s to port: %w", eipID, err)
	}

	if err := addNetworkingTags(d, config, "publicips"); err != nil {
		return diag.FromErr(err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcEIPV1Read(clientCtx, d, meta)
}

func resourceVpcEIPV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	eip, err := publicips.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "eIP")
	}
	bandwidthClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	bandWidth, err := bandwidths.Get(bandwidthClient, eip.BandwidthId).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching bandwidth: %w", err)
	}

	// Set public ip
	publicIP := []map[string]interface{}{
		{
			"type":       eip.Type,
			"ip_address": eip.PublicIpAddress,
			"ip_version": eip.IPVersion,
			"port_id":    eip.PortId,
			"name":       eip.Alias,
		},
	}
	if err := d.Set("publicip", publicIP); err != nil {
		return diag.FromErr(err)
	}

	// Set bandwidth
	bw := []map[string]interface{}{
		{
			"id":          bandWidth.ID,
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
	if err := d.Set("enterprise_project_id", eip.EnterpriseProjectId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_border_group", eip.PublicBorderGroup); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allow_share_bandwidth_types", eip.AllowShareBandwidthTypes); err != nil {
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
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	// Update bandwidth change
	if d.HasChange("bandwidth") {
		bandwidthClient, err := config.NetworkingV1Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf(errCreationV1Client, err)
		}
		var updateOpts bandwidths.UpdateOpts

		newBWList := d.Get("bandwidth").([]interface{})
		newMap := newBWList[0].(map[string]interface{})
		updateOpts.Size = newMap["size"].(int)
		updateOpts.Name = newMap["name"].(string)

		log.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

		eip, err := publicips.Get(client, d.Id())
		if err != nil {
			return common.CheckDeletedDiag(d, err, "eIP")
		}
		_, err = bandwidths.Update(bandwidthClient, eip.BandwidthId, updateOpts).Extract()
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
		var updateOpts publicips.UpdateOpts
		newIPList := d.Get("publicip").([]interface{})
		newMap := newIPList[0].(map[string]interface{})
		updateOpts.PortId = newMap["port_id"].(string)

		log.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err = publicips.Update(client, d.Id(), updateOpts)
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
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	if err := unbindToPort(ctx, d, d.Id(), client, timeout); err != nil {
		return fmterr.Errorf("error unbinding eip: %s to port: %w", d.Id(), err)
	}

	if err := publicips.Delete(client, d.Id()); err != nil {
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

func resourcePublicIP(d *schema.ResourceData) publicips.PublicIPRequest {
	publicIPRaw := d.Get("publicip").([]interface{})[0].(map[string]interface{})

	publicIpOpts := publicips.PublicIPRequest{
		Alias:     publicIPRaw["name"].(string),
		Type:      publicIPRaw["type"].(string),
		IpAddress: publicIPRaw["ip_address"].(string),
		IPVersion: publicIPRaw["ip_version"].(int),
	}
	return publicIpOpts
}

func resourceBandwidth(d *schema.ResourceData, config *cfg.Config) (publicips.BandWidth, error) {
	bandwidthRaw := d.Get("bandwidth").([]interface{})[0].(map[string]interface{})

	bandwidthOpts := publicips.BandWidth{
		ID:         bandwidthRaw["id"].(string),
		Name:       bandwidthRaw["name"].(string),
		Size:       bandwidthRaw["size"].(int),
		ShareType:  bandwidthRaw["share_type"].(string),
		ChargeMode: bandwidthRaw["charge_mode"].(string),
	}
	if bandwidthOpts.ShareType == "WHOLE" && bandwidthOpts.ID != "" {
		client, err := config.NetworkingV1Client(config.GetRegion(d))
		if err != nil {
			return publicips.BandWidth{}, fmt.Errorf(errCreationV1Client, err)
		}
		bandwidth, err := bandwidths.Get(client, bandwidthOpts.ID).Extract()
		if err != nil {
			return publicips.BandWidth{}, fmt.Errorf("error fetching shared bandwidth: %w", err)
		}
		bandwidthOpts.Name = bandwidth.Name
		bandwidthOpts.Size = bandwidth.Size
		bandwidthOpts.ChargeMode = bandwidth.ChargeMode
	}
	return bandwidthOpts, nil
}

func bindToPort(ctx context.Context, d *schema.ResourceData, eipID string, client *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})[0].(map[string]interface{})
	portID, ok := publicIPRaw["port_id"]
	if !ok || portID == "" {
		return nil
	}

	pd := portID.(string)
	log.Printf("[DEBUG] Bind eip: %s to port: %s", eipID, pd)

	updateOpts := publicips.UpdateOpts{PortId: pd}
	_, err := publicips.Update(client, eipID, updateOpts)
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

	if err := updateEIPPort(client, eipID, ""); err != nil {
		return fmt.Errorf("error unbinding port from EIP: %w", err)
	}
	return WaitForEIPActive(ctx, client, eipID, timeout)
}

func updateEIPPort(client *golangsdk.ServiceClient, eipID, portID string) error {
	if portID != "" {
		_, err := publicips.Update(client, eipID, publicips.UpdateOpts{PortId: portID})
		return err
	}

	body := map[string]interface{}{
		"publicip": map[string]interface{}{"port_id": ""},
	}
	_, err := client.Put(client.ServiceURL(client.ProjectID, "publicips", eipID), body, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return err
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
		eip, err := publicips.Get(client, eipID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return eip, "DELETED", nil
			}
			return nil, "", err
		}

		if eip.Status == "DOWN" || eip.Status == "ACTIVE" || eip.Status == "ELB" || eip.Status == "VPN" {
			return eip, "ACTIVE", nil
		}
		if eip.Status == "BIND_ERROR" || eip.Status == "ERROR" {
			return nil, "", fmt.Errorf("EIP status: %s", eip.Status)
		}

		return eip, "", nil
	}
}
