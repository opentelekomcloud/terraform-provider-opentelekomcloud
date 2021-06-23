package ims

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceImagesMemberV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagesMemberV2Create,
		ReadContext:   resourceImagesMemberV2Read,
		UpdateContext: resourceImagesMemberV2Update,
		DeleteContext: resourceImagesMemberV2Delete,

		Schema: map[string]*schema.Schema{
			"member": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"pending", "accepted", "rejected",
				}, false),
			},
			"vault_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceImagesMemberV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := members.CreateOpts{
		Member: d.Get("member").(string),
	}

	imageID := d.Get("image_id").(string)

	share, err := members.Create(client, imageID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error requesting share for private image: %w", err)
	}

	d.SetId(share.MemberID)

	return resourceImagesMemberV2Read(ctx, d, meta)
}

func resourceImagesMemberV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID := d.Get("image_id").(string)
	share, err := members.Get(client, imageID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "image"))
	}

	mErr := multierror.Append(
		d.Set("status", share.Status),
		d.Set("member", share.MemberID),
		d.Set("created_at", share.CreatedAt),
		d.Set("update_at", share.UpdatedAt),
		d.Set("image_id", share.ImageID),
		d.Set("schema", share.Schema),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceImagesMemberV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	updateOpts := members.UpdateOpts{
		Status:  d.Get("status").(string),
		VaultID: d.Get("vault_id").(string),
	}
	imageID := d.Get("image_id").(string)
	_, err = members.Update(client, imageID, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating share request: %w", err)
	}

	return resourceImagesMemberV2Read(ctx, d, meta)
}

func resourceImagesMemberV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID := d.Get("image_id").(string)
	if err := members.Delete(client, imageID, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting share request: %w", err)
	}

	d.SetId("")
	return nil
}
