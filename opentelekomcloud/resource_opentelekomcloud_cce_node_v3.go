package opentelekomcloud

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
)

func resourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodeV3Create,
		Read:   resourceCCENodeV3Read,
		Update: resourceCCENodeV3Update,
		Delete: resourceCCENodeV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:          schema.TypeMap,
				ConflictsWith: []string{"tags"},
				Optional:      true,
				ForceNew:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					}},
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					}},
			},
			"eip_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
				ConflictsWith: []string{
					"iptype", "bandwidth_charge_mode", "bandwidth_size", "sharetype",
				},
			},
			"eip_count": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"eip_ids"},
			},
			"iptype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5_bgp",
				}, true),
				ConflictsWith: []string{"eip_ids"},
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype", "bandwidth_charge_mode",
				},
			},
			"bandwidth_charge_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"traffic",
				}, true),
				ConflictsWith: []string{"eip_ids"},
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype", "bandwidth_charge_mode",
				},
			},
			"sharetype": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PER",
				}, true),
				ConflictsWith: []string{"eip_ids"},
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype", "bandwidth_charge_mode",
				},
			},
			"bandwidth_size": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"eip_ids"},
				RequiredWith: []string{
					"iptype", "bandwidth_size", "sharetype", "bandwidth_charge_mode",
				},
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"extend_param_charging_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ecs_performance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"max_pods": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"preinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"postinstall": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch v.(type) {
					case string:
						return installScriptHashSum(v.(string))
					default:
						return ""
					}
				},
			},
			"tags": {
				Type:          schema.TypeMap,
				ConflictsWith: []string{"labels"},
				Optional:      true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCCENodeLabelsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceCCENodeAnnotationsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func resourceCCENodeTags(d *schema.ResourceData) []tags.ResourceTag {
	tagRaw := d.Get("tags").(map[string]interface{})
	return expandResourceTags(tagRaw)
}

func resourceCCEDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	volumeRaw := d.Get("data_volumes").([]interface{})
	volumes := make([]nodes.VolumeSpec, len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:        rawMap["size"].(int),
			VolumeType:  rawMap["volumetype"].(string),
			ExtendParam: rawMap["extend_param"].(string),
		}
	}
	return volumes
}

func resourceCCERootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var nics nodes.VolumeSpec
	nicsRaw := d.Get("root_volume").([]interface{})
	if len(nicsRaw) == 1 {
		nics.Size = nicsRaw[0].(map[string]interface{})["size"].(int)
		nics.VolumeType = nicsRaw[0].(map[string]interface{})["volumetype"].(string)
		nics.ExtendParam = nicsRaw[0].(map[string]interface{})["extend_param"].(string)
	}
	return nics
}

func resourceCCEEipIDs(d *schema.ResourceData) []string {
	rawID := d.Get("eip_ids").(*schema.Set)
	id := make([]string, rawID.Len())
	for i, raw := range rawID.List() {
		id[i] = raw.(string)
	}
	return id
}

func resourceCCENodeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE Node client: %s", err)
	}

	var base64PreInstall, base64PostInstall string
	if v, ok := d.GetOk("preinstall"); ok {
		base64PreInstall = installScriptEncode(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		base64PostInstall = installScriptEncode(v.(string))
	}

	// eip_count and bandwidth_size parameters must be set simultaneously
	bandwidthSize := d.Get("bandwidth_size").(int)
	eipCount := d.Get("eip_count").(int)
	if bandwidthSize > 0 && eipCount == 0 {
		eipCount = 1
	}

	createOpts := nodes.CreateOpts{
		Kind:       "Node",
		ApiVersion: "v3",
		Metadata: nodes.CreateMetaData{
			Name:        d.Get("name").(string),
			Labels:      resourceCCENodeLabelsV2(d),
			Annotations: resourceCCENodeAnnotationsV2(d),
		},
		Spec: nodes.Spec{
			Flavor:      d.Get("flavor_id").(string),
			Az:          d.Get("availability_zone").(string),
			Login:       nodes.LoginSpec{SshKey: d.Get("key_pair").(string)},
			RootVolume:  resourceCCERootVolume(d),
			DataVolumes: resourceCCEDataVolume(d),
			PublicIP: nodes.PublicIPSpec{
				Ids:   resourceCCEEipIDs(d),
				Count: eipCount,
				Eip: nodes.EipSpec{
					IpType: d.Get("iptype").(string),
					Bandwidth: nodes.BandwidthOpts{
						ChargeMode: d.Get("bandwidth_charge_mode").(string),
						Size:       d.Get("bandwidth_size").(int),
						ShareType:  d.Get("sharetype").(string),
					},
				},
			},
			BillingMode: d.Get("billing_mode").(int),
			Count:       1,
			ExtendParam: nodes.ExtendParam{
				ChargingMode:       d.Get("extend_param_charging_mode").(int),
				EcsPerformanceType: d.Get("ecs_performance_type").(string),
				MaxPods:            d.Get("max_pods").(int),
				OrderID:            d.Get("order_id").(string),
				ProductID:          d.Get("product_id").(string),
				PublicKey:          d.Get("public_key").(string),
				PreInstall:         base64PreInstall,
				PostInstall:        base64PostInstall,
			},
			UserTags: resourceCCENodeTags(d),
		},
	}

	clusterId := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:     []string{"Available"},
		Refresh:    waitForClusterAvailable(nodeClient, clusterId),
		Timeout:    15 * time.Minute,
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	s, err := nodes.Create(nodeClient, clusterId, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveCreate(nodeClient, createOpts, clusterId, 403)
			if err == "fail" {
				return fmt.Errorf("Error creating OpenTelekomCloud Node")
			}
			s = retryNode
		} else {
			return fmt.Errorf("Error creating OpenTelekomCloud Node: %s", err)
		}
	}

	job, err := nodes.GetJobDetails(nodeClient, s.Status.JobID).ExtractJob()
	if err != nil {
		return fmt.Errorf("Error fetching OpenTelekomCloud Job Details: %s", err)
	}
	jobResourceId := job.Spec.SubJobs[0].Metadata.ID

	subJob, err := nodes.GetJobDetails(nodeClient, jobResourceId).ExtractJob()
	if err != nil {
		return fmt.Errorf("Error fetching OpenTelekomCloud Job Details: %s", err)
	}

	var nodeId string
	for _, s := range subJob.Spec.SubJobs {
		if s.Spec.Type == "CreateNodeVM" {
			nodeId = s.Spec.ResourceID
			break
		}
	}

	log.Printf("[DEBUG] Waiting for CCE Node (%s) to become available", s.Metadata.Name)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Build", "Installing"},
		Target:     []string{"Active"},
		Refresh:    waitForCceNodeActive(nodeClient, clusterId, nodeId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE Node: %s", err)
	}

	d.SetId(nodeId)
	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE Node client: %s", err)
	}
	clusterId := d.Get("cluster_id").(string)
	s, err := nodes.Get(nodeClient, clusterId, d.Id()).Extract()

	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving OpenTelekomCloud Node: %s", err)
	}

	d.Set("region", GetRegion(d, config))
	d.Set("name", s.Metadata.Name)
	d.Set("flavor_id", s.Spec.Flavor)
	d.Set("availability_zone", s.Spec.Az)
	d.Set("billing_mode", s.Spec.BillingMode)
	d.Set("extend_param_charging_mode", s.Spec.ExtendParam.ChargingMode)
	d.Set("ecs:performance_type", s.Spec.ExtendParam.PublicKey)
	d.Set("order_id", s.Spec.ExtendParam.OrderID)
	d.Set("product_id", s.Spec.ExtendParam.ProductID)
	d.Set("max_pods", s.Spec.ExtendParam.MaxPods)
	d.Set("ecs_performance_type", s.Spec.ExtendParam.EcsPerformanceType)
	d.Set("key_pair", s.Spec.Login.SshKey)

	// Spec.PublicIP field is empty in the response body even if eip was configured,
	// so we should not set the following attributes
	/*
		// set PublicIPSpec
		d.Set("eip_ids", s.Spec.PublicIP.Ids)
		d.Set("iptype", s.Spec.PublicIP.Eip.IpType)
		d.Set("bandwidth_charge_mode", s.Spec.PublicIP.Eip.Bandwidth.ChargeMode)
		d.Set("bandwidth_size", s.Spec.PublicIP.Eip.Bandwidth.Size)
		d.Set("sharetype", s.Spec.PublicIP.Eip.Bandwidth.ShareType)
	*/

	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_param"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving dataVolumes to state for OpenTelekomCloud Node (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":         s.Spec.RootVolume.Size,
			"volumetype":   s.Spec.RootVolume.VolumeType,
			"extend_param": s.Spec.RootVolume.ExtendParam,
		},
	}
	d.Set("root_volume", rootVolume)
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmt.Errorf("[DEBUG] Error saving root Volume to state for OpenTelekomCloud Node (%s): %s", d.Id(), err)
	}

	// set computed attributes
	serverId := s.Status.ServerID
	d.Set("server_id", serverId)
	d.Set("private_ip", s.Status.PrivateIP)
	d.Set("public_ip", s.Status.PublicIP)
	d.Set("status", s.Status.Phase)

	// fetch tags from ECS instance
	computeClient, err := config.computeV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
	}

	resourceTags, err := tags.Get(computeClient, "cloudservers", serverId).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching OpenTelekomCloud instance tags: %s", err)
	}

	tagMap := tagsToMap(resourceTags.Tags)
	// ignore "CCE-Dynamic-Provisioning-Node"
	delete(tagMap, "CCE-Dynamic-Provisioning-Node")
	if err := d.Set("tags", tagMap); err != nil {
		return fmt.Errorf("Error saving tags of cce node: %s", err)
	}

	return nil
}

func resourceCCENodeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE client: %s", err)
	}

	var updateOpts nodes.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Metadata.Name = d.Get("name").(string)

		clusterId := d.Get("cluster_id").(string)
		_, err = nodes.Update(nodeClient, clusterId, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenTelekomCloud cce node: %s", err)
		}
	}

	//update tags
	if d.HasChange("tags") {
		computeClient, err := config.computeV1Client(GetRegion(d, config))
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud compute client: %s", err)
		}

		serverId := d.Get("server_id").(string)
		tagErr := UpdateResourceTags(computeClient, d, "cloudservers", serverId)
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of CCE node %s: %s", d.Id(), tagErr)
		}
	}

	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE client: %s", err)
	}
	clusterId := d.Get("cluster_id").(string)
	err = nodes.Delete(nodeClient, clusterId, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"Deleted"},
		Refresh:    waitForCceNodeDelete(nodeClient, clusterId, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud CCE Node: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodeActive(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodes.Get(cceClient, clusterId, nodeId).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func waitForCceNodeDelete(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud CCE Node %s.\n", nodeId)

		r, err := nodes.Get(cceClient, clusterId, nodeId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud CCE Node %s", nodeId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud CCE Node %s still available.\n", nodeId)
		return r, r.Status.Phase, nil
	}
}

func waitForClusterAvailable(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[INFO] Waiting for OpenTelekomCloud Cluster to be available %s.\n", clusterId)
		n, err := clusters.Get(cceClient, clusterId).Extract()

		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func recursiveCreate(cceClient *golangsdk.ServiceClient, opts nodes.CreateOptsBuilder, ClusterID string, errCode int) (*nodes.Nodes, string) {
	if errCode == 403 {
		stateCluster := &resource.StateChangeConf{
			Target:     []string{"Available"},
			Refresh:    waitForClusterAvailable(cceClient, ClusterID),
			Timeout:    15 * time.Minute,
			Delay:      15 * time.Second,
			MinTimeout: 10 * time.Second,
		}
		_, stateErr := stateCluster.WaitForState()
		if stateErr != nil {
			log.Printf("[INFO] Cluster Unavailable %s.\n", stateErr)
		}
		s, err := nodes.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			//if err.(golangsdk.ErrUnexpectedResponseCode).Actual == 403 {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				return recursiveCreate(cceClient, opts, ClusterID, 403)
			} else {
				return s, "fail"
			}
		} else {
			return s, "success"
		}
	}
	return nil, "fail"
}

func installScriptHashSum(script string) string {
	// Check whether the preinstall/postinstall is not Base64 encoded.
	// Always calculate hash of base64 decoded value since we
	// check against double-encoding when setting it
	v, base64DecodeError := base64.StdEncoding.DecodeString(script)
	if base64DecodeError != nil {
		v = []byte(script)
	}

	hash := sha1.Sum(v)
	return hex.EncodeToString(hash[:])
}

func installScriptEncode(script string) string {
	if _, err := base64.StdEncoding.DecodeString(script); err != nil {
		return base64.StdEncoding.EncodeToString([]byte(script))
	}
	return script
}
