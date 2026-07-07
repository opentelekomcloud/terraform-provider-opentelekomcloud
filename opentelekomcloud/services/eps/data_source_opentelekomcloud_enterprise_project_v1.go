package eps

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/eps/v1/projects"
)

func DataSourceEnterpriseProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEnterpriseProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
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
	}
}

func dataSourceEnterpriseProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	listOpts := projects.ListOpts{
		Name:   d.Get("name").(string),
		ID:     d.Get("id").(string),
		Status: d.Get("status").(int),
	}
	p, err := projects.List(client, listOpts)

	if err != nil {
		return diag.Errorf("error retrieving enterprise projects: %s", err)
	}

	project, err := flattenEnterpriseProject(d, p.EnterpriseProjects)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(project.ID)
	mErr := multierror.Append(nil,
		d.Set("name", project.Name),
		d.Set("description", project.Description),
		d.Set("status", project.Status),
		d.Set("created_at", project.CreatedAt),
		d.Set("updated_at", project.UpdatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting enterprise project fields: %s", err)
	}

	return nil
}

func flattenEnterpriseProject(d *schema.ResourceData, projects []projects.EnterpriseProject) (
	*projects.EnterpriseProject, error) {
	if len(projects) < 1 {
		return nil, fmt.Errorf("your query returned no results." +
			" Please change your search criteria and try again")
	}

	if len(projects) > 1 {
		name := d.Get("name").(string)
		// There is no condition to find the target enterprise project.
		if name == "" {
			return nil, fmt.Errorf("your query returned more than one result." +
				" Please specify your enterprise project name or enterprise project id")
		}

		for _, v := range projects {
			if v.Name == name {
				// Use name to find the target enterprise project.
				return &v, nil
			}
		}
		return nil, fmt.Errorf("cannot find the target enterprise project through your input name." +
			" Please change your search criteria and try again")
	}
	return &projects[0], nil
}
