package iam

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityUserGroupMembershipV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityUserGroupMembershipV3Create,
		ReadContext:   resourceIdentityUserGroupMembershipV3Read,
		UpdateContext: resourceIdentityUserGroupMembershipV3Update,
		DeleteContext: resourceIdentityUserGroupMembershipV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"user": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"groups": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

const iamClientKey = "client-iam"

func resourceIdentityUserGroupMembershipV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	userID := d.Get("user").(string)
	d.SetId(userID)
	groupSet := d.Get("groups").(*schema.Set)
	mErr := &multierror.Error{}
	for _, group := range groupSet.List() {
		mErr = multierror.Append(mErr,
			users.AddToGroup(client, group.(string), userID).ExtractErr(),
		)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error assigning user to the group")
	}

	clientCtx := common.CtxWithClient(ctx, client, iamClientKey)
	return resourceIdentityUserGroupMembershipV3Read(clientCtx, d, meta)
}

func resourceIdentityUserGroupMembershipV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, iamClientKey, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	pages, err := users.ListGroups(client, d.Id()).AllPages()
	if err != nil {
		return fmterr.Errorf("error listing groups: %w", err)
	}
	gps, err := groups.ExtractGroups(pages)
	if err != nil {
		return fmterr.Errorf("error extracting group list: %w", err)
	}

	groupIDs := make([]string, len(gps))
	for i, g := range gps {
		groupIDs[i] = g.ID
	}
	if err := d.Set("groups", groupIDs); err != nil {
		return fmterr.Errorf("error setting group IDs: %w", err)
	}

	return nil
}

func resourceIdentityUserGroupMembershipV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	if !d.HasChange("groups") {
		return nil
	}

	oldG, newG := d.GetChange("groups")
	oldSet := oldG.(*schema.Set)
	newSet := newG.(*schema.Set)

	toAdd := newSet.Difference(oldSet)
	if err := addGroupsToUser(client, d.Id(), toAdd.List()); err != nil {
		return fmterr.Errorf("error adding new groups to the user: %w", err)
	}

	toRemove := oldSet.Difference(newSet)
	if err := removeGroupsFromUser(client, d.Id(), toRemove.List()); err != nil {
		return fmterr.Errorf("error removing groups from the user: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, iamClientKey)
	return resourceIdentityUserGroupMembershipV3Read(clientCtx, d, meta)
}

func resourceIdentityUserGroupMembershipV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	groupIDs := d.Get("groups").(*schema.Set)
	if err := removeGroupsFromUser(client, d.Id(), groupIDs.List()); err != nil {
		return fmterr.Errorf("error removing groups from the user: %w", err)
	}

	return nil
}

func addGroupsToUser(client *golangsdk.ServiceClient, userID string, groupIDs []interface{}) error {
	mErr := &multierror.Error{}
	for _, group := range groupIDs {
		mErr = multierror.Append(mErr,
			users.AddToGroup(client, group.(string), userID).ExtractErr(),
		)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}
	return nil
}

func removeGroupsFromUser(client *golangsdk.ServiceClient, userID string, groupIDs []interface{}) error {
	mErr := &multierror.Error{}
	for _, group := range groupIDs {
		err := users.RemoveFromGroup(client, group.(string), userID).ExtractErr()
		if _, ok := err.(golangsdk.ErrDefault404); ok { // group was already removed
			continue
		}
		mErr = multierror.Append(mErr, err)
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}
	return nil
}
