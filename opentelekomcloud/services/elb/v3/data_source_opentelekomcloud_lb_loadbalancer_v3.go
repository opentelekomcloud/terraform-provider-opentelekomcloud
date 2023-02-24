package v3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/loadbalancers"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLoadBalancerV3() *schema.Resource {
	return &schema.Resource{
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
		},
	}
}

func dataSourceLoadBalancerV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	id := d.Get("id").(string)
	if id != "" {
		lb, err := loadbalancers.Get(client, id).Extract()
		if err != nil {
			return fmterr.Errorf("error getting ELBv3 Load Balancer: %w", err)
		}
		return setLoadBalancerFields(d, meta, lb)
	}

	listOpts := loadbalancers.ListOpts{
		Name:            common.StrSlice(d.Get("name")),
		VpcID:           common.StrSlice(d.Get("router_id")),
		VipSubnetCidrID: common.StrSlice(d.Get("subnet_id")),
		L7FlavorID:      common.StrSlice(d.Get("l7_flavor")),
		L4FlavorID:      common.StrSlice(d.Get("l4_flavor")),
		VipAddress:      common.StrSlice(d.Get("vip_address")),
		VipPortID:       common.StrSlice(d.Get("vip_port_id")),
	}

	pages, err := loadbalancers.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing ELBv3 Load Balancer: %w", err)
	}
	lbList, err := loadbalancers.ExtractLoadbalancers(pages)
	if err != nil {
		return fmterr.Errorf("error extracting ELBv3 Load Balancer: %w", err)
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
