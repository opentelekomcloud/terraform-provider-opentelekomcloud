package eps

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/eps/v1/projects"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceEnterpriseProjectsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseProjectsV1Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_projects": {
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
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
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
				},
			},
		},
	}
}

func flattenEnterpriseProjectResponseBody(project projects.EnterpriseProject) map[string]interface{} {
	result := map[string]interface{}{
		"id":          project.ID,
		"name":        project.Name,
		"type":        project.Type,
		"status":      project.Status,
		"description": project.Description,
		"created_at":  project.CreatedAt,
		"updated_at":  project.UpdatedAt,
	}

	return result
}

func dataSourceEnterpriseProjectsV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	opts := projects.ListOpts{
		Name:    d.Get("name").(string),
		ID:      d.Get("enterprise_project_id").(string),
		Status:  d.Get("status").(int),
		SortKey: "name",
		SortDir: "asc",
	}
	p, err := projects.List(client, opts)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving enterprise projects")
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result := make([]interface{}, 0, len(p.EnterpriseProjects))
	for _, project := range p.EnterpriseProjects {
		result = append(result, flattenEnterpriseProjectResponseBody(project))
	}

	mErr := multierror.Append(nil,
		d.Set("enterprise_projects", result),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
