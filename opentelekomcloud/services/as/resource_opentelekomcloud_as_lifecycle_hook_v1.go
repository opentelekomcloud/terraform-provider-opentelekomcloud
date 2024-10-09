package as

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

	lifecyclehooks "github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/lifecycle_hooks"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceASLifecycleHook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASLifecycleHookCreate,
		ReadContext:   resourceASLifecycleHookRead,
		UpdateContext: resourceASLifecycleHookUpdate,
		DeleteContext: resourceASLifecycleHookDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_lifecycle_hook_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceASLifeCycleHookValidateName,
			},
			"scaling_lifecycle_hook_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					"INSTANCE_TERMINATING", "INSTANCE_LAUNCHING",
				}, true),
				Description: "Specifies the lifecycle hook type. Options: INSTANCE_TERMINATING, INSTANCE_LAUNCHING",
			},
			"default_result": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					"ABANDON", "CONTINUE",
				}, true),
				Description: "Specifies the default lifecycle hook callback operation. Options: CONTINUE, ABANDON(Default)",
				Default:     "ABANDON",
			},
			"default_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				ValidateFunc: validation.IntBetween(60, 86400),
				Description:  "Timeout duration, in seconds.",
				ForceNew:     false,
			},
			"notification_topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "URN of the SMN Topic",
			},
			"notification_metadata": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: resourceASLifeCycleHookValidateNotificationMetadata,
			},
			"notification_topic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceASLifecycleHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := lifecyclehooks.CreateOpts{
		LifecycleHookName:    d.Get("scaling_lifecycle_hook_name").(string),
		LifecycleHookType:    d.Get("scaling_lifecycle_hook_type").(string),
		DefaultResult:        d.Get("default_result").(string),
		DefaultTimeout:       d.Get("default_timeout").(int),
		NotificationTopicUrn: d.Get("notification_topic_urn").(string),
		NotificationMetadata: d.Get("notification_metadata").(string),
	}
	asGroupId := d.Get("scaling_group_id").(string)

	log.Printf("[DEBUG] Create AS Lifecycle Hook Options: %#v", createOpts)
	asLifeCycleHook, err := lifecyclehooks.Create(client, createOpts, asGroupId)
	if err != nil {
		return fmterr.Errorf("error creating AS Lifecycle Hook: %s", err)
	}
	log.Printf("[DEBUG] Create AS LifeCycle Hook %q Success!", asLifeCycleHook.LifecycleHookName)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceASLifecycleHookRead(clientCtx, d, meta)
}

func resourceASLifecycleHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	asLifecycleHookName := d.Get("scaling_lifecycle_hook_name").(string)
	asGroupId := d.Get("scaling_group_id").(string)
	asLifecycleHook, err := lifecyclehooks.Get(client, asGroupId, asLifecycleHookName)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud AutoScaling Lifecycle Hook: %s", err)
	}

	// set properties based on the read info
	mErr := multierror.Append(nil,
		d.Set("scaling_lifecycle_hook_name", asLifecycleHook.LifecycleHookName),
		d.Set("scaling_lifecycle_hook_type", asLifecycleHook.LifecycleHookType),
		d.Set("default_result", asLifecycleHook.DefaultResult),
		d.Set("default_timeout", asLifecycleHook.DefaultTimeout),
		d.Set("notification_topic_urn", asLifecycleHook.NotificationTopicUrn),
		d.Set("notification_metadata", asLifecycleHook.NotificationMetadata),
		d.Set("notification_topic_name", asLifecycleHook.NotificationTopicName),
		d.Set("create_time", asLifecycleHook.CreateTime),
		d.Set("scaling_group_id", asGroupId),
		d.Set("region", config.GetRegion(d)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceASLifecycleHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	asGroupId := d.Get("scaling_group_id").(string)
	asLifecycleHookName := d.Get("scaling_lifecycle_hook_name").(string)
	updateOpts := lifecyclehooks.UpdateOpts{
		LifecycleHookType:    d.Get("scaling_lifecycle_hook_type").(string),
		DefaultResult:        d.Get("default_result").(string),
		DefaultTimeout:       d.Get("default_timeout").(int),
		NotificationTopicUrn: d.Get("notification_topic_urn").(string),
		NotificationMetadata: d.Get("notification_metadata").(string),
	}
	_, err = lifecyclehooks.Update(client, updateOpts, asGroupId, asLifecycleHookName)
	if err != nil {
		return fmterr.Errorf("error updating AS Lifecycle Hook %q: %s", asLifecycleHookName, err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceASLifecycleHookRead(clientCtx, d, meta)
}

func resourceASLifecycleHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	asGroupId := d.Get("scaling_group_id").(string)
	asLifecycleHookName := d.Get("scaling_lifecycle_hook_name").(string)

	log.Printf("[DEBUG] Begin deleting AS Lifecycle Hook %q", asLifecycleHookName)
	err = lifecyclehooks.Delete(client, asGroupId, asLifecycleHookName)
	if err != nil {
		return fmterr.Errorf("error deleting AutoScaling Lifecycle Hook %q: %s", asLifecycleHookName, err)
	}

	return nil
}

func resourceASLifeCycleHookValidateName(v interface{}, _ string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 32 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must contain more than 1 and less than 32 characters", value))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z-_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters, hyphens, and underscores allowed in %q", value))
	}
	return
}

func resourceASLifeCycleHookValidateNotificationMetadata(v interface{}, _ string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 256 {
		errors = append(errors, fmt.Errorf("%q must contain less than 256 characters", value))
	}
	if regexp.MustCompile(`[<>&'()]`).MatchString(value) {
		errors = append(errors, fmt.Errorf("The characters < > & ' ( ) are not allowed in %q", value))
	}
	return
}
