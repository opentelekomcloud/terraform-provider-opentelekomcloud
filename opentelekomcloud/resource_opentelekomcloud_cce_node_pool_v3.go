package opentelekomcloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodepools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodes"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
)

var (
	// Cluster pool taint key and value is 1 to 63 characters starting with a letter or digit.
	// Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed.
	clusterPoolTaintRegex, _ = regexp.Compile("^[a-zA-Z0-9_.-]{1,63}$")
)

func resourceCCENodePoolV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodePoolV3Create,
		Read:   resourceCCENodePoolV3Read,
		Update: resourceCCENodePoolV3Update,
		Delete: resourceCCENodePoolV3Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: multipleCustomizeDiffs(
			validateVolumeType("root_volume.*.volumetype"),
			validateVolumeType("data_volumes.*.volumetype"),
			validateSubnet("subnet_id"),
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "random",
			},
			"os": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"root_volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(10, 32768),
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
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
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(100, 32768),
						},
						"volumetype": {
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": {
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"initial_node_count": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"k8s_tags": {
				Type:         schema.TypeMap,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateK8sTagsMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"user_tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"taints": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(clusterPoolTaintRegex, "Invalid key. "+
								"Cluster pool taint key is 1 to 63 characters starting with a letter or digit. "+
								"Only lowercase letters, digits, and hyphens (-) are allowed."),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(clusterPoolTaintRegex, "Invalid value. "+
								"Cluster pool taint value is 1 to 63 characters starting with a letter or digit. "+
								"Only letters, digits, hyphens (-), underscores (_), and periods (.) are allowed."),
						},
						"effect": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"NoSchedule", "PreferNoSchedule", "NoExecute",
							}, false),
						},
					}},
			},
			"key_pair": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"password": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"password", "key_pair"},
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"preinstall": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: getHashOrEmpty,
			},
			"postinstall": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				StateFunc: getHashOrEmpty,
			},
			"scale_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"min_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_node_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scale_down_cooldown_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"server_group_reference": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCCENodePoolUserTags(d *schema.ResourceData) []tags.ResourceTag {
	tagRaw := d.Get("user_tags").(map[string]interface{})
	return expandResourceTags(tagRaw)
}

func resourceCCENodePoolV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud CCE Node Pool client: %s", err)
	}

	var base64PreInstall, base64PostInstall string
	if v, ok := d.GetOk("preinstall"); ok {
		base64PreInstall = installScriptEncode(v.(string))
	}
	if v, ok := d.GetOk("postinstall"); ok {
		base64PostInstall = installScriptEncode(v.(string))
	}
	var loginSpec nodes.LoginSpec
	if hasFilledOpt(d, "key_pair") {
		loginSpec = nodes.LoginSpec{SshKey: d.Get("key_pair").(string)}
	}
	if hasFilledOpt(d, "password") {
		loginSpec = nodes.LoginSpec{
			UserPassword: nodes.UserPassword{
				Username: "root",
				Password: d.Get("password").(string),
			},
		}
	}

	createOpts := nodepools.CreateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.CreateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.CreateSpec{
			InitialNodeCount: d.Get("initial_node_count").(int),
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scale_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
			NodeManagement: nodepools.NodeManagementSpec{
				ServerGroupReference: d.Get("server_group_reference").(string),
			},
			NodeTemplate: nodes.Spec{
				Flavor:      d.Get("flavor").(string),
				Az:          d.Get("availability_zone").(string),
				Os:          d.Get("os").(string),
				Login:       loginSpec,
				RootVolume:  resourceCCERootVolume(d),
				DataVolumes: resourceCCEDataVolume(d),
				BillingMode: 0,
				Count:       1,
				NodeNicSpec: nodes.NodeNicSpec{
					PrimaryNic: nodes.PrimaryNic{
						SubnetId: d.Get("subnet_id").(string),
					},
				},
				ExtendParam: nodes.ExtendParam{
					PreInstall:  base64PreInstall,
					PostInstall: base64PostInstall,
				},
				Taints:   resourceCCENodeTaints(d),
				K8sTags:  resourceCCENodeK8sTags(d),
				UserTags: resourceCCENodePoolUserTags(d),
			},
		},
	}

	clusterId := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Target:     []string{"Available"},
		Refresh:    waitForClusterAvailable(nodePoolClient, clusterId),
		Timeout:    15 * time.Minute,
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateCluster.WaitForState()

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	s, err := nodepools.Create(nodePoolClient, clusterId, createOpts).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault403); ok {
			retryNode, err := recursiveNodePoolCreate(nodePoolClient, createOpts, clusterId, 403)
			if err == "fail" {
				return fmt.Errorf("error creating Open Telekom Cloud CCE Node Pool")
			}
			s = retryNode
		} else {
			return fmt.Errorf("error creating Open Telekom Cloud CCE Node Pool: %s", err)
		}
	}

	if len(s.Metadata.Id) == 0 {
		return fmt.Errorf("error fetching CreateNodePool id")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Synchronizing"},
		Target:       []string{""},
		Refresh:      waitForCceNodePoolActive(nodePoolClient, clusterId, s.Metadata.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error creating Open Telekom Cloud CCE Node Pool: %s", err)
	}

	d.SetId(s.Metadata.Id)
	return resourceCCENodePoolV3Read(d, meta)
}

func resourceCCENodePoolV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Open Telekom Cloud CCE Node Pool client: %s", err)
	}
	clusterId := d.Get("cluster_id").(string)
	s, err := nodepools.Get(nodePoolClient, clusterId, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("error retrieving Open Telekom Cloud CCE Node Pool: %s", err)
	}

	me := multierror.Append(nil,
		d.Set("name", s.Metadata.Name),
		d.Set("flavor", s.Spec.NodeTemplate.Flavor),
		d.Set("availability_zone", s.Spec.NodeTemplate.Az),
		d.Set("os", s.Spec.NodeTemplate.Os),
		d.Set("key_pair", s.Spec.NodeTemplate.Login.SshKey),
		d.Set("initial_node_count", s.Spec.InitialNodeCount),
		d.Set("scale_enable", s.Spec.Autoscaling.Enable),
		d.Set("min_node_count", s.Spec.Autoscaling.MinNodeCount),
		d.Set("max_node_count", s.Spec.Autoscaling.MaxNodeCount),
		d.Set("scale_down_cooldown_time", s.Spec.Autoscaling.ScaleDownCooldownTime),
		d.Set("priority", s.Spec.Autoscaling.Priority),
	)
	if err := me.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting CCE Node Pool attributes (%s): %s", d.Id(), err)
	}

	k8sTags := map[string]string{}
	for key, val := range s.Spec.NodeTemplate.K8sTags {
		if strings.Contains(key, "cce.cloud.com") {
			continue
		}
		k8sTags[key] = val
	}
	if err := d.Set("k8s_tags", k8sTags); err != nil {
		return fmt.Errorf("[DEBUG] Error saving k8s_tags to state for Open Telekom Cloud CCE Node Pool (%s): %s", d.Id(), err)
	}

	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.NodeTemplate.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_param"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving dataVolumes to state for Open Telekom Cloud CCE Node Pool (%s): %s", d.Id(), err)
	}

	rootVolume := []map[string]interface{}{
		{
			"size":         s.Spec.NodeTemplate.RootVolume.Size,
			"volumetype":   s.Spec.NodeTemplate.RootVolume.VolumeType,
			"extend_param": s.Spec.NodeTemplate.RootVolume.ExtendParam,
		},
	}
	if err := d.Set("root_volume", rootVolume); err != nil {
		return fmt.Errorf("[DEBUG] Error saving rootVolume to state for Open Telekom Cloud CCE Node Pool (%s): %s", d.Id(), err)
	}

	if err := d.Set("status", s.Status.Phase); err != nil {
		return fmt.Errorf("[DEBUG] Error saving status to state for Open Telekom Cloud CCE Node Pool (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceCCENodePoolV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Open Telekom Cloud CCE client: %s", err)
	}
	updateOpts := nodepools.UpdateOpts{
		Kind:       "NodePool",
		ApiVersion: "v3",
		Metadata: nodepools.UpdateMetaData{
			Name: d.Get("name").(string),
		},
		Spec: nodepools.UpdateSpec{
			InitialNodeCount: d.Get("initial_node_count").(int),
			Autoscaling: nodepools.AutoscalingSpec{
				Enable:                d.Get("scale_enable").(bool),
				MinNodeCount:          d.Get("min_node_count").(int),
				MaxNodeCount:          d.Get("max_node_count").(int),
				ScaleDownCooldownTime: d.Get("scale_down_cooldown_time").(int),
				Priority:              d.Get("priority").(int),
			},
		},
	}
	clusterId := d.Get("cluster_id").(string)
	_, err = nodepools.Update(nodePoolClient, clusterId, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating Open Telekom Cloud CCE Node Pool: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Synchronizing"},
		Target:     []string{""},
		Refresh:    waitForCceNodePoolActive(nodePoolClient, clusterId, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error creating Open Telekom Cloud CCE Node Pool: %s", err)
	}

	return resourceCCENodePoolV3Read(d, meta)
}
func resourceCCENodePoolV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodePoolClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating Open Telekom Cloud CCE client: %s", err)
	}
	clusterId := d.Get("cluster_id").(string)
	err = nodepools.Delete(nodePoolClient, clusterId, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting Open Telekom Cloud CCE Node Pool: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Deleting"},
		Target:       []string{"Deleted"},
		Refresh:      waitForCceNodePoolDelete(nodePoolClient, clusterId, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error deleting Open Telekom Cloud CCE Node Pool: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodePoolActive(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()
		if err != nil {
			return nil, "", err
		}
		return n, n.Status.Phase, nil
	}
}

func waitForCceNodePoolDelete(cceClient *golangsdk.ServiceClient, clusterId, nodePoolId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete Open Telekom Cloud CCE Node Pool %s.\n", nodePoolId)

		r, err := nodepools.Get(cceClient, clusterId, nodePoolId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted Open Telekom Cloud CCE Node Pool %s", nodePoolId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		log.Printf("[DEBUG] Open Telekom Cloud Node Pool %s still available.\n", nodePoolId)
		return r, r.Status.Phase, nil
	}
}

func recursiveNodePoolCreate(cceClient *golangsdk.ServiceClient, opts nodepools.CreateOptsBuilder, ClusterID string, errCode int) (*nodepools.NodePool, string) {
	if errCode == 403 {
		stateCluster := &resource.StateChangeConf{
			Target:     []string{"Available"},
			Refresh:    waitForClusterAvailable(cceClient, ClusterID),
			Timeout:    15 * time.Minute,
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, stateErr := stateCluster.WaitForState()
		if stateErr != nil {
			log.Printf("[INFO] Cluster Unavailable %s.\n", stateErr)
		}
		s, err := nodepools.Create(cceClient, ClusterID, opts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault403); ok {
				return recursiveNodePoolCreate(cceClient, opts, ClusterID, 403)
			} else {
				return s, "fail"
			}
		} else {
			return s, "success"
		}
	}
	return nil, "fail"
}
