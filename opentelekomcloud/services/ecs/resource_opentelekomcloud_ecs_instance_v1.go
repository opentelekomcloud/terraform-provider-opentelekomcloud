package ecs

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/secgroups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/cloudservers"
	tags "github.com/opentelekomcloud/gophertelekomcloud/openstack/ecs/v1/cloudservertags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceEcsInstanceV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceEcsInstanceV1Create,
		Read:   resourceEcsInstanceV1Read,
		Update: resourceEcsInstanceV1Update,
		Delete: resourceEcsInstanceV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: common.MultipleCustomizeDiffs(
			common.ValidateVPC("vpc_id"),
			common.ValidateVolumeType("system_disk_type"),
			common.ValidateVolumeType("data_disks.*.type"),
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
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
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"nics": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 12,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"system_disk_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"system_disk_size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"data_disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 23,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"snapshot_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
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
			"tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ValidateFunc: common.ValidateECSTagValue,
			},
			"auto_recovery": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"delete_disks_on_termination": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceEcsInstanceV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV1 client: %s", err)
	}

	createOpts := &cloudservers.CreateOpts{
		Name:             d.Get("name").(string),
		ImageRef:         d.Get("image_id").(string),
		FlavorRef:        d.Get("flavor").(string),
		KeyName:          d.Get("key_name").(string),
		VpcId:            d.Get("vpc_id").(string),
		SecurityGroups:   resourceInstanceSecGroupsV1(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		Nics:             resourceInstanceNicsV1(d),
		RootVolume:       resourceInstanceRootVolumeV1(d),
		DataVolumes:      resourceInstanceDataVolumesV1(d),
		AdminPass:        d.Get("password").(string),
		UserData:         []byte(d.Get("user_data").(string)),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	n, err := cloudservers.Create(client, createOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud server: %s", err)
	}

	if err := cloudservers.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	entity, err := cloudservers.GetJobEntity(client, n.JobID, "server_id")
	if err != nil {
		return err
	}

	if id, ok := entity.(string); ok {
		d.SetId(id)

		if common.HasFilledOpt(d, "tags") {
			tagMap := d.Get("tags").(map[string]interface{})
			log.Printf("[DEBUG] Setting tags: %v", tagMap)
			err = SetTagForInstance(d, meta, id, tagMap)
			if err != nil {
				log.Printf("[WARN] Error setting tags of instance:%s, err=%s", id, err)
			}
		}

		if common.HasFilledOpt(d, "auto_recovery") {
			ar := d.Get("auto_recovery").(bool)
			log.Printf("[DEBUG] Set auto recovery of instance to %t", ar)
			err = setAutoRecoveryForInstance(d, meta, id, ar)
			if err != nil {
				log.Printf("[WARN] Error setting auto recovery of instance:%s, err=%s", id, err)
			}
		}

		return resourceEcsInstanceV1Read(d, meta)
	}

	return fmt.Errorf("unexpected conversion error in resourceEcsInstanceV1Create")
}

func resourceEcsInstanceV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV1 client: %s", err)
	}

	server, err := cloudservers.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "server")
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", d.Id(), server)

	mErr := multierror.Append(
		d.Set("name", server.Name),
		d.Set("image_id", server.Image.ID),
		d.Set("flavor", server.Flavor.ID),
		d.Set("password", d.Get("password")),
		d.Set("key_name", server.KeyName),
		d.Set("vpc_id", server.Metadata.VpcID),
		d.Set("availability_zone", server.AvailabilityZone),
	)
	var secGrpIDs []string
	for _, sg := range server.SecurityGroups {
		secGrpIDs = append(secGrpIDs, sg.ID)
	}
	mErr = multierror.Append(mErr,
		d.Set("security_groups", secGrpIDs),
	)

	// Get the instance network and address information
	nics := flattenInstanceNicsV1(d, meta, server.Addresses)
	mErr = multierror.Append(mErr,
		d.Set("nics", nics),
	)

	// Set instance tags
	tagList, err := tags.Get(client, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error fetching OpenTelekomCloud instance tags: %s", err)
	}

	tagMap := make(map[string]string)
	for _, val := range tagList.Tags {
		tagMap[val.Key] = val.Value
	}
	if err := d.Set("tags", tagMap); err != nil {
		return fmt.Errorf("[DEBUG] Error saving tag to state for OpenTelekomCloud instance (%s): %s", d.Id(), err)
	}

	ar, err := resourceECSAutoRecoveryV1Read(d, meta, d.Id())
	if err != nil && !common.IsResourceNotFound(err) {
		return fmt.Errorf("error reading auto recovery of instance:%s, err=%s", d.Id(), err)
	}
	mErr = multierror.Append(mErr,
		d.Set("auto_recovery", ar),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting ECS attrbutes: %s", err)
	}

	return nil
}

func resourceEcsInstanceV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	var updateOpts servers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if updateOpts != (servers.UpdateOpts{}) {
		_, err := servers.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating OpenTelekomCloud server: %s", err)
		}
	}

	if d.HasChange("security_groups") {
		oldSGRaw, newSGRaw := d.GetChange("security_groups")
		oldSGSet := oldSGRaw.(*schema.Set)
		newSGSet := newSGRaw.(*schema.Set)
		secGroupsToAdd := newSGSet.Difference(oldSGSet)
		secGroupsToRemove := oldSGSet.Difference(newSGSet)

		log.Printf("[DEBUG] Security groups to add: %v", secGroupsToAdd)

		log.Printf("[DEBUG] Security groups to remove: %v", secGroupsToRemove)

		for _, sg := range secGroupsToRemove.List() {
			err := secgroups.RemoveServer(client, d.Id(), sg.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					continue
				}

				return fmt.Errorf("error removing security group (%s) from OpenTelekomCloud server (%s): %s", sg, d.Id(), err)
			} else {
				log.Printf("[DEBUG] Removed security group (%s) from instance (%s)", sg, d.Id())
			}
		}

		for _, sg := range secGroupsToAdd.List() {
			err := secgroups.AddServer(client, d.Id(), sg.(string)).ExtractErr()
			if err != nil && err.Error() != "EOF" {
				return fmt.Errorf("error adding security group (%s) to OpenTelekomCloud server (%s): %s", sg, d.Id(), err)
			}
			log.Printf("[DEBUG] Added security group (%s) to instance (%s)", sg, d.Id())
		}
	}

	if d.HasChange("flavor") {
		newFlavorId := d.Get("flavor").(string)

		resizeOpts := &servers.ResizeOpts{
			FlavorRef: newFlavorId,
		}
		log.Printf("[DEBUG] Resize configuration: %#v", resizeOpts)
		err := servers.Resize(client, d.Id(), resizeOpts).ExtractErr()
		if err != nil {
			return fmt.Errorf("error resizing OpenTelekomCloud server: %s", err)
		}

		// Wait for the instance to finish resizing.
		log.Printf("[DEBUG] Waiting for instance (%s) to finish resizing", d.Id())

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"RESIZE"},
			Target:     []string{"VERIFY_RESIZE"},
			Refresh:    ServerV2StateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error waiting for instance (%s) to resize: %s", d.Id(), err)
		}

		// Confirm resize.
		log.Printf("[DEBUG] Confirming resize")
		err = servers.ConfirmResize(client, d.Id()).ExtractErr()
		if err != nil {
			return fmt.Errorf("error confirming resize of OpenTelekomCloud server: %s", err)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"VERIFY_RESIZE"},
			Target:     []string{"ACTIVE"},
			Refresh:    ServerV2StateRefreshFunc(client, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("error waiting for instance (%s) to confirm resize: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		computeClient, err := config.ComputeV1Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud compute v1 client: %s", err)
		}
		oldTags, err := tags.Get(computeClient, d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("error fetching OpenTelekomCloud instance tags: %s", err)
		}
		if len(oldTags.Tags) > 0 {
			deleteOpts := tags.BatchOpts{Action: tags.ActionDelete, Tags: oldTags.Tags}
			deleteTags := tags.BatchAction(computeClient, d.Id(), deleteOpts)
			if deleteTags.Err != nil {
				return fmt.Errorf("error updating OpenTelekomCloud instance tags: %s", deleteTags.Err)
			}
		}

		if common.HasFilledOpt(d, "tags") {
			tagMap := d.Get("tags").(map[string]interface{})
			if len(tagMap) > 0 {
				log.Printf("[DEBUG] Setting tags: %v", tagMap)
				err = SetTagForInstance(d, meta, d.Id(), tagMap)
				if err != nil {
					return fmt.Errorf("error updating tags of instance:%s, err:%s", d.Id(), err)
				}
			}
		}
	}

	if d.HasChange("auto_recovery") {
		ar := d.Get("auto_recovery").(bool)
		log.Printf("[DEBUG] Update auto recovery of instance to %t", ar)
		err = setAutoRecoveryForInstance(d, meta, d.Id(), ar)
		if err != nil {
			return fmt.Errorf("error updating auto recovery of instance:%s, err:%s", d.Id(), err)
		}
	}

	return resourceEcsInstanceV1Read(d, meta)
}

func resourceEcsInstanceV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.ComputeV1Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	var serverRequests []cloudservers.Server
	server := cloudservers.Server{
		Id: d.Id(),
	}
	serverRequests = append(serverRequests, server)

	deleteOpts := cloudservers.DeleteOpts{
		Servers:      serverRequests,
		DeleteVolume: d.Get("delete_disks_on_termination").(bool),
	}

	n, err := cloudservers.Delete(client, deleteOpts).ExtractJobResponse()
	if err != nil {
		return fmt.Errorf("error deleting OpenTelekomCloud server: %s", err)
	}

	if err := cloudservers.WaitForJobSuccess(client, int(d.Timeout(schema.TimeoutCreate)/time.Second), n.JobID); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func resourceInstanceNicsV1(d *schema.ResourceData) []cloudservers.Nic {
	var nicRequests []cloudservers.Nic

	nics := d.Get("nics").([]interface{})
	for i := range nics {
		nic := nics[i].(map[string]interface{})
		nicRequest := cloudservers.Nic{
			SubnetId:  nic["network_id"].(string),
			IpAddress: nic["ip_address"].(string),
		}

		nicRequests = append(nicRequests, nicRequest)
	}
	return nicRequests
}

func resourceInstanceRootVolumeV1(d *schema.ResourceData) cloudservers.RootVolume {
	diskType := d.Get("system_disk_type").(string)
	if diskType == "" {
		diskType = "SATA"
	}
	volRequest := cloudservers.RootVolume{
		VolumeType: diskType,
		Size:       d.Get("system_disk_size").(int),
	}
	return volRequest
}

func resourceInstanceDataVolumesV1(d *schema.ResourceData) []cloudservers.DataVolume {
	var volRequests []cloudservers.DataVolume

	vols := d.Get("data_disks").([]interface{})
	for i := range vols {
		vol := vols[i].(map[string]interface{})
		volRequest := cloudservers.DataVolume{
			VolumeType: vol["type"].(string),
			Size:       vol["size"].(int),
		}
		if vol["snapshot_id"] != "" {
			extendParam := cloudservers.VolumeExtendParam{
				SnapshotId: vol["snapshot_id"].(string),
			}
			volRequest.Extendparam = &extendParam
		}

		volRequests = append(volRequests, volRequest)
	}
	return volRequests
}

func resourceInstanceSecGroupsV1(d *schema.ResourceData) []cloudservers.SecurityGroup {
	rawSecGroups := d.Get("security_groups").(*schema.Set).List()
	secGroups := make([]cloudservers.SecurityGroup, len(rawSecGroups))
	for i, raw := range rawSecGroups {
		secGroups[i] = cloudservers.SecurityGroup{
			ID: raw.(string),
		}
	}
	return secGroups
}

func flattenInstanceNicsV1(
	d *schema.ResourceData, meta interface{}, addresses map[string][]cloudservers.Address) []map[string]interface{} {

	config := meta.(*cfg.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		log.Printf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	var network string
	var nics []map[string]interface{}
	// Loop through all networks and addresses.
	for _, addrs := range addresses {
		for _, addr := range addrs {
			// Skip if not fixed ip
			if addr.Type != "fixed" {
				continue
			}

			p, err := ports.Get(networkingClient, addr.PortID).Extract()
			if err != nil {
				network = ""
				log.Printf("[DEBUG] flattenInstanceNicsV1: failed to fetch port %s", addr.PortID)
			} else {
				network = p.NetworkID
			}

			v := map[string]interface{}{
				"network_id":  network,
				"ip_address":  addr.Addr,
				"mac_address": addr.MacAddr,
			}
			nics = append(nics, v)
		}
	}

	log.Printf("[DEBUG] flattenInstanceNicsV1: %#v", nics)
	return nics
}
