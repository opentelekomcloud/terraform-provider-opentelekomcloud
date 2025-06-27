package er

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/instance"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceErInstancesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceErInstancesV3Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asn": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
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
						"enable_default_propagation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable_default_association": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"auto_accept_shared_attachments": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"default_propagation_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_association_route_table_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zones": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildInstanceListOpts(d *schema.ResourceData) instance.ListOpts {
	return instance.ListOpts{
		State:   common.StringSliceIgnoreEmpty(d.Get("status").(string)),
		ID:      common.StringSliceIgnoreEmpty(d.Get("instance_id").(string)),
		SortKey: []string{"name"},
	}
}

func dataSourceErInstancesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	resp, err := instance.List(client, buildInstanceListOpts(d))
	if err != nil {
		return diag.Errorf("error retrieving OpenTelekomCloud ER v3 instances: %s", err)
	}
	instances, err := filterInstances(d, resp.Instances)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("instances", flattenInstances(instances)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving OpenTelekomCloud ER v3 instance data source fields: %s", mErr)
	}
	return nil
}

func flattenInstances(instances []instance.RouterInstance) []map[string]interface{} {
	if len(instances) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(instances))
	for i, item := range instances {
		result[i] = map[string]interface{}{
			"id":                                 item.ID,
			"asn":                                item.Asn,
			"name":                               item.Name,
			"description":                        item.Description,
			"status":                             item.State,
			"tags":                               common.TagsToMap(item.Tags),
			"created_at":                         common.FormatTimeStampRFC3339(common.ConvertTimeStrToNanoTimestamp(item.CreatedAt)/1000, false),
			"updated_at":                         common.FormatTimeStampRFC3339(common.ConvertTimeStrToNanoTimestamp(item.UpdatedAt)/1000, false),
			"enable_default_propagation":         item.EnableDefaultPropagation,
			"enable_default_association":         item.EnableDefaultAssociation,
			"auto_accept_shared_attachments":     item.AutoAcceptSharedAttachments,
			"default_propagation_route_table_id": item.DefaultPropagationRouteTableID,
			"default_association_route_table_id": item.DefaultAssociationRouteTableID,
			"availability_zones":                 item.AvailabilityZoneIDs,
		}
	}
	return result
}

func filterInstances(d *schema.ResourceData, instances []instance.RouterInstance) ([]instance.RouterInstance, error) {
	filter := map[string]interface{}{}
	if name, ok := d.GetOk("name"); ok {
		filter["Name"] = name
	}
	filterResult, err := common.FilterSliceWithField(instances, filter)
	if err != nil {
		return nil, fmt.Errorf("error filting instance list: %s", err)
	}

	result := make([]instance.RouterInstance, 0, len(filterResult))
	tagFilter := d.Get("tags").(map[string]interface{})
	for _, val := range filterResult {
		item := val.(instance.RouterInstance)
		tagmap := common.TagsToMap(item.Tags)
		if !common.MapContains(tagmap, tagFilter) {
			continue
		}
		result = append(result, item)
	}
	return result, nil
}
