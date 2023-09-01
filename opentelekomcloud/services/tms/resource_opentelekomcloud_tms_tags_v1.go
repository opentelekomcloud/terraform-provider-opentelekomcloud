package tms

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func ResourceTmsTagV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmsTagV1Create,
		DeleteContext: resourceTmsTagV1Delete,
		ReadContext:   resourceTmsTagV1Read,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fffA-Za-z0-9-_]+$"),
									"The key can only consist of letters, digits, underscores (_) and hyphens (-)."),
								validation.StringLenBetween(1, 36),
							),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fffA-Za-z0-9-_.]+$"),
									"The key can only consist of letters, digits, periods (.)underscores (_) and hyphens (-)."),
								validation.StringLenBetween(1, 43),
							),
						},
					},
				},
			},
		},
	}
}

func resourceTmsTagV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TmsV1Client()
	if err != nil {
		return fmterr.Errorf("Error creating Opentelekomcloud TMS client: %s", err)
	}

	var tagIds []string
	var predefineTags []tags.Tag
	tagsRaw := d.Get("tags").([]interface{})
	for _, v := range tagsRaw {
		tag := v.(map[string]interface{})
		predefineTag := tags.Tag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
		predefineTags = append(predefineTags, predefineTag)
		tagId := fmt.Sprintf("%s:%s", tag["key"], tag["value"])
		tagIds = append(tagIds, tagId)
	}

	createOpts := &tags.BatchOpts{
		Tags:   predefineTags,
		Action: tags.ActionCreate,
	}

	log.Printf("[DEBUG] Create TMS tag options: %#v", createOpts)
	_, err = tags.BatchAction(client, "", createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("Error creating Opentelekomcloud TMS tags: %s", err)
	}

	d.SetId(hashcode.Strings(tagIds))
	return resourceTmsTagV1Read(ctx, d, meta)
}

func resourceTmsTagV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TmsV1Client()
	if err != nil {
		return fmterr.Errorf("Error creating Opentelekomcloud TMS client: %s", err)
	}

	allTags, err := tags.Get(client).Extract()
	if err != nil {
		return fmterr.Errorf("Error listing TMS predefined tags: %s", err)
	}

	// Check if the requested tag is missing on cloud side
	var tagList []map[string]interface{}
	tagsRaw := d.Get("tags").([]interface{})
	for _, v := range tagsRaw {
		tag := v.(map[string]interface{})
		key := tag["key"].(string)
		value := tag["value"].(string)

		for _, t := range allTags.Tags {
			if key == t.Key && value == t.Value {
				tagFound := map[string]interface{}{
					"key":   key,
					"value": value,
				}
				tagList = append(tagList, tagFound)
			}
		}
	}
	if err = d.Set("tags", tagList); err != nil {
		return fmterr.Errorf("Error setting TMS tags: %s", err)
	}

	return nil
}

func resourceTmsTagV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TmsV1Client()
	if err != nil {
		return fmterr.Errorf("Error creating Opentelekomcloud TMS client")
	}

	var predefineTags []tags.Tag
	tagsRaw := d.Get("tags").([]interface{})
	if len(tagsRaw) == 0 {
		log.Printf("[DEBUG] TMS tags are empty, no need to issue delete request")
		return nil
	}
	for _, v := range tagsRaw {
		tag := v.(map[string]interface{})
		predefineTag := tags.Tag{
			Key:   tag["key"].(string),
			Value: tag["value"].(string),
		}
		predefineTags = append(predefineTags, predefineTag)
	}

	deleteOpts := &tags.BatchOpts{
		Tags:   predefineTags,
		Action: tags.ActionDelete,
	}

	log.Printf("[DEBUG] Delete TMS tag options: %#v", deleteOpts)
	_, err = tags.BatchAction(client, "", deleteOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error deleting TMS tag: %s", err)
	}

	return nil
}
