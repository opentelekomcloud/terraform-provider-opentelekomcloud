package ims

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v1/others"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/images"

	tagCommon "github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceImsImageV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsImageV2Create,
		ReadContext:   resourceImsImageV2Read,
		UpdateContext: resourceImsImageV2Update,
		DeleteContext: resourceImagesImageV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"max_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"min_ram": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
			},
			// instance_id is required for creating an image from an ECS
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"image_url", "volume_id"},
			},
			"volume_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id", "image_url"},
			},
			// image_url and min_disk are required for creating an image from an OBS
			"image_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id", "volume_id"},
			},
			"min_disk": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"instance_id"},
			},
			// following are valid for creating an image from an OBS
			"os_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_config": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"cmk_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ECS", "FusionCompute", "BMS", "Ironic",
				}, true),
			},
			// following are additional attributus
			"visibility": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceContainerImageTags(d *schema.ResourceData) []tagCommon.ResourceTag {
	var tagList []tagCommon.ResourceTag

	imageTags := d.Get("tags").(map[string]interface{})
	for key, val := range imageTags {
		tagRequest := tagCommon.ResourceTag{
			Key:   key,
			Value: val.(string),
		}
		tagList = append(tagList, tagRequest)
	}

	return tagList
}

func resourceImsImageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}
	client1, err := config.ImageV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client v1: %s", err)
	}

	if !common.HasFilledOpt(d, "instance_id") && !common.HasFilledOpt(d, "image_url") && !common.HasFilledOpt(d, "volume_id") {
		return fmterr.Errorf("error creating OpenTelekomCloud IMS: " +
			"Either 'instance_id', 'volume_id' or 'image_url' must be specified")
	}

	var jobId *string
	imageTags := resourceContainerImageTags(d)

	switch {
	case common.HasFilledOpt(d, "instance_id"):
		createOpts := images.CreateImageFromECSOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			InstanceId:  d.Get("instance_id").(string),
			MaxRam:      d.Get("max_ram").(int),
			MinRam:      d.Get("min_ram").(int),
			ImageTags:   imageTags,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		jobId, err = images.CreateImageFromECS(client, createOpts)
	case common.HasFilledOpt(d, "volume_id"):
		createOpts := images.CreateImageFromDiskOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			VolumeId:    d.Get("volume_id").(string),
			OsVersion:   d.Get("os_version").(string),
			Type:        d.Get("type").(string),
			MaxRam:      d.Get("max_ram").(int),
			MinRam:      d.Get("min_ram").(int),
			ImageTags:   imageTags,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		jobId, err = images.CreateImageFromDisk(client, createOpts)
	case common.HasFilledOpt(d, "image_url"):
		if !common.HasFilledOpt(d, "min_disk") {
			return fmterr.Errorf("error creating OpenTelekomCloud IMS: 'min_disk' must be specified")
		}

		createOpts := images.CreateImageFromOBSOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			ImageUrl:    d.Get("image_url").(string),
			MinDisk:     d.Get("min_disk").(int),
			MaxRam:      d.Get("max_ram").(int),
			MinRam:      d.Get("min_ram").(int),
			OsVersion:   d.Get("os_version").(string),
			IsConfig:    d.Get("is_config").(bool),
			CmkId:       d.Get("cmk_id").(string),
			Type:        d.Get("type").(string),
			ImageTags:   imageTags,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		jobId, err = images.CreateImageFromOBS(client, createOpts)
	}

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud IMS: %s", err)
	}
	log.Printf("[INFO] IMS Job ID: %s", *jobId)

	// Wait for the ims to become available.
	log.Printf("[DEBUG] Waiting for IMS to become available")
	err = others.WaitForJob(client1, *jobId, int(d.Timeout(schema.TimeoutCreate)/time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := others.ShowJob(client1, *jobId)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] IMS ID: %s", entity.Entities.ImageId)
	// Store the ID now
	d.SetId(entity.Entities.ImageId)
	return resourceImsImageV2Read(ctx, d, meta)
}

func GetCloudImage(client *golangsdk.ServiceClient, id string) (*images.ImageInfo, error) {
	imgs, err := images.ListImages(client, images.ListImagesOpts{Id: id})
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve images: %s", err)
	}
	if len(imgs) < 1 {
		return nil, nil
	}
	img := imgs[0]
	log.Printf("[DEBUG] Retrieved Image %s: %#v", id, img)
	return &img, nil
}

func resourceImsImageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	img, err := GetCloudImage(client, d.Id())
	if err != nil {
		return fmterr.Errorf("image %s not found: %s", d.Id(), err)
	}
	log.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	mErr := multierror.Append(
		d.Set("name", img.Name),
		d.Set("visibility", img.Visibility),
		d.Set("file", img.File),
		d.Set("data_origin", img.DataOrigin),
		d.Set("disk_format", img.DiskFormat),
		d.Set("image_size", img.ImageSize),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// Set image tags
	tagList, err := tags.ListImageTags(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud image tags: %s", err)
	}

	tagmap := make(map[string]string)
	for _, val := range tagList {
		tagmap[val.Key] = val.Value
	}
	if err := d.Set("tags", tagmap); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving tags for OpenTelekomCloud image (%s): %s", d.Id(), err)
	}
	return nil
}

func setTagForImage(d *schema.ResourceData, meta interface{}, imageID string, tagmap map[string]interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	taglist := make([]tagCommon.ResourceTag, 0)
	for k, v := range tagmap {
		tag := tagCommon.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	err = tags.BatchAddOrDeleteTags(client, tags.BatchAddOrDeleteTagsOpts{
		ImageId: imageID,
		Action:  "create",
		Tags:    taglist,
	})
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud image tags: %s", err)
	}

	return nil
}

func resourceImsImageV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := make([]images.UpdateImageOpts, 0)
		v := images.UpdateImageOpts{
			Op:    "replace",
			Path:  "/name",
			Value: d.Get("name").(string),
		}
		updateOpts = append(updateOpts, v)

		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
		_, err = images.UpdateImage(client, d.Id(), updateOpts)
		if err != nil {
			return fmterr.Errorf("error updating image: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, err := tags.ListImageTags(client, d.Id())
		if err != nil {
			return fmterr.Errorf("error fetching OpenTelekomCloud image tags: %s", err)
		}

		if len(oldTags) > 0 {
			err = tags.BatchAddOrDeleteTags(client, tags.BatchAddOrDeleteTagsOpts{
				ImageId: d.Id(),
				Action:  "delete",
				Tags:    oldTags,
			})
			if err != nil {
				return fmterr.Errorf("error deleting OpenTelekomCloud image tags: %s", err)
			}
		}

		if common.HasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForImage(d, meta, d.Id(), tagmap)
				if err != nil {
					return fmterr.Errorf("error updating OpenTelekomCloud tags of image:%s", err)
				}
			}
		}
	}

	return resourceImsImageV2Read(ctx, d, meta)
}
