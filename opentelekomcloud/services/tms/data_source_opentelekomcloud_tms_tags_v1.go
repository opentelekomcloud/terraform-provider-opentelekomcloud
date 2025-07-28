package tms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceTmsTagV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsTagV1Read,

		Schema: map[string]*schema.Schema{
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
	client, err := config.TmsV1Client()
	if err != nil {
		return fmterr.Errorf("Error creating Opentelekomcloud TMS client: %s", err)
	}

	allTags, err := tags.Get(client).Extract()
	if err != nil {
		return fmterr.Errorf("Error listing TMS predefined tags: %s", err)
	}

	var tagList []map[string]interface{}
	for _, t := range allTags.Tags {
		tag := map[string]interface{}{
			"key":   t.Key,
			"value": t.Value,
		}
		tagList = append(tagList, tag)
	}

	if err = d.Set("tags", tagList); err != nil {
		return fmterr.Errorf("Error setting TMS tags: %s", err)
	}

	return nil
}
