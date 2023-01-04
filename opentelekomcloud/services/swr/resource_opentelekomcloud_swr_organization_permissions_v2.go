package swr

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/swr/v2/organizations"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSwrOrganizationPermissionsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrOrganizationPermissionsV2Create,
		ReadContext:   resourceSwrOrganizationPermissionsV2Read,
		UpdateContext: resourceSwrOrganizationPermissionsV2Update,
		DeleteContext: resourceSwrOrganizationPermissionsV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"auth": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceSwrOrganizationPermissionsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := organizations.Auth{
		UserID:   d.Get("user_id").(string),
		Username: d.Get("username").(string),
		Auth:     d.Get("auth").(int),
	}
	err = organizations.CreatePermissions(client, organization(d), []organizations.Auth{opts})
	if err != nil {
		return fmterr.Errorf("error creating organization permissions: %w", err)
	}
	d.SetId(opts.UserID)
	return resourceSwrOrganizationPermissionsV2Read(ctx, d, meta)
}

func resourceSwrOrganizationPermissionsV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	perms, err := organizations.GetPermissions(client, organization(d))
	if err != nil {
		return fmterr.Errorf("error getting organization permissions: %w", err)
	}
	var found *organizations.Auth
	for _, auth := range perms.OthersAuth {
		if auth.UserID == d.Id() {
			found = &auth
			break
		}
	}
	if found == nil {
		return diag.Errorf("no permissions for user %s are found", d.Id())
	}

	mErr := multierror.Append(
		d.Set("auth", found.Auth), // username/id shouldn't change, I suppose
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting permissions fields: %w", err)
	}

	return nil
}

func resourceSwrOrganizationPermissionsV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	opts := organizations.Auth{
		UserID:   d.Get("user_id").(string),
		Username: d.Get("username").(string),
		Auth:     d.Get("auth").(int),
	}
	err = organizations.UpdatePermissions(client, organization(d), []organizations.Auth{opts})
	if err != nil {
		return fmterr.Errorf("error updating organization permissions: %w", err)
	}

	return resourceSwrOrganizationPermissionsV2Read(ctx, d, meta)
}

func resourceSwrOrganizationPermissionsV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SwrV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ClientError, err)
	}

	err = organizations.DeletePermissions(client, organization(d), d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting organization permissions: %w", err)
	}

	d.SetId("")
	return nil
}
