package opentelekomcloud

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	bms "github.com/huaweicloud/golangsdk/openstack/bms/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/bms/v2/tags"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/keypairs"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/secgroups"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/startstop"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/flavors"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/images"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	ecstags "github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservertags"
)

func resourceComputeBMSInstanceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeBMSInstanceV2Create,
		Read:   resourceComputeBMSInstanceV2Read,
		Update: resourceComputeBMSInstanceV2Update,
		Delete: resourceComputeBMSInstanceV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"image_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"image_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_BMS_FLAVOR_NAME", nil),
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_BMS_FLAVOR_NAME", nil),
			},

			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// just stash the hash for state & diff comparisons
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						hash := sha1.Sum([]byte(v.(string)))
						return hex.EncodeToString(hash[:])
					default:
						return ""
					}
				},
			},
			"security_groups": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v4": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"fixed_ip_v6": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_network": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"config_drive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"admin_pass": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
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
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: validateECSTagValue,
			},
			"stop_before_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kernel_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_device": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"destination_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"boot_index": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"guest_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func bmsTagsCreate(client *golangsdk.ServiceClient, server_id string) error {

	createOpts := tags.CreateOpts{
		Tag: []string{"__type_baremetal"},
	}

	_, err := tags.Create(client, server_id, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Tags: %s", err)
	}
	return nil
}

func resourceComputeBMSInstanceV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2HWClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
	}

	var createOpts servers.CreateOptsBuilder

	// Determines the Image ID using the following rules:
	// If an image_id was specified, use it.
	// If an image_name was specified, look up the image ID, report if error.
	imageId, err := getImageId(computeClient, d)
	if err != nil {
		return err
	}

	flavorId, err := getFlavorId(computeClient, d)
	if err != nil {
		return err
	}

	// Build a list of networks with the information given upon creation.
	// Error out if an invalid network configuration was used.
	allInstanceNetworks, err := getAllServerNetwork(d, meta)
	if err != nil {
		return err
	}

	// Build a []servers.Network to pass into the create options.
	networks := expandBmsInstanceNetworks(allInstanceNetworks)

	createOpts = &servers.CreateOpts{
		Name:             d.Get("name").(string),
		ImageRef:         imageId,
		FlavorRef:        flavorId,
		SecurityGroups:   resourceBmsInstanceSecGroupsV2(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		Networks:         networks,
		Metadata:         resourceBmsInstanceMetadataV2(d),
		AdminPass:        d.Get("admin_pass").(string),
		UserData:         []byte(d.Get("user_data").(string)),
	}

	if keyName, ok := d.Get("key_pair").(string); ok && keyName != "" {
		createOpts = &keypairs.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			KeyName:           keyName,
		}
	}

	if vL, ok := d.GetOk("block_device"); ok {
		blockDevices, err := resourceInstanceBlockDevicesV2(d, vL.([]interface{}))
		if err != nil {
			return err
		}

		createOpts = &bootfromvolume.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			BlockDevice:       blockDevices,
		}
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	var server *servers.Server
	if _, ok := d.GetOk("block_device"); ok {
		server, err = bootfromvolume.Create(computeClient, createOpts).Extract()
	} else {
		server, err = servers.Create(computeClient, createOpts).Extract()
	}
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud server: %s", err)
	}
	log.Printf("[INFO] Instance ID: %s", server.ID)

	// Store the ID now
	d.SetId(server.ID)

	// Set bms sepcific tag
	bmsClient, err := config.bmsClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud bms client: %s", err)
	}
	err = bmsTagsCreate(bmsClient, d.Id())
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud bms tag: %s", err)
	}

	// Wait for the instance to become running so we can get some attributes
	// that aren't available until later.
	log.Printf(
		"[DEBUG] Waiting for instance (%s) to become running",
		server.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    BmsServerV2StateRefreshFunc(computeClient, server.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			server.ID, err)
	}

	if hasFilledOpt(d, "tags") {
		tagmap := d.Get("tags").(map[string]interface{})
		log.Printf("[DEBUG] Setting tags: %v", tagmap)
		err = setTagForInstance(d, meta, server.ID, tagmap)
		if err != nil {
			log.Printf("[WARN] Error setting tags of bms server:%s, err=%s", server.ID, err)
		}
	}

	return resourceComputeBMSInstanceV2Read(d, meta)
}

func resourceComputeBMSInstanceV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2HWClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
	}

	server, err := bms.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "server")
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	d.Set("name", server.Name)

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(d, meta)
	if err != nil {
		return err
	}

	// Determine the best IPv4 and IPv6 addresses to access the instance with
	hostv4, hostv6 := getInstanceAccessAddresses(d, networks)

	// AccessIPv4/v6 isn't standard in OpenTelekomCloud, but there have been reports
	// of them being used in some environments.
	if server.AccessIPv4 != "" && hostv4 == "" {
		hostv4 = server.AccessIPv4
	}

	if server.AccessIPv6 != "" && hostv6 == "" {
		hostv6 = server.AccessIPv6
	}

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)
	d.Set("access_ip_v6", hostv6)

	// Determine the best IP address to use for SSH connectivity.
	// Prefer IPv4 over IPv6.
	var preferredSSHAddress string
	if hostv4 != "" {
		preferredSSHAddress = hostv4
	} else if hostv6 != "" {
		preferredSSHAddress = hostv6
	}

	if preferredSSHAddress != "" {
		// Initialize the connection info
		d.SetConnInfo(map[string]string{
			"type": "ssh",
			"host": preferredSSHAddress,
		})
	}

	d.Set("metadata", server.Metadata)
	d.Set("flavor_id", server.Flavor.ID)

	flavor, err := flavors.Get(computeClient, server.Flavor.ID).Extract()
	if err != nil {
		return err
	}
	d.Set("flavor_name", flavor.Name)

	// Set the instance's image information appropriately
	if err := setImageInformations(computeClient, server, d); err != nil {
		return err
	}

	d.Set("availability_zone", server.AvailabilityZone)
	d.Set("tenant_id", server.TenantID)
	d.Set("host_status", server.HostStatus)
	d.Set("host_id", server.HostID)
	d.Set("kernel_id", server.KernelId)
	d.Set("user_id", server.UserID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeBMSInstanceV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2HWClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
	}

	var updateOpts servers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if updateOpts != (servers.UpdateOpts{}) {
		_, err := servers.Update(computeClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud server: %s", err)
		}
	}

	if d.HasChange("metadata") {
		oldMetadata, newMetadata := d.GetChange("metadata")
		var metadataToDelete []string

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata.(map[string]interface{}) {
			var found bool
			for newKey := range newMetadata.(map[string]interface{}) {
				if oldKey == newKey {
					found = true
				}
			}

			if !found {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		for _, key := range metadataToDelete {
			err := servers.DeleteMetadatum(computeClient, d.Id(), key).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error deleting metadata (%s) from server (%s): %s", key, d.Id(), err)
			}
		}

		// Update existing metadata and add any new metadata.
		metadataOpts := make(servers.MetadataOpts)
		for k, v := range newMetadata.(map[string]interface{}) {
			metadataOpts[k] = v.(string)
		}

		_, err := servers.UpdateMetadata(computeClient, d.Id(), metadataOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud server (%s) metadata: %s", d.Id(), err)
		}
	}

	if d.HasChange("security_groups") {
		oldSGRaw, newSGRaw := d.GetChange("security_groups")
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secgroupsToAdd := newSGSet.Difference(oldSGSet)
		secgroupsToRemove := oldSGSet.Difference(newSGSet)

		log.Printf("[DEBUG] Security groups to add: %v", secgroupsToAdd)

		log.Printf("[DEBUG] Security groups to remove: %v", secgroupsToRemove)

		for _, g := range secgroupsToRemove.List() {
			err := secgroups.RemoveServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					continue
				}

				return fmt.Errorf("Error removing security group (%s) from OpenTelekomCloud server (%s): %s", g, d.Id(), err)
			} else {
				log.Printf("[DEBUG] Removed security group (%s) from instance (%s)", g, d.Id())
			}
		}

		for _, g := range secgroupsToAdd.List() {
			err := secgroups.AddServer(computeClient, d.Id(), g.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return fmt.Errorf("Error adding security group (%s) to OpenTelekomCloud server (%s): %s", g, d.Id(), err)
			}
			log.Printf("[DEBUG] Added security group (%s) to instance (%s)", g, d.Id())
		}
	}

	if d.HasChange("admin_pass") {
		if newPwd, ok := d.Get("admin_pass").(string); ok {
			err := servers.ChangeAdminPassword(computeClient, d.Id(), newPwd).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error changing admin password of OpenTelekomCloud server (%s): %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("tags") {
		computeClient, err := config.computeV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud compute v1 client: %s", err)
		}
		oldTags, err := ecstags.Get(computeClient, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Error fetching OpenTelekomCloud instance tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteopts := ecstags.BatchOpts{Action: ecstags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := ecstags.BatchAction(computeClient, d.Id(), deleteopts)
			if deleteTags.Err != nil {
				return fmt.Errorf("Error updating OpenTelekomCloud instance tags: %s", deleteTags.Err)
			}
		}

		if hasFilledOpt(d, "tags") {
			tagmap := d.Get("tags").(map[string]interface{})
			if len(tagmap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagmap)
				err = setTagForInstance(d, meta, d.Id(), tagmap)
				if err != nil {
					return fmt.Errorf("Error updating tags of instance:%s, err:%s", d.Id(), err)
				}
			}
		}
	}

	if d.HasChange("flavor_id") || d.HasChange("flavor_name") {
		var newFlavorId string
		var err error
		if d.HasChange("flavor_id") {
			newFlavorId = d.Get("flavor_id").(string)
		} else {
			newFlavorName := d.Get("flavor_name").(string)
			newFlavorId, err = flavors.IDFromName(computeClient, newFlavorName)
			if err != nil {
				return err
			}
		}

		resizeOpts := &servers.ResizeOpts{
			FlavorRef: newFlavorId,
		}
		log.Printf("[DEBUG] Resize configuration: %#v", resizeOpts)
		err = servers.Resize(computeClient, d.Id(), resizeOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error resizing OpenTelekomCloud server: %s", err)
		}

		// Wait for the instance to finish resizing.
		log.Printf("[DEBUG] Waiting for instance (%s) to finish resizing", d.Id())

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"RESIZE"},
			Target:     []string{"VERIFY_RESIZE"},
			Refresh:    BmsServerV2StateRefreshFunc(computeClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for instance (%s) to resize: %s", d.Id(), err)
		}

		// Confirm resize.
		log.Printf("[DEBUG] Confirming resize")
		err = servers.ConfirmResize(computeClient, d.Id()).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error confirming resize of OpenTelekomCloud server: %s", err)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"VERIFY_RESIZE"},
			Target:     []string{"ACTIVE"},
			Refresh:    BmsServerV2StateRefreshFunc(computeClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for instance (%s) to confirm resize: %s", d.Id(), err)
		}
	}

	return resourceComputeBMSInstanceV2Read(d, meta)
}

func resourceComputeBMSInstanceV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.computeV2HWClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
	}

	if d.Get("stop_before_destroy").(bool) {
		err = startstop.Stop(computeClient, d.Id()).ExtractErr()
		if err != nil {
			log.Printf("[WARN] Error stopping OpenTelekomCloud instance: %s", err)
		} else {
			stopStateConf := &resource.StateChangeConf{
				Pending:    []string{"ACTIVE"},
				Target:     []string{"SHUTOFF"},
				Refresh:    BmsServerV2StateRefreshFunc(computeClient, d.Id()),
				Timeout:    3 * time.Minute,
				Delay:      10 * time.Second,
				MinTimeout: 3 * time.Second,
			}
			log.Printf("[DEBUG] Waiting for instance (%s) to stop", d.Id())
			_, err = stopStateConf.WaitForState()
			if err != nil {
				log.Printf("[WARN] Error waiting for instance (%s) to stop: %s, proceeding to delete", d.Id(), err)
			}
		}
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud Instance %s", d.Id())
	err = servers.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud server: %s", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "SHUTOFF"},
		Target:     []string{"DELETED", "SOFT_DELETED"},
		Refresh:    BmsServerV2StateRefreshFunc(computeClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	d.SetId("")
	return nil
}

// ServerV2StateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// an OpenTelekomCloud instance.
func BmsServerV2StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := servers.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return s, "DELETED", nil
			}
			return nil, "", err
		}

		return s, s.Status, nil
	}
}

func resourceBmsInstanceSecGroupsV2(d *schema.ResourceData) []string {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secgroups := make([]string, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secgroups[i] = raw.(string)
	}
	return secgroups
}

func resourceBmsInstanceMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func getImageId(computeClient *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	// If block_device was used, an Image does not need to be specified, unless an image/local
	// combination was used. This emulates normal boot behavior. Otherwise, ignore the image altogether.
	if vL, ok := d.GetOk("block_device"); ok {
		needImage := false
		for _, v := range vL.([]interface{}) {
			vM := v.(map[string]interface{})
			if vM["source_type"] == "image" && vM["destination_type"] == "local" {
				needImage = true
			}
		}
		if !needImage {
			return "", nil
		}
	}

	if imageId := d.Get("image_id").(string); imageId != "" {
		return imageId, nil
	} else {
		// try the OS_IMAGE_ID environment variable
		if v := os.Getenv("OS_IMAGE_ID"); v != "" {
			return v, nil
		}
	}

	imageName := d.Get("image_name").(string)
	if imageName == "" {
		// try the OS_IMAGE_NAME environment variable
		if v := os.Getenv("OS_IMAGE_NAME"); v != "" {
			imageName = v
		}
	}

	if imageName != "" {
		imageId, err := images.IDFromName(computeClient, imageName)
		if err != nil {
			return "", err
		}
		return imageId, nil
	}

	return "", fmt.Errorf("Neither a image ID, or image name were able to be determined.")
}

func setImageInformations(computeClient *golangsdk.ServiceClient, server *bms.Server, d *schema.ResourceData) error {
	imageId := server.Image.ID
	if imageId != "" {
		d.Set("image_id", imageId)
		if image, err := images.Get(computeClient, imageId).Extract(); err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				// If the image name can't be found, set the value to "Image not found".
				// The most likely scenario is that the image no longer exists in the Image Service
				// but the instance still has a record from when it existed.
				d.Set("image_name", "Image not found")
				return nil
			}
			return err
		} else {
			d.Set("image_name", image.Name)
		}
	}

	return nil
}

func getFlavorId(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	flavorId := d.Get("flavor_id").(string)

	if flavorId != "" {
		return flavorId, nil
	}

	flavorName := d.Get("flavor_name").(string)
	return flavors.IDFromName(client, flavorName)
}
