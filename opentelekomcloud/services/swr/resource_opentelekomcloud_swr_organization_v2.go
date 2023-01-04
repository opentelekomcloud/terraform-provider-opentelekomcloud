package swr

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/organizations"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSwrOrganizationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrganizationCreate,
		ReadContext:   resourceOrganizationRead,
		DeleteContext: resourceOrganizationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[a-z][a-z0-9._-]+[a-z0-9]+$`),
						"Only lowercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
					),
					validation.StringDoesNotMatch(
						regexp.MustCompile(`_{3,}?|\.{2,}?|-{2,}?`),
						"Periods, underscores, and hyphens cannot be placed next to each other. A maximum of two consecutive underscores are allowed.",
					),
				),
			},

			"organization_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceOrganizationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}
	name := d.Get("name").(string)
	opts := organizations.CreateOpts{
		Namespace: name,
	}
	err = organizations.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating SWR organization: %w", err)
	}
	d.SetId(name)

	return resourceOrganizationRead(ctx, d, meta)
}

func resourceOrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}
	org, err := organizations.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error reading SWR organization: %w", err)
	}
	mErr := multierror.Append(
		d.Set("name", org.Name),
		d.Set("organization_id", org.ID),
		d.Set("creator_name", org.CreatorName),
		d.Set("auth", org.Auth),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting SWR organization fields: %w", err)
	}
	return nil
}

func resourceOrganizationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}
	err = organizations.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting SWR organization: %w", err)
	}
	return nil
}
