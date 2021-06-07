package lts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/loggroups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/logtopics"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceLTSTopicV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicV2Create,
		ReadContext:   resourceTopicV2Read,
		DeleteContext: resourceTopicV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceTopicV2Import,
		},

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"topic_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"index_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceTopicV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	createOpts := &logtopics.CreateOpts{
		LogTopicName: d.Get("topic_name").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	topicCreate, err := logtopics.Create(client, groupId, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating log topic: %s", err)
	}

	d.SetId(topicCreate.ID)
	return resourceTopicV2Read(ctx, d, meta)
}

func resourceTopicV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	topic, err := logtopics.Get(client, groupId, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Error getting OpenTelekomCloud log topic %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved log topic %s: %#v", d.Id(), topic)
	if topic.ID != "" {
		d.SetId(topic.ID)
	}
	d.Set("topic_name", topic.Name)
	d.Set("index_enabled", topic.IndexEnabled)
	return nil
}

func resourceTopicV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	err = logtopics.Delete(client, groupId, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "Error deleting log topic"))
	}

	d.SetId("")
	return nil
}

func resourceTopicV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for logtank topic. Format must be <group id>/<topic id>")
		return nil, err
	}

	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := parts[0]
	topicId := parts[1]
	log.Printf("[DEBUG] Import log topic %s / %s", groupId, topicId)

	// check the parent logtank group whether exists.
	_, err = loggroups.Get(client, groupId).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error importing OpenTelekomCloud log topic %s: %s", topicId, err)
	}

	d.SetId(topicId)
	d.Set("group_id", groupId)

	return []*schema.ResourceData{d}, nil
}
