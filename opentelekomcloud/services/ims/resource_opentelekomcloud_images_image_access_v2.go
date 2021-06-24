package ims

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/members"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceImagesImageAccessV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagesImageAccessV2Create,
		ReadContext:   resourceImagesImageAccessV2Read,
		UpdateContext: resourceImagesImageAccessV2Update,
		DeleteContext: resourceImagesImageAccessV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"member_id": {
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

func resourceImagesImageAccessV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	memberID := d.Get("member_id").(string)
	imageID := d.Get("image_id").(string)

	createOpts := members.CreateOpts{
		Member: memberID,
	}
	_, err = members.Create(client, imageID, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error requesting share for private image: %w", err)
	}

	id := fmt.Sprintf("%s/%s", imageID, memberID)
	d.SetId(id)

	status := d.Get("status").(string)
	if status != "" {
		opts := members.UpdateOpts{
			Status: status,
		}

		_, err := members.Update(client, imageID, memberID, opts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating the image status: %w", err)
		}
	}

	return resourceImagesImageAccessV2Read(ctx, d, meta)
}

func resourceImagesImageAccessV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	member, err := members.Get(client, imageID, memberID).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "image_access_v2"))
	}

	mErr := multierror.Append(
		d.Set("status", member.Status),
		d.Set("member_id", member.MemberID),
		d.Set("created_at", member.CreatedAt),
		d.Set("update_at", member.UpdatedAt),
		d.Set("image_id", member.ImageID),
		d.Set("schema", member.Schema),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceImagesImageAccessV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts := members.UpdateOpts{
		Status: d.Get("status").(string),
	}

	_, err = members.Update(client, imageID, memberID, updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating share request: %w", err)
	}

	return resourceImagesImageAccessV2Read(ctx, d, meta)
}

func resourceImagesImageAccessV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID, memberID, err := resourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := members.Delete(client, imageID, memberID).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting share request: %w", err)
	}

	d.SetId("")
	return nil
}
