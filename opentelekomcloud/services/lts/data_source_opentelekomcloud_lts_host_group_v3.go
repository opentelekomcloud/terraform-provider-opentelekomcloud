package lts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	hg "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/host-groups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceHostGroupV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostGroupV3Read,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ids": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"agent_access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceHostGroupV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.LtsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	requestResp, err := hg.List(client, hg.ListOpts{})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud LTS v3 host group")
	}

	var hostGroups []map[string]interface{}
	for _, group := range requestResp.Result {
		if d.Get("id").(string) != "" && group.ID != d.Get("id").(string) {
			continue
		}
		if d.Get("name").(string) != "" && group.Name != d.Get("name").(string) {
			continue
		}

		tagsMap := make(map[string]string)
		for _, tag := range group.Tags {
			tagsMap[tag.Key] = tag.Value
		}

		hostGroup := map[string]interface{}{
			"id":                group.ID,
			"name":              group.Name,
			"type":              group.Type,
			"host_ids":          group.HostIdList,
			"agent_access_type": group.AgentAccessType,
			"labels":            group.Labels,
			"tags":              tagsMap,
			"created_at":        common.FormatTimeStampRFC3339(group.CreatedAt/1000, false),
			"updated_at":        common.FormatTimeStampRFC3339(group.UpdatedAt/1000, false),
		}
		hostGroups = append(hostGroups, hostGroup)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("host_groups", hostGroups),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}
