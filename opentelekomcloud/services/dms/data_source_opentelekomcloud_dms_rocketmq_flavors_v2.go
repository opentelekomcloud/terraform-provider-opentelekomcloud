package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/products"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQFlavorsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQFlavorsV2Read,

		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"arch_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"billing_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arch_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"charging_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"qingtian_incompatible": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ios": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_spec_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"availability_zones": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"unavailability_zones": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"properties": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"engine_versions": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"flavor_alias": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"broker_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"core": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rcu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_tps_per_rcu": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_topic": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_topic_per_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_consumer": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_consumer_per_broker": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_storage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_storage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_storage_per_node": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"max_storage_per_node": {
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

func dataSourceDmsRocketMQFlavorsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	resp, err := products.List(client, products.ListOpts{
		Engine:    "reliability",
		ProductId: d.Get("flavor_id").(string),
	})
	if err != nil {
		return diag.Errorf("error getting DMS RocketMQ flavor list: %s", err)
	}
	log.Printf("[DEBUG] The response of DMS RocketMQ flavor list request is: %#v", resp)

	// There is no charging_mode in this schema, so filterFlavors doesn't filter by it.
	filtered, ids := filterFlavors(d, resp.Products)
	d.SetId(hashcode.Strings(ids))

	flavors := make([]map[string]interface{}, len(filtered))
	for i, v := range filtered {
		flavors[i] = map[string]interface{}{
			"id":                    v.ProductId,
			"type":                  v.Type,
			"vm_specification":      v.ECSFlavorId,
			"billing_code":          v.BillingCode,
			"arch_types":            v.ArchTypes,
			"charging_modes":        v.ChargingMode,
			"qingtian_incompatible": v.QingtianIncompatible,
			"ios":                   flattenIOs(v.IOS),
			"properties":            flattenRocketMQFlavorProperties(v.Properties),
		}
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("versions", resp.Versions),
		d.Set("flavors", flavors),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRocketMQFlavorProperties(p products.EngineProperties) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"engine_versions":         p.EngineVersions,
			"flavor_alias":            p.ProductAlias,
			"broker_num":              convertStringNumberIgnoreErr(p.BrokerNum),
			"min_broker":              convertStringNumberIgnoreErr(p.MinBroker),
			"max_broker":              convertStringNumberIgnoreErr(p.MaxBroker),
			"core":                    convertStringNumberIgnoreErr(p.Core),
			"rcu":                     convertStringNumberIgnoreErr(p.RCU),
			"max_tps_per_rcu":         convertStringNumberIgnoreErr(p.MaxTPSPerRCU),
			"max_topic":               convertStringNumberIgnoreErr(p.MaxTopic),
			"max_topic_per_broker":    convertStringNumberIgnoreErr(p.MaxTopicPerBroker),
			"max_consumer":            convertStringNumberIgnoreErr(p.MaxConsumer),
			"max_consumer_per_broker": convertStringNumberIgnoreErr(p.MaxConsumerPerBroker),
			"min_storage":             convertStringNumberIgnoreErr(p.MinStorage),
			"max_storage":             convertStringNumberIgnoreErr(p.MaxStorage),
			"min_storage_per_node":    convertStringNumberIgnoreErr(p.MinStoragePerNode),
			"max_storage_per_node":    convertStringNumberIgnoreErr(p.MaxStoragePerNode),
		},
	}
}
