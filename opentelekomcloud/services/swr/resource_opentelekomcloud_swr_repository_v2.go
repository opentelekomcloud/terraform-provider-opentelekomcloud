package swr

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/repositories"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSwrRepositoryV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRepositoryCreate,
		ReadContext:   resourceRepositoryRead,
		UpdateContext: resourceRepositoryUpdate,
		DeleteContext: resourceRepositoryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRepositoryImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(
						regexp.MustCompile(`^[a-z0-9][a-z0-9\/._-]+[a-z0-9]+$`),
						"Only lowercase letters, digits, periods (.), underscores (_), slashes (/), and hyphens (-) are allowed.",
					),
					validation.StringDoesNotMatch(
						regexp.MustCompile(`_{3,}?|\.{2,}?|-{2,}?`),
						"Periods, underscores, and hyphens cannot be placed next to each other. A maximum of two consecutive underscores are allowed.",
					),
				),
			},
			"repository_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"app_server", "linux", "framework_app", "database", "lang", "other", "windows", "arm"},
					false,
				),
			},
			"is_public": {
				Type:     schema.TypeBool,
				Required: true,
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
	}
}

func resourceRepositoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := repositories.CreateOpts{
		Namespace:   d.Get("organization").(string),
		Repository:  d.Get("name").(string),
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
	}

	err = repositories.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating repository: %w", err)
	}
	d.SetId(opts.Repository)

	return resourceRepositoryRead(ctx, d, meta)
}

func resourceRepositoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	repo, err := repositories.Get(client, organization(d), repository(d.Id()))
	if err != nil {
		return fmterr.Errorf("error reading repository: %w", err)
	}

	mErr := multierror.Append(
		d.Set("name", repo.Name),
		d.Set("repository_id", repo.ID),
		d.Set("description", repo.Description),
		d.Set("category", repo.Category),
		d.Set("is_public", repo.IsPublic),
		d.Set("path", repo.Path),
		d.Set("internal_path", repo.InternalPath),
		d.Set("num_images", repo.NumImages),
		d.Set("size", repo.Size),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting resource fields: %w", err)
	}

	return nil
}

func resourceRepositoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := repositories.UpdateOpts{
		Namespace:   d.Get("organization").(string),
		Repository:  d.Get("repository").(string),
		Category:    d.Get("category").(string),
		Description: d.Get("description").(string),
		IsPublic:    d.Get("is_public").(bool),
	}
	err = repositories.Update(client, opts)
	if err != nil {
		return fmterr.Errorf("error updating repository: %w", err)
	}

	return resourceRepositoryRead(ctx, d, meta)
}

func resourceRepositoryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	err = repositories.Delete(client, organization(d), repository(d.Id()))
	if err != nil {
		fmterr.Errorf("error deleting repository: %w", err)
	}

	return nil
}

func resourceRepositoryImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for SWR repository import: format must be <organization>/<repository>")
		return nil, err
	}
	org := parts[0]
	repo := parts[1]
	d.SetId(repo)
	if err := d.Set("organization", org); err != nil {
		return nil, err
	}
	return schema.ImportStatePassthroughContext(ctx, d, meta)
}
