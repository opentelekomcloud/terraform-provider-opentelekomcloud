package er

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/association"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceAssociationsV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssociationsRead,

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
			"associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     associationsSchema(),
			},
		},
	}
}

func associationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
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

func buildAssociationsistOpts(d *schema.ResourceData) association.ListOpts {
	opts := association.ListOpts{
		RouterId:     d.Get("instance_id").(string),
		RouteTableId: d.Get("route_table_id").(string),
	}
	if attachmentId, ok := d.GetOk("attachment_id"); ok {
		opts.AttachmentId = []string{attachmentId.(string)}
	}

	if attachmentType, ok := d.GetOk("attachment_type"); ok {
		opts.ResourceType = []string{attachmentType.(string)}
	}

	if status, ok := d.GetOk("status"); ok {
		opts.State = []string{status.(string)}
	}

	return opts
}

func flattenAssociations(all []association.Association) []map[string]interface{} {
	if len(all) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(all))
	for i, assoc := range all {
		result[i] = map[string]interface{}{
			"id":              assoc.ID,
			"route_table_id":  assoc.RouteTableID,
			"attachment_id":   assoc.AttachmentID,
			"attachment_type": assoc.ResourceType,
			"resource_id":     assoc.ResourceID,
			"status":          assoc.State,
			"created_at":      assoc.CreatedAt,
			"updated_at":      assoc.UpdatedAt,
		}
	}
	return result
}

func dataSourceAssociationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ErV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	opts := buildAssociationsistOpts(d)
	resp, err := association.List(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving associations: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("associations", flattenAssociations(resp.Associations)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
