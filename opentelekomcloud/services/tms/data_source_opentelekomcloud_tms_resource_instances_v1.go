package tms

import (
	"context"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	rt "github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/resource-tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceTmsResourceInstancesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmsResourceInstancesV1Read,

		Schema: map[string]*schema.Schema{
			"resource_types": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"without_any_tag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTmsResourceInstancesV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, tmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.TmsV1Client()
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	listOpts := rt.ListResourceOpts{
		ResourceTypes: common.ExpandToStringList(d.Get("resource_types").([]interface{})),
		ProjectId:     d.Get("project_id").(string),
		Tags:          getTMSTags(d),
		WithoutAnyTag: pointerto.Bool(d.Get("without_any_tag").(bool)),
	}

	allResourceInstances, err := rt.ListResources(client, listOpts)
	if err != nil {
		return fmterr.Errorf("Error listing Resources: %s", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	var resourceList []map[string]interface{}
	for _, t := range allResourceInstances {
		resource := map[string]interface{}{
			"project_id":    t.ProjectId,
			"project_name":  t.ProjectName,
			"resource_id":   t.ResourceId,
			"resource_name": t.ResourceName,
			"resource_type": t.ResourceType,
			"tags":          setTMSResourceTags(t.Tags),
		}
		resourceList = append(resourceList, resource)
	}

	if err = d.Set("resources", resourceList); err != nil {
		return fmterr.Errorf("Error setting TMS tagged resources: %s", err)
	}
	return nil
}

func getTMSTags(d *schema.ResourceData) []rt.ListResourceTag {
	tagsRaw := d.Get("tags").([]interface{})

	var tags []rt.ListResourceTag
	for _, v := range tagsRaw {
		tagRaw := v.(map[string]interface{})
		tags = append(tags, rt.ListResourceTag{
			Key:    tagRaw["key"].(string),
			Values: common.ExpandToStringList(tagRaw["values"].([]interface{})),
		})
	}
	return tags
}
func setTMSResourceTags(tags []rt.ResourceTag) map[string]string {
	result := make(map[string]string)
	for _, val := range tags {
		result[val.Key] = common.StringPtrToString(val.Value)
	}

	return result
}
