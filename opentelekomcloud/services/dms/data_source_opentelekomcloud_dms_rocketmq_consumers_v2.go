package dms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/groups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceDmsRocketMQConsumersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQConsumersV2Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_detail": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"online": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subscription_consistency": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"clients": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"language":    {Type: schema.TypeString, Computed: true},
						"version":     {Type: schema.TypeString, Computed: true},
						"client_id":   {Type: schema.TypeString, Computed: true},
						"client_addr": {Type: schema.TypeString, Computed: true},
						"subscriptions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"topic":      {Type: schema.TypeString, Computed: true},
									"type":       {Type: schema.TypeString, Computed: true},
									"expression": {Type: schema.TypeString, Computed: true},
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

func dataSourceDmsRocketMQConsumersV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.DmsV2Client(region)
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	group := d.Get("group").(string)
	opts := groups.ListClientsOpts{Limit: 50, IsDetail: d.Get("is_detail").(bool)}
	var resp *groups.ListClientsResponse
	var all []groups.Client
	for {
		resp, err = groups.ListClients(client, instanceID, group, opts)
		if err != nil {
			return diag.Errorf("error getting consumers of DMS RocketMQ consumer group (%s): %s", group, err)
		}
		all = append(all, resp.Clients...)
		if len(resp.Clients) == 0 || len(all) >= resp.Total {
			break
		}
		opts.Offset += len(resp.Clients)
	}

	ids := []string{instanceID, group}
	clients := make([]map[string]interface{}, len(all))
	for i, c := range all {
		subscriptions := make([]map[string]interface{}, len(c.Subscriptions))
		for j, s := range c.Subscriptions {
			subscriptions[j] = map[string]interface{}{
				"topic":      s.Topic,
				"type":       s.Type,
				"expression": s.Expression,
			}
		}
		ids = append(ids, c.ClientID)
		clients[i] = map[string]interface{}{
			"language":      c.Language,
			"version":       c.Version,
			"client_id":     c.ClientID,
			"client_addr":   c.ClientAddr,
			"subscriptions": subscriptions,
		}
	}
	log.Printf("[DEBUG] DMS RocketMQ consumers of group %s: %v", group, ids[2:])
	d.SetId(hashcode.Strings(ids))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("online", resp.Online),
		d.Set("subscription_consistency", resp.SubscriptionConsistency),
		d.Set("clients", clients),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
