package taurusdb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/proxy"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func DataSourceTaurusDBV3MysqlProxies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBMysqlProxiesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proxy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_threshold_in_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"elb_vip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transaction_split": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     nodesElem(),
						},
						"master_node_weight": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"readonly_nodes_weight": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func nodesElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen_flag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBMysqlProxiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	instanceId := d.Get("instance_id").(string)

	proxies, err := proxy.List(client, instanceId, proxy.ListOpts{})
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB MySQL proxies: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	proxyList := make([]map[string]interface{}, len(proxies))
	for i, p := range proxies {
		proxyList[i] = map[string]interface{}{
			"id":                         p.Proxy.PoolId,
			"name":                       p.Proxy.Name,
			"flavor":                     p.Proxy.FlavorRef,
			"port":                       p.Proxy.Port,
			"status":                     p.Proxy.Status,
			"delay_threshold_in_seconds": p.Proxy.DelayThresholdInSeconds,
			"node_num":                   p.Proxy.NodeNum,
			"ram":                        p.Proxy.Ram,
			"mode":                       p.Proxy.Mode,
			"elb_vip":                    p.Proxy.ElbVip,
			"vcpus":                      p.Proxy.Vcpus,
			"transaction_split":          convertTransactionSplit(p.Proxy.TransactionSplit),
			"address":                    p.Proxy.Address,
			"nodes":                      flattenProxyNodes(p.Proxy.Nodes),
			"master_node_weight":         flattenMasterNode(p.MasterNode),
			"readonly_nodes_weight":      flattenReadonlyNodes(p.ReadonlyNodes),
		}
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("proxy_list", proxyList),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting TaurusDB MySQL proxies fields: %s", err)
	}

	return nil
}

func convertTransactionSplit(value string) string {
	if value == "true" {
		return "ON"
	}
	return "OFF"
}

func flattenMasterNode(node proxy.MysqlProxyNodeV3) []map[string]interface{} {
	if node.Id == "" {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":     node.Id,
			"name":   node.Name,
			"weight": node.Weight,
		},
	}
}

func flattenReadonlyNodes(nodes []proxy.MysqlProxyNodeV3) []map[string]interface{} {
	if len(nodes) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(nodes))
	for i, node := range nodes {
		result[i] = map[string]interface{}{
			"id":     node.Id,
			"name":   node.Name,
			"weight": node.Weight,
		}
	}
	return result
}
