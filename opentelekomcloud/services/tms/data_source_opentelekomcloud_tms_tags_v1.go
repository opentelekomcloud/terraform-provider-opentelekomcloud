package tms

import (
	"context"
	"strings"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceTmsTagV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsTagV1Read,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTmsTagV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV1Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	allTags, err := tags.Get(client).Extract()
	if err != nil {
		return fmterr.Errorf("Error listing TMS predefined tags: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	var tagList []map[string]interface{}
	for _, t := range allTags.Tags {
		tag := map[string]interface{}{
			"key":   t.Key,
			"value": t.Value,
		}
		tagList = append(tagList, tag)
	}

	filter := map[string]interface{}{}
	if key, ok := d.GetOk("key"); ok {
		filter["key"] = key
	}
	if value, ok := d.GetOk("value"); ok {
		filter["value"] = value
	}
	filterResult, err := filterTags(tagList, filter)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("tags", filterResult); err != nil {
		return fmterr.Errorf("Error setting TMS tags: %s", err)
	}

	return nil
}

func filterTags(tags []map[string]interface{}, filter map[string]interface{}) ([]map[string]interface{}, error) {
	var filtered []map[string]interface{}

	if len(filter) == 0 {
		return tags, nil
	}

	for _, tag := range tags {
		match := true
		if val, ok := filter["key"]; ok {
			if !strings.Contains(strings.ToLower(tag["key"].(string)), strings.ToLower(val.(string))) {
				match = false
			}
		}
		if val, ok := filter["value"]; ok {
			if !strings.Contains(strings.ToLower(tag["value"].(string)), strings.ToLower(val.(string))) {
				match = false
			}
		}
		if match {
			filtered = append(filtered, tag)
		}
	}

	return filtered, nil
}
