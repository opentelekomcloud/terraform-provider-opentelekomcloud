package opentelekomcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/floatingips"
)

func resourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCEClusterV3Create,
		Read:   resourceCCEClusterV3Read,
		Update: resourceCCEClusterV3Update,
		Delete: resourceCCEClusterV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				ForceNew: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"highway_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"extend_param": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"container_network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"container_network_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"authentication_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "x509",
			},
			"multi_az": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"eip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIP,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_otc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_authority_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"certificate_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_key_data": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceClusterLabelsV3(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceClusterAnnotationsV3(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceClusterExtendParamV3(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("extend_param").(map[string]interface{}) {
		m[key] = val.(string)
	}
	if multi_az, ok := d.GetOk("multi_az"); ok && multi_az == true {
		m["clusterAZ"] = "multi_az"
	}
	if eip, ok := d.GetOk("eip"); ok {
		m["clusterExternalIP"] = eip.(string)
	}
	return m
}

func resourceCCEClusterV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to create opentelekomcloud CCE client : %s", err)
	}
	if d.Get("eip").(string) != "" {
		fipId, err := resourceFloatingIPV2Exists(d, meta, d.Get("eip").(string))
		if err != nil {
			return fmt.Errorf("Error retrieving the eip: %s", err)
		}
		if fipId == "" {
			return fmt.Errorf("The specified EIP %s does not exist", d.Get("eip").(string))
		}
	}

	createOpts := clusters.CreateOpts{
		Kind:       "Cluster",
		ApiVersion: "v3",
		Metadata: clusters.CreateMetaData{Name: d.Get("name").(string),
			Labels:      resourceClusterLabelsV3(d),
			Annotations: resourceClusterAnnotationsV3(d)},
		Spec: clusters.Spec{
			Type:        d.Get("cluster_type").(string),
			Flavor:      d.Get("flavor_id").(string),
			Version:     d.Get("cluster_version").(string),
			Description: d.Get("description").(string),
			HostNetwork: clusters.HostNetworkSpec{VpcId: d.Get("vpc_id").(string),
				SubnetId:      d.Get("subnet_id").(string),
				HighwaySubnet: d.Get("highway_subnet_id").(string)},
			ContainerNetwork: clusters.ContainerNetworkSpec{Mode: d.Get("container_network_type").(string),
				Cidr: d.Get("container_network_cidr").(string)},
			Authentication: clusters.AuthenticationSpec{Mode: d.Get("authentication_mode").(string),
				AuthenticatingProxy: make(map[string]string)},
			BillingMode: d.Get("billing_mode").(int),
			ExtendParam: resourceClusterExtendParamV3(d),
		},
	}

	create, err := clusters.Create(cceClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for opentelekomcloud CCE cluster (%s) to become available", create.Metadata.Id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Creating"},
		Target:     []string{"Available"},
		Refresh:    waitForCCEClusterActive(cceClient, create.Metadata.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud CCE cluster: %s", err)
	}
	d.SetId(create.Metadata.Id)

	return resourceCCEClusterV3Read(d, meta)

}

func resourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud CCE client: %s", err)
	}

	n, err := clusters.Get(cceClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving opentelekomcloud CCE: %s", err)
	}

	d.Set("id", n.Metadata.Id)
	d.Set("name", n.Metadata.Name)
	d.Set("status", n.Status.Phase)
	d.Set("flavor_id", n.Spec.Flavor)
	d.Set("cluster_type", n.Spec.Type)
	d.Set("description", n.Spec.Description)
	d.Set("billing_mode", n.Spec.BillingMode)
	d.Set("vpc_id", n.Spec.HostNetwork.VpcId)
	d.Set("subnet_id", n.Spec.HostNetwork.SubnetId)
	d.Set("highway_subnet_id", n.Spec.HostNetwork.HighwaySubnet)
	d.Set("container_network_type", n.Spec.ContainerNetwork.Mode)
	d.Set("container_network_cidr", n.Spec.ContainerNetwork.Cidr)
	d.Set("authentication_mode", n.Spec.Authentication.Mode)
	d.Set("internal", n.Status.Endpoints[0].Internal)
	d.Set("external", n.Status.Endpoints[0].External)
	d.Set("external_otc", n.Status.Endpoints[0].ExternalOTC)
	d.Set("region", GetRegion(d, config))
	if n.Status.Endpoints[0].External != "" {
		eip := strings.Split(n.Status.Endpoints[0].External, "//")
		eip = strings.Split(eip[1], ":")
		d.Set("eip", eip[0])
	} else {
		d.Set("eip", "")
	}

	cert, err := clusters.GetCert(cceClient, d.Id()).Extract()
	if err != nil {
		log.Printf("Error retrieving opentelekomcloud CCE cluster cert: %s", err)
	}

	//Set Certificate Clusters
	var clusterList []map[string]interface{}
	for _, clusterObj := range cert.Clusters {
		clusterCert := make(map[string]interface{})
		clusterCert["name"] = clusterObj.Name
		clusterCert["server"] = clusterObj.Cluster.Server
		clusterCert["certificate_authority_data"] = clusterObj.Cluster.CertAuthorityData
		clusterList = append(clusterList, clusterCert)
	}
	d.Set("certificate_clusters", clusterList)

	//Set Certificate Users
	var userList []map[string]interface{}
	for _, userObj := range cert.Users {
		userCert := make(map[string]interface{})
		userCert["name"] = userObj.Name
		userCert["client_certificate_data"] = userObj.User.ClientCertData
		userCert["client_key_data"] = userObj.User.ClientKeyData
		userList = append(userList, userCert)
	}
	d.Set("certificate_users", userList)

	return nil
}

func resourceCCEClusterV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud CCE Client: %s", err)
	}

	var updateOpts clusters.UpdateOpts

	if d.HasChange("description") {
		updateOpts.Spec.Description = d.Get("description").(string)
		_, err = clusters.Update(cceClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating opentelekomcloud CCE: %s", err)
		}
	}

	if d.HasChange("eip") {
		oldEip, newEip := d.GetChange("eip")
		oldEipStr := oldEip.(string)
		newEipStr := newEip.(string)
		var fipId string
		if newEipStr != "" {
			fipId, err = resourceFloatingIPV2Exists(d, meta, newEipStr)
			if err != nil {
				return fmt.Errorf("Error retrieving the eip: %s", err)
			}
			if fipId == "" {
				return fmt.Errorf("The specified EIP %s does not exist", newEipStr)
			}
		}
		if oldEipStr != "" {
			updateIpOpts := clusters.UpdateIpOpts{
				Action: "unbind",
			}
			err = clusters.UpdateMasterIp(cceClient, d.Id(), updateIpOpts).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error unbinding EIP to opentelekomcloud CCE: %s", err)
			}
		}
		if newEipStr != "" {
			updateIpOpts := clusters.UpdateIpOpts{
				Action:    "bind",
				ElasticIp: newEipStr,
			}
			updateIpOpts.Spec.ID = fipId
			err = clusters.UpdateMasterIp(cceClient, d.Id(), updateIpOpts).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error binding EIP to opentelekomcloud CCE: %s", err)
			}
		}
	}

	return resourceCCEClusterV3Read(d, meta)
}

func resourceCCEClusterV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating opentelekomcloud CCE Client: %s", err)
	}
	err = clusters.Delete(cceClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting opentelekomcloud CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting", "Available", "Unavailable"},
		Target:     []string{"Deleted"},
		Refresh:    waitForCCEClusterDelete(cceClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error deleting opentelekomcloud CCE cluster: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCCEClusterActive(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := clusters.Get(cceClient, clusterId).Extract()
		if err != nil {
			return nil, "", err
		}

		return n, n.Status.Phase, nil
	}
}

func waitForCCEClusterDelete(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete  CCE cluster %s.\n", clusterId)

		r, err := clusters.Get(cceClient, clusterId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted opentelekomcloud CCE cluster %s", clusterId)
				return r, "Deleted", nil
			}
		}
		if r.Status.Phase == "Deleting" {
			return r, "Deleting", nil
		}
		log.Printf("[DEBUG] opentelekomcloud CCE cluster %s still available.\n", clusterId)
		return r, "Available", nil
	}
}

func resourceFloatingIPV2Exists(d *schema.ResourceData, meta interface{}, floatingIP string) (string, error) {
	config := meta.(*Config)
	networkClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating opentelekomcloud networking Client: %s", err)
	}
	listOpts := floatingips.ListOpts{
		FloatingIP: floatingIP,
	}
	allPages, err := floatingips.List(networkClient, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allFips, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return "", err
	}

	if len(allFips) == 0 {
		return "", nil
	}

	return allFips[0].ID, nil
}
