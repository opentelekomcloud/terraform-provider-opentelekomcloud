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

	v2 "github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/images"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/cloudimages"
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

func resourceContainerImageTags(d *schema.ResourceData) []cloudimages.ImageTag {
	var tagList []cloudimages.ImageTag

	imageTags := d.Get("tags").(map[string]interface{})
	for key, val := range imageTags {
		tagRequest := cloudimages.ImageTag{
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

	if !common.HasFilledOpt(d, "instance_id") && !common.HasFilledOpt(d, "image_url") && !common.HasFilledOpt(d, "volume_id") {
		return fmterr.Errorf("error creating OpenTelekomCloud IMS: " +
			"Either 'instance_id', 'volume_id' or 'image_url' must be specified")
	}

	var v *cloudimages.JobResponse
	imageTags := resourceContainerImageTags(d)
	switch {
	case common.HasFilledOpt(d, "instance_id"):
		createOpts := &cloudimages.CreateByServerOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			InstanceId:  d.Get("instance_id").(string),
			MaxRam:      d.Get("max_ram").(int),
			MinRam:      d.Get("min_ram").(int),
			ImageTags:   imageTags,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = cloudimages.CreateImageByServer(client, createOpts).ExtractJobResponse()
	case common.HasFilledOpt(d, "volume_id"):
		createOpts := &cloudimages.CreateByVolumeOpts{
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
		v, err = cloudimages.CreateImageByVolume(client, createOpts).ExtractJobResponse()
	case common.HasFilledOpt(d, "image_url"):
		if !common.HasFilledOpt(d, "min_disk") {
			return fmterr.Errorf("error creating OpenTelekomCloud IMS: 'min_disk' must be specified")
		}

		createOpts := &cloudimages.CreateByOBSOpts{
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
		v, err = cloudimages.CreateImageByOBS(client, createOpts).ExtractJobResponse()
	}

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud IMS: %s", err)
	}
	log.Printf("[INFO] IMS Job ID: %s", v.JobID)

	// Wait for the ims to become available.
	log.Printf("[DEBUG] Waiting for IMS to become available")
	err = cloudimages.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := cloudimages.GetJobEntity(client, v.JobID, "image_id")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] IMS ID: %s", id)
		// Store the ID now
		d.SetId(id)
		return resourceImsImageV2Read(ctx, d, meta)
	}
	return fmterr.Errorf("unexpected conversion error in resourceImsImageV2Create.")
}

func GetCloudImage(client *golangsdk.ServiceClient, id string) (*cloudimages.Image, error) {
	listOpts := &cloudimages.ListOpts{
		ID:    id,
		Limit: 1,
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, fmt.Errorf("unable to find images %s: Maybe not existed", id)
	}

	img := allImages[0]
	if img.ID != id {
		return nil, fmt.Errorf("unexpected images ID")
	}
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
	tagList, err := tags.Get(client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud image tags: %s", err)
	}

	tagmap := make(map[string]string)
	for _, val := range tagList.Tags {
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

	rId := imageID
	taglist := make([]tags.Tag, 0)
	for k, v := range tagmap {
		tag := tags.Tag{
			Key:   k,
			Value: v.(string),
		}
		taglist = append(taglist, tag)
	}

	createOpts := tags.BatchOpts{Action: tags.ActionCreate, Tags: taglist}
	createTags := tags.BatchAction(client, rId, createOpts)
	if createTags.Err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud image tags: %s", createTags.Err)
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
		updateOpts := make(v2.UpdateOpts, 0)
		v := v2.ReplaceImageName{NewName: d.Get("name").(string)}
		updateOpts = append(updateOpts, v)

		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
		_, err = v2.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("error updating image: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, err := tags.Get(client, d.Id()).Extract()
		if err != nil {
			return fmterr.Errorf("error fetching OpenTelekomCloud image tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteopts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(client, d.Id(), deleteopts)
			if deleteTags.Err != nil {
				return fmterr.Errorf("error deleting OpenTelekomCloud image tags: %s", deleteTags.Err)
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
