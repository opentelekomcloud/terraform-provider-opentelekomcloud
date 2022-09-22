package dcs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/configs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/instances"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/whitelists"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDcsInstanceV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsInstancesV1Create,
		ReadContext:   resourceDcsInstancesV1Read,
		UpdateContext: resourceDcsInstancesV1Update,
		DeleteContext: resourceDcsInstancesV1Delete,

		CustomizeDiff: validateEngine,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"save_days": {
				Type:       schema.TypeInt,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use `backup_policy` instead",
			},
			"backup_type": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Please use `backup_policy` instead",
			},
			"begin_at": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
			},
			"period_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"begin_at", "backup_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
			},
			"backup_at": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"period_type", "begin_at", "save_days", "backup_type"},
				Deprecated:   "Please use `backup_policy` instead",
				Elem:         &schema.Schema{Type: schema.TypeInt},
			},
			"backup_policy": {
				Type:          schema.TypeList,
				Optional:      true,
				ConflictsWith: []string{"backup_type", "begin_at", "period_type", "backup_at", "save_days"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"save_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"backup_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"begin_at": {
							Type:     schema.TypeString,
							Required: true,
						},
						"period_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"backup_at": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameter_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"parameter_value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"internal_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_whitelist": {
				Type:         schema.TypeBool,
				Computed:     true,
				Optional:     true,
				RequiredWith: []string{"whitelist"},
			},
			"whitelist": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				RequiredWith: []string{"enable_whitelist"},
			},
		},
	}
}

func formatAts(src []interface{}) []int {
	res := make([]int, len(src))
	for i, at := range src {
		res[i] = at.(int)
	}
	return res
}

func getInstanceBackupPolicy(d *schema.ResourceData) *instances.InstanceBackupPolicy {
	var instanceBackupPolicy *instances.InstanceBackupPolicy
	if _, ok := d.GetOk("backup_policy"); !ok { // deprecated branch
		backupAts := d.Get("backup_at").([]interface{})
		if len(backupAts) == 0 {
			return nil
		}
		instanceBackupPolicy = &instances.InstanceBackupPolicy{
			SaveDays:   d.Get("save_days").(int),
			BackupType: d.Get("backup_type").(string),
			PeriodicalBackupPlan: instances.PeriodicalBackupPlan{
				BeginAt:    d.Get("begin_at").(string),
				PeriodType: d.Get("period_type").(string),
				BackupAt:   formatAts(backupAts),
			},
		}
		return instanceBackupPolicy
	}

	backupPolicyList := d.Get("backup_policy").([]interface{})
	if len(backupPolicyList) == 0 {
		return nil
	}
	backupPolicy := backupPolicyList[0].(map[string]interface{})
	backupAts := backupPolicy["backup_at"].([]interface{})
	instanceBackupPolicy = &instances.InstanceBackupPolicy{
		SaveDays:   backupPolicy["save_days"].(int),
		BackupType: backupPolicy["backup_type"].(string),
		PeriodicalBackupPlan: instances.PeriodicalBackupPlan{
			BeginAt:    backupPolicy["begin_at"].(string),
			PeriodType: backupPolicy["period_type"].(string),
			BackupAt:   formatAts(backupAts),
		},
	}

	return instanceBackupPolicy
}

func getInstanceRedisConfiguration(d *schema.ResourceData) []configs.RedisConfig {
	redisConfigRaw := d.Get("configuration").([]interface{})
	if len(redisConfigRaw) == 0 {
		return nil
	}
	var redisConfigList []configs.RedisConfig
	for _, v := range redisConfigRaw {
		configuration := v.(map[string]interface{})
		redisConfig := configs.RedisConfig{
			ParamID:    configuration["parameter_id"].(string),
			ParamName:  configuration["parameter_name"].(string),
			ParamValue: configuration["parameter_value"].(string),
		}
		redisConfigList = append(redisConfigList, redisConfig)
	}

	return redisConfigList
}

func getInstanceWhitelistOpts(d *schema.ResourceData) whitelists.WhitelistOpts {
	var whitelistOpts whitelists.WhitelistOpts
	enabled := d.Get("enable_whitelist").(bool)
	whitelist := d.Get("whitelist").(*schema.Set).List()
	whitelistOpts.Enable = &enabled
	if len(whitelist) == 0 {
		whitelistOpts.Enable = &enabled
		whitelistOpts.Groups = []whitelists.WhitelistGroupOpts{}
		return whitelistOpts
	}

	for _, v := range whitelist {
		group := v.(map[string]interface{})
		groupOpts := whitelists.WhitelistGroupOpts{
			GroupName: group["group_name"].(string),
		}

		ipList := group["ip_list"].([]interface{})
		var refinedIpList []string

		for _, s := range ipList {
			refinedIpList = append(refinedIpList, s.(string))
		}
		groupOpts.IPList = refinedIpList

		whitelistOpts.Groups = append(whitelistOpts.Groups, groupOpts)
	}

	return whitelistOpts
}

func resourceDcsInstancesV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DcsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	noPasswordAccess := "true"
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		noPasswordAccess = "false"
	}
	createOpts := &instances.CreateOps{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Engine:               d.Get("engine").(string),
		EngineVersion:        d.Get("engine_version").(string),
		Capacity:             d.Get("capacity").(float64),
		NoPasswordAccess:     noPasswordAccess,
		Password:             d.Get("password").(string),
		AccessUser:           d.Get("access_user").(string),
		VPCID:                d.Get("vpc_id").(string),
		SecurityGroupID:      d.Get("security_group_id").(string),
		SubnetID:             d.Get("subnet_id").(string),
		AvailableZones:       common.GetAllAvailableZones(d),
		ProductID:            d.Get("product_id").(string),
		InstanceBackupPolicy: getInstanceBackupPolicy(d),
		MaintainBegin:        d.Get("maintain_begin").(string),
		MaintainEnd:          d.Get("maintain_end").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating DCS instance: %w", err)
	}
	log.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    dcsInstancesV1StateRefreshFunc(client, v.InstanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	updateOpts := configs.UpdateOpts{
		RedisConfigs: getInstanceRedisConfiguration(d),
	}
	if len(updateOpts.RedisConfigs) > 0 {
		if err := configs.Update(client, d.Id(), updateOpts).ExtractErr(); err != nil {
			return fmterr.Errorf("error updating redis configuration of DCS instance: %w", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"UPDATING"},
			Target:  []string{"SUCCESS"},
			Refresh: dcsInstanceV1ConfigStateRefreshFunc(client, d.Id()),
			Timeout: d.Timeout(schema.TimeoutCreate),
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for instance (%s) to delete: %w", d.Id(), err)
		}
	}

	if _, ok := d.GetOk("whitelist"); ok {
		whitelistOpts := getInstanceWhitelistOpts(d)
		if err := whitelists.Put(client, d.Id(), whitelistOpts).ExtractErr(); err != nil {
			return fmterr.Errorf("error updating redis whitelist of DCS instance: %w", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"UPDATING"},
			Target:  []string{"SUCCESS"},
			Refresh: dcsInstanceV1WhitelistRefreshFunc(client, d.Id(), nil),
			Timeout: d.Timeout(schema.TimeoutCreate),
			Delay:   4 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for instance (%s) to update whitelist: %w", d.Id(), err)
		}
	}

	return resourceDcsInstancesV1Read(ctx, d, meta)
}

func resourceDcsInstancesV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DcsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}
	v, err := instances.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] DCS instance %s: %+v", d.Id(), v)

	capacity := float64(v.Capacity)
	if v.Capacity == 0 {
		capacity, _ = strconv.ParseFloat(v.CapacityMinor, 32)
	}

	mErr := multierror.Append(
		d.Set("name", v.Name),
		d.Set("engine", v.Engine),
		d.Set("capacity", capacity),
		d.Set("used_memory", v.UsedMemory),
		d.Set("max_memory", v.MaxMemory),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("description", v.Description),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("internal_version", v.InternalVersion),
		d.Set("vpc_id", v.VPCID),
		d.Set("vpc_name", v.VPCName),
		d.Set("created_at", v.CreatedAt),
		d.Set("product_id", v.ProductID),
		d.Set("subnet_id", v.SubnetID),
		d.Set("subnet_name", v.SubnetName),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("order_id", v.OrderID),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("access_user", v.AccessUser),
		d.Set("ip", v.IP),
	)

	if v.EngineVersion == "3.0" {
		mErr = multierror.Append(
			d.Set("security_group_id", v.SecurityGroupID),
			d.Set("security_group_name", v.SecurityGroupName),
		)
	} else {
		w, err := whitelists.Get(client, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		if w.InstanceID != "" {
			var whitelistGroups []map[string]interface{}
			for _, group := range w.Groups {
				ipList := make([]string, len(group.IPList))
				copy(ipList, group.IPList)
				resourceMap := map[string]interface{}{
					"group_name": group.GroupName,
					"ip_list":    ipList,
				}
				whitelistGroups = append(whitelistGroups, resourceMap)
			}

			mErr = multierror.Append(
				d.Set("enable_whitelist", w.Enable),
				d.Set("whitelist", whitelistGroups),
			)
		}
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDcsInstancesV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DcsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}
	var updateOpts instances.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("maintain_begin") {
		updateOpts.MaintainBegin = d.Get("maintain_begin").(string)
	}
	if d.HasChange("maintain_end") {
		updateOpts.MaintainEnd = d.Get("maintain_end").(string)
	}
	if d.HasChange("security_group_id") {
		updateOpts.SecurityGroupID = d.Get("security_group_id").(string)
	}
	if d.HasChange("backup_policy") {
		updateOpts.InstanceBackupPolicy = getInstanceBackupPolicy(d)
	}

	err = instances.Update(client, d.Id(), updateOpts).Err
	if err != nil {
		return fmterr.Errorf("error updating DCS Instance: %w", err)
	}

	var updateConfigOpts configs.UpdateOpts
	if d.HasChange("configuration") {
		updateConfigOpts.RedisConfigs = getInstanceRedisConfiguration(d)
	}
	if len(updateConfigOpts.RedisConfigs) > 0 {
		if err := configs.Update(client, d.Id(), updateConfigOpts).ExtractErr(); err != nil {
			return fmterr.Errorf("error updating redis config of DCS instance: %w", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"UPDATING"},
			Target:  []string{"SUCCESS"},
			Refresh: dcsInstanceV1ConfigStateRefreshFunc(client, d.Id()),
			Timeout: d.Timeout(schema.TimeoutUpdate),
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for instance (%s) to delete: %w", d.Id(), err)
		}
	}

	if d.HasChange("enable_whitelist") || d.HasChange("whitelist") {
		getResp, err := whitelists.Get(client, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		whitelistOpts := getInstanceWhitelistOpts(d)
		if err := whitelists.Put(client, d.Id(), whitelistOpts).ExtractErr(); err != nil {
			return fmterr.Errorf("error updating redis whitelist of DCS instance: %w", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"UPDATING"},
			Target:  []string{"SUCCESS"},
			Refresh: dcsInstanceV1WhitelistRefreshFunc(client, d.Id(), getResp),
			Timeout: d.Timeout(schema.TimeoutCreate),
			Delay:   4 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for instance (%s) to update whitelist: %w", d.Id(), err)
		}
	}

	return resourceDcsInstancesV1Read(ctx, d, meta)
}

func resourceDcsInstancesV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DcsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	_, err = instances.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DCS instance")
	}

	err = instances.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting DCS instance: %w", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "RUNNING"},
		Target:     []string{"DELETED"},
		Refresh:    dcsInstancesV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to delete: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] DCS instance %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func dcsInstancesV1StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}

		return v, v.Status, nil
	}
}

func dcsInstanceV1ConfigStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := configs.List(client, instanceID).Extract()
		if err != nil {
			return nil, "", err
		}
		return v, v.ConfigStatus, nil
	}
}

func dcsInstanceV1WhitelistRefreshFunc(client *golangsdk.ServiceClient, instanceID string, getResp *whitelists.Whitelist) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := whitelists.Get(client, instanceID).Extract()
		if err != nil {
			return nil, "", err
		}
		// first logical condition - on resource creation, second - on update
		if getResp == nil && len(v.Groups) == 0 || v == getResp {
			return v, "UPDATING", nil
		}
		return v, "SUCCESS", nil
	}
}

func validateEngine(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	mErr := &multierror.Error{}
	engineVersion := d.Get("engine_version").(string)

	if _, ok := d.GetOk("security_group_id"); !ok && engineVersion == "3.0" {
		mErr = multierror.Append(mErr, fmt.Errorf("'security_group_id' should be set with engine_version==3.0"))
	} else if ok && engineVersion != "3.0" {
		mErr = multierror.Append(mErr, fmt.Errorf("DCS Redis 4.0 and 5.0 instances do not support security groups"))
	}

	if _, ok := d.GetOk("whitelist"); ok && engineVersion == "3.0" {
		mErr = multierror.Append(mErr, fmt.Errorf("DCS Redis 3.0 instance does not support whitelisting"))
	}

	return mErr.ErrorOrNil()
}
