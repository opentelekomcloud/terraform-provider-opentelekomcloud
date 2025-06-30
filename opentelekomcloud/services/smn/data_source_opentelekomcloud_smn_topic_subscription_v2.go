package smn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/subscriptions"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSmnTopicSubscriptionV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSmnTopicSubscriptionReadV2,

		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSmnTopicSubscriptionReadV2(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	protocol := d.Get("protocol").(string)
	endpoint := d.Get("endpoint").(string)

	subList, err := subscriptions.ListTopic(client, subscriptions.ListTopicOpts{
		TopicUrn: d.Get("topic_urn").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	if len(subList) < 1 {
		return fmterr.Errorf("your query returned no results." +
			" Please change your search criteria and try again")
	}

	var filteredSubs []subscriptions.Subscription

	for _, sub := range subList {
		if protocol != "" && sub.Protocol != protocol {
			continue
		}
		if endpoint != "" && sub.Endpoint != endpoint {
			continue
		}
		filteredSubs = append(filteredSubs, sub)
	}

	if len(filteredSubs) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	d.SetId(filteredSubs[0].SubscriptionUrn)

	mErr := multierror.Append(
		d.Set("topic_urn", filteredSubs[0].TopicUrn),
		d.Set("protocol", filteredSubs[0].Protocol),
		d.Set("status", filteredSubs[0].Status),
		d.Set("endpoint", filteredSubs[0].Endpoint),
		d.Set("owner", filteredSubs[0].Owner),
		d.Set("description", filteredSubs[0].Remark),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
