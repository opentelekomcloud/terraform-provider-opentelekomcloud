package opentelekomcloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/clusters"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"
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
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressSmartVersionDiff,
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
			"authenticating_proxy_ca": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressEquivalentBase64,
			},
			"kubernetes_svc_ip_range": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"kube_proxy_mode": { // can't be set via API currently
				Type:     schema.TypeString,
				Computed: true,
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
	if multiAZ, ok := d.GetOk("multi_az"); ok && multiAZ == true {
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
		return fmt.Errorf("unable to create opentelekomcloud CCE client : %s", err)
	}
	if d.Get("eip").(string) != "" {
		fipId, err := resourceFloatingIPV2Exists(d, meta, d.Get("eip").(string))
		if err != nil {
			return fmt.Errorf("error retrieving the eip: %s", err)
		}
		if fipId == "" {
			return fmt.Errorf("the specified EIP %s does not exist", d.Get("eip").(string))
		}
	}

	authProxy := make(map[string]string)
	if ca, ok := d.GetOk("authenticating_proxy_ca"); ok {
		authProxy["ca"] = base64IfNot(ca.(string))
	}

	createOpts := clusters.CreateOpts{
		Kind:       "Cluster",
		ApiVersion: "v3",
		Metadata: clusters.CreateMetaData{
			Name:        d.Get("name").(string),
			Labels:      resourceClusterLabelsV3(d),
			Annotations: resourceClusterAnnotationsV3(d),
		},
		Spec: clusters.Spec{
			Type:        d.Get("cluster_type").(string),
			Flavor:      d.Get("flavor_id").(string),
			Version:     d.Get("cluster_version").(string),
			Description: d.Get("description").(string),
			HostNetwork: clusters.HostNetworkSpec{
				VpcId:         d.Get("vpc_id").(string),
				SubnetId:      d.Get("subnet_id").(string),
				HighwaySubnet: d.Get("highway_subnet_id").(string),
			},
			ContainerNetwork: clusters.ContainerNetworkSpec{
				Mode: d.Get("container_network_type").(string),
				Cidr: d.Get("container_network_cidr").(string),
			},
			Authentication: clusters.AuthenticationSpec{
				Mode:                d.Get("authentication_mode").(string),
				AuthenticatingProxy: authProxy,
			},
			BillingMode:          d.Get("billing_mode").(int),
			ExtendParam:          resourceClusterExtendParamV3(d),
			KubernetesSvcIpRange: d.Get("kubernetes_svc_ip_range").(string),
			KubeProxyMode:        d.Get("kube_proxy_mode").(string),
		},
	}

	create, err := clusters.Create(cceClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud Cluster: %s", err)
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
		return fmt.Errorf("error creating OpenTelekomCloud CCE cluster: %s", err)
	}
	d.SetId(create.Metadata.Id)

	return resourceCCEClusterV3Read(d, meta)
}

func resourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud CCE client: %s", err)
	}

	cluster, err := clusters.Get(cceClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error retrieving opentelekomcloud CCE: %s", err)
	}

	authProxyCA, ok := cluster.Spec.Authentication.AuthenticatingProxy["ca"]
	if !ok {
		return fmt.Errorf("error reading authenticating proxy CA property")
	}
	b64decodedCA, err := base64.StdEncoding.DecodeString(authProxyCA)
	if err != nil {
		return fmt.Errorf("error decoding auth proxy CA: %s", err)
	}
	authProxyCA = string(b64decodedCA)

	eip := ""
	if cluster.Status.Endpoints[0].External != "" {
		endpointURL, err := url.Parse(cluster.Status.Endpoints[0].External)
		if err != nil {
			return fmt.Errorf("error parsing endpoint URL: %s", err)
		}
		eip = endpointURL.Hostname()
	}

	mErr := multierror.Append(nil,
		d.Set("name", cluster.Metadata.Name),
		d.Set("status", cluster.Status.Phase),
		d.Set("flavor_id", cluster.Spec.Flavor),
		d.Set("cluster_type", cluster.Spec.Type),
		d.Set("cluster_version", cluster.Spec.Version),
		d.Set("description", cluster.Spec.Description),
		d.Set("billing_mode", cluster.Spec.BillingMode),
		d.Set("vpc_id", cluster.Spec.HostNetwork.VpcId),
		d.Set("subnet_id", cluster.Spec.HostNetwork.SubnetId),
		d.Set("highway_subnet_id", cluster.Spec.HostNetwork.HighwaySubnet),
		d.Set("container_network_type", cluster.Spec.ContainerNetwork.Mode),
		d.Set("container_network_cidr", cluster.Spec.ContainerNetwork.Cidr),
		d.Set("authentication_mode", cluster.Spec.Authentication.Mode),
		d.Set("authenticating_proxy_ca", authProxyCA),
		d.Set("kubernetes_svc_ip_range", cluster.Spec.KubernetesSvcIpRange),
		d.Set("kube_proxy_mode", cluster.Spec.KubeProxyMode),
		d.Set("internal", cluster.Status.Endpoints[0].Internal),
		d.Set("external", cluster.Status.Endpoints[0].External),
		d.Set("external_otc", cluster.Status.Endpoints[0].ExternalOTC),
		d.Set("region", GetRegion(d, config)),
		d.Set("eip", eip),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting cce cluster fields: %s", err)
	}

	cert, err := clusters.GetCert(cceClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error retrieving opentelekomcloud CCE cluster cert: %s", err)
	}

	// Set Certificate Clusters
	var clusterList []map[string]interface{}
	for _, clusterObj := range cert.Clusters {
		clusterCert := map[string]interface{}{
			"name":                       clusterObj.Name,
			"server":                     clusterObj.Cluster.Server,
			"certificate_authority_data": clusterObj.Cluster.CertAuthorityData,
		}
		clusterList = append(clusterList, clusterCert)
	}
	if err := d.Set("certificate_clusters", clusterList); err != nil {
		return err
	}

	// Set Certificate Users
	var userList []map[string]interface{}
	for _, userObj := range cert.Users {
		userCert := map[string]interface{}{
			"name":                    userObj.Name,
			"client_certificate_data": userObj.User.ClientCertData,
			"client_key_data":         userObj.User.ClientKeyData,
		}
		userList = append(userList, userCert)
	}
	if err := d.Set("certificate_users", userList); err != nil {
		return err
	}

	return nil
}

func resourceCCEClusterV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud CCE Client: %s", err)
	}

	var updateOpts clusters.UpdateOpts

	if d.HasChange("description") {
		updateOpts.Spec.Description = d.Get("description").(string)
		_, err = clusters.Update(cceClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("error updating opentelekomcloud CCE: %s", err)
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
				return fmt.Errorf("error retrieving the eip: %s", err)
			}
			if fipId == "" {
				return fmt.Errorf("the specified EIP %s does not exist", newEipStr)
			}
		}
		if oldEipStr != "" {
			updateIpOpts := clusters.UpdateIpOpts{
				Action: "unbind",
			}
			err = clusters.UpdateMasterIp(cceClient, d.Id(), updateIpOpts).ExtractErr()
			if err != nil {
				return fmt.Errorf("error unbinding EIP to opentelekomcloud CCE: %s", err)
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
				return fmt.Errorf("error binding EIP to opentelekomcloud CCE: %s", err)
			}
		}
	}

	return resourceCCEClusterV3Read(d, meta)
}

func resourceCCEClusterV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud CCE Client: %s", err)
	}
	err = clusters.Delete(cceClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("error deleting opentelekomcloud CCE Cluster: %s", err)
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
		return fmt.Errorf("error deleting opentelekomcloud CCE cluster: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCCEClusterActive(cceClient *golangsdk.ServiceClient, clusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := clusters.Get(cceClient, clusterId).Extract()
		if err != nil {
			return nil, "", fmt.Errorf("error waiting for CCE cluster to become active: %s", err)
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
			return nil, "", fmt.Errorf("error waiting CCE cluster to become deleted: %s", err)
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
		return "", fmt.Errorf("error creating opentelekomcloud networking Client: %s", err)
	}
	listOpts := floatingips.ListOpts{
		FloatingIP: floatingIP,
	}
	allPages, err := floatingips.List(networkClient, listOpts).AllPages()
	if err != nil {
		return "", fmt.Errorf("error listing floating IPs: %s", err)
	}

	allFips, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return "", fmt.Errorf("error extracting floating IPs: %s", err)
	}

	if len(allFips) == 0 {
		return "", nil
	}

	return allFips[0].ID, nil
}
