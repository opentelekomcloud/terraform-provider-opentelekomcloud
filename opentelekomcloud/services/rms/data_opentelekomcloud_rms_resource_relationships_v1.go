package rms

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rms/relations"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRmsResourceRelationshipsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRmsResourceRelationshipsRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"in", "out",
				}, false),
			},
			"relations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"relation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"from_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"to_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRmsResourceRelationshipsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, rmsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.RmsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationRMSV1Client, err)
	}

	listOpts := relations.ListAllOpts{
		DomainId:   GetRmsDomainId(client, config),
		ResourceId: d.Get("resource_id").(string),
		Direction:  d.Get("direction").(string),
	}

	relations, err := relations.ListRelations(client, listOpts)

	if err != nil {
		return diag.Errorf("error getting the relations list from server: %s", err)
	}

	var flattenRelations []map[string]interface{}

	for _, item := range relations {
		query := map[string]interface{}{
			"relation_type":      item.RelationType,
			"from_resource_type": item.FromResourceType,
			"to_resource_type":   item.ToResourceType,
			"from_resource_id":   item.FromResourceId,
			"to_resource_id":     item.ToResourceId,
		}
		flattenRelations = append(flattenRelations, query)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("relations", flattenRelations),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
