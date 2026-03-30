package swr

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/repositories"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceSwrRepositoryV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRepositoryRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"repositories": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repository_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"num_images": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRepositoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := repositories.ListOpts{
		Name: repository(d.Get("name").(string)),
	}

	pages, err := repositories.List(client, opts).AllPages()
	if err != nil {
		return fmterr.Errorf("error fetching repositories: %w", err)
	}
	repos, err := repositories.ExtractRepositories(pages)
	if err != nil {
		return fmterr.Errorf("error extracting repositories: %w", err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("repositories", setRepositories(repos)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setRepositories(reposInResp []repositories.ImageRepository) []map[string]interface{} {
	var repos []map[string]interface{}
	for _, repoInResp := range reposInResp {
		repo := map[string]interface{}{
			"name":          repoInResp.Name,
			"repository_id": repoInResp.ID,
			"description":   repoInResp.Description,
			"category":      repoInResp.Category,
			"is_public":     repoInResp.IsPublic,
			"path":          repoInResp.Path,
			"internal_path": repoInResp.InternalPath,
			"num_images":    repoInResp.NumImages,
			"size":          repoInResp.Size,
		}
		repos = append(repos, repo)
	}
	return repos
}
