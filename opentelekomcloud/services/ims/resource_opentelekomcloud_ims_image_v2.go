package ims

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	imageservice_v2 "github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/images"
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
				ConflictsWith: []string{"image_url"},
			},
			// image_url and min_disk are required for creating an image from an OBS
			"image_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"instance_id"},
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
		},
	}
}

func resourceContainerImageTags(d *schema.ResourceData) []cloudimages.ImageTag {
	var tags []cloudimages.ImageTag

	image_tags := d.Get("tags").(map[string]interface{})
	for key, val := range image_tags {
		tagRequest := cloudimages.ImageTag{
			Key:   key,
			Value: val.(string),
		}
		tags = append(tags, tagRequest)
	}
	return tags
}

func resourceImsImageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ims_Client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud image client: %s", err)
	}

	if !common.HasFilledOpt(d, "instance_id") && !common.HasFilledOpt(d, "image_url") {
		return fmterr.Errorf("Error creating OpenTelekomCloud IMS: " +
			"Either 'instance_id' or 'image_url' must be specified")
	}

	v := new(cloudimages.JobResponse)
	image_tags := resourceContainerImageTags(d)
	if common.HasFilledOpt(d, "instance_id") {
		createOpts := &cloudimages.CreateByServerOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			InstanceId:  d.Get("instance_id").(string),
			MaxRam:      d.Get("max_ram").(int),
			MinRam:      d.Get("min_ram").(int),
			ImageTags:   image_tags,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = cloudimages.CreateImageByServer(ims_Client, createOpts).ExtractJobResponse()
	} else {
		if !common.HasFilledOpt(d, "min_disk") {
			return fmterr.Errorf("Error creating OpenTelekomCloud IMS: 'min_disk' must be specified")
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
			ImageTags:   image_tags,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = cloudimages.CreateImageByOBS(ims_Client, createOpts).ExtractJobResponse()
	}

	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud IMS: %s", err)
	}
	log.Printf("[INFO] IMS Job ID: %s", v.JobID)

	// Wait for the ims to become available.
	log.Printf("[DEBUG] Waiting for IMS to become available")
	err = cloudimages.WaitForJobSuccess(ims_Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := cloudimages.GetJobEntity(ims_Client, v.JobID, "image_id")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] IMS ID: %s", id)
		// Store the ID now
		d.SetId(id)
		return resourceImsImageV2Read(ctx, d, meta)
	}
	return fmterr.Errorf("Unexpected conversion error in resourceImsImageV2Create.")
}

func GetCloudImage(client *golangsdk.ServiceClient, id string) (*cloudimages.Image, error) {
	listOpts := &cloudimages.ListOpts{
		ID:    id,
		Limit: 1,
	}
	allPages, err := cloudimages.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("Unable to query images: %s", err)
	}

	allImages, err := cloudimages.ExtractImages(allPages)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve images: %s", err)
	}

	if len(allImages) < 1 {
		return nil, fmt.Errorf("Unable to find images %s: Maybe not existed", id)
	}

	img := allImages[0]
	if img.ID != id {
		return nil, fmt.Errorf("Unexpected images ID")
	}
	log.Printf("[DEBUG] Retrieved Image %s: %#v", id, img)
	return &img, nil
}

func resourceImsImageV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ims_Client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud image client: %s", err)
	}

	img, err := GetCloudImage(ims_Client, d.Id())
	if err != nil {
		return fmterr.Errorf("Image %s not found: %s", d.Id(), err)
	}
	log.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	d.Set("name", img.Name)
	d.Set("visibility", img.Visibility)
	d.Set("file", img.File)
	d.Set("schema", img.Schema)
	d.Set("data_origin", img.DataOrigin)
	d.Set("disk_format", img.DiskFormat)
	d.Set("image_size", img.ImageSize)

	// Set image tags
	Taglist, err := tags.Get(ims_Client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("Error fetching OpenTelekomCloud image tags: %s", err)
	}

	tagmap := make(map[string]string)
	for _, val := range Taglist.Tags {
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
		return fmt.Errorf("Error creating OpenTelekomCloud image client: %s", err)
	}

	rId := imageID
	taglist := []tags.Tag{}
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
		return fmt.Errorf("Error creating OpenTelekomCloud image tags: %s", createTags.Err)
	}

	return nil
}

func resourceImsImageV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ims_Client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("Error creating OpenTelekomCloud image client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := make(imageservice_v2.UpdateOpts, 0)
		v := imageservice_v2.ReplaceImageName{NewName: d.Get("name").(string)}
		updateOpts = append(updateOpts, v)

		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
		_, err = imageservice_v2.Update(ims_Client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmterr.Errorf("Error updating image: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, err := tags.Get(ims_Client, d.Id()).Extract()
		if err != nil {
			return fmterr.Errorf("Error fetching OpenTelekomCloud image tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteopts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(ims_Client, d.Id(), deleteopts)
			if deleteTags.Err != nil {
				return fmterr.Errorf("Error deleting OpenTelekomCloud image tags: %s", deleteTags.Err)
			}
		}

		if common.HasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForImage(d, meta, d.Id(), tagmap)
				if err != nil {
					return fmterr.Errorf("Error updating OpenTelekomCloud tags of image:%s", err)
				}
			}
		}
	}

	return resourceImsImageV2Read(ctx, d, meta)
}
