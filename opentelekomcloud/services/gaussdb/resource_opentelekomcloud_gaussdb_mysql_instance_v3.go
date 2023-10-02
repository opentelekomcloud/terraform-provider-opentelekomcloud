package gaussdb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	v3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/gaussdb/v3"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gaussdb/v3/backup"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gaussdb/v3/instance"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceGaussDBInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBInstanceV3Create,
		UpdateContext: resourceGaussDBInstanceV3Update,
		ReadContext:   resourceGaussDBInstanceV3Read,
		DeleteContext: resourceGaussDBInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
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
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"configuration_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dedicated_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dedicated_resource_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"read_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "UTC+08:00",
			},
			"availability_zone_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "single",
				ValidateFunc: validation.StringInSlice([]string{
					"single", "multi",
				}, true),
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"gaussdb-mysql",
							}, true),
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_write_ip": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"private_read_ips": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"az_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_ref": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_connections": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcpus": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ram": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"need_restart": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBDataStore(d *schema.ResourceData) instance.Datastore {
	var db instance.Datastore

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
	} else {
		db.Type = "gaussdb-mysql"
		db.Version = "8.0"
	}
	return db
}

func GaussDBInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instance.GetInstance(client, instanceID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		if v.Id == "" {
			return v, "DELETED", nil
		}
		return v, v.Status, nil
	}
}

func resourceGaussDBInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GaussDBV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GaussDB client: %s ", err)
	}

	createOpts := instance.CreateInstanceOpts{
		Name:                d.Get("name").(string),
		FlavorRef:           d.Get("flavor").(string),
		Region:              config.GetRegion(d),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		DedicatedResourceId: d.Get("dedicated_resource_id").(string),
		TimeZone:            d.Get("time_zone").(string),
		SlaveCount:          pointerto.Int(d.Get("read_replicas").(int)),
		Mode:                "Cluster",
		Datastore:           resourceGaussDBDataStore(d),
	}

	azMode := d.Get("availability_zone_mode").(string)
	createOpts.AvailabilityZoneMode = azMode
	if azMode == "multi" {
		v, exist := d.GetOk("master_availability_zone")
		if !exist {
			return fmterr.Errorf("missing master_availability_zone in a multi availability zone mode")
		}
		createOpts.MasterAvailabilityZone = v.(string)
	}

	if _, ok := d.GetOk("backup_strategy"); ok {
		var backupOpts instance.BackupStrategy
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keepDays := rawMap["keep_days"].(string)
		backupOpts.KeepDays = keepDays
		backupOpts.StartTime = rawMap["start_time"].(string)
		createOpts.BackupStrategy = &backupOpts
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	createOpts.Password = d.Get("password").(string)

	inst, err := instance.CreateInstance(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating GaussDB instance : %s", err)
	}

	d.SetId(inst.Instance.Id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      GaussDBInstanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        180 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", d.Id(), err)
	}

	return resourceGaussDBInstanceV3Read(ctx, d, meta)
}

func resourceGaussDBInstanceV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GaussDBV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GaussDB client: %s ", err)
	}

	inst, err := instance.GetInstance(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "GaussDB instance")
	}
	if inst.Id == "" {
		return fmterr.Errorf("error retrieving OpenTelekomCloud GaussDB instance: %s ", err)
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", d.Id(), inst)

	// set data store
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": inst.Datastore.Version,
		"engine":  inst.Datastore.Type,
	}
	dbList = append(dbList, db)

	port, err := strconv.Atoi(inst.Port)
	if err != nil {
		common.CheckDeletedDiag(d, err, "incorrect port format")
	}

	allNodes := *inst.Nodes

	mErr := multierror.Append(
		d.Set("region", inst.Region),
		d.Set("project_id", inst.ProjectId),
		d.Set("name", inst.Name),
		d.Set("status", inst.Status),
		d.Set("mode", inst.Type),
		d.Set("vpc_id", inst.VpcId),
		d.Set("flavor", allNodes[0].FlavorRef),
		d.Set("subnet_id", inst.SubnetId),
		d.Set("security_group_id", inst.SecurityGroupId),
		d.Set("configuration_id", inst.ConfigurationId),
		d.Set("dedicated_resource_id", inst.DedicatedResourceId),
		d.Set("db_user_name", inst.DbUserName),
		d.Set("time_zone", inst.TimeZone),
		d.Set("availability_zone_mode", inst.AzMode),
		d.Set("master_availability_zone", inst.MasterAzCode),
		d.Set("port", port),
		d.Set("datastore", dbList),
		d.Set("alias", inst.Alias),
		d.Set("charging_mode", inst.ChargeInfo.ChargeMode),
		d.Set("node_count", inst.NodeCount),
		d.Set("created", inst.Created),
		d.Set("updated", inst.Updated),
		d.Set("public_ip", inst.PublicIps),
		d.Set("private_write_ip", []string{}),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	if len(inst.PrivateIps) > 0 {
		var ipList []string
		copy(ipList, inst.PrivateIps)
		mErr = multierror.Append(d.Set("private_write_ip", ipList))
	}
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// set nodes
	slaveCount := 0
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range *inst.Nodes {
		node := map[string]interface{}{
			"id":              raw.Id,
			"name":            raw.Name,
			"status":          raw.Status,
			"type":            raw.Type,
			"port":            raw.Port,
			"az_code":         raw.AzCode,
			"region_code":     raw.RegionCode,
			"created":         raw.Created,
			"updated":         raw.Updated,
			"flavor_ref":      raw.FlavorRef,
			"max_connections": raw.MaxConnections,
			"vcpus":           raw.Vcpus,
			"ram":             raw.Ram,
			"need_restart":    raw.NeedRestart,
			"priority":        raw.Priority,
		}
		if len(raw.PrivateReadIps) > 0 {
			node["private_read_ips"] = raw.PrivateReadIps
		}
		nodesList = append(nodesList, node)
		if raw.Type == "slave" && (raw.Status == "ACTIVE" || raw.Status == "BACKING UP") {
			slaveCount += 1
		}
	}

	mErr = multierror.Append(
		d.Set("nodes", nodesList),
		d.Set("read_replicas", slaveCount),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// set backup_strategy
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": inst.BackupStrategy.StartTime,
	}
	if days, err := strconv.Atoi(inst.BackupStrategy.KeepDays); err == nil {
		backupStrategy["keep_days"] = days
	}
	backupStrategyList[0] = backupStrategy
	mErr = multierror.Append(d.Set("backup_strategy", backupStrategyList))
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGaussDBInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GaussDBV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GaussDB client: %s ", err)
	}

	if d.HasChange("name") {
		newName := d.Get("name").(string)
		updateNameOpts := instance.UpdateNameOpts{
			Name:       newName,
			InstanceId: d.Id(),
		}
		log.Printf("[DEBUG] Update Name Options: %+v", updateNameOpts)

		n, err := instance.UpdateName(client, updateNameOpts)
		if err != nil {
			return fmterr.Errorf("error updating name for instance %s: %s ", d.Id(), err)
		}

		if _, err = waitForGaussJob(client, *n, int(d.Timeout(schema.TimeoutUpdate)/time.Second)); err != nil {
			return diag.FromErr(err)
		}
		log.Printf("[DEBUG] Updated Name to %s for instance %s", newName, d.Id())
	}

	if d.HasChange("password") {
		newPass := d.Get("password").(string)
		updatePassOpts := instance.ResetPwdOpts{
			Password:   newPass,
			InstanceId: d.Id(),
		}

		err = instance.ResetPassword(client, updatePassOpts)
		if err != nil {
			return fmterr.Errorf("error updating password for instance %s: %s ", d.Id(), err)
		}
		log.Printf("[DEBUG] Updated Password for instance %s", d.Id())
	}

	if d.HasChange("flavor") {
		newFlavor := d.Get("flavor").(string)
		resizeOpts := instance.UpdateSpecOpts{
			InstanceId: d.Id(),
			ResizeFlavor: instance.ResizeFlavor{
				SpecCode: newFlavor,
			},
		}
		log.Printf("[DEBUG] Update Flavor Options: %+v", resizeOpts)

		n, err := instance.UpdateInstance(client, resizeOpts)
		if err != nil {
			return fmterr.Errorf("error updating flavor for instance %s: %s ", d.Id(), err)
		}

		// wait for job success
		if n.JobId != "" {
			if _, err := waitForGaussJob(client, n.JobId, int(d.Timeout(schema.TimeoutUpdate)/time.Second)); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("read_replicas") {
		oldNum, newNum := d.GetChange("read_replicas")
		if newNum.(int) > oldNum.(int) {
			expand_size := newNum.(int) - oldNum.(int)
			var priorities []int
			for i := 0; i < expand_size; i++ {
				priorities = append(priorities, 1)
			}
			createReplicaOpts := instance.CreateNodeOpts{
				Priorities: priorities,
				InstanceId: d.Id(),
			}
			log.Printf("[DEBUG] Create Replica Options: %+v", createReplicaOpts)

			n, err := instance.CreateReplica(client, createReplicaOpts)
			if err != nil {
				return fmterr.Errorf("error creating read replicas for instance %s: %s ", d.Id(), err)
			}

			// wait for job success
			if n.JobId != "" {
				jobList := strings.Split(n.JobId, ",")
				log.Printf("[DEBUG] Create Replica Jobs: %#v", jobList)
				for i := 0; i < len(jobList); i++ {
					jobId := jobList[i]
					log.Printf("[DEBUG] Waiting for job: %s", jobId)
					if _, err = waitForGaussJob(client, jobId, int(d.Timeout(schema.TimeoutUpdate)/time.Second)); err != nil {
						return diag.FromErr(err)
					}
				}
			}
		}
		if newNum.(int) < oldNum.(int) {
			shrinkSize := oldNum.(int) - newNum.(int)

			slaveNodes := []string{}
			nodes := d.Get("nodes").([]interface{})
			for _, nodeRaw := range nodes {
				node := nodeRaw.(map[string]interface{})
				if node["type"].(string) == "slave" && node["status"] == "ACTIVE" {
					slaveNodes = append(slaveNodes, node["id"].(string))
				}
			}
			log.Printf("[DEBUG] Slave Nodes: %+v", slaveNodes)
			if len(slaveNodes) <= shrinkSize {
				return fmterr.Errorf("error deleting read replicas for instance %s: Shrink Size is bigger than active slave nodes", d.Id())
			}
			for i := 0; i < shrinkSize; i++ {
				n, err := instance.DeleteReplica(client, d.Id(), slaveNodes[i])
				if err != nil {
					return fmterr.Errorf("error creating read replica %s for instance %s: %s ", slaveNodes[i], d.Id(), err)
				}

				if _, err := waitForGaussJob(client, *n, int(d.Timeout(schema.TimeoutUpdate)/time.Second)); err != nil {
					return diag.FromErr(err)
				}
				log.Printf("[DEBUG] Deleted Read Replica: %s", slaveNodes[i])
			}
		}
	}

	if d.HasChange("backup_strategy") {
		updateOpts := backup.UpdatePolicyOpts{
			InstanceId:   d.Id(),
			BackupPolicy: backup.UpdateBackupPolicy{},
		}
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keepDays := rawMap["keep_days"].(int)
		updateOpts.BackupPolicy.KeepDays = keepDays
		updateOpts.BackupPolicy.StartTime = rawMap["start_time"].(string)
		// Fixed to "1,2,3,4,5,6,7"
		updateOpts.BackupPolicy.Period = "1,2,3,4,5,6,7"
		log.Printf("[DEBUG] Update backup_strategy: %#v", updateOpts)

		_, err = backup.UpdatePolicy(client, updateOpts)
		if err != nil {
			return fmterr.Errorf("error updating backup_strategy: %s", err)
		}
	}

	return resourceGaussDBInstanceV3Read(ctx, d, meta)
}

func resourceGaussDBInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GaussDBV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud GaussDB client: %s ", err)
	}

	_, err = instance.DeleteInstance(client, d.Id())
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "GaussDB instance"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:     []string{"DELETED"},
		Refresh:    GaussDBInstanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"error waiting for instance (%s) to be deleted: %s ",
			d.Id(), err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", d.Id())
	return nil
}

func waitForGaussJob(client *golangsdk.ServiceClient, jobId string, timeout int) (bool, error) {
	err := golangsdk.WaitFor(timeout, func() (bool, error) {
		cur, err := v3.ShowJobInfo(client, jobId)
		if err != nil {
			return false, err
		}

		if cur.Status == "Completed" {
			return true, nil
		}

		if cur.Status == "Failed" {
			return false, fmt.Errorf("job %s failed: %s", jobId, cur.FailReason)
		}

		return false, nil
	})

	return false, err
}
