package ims

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v1/others"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/images"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceImsDataImageV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImsDataImageV2Create,
		ReadContext:   resourceImsDataImageV2Read,
		UpdateContext: resourceImsDataImageV2Update,
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
			"volume_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"image_url"},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			// image_url and min_disk are required for creating an image from an OBS
			"image_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"volume_id"},
			},
			"min_disk": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"volume_id"},
			},
			// following are valid for creating an image from an OBS
			"os_type": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"volume_id"},
				ValidateFunc: validation.StringInSlice([]string{
					"Windows", "Linux",
				}, true),
			},
			"cmk_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"volume_id"},
			},
			// following are additional attributes
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

func resourceImsDataImageV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}
	client1, err := config.ImageV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client v1: %s", err)
	}

	if !common.HasFilledOpt(d, "volume_id") && !common.HasFilledOpt(d, "image_url") {
		return fmterr.Errorf("error creating OpenTelekomCloud IMS: " +
			"Either 'volume_id' or 'image_url' must be specified")
	}

	var v *string
	if common.HasFilledOpt(d, "volume_id") {
		var dataImages []images.ECSDataImage
		dataImageOpts := images.ECSDataImage{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			VolumeId:    d.Get("volume_id").(string),
		}

		dataImages = append(dataImages, dataImageOpts)
		createOpts := images.CreateImageFromECSOpts{
			DataImages: dataImages,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = images.CreateImageFromECS(client, createOpts)
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud IMS: %s", err)
		}
	} else {
		if !common.HasFilledOpt(d, "min_disk") {
			return fmterr.Errorf("error creating OpenTelekomCloud IMS: 'min_disk' must be specified")
		}

		v1Client, err := config.ImageV1Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
		}

		createOpts := images.CreateImageFromOBSOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			ImageUrl:    d.Get("image_url").(string),
			MinDisk:     d.Get("min_disk").(int),
			OsType:      d.Get("os_type").(string),
			CmkId:       d.Get("cmk_id").(string),
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = images.CreateImageFromOBS(v1Client, createOpts)
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud IMS: %s", err)
		}
	}
	log.Printf("[INFO] IMS Job ID: %s", *v)

	// Wait for the ims to become available.
	log.Printf("[DEBUG] Waiting for IMS to become available")
	err = others.WaitForJob(client1, *v, int(d.Timeout(schema.TimeoutCreate)/time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := others.ShowJob(client1, *v)
	if err != nil {
		return diag.FromErr(err)
	}

	if entity.JobId != "" {
		log.Printf("[INFO] IMS ID: %s", entity.Entities.ImageId)
		// Store the ID now
		d.SetId(entity.Entities.ImageId)

		if common.HasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForImage(d, meta, entity.Entities.SubJobsResult[0].Entities.ImageId, tagmap)
				if err != nil {
					return fmterr.Errorf("error setting OpenTelekomCloud tags of image:%s", err)
				}
			}
		}
		return resourceImsDataImageV2Read(ctx, d, meta)
	}
	return fmterr.Errorf("unexpected conversion error in resourceImsDataImageV2Create.")
}

func resourceImsDataImageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		d.Set("description", img.Description),
		d.Set("visibility", img.Visibility),
		d.Set("data_origin", img.DataOrigin),
		d.Set("disk_format", img.DiskFormat),
		d.Set("image_size", img.ImageSize),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// Set image tags
	taglist, err := tags.ListImageTags(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud image tags: %s", err)
	}

	tagMap := make(map[string]string)
	for _, val := range taglist {
		tagMap[val.Key] = val.Value
	}
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving tags for OpenTelekomCloud image (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceImsDataImageV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud image client: %s", err)
	}

	var updateOpts []images.UpdateImageOpts
	if d.HasChange("name") {
		v := images.UpdateImageOpts{
			Op:    "replace",
			Path:  "/name",
			Value: d.Get("name").(string),
		}
		updateOpts = append(updateOpts, v)
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
			deleteOpts := tags.BatchAddOrDeleteTagsOpts{
				ImageId: d.Id(),
				Action:  "delete",
				Tags:    oldTags,
			}
			err = tags.BatchAddOrDeleteTags(client, deleteOpts)
			if err != nil {
				return fmterr.Errorf("error deleting OpenTelekomCloud image tags: %s", err)
			}
		}

		if common.HasFilledOpt(d, "tags") {
			tagMap := d.Get("tags").(map[string]interface{})
			if len(tagMap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagMap)
				err = setTagForImage(d, meta, d.Id(), tagMap)
				if err != nil {
					return fmterr.Errorf("error updating OpenTelekomCloud tags of image:%s", err)
				}
			}
		}
	}

	return resourceImsDataImageV2Read(ctx, d, meta)
}
