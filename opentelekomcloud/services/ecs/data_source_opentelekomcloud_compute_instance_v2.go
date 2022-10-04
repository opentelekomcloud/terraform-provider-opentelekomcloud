package ecs

import (
	"context"
	"crypto/rsa"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"golang.org/x/crypto/ssh"

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
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"ssh_private_key_path": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"encrypted_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
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
			"admin_pass": {
				Type:     schema.TypeString,
				Computed: true,
				// Sensitive: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"power_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func dataSourceComputeInstanceV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	log.Print("[DEBUG] Creating compute client")
	client, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateV2Client, err)
	}

	var allServers []servers.Server
	if serverId := d.Get("id").(string); serverId != "" {
		server, err := servers.Get(client, serverId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmterr.Errorf("no server found")
			}
			return fmterr.Errorf("unable to retrieve OpenTelekomCloud %s server: %w", serverId, err)
		}

		allServers = append(allServers, *server)
	} else {
		serverName := d.Get("name").(string)

		log.Printf("[DEBUG] Attempting to retrieve server %s", serverName)
		allPages, err := servers.List(client, servers.ListOpts{Name: serverName}).AllPages()
		if err != nil {
			return fmterr.Errorf("unable to retrieve OpenTelekomCloud servers: %w", err)
		}
		allServers, err = servers.ExtractServers(allPages)
		if err != nil {
			return fmterr.Errorf("unable to retrieve OpenTelekomCloud servers: %w", err)
		}
	}

	if len(allServers) < 1 {
		return fmterr.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	server := allServers[0]
	log.Printf("[DEBUG] Retrieved server %s: %+v", server.Name, server)

	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "server"))
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", server.ID, server)

	d.SetId(server.ID)

	mErr := multierror.Append(
		d.Set("name", server.Name),
		d.Set("image_id", server.Image["ID"]),
	)

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

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg["name"].(string))
	}

	log.Printf("[DEBUG] Setting security groups: %+v", secGrpNames)

	flavorID, ok := server.Flavor["id"].(string)
	if !ok {
		return diag.Errorf("Error setting server's flavor: %v", server.Flavor)
	}
	flavor, err := flavors.Get(client, flavorID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	// Set the instance's image information appropriately
	if err := setImageInformation(client, &server, d); err != nil {
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

	// Set the current power_state
	currentStatus := strings.ToLower(server.Status)
	switch currentStatus {
	case "active", "shutoff", "error", "migrating", "shelved_offloaded", "shelved":
		mErr = multierror.Append(mErr, d.Set("power_state", currentStatus))
	default:
		return fmterr.Errorf("invalid power_state for instance %s: %s", d.Id(), server.Status)
	}

	// Set win instance password
	if v, ok := d.GetOk("ssh_private_key_path"); ok {
		readFile, err := os.ReadFile(v.(string))
		if err != nil {
			return fmterr.Errorf("error reading private key file: %w", err)
		}
		privateKey, err := ssh.ParseRawPrivateKey(readFile)
		if err != nil {
			return fmterr.Errorf("error parsing private key: %w", err)
		}
		pass, err := servers.GetPassword(client, d.Id()).ExtractPassword(privateKey.(*rsa.PrivateKey))
		if err != nil {
			return fmterr.Errorf("error getting password: %w", err)
		}
		mErr = multierror.Append(mErr, d.Set("password", pass))
	} else {
		pass, err := servers.GetPassword(client, d.Id()).ExtractPassword(nil)
		if err != nil {
			return fmterr.Errorf("error getting password: %w", err)
		}
		mErr = multierror.Append(mErr, d.Set("encrypted_password", pass))
	}

	mErr = multierror.Append(mErr,
		d.Set("network", networks),
		d.Set("access_ip_v4", hostv4),
		d.Set("access_ip_v6", hostv6),
		d.Set("metadata", server.Metadata),
		d.Set("security_groups", secGrpNames),
		d.Set("flavor_id", flavorID),
		d.Set("key_pair", server.KeyName),
		d.Set("flavor_name", flavor.Name),
		d.Set("availability_zone", serverWithAZ.AvailabilityZone),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	computeClient, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud ComputeV1 client: %w", err)
	}
	// save tags
	resourceTags, err := tags.Get(computeClient, "cloudservers", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud CloudServers tags: %w", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	mErr = multierror.Append(mErr, d.Set("tags", tagMap))

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting opentelekomcloud_compute_instance_v2 values: %w", err)
	}

	return nil
}
