package bms

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/bms/v2/servers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceBMSServersV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBMSServersV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_ip_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
			"security_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"locked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"config_drive": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kernel_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hypervisor_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBMSServersV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	bmsClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating compute v2 client: %w", err)
	}

	listServerOpts := servers.ListOpts{
		ID:         d.Id(),
		Name:       d.Get("name").(string),
		Status:     d.Get("status").(string),
		UserID:     d.Get("user_id").(string),
		KeyName:    d.Get("key_name").(string),
		FlavorID:   d.Get("flavor_id").(string),
		ImageID:    d.Get("image_id").(string),
		HostStatus: d.Get("host_status").(string),
	}
	pages, err := servers.List(bmsClient, listServerOpts)

	if err != nil {
		return fmterr.Errorf("unable to retrieve bms server: %s", err)
	}

	if len(pages) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}
	if len(pages) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}
	server := pages[0]

	log.Printf("[INFO] Retrieved BMS Server using given filter %s: %+v", server.ID, server)
	d.SetId(server.ID)

	var secGroups []map[string]interface{}
	for _, value := range server.SecurityGroups {
		mapping := map[string]interface{}{
			"name": value.Name,
		}
		secGroups = append(secGroups, mapping)
	}

	mErr := multierror.Append(
		d.Set("user_id", server.UserID),
		d.Set("name", server.Name),
		d.Set("status", server.Status),
		d.Set("host_status", server.HostStatus),
		d.Set("host_id", server.HostID),
		d.Set("flavor_id", server.Flavor.ID),
		d.Set("metadata", server.Metadata),
		d.Set("tenant_id", server.TenantID),
		d.Set("image_id", server.Image.ID),
		d.Set("access_ip_v4", server.AccessIPv4),
		d.Set("access_ip_v6", server.AccessIPv6),
		d.Set("progress", server.Progress),
		d.Set("key_name", server.KeyName),
		d.Set("security_groups", secGroups),
		d.Set("locked", server.Locked),
		d.Set("config_drive", server.ConfigDrive),
		d.Set("availability_zone", server.AvailabilityZone),
		d.Set("description", server.Description),
		d.Set("kernel_id", server.KernelId),
		d.Set("hypervisor_hostname", server.HypervisorHostName),
		d.Set("instance_name", server.InstanceName),
		d.Set("tags", server.Tags),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	networks, err := flattenServerNetwork(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network", networks); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving network to state for OpenTelekomCloud server (%s): %s", d.Id(), err)
	}

	return nil
}
