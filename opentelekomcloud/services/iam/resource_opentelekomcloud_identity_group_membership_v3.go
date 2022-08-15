package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIdentityGroupMembershipV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityGroupMembershipV3Create,
		ReadContext:   resourceIdentityGroupMembershipV3Read,
		UpdateContext: resourceIdentityGroupMembershipV3Update,
		DeleteContext: resourceIdentityGroupMembershipV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"users": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIdentityGroupMembershipV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	group := d.Get("group").(string)
	userList := common.ExpandToStringSlice(d.Get("users").(*schema.Set).List())

	if err := addUsersToGroup(identityClient, group, userList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())

	return resourceIdentityGroupMembershipV3Read(ctx, d, meta)
}

func resourceIdentityGroupMembershipV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}
	group := d.Get("group").(string)
	userList := d.Get("users").(*schema.Set)
	var ul []string

	allPages, err := users.ListInGroup(identityClient, group, users.ListOpts{}).AllPages()
	if err != nil {
		if _, b := err.(golangsdk.ErrDefault404); b {
			d.SetId("")
			return nil
		}
		return fmterr.Errorf("unable to query groups: %s", err)
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve users: %s", err)
	}

	for _, u := range allUsers {
		if userList.Contains(u.ID) {
			ul = append(ul, u.ID)
		}
	}

	if err := d.Set("users", ul); err != nil {
		return fmterr.Errorf("error setting user list from IAM (%s), error: %s", group, err)
	}

	return nil
}

func resourceIdentityGroupMembershipV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	if d.HasChange("users") {
		group := d.Get("group").(string)

		o, n := d.GetChange("users")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}

		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := common.ExpandToStringSlice(os.Difference(ns).List())
		add := common.ExpandToStringSlice(ns.Difference(os).List())

		if err := removeUsersFromGroup(identityClient, group, remove); err != nil && common.IsResourceNotFound(err) {
			return fmterr.Errorf("error update user-group-membership: %s", err)
		}

		if err := addUsersToGroup(identityClient, group, add); err != nil {
			return fmterr.Errorf("error update user-group-membership: %s", err)
		}
	}

	return resourceIdentityGroupMembershipV3Read(ctx, d, meta)
}

func resourceIdentityGroupMembershipV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	identityClient, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.IdentityV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(clientCreationFail, err)
	}

	group := d.Get("group").(string)
	userSlice := common.ExpandToStringSlice(d.Get("users").(*schema.Set).List())

	if err := removeUsersFromGroup(identityClient, group, userSlice); err != nil {
		return fmterr.Errorf("error delete user-group-membership: %s", err)
	}

	d.SetId("")
	return nil
}

func addUsersToGroup(identityClient *golangsdk.ServiceClient, group string, userList []string) error {
	for _, u := range userList {
		if r := users.AddToGroup(identityClient, group, u).ExtractErr(); r != nil {
			return fmt.Errorf("error add user %s to group %s: %s ", u, group, r)
		}
	}
	return nil
}

func removeUsersFromGroup(identityClient *golangsdk.ServiceClient, group string, userList []string) error {
	for _, u := range userList {
		if r := users.RemoveFromGroup(identityClient, group, u).ExtractErr(); r != nil {
			return fmt.Errorf("error remove user %s from group %s: %s", u, group, r)
		}
	}
	return nil
}

// func checkMembership(identityClient *golangsdk.ServiceClient, group string, user string)  error {
