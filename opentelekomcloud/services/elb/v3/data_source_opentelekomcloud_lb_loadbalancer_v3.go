package v3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLoadBalancerV3() *schema.Resource {
	dataSource := &schema.Resource{
		ReadContext: dataSourceLoadBalancerV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vip_port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_target_enable": {
				Type:     schema.TypeBool,
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
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth_charge_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth_share_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeMap,
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
			"provisioning_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"guaranteed": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ipv6_vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_vip_port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ipv6_vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"l4_scale_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l7_scale_flavor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},
			"elb_subnet_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"member_device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}

	resourceSchema := ResourceLoadBalancerV3().Schema
	for _, key := range []string{
		"provider_name", "project_id", "pools", "listeners", "ipv6_bandwidth_id",
		"frozen_scene", "billing_info", "public_border_group", "waf_failure_action",
		"charge_mode", "protection_reason", "loadbalancer_type", "gateway_flavor_id",
		"instance_type", "instance_id", "log_group_id", "log_topic_id",
		"service_lb_mode", "eips", "global_eips", "autoscaling", "custom_qos_limit",
		"proxy_protocol_extensions",
	} {
		field := *resourceSchema[key]
		field.Optional = false
		field.Required = false
		field.ForceNew = false
		field.Computed = true
		dataSource.Schema[key] = &field
	}
	return dataSource
}

func dataSourceLoadBalancerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	id := d.Get("id").(string)
	if id != "" {
		lb, err := loadbalancers.Get(client, id)
		if err != nil {
			return fmterr.Errorf("error getting ELBv3 Load Balancer: %w", err)
		}
		return setLoadBalancerFields(d, meta, lb)
	}

	listOpts := loadbalancers.ListOpts{
		Name:                 common.StrSlice(d.Get("name")),
		Description:          common.StrSlice(d.Get("description")),
		ProvisioningStatus:   common.StrSlice(d.Get("provisioning_status")),
		OperatingStatus:      common.StrSlice(d.Get("operating_status")),
		VpcID:                common.StrSlice(d.Get("router_id")),
		VipSubnetCidrID:      common.StrSlice(d.Get("subnet_id")),
		L7FlavorID:           common.StrSlice(d.Get("l7_flavor")),
		L4FlavorID:           common.StrSlice(d.Get("l4_flavor")),
		L7ScaleFlavorID:      common.StrSlice(d.Get("l7_scale_flavor")),
		L4ScaleFlavorID:      common.StrSlice(d.Get("l4_scale_flavor")),
		VipAddress:           common.StrSlice(d.Get("vip_address")),
		VipPortID:            common.StrSlice(d.Get("vip_port_id")),
		IpV6VipAddress:       common.StrSlice(d.Get("ipv6_vip_address")),
		IpV6VipPortID:        common.StrSlice(d.Get("ipv6_vip_port_id")),
		IpV6VipSubnetID:      common.StrSlice(d.Get("ipv6_vip_subnet_id")),
		Eips:                 common.StrSlice(d.Get("eip_id")),
		PublicIps:            common.StrSlice(d.Get("public_ip_id")),
		AvailabilityZoneList: common.ExpandToStringSlice(d.Get("availability_zones").(*schema.Set).List()),
		EnterpriseProjectID:  common.StrSlice(d.Get("enterprise_project_id")),
		IPVersion:            intSlice(d.Get("ip_version")),
		ElbSubnetType:        common.StrSlice(d.Get("elb_subnet_type")),
		ProtectionStatus:     common.StrSlice(d.Get("protection_status")),
		MemberDeviceID:       common.StrSlice(d.Get("member_device_id")),
		MemberAddress:        common.StrSlice(d.Get("member_address")),
	}
	if common.IsAttrSet(d, "guaranteed") {
		value := d.Get("guaranteed").(bool)
		listOpts.Guaranteed = &value
	}
	if common.IsAttrSet(d, "deletion_protection") {
		value := d.Get("deletion_protection").(bool)
		listOpts.DeletionProtectionEnable = &value
	}

	lbList, err := loadbalancers.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error listing ELBv3 Load Balancer: %w", err)
	}

	if len(lbList) > 1 {
		return common.DataSourceTooManyDiag
	}
	if len(lbList) < 1 {
		return common.DataSourceTooFewDiag
	}

	lb := &lbList[0]
	return setLoadBalancerFields(d, meta, lb)
}

func intSlice(value interface{}) []int {
	if value == nil {
		return nil
	}
	v, ok := value.(int)
	if !ok || v == 0 {
		return nil
	}
	return []int{v}
}
