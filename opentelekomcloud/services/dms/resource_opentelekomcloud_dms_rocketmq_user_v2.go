package dms

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/users"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const rocketMQUserNotFound = "DMS.00500972"

func ResourceDmsRocketMQUserV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQUserV2Create,
		ReadContext:   resourceDmsRocketMQUserV2Read,
		UpdateContext: resourceDmsRocketMQUserV2Update,
		DeleteContext: resourceDmsRocketMQUserV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z][\w-]{6,63}$`),
					"must start with a letter, consist of 7 to 64 characters and contain only letters, digits, hyphens (-) and underscores (_)"),
			},
			"secret_key": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(8, 32),
			},
			"white_remote_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_topic_perm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_group_perm": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsRocketMQUserV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	createOpts := users.CreateOpts{
		AccessKey:          d.Get("access_key").(string),
		SecretKey:          d.Get("secret_key").(string),
		WhiteRemoteAddress: d.Get("white_remote_address").(string),
	}

	created, err := users.Create(client, instanceID, createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS RocketMQ user: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceID, created.AccessKey))

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQUserV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQUserV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, name, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	v, err := users.Get(client, instanceID, name)
	if err != nil {
		return common.CheckDeletedDiag(d, convertRocketMQNotFound(err, rocketMQUserNotFound), "DMS RocketMQ user")
	}
	log.Printf("[DEBUG] Get DMS RocketMQ user: %s", v.AccessKey)

	mErr := multierror.Append(
		d.Set("instance_id", instanceID),
		d.Set("access_key", v.AccessKey),
		d.Set("secret_key", v.SecretKey),
		d.Set("white_remote_address", v.WhiteRemoteAddress),
		d.Set("admin", v.Admin),
		d.Set("default_topic_perm", v.DefaultTopicPerm),
		d.Set("default_group_perm", v.DefaultGroupPerm),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQUserV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, name, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// The API rejects an empty secret key, so the current one is always passed.
	_, err = users.Update(client, instanceID, name, users.UpdateOpts{
		SecretKey:          d.Get("secret_key").(string),
		WhiteRemoteAddress: pointerto.String(d.Get("white_remote_address").(string)),
	})
	if err != nil {
		return diag.Errorf("error updating DMS RocketMQ user (%s): %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQUserV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQUserV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, name, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err = users.Delete(client, instanceID, name); err != nil {
		return common.CheckDeletedDiag(d, convertRocketMQNotFound(err, rocketMQUserNotFound), "error deleting DMS RocketMQ user")
	}
	return nil
}
