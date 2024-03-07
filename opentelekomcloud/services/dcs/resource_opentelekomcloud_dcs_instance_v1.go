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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/configs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/lifecycle"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v1/others"
	dcsTags "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcs/v2/tags"
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
			StateContext: resourceDcsInstanceV1ImportState,
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
				ValidateFunc: validation.StringInSlice([]string{
					"3.0", "4.0", "5.0", "6.0",
				}, false),
			},
			"capacity": {
				Type:     schema.TypeFloat,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				ForceNew:  true,
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
			"private_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"no_password_access": {
				Type:     schema.TypeString,
				Computed: true,
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
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
			},
			"tags": common.TagsSchema(),
		},
	}
}

func buildDcsTags(tagsMap map[string]interface{}) []tags.ResourceTag {
	tagsList := make([]tags.ResourceTag, 0, len(tagsMap))
	for k, v := range tagsMap {
		tag := tags.ResourceTag{
			Key:   k,
			Value: v.(string),
		}
		tagsList = append(tagsList, tag)
	}
	return tagsList
}

func formatAts(src []interface{}) []int {
	res := make([]int, len(src))
	for i, at := range src {
		res[i] = at.(int)
	}
	return res
}

func getInstanceBackupPolicy(d *schema.ResourceData) *lifecycle.InstanceBackupPolicy {
	var instanceBackupPolicy *lifecycle.InstanceBackupPolicy
	if _, ok := d.GetOk("backup_policy"); !ok { // deprecated branch
		backupAts := d.Get("backup_at").([]interface{})
		if len(backupAts) == 0 {
			return nil
		}
		instanceBackupPolicy = &lifecycle.InstanceBackupPolicy{
			SaveDays:   d.Get("save_days").(int),
			BackupType: d.Get("backup_type").(string),
			PeriodicalBackupPlan: lifecycle.PeriodicalBackupPlan{
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
	instanceBackupPolicy = &lifecycle.InstanceBackupPolicy{
		SaveDays:   backupPolicy["save_days"].(int),
		BackupType: backupPolicy["backup_type"].(string),
		PeriodicalBackupPlan: lifecycle.PeriodicalBackupPlan{
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
	if d.Get("password").(string) != "" {
		noPasswordAccess = "false"
	}
	productId := d.Get("product_id").(string)
	var specCode string
	products, err := others.GetProducts(client)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Dcs get products : %+v", products)
	for _, pd := range products {
		if productId != "" && pd.ProductID == productId {
			specCode = pd.SpecCode
		}
	}
	createOpts := lifecycle.CreateOps{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Engine:               d.Get("engine").(string),
		EngineVersion:        d.Get("engine_version").(string),
		Capacity:             d.Get("capacity").(float64),
		NoPasswordAccess:     noPasswordAccess,
		Password:             d.Get("password").(string),
		VPCId:                d.Get("vpc_id").(string),
		SecurityGroupID:      d.Get("security_group_id").(string),
		SubnetID:             d.Get("subnet_id").(string),
		AvailableZones:       common.GetAllAvailableZones(d),
		SpecCode:             specCode,
		InstanceBackupPolicy: getInstanceBackupPolicy(d),
		MaintainBegin:        d.Get("maintain_begin").(string),
		MaintainEnd:          d.Get("maintain_end").(string),
		Tags:                 buildDcsTags(d.Get("tags").(map[string]interface{})),
	}

	if ip, ok := d.GetOk("private_ip"); ok {
		createOpts.PrivateIps = []string{ip.(string)}
		log.Printf("[DEBUG] private ip: %#v", createOpts.PrivateIps[0])
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	instanceID, err := lifecycle.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating DCS instance: %w", err)
	}
	log.Printf("[INFO] instance ID: %s", instanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    dcsInstancesV1StateRefreshFunc(client, instanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", instanceID, err)
	}

	// Store the instance ID now
	d.SetId(instanceID)

	updateOpts := configs.UpdateOpts{
		RedisConfigs: getInstanceRedisConfiguration(d),
	}
	if len(updateOpts.RedisConfigs) > 0 {
		if err := configs.Update(client, d.Id(), updateOpts); err != nil {
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
		if err := whitelists.Put(client, d.Id(), whitelistOpts); err != nil {
			return fmterr.Errorf("error updating redis whitelist of DCS instance: %w", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending: []string{"UPDATING"},
			Target:  []string{"SUCCESS"},
			Refresh: dcsInstanceV1WhitelistRefreshFunc(client, d.Id()),
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
	v, err := lifecycle.Get(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] DCS instance %s: %+v", d.Id(), v)

	capacity := float64(v.Capacity)
	if v.Capacity == 0 {
		capacity, _ = strconv.ParseFloat(v.CapacityMinor, 32)
	}
	var productId string
	products, err := others.GetProducts(client)
	if err != nil {
		return diag.FromErr(err)
	}
	for _, pd := range products {
		if pd.SpecCode == v.SpecCode {
			productId = pd.ProductID
		}
	}

	mErr := multierror.Append(
		d.Set("name", v.Name),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
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
		d.Set("product_id", productId),
		d.Set("subnet_id", v.SubnetID),
		d.Set("subnet_name", v.SubnetName),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("order_id", v.OrderID),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("ip", v.IP),
		d.Set("no_password_access", v.NoPasswordAccess),
	)

	if v.EngineVersion == "3.0" {
		mErr = multierror.Append(
			d.Set("security_group_id", v.SecurityGroupID),
			d.Set("security_group_name", v.SecurityGroupName),
		)
	} else {
		w, err := whitelists.Get(client, d.Id())
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
			if err := mErr.ErrorOrNil(); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if resourceTags, err := tags.Get(client, "instances", d.Id()).Extract(); err == nil {
		tagMap := common.TagsToMap(resourceTags)
		if err := d.Set("tags", tagMap); err != nil {
			return diag.Errorf("[DEBUG] error saving tags for OpenTelekomCloud DCS instance (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] fetching tags of OpenTelekomCloud DCS instance failed: %s", err)
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
	var updateOpts lifecycle.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = description
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

	err = lifecycle.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating DCS Instance: %w", err)
	}

	var updateConfigOpts configs.UpdateOpts
	if d.HasChange("configuration") {
		updateConfigOpts.RedisConfigs = getInstanceRedisConfiguration(d)
	}
	if len(updateConfigOpts.RedisConfigs) > 0 {
		if err := configs.Update(client, d.Id(), updateConfigOpts); err != nil {
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

	if d.HasChanges("enable_whitelist", "whitelist") {
		enable := d.Get("enable_whitelist").(bool)
		whitelistOpts := getInstanceWhitelistOpts(d)
		if err := whitelists.Put(client, d.Id(), whitelistOpts); err != nil {
			return fmterr.Errorf("error updating redis whitelist of DCS instance: %w", err)
		}
		stateConf := &resource.StateChangeConf{
			Pending:      []string{strconv.FormatBool(!enable)},
			Target:       []string{strconv.FormatBool(enable)},
			Refresh:      refreshForWhiteListEnableStatus(client, d.Id()),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        5 * time.Second,
			PollInterval: 5 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf("error waiting for instance (%s) to update whitelist: %w", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		oldVal, newVal := d.GetChange("tags")
		err = updateDcsTags(client, d.Id(), oldVal.(map[string]interface{}), newVal.(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
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

	_, err = lifecycle.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DCS instance")
	}

	err = lifecycle.Delete(client, d.Id())
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
		v, err := lifecycle.Get(client, instanceID)
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
		v, err := configs.List(client, instanceID)
		if err != nil {
			return nil, "", err
		}
		return v, v.ConfigStatus, nil
	}
}

func refreshForWhiteListEnableStatus(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := whitelists.Get(client, instanceID)
		if err != nil {
			return nil, "", err
		}
		return r, strconv.FormatBool(r.Enable), nil
	}
}

func dcsInstanceV1WhitelistRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := whitelists.Get(client, instanceID)
		if err != nil {
			return nil, "", err
		}
		return v, "SUCCESS", nil
	}
}

func validateEngine(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	engineVersion := d.Get("engine_version").(string)

	if _, ok := d.GetOk("whitelist"); ok && engineVersion == "3.0" {
		return fmt.Errorf("DCS Redis 3.0 instance does not support whitelisting")
	}

	return nil
}

func resourceDcsInstanceV1ImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*cfg.Config)
	client, err := config.DcsV1Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating ComputeV2 client: %w", err)
	}

	results := make([]*schema.ResourceData, 1)
	if diagRead := resourceDcsInstancesV1Read(ctx, d, meta); diagRead.HasError() {
		return nil, fmt.Errorf("error reading opentelekomcloud_dcs_instance_v1 %s: %s", d.Id(), diagRead[0].Summary)
	}

	instance, err := lifecycle.Get(client, d.Id())
	if err != nil {
		return nil, fmt.Errorf("unable to get instance %s: %s", d.Id(), err)
	}
	if err := d.Set("available_zones", instance.AvailableZones); err != nil {
		return nil, fmt.Errorf("error setting available zones")
	}
	if err := d.Set("used_memory", instance.UsedMemory); err != nil {
		return nil, fmt.Errorf("error setting used memory")
	}
	var backup []map[string]interface{}
	backupPolicy := make(map[string]interface{})
	backupPolicy["backup_type"] = instance.InstanceBackupPolicy.Policy.BackupType
	backupPolicy["save_days"] = instance.InstanceBackupPolicy.Policy.SaveDays
	backupPolicy["begin_at"] = instance.InstanceBackupPolicy.Policy.PeriodicalBackupPlan.BeginAt
	backupPolicy["period_type"] = instance.InstanceBackupPolicy.Policy.PeriodicalBackupPlan.PeriodType
	var backupAts []int
	backupAts = append(backupAts, instance.InstanceBackupPolicy.Policy.PeriodicalBackupPlan.BackupAt...)
	backupPolicy["backup_at"] = backupAts

	backup = append(backup, backupPolicy)
	if err := d.Set("backup_policy", backup); err != nil {
		return nil, fmt.Errorf("error setting backup policy")
	}

	results[0] = d

	return results, nil
}

func updateDcsTags(c *golangsdk.ServiceClient, id string, oldVal, newVal map[string]interface{}) error {
	if len(oldVal) > 0 {
		tagList := buildDcsTags(oldVal)
		err := dcsTags.Delete(c, id, tagList)
		if err != nil {
			return err
		}
	}
	if len(newVal) > 0 {
		tagList := buildDcsTags(newVal)
		err := dcsTags.Create(c, id, tagList)
		if err != nil {
			return err
		}
	}
	return nil
}
