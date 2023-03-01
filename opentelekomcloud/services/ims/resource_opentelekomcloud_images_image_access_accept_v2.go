package ims

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/image/v2/members"
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
	status := d.Get("status").(string)
	pendingStatus := "pending"
	shareState := &resource.StateChangeConf{
		Target:  []string{pendingStatus},
		Refresh: waitForImageRequestStatus(client, imageID, memberID, pendingStatus),
		Timeout: 1 * time.Minute,
	}
	_, err = shareState.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for `%s` status: %w", pendingStatus, err)
	}

	_, err = members.Update(client, members.UpdateOpts{
		ImageId:  imageID,
		MemberId: memberID,
		Status:   status,
	})
	if err != nil {
		return fmterr.Errorf("error setting a member status to the image share: %w", err)
	}
	state := &resource.StateChangeConf{
		Target:  []string{status},
		Refresh: waitForImageRequestStatus(client, imageID, memberID, status),
		Timeout: 1 * time.Minute,
	}
	_, err = state.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for `%s` status: %w", status, err)
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

	member, err := members.Get(client, members.MemberOpts{
		ImageId:  imageID,
		MemberId: memberID,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "image_access_accept_v2")
	}

	mErr := multierror.Append(
		d.Set("status", member.Status),
		d.Set("member_id", member.MemberId),
		d.Set("created_at", member.CreatedAt),
		d.Set("image_id", member.ImageId),
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

	status := d.Get("status").(string)
	opts := members.UpdateOpts{
		ImageId:  imageID,
		MemberId: memberID,
		Status:   status,
	}
	_, err = members.Update(client, opts)
	if err != nil {
		return fmterr.Errorf("error updating the image with the member: %w", err)
	}
	state := &resource.StateChangeConf{
		Target:  []string{status},
		Refresh: waitForImageRequestStatus(client, imageID, memberID, status),
		Timeout: 1 * time.Minute,
	}
	_, err = state.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for `%s` status: %w", status, err)
	}

	return resourceImagesImageAccessAcceptV2Read(ctx, d, meta)
}

func resourceImagesImageAccessAcceptV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	status := "rejected"
	opts := members.UpdateOpts{
		ImageId:  imageID,
		MemberId: memberID,
		Status:   status,
	}
	if _, err := members.Update(client, opts); err != nil {
		return common.CheckDeletedDiag(d, err, "image_access_accept_v2")
	}
	state := &resource.StateChangeConf{
		Target:  []string{status},
		Refresh: waitForImageRequestStatus(client, imageID, memberID, status),
		Timeout: 1 * time.Minute,
	}
	_, err = state.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for `%s` status: %w", status, err)
	}

	return nil
}
