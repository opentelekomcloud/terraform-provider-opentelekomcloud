package v3

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/pools"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceLBPoolV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBPoolV3Read,

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
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lb_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP", "QUIC_CID"},
					false,
				),
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"TCP", "UDP", "HTTP", "QUIC", "HTTPS"},
					true,
				),
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"session_persistence": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cookie_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"persistence_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"slow_start": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"healthmonitor_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"loadbalancer_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"member_deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"member_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"member_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"nonProtection", "consoleProtection",
				}, false),
			},
			"protection_reason": {
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

func buildLBPoolV3ListOpts(d *schema.ResourceData, memberDeletionProtection *bool) pools.ListOpts {
	return pools.ListOpts{
		Description:                    common.StrSlice(d.Get("description")),
		HealthMonitorID:                common.StrSlice(d.Get("healthmonitor_id")),
		LBMethod:                       common.StrSlice(d.Get("lb_algorithm")),
		Protocol:                       common.StrSlice(d.Get("protocol")),
		Name:                           common.StrSlice(d.Get("name")),
		LoadbalancerID:                 common.StrSlice(d.Get("loadbalancer_id")),
		IPVersion:                      common.StrSlice(d.Get("ip_version")),
		MemberAddress:                  common.StrSlice(d.Get("member_address")),
		MemberDeviceID:                 common.StrSlice(d.Get("member_device_id")),
		ListenerID:                     common.StrSlice(d.Get("listener_id")),
		MemberInstanceID:               common.StrSlice(d.Get("member_instance_id")),
		VpcID:                          common.StrSlice(d.Get("vpc_id")),
		Type:                           common.StrSlice(d.Get("type")),
		ProtectionStatus:               common.StrSlice(d.Get("protection_status")),
		MemberDeletionProtectionEnable: memberDeletionProtection,
	}
}

func dataSourceLBPoolV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if id := d.Get("id").(string); id != "" {
		pool, err := pools.Get(client, id)
		if err != nil {
			return fmterr.Errorf("error getting ELB v3 pool: %w", err)
		}
		return setLBPoolV3Fields(d, pool)
	}

	var memberDeletionProtection *bool
	if common.IsAttrSet(d, "member_deletion_protection") {
		value := d.Get("member_deletion_protection").(bool)
		memberDeletionProtection = &value
	}
	poolList, err := pools.List(client, buildLBPoolV3ListOpts(d, memberDeletionProtection))
	if err != nil {
		return fmterr.Errorf("error listing ELB v3 pools: %w", err)
	}
	if len(poolList) < 1 {
		return common.DataSourceTooFewDiag
	}
	if len(poolList) > 1 {
		return common.DataSourceTooManyDiag
	}
	return setLBPoolV3Fields(d, &poolList[0])
}
