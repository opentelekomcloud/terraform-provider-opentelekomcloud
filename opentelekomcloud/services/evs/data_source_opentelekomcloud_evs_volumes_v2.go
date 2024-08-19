package evs

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/evs/v2/cloudvolumes"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"
)

func DataSourceEvsVolumesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEvsVolumesV2Read,

		Schema: map[string]*schema.Schema{
			"volume_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_type_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shareable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attachments": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"attached_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"attached_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"device_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"server_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bootable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"shareable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"wwn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQueryOpts(d *schema.ResourceData, cfg *cfg.Config) cloudvolumes.ListOpts {
	result := cloudvolumes.ListOpts{
		ID:               d.Get("volume_id").(string),
		Name:             d.Get("name").(string),
		VolumeTypeID:     d.Get("volume_type_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
		ServerID:         d.Get("server_id").(string),
		Status:           d.Get("status").(string),
	}
	if val, ok := d.GetOk("shareable"); ok {
		result.Multiattach = val.(bool)
	}
	return result
}

func sourceEvsAttachment(attachements []cloudvolumes.Attachment, mode string) []map[string]interface{} {
	result := make([]map[string]interface{}, len(attachements))
	for i, attachement := range attachements {
		result[i] = map[string]interface{}{
			"id":            attachement.AttachmentID,
			"attached_at":   attachement.AttachedAt,
			"attached_mode": mode,
			"device_name":   attachement.Device,
			"server_id":     attachement.ServerID,
		}
	}
	return result
}

func sourceEvsVolumes(volumes []cloudvolumes.Volume) ([]map[string]interface{}, []string, error) {
	result := make([]map[string]interface{}, len(volumes))
	ids := make([]string, len(volumes))

	for i, volume := range volumes {
		vMap := map[string]interface{}{
			"id":                volume.ID,
			"attachments":       sourceEvsAttachment(volume.Attachments, volume.Metadata.AttachedMode),
			"availability_zone": volume.AvailabilityZone,
			"description":       volume.Description,
			"volume_type":       volume.VolumeType,
			"name":              volume.Name,
			"service_type":      volume.ServiceType,
			"shareable":         volume.Multiattach,
			"size":              volume.Size,
			"status":            volume.Status,
			"create_at":         volume.CreatedAt,
			"update_at":         volume.UpdatedAt,
			"tags":              volume.Tags,
			"wwn":               volume.WWN,
		}
		bootable, err := strconv.ParseBool(volume.Bootable)
		if err != nil {
			return nil, nil, fmt.Errorf("the bootable of volume (%s) connot be converted from boolen to string",
				volume.ID)
		}

		vMap["bootable"] = bootable

		result[i] = vMap
		ids[i] = volume.ID
	}
	return result, ids, nil
}

func dataSourceEvsVolumesV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.BlockStorageV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating EVS v2 client: %s", err)
	}

	volumes, err := cloudvolumes.List(client, buildQueryOpts(d, config))
	if err != nil {
		return diag.Errorf("error getting the EVS volume list form server: %s", err)
	}

	filter := d.Get("tags").(map[string]interface{})
	filterVolumes := filterVolumeListByTags(volumes, filter)
	log.Printf("filter %d EVS volumes from %d through options %v", len(filterVolumes), len(volumes), filter)

	vMap, ids, err := sourceEvsVolumes(filterVolumes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil, d.Set("volumes", vMap))
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving the detailed information of the EVS disks to state: %s", mErr)
	}
	return nil
}

func filterVolumeListByTags(volumes []cloudvolumes.Volume, filter map[string]interface{}) []cloudvolumes.Volume {
	result := make([]cloudvolumes.Volume, 0, len(volumes))
	for _, volume := range volumes {
		if common.MapContains(volume.Tags, filter) {
			result = append(result, volume)
		}
	}
	return result
}
