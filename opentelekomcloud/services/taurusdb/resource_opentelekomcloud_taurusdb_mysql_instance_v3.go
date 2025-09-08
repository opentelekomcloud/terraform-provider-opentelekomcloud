package taurusdb

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
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/backup"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/instance"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/taurus/v3/job"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceTaurusDbV3Instance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDbV3InstanceCreate,
		UpdateContext: resourceTaurusDbV3InstanceUpdate,
		ReadContext:   resourceTaurusDbV3InstanceRead,
		DeleteContext: resourceTaurusDbV3InstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
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
				ForceNew: true,
				Optional: true,
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"table_name_case_sensitivity": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"read_replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
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
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"private_write_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"seconds_level_monitoring_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"seconds_level_monitoring_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"seconds_level_monitoring_enabled"},
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
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_dns_name": {
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
						"private_read_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTaurusDbV3DataStore(d *schema.ResourceData) instance.DataStoreOpt {
	var db instance.DataStoreOpt

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

func TaurusDbInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instance.Get(client, instanceID)
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

func resourceTaurusDbV3InstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	createOpts := instance.CreateOpts{
		Name:                d.Get("name").(string),
		Flavor:              d.Get("flavor").(string),
		Region:              config.GetRegion(d),
		VpcId:               d.Get("vpc_id").(string),
		SubnetId:            d.Get("subnet_id").(string),
		SecurityGroupId:     d.Get("security_group_id").(string),
		ConfigurationId:     d.Get("configuration_id").(string),
		DedicatedResourceId: d.Get("dedicated_resource_id").(string),
		TimeZone:            d.Get("time_zone").(string),
		SlaveCount:          d.Get("read_replicas").(int),
		EnterpriseProjectId: pointerto.String(d.Get("enterprise_project_id").(string)),
		Mode:                "Cluster",
		DataStore:           resourceTaurusDbV3DataStore(d),
	}

	if d.Get("table_name_case_sensitivity").(bool) {
		lowerCaseTableNames := 0
		createOpts.LowerCaseTableNames = &lowerCaseTableNames
	}

	azMode := d.Get("availability_zone_mode").(string)
	createOpts.AZMode = azMode
	if azMode == "multi" {
		v, exist := d.GetOk("master_availability_zone")
		if !exist {
			return diag.Errorf("missing master_availability_zone in a multi availability zone mode")
		}
		createOpts.MasterAZ = v.(string)
	}

	if common.HasFilledOpt(d, "volume_size") {
		volume := &instance.VolumeOpt{
			Size: d.Get("volume_size").(int),
		}
		createOpts.Volume = volume
	}

	log.Printf("[DEBUG] create options: %#v", createOpts)
	createOpts.Password = d.Get("password").(string)

	instance, err := instance.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating TaurusDB instance : %s", err)
	}

	id := instance.Instance.Id
	d.SetId(id)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      TaurusDbInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        180 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to become ready: %s",
			id, err)
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(120 * time.Second)

	stateConf = &resource.StateChangeConf{
		Pending:      []string{"BUILD", "BACKING UP"},
		Target:       []string{"ACTIVE"},
		Refresh:      TaurusDbInstanceStateRefreshFunc(client, id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to become ready: %s", id, err)
	}

	if _, ok := d.GetOk("backup_strategy"); ok {
		if err = updateInstanceBackupStrategy(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("port"); ok {
		if err = updatePort(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("seconds_level_monitoring_enabled"); ok {
		if err = updatesSecondsLevelMonitoring(ctx, client, d, schema.TimeoutCreate); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaurusDbV3InstanceRead(ctx, d, meta)
}

func resourceTaurusDbV3InstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.TaurusDBV3Client(region)
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s", err)
	}
	var mErr *multierror.Error

	instanceID := d.Id()
	instance, err := instance.Get(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving TaurusDb MySQL instance")
	}
	if instance.Id == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving TaurusDb MySQL instance")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", instance.Name),
		d.Set("status", instance.Status),
		d.Set("mode", instance.Type),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("configuration_id", instance.ConfigurationId),
		d.Set("dedicated_resource_id", instance.DedicatedResourceId),
		d.Set("db_user_name", instance.DbUserName),
		d.Set("time_zone", instance.TimeZone),
		d.Set("availability_zone_mode", instance.AZMode),
		d.Set("master_availability_zone", instance.MasterAZ),
		d.Set("created_at", instance.Created),
		d.Set("updated_at", instance.Updated),
		d.Set("enterprise_project_id", instance.EnterpriseProjectId),
	)

	if dbPort, err := strconv.Atoi(instance.Port); err == nil {
		mErr = multierror.Append(mErr, d.Set("port", dbPort))
	}
	if len(instance.PrivateIps) > 0 {
		mErr = multierror.Append(mErr, d.Set("private_write_ip", instance.PrivateIps[0]))
	}

	mErr = multierror.Append(mErr, setDatastore(d, instance.DataStore))
	mErr = multierror.Append(mErr, setNodes(d, instance.Nodes)...)
	mErr = multierror.Append(mErr, setBackupStrategy(d, instance.BackupStrategy))
	mErr = multierror.Append(mErr, setSecondsLevelMonitoring(d, client, instanceID)...)

	return diag.FromErr(mErr.ErrorOrNil())
}

func setNodes(d *schema.ResourceData, nodes []instance.Nodes) []error {
	flavor := ""
	slaveCount := 0
	volumeSize := 0
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, raw := range nodes {
		node := map[string]interface{}{
			"id":                raw.Id,
			"name":              raw.Name,
			"status":            raw.Status,
			"type":              raw.Type,
			"availability_zone": raw.AvailabilityZone,
		}
		if len(raw.PrivateIps) > 0 {
			node["private_read_ip"] = raw.PrivateIps[0]
		}
		if raw.Volume.Size > 0 {
			volumeSize = int(raw.Volume.Size)
		}
		nodesList = append(nodesList, node)
		if raw.Type == "slave" && (raw.Status == "ACTIVE" || raw.Status == "BACKING UP") {
			slaveCount++
		}
		if flavor == "" {
			flavor = raw.Flavor
		}
	}
	var errs []error
	errs = append(errs, d.Set("nodes", nodesList))
	errs = append(errs, d.Set("read_replicas", slaveCount))
	errs = append(errs, d.Set("volume_size", volumeSize))
	if flavor != "" {
		log.Printf("[DEBUG] node flavor: %s", flavor)
		errs = append(errs, d.Set("flavor", flavor))
	}
	return errs
}

func setDatastore(d *schema.ResourceData, datastore instance.DataStore) error {
	dbList := make([]map[string]interface{}, 1)
	db := map[string]interface{}{
		"version": normalizeVersion(datastore.Version),
	}

	engine := datastore.Type
	if engine == "GaussDB(for MySQL)" {
		engine = "gaussdb-mysql"
	}
	db["engine"] = engine
	dbList[0] = db
	return d.Set("datastore", dbList)
}

func normalizeVersion(v string) string {
	parts := strings.Split(v, ".")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s.%s", parts[0], parts[1])
	}
	return v
}

func setBackupStrategy(d *schema.ResourceData, strategy instance.BackupStrategy) error {
	backupStrategyList := make([]map[string]interface{}, 1)
	backupStrategy := map[string]interface{}{
		"start_time": strategy.StartTime,
	}
	if days, err := strconv.Atoi(strategy.KeepDays); err == nil {
		backupStrategy["keep_days"] = days
	}
	backupStrategyList[0] = backupStrategy
	return d.Set("backup_strategy", backupStrategyList)
}

func setSecondsLevelMonitoring(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) []error {
	resp, err := instance.GetSecondLevelMonitoring(client, instanceId)
	if err != nil {
		log.Printf("[WARN] query instance %s seconds level monitoring failed: %s", instanceId, err)
		return nil
	}
	var errs []error
	errs = append(errs, d.Set("seconds_level_monitoring_enabled", resp.MonitorSwitch))
	errs = append(errs, d.Set("seconds_level_monitoring_period", resp.Period))
	return errs
}

func resourceTaurusDbV3InstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := config.TaurusDBV3Client(region)
	if err != nil {
		return diag.Errorf("error creating TaurusDb client: %s ", err)
	}

	if d.HasChange("name") {
		if err = updateInstanceName(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("password") {
		newPass := d.Get("password").(string)
		err = instance.UpdatePass(client, d.Id(), newPass)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("flavor") {
		if err = updateInstanceFlavor(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("read_replicas") {
		if err = updateInstanceReadReplica(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("backup_strategy") {
		if err = updateInstanceBackupStrategy(ctx, client, d, schema.TimeoutUpdate); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("port") {
		err = updatePort(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("seconds_level_monitoring_enabled", "seconds_level_monitoring_period") {
		err = updatesSecondsLevelMonitoring(ctx, client, d, schema.TimeoutUpdate)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaurusDbV3InstanceRead(ctx, d, meta)
}

func resourceTaurusDbV3InstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.TaurusDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating TaurusDbV3 client: %s ", err)
	}

	instanceId := d.Id()

	_, err = instance.Delete(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "TaurusDB instance")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "BACKING UP", "FAILED"},
		Target:     []string{"DELETED"},
		Refresh:    TaurusDbInstanceStateRefreshFunc(client, instanceId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to be deleted: %s ", instanceId, err)
	}
	log.Printf("[DEBUG] successfully deleted instance %s", instanceId)
	return nil
}

func updateInstanceName(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	newName := d.Get("name").(string)

	retryFunc := func() (interface{}, bool, error) {
		res, err := instance.UpdateName(client, d.Id(), newName)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     TaurusDbInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating name for instance %s: %s ", d.Id(), err)
	}

	job := r.(*string)
	return checkTaurusDbV3MySQLJobFinish(ctx, client, *job, d.Timeout(schema.TimeoutUpdate))
}

func updateInstanceFlavor(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	newFlavor := d.Get("flavor").(string)

	retryFunc := func() (interface{}, bool, error) {
		res, err := instance.Resize(client, d.Id(), newFlavor)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     TaurusDbInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating flavor for instance %s: %s ", d.Id(), err)
	}
	job := r.(*string)

	if *job != "" {
		if err = checkTaurusDbV3MySQLJobFinish(ctx, client, *job, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	return nil
}

func updateInstanceReadReplica(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldNum, newNum := d.GetChange("read_replicas")
	if newNum.(int) > oldNum.(int) {
		if err := createInstanceReadReplica(ctx, client, d, newNum.(int), oldNum.(int)); err != nil {
			return err
		}
	}
	if newNum.(int) < oldNum.(int) {
		if err := deleteInstanceReadReplica(ctx, client, d, newNum.(int), oldNum.(int)); err != nil {
			return err
		}
	}
	return nil
}

func createInstanceReadReplica(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	newNum, oldNum int) error {
	expandSize := newNum - oldNum
	priorities := make([]int, 0)
	for i := 0; i < expandSize; i++ {
		priorities = append(priorities, 1)
	}
	createReplicaOpts := instance.CreateReplicaOpts{
		Priorities: priorities,
	}
	retryFunc := func() (interface{}, bool, error) {
		res, err := instance.CreateReplica(client, d.Id(), createReplicaOpts)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     TaurusDbInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error creating read replicas for instance %s: %s ", d.Id(), err)
	}
	job := r.(*string)

	if *job != "" {
		jobList := strings.Split(*job, ",")
		log.Printf("[DEBUG] create replica jobs: %#v", jobList)
		for i := 0; i < len(jobList); i++ {
			jobId := jobList[i]
			log.Printf("[DEBUG] waiting for job: %s", jobId)
			if err = checkTaurusDbV3MySQLJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteInstanceReadReplica(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	newNum, oldNum int) error {
	shrinkSize := oldNum - newNum
	slaveNodes := make([]string, 0)
	nodes := d.Get("nodes").([]interface{})
	for _, nodeRaw := range nodes {
		node := nodeRaw.(map[string]interface{})
		if node["type"].(string) == "slave" && node["status"] == "ACTIVE" {
			slaveNodes = append(slaveNodes, node["id"].(string))
		}
	}
	log.Printf("[DEBUG] Slave Nodes: %+v", slaveNodes)
	if len(slaveNodes) <= shrinkSize {
		return fmt.Errorf("error deleting read replicas for instance %s: Shrink Size is bigger than active slave nodes", d.Id())
	}
	for i := 0; i < shrinkSize; i++ {
		retryFunc := func() (interface{}, bool, error) {
			res, err := instance.DeleteReplica(client, d.Id(), slaveNodes[i])
			retry, err := handleMultiOperationsError(err)
			return res, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     TaurusDbInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"ACTIVE"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return fmt.Errorf("error deleting read replicas for instance %s: %s ", d.Id(), err)
		}
		job := r.(*string)

		if *job != "" {
			if err = checkTaurusDbV3MySQLJobFinish(ctx, client, *job, d.Timeout(schema.TimeoutUpdate)); err != nil {
				return err
			}
		}
	}

	instance, err := instance.Get(client, d.Id())
	if err != nil {
		return err
	}
	slaveCount := 0
	for _, raw := range instance.Nodes {
		if raw.Type == "slave" && (raw.Status == "ACTIVE" || raw.Status == "BACKING UP") {
			slaveCount++
		}
	}
	if newNum != slaveCount {
		return fmt.Errorf("error updating read_replicas for instance %s: order failed", d.Id())
	}
	return nil
}

func updateInstanceBackupStrategy(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	var updateOpts backup.UpdatePolicyOpts
	backupRaw := d.Get("backup_strategy").([]interface{})
	rawMap := backupRaw[0].(map[string]interface{})
	keepDays := rawMap["keep_days"].(int)
	updateOpts.KeepDays = keepDays
	updateOpts.StartTime = rawMap["start_time"].(string)

	updateOpts.Period = "1,2,3,4,5,6,7"
	updateOpts.InstanceId = d.Id()
	log.Printf("[DEBUG] update backup_strategy: %#v", updateOpts)

	retryFunc := func() (interface{}, bool, error) {
		_, err := backup.UpdatePolicy(client, updateOpts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     taurusDbV3MysqlDatabaseStatusRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating backup_strategy: %s", err)
	}
	return nil
}

func updatePort(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout string) error {
	retryFunc := func() (interface{}, bool, error) {
		res, err := instance.UpdatePort(client, d.Id(), d.Get("port").(int))
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     taurusDbV3MysqlDatabaseStatusRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating port for instance %s: %s ", d.Id(), err)
	}

	job := r.(*string)
	return checkTaurusDbV3MySQLJobFinish(ctx, client, *job, d.Timeout(timeout))
}

func updatesSecondsLevelMonitoring(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout string) error {
	opts := instance.UpdateSecondLevelMonitoringOpts{
		InstanceId:    d.Id(),
		MonitorSwitch: d.Get("seconds_level_monitoring_enabled").(bool),
		Period:        d.Get("seconds_level_monitoring_period").(int),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := instance.UpdateSecondLevelMonitoring(client, opts)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     taurusDbV3MysqlDatabaseStatusRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(timeout),
		DelayTimeout: 30 * time.Second,
		PollInterval: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating seconds level monitoring for instance %s: %s ", d.Id(), err)
	}

	job := r.(*string)
	return checkTaurusDbV3MySQLJobFinish(ctx, client, *job, d.Timeout(timeout))
}

func checkTaurusDbV3MySQLJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending", "Running"},
		Target:       []string{"Completed"},
		Refresh:      taurusDbV3MysqlDatabaseStatusRefreshFunc(client, jobID),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for TaurusDbV3 MySQL instance job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func taurusDbV3MysqlDatabaseStatusRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobStatus, err := job.GetJobStatus(client, jobId)
		if err != nil {
			return nil, "Failed", err
		}

		return jobStatus, jobStatus.Job.Status, nil
	}
}
