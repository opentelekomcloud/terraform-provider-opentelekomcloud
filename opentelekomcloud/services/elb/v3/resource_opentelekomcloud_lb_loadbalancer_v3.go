package v3

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
				Set:      schema.HashString,
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
						"ip_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"5_bgp", "5_mailbgp", "5_gray",
							}, false),
						},
						"bandwidth_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"bandwidth_size": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
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
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
				ForceNew:     true,
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
		},
	}
}

func getPublicIp(d *schema.ResourceData) *loadbalancers.PublicIp {
	publicIpRaw := d.Get("public_ip").([]interface{})
	if len(publicIpRaw) == 0 {
		return nil
	}
	publicIpElement := publicIpRaw[0].(map[string]interface{})

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
	createOpts := loadbalancers.CreateOpts{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		VipAddress:           d.Get("vip_address").(string),
		VipSubnetCidrID:      d.Get("subnet_id").(string),
		L4Flavor:             d.Get("l4_flavor").(string),
		VpcID:                d.Get("router_id").(string),
		AvailabilityZoneList: common.ExpandToStringSlice(d.Get("availability_zones").(*schema.Set).List()),
		Tags:                 common.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
		AdminStateUp:         &adminStateUp,
		L7Flavor:             d.Get("l7_flavor").(string),
		PublicIp:             getPublicIp(d),
		ElbSubnetIDs:         common.ExpandToStringSlice(d.Get("network_ids").(*schema.Set).List()),
		IpTargetEnable:       &ipTargetEnable,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	lb, err := loadbalancers.Create(client, createOpts).Extract()
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

	lb, err := loadbalancers.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "loadbalancerV3"))
	}

	log.Printf("[DEBUG] Retrieved loadbalancer %s: %#v", d.Id(), lb)

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
		d.Set("created_at", lb.CreatedAt),
		d.Set("updated_at", lb.UpdatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	tagMap := common.TagsToMap(lb.Tags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud LoadBalancerV3: %s", err)
	}

	return nil
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
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}
	if d.HasChange("network_ids") {
		updateOpts.ElbSubnetIDs = common.ExpandToStringSlice(d.Get("network_ids").(*schema.Set).List())
	}
	if d.HasChange("vip_address") {
		updateOpts.VipAddress = d.Get("vip_address").(string)
	}
	if d.HasChange("l7_flavor") {
		updateOpts.L7Flavor = d.Get("l7_flavor").(string)
	}
	if d.HasChange("l4_flavor") {
		updateOpts.L4Flavor = d.Get("l4_flavor").(string)
	}
	if d.HasChange("subnet_id") {
		subnetID := d.Get("subnet_id").(string)
		updateOpts.VipSubnetCidrID = &subnetID
	}
	if d.HasChange("ip_target_enable") {
		ipTargetEnable := d.Get("ip_target_enable").(bool)
		updateOpts.IpTargetEnable = &ipTargetEnable
	}

	log.Printf("[DEBUG] Updating loadbalancer %s with options: %#v", d.Id(), updateOpts)
	_, err = loadbalancers.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("unable to update LoadBalancerV3 %s: %s", d.Id(), err)
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

	lb, err := loadbalancers.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting loadbalancer %s", d.Id())
	if err := loadbalancers.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("unable to delete LoadBalancerV3 %s: %s", d.Id(), err)
	}

	if len(lb.PublicIps) > 0 {
		config := meta.(*cfg.Config)
		nwV1Client, err := config.NetworkingV1Client(config.GetRegion(d))
		if err != nil {
			return diag.FromErr(err)
		}
		ipIdToDelete := lb.PublicIps[0].PublicIpID
		if err := eips.Delete(nwV1Client, ipIdToDelete).ExtractErr(); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}
