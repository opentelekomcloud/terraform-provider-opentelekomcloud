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

func ResourceImagesImageAccessAcceptV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagesImageAccessAcceptV2Create,
		ReadContext:   resourceImagesImageAccessAcceptV2Read,
		UpdateContext: resourceImagesImageAccessAcceptV2Update,
		DeleteContext: resourceImagesImageAccessAcceptV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"accepted", "rejected", "pending",
				}, false),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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

func resourceImagesImageAccessAcceptV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID := d.Get("image_id").(string)
	memberID := d.Get("member_id").(string)

	// accept status on the consumer side
	opts := members.UpdateOpts{
		Status: d.Get("status").(string),
	}
	_, err = members.Update(client, imageID, memberID, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error setting a member status to the image share: %w", err)
	}

	id := fmt.Sprintf("%s/%s", imageID, memberID)
	d.SetId(id)

	return resourceImagesImageAccessAcceptV2Read(ctx, d, meta)
}

func resourceImagesImageAccessAcceptV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID, memberID, err := ResourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	member, err := members.Get(client, imageID, memberID).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "image_access_accept_v2"))
	}

	mErr := multierror.Append(
		d.Set("status", member.Status),
		d.Set("member_id", member.MemberID),
		d.Set("created_at", member.CreatedAt),
		d.Set("image_id", member.ImageID),
		d.Set("schema", member.Schema),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceImagesImageAccessAcceptV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID, memberID, err := ResourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	opts := members.UpdateOpts{
		Status: d.Get("status").(string),
	}
	_, err = members.Update(client, imageID, memberID, opts).Extract()
	if err != nil {
		return fmterr.Errorf("Error updating the image with the member: %w", err)
	}

	return resourceImagesImageAccessAcceptV2Read(ctx, d, meta)
}

func resourceImagesImageAccessAcceptV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	imageID, memberID, err := ResourceImagesImageAccessV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// reject status on the consumer side
	opts := members.UpdateOpts{
		Status: "rejected",
	}
	if err := members.Update(client, imageID, memberID, opts).Err; err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "image_access_accept_v2"))
	}

	return nil
}
