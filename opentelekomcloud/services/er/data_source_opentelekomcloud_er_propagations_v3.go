package er

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/propagation"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourcePropagationsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePropagationsRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"propagations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     propagationSchema(),
			},
		},
	}
}

func propagationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attachment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
	}
}

func buildPropagationListOpts(d *schema.ResourceData) propagation.ListOpts {
	opts := propagation.ListOpts{}
	if attachmentId, ok := d.GetOk("attachment_id"); ok {
		opts.AttachmentId = []string{attachmentId.(string)}
	}

	if attachmentType, ok := d.GetOk("attachment_type"); ok {
		opts.ResourceType = []string{attachmentType.(string)}
	}

	if status, ok := d.GetOk("status"); ok {
		opts.State = []string{status.(string)}
	}

	opts.RouterId = d.Get("instance_id").(string)
	opts.RouteTableId = d.Get("route_table_id").(string)

	return opts
}

func flattenPropagations(all []propagation.Propagation) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, prop := range all {
		result[i] = map[string]interface{}{
			"id":              prop.ID,
			"instance_id":     prop.ErID,
			"route_table_id":  prop.RouteTableID,
			"attachment_id":   prop.AttachmentID,
			"attachment_type": prop.ResourceType,
			"resource_id":     prop.ResourceID,
			"status":          prop.State,
			"created_at":      prop.CreatedAt,
			"updated_at":      prop.UpdatedAt,
		}
	}
	return result
}

func dataSourcePropagationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	opts := buildPropagationListOpts(d)

	resp, err := propagation.List(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving propagations: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("propagations", flattenPropagations(resp.Propagations)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of ER propagations: %s", mErr)
	}
	return nil
}
