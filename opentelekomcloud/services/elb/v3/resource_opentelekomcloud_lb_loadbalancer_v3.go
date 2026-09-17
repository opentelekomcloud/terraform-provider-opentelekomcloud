package v3

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/bandwidths"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/eips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLoadBalancerV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerV3Create,
		ReadContext:   resourceLoadBalancerV3Read,
		UpdateContext: resourceLoadBalancerV3Update,
		DeleteContext: resourceLoadBalancerV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"router_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"subnet_id"},
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"router_id"},
			},
			"network_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_target_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"l4_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l7_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"admin_state_up": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
				ValidateFunc: common.ValidateTrueOnly,
			},
			"public_ip": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"_managed": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ExactlyOneOf: []string{"public_ip.0.id"},
							ForceNew:     true,
						},
						"bandwidth_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							RequiredWith: []string{"public_ip.0.ip_type"},
							ForceNew:     true,
						},
						"bandwidth_size": {
							Type:                  schema.TypeInt,
							Optional:              true,
							Computed:              true,
							ForceNew:              true,
							DiffSuppressOnRefresh: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old >= new
							},
							RequiredWith: []string{"public_ip.0.ip_type"},
							ValidateFunc: validation.IntBetween(0, 99999),
						},
						"bandwidth_charge_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "traffic",
							ValidateFunc: validation.StringInSlice([]string{
								"traffic",
							}, false),
						},
						"bandwidth_share_type": {
							Type:                  schema.TypeString,
							Optional:              true,
							Computed:              true,
							ForceNew:              true,
							DiffSuppressOnRefresh: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								if old == "WHOLE" && new == "PER" {
									return true
								}
								return false
							},
							RequiredWith: []string{"public_ip.0.ip_type"},
							ValidateFunc: validation.StringInSlice([]string{
								"PER", "WHOLE",
							}, false),
						},
					},
				},
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: common.ValidateTags,
			},
			"vip_port_id": {
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
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ipv6_vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_bandwidth_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"guaranteed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"waf_failure_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pools": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listeners": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ipv6_vip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_vip_port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l4_scale_flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"l7_scale_flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"elb_subnet_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen_scene": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"loadbalancer_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"log_topic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_lb_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eips":        loadBalancerIPInfoSchema("eip_id", "eip_address"),
			"global_eips": loadBalancerIPInfoSchema("global_eip_id", "global_eip_address"),
			"autoscaling": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"enable":           {Type: schema.TypeBool, Computed: true},
					"min_l7_flavor_id": {Type: schema.TypeString, Computed: true},
				}},
			},
			"custom_qos_limit": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"l4_connection": {Type: schema.TypeInt, Computed: true},
					"l4_cps":        {Type: schema.TypeInt, Computed: true},
					"l7_connection": {Type: schema.TypeInt, Computed: true},
					"l7_cps":        {Type: schema.TypeInt, Computed: true},
				}},
			},
			"proxy_protocol_extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"vip_address":      {Type: schema.TypeString, Computed: true},
					"ipv6_vip_address": {Type: schema.TypeString, Computed: true},
					"endpoint_id":      {Type: schema.TypeString, Computed: true},
					"endpoint_service_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				}},
			},
		},
	}
}

func loadBalancerIPInfoSchema(idField, addressField string) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			idField:      {Type: schema.TypeString, Computed: true},
			addressField: {Type: schema.TypeString, Computed: true},
			"ip_version": {Type: schema.TypeInt, Computed: true},
		}},
	}
}

func getPublicIp(d *schema.ResourceData) *loadbalancers.PublicIp {
	publicIpRaw := d.Get("public_ip").([]interface{})
	if len(publicIpRaw) == 0 {
		return nil
	}
	publicIpElement := publicIpRaw[0].(map[string]interface{})
	publicIpElement["_managed"] = true
	_ = d.Set("public_ip", publicIpRaw)

	publicIpOpts := &loadbalancers.PublicIp{
		NetworkType: publicIpElement["ip_type"].(string),
		Bandwidth: loadbalancers.Bandwidth{
			Name:       publicIpElement["bandwidth_name"].(string),
			Size:       publicIpElement["bandwidth_size"].(int),
			ChargeMode: publicIpElement["bandwidth_charge_mode"].(string),
			ShareType:  publicIpElement["bandwidth_share_type"].(string),
		},
	}
	return publicIpOpts
}

func resourceLoadBalancerV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	ipTargetEnable := d.Get("ip_target_enable").(bool)
	deletionProtection := d.Get("deletion_protection").(bool)
	createOpts := loadbalancers.CreateOpts{
		Name:                     d.Get("name").(string),
		Description:              d.Get("description").(string),
		VipAddress:               d.Get("vip_address").(string),
		VipSubnetCidrID:          d.Get("subnet_id").(string),
		L4Flavor:                 d.Get("l4_flavor").(string),
		VpcID:                    d.Get("router_id").(string),
		AvailabilityZoneList:     common.ExpandToStringSlice(d.Get("availability_zones").(*schema.Set).List()),
		Tags:                     common.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		AdminStateUp:             &adminStateUp,
		L7Flavor:                 d.Get("l7_flavor").(string),
		ElbSubnetIDs:             common.ExpandToStringSlice(d.Get("network_ids").(*schema.Set).List()),
		IpTargetEnable:           &ipTargetEnable,
		DeletionProtectionEnable: &deletionProtection,
		IpV6VipSubnetID:          d.Get("ipv6_vip_subnet_id").(string),
		EnterpriseProjectID:      d.Get("enterprise_project_id").(string),
		ChargeMode:               d.Get("charge_mode").(string),
		ProtectionStatus:         d.Get("protection_status").(string),
		ProtectionReason:         d.Get("protection_reason").(string),
	}
	if bandwidthID := d.Get("ipv6_bandwidth_id").(string); bandwidthID != "" {
		createOpts.IPV6Bandwidth = &loadbalancers.BandwidthRef{ID: bandwidthID}
	}
	if common.IsAttrSet(d, "guaranteed") {
		guaranteed := d.Get("guaranteed").(bool)
		createOpts.Guaranteed = &guaranteed
	}

	// currently API supports only a single EIP
	if id, ok := d.GetOk("public_ip.0.id"); ok {
		createOpts.PublicIpIDs = []string{id.(string)}
		_ = d.Set("public_ip", []map[string]interface{}{{"_managed": false}})
	} else {
		createOpts.PublicIp = getPublicIp(d)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	lb, err := loadbalancers.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating LoadBalancerV3: %w", err)
	}

	d.SetId(lb.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLoadBalancerV3Read(clientCtx, d, meta)
}

func resourceLoadBalancerV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	lb, err := loadbalancers.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "loadbalancerV3")
	}

	log.Printf("[DEBUG] Retrieved loadbalancer %s: %#v", d.Id(), lb)

	return setLoadBalancerFields(d, meta, lb)
}

func resourceLoadBalancerV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var updateOpts loadbalancers.UpdateOpts
	updateRequired := false
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
		updateRequired = true
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
		updateRequired = true
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
		updateRequired = true
	}
	if d.HasChange("network_ids") {
		updateOpts.ElbSubnetIDs = common.ExpandToStringSlice(d.Get("network_ids").(*schema.Set).List())
		updateRequired = true
	}
	if d.HasChange("vip_address") {
		updateOpts.VipAddress = d.Get("vip_address").(string)
		updateRequired = true
	}
	if d.HasChange("l7_flavor") {
		updateOpts.L7Flavor = d.Get("l7_flavor").(string)
		updateRequired = true
	}
	if d.HasChange("l4_flavor") {
		updateOpts.L4Flavor = d.Get("l4_flavor").(string)
		updateRequired = true
	}
	if d.HasChange("subnet_id") {
		subnetID := d.Get("subnet_id").(string)
		updateOpts.VipSubnetCidrID = &subnetID
		updateRequired = true
	}
	if d.HasChange("ip_target_enable") {
		ipTargetEnable := d.Get("ip_target_enable").(bool)
		updateOpts.IpTargetEnable = &ipTargetEnable
		updateRequired = true
	}
	if d.HasChange("deletion_protection") {
		deletionProtection := d.Get("deletion_protection").(bool)
		updateOpts.DeletionProtectionEnable = &deletionProtection
		updateRequired = true
	}
	if d.HasChange("ipv6_vip_subnet_id") {
		subnetID := d.Get("ipv6_vip_subnet_id").(string)
		updateOpts.IpV6VipSubnetID = &subnetID
		updateRequired = true
	}
	if d.HasChange("ipv6_bandwidth_id") {
		updateOpts.IpV6Bandwidth = &loadbalancers.BandwidthRef{ID: d.Get("ipv6_bandwidth_id").(string)}
		updateRequired = true
	}
	if d.HasChange("protection_status") {
		updateOpts.ProtectionStatus = d.Get("protection_status").(string)
		updateRequired = true
	}
	if d.HasChange("protection_reason") {
		updateOpts.ProtectionReason = d.Get("protection_reason").(string)
		updateRequired = true
	}

	if updateRequired {
		log.Printf("[DEBUG] Updating loadbalancer %s with options: %#v", d.Id(), updateOpts)
		_, err = loadbalancers.Update(client, d.Id(), updateOpts)
		if err != nil {
			return fmterr.Errorf("unable to update LoadBalancerV3 %s: %s", d.Id(), err)
		}
	}

	// update tags by calling v2 api
	if d.HasChange("tags") {
		elbV2Client, err := config.ElbV2Client(config.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating ELB 2.0 client: %s", err)
		}
		tagErr := common.UpdateResourceTags(elbV2Client, d, "loadbalancers", d.Id())
		if tagErr != nil {
			return diag.Errorf("unable to update tags for LoadBalancerV3:%s, err:%s", d.Id(), tagErr)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceLoadBalancerV3Read(clientCtx, d, meta)
}

func resourceLoadBalancerV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting loadbalancer %s", d.Id())
	if err := loadbalancers.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("unable to delete LoadBalancerV3 %s: %s", d.Id(), err)
	}

	if d.Get("public_ip.#").(int) > 0 && d.Get("public_ip.0._managed").(bool) {
		config := meta.(*cfg.Config)
		nwV1Client, err := config.NetworkingV1Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
		}
		ipID := d.Get("public_ip.0.id").(string)
		if err := eips.Delete(nwV1Client, ipID).ExtractErr(); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func setLoadBalancerFields(d *schema.ResourceData, meta interface{}, lb *loadbalancers.LoadBalancer) diag.Diagnostics {
	d.SetId(lb.ID)
	publicIpInfo := make([]map[string]interface{}, len(lb.PublicIps))
	if len(lb.PublicIps) > 0 {
		info, err := getPublicIpInfo(d, meta, lb.PublicIps[0].PublicIpID)
		if err != nil {
			return diag.FromErr(err)
		}
		if v, ok := d.GetOk("public_ip.0._managed"); ok {
			info["_managed"] = v
		}
		publicIpInfo[0] = info
	}
	tagMap := common.TagsToMap(lb.Tags)
	pools := make([]string, len(lb.Pools))
	for i, pool := range lb.Pools {
		pools[i] = pool.ID
	}
	listeners := make([]string, len(lb.Listeners))
	for i, listener := range lb.Listeners {
		listeners[i] = listener.ID
	}
	eips := make([]map[string]interface{}, len(lb.Eips))
	for i, eip := range lb.Eips {
		eips[i] = map[string]interface{}{
			"eip_id": eip.EipID, "eip_address": eip.EipAddress, "ip_version": eip.IpVersion,
		}
	}
	globalEips := make([]map[string]interface{}, len(lb.GlobalEips))
	for i, eip := range lb.GlobalEips {
		globalEips[i] = map[string]interface{}{
			"global_eip_id": eip.GlobalEipID, "global_eip_address": eip.GlobalEipAddress, "ip_version": eip.IpVersion,
		}
	}
	proxyExtensions := make([]map[string]interface{}, len(lb.ProxyProtocolExtensions))
	for i, extension := range lb.ProxyProtocolExtensions {
		proxyExtensions[i] = map[string]interface{}{
			"vip_address": extension.VipAddress, "ipv6_vip_address": extension.IpV6VipAddress,
			"endpoint_id": extension.Extension.EpID, "endpoint_service_id": extension.Extension.EpServiceID,
		}
	}

	mErr := multierror.Append(
		d.Set("name", lb.Name),
		d.Set("description", lb.Description),
		d.Set("vip_address", lb.VipAddress),
		d.Set("vip_port_id", lb.VipPortID),
		d.Set("admin_state_up", lb.AdminStateUp),
		d.Set("router_id", lb.VpcID),
		d.Set("subnet_id", lb.VipSubnetCidrID),
		d.Set("ip_target_enable", lb.IpTargetEnable),
		d.Set("l4_flavor", lb.L4FlavorID),
		d.Set("l7_flavor", lb.L7FlavorID),
		d.Set("availability_zones", lb.AvailabilityZoneList),
		d.Set("network_ids", lb.ElbSubnetIDs),
		d.Set("public_ip", publicIpInfo),
		d.Set("tags", tagMap),
		d.Set("created_at", lb.CreatedAt),
		d.Set("updated_at", lb.UpdatedAt),
		d.Set("deletion_protection", lb.DeletionProtectionEnable),
		d.Set("ipv6_vip_subnet_id", lb.IpV6VipSubnetID),
		d.Set("ipv6_bandwidth_id", lb.IpV6Bandwidth.ID),
		d.Set("guaranteed", lb.Guaranteed),
		d.Set("enterprise_project_id", lb.EnterpriseProjectID),
		d.Set("waf_failure_action", lb.WafFailureAction),
		d.Set("charge_mode", lb.ChargeMode),
		d.Set("protection_status", lb.ProtectionStatus),
		d.Set("protection_reason", lb.ProtectionReason),
		d.Set("provisioning_status", lb.ProvisioningStatus),
		d.Set("operating_status", lb.OperatingStatus),
		d.Set("provider_name", lb.Provider),
		d.Set("project_id", lb.ProjectID),
		d.Set("pools", pools),
		d.Set("listeners", listeners),
		d.Set("ipv6_vip_address", lb.IpV6VipAddress),
		d.Set("ipv6_vip_port_id", lb.IpV6VipPortID),
		d.Set("l4_scale_flavor", lb.L4ScaleFlavorID),
		d.Set("l7_scale_flavor", lb.L7ScaleFlavorID),
		d.Set("elb_subnet_type", lb.ElbSubnetType),
		d.Set("frozen_scene", lb.FrozenScene),
		d.Set("billing_info", lb.BillingInfo),
		d.Set("public_border_group", lb.PublicBorderGroup),
		d.Set("loadbalancer_type", lb.LoadbalancerType),
		d.Set("gateway_flavor_id", lb.GwFlavorID),
		d.Set("instance_type", lb.InstanceType),
		d.Set("instance_id", lb.InstanceID),
		d.Set("log_group_id", lb.LogGroupID),
		d.Set("log_topic_id", lb.LogTopicID),
		d.Set("service_lb_mode", lb.ServiceLBMode),
		d.Set("eips", eips),
		d.Set("global_eips", globalEips),
		d.Set("autoscaling", []map[string]interface{}{{
			"enable": lb.Autoscaling.Enable, "min_l7_flavor_id": lb.Autoscaling.MinL7FlavorID,
		}}),
		d.Set("custom_qos_limit", []map[string]interface{}{{
			"l4_connection": lb.CustomQosLimit.L4.Connection, "l4_cps": lb.CustomQosLimit.L4.CPS,
			"l7_connection": lb.CustomQosLimit.L7.Connection, "l7_cps": lb.CustomQosLimit.L7.CPS,
		}}),
		d.Set("proxy_protocol_extensions", proxyExtensions),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getPublicIpInfo(d *schema.ResourceData, meta interface{}, publicIpID string) (map[string]interface{}, error) {
	config := meta.(*cfg.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %w", err)
	}
	floatingIP, err := eips.Get(client, publicIpID).Extract()
	if err != nil {
		return nil, err
	}
	bandwidth, err := bandwidths.Get(client, floatingIP.BandwidthID).Extract()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":                    floatingIP.ID,
		"address":               floatingIP.PublicAddress,
		"ip_type":               floatingIP.Type,
		"bandwidth_name":        bandwidth.Name,
		"bandwidth_size":        bandwidth.Size,
		"bandwidth_charge_mode": bandwidth.ChargeMode,
		"bandwidth_share_type":  bandwidth.ShareType,
	}, nil
}
