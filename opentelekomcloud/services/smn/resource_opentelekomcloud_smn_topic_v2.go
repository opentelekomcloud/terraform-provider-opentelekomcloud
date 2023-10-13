package smn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/topics"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicCreate,
		ReadContext:   resourceTopicRead,
		UpdateContext: resourceTopicUpdate,
		DeleteContext: resourceTopicDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
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

func resourceTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud smn client: %s", err)
	}

	createOpts := topics.CreateOps{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	topic, err := topics.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error getting topic from result: %s", err)
	}
	log.Printf("[DEBUG] Create : topic.TopicUrn %s", topic.TopicUrn)

	if common.HasFilledOpt(d, "tags") {
		tagClient, err := config.SmnV2TagClient(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud smn tags client: %s", err)
		}
		tagClient.MoreHeaders = map[string]string{
			"X-SMN-RESOURCEID-TYPE": "name",
		}
		tagRaw := d.Get("tags").(map[string]interface{})
		if len(tagRaw) > 0 {
			tagList := common.ExpandResourceTags(tagRaw)
			if err := tags.Create(tagClient, "smn_topic", d.Get("name").(string), tagList).ExtractErr(); err != nil {
				return fmterr.Errorf("error setting tags of SMN topic: %w", err)
			}
		}
	}

	if topic.TopicUrn != "" {
		d.SetId(topic.TopicUrn)
		return resourceTopicRead(ctx, d, meta)
	}

	return fmterr.Errorf("unexpected conversion error in resourceTopicCreate.")
}

func resourceTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud smn client: %s", err)
	}

	topicUrn := d.Id()
	topic, err := topics.Get(client, topicUrn).ExtractGet()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "topic")
	}

	log.Printf("[DEBUG] Retrieved topic %s: %#v", topicUrn, topic)

	mErr := multierror.Append(
		d.Set("topic_urn", topic.TopicUrn),
		d.Set("display_name", topic.DisplayName),
		d.Set("name", topic.Name),
		d.Set("push_policy", topic.PushPolicy),
		d.Set("update_time", topic.UpdateTime),
		d.Set("create_time", topic.CreateTime),
	)

	// read tags
	tagClient, err := config.SmnV2TagClient(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud smn tags client: %s", err)
	}
	if err != nil {
		return diag.Errorf("error creating SMN tag client: %s", err)
	}
	tagClient.MoreHeaders = map[string]string{
		"X-SMN-RESOURCEID-TYPE": "name",
	}
	if resourceTags, err := tags.Get(tagClient, "smn_topic", d.Get("name").(string)).Extract(); err == nil {
		tagMap := common.TagsToMap(resourceTags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud SMN topic: %s", err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceTopicDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Deleting topic %s", d.Id())

	id := d.Id()
	result := topics.Delete(client, id)
	if result.Err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Successfully deleted topic %s", id)
	return nil
}

func resourceTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Updating topic %s", d.Id())
	id := d.Id()

	var updateOpts topics.UpdateOps
	if d.HasChange("display_name") {
		updateOpts.DisplayName = d.Get("display_name").(string)
	}
	if d.HasChange("tags") {
		tagClient, err := config.SmnV2TagClient(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud smn tags client: %s", err)
		}
		tagClient.MoreHeaders = map[string]string{
			"X-SMN-RESOURCEID-TYPE": "name",
		}
		if err := common.UpdateResourceTags(tagClient, d, "smn_topic", d.Get("name").(string)); err != nil {
			return fmterr.Errorf("error updating tags of SMN topic %s: %s", d.Id(), err)
		}
	}
	topic, err := topics.Update(client, updateOpts, id).Extract()
	if err != nil {
		return fmterr.Errorf("error updating topic from result: %s", err)
	}

	log.Printf("[DEBUG] Update : topic.TopicUrn: %s", topic.TopicUrn)
	if topic.TopicUrn != "" {
		d.SetId(topic.TopicUrn)
		return resourceTopicRead(ctx, d, meta)
	}
	return nil
}
