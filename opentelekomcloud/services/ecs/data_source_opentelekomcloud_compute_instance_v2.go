package ecs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/flavors"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
)

func DataSourceComputeInstanceV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstanceV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// just stash the hash for state & diff comparisons
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"access_ip_v4": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip_v6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceComputeInstanceV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	log.Print("[DEBUG] Creating compute client")
	client, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("Error creating compute client: %s", err)
	}

	id := d.Get("id").(string)
	log.Printf("[DEBUG] Attempting to retrieve server %s", id)
	server, err := servers.Get(client, id).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "server"))
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", id, server)

	d.SetId(server.ID)

	mErr := multierror.Append(
		d.Set("name", server.Name),
		d.Set("image_id", server.Image["ID"]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// Get the instance network and address information
	networks, err := FlattenInstanceNetworks(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determine the best IPv4 and IPv6 addresses to access the instance with
	hostv4, hostv6 := GetInstanceAccessAddresses(d, networks)

	// AccessIPv4/v6 isn't standard in OpenStack, but there have been reports
	// of them being used in some environments.
	if server.AccessIPv4 != "" && hostv4 == "" {
		hostv4 = server.AccessIPv4
	}

	if server.AccessIPv6 != "" && hostv6 == "" {
		hostv6 = server.AccessIPv6
	}

	log.Printf("[DEBUG] Setting networks: %+v", networks)

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)
	d.Set("access_ip_v6", hostv6)

	d.Set("metadata", server.Metadata)

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg["name"].(string))
	}

	log.Printf("[DEBUG] Setting security groups: %+v", secGrpNames)

	d.Set("security_groups", secGrpNames)

	flavorID, ok := server.Flavor["id"].(string)
	if !ok {
		return diag.Errorf("Error setting server's flavor: %v", server.Flavor)
	}
	d.Set("flavor_id", flavorID)

	d.Set("key_pair", server.KeyName)
	flavor, err := flavors.Get(client, flavorID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("flavor_name", flavor.Name)

	// Set the instance's image information appropriately
	if err := setImageInformation(client, server, d); err != nil {
		return diag.FromErr(err)
	}

	// Build a custom struct for the availability zone extension
	var serverWithAZ struct {
		servers.Server
		availabilityzones.ServerAvailabilityZoneExt
	}

	// Do another Get so the above work is not disturbed.
	err = servers.Get(client, d.Id()).ExtractInto(&serverWithAZ)
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "server"))
	}
	// Set the availability zone
	d.Set("availability_zone", serverWithAZ.AvailabilityZone)

	// Set the region
	d.Set("region", config.GetRegion(d))

	// Set the current power_state
	currentStatus := strings.ToLower(server.Status)
	switch currentStatus {
	case "active", "shutoff", "error", "migrating", "shelved_offloaded", "shelved":
		d.Set("power_state", currentStatus)
	default:
		return diag.Errorf("Invalid power_state for instance %s: %s", d.Id(), server.Status)
	}

	// Populate tags
	instanceTags, err := tags.List(client, server.ID).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get tags for openstack_compute_instance_v2: %s", err)
	} else {
		d.Set("tags", instanceTags)
	}

	return nil
}
