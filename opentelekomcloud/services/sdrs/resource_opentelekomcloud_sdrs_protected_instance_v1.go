package sdrs

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/sdrs/v1/protectedinstances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSdrsProtectedInstanceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSdrsProtectedInstanceV1Create,
		ReadContext:   resourceSdrsProtectedInstanceV1Read,
		UpdateContext: resourceSdrsProtectedInstanceV1Update,
		DeleteContext: resourceSdrsProtectedInstanceV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"priority_station": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_target_server": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"delete_target_eip": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func resourceSdrsProtectedInstanceV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := protectedinstances.CreateOpts{
		GroupID:     d.Get("group_id").(string),
		ServerID:    d.Get("server_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		SubnetID:    d.Get("subnet_id").(string),
		IpAddress:   d.Get("ip_address").(string),
	}

	job, err := protectedinstances.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomcomCloud SDRS Protected Instance: %w", err)
	}

	createTimeout := int(d.Timeout(schema.TimeoutCreate) / time.Second)
	if err := protectedinstances.WaitForJobSuccess(client, createTimeout, job.JobID); err != nil {
		return diag.FromErr(err)
	}

	instanceID, err := protectedinstances.GetJobEntity(client, job.JobID, "protected_instance_id")
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(instanceID.(string))

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "protected-instances", d.Id(), tagList).ExtractErr(); err != nil {
			return fmterr.Errorf("error setting tags of SDRS Protected Instance: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, sdrsClientV1)
	return resourceSdrsProtectedInstanceV1Read(clientCtx, d, meta)
}

func resourceSdrsProtectedInstanceV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	instance, err := protectedinstances.Get(client, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud SDRS Protected Instance: %s", err)
	}

	mErr := multierror.Append(
		d.Set("name", instance.Name),
		d.Set("description", instance.Description),
		d.Set("group_id", instance.GroupID),
		d.Set("server_id", instance.SourceServer),
		d.Set("created_at", instance.CreatedAt),
		d.Set("updated_at", instance.UpdatedAt),
		d.Set("target_id", instance.TargetServer),
		d.Set("priority_station", instance.PriorityStation),
	)

	// save tags
	resourceTags, err := tags.Get(client, "protected-instances", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud SDRS Protected Instance tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud SDRS Protected Instance: %s", err)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceSdrsProtectedInstanceV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	var updateOpts protectedinstances.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "protected-instances", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of SDRS Protected Instance %s: %s", d.Id(), err)
		}
	}

	_, err = protectedinstances.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud SDRS Protected Instance: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, sdrsClientV1)
	return resourceSdrsProtectedInstanceV1Read(clientCtx, d, meta)
}

func resourceSdrsProtectedInstanceV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, sdrsClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.SdrsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	deleteTargetServer := d.Get("delete_target_server").(bool)
	deleteTargetEIP := d.Get("delete_target_eip").(bool)

	deleteOpts := protectedinstances.DeleteOpts{
		DeleteTargetServer: &deleteTargetServer,
		DeleteTargetEip:    &deleteTargetEIP,
	}

	job, err := protectedinstances.Delete(client, d.Id(), deleteOpts).ExtractJobResponse()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud SDRS Protected Instance: %w", err)
	}

	deleteTimeout := int(d.Timeout(schema.TimeoutDelete) / time.Second)
	if err := protectedinstances.WaitForJobSuccess(client, deleteTimeout, job.JobID); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
