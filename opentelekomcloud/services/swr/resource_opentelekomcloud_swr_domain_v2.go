package swr

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/domains"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSwrDomainV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrDomainCreate,
		ReadContext:   resourceSwrDomainRead,
		UpdateContext: resourceSwrDomainUpdate,
		DeleteContext: resourceSwrDomainDelete,

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"read",
				}, false),
			},
			"deadline": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSwrDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := domains.CreateOpts{
		Namespace:    d.Get("organization").(string),
		Repository:   repository(d.Get("repository").(string)),
		AccessDomain: d.Get("access_domain").(string),
		Permit:       d.Get("permission").(string),
		Deadline:     d.Get("deadline").(string),
		Description:  d.Get("description").(string),
	}

	err = domains.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating domain: %w", err)
	}
	d.SetId(opts.AccessDomain)

	return resourceSwrDomainRead(ctx, d, meta)
}

func resourceSwrDomainRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := domains.GetOpts{
		Namespace:    d.Get("organization").(string),
		Repository:   repository(d.Get("repository").(string)),
		AccessDomain: d.Get("access_domain").(string),
	}
	domain, err := domains.Get(client, opts)
	if err != nil {
		return fmterr.Errorf("error reading domain: %w", err)
	}

	mErr := multierror.Append(
		d.Set("access_domain", strings.ToUpper(domain.AccessDomain)),
		d.Set("repository", domain.Repository),
		d.Set("organization", domain.Namespace),
		d.Set("description", domain.Description),
		d.Set("status", domain.Status),
		d.Set("permission", domain.Permit),
		d.Set("created", domain.Created),
		d.Set("updated", domain.Updated),
		d.Set("creator_id", domain.CreatorID),
		d.Set("creator_name", domain.CreatorName),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting resource fields: %w", err)
	}

	return nil
}

func resourceSwrDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := domains.UpdateOpts{
		Namespace:   d.Get("organization").(string),
		Repository:  repository(d.Get("repository").(string)),
		Permit:      d.Get("permission").(string),
		Deadline:    d.Get("deadline").(string),
		Description: d.Get("description").(string),
	}

	err = domains.Update(client, opts)
	if err != nil {
		return fmterr.Errorf("error updating domain: %w", err)
	}

	return resourceSwrDomainRead(ctx, d, meta)
}

func resourceSwrDomainDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := domains.GetOpts{
		Namespace:    d.Get("organization").(string),
		Repository:   repository(d.Get("repository").(string)),
		AccessDomain: d.Get("access_domain").(string),
	}
	err = domains.Delete(client, opts)
	if err != nil {
		fmterr.Errorf("error deleting domain: %w", err)
	}

	return nil
}
