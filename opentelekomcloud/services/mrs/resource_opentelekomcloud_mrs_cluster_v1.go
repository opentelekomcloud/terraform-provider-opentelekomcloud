package mrs

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/mrs/v1/cluster"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceMRSClusterV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterV1Create,
		ReadContext:   resourceClusterV1Read,
		UpdateContext: resourceClusterV1Update,
		DeleteContext: resourceClusterV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"billing_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"master_node_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"master_node_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"core_node_num": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"core_node_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"available_zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Required: true,
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
			"cluster_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MRS 1.6.3", "MRS 1.7.2", "MRS 1.9.2", "MRS 2.1.0", "MRS 3.0.2",
				}, false),
			},
			"cluster_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SATA", "SAS", "SSD",
				}, false),
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"master_data_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SATA", "SAS", "SSD",
				}, false),
			},
			"master_data_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"master_data_volume_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"core_data_volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SATA", "SAS", "SSD",
				}, false),
			},
			"core_data_volume_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"core_data_volume_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"node_public_cert_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"safe_mode": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"cluster_admin_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log_collection": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"component_list": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"component_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"component_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"add_jobs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"job_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"jar_path": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"arguments": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"input": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"output": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"job_log": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"shutdown_cluster": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"file_action": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"submit_job_once_cluster_run": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"hql": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"hive_script_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"bootstrap_scripts": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"parameters": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"nodes": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"active_master": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"before_component_start": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"fail_action": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hadoop_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_node_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_first": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slave_security_groups_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_groups_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_alternate_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_node_spec_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_node_spec_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"master_node_product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core_node_product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vnc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fee": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getAllClusterComponents(d *schema.ResourceData) []cluster.ComponentList {
	var componentOpts []cluster.ComponentList

	components := d.Get("component_list").(*schema.Set)
	for _, v := range components.List() {
		component := v.(map[string]interface{})

		componentOpts = append(componentOpts, cluster.ComponentList{
			ComponentName: component["component_name"].(string),
		})
	}

	return componentOpts
}

func getAllClusterJobs(d *schema.ResourceData) []cluster.AddJobs {
	var jobOpts []cluster.AddJobs

	jobs := d.Get("add_jobs").([]interface{})
	for _, v := range jobs {
		job := v.(map[string]interface{})

		shutDown := job["shutdown_cluster"].(bool)
		submitJob := job["submit_job_once_cluster_run"].(bool)
		jobOpts = append(jobOpts, cluster.AddJobs{
			JobType:                 job["job_type"].(int),
			JobName:                 job["job_name"].(string),
			JarPath:                 job["jar_path"].(string),
			Arguments:               job["arguments"].(string),
			Input:                   job["input"].(string),
			Output:                  job["output"].(string),
			JobLog:                  job["job_log"].(string),
			ShutdownCluster:         &shutDown,
			FileAction:              job["file_action"].(string),
			SubmitJobOnceClusterRun: &submitJob,
			Hql:                     job["hql"].(string),
			HiveScriptPath:          job["hive_script_path"].(string),
		})
	}

	return jobOpts
}

func getAllClusterScripts(d *schema.ResourceData) []cluster.BootstrapScript {
	var scriptOpts []cluster.BootstrapScript

	scripts := d.Get("bootstrap_scripts").([]interface{})
	for _, v := range scripts {
		script := v.(map[string]interface{})

		var nodes []string
		nodesRaw := script["nodes"].([]interface{})
		if len(nodesRaw) > 0 {
			for _, n := range nodesRaw {
				nodes = append(nodes, n.(string))
			}
		}

		activeMaster := script["active_master"].(bool)
		beforeComponent := script["before_component_start"].(bool)
		scriptOpts = append(scriptOpts, cluster.BootstrapScript{
			Name:                 script["name"].(string),
			Uri:                  script["uri"].(string),
			Parameters:           script["parameters"].(string),
			Nodes:                nodes,
			ActiveMaster:         &activeMaster,
			BeforeComponentStart: &beforeComponent,
			FailAction:           script["fail_action"].(string),
		})
	}

	return scriptOpts
}

func resourceClusterV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}
	nwClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV1 client: %s", err)
	}

	// Get vpc name
	vpc, err := vpcs.Get(nwClient, d.Get("vpc_id").(string)).Extract()
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud VPC: %s", err)
	}
	// Get subnet name
	subnet, err := subnets.Get(nwClient, d.Get("subnet_id").(string)).Extract()
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud Subnet: %s", err)
	}

	createOpts := &cluster.CreateOpts{
		BillingType:           d.Get("billing_type").(int),
		DataCenter:            config.GetRegion(d),
		MasterNodeNum:         d.Get("master_node_num").(int),
		MasterNodeSize:        d.Get("master_node_size").(string),
		CoreNodeNum:           d.Get("core_node_num").(int),
		CoreNodeSize:          d.Get("core_node_size").(string),
		AvailableZoneId:       d.Get("available_zone_id").(string),
		ClusterName:           d.Get("cluster_name").(string),
		Vpc:                   vpc.Name,
		VpcId:                 d.Get("vpc_id").(string),
		SubnetId:              d.Get("subnet_id").(string),
		SubnetName:            subnet.Name,
		ClusterVersion:        d.Get("cluster_version").(string),
		ClusterType:           pointerto.Int(d.Get("cluster_type").(int)),
		VolumeType:            d.Get("volume_type").(string),
		VolumeSize:            d.Get("volume_size").(int),
		MasterDataVolumeType:  d.Get("master_data_volume_type").(string),
		MasterDataVolumeSize:  d.Get("master_data_volume_size").(int),
		MasterDataVolumeCount: d.Get("master_data_volume_count").(int),
		CoreDataVolumeType:    d.Get("core_data_volume_type").(string),
		CoreDataVolumeSize:    d.Get("core_data_volume_size").(int),
		CoreDataVolumeCount:   d.Get("core_data_volume_count").(int),
		NodePublicCertName:    d.Get("node_public_cert_name").(string),
		SafeMode:              d.Get("safe_mode").(int),
		ClusterAdminSecret:    d.Get("cluster_admin_secret").(string),
		LogCollection:         pointerto.Int(d.Get("log_collection").(int)),
		ComponentList:         getAllClusterComponents(d),
		AddJobs:               getAllClusterJobs(d),
		BootstrapScripts:      getAllClusterScripts(d),
		LoginMode:             pointerto.Int(1),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	mrsCluster, err := cluster.Create(client, *createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Cluster: %s", err)
	}
	d.SetId(mrsCluster.ClusterId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"starting"},
		Target:     []string{"running"},
		Refresh:    clusterStateRefreshFunc(client, mrsCluster.ClusterId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for cluster (%s) to become ready: %s ", mrsCluster.ClusterId, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := common.ExpandResourceTags(tagRaw)
		if err := tags.Create(client, "clusters", mrsCluster.ClusterId, tagList).ExtractErr(); err != nil {
			return fmterr.Errorf("error setting tags of MRS cluster %s: %s", mrsCluster.ClusterId, err)
		}
	}

	return resourceClusterV1Read(ctx, d, meta)
}

func resourceClusterV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	mrsCluster, err := cluster.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Cluster")
	}
	log.Printf("[DEBUG] Retrieved Cluster %s: %#v", d.Id(), mrsCluster)

	masterNodeNum, err := strconv.Atoi(mrsCluster.MasterNodeNum)
	if err != nil {
		return fmterr.Errorf("error converting MasterNodeNum: %s", err)
	}
	coreNodeNum, err := strconv.Atoi(mrsCluster.CoreNodeNum)
	if err != nil {
		return fmterr.Errorf("error converting CoreNodeNum: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("order_id", mrsCluster.OrderId),
		d.Set("cluster_id", mrsCluster.ClusterId),
		d.Set("available_zone_name", mrsCluster.AzName),
		d.Set("available_zone_id", mrsCluster.AzId),
		d.Set("cluster_version", mrsCluster.ClusterVersion),
		d.Set("master_node_num", masterNodeNum),
		d.Set("core_node_num", coreNodeNum),
		d.Set("cluster_name", mrsCluster.ClusterName),
		d.Set("core_node_size", mrsCluster.CoreNodeSize),
		d.Set("master_data_volume_type", mrsCluster.MasterDataVolumeType),
		d.Set("master_data_volume_size", mrsCluster.MasterDataVolumeSize),
		d.Set("master_data_volume_count", mrsCluster.MasterDataVolumeCount),
		d.Set("core_data_volume_type", mrsCluster.CoreDataVolumeType),
		d.Set("core_data_volume_size", mrsCluster.CoreDataVolumeSize),
		d.Set("core_data_volume_count", mrsCluster.CoreDataVolumeCount),
		d.Set("node_public_cert_name", mrsCluster.NodePublicCertName),
		d.Set("safe_mode", mrsCluster.SafeMode),
		d.Set("master_node_size", mrsCluster.MasterNodeSize),
		d.Set("instance_id", mrsCluster.InstanceId),
		d.Set("hadoop_version", mrsCluster.HadoopVersion),
		d.Set("master_node_ip", mrsCluster.MasterNodeIp),
		d.Set("external_ip", mrsCluster.ExternalIp),
		d.Set("private_ip_first", mrsCluster.PrivateIpFirst),
		d.Set("internal_ip", mrsCluster.InternalIp),
		d.Set("slave_security_groups_id", mrsCluster.SlaveSecurityGroupsId),
		d.Set("security_groups_id", mrsCluster.SecurityGroupsId),
		d.Set("external_alternate_ip", mrsCluster.ExternalAlternateIp),
		d.Set("master_node_spec_id", mrsCluster.MasterNodeSpecId),
		d.Set("core_node_spec_id", mrsCluster.CoreNodeSpecId),
		d.Set("master_node_product_id", mrsCluster.MasterNodeProductId),
		d.Set("core_node_product_id", mrsCluster.CoreNodeProductId),
		d.Set("vnc", mrsCluster.Vnc),
		d.Set("fee", mrsCluster.Fee),
		d.Set("deployment_id", mrsCluster.DeploymentId),
		d.Set("cluster_state", mrsCluster.ClusterState),
		d.Set("error_info", mrsCluster.ErrorInfo),
		d.Set("remark", mrsCluster.Remark),
		d.Set("tenant_id", mrsCluster.TenantId),
	)

	updateAt, err := strconv.ParseInt(mrsCluster.UpdateAt, 10, 64)
	if err != nil {
		return fmterr.Errorf("error converting UpdateAt: %s", err)
	}

	createAt, err := strconv.ParseInt(mrsCluster.CreateAt, 10, 64)
	if err != nil {
		return fmterr.Errorf("error converting CreateAt: %s", err)
	}

	chargingStartTime, err := strconv.ParseInt(mrsCluster.ChargingStartTime, 10, 64)
	if err != nil {
		return fmterr.Errorf("error converting ChargingStartTime: %s", err)
	}

	mErr = multierror.Append(mErr,
		d.Set("update_at", time.Unix(updateAt, 0).String()),
		d.Set("create_at", time.Unix(createAt, 0).String()),
		d.Set("charging_start_time", time.Unix(chargingStartTime, 0).String()),
	)

	components := make([]map[string]interface{}, len(mrsCluster.ComponentList))
	for i, attachment := range mrsCluster.ComponentList {
		components[i] = make(map[string]interface{})
		components[i]["component_id"] = attachment.ComponentId
		components[i]["component_name"] = attachment.ComponentName
		components[i]["component_version"] = attachment.ComponentVersion
		components[i]["component_desc"] = attachment.ComponentDesc
	}

	scripts := make([]map[string]interface{}, len(mrsCluster.BootstrapScripts))
	for i, script := range mrsCluster.BootstrapScripts {
		scripts[i] = make(map[string]interface{})
		scripts[i]["name"] = script.Name
		scripts[i]["uri"] = script.Uri
		scripts[i]["parameters"] = script.Parameters
		scripts[i]["nodes"] = script.Nodes
		scripts[i]["active_master"] = script.ActiveMaster
		scripts[i]["before_component_start"] = script.BeforeComponentStart
		scripts[i]["fail_action"] = script.FailAction
	}

	mErr = multierror.Append(mErr,
		d.Set("component_list", components),
		d.Set("bootstrap_scripts", scripts),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// save tags
	resourceTags, err := tags.Get(client, "clusters", d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error fetching OpenTelekomCloud MRS Cluster tags: %s", err)
	}
	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmterr.Errorf("error saving tags for OpenTelekomCloud MRS Cluster: %s", err)
	}

	return nil
}

func resourceClusterV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	// update tags
	if d.HasChange("tags") {
		if err := common.UpdateResourceTags(client, d, "clusters", d.Id()); err != nil {
			return fmterr.Errorf("error updating tags of MRS cluster %s: %s", d.Id(), err)
		}
	}

	return resourceClusterV1Read(ctx, d, meta)
}

func resourceClusterV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.MrsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(ErrCreationClient, err)
	}

	if err := cluster.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for Cluster (%s) to be terminated", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"running", "terminating"},
		Target:     []string{"terminated"},
		Refresh:    clusterStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for Cluster (%s) to be terminated: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func clusterStateRefreshFunc(client *golangsdk.ServiceClient, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := cluster.Get(client, clusterID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return n, "DELETED", nil
			}
			return nil, "", err
		}

		return n, n.ClusterState, nil
	}
}
