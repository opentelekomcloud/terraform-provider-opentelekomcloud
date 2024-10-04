package vpn

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/gateway"
	vpntags "github.com/opentelekomcloud/gophertelekomcloud/openstack/evpn/v5/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceEnterpriseVpnGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEvpnGatewayCreate,
		UpdateContext: resourceEvpnGatewayUpdate,
		ReadContext:   resourceEvpnGatewayRead,
		DeleteContext: resourceEvpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vpc",
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vpc", "er",
				}, false),
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"local_subnets": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Computed: true,
			},
			"connect_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"er_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ha_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"active-active", "active-standby"}, false),
			},
			"eip1": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Elem:         GatewayEipSchema(),
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"eip2"},
			},
			"eip2": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Elem:         GatewayEipSchema(),
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"eip1"},
			},
			"access_vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"access_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  64512,
				ForceNew: true,
			},
			"access_private_ip_1": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"access_private_ip_2"},
			},
			"access_private_ip_2": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				RequiredWith: []string{"access_private_ip_1"},
			},
			"tags": common.TagsSchema(),
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
			"er_attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_connection_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_connection_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func GatewayEipSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"bandwidth", "traffic",
				}, false),
			},

			"bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceEvpnGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	var zones []string
	azRaw := d.Get("availability_zones").([]interface{})
	for _, az := range azRaw {
		zones = append(zones, az.(string))
	}
	gatewayTags := d.Get("tags").(map[string]interface{})
	var tagSlice []tags.ResourceTag
	for k, v := range gatewayTags {
		tagSlice = append(tagSlice, tags.ResourceTag{Key: k, Value: v.(string)})
	}
	createOpts := gateway.CreateOpts{
		Name:                d.Get("name").(string),
		NetworkType:         d.Get("network_type").(string),
		AttachmentType:      d.Get("attachment_type").(string),
		ErId:                d.Get("er_id").(string),
		VpcId:               d.Get("vpc_id").(string),
		LocalSubnets:        buildLocalSubnets(d),
		ConnectSubnet:       d.Get("connect_subnet").(string),
		BgpAsn:              pointerto.Int(d.Get("asn").(int)),
		Flavor:              d.Get("flavor").(string),
		AvailabilityZoneIds: zones,
		Eip1:                buildCreateEvpnGatewayEIP(d, "eip1"),
		Eip2:                buildCreateEvpnGatewayEIP(d, "eip2"),
		AccessVpcId:         d.Get("access_vpc_id").(string),
		AccessSubnetId:      d.Get("access_subnet_id").(string),
		HaMode:              d.Get("ha_mode").(string),
		AccessPrivateIp1:    d.Get("access_private_ip_1").(string),
		AccessPrivateIp2:    d.Get("access_private_ip_2").(string),
		Tags:                tagSlice,
	}

	n, err := gateway.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud EVPN gateway: %w", err)
	}

	d.SetId(n.ID)
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"ACTIVE"},
		Refresh:      waitForGatewayActive(client, n.ID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
		MinTimeout:   3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for OpenTelekomCloud EVPN gateway (%s) to become ACTIVE: %w", n.ID, err)
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV5)
	return resourceEvpnGatewayRead(clientCtx, d, meta)
}

func buildLocalSubnets(d *schema.ResourceData) []string {
	var localSubnets []string
	subRaw := d.Get("local_subnets").([]interface{})
	for _, s := range subRaw {
		localSubnets = append(localSubnets, s.(string))
	}
	return localSubnets
}

func buildCreateEvpnGatewayEIP(d *schema.ResourceData, param string) *gateway.Eip {
	if rawArray, ok := d.Get(param).([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		eip := &gateway.Eip{
			ID:            raw["id"].(string),
			Type:          raw["type"].(string),
			ChargeMode:    raw["charge_mode"].(string),
			BandwidthSize: raw["bandwidth_size"].(int),
			BandwidthName: raw["bandwidth_name"].(string),
		}
		return eip
	}
	return nil
}

func resourceEvpnGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}
	gw, err := gateway.Get(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud EVPN gateway (%s): %s", d.Id(), err)
	}

	tagsMap := make(map[string]string)
	for _, tag := range gw.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	mErr := multierror.Append(
		nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("attachment_type", gw.AttachmentType),
		d.Set("availability_zones", gw.AvailabilityZoneIds),
		d.Set("asn", gw.BgpAsn),
		d.Set("connect_subnet", gw.ConnectSubnet),
		d.Set("created_at", gw.CreatedAt),
		d.Set("flavor", gw.Flavor),
		d.Set("local_subnets", gw.LocalSubnets),
		d.Set("ha_mode", gw.HaMode),
		d.Set("eip1", flattenEvpGatewayResponseEip(gw.Eip1)),
		d.Set("name", gw.Name),
		d.Set("eip2", flattenEvpGatewayResponseEip(gw.Eip2)),
		d.Set("status", gw.Status),
		d.Set("updated_at", gw.UpdatedAt),
		d.Set("used_connection_group", gw.UsedConnectionGroup),
		d.Set("used_connection_number", gw.UsedConnectionNumber),
		d.Set("vpc_id", gw.VpcId),
		d.Set("access_vpc_id", gw.AccessVpcId),
		d.Set("access_subnet_id", gw.AccessSubnetId),
		d.Set("er_id", gw.ErId),
		d.Set("network_type", gw.NetworkType),
		d.Set("access_private_ip_1", gw.AccessPrivateIp1),
		d.Set("access_private_ip_2", gw.AccessPrivateIp2),
		d.Set("tags", tagsMap),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEvpGatewayResponseEip(resp gateway.EipResp) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"bandwidth_id":   resp.BandwidthId,
			"bandwidth_name": resp.BandwidthName,
			"bandwidth_size": resp.BandwidthSize,
			"charge_mode":    resp.ChargeMode,
			"id":             resp.ID,
			"ip_address":     resp.IpAddress,
			"ip_version":     resp.IpVersion,
			"type":           resp.Type,
		},
	}
	return rst
}

func resourceEvpnGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	updateChanges := []string{
		"name",
		"local_subnets",
		"eip1",
		"eip2",
	}

	if d.HasChanges(updateChanges...) {
		opts := gateway.UpdateOpts{
			GatewayID:    d.Id(),
			Name:         d.Get("name").(string),
			LocalSubnets: buildLocalSubnets(d),
			Eip1:         buildCreateEvpnGatewayEIP(d, "eip1"),
			Eip2:         buildCreateEvpnGatewayEIP(d, "eip2"),
		}
		_, err = gateway.Update(client, opts)
		if err != nil {
			return diag.Errorf("error updating OpenTelekomCloud EVPN gateway: %s", err)
		}
	}

	// update tags
	if d.HasChange("tags") {
		if err = updateTags(client, d, "vpn-gateway", d.Id()); err != nil {
			return diag.Errorf("error updating tags of OpenTelekomCloud EVPN gateway (%s): %s", d.Id(), err)
		}
	}
	clientCtx := common.CtxWithClient(ctx, client, keyClientV5)
	return resourceEvpnGatewayRead(clientCtx, d, meta)
}

func updateTags(client *golangsdk.ServiceClient, d *schema.ResourceData, resourceType, id string) error {
	if d.HasChange("tags") {
		oldMapRaw, newMapRaw := d.GetChange("tags")
		oldMap := oldMapRaw.(map[string]interface{})
		newMap := newMapRaw.(map[string]interface{})

		// remove old tags
		if len(oldMap) > 0 {
			tagList := common.ExpandResourceTags(oldMap)
			err := vpntags.Delete(client, resourceType, id, vpntags.TagsOpts{Tags: tagList})
			if err != nil {
				return err
			}
		}

		// set new tags
		if len(newMap) > 0 {
			tagList := common.ExpandResourceTags(newMap)
			err := vpntags.Create(client, resourceType, id, vpntags.TagsOpts{Tags: tagList})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceEvpnGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV5, func() (*golangsdk.ServiceClient, error) {
		return config.EvpnV5Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV5Client, err)
	}

	err = gateway.Delete(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud EVPN gateway")
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: waitForGatewayDeletion(client, d.Id()),
		Timeout: d.Timeout(schema.TimeoutCreate),
		Delay:   10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for deleting OpenTelekomCloud EVPN gateway (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func waitForGatewayActive(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := gateway.Get(client, id)
		if err != nil {
			return nil, "", err
		}
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}
		return n, "CREATING", nil
	}
}

func waitForGatewayDeletion(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := gateway.Get(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] The OpenTelekomCloud EVPN gateway has been deleted (ID:%s).", id)
				return r, "DELETED", nil
			}
			return nil, "ERROR", err
		}
		switch r.Status {
		case "ACTIVE", "PENDING_DELETE":
			return r, "DELETING", nil
		default:
			err = fmt.Errorf("error deleting OpenTelekomCloud EVPN gateway[%s]. "+
				"Unexpected status: %v", r.ID, r.Status)
			return r, "ERROR", err
		}
	}
}
