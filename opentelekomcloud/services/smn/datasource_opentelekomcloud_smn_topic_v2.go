package smn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/topics"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceTopic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTopicRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"push_policy": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	name := d.Get("name").(string)

	allTopics, err := topics.List(client, topics.ListOpts{})
	if err != nil {
		return diag.FromErr(err)
	}

	var filterdTopics []topics.Topic
	for _, topic := range allTopics {
		if name != "" && topic.Name != name {
			continue
		}
		filterdTopics = append(filterdTopics, topic)
	}
	if len(filterdTopics) < 1 {
		return fmterr.Errorf("your query returned no results. Please use exact topic name and try again.")
	}

	desiredTopic := filterdTopics[0]
	d.SetId(desiredTopic.TopicUrn)

	log.Printf("[DEBUG] Retrieved topic %s: %#v", desiredTopic.TopicUrn, desiredTopic)

	mErr := multierror.Append(
		d.Set("topic_urn", desiredTopic.TopicUrn),
		d.Set("display_name", desiredTopic.DisplayName),
		d.Set("name", desiredTopic.Name),
		d.Set("push_policy", desiredTopic.PushPolicy),
		d.Set("update_time", desiredTopic.UpdateTime),
		d.Set("create_time", desiredTopic.CreateTime),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
