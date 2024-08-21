package dds

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/backups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/instances"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/job"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsInstanceV3Create,
		ReadContext:   resourceDdsInstanceV3Read,
		UpdateContext: resourceDdsInstanceV3Update,
		DeleteContext: resourceDdsInstanceV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(90 * time.Minute),
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateName,
			},
			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"security_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"disk_encryption_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"num": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 16),
						},
						"storage": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntDivisibleBy(10),
						},
						"spec_code": {
							Type:     schema.TypeString,
							Required: true,
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
							Required: true,
						},
						"period": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"db_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pay_mode": {
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
						"role": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDdsDataStore(d *schema.ResourceData) instances.DataStore {
	var dataStore instances.DataStore
	datastoreRaw := d.Get("datastore").([]interface{})
	log.Printf("[DEBUG] datastoreRaw: %+v", datastoreRaw)
	if len(datastoreRaw) == 1 {
		dataStore.Type = datastoreRaw[0].(map[string]interface{})["type"].(string)
		dataStore.Version = datastoreRaw[0].(map[string]interface{})["version"].(string)
		dataStore.StorageEngine = datastoreRaw[0].(map[string]interface{})["storage_engine"].(string)
	}
	log.Printf("[DEBUG] datastore: %+v", dataStore)
	return dataStore
}

func resourceDdsFlavors(d *schema.ResourceData) []instances.Flavor {
	var flavors []instances.Flavor
	flavorRaw := d.Get("flavor").([]interface{})
	log.Printf("[DEBUG] flavorRaw: %+v", flavorRaw)
	for i := range flavorRaw {
		flavor := flavorRaw[i].(map[string]interface{})
		flavorReq := instances.Flavor{
			Type:     flavor["type"].(string),
			Num:      flavor["num"].(int),
			Storage:  flavor["storage"].(string),
			Size:     flavor["size"].(int),
			SpecCode: flavor["spec_code"].(string),
		}
		flavors = append(flavors, flavorReq)
	}
	log.Printf("[DEBUG] flavors: %+v", flavors)
	return flavors
}

func resourceDdsBackupStrategy(d *schema.ResourceData) instances.BackupStrategy {
	var backupStrategy instances.BackupStrategy
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	startTime := "00:00-01:00"
	keepDays := 7
	period := "1,2,3,4,5,6,7"
	if len(backupStrategyRaw) == 1 {
		startTime = backupStrategyRaw[0].(map[string]interface{})["start_time"].(string)
		keepDays = backupStrategyRaw[0].(map[string]interface{})["keep_days"].(int)
		if periodRaw := backupStrategyRaw[0].(map[string]interface{})["period"].(string); periodRaw != "" {
			period = periodRaw
		}
	}
	backupStrategy.KeepDays = &keepDays
	backupStrategy.StartTime = startTime
	backupStrategy.Period = period
	return backupStrategy
}

func instanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := instances.ListInstanceOpts{
			Id: instanceID,
		}
		instancesList, err := instances.List(client, opts)
		if err != nil {
			return nil, "", err
		}

		if instancesList.TotalCount == 0 || len(instancesList.Instances) == 0 {
			var instance instances.InstanceResponse
			return instance, "deleted", nil
		}
		ddsInstances := instancesList.Instances

		status := ddsInstances[0].Status
		actions := ddsInstances[0].Actions
		// wait for updating
		if status == "normal" && len(actions) > 0 {
			status = "updating"
		}
		return ddsInstances[0], status, nil
	}
}

func resourceDdsInstanceV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := instances.CreateOpts{
		Name:             d.Get("name").(string),
		DataStore:        resourceDdsDataStore(d),
		Region:           config.GetRegion(d),
		AvailabilityZone: d.Get("availability_zone").(string),
		VpcId:            d.Get("vpc_id").(string),
		SubnetId:         d.Get("subnet_id").(string),
		SecurityGroupId:  d.Get("security_group_id").(string),
		Password:         d.Get("password").(string),
		DiskEncryptionId: d.Get("disk_encryption_id").(string),
		Mode:             d.Get("mode").(string),
		Flavor:           resourceDdsFlavors(d),
		BackupStrategy:   resourceDdsBackupStrategy(d),
		Tags:             ddsTags(d),
	}
	if d.Get("ssl").(bool) {
		createOpts.Ssl = "1"
	} else {
		createOpts.Ssl = "0"
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	if val, ok := d.GetOk("port"); ok {
		createOpts.Port = strconv.Itoa(val.(int))
	}

	instance, err := instances.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting instance from result: %w", err)
	}
	log.Printf("[DEBUG] Create instance %s: %#v", instance.Id, instance)

	d.SetId(instance.Id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating"},
		Target:     []string{"normal"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to become ready: %w", d.Id(), err)
	}

	// since the POST method has no `period`, update backup strategy for it
	backupStrategyRaw := d.Get("backup_strategy").([]interface{})
	if len(backupStrategyRaw) == 1 {
		period := backupStrategyRaw[0].(map[string]interface{})["period"].(string)
		if period != "" && !isEqualPeriod(period, "1,2,3,4,5,6,7") {
			if err := createBackupStrategy(ctx, client, d); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceDdsInstanceV3Read(clientCtx, d, meta)
}

func ddsTags(d *schema.ResourceData) []tags.ResourceTag {
	vaultTags := d.Get("tags").(map[string]interface{})
	var tagSlice []tags.ResourceTag
	for k, v := range vaultTags {
		tagSlice = append(tagSlice, tags.ResourceTag{Key: k, Value: v.(string)})
	}
	return tagSlice
}

func resourceDdsInstanceV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	listOpts := instances.ListInstanceOpts{
		Id: d.Id(),
	}
	instancesList, err := instances.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching DDS instance: %w", err)
	}
	if instancesList.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", d.Id())
		d.SetId("")
		return nil
	}
	ddsInstances := instancesList.Instances
	instance := ddsInstances[0]

	log.Printf("[DEBUG] Retrieved instance %s: %#v", d.Id(), instance)

	tagsMap := make(map[string]string)
	for _, tag := range instance.Tags {
		tagsMap[tag.Key] = tag.Value
	}

	mErr := multierror.Append(nil,
		d.Set("region", instance.Region),
		d.Set("name", instance.Name),
		d.Set("vpc_id", instance.VpcId),
		d.Set("subnet_id", instance.SubnetId),
		d.Set("security_group_id", instance.SecurityGroupId),
		d.Set("disk_encryption_id", instance.DiskEncryptionId),
		d.Set("mode", instance.Mode),
		d.Set("db_username", instance.DbUserName),
		d.Set("status", instance.Status),
		d.Set("port", instance.Port),
		d.Set("pay_mode", instance.PayMode),
		d.Set("created_at", instance.Created),
		d.Set("updated_at", instance.Updated),
		d.Set("time_zone", instance.TimeZone),
		d.Set("tags", tagsMap),
	)

	sslEnable := true
	if instance.Ssl == 0 {
		sslEnable = false
	}
	mErr = multierror.Append(
		mErr,
		d.Set("ssl", sslEnable),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting DDSv3 multiple opts: %w", err)
	}

	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":           instance.DataStore.Type,
		"version":        instance.DataStore.Version,
		"storage_engine": instance.Engine,
	}
	datastoreList = append(datastoreList, datastore)
	if err := d.Set("datastore", datastoreList); err != nil {
		return fmterr.Errorf("error setting DDSv3 datastore opts: %w", err)
	}

	backupStrategyResp, err := backups.GetBackupPolicy(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": backupStrategyResp.StartTime,
		"keep_days":  backupStrategyResp.KeepDays,
		"period":     backupStrategyResp.Period,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	if err := d.Set("backup_strategy", backupStrategyList); err != nil {
		return fmterr.Errorf("error setting DDSv3 backup_strategy opts: %w", err)
	}

	// save nodes attribute
	err = d.Set("nodes", flattenDdsInstanceV3Nodes(instance))
	if err != nil {
		return fmterr.Errorf("error setting nodes of DDSv3 instance: %w", err)
	}

	return nil
}

func resourceDdsInstanceV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	if d.HasChange("name") {
		err := instances.UpdateName(client, instances.UpdateNameOpt{
			InstanceId:      d.Id(),
			NewInstanceName: d.Get("name").(string),
		})
		if err != nil {
			return fmterr.Errorf("error updating instance name: %w", err)
		}
	}

	if d.HasChange("password") {
		err := instances.ChangePassword(client, instances.ChangePasswordOpt{
			InstanceId: d.Id(),
			UserPwd:    d.Get("password").(string),
		})
		if err != nil {
			return fmterr.Errorf("error updating instance password: %w", err)
		}
	}

	if d.HasChange("ssl") {
		opt := instances.SSLOpt{
			InstanceId: d.Id(),
		}
		if d.Get("ssl").(bool) {
			opt.SSL = "1"
		} else {
			opt.SSL = "0"
		}
		_, err = instances.SwitchSSL(client, opt)
		if err != nil {
			return fmterr.Errorf("error updating ssl: %w", err)
		}
	}

	if d.HasChange("security_group_id") {
		_, err = instances.ModifySG(client, instances.ModifySGOpt{
			InstanceId:      d.Id(),
			SecurityGroupId: d.Get("security_group_id").(string),
		})
		if err != nil {
			return fmterr.Errorf("error updating security group: %w", err)
		}
	}

	if err := resourceDdsInstanceWaitUpdate(ctx, client, d); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("tags") {
		tagErr := common.UpdateResourceTags(client, d, "instances", d.Id())
		if tagErr != nil {
			return fmterr.Errorf("error updating tags of DDS instance:%s, err:%s", d.Id(), tagErr)
		}
	}

	if d.HasChange("flavor") {
		for i := range d.Get("flavor").([]interface{}) {
			volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
			numIndex := fmt.Sprintf("flavor.%d.num", i)
			specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)

			if d.HasChange(volumeSizeIndex) {
				err := flavorSizeUpdate(ctx, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			if d.HasChange(specCodeIndex) {
				err := flavorSpecCodeUpdate(ctx, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
			if d.HasChange(numIndex) {
				err := flavorNumUpdate(ctx, client, d, i)
				if err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	if d.HasChange("backup_strategy") {
		backupStrategy := resourceDdsBackupStrategy(d)
		err := backups.SetBackupPolicy(client, backups.ModifyBackupPolicyOpts{
			InstanceId:   d.Id(),
			BackupPolicy: &backupStrategy,
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("port") {
		retryFunc := func() (interface{}, bool, error) {
			resp, err := instances.ModifyPort(client, instances.ModifyPortOpt{
				Port:       d.Get("port").(int),
				InstanceId: d.Id(),
			})
			retry, err := handleMultiOperationsError(err)
			return resp, retry, err
		}
		r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"normal"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			return diag.Errorf("error updating database access port: %s", err)
		}
		jobId := r.(*string)
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"Running"},
			Target:       []string{"Completed"},
			Refresh:      JobStateRefreshFunc(client, *jobId),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			PollInterval: 10 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for the job (%s) completed: %s ", *jobId, err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceDdsInstanceV3Read(clientCtx, d, meta)
}

func resourceDdsInstanceV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	_, err = instances.Delete(client, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"normal", "abnormal", "frozen", "createfail", "enlargefail", "data_disk_full"},
		Target:     []string{"deleted"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to be deleted: %w", d.Id(), err)
	}
	log.Printf("[DEBUG] Successfully deleted instance %s", d.Id())
	return nil
}

func flattenDdsInstanceV3Nodes(dds instances.InstanceResponse) interface{} {
	var nodesList []map[string]interface{}
	for _, group := range dds.Groups {
		groupType := group.Type
		for _, Node := range group.Nodes {
			node := map[string]interface{}{
				"type":       groupType,
				"id":         Node.Id,
				"name":       Node.Name,
				"role":       Node.Role,
				"status":     Node.Status,
				"private_ip": Node.PrivateIP,
				"public_ip":  Node.PublicIP,
			}
			nodesList = append(nodesList, node)
		}
	}
	return nodesList
}

func flavorSizeUpdate(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	oldSizeRaw, newSizeRaw := d.GetChange(volumeSizeIndex)
	oldSize := oldSizeRaw.(int)
	newSize := newSizeRaw.(int)
	if newSize < oldSize {
		return fmt.Errorf("error updating instance: the new size(%d) must be greater than the old size(%d)", newSize, oldSize)
	}
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "replica" && groupType != "single" && groupType != "shard" {
		return fmt.Errorf("error updating instance: %s does not support scaling up storage space", groupType)
	}

	if groupType == "shard" {
		groupIDs, err := getDdsInstanceV3ShardGroupID(client, d)
		if err != nil {
			return err
		}

		for _, groupID := range groupIDs {
			updateVolumeOpts := instances.ScaleStorageOpt{
				GroupId:    groupID,
				Size:       fmt.Sprintf("%d", newSize),
				InstanceId: d.Id(),
			}

			_, err := instances.ScaleStorage(client, updateVolumeOpts)
			if err != nil {
				return err
			}

			if err := resourceDdsInstanceWaitUpdate(ctx, client, d); err != nil {
				return err
			}
		}
	} else {
		updateVolumeOpts := instances.ScaleStorageOpt{
			Size:       fmt.Sprintf("%d", newSize),
			InstanceId: d.Id(),
		}

		_, err := instances.ScaleStorage(client, updateVolumeOpts)
		if err != nil {
			return err
		}

		if err := resourceDdsInstanceWaitUpdate(ctx, client, d); err != nil {
			return err
		}
	}
	return nil
}

func flavorNumUpdate(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType != "mongos" && groupType != "shard" {
		return fmt.Errorf("error updating instance: %s does not support adding nodes", groupType)
	}
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	volumeSizeIndex := fmt.Sprintf("flavor.%d.size", i)
	volumeSize := d.Get(volumeSizeIndex).(int)
	numIndex := fmt.Sprintf("flavor.%d.num", i)
	oldNumRaw, newNumRaw := d.GetChange(numIndex)
	oldNum := oldNumRaw.(int)
	newNum := newNumRaw.(int)
	if newNum < oldNum {
		return fmt.Errorf("error updating instance: the new num(%d) must be greater than the old num(%d)", newNum, oldNum)
	}

	var updateNodeNumOpts instances.AddNodeOpts
	if groupType == "mongos" {
		updateNodeNumOpts = instances.AddNodeOpts{
			Type:       groupType,
			SpecCode:   d.Get(specCodeIndex).(string),
			Num:        newNum - oldNum,
			InstanceId: d.Id(),
		}
	} else {
		volume := instances.VolumeNode{
			Size: volumeSize,
		}
		updateNodeNumOpts = instances.AddNodeOpts{
			Type:       groupType,
			SpecCode:   d.Get(specCodeIndex).(string),
			Num:        newNum - oldNum,
			Volume:     &volume,
			InstanceId: d.Id(),
		}
	}

	_, err := instances.AddNode(client, updateNodeNumOpts)
	if err != nil {
		return err
	}

	if err := resourceDdsInstanceWaitUpdate(ctx, client, d); err != nil {
		return err
	}

	return nil
}

func flavorSpecCodeUpdate(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, i int) error {
	specCodeIndex := fmt.Sprintf("flavor.%d.spec_code", i)
	groupTypeIndex := fmt.Sprintf("flavor.%d.type", i)
	groupType := d.Get(groupTypeIndex).(string)
	if groupType == "config" {
		return fmt.Errorf("error updating instance: %s does not support updating spec_code", groupType)
	}
	switch groupType {
	case "mongos":
		nodeIDs, err := getDdsInstanceV3MongosNodeID(client, d)
		if err != nil {
			return err
		}
		for _, ID := range nodeIDs {
			updateSpecOpts := instances.ModifySpecOpt{
				TargetType:     "mongos",
				TargetId:       ID,
				TargetSpecCode: d.Get(specCodeIndex).(string),
				InstanceId:     d.Id(),
			}

			_, err = instances.ModifySpec(client, updateSpecOpts)
			if err != nil {
				return err
			}

			if err := resourceDdsInstanceWaitUpdate(ctx, client, d); err != nil {
				return err
			}
		}
	case "shard":
		groupIDs, err := getDdsInstanceV3ShardGroupID(client, d)
		if err != nil {
			return err
		}

		for _, ID := range groupIDs {
			updateSpecOpts := instances.ModifySpecOpt{
				TargetType:     "shard",
				TargetId:       ID,
				TargetSpecCode: d.Get(specCodeIndex).(string),
				InstanceId:     d.Id(),
			}

			_, err = instances.ModifySpec(client, updateSpecOpts)
			if err != nil {
				return err
			}

			if err := resourceDdsInstanceWaitUpdate(ctx, client, d); err != nil {
				return err
			}
		}
	default:
		updateSpecOpts := instances.ModifySpecOpt{
			TargetId:       d.Id(),
			TargetSpecCode: d.Get(specCodeIndex).(string),
			InstanceId:     d.Id(),
		}

		_, err := instances.ModifySpec(client, updateSpecOpts)
		if err != nil {
			return err
		}

		if err := resourceDdsInstanceWaitUpdate(ctx, client, d); err != nil {
			return err
		}
	}
	return nil
}

func getDdsInstanceV3ShardGroupID(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]string, error) {
	groupIDs := make([]string, 0)

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	ddsInstances, err := instances.List(client, opts)
	if err != nil {
		return groupIDs, fmt.Errorf("error fetching DDS instance: %s", err)
	}
	if ddsInstances.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		return groupIDs, nil
	}

	log.Printf("[DEBUG] Retrieved instance, id: %s", instanceID)

	for _, group := range ddsInstances.Instances[0].Groups {
		if group.Type == "shard" {
			groupIDs = append(groupIDs, group.Id)
		}
	}

	return groupIDs, nil
}

func getDdsInstanceV3MongosNodeID(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]string, error) {
	nodeIDs := make([]string, 0)

	instanceID := d.Id()
	opts := instances.ListInstanceOpts{
		Id: instanceID,
	}
	ddsInstances, err := instances.List(client, opts)
	if err != nil {
		return nodeIDs, fmt.Errorf("error fetching DDS instance: %s", err)
	}

	if ddsInstances.TotalCount == 0 {
		log.Printf("[WARN] DDS instance (%s) was not found", instanceID)
		return nodeIDs, nil
	}

	log.Printf("[DEBUG] Retrieved instance, id: %s", instanceID)

	for _, group := range ddsInstances.Instances[0].Groups {
		if group.Type == "mongos" {
			for _, node := range group.Nodes {
				nodeIDs = append(nodeIDs, node.Id)
			}
		}
	}

	return nodeIDs, nil
}

func resourceDdsInstanceWaitUpdate(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating", "updating"},
		Target:     []string{"normal"},
		Refresh:    instanceStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      15 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for instance (%s) to become ready: %w", d.Id(), err)
	}
	return nil
}

func isEqualPeriod(old, new string) bool {
	if len(old) != len(new) {
		return false
	}
	oldArray := strings.Split(old, ",")
	newArray := strings.Split(new, ",")
	sort.Strings(oldArray)
	sort.Strings(newArray)

	return reflect.DeepEqual(oldArray, newArray)
}

func createBackupStrategy(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	retryFunc := func() (interface{}, bool, error) {
		strategy := resourceDdsBackupStrategy(d)
		opts := backups.ModifyBackupPolicyOpts{
			InstanceId:   d.Id(),
			BackupPolicy: &strategy,
		}
		err := backups.SetBackupPolicy(client, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("error creating backup strategy of the DDS instance: %s ", err)
	}

	return nil
}

func JobStateRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := job.Get(client, jobId)
		if err != nil {
			return nil, "", err
		}

		return resp, resp.Status, nil
	}
}
