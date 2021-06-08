package ims

import (
	"context"
	"log"
	"regexp"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/images"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceImagesImageV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesImageV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"size_min": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"size_max": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "name",
			},
			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "asc",
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, false),
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"most_recent": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"properties": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			// Computed values
			"container_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk_gb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"min_ram_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
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
			"file": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func dataSourceImagesImageV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.ImageV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud IMSv2 client: %w", err)
	}

	visibility := resourceImagesImageV2VisibilityFromString(d.Get("visibility").(string))

	listOpts := images.ListOpts{
		Name:       d.Get("name").(string),
		Visibility: visibility,
		Owner:      d.Get("owner").(string),
		Status:     images.ImageStatusActive,
		SizeMin:    int64(d.Get("size_min").(int)),
		SizeMax:    int64(d.Get("size_max").(int)),
		SortKey:    d.Get("sort_key").(string),
		SortDir:    d.Get("sort_direction").(string),
		Tag:        d.Get("tag").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var image images.Image
	allPages, err := images.List(client, listOpts).AllPages()
	if err != nil {
		return fmterr.Errorf("unable to query images: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return fmterr.Errorf("unable to retrieve images: %s", err)
	}

	var filteredImages []images.Image
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		r := regexp.MustCompile(nameRegex.(string))
		for _, image := range allImages {
			if r.MatchString(image.Name) {
				filteredImages = append(filteredImages, image)
			}
		}
		allImages = filteredImages
	}

	properties := d.Get("properties").(map[string]interface{})
	imageProperties := resourceImagesImageV2ExpandProperties(properties)
	if len(filteredImages) > 1 && len(imageProperties) > 0 {
		for _, image := range filteredImages {
			if len(image.Properties) > 0 {
				match := true
				for searchKey, searchValue := range imageProperties {
					imageValue, ok := image.Properties[searchKey]
					if !ok {
						match = false
						break
					}

					if searchValue != imageValue {
						match = false
						break
					}
				}

				if match {
					filteredImages = append(filteredImages, image)
				}
			}
		}
		allImages = filteredImages
	}

	if len(allImages) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again")
	}

	if len(allImages) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] Multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			image = mostRecentImage(allImages)
		} else {
			return fmterr.Errorf("your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true")
		}
	} else {
		image = allImages[0]
	}

	log.Printf("[DEBUG] Single Image found: %s", image.ID)
	d.SetId(image.ID)

	mErr := multierror.Append(nil,
		d.Set("name", image.Name),
		d.Set("tags", image.Tags),
		d.Set("container_format", image.ContainerFormat),
		d.Set("disk_format", image.DiskFormat),
		d.Set("min_disk_gb", image.MinDiskGigabytes),
		d.Set("min_ram_mb", image.MinRAMMegabytes),
		d.Set("owner", image.Owner),
		d.Set("protected", image.Protected),
		d.Set("visibility", image.Visibility),
		d.Set("checksum", image.Checksum),
		d.Set("size_bytes", image.SizeBytes),
		d.Set("created_at", image.CreatedAt.String()),
		d.Set("updated_at", image.UpdatedAt.String()),
		d.Set("file", image.File),
		d.Set("schema", image.Schema),
		d.Set("metadata", image.Metadata),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

type imageSort []images.Image

func (a imageSort) Len() int {
	return len(a)
}

func (a imageSort) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a imageSort) Less(i, j int) bool {
	itime := a[i].CreatedAt
	jtime := a[j].CreatedAt
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []images.Image) images.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
