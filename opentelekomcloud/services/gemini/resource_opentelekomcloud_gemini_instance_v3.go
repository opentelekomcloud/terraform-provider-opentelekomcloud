package gemini

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/backup"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/instance"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/gemini/v3/template"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

type defaultValues struct {
	Mode      string
	dbType    string
	dbVersion string
	logName   string
}

func ResourceGeminiDBInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBCassandraInstanceCreate,
		ReadContext:   resourceGeminiDBInstanceV3Read,
		UpdateContext: resourceGaussDBCassandraInstanceUpdate,
		DeleteContext: resourceGeminiDBInstanceV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3,
			},
			"volume_size": {
				Type:     schema.TypeInt,
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
			},
			"configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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
						"storage_engine": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
							Computed: true,
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
			"db_user_name": {
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
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"support_reduce": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"period": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"tags": common.TagsSchema(),
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGaussDBCassandraInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaults := defaultValues{
		Mode:      "Cluster",
		dbType:    "cassandra",
		dbVersion: "3.11",
		logName:   "cassandra",
	}
	return resourceGeminiDBInstanceV3Create(ctx, d, meta, defaults)
}

func resourceGeminiDBDataStore(d *schema.ResourceData, defaults defaultValues) instance.DataStoreOpt {
	var db instance.DataStoreOpt

	datastoreRaw := d.Get("datastore").([]interface{})
	if len(datastoreRaw) == 1 {
		datastore := datastoreRaw[0].(map[string]interface{})
		db.Type = datastore["engine"].(string)
		db.Version = datastore["version"].(string)
		db.StorageEngine = datastore["storage_engine"].(string)
	} else {
		db.Type = defaults.dbType
		db.Version = defaults.dbVersion
		db.StorageEngine = "rocksDB"
	}
	return db
}

func resourceGeminiDBBackupStrategy(d *schema.ResourceData) *instance.BackupStrategyOpt {
	if _, ok := d.GetOk("backup_strategy"); ok {
		opt := &instance.BackupStrategyOpt{
			StartTime: d.Get("backup_strategy.0.start_time").(string),
		}
		// The default value of keepdays is 7, but empty value of keepdays will be converted to 0.
		if v, ok := d.GetOk("backup_strategy.0.keep_days"); ok {
			opt.KeepDays = strconv.Itoa(v.(int))
		}
		return opt
	}
	return nil
}

func resourceGeminiDBFlavor(d *schema.ResourceData) []instance.FlavorOpt {
	var flavorList []instance.FlavorOpt
	flavor := instance.FlavorOpt{
		Num:      strconv.Itoa(d.Get("node_num").(int)),
		Size:     strconv.Itoa(d.Get("volume_size").(int)),
		Storage:  "ULTRAHIGH",
		SpecCode: d.Get("flavor").(string),
	}
	flavorList = append(flavorList, flavor)
	return flavorList
}

func GeminiDBInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		inst, err := GetInstanceByID(client, instanceID)

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return inst, "deleted", nil
			}
			return nil, "", err
		}

		return inst, inst.Status, nil
	}
}

func resourceGeminiDBInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}, defaults defaultValues) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	createOpts := instance.CreateOpts{
		Name:             d.Get("name").(string),
		Region:           config.GetRegion(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		VpcId:            d.Get("vpc_id").(string),
		SubnetId:         d.Get("subnet_id").(string),
		SecurityGroupId:  d.Get("security_group_id").(string),
		ConfigurationId:  d.Get("configuration_id").(string),
		Mode:             defaults.Mode,
		Flavor:           resourceGeminiDBFlavor(d),
		DataStore:        resourceGeminiDBDataStore(d, defaults),
		BackupStrategy:   resourceGeminiDBBackupStrategy(d),
	}
	if ssl := d.Get("ssl").(bool); ssl {
		createOpts.SslOption = pointerto.String("1")
	}

	log.Printf("[DEBUG] create options: %#v", createOpts)
	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	inst, err := instance.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating GeminiDB instance : %s", err)
	}

	d.SetId(inst.Id)
	// waiting for the instance to become ready
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"creating"},
		Target:       []string{"normal"},
		Refresh:      GeminiDBInstanceStateRefreshFunc(client, inst.Id),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        120 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to become ready: %s",
			inst.Id, err)
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := common.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(client, "instances", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of GeminiDB %s: %s", d.Id(), tagErr)
		}
	}

	// This is a workaround to avoid db connection issue
	time.Sleep(360 * time.Second) // lintignore:R018

	return resourceGeminiDBInstanceV3Read(ctx, d, meta)
}

func resourceGeminiDBInstanceV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	instanceID := d.Id()
	inst, err := GetInstanceByID(client, instanceID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "GeminiDB")
	}

	if inst.Id == "" {
		d.SetId("")
		log.Printf("[WARN] failed to fetch GeminiDB instance: deleted")
		return nil
	}

	log.Printf("[DEBUG] retrieved instance %s: %#v", instanceID, inst)

	mErr := multierror.Append(
		d.Set("name", inst.Name),
		d.Set("region", inst.Region),
		d.Set("status", inst.Status),
		d.Set("vpc_id", inst.VpcId),
		d.Set("subnet_id", inst.SubnetId),
		d.Set("security_group_id", inst.SecurityGroupId),
		d.Set("mode", inst.Mode),
		d.Set("db_user_name", inst.DbUserName),
		d.Set("tags", d.Get("tags")),
	)

	if dbPort, err := strconv.Atoi(inst.Port); err == nil {
		mErr = multierror.Append(mErr, d.Set("port", dbPort))
	}

	dbList := make([]map[string]interface{}, 0, 1)
	db := map[string]interface{}{
		"engine":         inst.DataStore.Type,
		"version":        inst.DataStore.Version,
		"storage_engine": inst.Engine,
	}
	dbList = append(dbList, db)
	mErr = multierror.Append(mErr, d.Set("datastore", dbList))

	specCode := ""
	wrongFlavor := "Inconsistent Flavor"
	ipsList := []string{}
	nodesList := make([]map[string]interface{}, 0, 1)
	for _, group := range inst.Groups {
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"id":             Node.Id,
				"name":           Node.Name,
				"status":         Node.Status,
				"private_ip":     Node.PrivateIp,
				"support_reduce": Node.SupportReduce,
			}
			if specCode == "" {
				specCode = Node.SpecCode
			} else if specCode != Node.SpecCode && specCode != wrongFlavor {
				specCode = wrongFlavor
			}
			nodesList = append(nodesList, node)
			// Only return Node private ips which doesn't support reduce
			if !Node.SupportReduce {
				ipsList = append(ipsList, Node.PrivateIp)
			}
		}
		volumeSize, err := strconv.Atoi(group.Volume.Size)
		if err != nil {
			return diag.FromErr(err)
		}
		mErr = multierror.Append(mErr, d.Set("volume_size", volumeSize))

		if specCode != "" {
			log.Printf("[DEBUG] node specCode: %s", specCode)
			mErr = multierror.Append(mErr, d.Set("flavor", specCode))
		}
	}
	mErr = multierror.Append(
		mErr,
		d.Set("nodes", nodesList),
		d.Set("private_ips", ipsList),
		d.Set("node_num", len(nodesList)),
	)

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": inst.BackupStrategy.StartTime,
		"keep_days":  inst.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	mErr = multierror.Append(mErr, d.Set("backup_strategy", backupStrategyList))

	// save geminidb tags
	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagmap := common.TagsToMap(resourceTags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for geminidb (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] error fetching tags of geminidb (%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeminiDBInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}

	instanceId := d.Id()

	_, err = instance.Delete(client, instanceId)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"normal", "abnormal", "creating", "createfail", "enlargefail", "data_disk_full"},
		Target:       []string{"deleted"},
		Refresh:      GeminiDBInstanceStateRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        15 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to be deleted: %s ",
			instanceId, err)
	}

	return nil
}

func resourceGaussDBCassandraInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	defaults := defaultValues{
		Mode:      "Cluster",
		dbType:    "cassandra",
		dbVersion: "3.11",
		logName:   "cassandra",
	}
	return resourceGeminiDBInstanceV3Update(ctx, d, meta, defaults)
}

func resourceGeminiDBInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}, defaults defaultValues) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.GeminiDBV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s ", err)
	}
	// update tags
	instanceId := d.Id()
	if d.HasChange("tags") {
		tagErr := common.UpdateResourceTags(client, d, "instances", instanceId)
		if tagErr != nil {
			return diag.Errorf("error updating tags of GeminiDB %q: %s", instanceId, tagErr)
		}
	}

	if d.HasChange("name") {
		updateNameOpts := instance.RenameInstanceOpts{
			InstanceID: instanceId,
			Name:       d.Get("name").(string),
		}

		err := instance.RenameInstance(client, updateNameOpts)
		if err != nil {
			return diag.Errorf("error updating name for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}
	}

	if d.HasChange("password") {
		updatePassOpts := instance.ResetPasswordOpts{
			InstanceId: instanceId,
			Password:   d.Get("password").(string),
		}

		err := instance.ResetPassword(client, updatePassOpts)
		if err != nil {
			return diag.Errorf("error updating password for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}
	}

	if d.HasChange("configuration_id") {
		instanceIds := []string{d.Id()}
		configId := d.Get("configuration_id").(string)
		applyOpts := template.ApplyOpts{
			ConfigId:    configId,
			InstanceIds: instanceIds,
		}

		ret, err := template.Apply(client, applyOpts)
		if err != nil || !ret.Success {
			return diag.Errorf("error updating configuration_id for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"SET_CONFIGURATION"},
			Target:     []string{"available"},
			Refresh:    GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "SET_CONFIGURATION"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
		}

		// Compare the target configuration and the instance configuration
		tmpl, err := template.Get(client, configId)
		if err != nil {
			return diag.Errorf("error fetching configuration %s: %s", configId, err)
		}
		configParams := tmpl.ConfigurationParameters
		log.Printf("[DEBUG] configuration parameters %#v", configParams)

		instanceConfig, err := template.GetInstanceParameters(client, instanceId)
		if err != nil {
			return diag.Errorf("error fetching instance configuration for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}
		instanceConfigParams := instanceConfig.ConfigurationParameters
		log.Printf("[DEBUG] instance configuration parameters %#v", instanceConfigParams)

		if len(configParams) != len(instanceConfigParams) {
			return diag.Errorf("error updating configuration for instance: %s", instanceId)
		}
		for i := range configParams {
			if !configParams[i].Readonly && configParams[i] != instanceConfigParams[i] {
				return diag.Errorf("error updating configuration for instance: %s", instanceId)
			}
		}
	}

	if d.HasChange("volume_size") {
		extendOpts := instance.ExtendVolumeOpts{
			InstanceId: instanceId,
			Size:       d.Get("volume_size").(int),
		}

		_, err = instance.ExtendVolume(client, extendOpts)
		if err != nil {
			return diag.Errorf("error extending gaussdb_%s_instance %s size: %s", defaults.logName, instanceId, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"RESIZE_VOLUME"},
			Target:     []string{"available"},
			Refresh:    GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "RESIZE_VOLUME"),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			MinTimeout: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
		}
	}

	if d.HasChange("node_num") {
		old, newnum := d.GetChange("node_num")
		if newnum.(int) > old.(int) {
			// Enlarge Nodes
			expandSize := newnum.(int) - old.(int)
			enlargeNodeOpts := instance.EnlargeNodeOpts{
				InstanceId: instanceId,
				Num:        expandSize,
			}

			log.Printf("[DEBUG] enlarge node options: %+v", enlargeNodeOpts)

			_, err = instance.EnlargeNode(client, enlargeNodeOpts)
			if err != nil {
				return diag.Errorf("error enlarging gaussdb_%s_instance %s node size: %s", defaults.logName, instanceId, err)
			}

			// 2. wait instance status
			stateConf := &resource.StateChangeConf{
				Pending:      []string{"GROWING"},
				Target:       []string{"available"},
				Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "GROWING"),
				Timeout:      d.Timeout(schema.TimeoutUpdate),
				Delay:        15 * time.Second,
				PollInterval: 20 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf(
					"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
			}
		}
		if newnum.(int) < old.(int) {
			if defaults.dbType == "influxdb" {
				return diag.Errorf("shrinking gaussdb %s instance node size is not allowed", defaults.logName)
			}
			// Reduce Nodes
			shrinkSize := old.(int) - newnum.(int)
			// the API accepts maxinum num of 10
			reduceNum := 10
			loopSize := shrinkSize / reduceNum
			lastNum := shrinkSize % reduceNum
			if lastNum > 0 {
				loopSize++
			}

			for i := 0; i < loopSize; i++ {
				if lastNum > 0 && (i == loopSize-1) {
					reduceNum = lastNum
				}
				reduceNodeOpts := instance.ReduceNodeOpts{
					InstanceID: instanceId,
					Num:        pointerto.Int(reduceNum),
				}
				log.Printf("[DEBUG] reduce node options: %+v", reduceNodeOpts)

				_, err = instance.ReduceNode(client, reduceNodeOpts)
				if err != nil {
					return diag.Errorf("error shrinking gaussdb %s instance %s node size: %s", defaults.logName, instanceId, err)
				}

				stateConf := &resource.StateChangeConf{
					Pending:      []string{"REDUCING"},
					Target:       []string{"available"},
					Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "REDUCING"),
					Timeout:      d.Timeout(schema.TimeoutUpdate),
					Delay:        15 * time.Second,
					PollInterval: 20 * time.Second,
				}

				_, err = stateConf.WaitForStateContext(ctx)
				if err != nil {
					return diag.Errorf(
						"error waiting for gaussdb %s instance %s to become ready: %s", defaults.logName, instanceId, err)
				}
			}
		}
	}

	if d.HasChange("flavor") {
		inst, err := GetInstanceByID(client, instanceId)
		if err != nil {
			return diag.Errorf(
				"error fetching gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}

		specCode := ""
		for _, action := range inst.Actions {
			if action == "RESIZE_FLAVOR" {
				// Wait here if the instance already in RESIZE_FLAVOR state
				stateConf := &resource.StateChangeConf{
					Pending:      []string{"RESIZE_FLAVOR"},
					Target:       []string{"available"},
					Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "RESIZE_FLAVOR"),
					Timeout:      d.Timeout(schema.TimeoutUpdate),
					PollInterval: 20 * time.Second,
				}

				_, err = stateConf.WaitForStateContext(ctx)
				if err != nil {
					return diag.Errorf(
						"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
				}

				inst, err := GetInstanceByID(client, instanceId)
				if err != nil {
					return diag.Errorf(
						"error fetching gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
				}

				// Fetch node flavor
				wrongFlavor := "Inconsistent Flavor"
				for _, group := range inst.Groups {
					for _, Node := range group.Nodes {
						if specCode == "" {
							specCode = Node.SpecCode
						} else if specCode != Node.SpecCode && specCode != wrongFlavor {
							specCode = wrongFlavor
						}
					}
				}
				break
			}
		}

		flavor := d.Get("flavor").(string)
		if specCode != flavor {
			log.Printf("[DEBUG] inconsistent node specCode: %s, flavor: %s", specCode, flavor)
			// Do resize action
			resizeOpts := instance.ResizeInstanceOpts{
				InstanceID:     d.Id(),
				TargetSpecCode: d.Get("flavor").(string),
			}

			_, err = instance.ResizeInstance(client, resizeOpts)
			if err != nil {
				return diag.Errorf("error resizing gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
			}

			stateConf := &resource.StateChangeConf{
				Pending:      []string{"RESIZE_FLAVOR"},
				Target:       []string{"available"},
				Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "RESIZE_FLAVOR"),
				Timeout:      d.Timeout(schema.TimeoutUpdate),
				PollInterval: 20 * time.Second,
			}

			_, err = stateConf.WaitForStateContext(ctx)
			if err != nil {
				return diag.Errorf(
					"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
			}
		}
	}

	if d.HasChange("security_group_id") {
		updateSgOpts := instance.ChangeSecGroupOpts{
			InstanceId:      d.Id(),
			SecurityGroupId: d.Get("security_group_id").(string),
		}

		_, err = instance.ChangeSecGroup(client, updateSgOpts)
		if err != nil {
			return diag.Errorf("error updating security group for gaussdb_%s_instance %s: %s", defaults.logName, instanceId, err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"MODIFY_SECURITYGROUP"},
			Target:       []string{"available"},
			Refresh:      GeminiDBInstanceUpdateRefreshFunc(client, instanceId, "MODIFY_SECURITYGROUP"),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"error waiting for gaussdb_%s_instance %s to become ready: %s", defaults.logName, instanceId, err)
		}
	}

	if d.HasChange("backup_strategy") {
		var updateOpts backup.SetBackupPolicyOpts
		backupRaw := d.Get("backup_strategy").([]interface{})
		rawMap := backupRaw[0].(map[string]interface{})
		keepDays := rawMap["keep_days"].(int)
		updateOpts.BackupPolicy.KeepDays = keepDays
		updateOpts.BackupPolicy.StartTime = rawMap["start_time"].(string)
		// Fixed to "1,2,3,4,5,6,7"
		updateOpts.BackupPolicy.Period = "1,2,3,4,5,6,7"
		updateOpts.InstanceId = instanceId
		log.Printf("[DEBUG] update backup_strategy: %#v", updateOpts)

		err = backup.SetBackupPolicy(client, updateOpts)
		if err != nil {
			return diag.Errorf("error updating backup_strategy: %s", err)
		}
	}

	return resourceGeminiDBInstanceV3Read(ctx, d, meta)
}

func GeminiDBInstanceUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID, state string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		inst, err := GetInstanceByID(client, instanceID)

		if err != nil {
			return nil, "", err
		}
		if inst.Id == "" {
			return inst, "deleted", nil
		}
		for _, action := range inst.Actions {
			if state == "REDUCING" {
				if action == "REDUCING" || action == "PERIOD_RESOURCE_DELETE" {
					return inst, state, nil
				}
			} else {
				if action == state {
					return inst, state, nil
				}
			}
		}

		return inst, "available", nil
	}
}

func GetInstanceByID(client *golangsdk.ServiceClient, instanceId string) (*instance.ListResult, error) {
	instances, err := instance.ListGeminiDB(client, instance.ListGeminiDBOpts{
		Id: instanceId,
	})
	if err != nil {
		return nil, err
	}

	if len(instances) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return &instances[0], nil
}
