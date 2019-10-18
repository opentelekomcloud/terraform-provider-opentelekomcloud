package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	imageservice_v2 "github.com/huaweicloud/golangsdk/openstack/imageservice/v2/images"
	"github.com/huaweicloud/golangsdk/openstack/ims/v2/cloudimages"
	"github.com/huaweicloud/golangsdk/openstack/ims/v2/tags"
)

func resourceImsDataImageV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceImsDataImageV2Create,
		Read:   resourceImsDataImageV2Read,
		Update: resourceImsDataImageV2Update,
		Delete: resourceImagesImageV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceImsDataImageV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ims_Client, err := config.imageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud image client: %s", err)
	}

	if !hasFilledOpt(d, "volume_id") && !hasFilledOpt(d, "image_url") {
		return fmt.Errorf("Error creating OpenTelekomCloud IMS: " +
			"Either 'volume_id' or 'image_url' must be specified")
	}

	v := new(cloudimages.JobResponse)
	if hasFilledOpt(d, "volume_id") {
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

		v, err = cloudimages.CreateImageByServer(ims_Client, createOpts).ExtractJobResponse()
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud IMS: %s", err)
		}
	} else {
		if !hasFilledOpt(d, "min_disk") {
			return fmt.Errorf("Error creating OpenTelekomCloud IMS: 'min_disk' must be specified")
		}

		imsV1_Client, err := config.imageV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud image client: %s", err)
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
		v, err = cloudimages.CreateDataImageByOBS(imsV1_Client, createOpts).ExtractJobResponse()
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud IMS: %s", err)
		}
	}

	log.Printf("[INFO] IMS Job ID: %s", v.JobID)

	// Wait for the ims to become available.
	log.Printf("[DEBUG] Waiting for IMS to become available")
	err = cloudimages.WaitForJobSuccess(ims_Client, int(d.Timeout(schema.TimeoutCreate)/time.Second), v.JobID)
	if err != nil {
		return err
	}

	entity, err := cloudimages.GetJobEntity(ims_Client, v.JobID, "__data_images")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		log.Printf("[INFO] IMS ID: %s", id)
		// Store the ID now
		d.SetId(id)

		if hasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForImage(d, meta, id, tagmap)
				if err != nil {
					return fmt.Errorf("Error setting OpenTelekomCloud tags of image:%s", err)
				}
			}
		}
		return resourceImsDataImageV2Read(d, meta)
	}
	return fmt.Errorf("Unexpected conversion error in resourceImsDataImageV2Create.")
}

func resourceImsDataImageV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ims_Client, err := config.imageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud image client: %s", err)
	}

	img, err := getCloudimage(ims_Client, d.Id())
	if err != nil {
		return fmt.Errorf("Image %s not found: %s", d.Id(), err)
	}
	log.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	d.Set("name", img.Name)
	d.Set("description", img.Description)
	d.Set("visibility", img.Visibility)
	d.Set("data_origin", img.DataOrigin)
	d.Set("disk_format", img.DiskFormat)
	d.Set("image_size", img.ImageSize)

	// Set image tags
	Taglist, err := tags.Get(ims_Client, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching OpenTelekomCloud image tags: %s", err)
	}

	tagmap := make(map[string]string)
	for _, val := range Taglist.Tags {
		tagmap[val.Key] = val.Value
	}
	if err := d.Set("tags", tagmap); err != nil {
		return fmt.Errorf("[DEBUG] Error saving tags for OpenTelekomCloud image (%s): %s", d.Id(), err)
	}
	return nil
}

func resourceImsDataImageV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ims_Client, err := config.imageV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud image client: %s", err)
	}

	updateOpts := make(imageservice_v2.UpdateOpts, 0)

	if d.HasChange("name") {
		v := imageservice_v2.ReplaceImageName{NewName: d.Get("name").(string)}
		updateOpts = append(updateOpts, v)

		log.Printf("[DEBUG] Update Options: %#v", updateOpts)

		_, err = imageservice_v2.Update(ims_Client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating image: %s", err)
		}
	}

	if d.HasChange("tags") {
		oldTags, err := tags.Get(ims_Client, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Error fetching OpenTelekomCloud image tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteopts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(ims_Client, d.Id(), deleteopts)
			if deleteTags.Err != nil {
				return fmt.Errorf("Error deleting OpenTelekomCloud image tags: %s", deleteTags.Err)
			}
		}

		if hasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForImage(d, meta, d.Id(), tagmap)
				if err != nil {
					return fmt.Errorf("Error updating OpenTelekomCloud tags of image:%s", err)
				}
			}
		}
	}

	return resourceImsDataImageV2Read(d, meta)
}
