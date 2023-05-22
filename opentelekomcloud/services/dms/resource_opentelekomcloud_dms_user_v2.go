package dms

import (
	"context"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsUsersV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsUsersV2Create,
		ReadContext:   resourceDmsUsersV2Read,
		UpdateContext: resourceDmsUsersV2Update,
		DeleteContext: resourceDmsUsersV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(5, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[\w\-.]+$`),
						"Only lowercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
					),
				),
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_app": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceDmsUsersV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	createOpts := users.CreateOpts{
		UserName:   d.Get("username").(string),
		UserPasswd: d.Get("password").(string),
	}

	instanceId := d.Get("instance_id").(string)

	err = users.Create(client, instanceId, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DMSv2 user: %w", err)
	}

	// Store the instance ID == username
	d.SetId(createOpts.UserName)

	return resourceDmsUsersV2Read(ctx, d, meta)
}

func resourceDmsUsersV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	readUser, err := getUserFromList(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS user")
	}

	if readUser.UserName == "" {
		return fmterr.Errorf("User %s doesn't exist", d.Id())
	}

	mErr := multierror.Append(
		d.Set("username", readUser.UserName),
		d.Set("role", readUser.Role),
		d.Set("default_app", readUser.DefaultApp),
		d.Set("creation_time", readUser.CreatedTime),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDmsUsersV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if d.HasChange("password") {
		password := d.Get("password").(string)
		err = users.ResetPassword(client, d.Get("instance_id").(string), d.Id(), password)
		if err != nil {
			return fmterr.Errorf("error updating OpenTelekomCloud DMSv2 User password: %s", err)
		}
	}

	return resourceDmsUsersV2Read(ctx, d, meta)
}

func resourceDmsUsersV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	deleteOpts := users.DeleteOpts{
		Action: "delete",
		Users: []string{
			d.Id(),
		},
	}

	err = users.Delete(client, d.Get("instance_id").(string), deleteOpts)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud DMSv2 instance: %w", err)
	}

	log.Printf("[DEBUG] DMS user %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func getUserFromList(client *golangsdk.ServiceClient, d *schema.ResourceData) (users.Users, error) {
	var readUser users.Users
	v, err := users.List(client, d.Get("instance_id").(string))
	if err != nil {
		return readUser, err
	}

	for _, user := range v {
		if user.UserName == d.Id() {
			readUser = user
			break
		}
	}
	return readUser, nil
}
