package ims

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	v2 "github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/images"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/cloudimages"
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

	if !common.HasFilledOpt(d, "volume_id") && !common.HasFilledOpt(d, "image_url") {
		return fmterr.Errorf("error creating OpenTelekomCloud IMS: " +
			"Either 'volume_id' or 'image_url' must be specified")
	}

	var v *cloudimages.JobResponse
	if common.HasFilledOpt(d, "volume_id") {
		var dataImages []cloudimages.DataImage
		dataImageOpts := cloudimages.DataImage{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			VolumeId:    d.Get("volume_id").(string),
		}

		dataImages = append(dataImages, dataImageOpts)
		createOpts := &cloudimages.CreateDataImageByServerOpts{
			DataImages: dataImages,
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)

		v, err = cloudimages.CreateImageByServer(client, createOpts).ExtractJobResponse()
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

		createOpts := &cloudimages.CreateDataImageByOBSOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			ImageUrl:    d.Get("image_url").(string),
			MinDisk:     d.Get("min_disk").(int),
			OsType:      d.Get("os_type").(string),
			CmkId:       d.Get("cmk_id").(string),
		}
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		v, err = cloudimages.CreateDataImageByOBS(v1Client, createOpts).ExtractJobResponse()
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud IMS: %s", err)
		}
	}

	log.Printf("[INFO] IMS Job ID: %s", v.JobID)

	// Wait for the ims to become available.
	log.Printf("[DEBUG] Waiting for IMS to become available")
	err = cloudimages.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return diag.FromErr(err)
	}

	entity, err := cloudimages.GetJobEntity(client, v.JobID, "__data_images")
	if err != nil {
		return diag.FromErr(err)
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] IMS ID: %s", id)
		// Store the ID now
		d.SetId(id)

		if common.HasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForImage(d, meta, id, tagmap)
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
	taglist, err := tags.Get(client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud image tags: %s", err)
	}

	tagMap := make(map[string]string)
	for _, val := range taglist.Tags {
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

	updateOpts := make(v2.UpdateOpts, 0)

	if d.HasChange("name") {
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
			deleteOpts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(client, d.Id(), deleteOpts)
			if deleteTags.Err != nil {
				return fmterr.Errorf("error deleting OpenTelekomCloud image tags: %s", deleteTags.Err)
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
