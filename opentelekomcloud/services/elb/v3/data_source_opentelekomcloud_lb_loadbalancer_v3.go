package v3

import (
	"context"

	"github.com/hashicorp/go-multierror"
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"max_connections": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"qps": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeInt,
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
		return setLoadBalancerFields(d, lb)
	}

	listOpts := loadbalancers.ListOpts{}
	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = []string{v.(string)}
	}
	if v, ok := d.GetOk("router_id"); ok {
		listOpts.VpcID = []string{v.(string)}
	}
	if v, ok := d.GetOk("subnet_id"); ok {
		listOpts.VipSubnetCidrID = []string{v.(string)}
	}
	if v, ok := d.GetOk("l7_flavor_id"); ok {
		listOpts.L7FlavorID = []string{v.(string)}
	}
	if v, ok := d.GetOk("l4_flavor_id"); ok {
		listOpts.L4FlavorID = []string{v.(string)}
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
	return setLoadBalancerFields(d, lb)
}

func setLoadBalancerFields(d *schema.ResourceData, lb *loadbalancers.LoadBalancer) diag.Diagnostics {
	d.SetId(lb.ID)
	publicIpInfo := make([]map[string]interface{}, len(lb.PublicIps))
	if len(lb.PublicIps) > 0 {
		info := d.Get("public_ip.0").(map[string]interface{})
		info["id"] = lb.PublicIps[0].PublicIpID
		info["address"] = lb.PublicIps[0].PublicIpAddress
		publicIpInfo[0] = info
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
