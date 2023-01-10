package lts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/streams"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceLTSTopicV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicV2Create,
		ReadContext:   resourceTopicV2Read,
		DeleteContext: resourceTopicV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTopicV2Import,
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

			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceTopicV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	createOpts := streams.CreateOpts{
		LogStreamName: d.Get("topic_name").(string),
		GroupId:       groupId,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	topicCreate, err := streams.CreateLogStream(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating log topic: %s", err)
	}

	d.SetId(topicCreate)
	return resourceTopicV2Read(ctx, d, meta)
}

func resourceTopicV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	allTopics, err := streams.ListLogStream(client, groupId)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud log topic %s: %s", d.Id(), err)
	}

	var stream streams.LogStream
	for _, topic := range allTopics {
		if topic.LogStreamId == d.Id() {
			stream = topic
			break
		}
	}

	if stream.LogStreamId == "" {
		return fmterr.Errorf("OpenTelekomCloud log stream %s was not found", d.Id())
	}

	log.Printf("[DEBUG] Retrieved log topic %s: %#v", d.Id(), stream)
	d.SetId(stream.LogStreamId)

	mErr := multierror.Append(
		d.Set("topic_name", stream.LogStreamName),
		d.Set("creation_time", stream.CreationTime),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceTopicV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	err = streams.DeleteLogStream(client, streams.DeleteOpts{
		GroupId:  groupId,
		StreamId: d.Id(),
	})
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault400); ok {
			d.SetId("")
			return nil
		} else {
			return common.CheckDeletedDiag(d, err, "Error deleting log topic")
		}
	}

	d.SetId("")
	return nil
}

func resourceTopicV2Import(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for logtank topic. Format must be <group id>/<topic id>")
		return nil, err
	}

	config := meta.(*cfg.Config)
	client, err := config.LtsV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud LTS client: %s", err)
	}

	groupId := parts[0]
	topicId := parts[1]
	log.Printf("[DEBUG] Import log topic %s / %s", groupId, topicId)

	// check the parent logtank group whether exists.
	_, err = groups.ListLogGroups(client)
	if err != nil {
		return nil, fmt.Errorf("error importing OpenTelekomCloud log topic %s: %s", topicId, err)
	}

	d.SetId(topicId)

	if err := d.Set("group_id", groupId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
