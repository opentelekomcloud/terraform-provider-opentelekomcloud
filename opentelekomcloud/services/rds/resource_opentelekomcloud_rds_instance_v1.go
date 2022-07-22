package rds

import (
	"context"
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
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v1/instances"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v1/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRdsInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         false,
				DiffSuppressFunc: common.SuppressRdsNameDiffs,
			},

			"datastore": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"PostgreSQL", "SQLServer", "MySQL"}, true),
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					}},
			},

			"flavorref": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"volume": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},
					}},
			},

			"availabilityzone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"vpc": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nics": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnetid": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					}},
			},

			"securitygroup": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					}},
			},

			"dbport": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"backupstrategy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"starttime": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"keepdays": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					}},
			},

			"dbrtpd": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ha": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"replicationmode": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"async", "sync", "semisync"}, true),
						},
					}},
			},
			"tag": common.TagsSchema(),
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
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
		},
	}
}

func resourceInstanceDataStore(d *schema.ResourceData) instances.DataStoreOps {
	var dataStore instances.DataStoreOps
	datastoreRaw := d.Get("datastore").([]interface{})
	log.Printf("[DEBUG] datastoreRaw: %+v", datastoreRaw)
	if len(datastoreRaw) == 1 {
		dataStore.Type = datastoreRaw[0].(map[string]interface{})["type"].(string)
		dataStore.Version = datastoreRaw[0].(map[string]interface{})["version"].(string)
	}
	log.Printf("[DEBUG] datastore: %+v", dataStore)
	return dataStore
}

func resourceInstanceVolume(d *schema.ResourceData) instances.VolumeOps {
	var volume instances.VolumeOps
	volumeRaw := d.Get("volume").([]interface{})
	log.Printf("[DEBUG] volumeRaw: %+v", volumeRaw)
	if len(volumeRaw) == 1 {
		volume.Type = volumeRaw[0].(map[string]interface{})["type"].(string)
		volume.Size = volumeRaw[0].(map[string]interface{})["size"].(int)
	}
	log.Printf("[DEBUG] volume: %+v", volume)
	return volume
}

func resourceInstanceNics(d *schema.ResourceData) instances.NicsOps {
	var nics instances.NicsOps
	nicsRaw := d.Get("nics").([]interface{})
	log.Printf("[DEBUG] nicsRaw: %+v", nicsRaw)
	if len(nicsRaw) == 1 {
		nics.SubnetId = nicsRaw[0].(map[string]interface{})["subnetid"].(string)
	}
	log.Printf("[DEBUG] nics: %+v", nics)
	return nics
}

func resourceInstanceSecurityGroup(d *schema.ResourceData) instances.SecurityGroupOps {
	var securityGroup instances.SecurityGroupOps
	SecurityGroupRaw := d.Get("securitygroup").([]interface{})
	log.Printf("[DEBUG] SecurityGroupOpsRaw: %+v", SecurityGroupRaw)
	if len(SecurityGroupRaw) == 1 {
		securityGroup.Id = SecurityGroupRaw[0].(map[string]interface{})["id"].(string)
	}
	log.Printf("[DEBUG] securityGroup: %+v", securityGroup)
	return securityGroup
}

func resourceInstanceBackupStrategy(d *schema.ResourceData) instances.BackupStrategyOps {
	var backupStrategy instances.BackupStrategyOps
	backupStrategyRaw := d.Get("backupstrategy").([]interface{})
	log.Printf("[DEBUG] backupStrategyRaw: %+v", backupStrategyRaw)
	if len(backupStrategyRaw) == 1 {
		backupStrategy.StartTime = backupStrategyRaw[0].(map[string]interface{})["starttime"].(string)
		backupStrategy.KeepDays = backupStrategyRaw[0].(map[string]interface{})["keepdays"].(int)
	} else {
		backupStrategy.StartTime = "00:00:00"
		backupStrategy.KeepDays = 0
	}
	log.Printf("[DEBUG] backupStrategy: %+v", backupStrategy)
	return backupStrategy
}

func resourceInstanceHa(d *schema.ResourceData) instances.HaOps {
	var ha instances.HaOps
	haRaw := d.Get("ha").([]interface{})
	log.Printf("[DEBUG] haRaw: %+v", haRaw)
	if len(haRaw) == 1 {
		ha.Enable = haRaw[0].(map[string]interface{})["enable"].(bool)
		if ha.Enable {
			ha.ReplicationMode = haRaw[0].(map[string]interface{})["replicationmode"].(string)
		}
	} else {
		ha.Enable = false
	}
	log.Printf("[DEBUG] ha: %+v", ha)
	return ha
}

func InstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}

		return instance, instance.Status, nil
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud rds client: %s ", err)
	}

	createOpts := instances.CreateOps{
		Name:             d.Get("name").(string),
		DataStore:        resourceInstanceDataStore(d),
		FlavorRef:        d.Get("flavorref").(string),
		Volume:           resourceInstanceVolume(d),
		Region:           config.GetRegion(d),
		AvailabilityZone: d.Get("availabilityzone").(string),
		Vpc:              d.Get("vpc").(string),
		Nics:             resourceInstanceNics(d),
		SecurityGroup:    resourceInstanceSecurityGroup(d),
		DbPort:           d.Get("dbport").(string),
		BackupStrategy:   resourceInstanceBackupStrategy(d),
		DbRtPd:           d.Get("dbrtpd").(string),
		Ha:               resourceInstanceHa(d),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	instance, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error getting instance from result: %s ", err)
	}
	log.Printf("[DEBUG] Create : instance %s: %#v", instance.ID, instance)

	d.SetId(instance.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    InstanceStateRefreshFunc(client, instance.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"Error waiting for instance (%s) to become ready: %s ",
			instance.ID, err)
	}

	if common.HasFilledOpt(d, "tag") {
		tagClient, err := config.RdsTagV1Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud rds tag client: %s ", err)
		}
		tagmap := d.Get("tag").(map[string]interface{})
		log.Printf("[DEBUG] Setting tag(key/value): %v", tagmap)
		for key, val := range tagmap {
			tagOpts := tags.CreateOpts{
				Key:   key,
				Value: val.(string),
			}
			err = tags.Create(tagClient, instance.ID, tagOpts).ExtractErr()
			if err != nil {
				log.Printf("[WARN] Error setting tag(key/value) of instance:%s, err=%s", instance.ID, err)
			}
		}
	}

	if instance.ID != "" {
		return resourceInstanceRead(ctx, d, meta)
	}
	return fmterr.Errorf("unexpected conversion error in resourceInstanceCreate. ")
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud rds client: %s", err)
	}

	instanceID := d.Id()
	instance, err := instances.Get(client, instanceID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "instance")
	}

	log.Printf("[DEBUG] Retrieved instance %s: %#v", instanceID, instance)

	if instance.Name != "" {
		nameList := strings.Split(instance.Name, "-"+instance.DataStore.Type)
		log.Printf("[DEBUG] Retrieved nameList %#v", nameList)
		if len(nameList) > 0 {
			_ = d.Set("name", nameList[0])
		}
	}
	mErr := multierror.Append(
		d.Set("hostname", instance.HostName),
		d.Set("type", instance.Type),
		d.Set("region", instance.Region),
		d.Set("availabilityzone", instance.AvailabilityZone),
		d.Set("vpc", instance.Vpc),
		d.Set("status", instance.Status),
	)

	nicsList := make([]map[string]interface{}, 0, 1)
	nics := map[string]interface{}{
		"subnetid": instance.Nics.SubnetId,
	}
	nicsList = append(nicsList, nics)
	log.Printf("[DEBUG] nicsList: %+v", nicsList)
	if err := d.Set("nics", nicsList); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving nics to Rds instance (%s): %s", d.Id(), err)
	}

	securitygroupList := make([]map[string]interface{}, 0, 1)
	securitygroup := map[string]interface{}{
		"id": instance.SecurityGroup.Id,
	}
	securitygroupList = append(securitygroupList, securitygroup)
	log.Printf("[DEBUG] securitygroupList: %+v", securitygroupList)
	if err := d.Set("securitygroup", securitygroupList); err != nil {
		return fmterr.Errorf("[DEBUG] Error saving securitygroup to Rds instance (%s): %s", d.Id(), err)
	}

	mErr = multierror.Append(mErr, d.Set("flavorref", instance.Flavor.Id))

	volumeList := make([]map[string]interface{}, 0, 1)
	volume := map[string]interface{}{
		"type": instance.Volume.Type,
		"size": instance.Volume.Size,
	}
	volumeList = append(volumeList, volume)
	if err := d.Set("volume", volumeList); err != nil {
		return fmterr.Errorf(
			"[DEBUG] Error saving volume to Rds instance (%s): %s", d.Id(), err)
	}

	mErr = multierror.Append(mErr, d.Set("dbport", strconv.Itoa(instance.DbPort)))

	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":    instance.DataStore.Type,
		"version": instance.DataStore.Version,
	}
	datastoreList = append(datastoreList, datastore)
	if err := d.Set("datastore", datastoreList); err != nil {
		return fmterr.Errorf(
			"[DEBUG] Error saving datastore to Rds instance (%s): %s", d.Id(), err)
	}

	mErr = multierror.Append(mErr,
		d.Set("updated", instance.Updated),
		d.Set("created", instance.Created),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	// set instance tag
	if _, ok := d.GetOk("tag"); ok {
		tagClient, err := config.RdsTagV1Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud rds tag client: %#v", err)
		}
		taglist, err := tags.Get(tagClient, d.Id()).Extract()
		if err != nil {
			return fmterr.Errorf("error fetching OpenTelekomCloud rds instance tags: %s", err)
		}

		tagmap := make(map[string]string)
		for _, val := range taglist.Tags {
			tagmap[val.Key] = val.Value
		}
		if err := d.Set("tag", tagmap); err != nil {
			return fmterr.Errorf("[DEBUG] Error saving tag to state for OpenTelekomCloud rds instance (%s): %s", d.Id(), err)
		}
	}
	return nil
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud rds client: %s ", err)
	}

	log.Printf("[DEBUG] Deleting Instance %s", d.Id())

	id := d.Id()
	result := instances.Delete(client, id)
	if result.Err != nil {
		return diag.FromErr(err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    InstanceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      15 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"Error waiting for instance (%s) to be deleted: %s ",
			id, err)
	}
	time.Sleep(80 * time.Second)
	log.Printf("[DEBUG] Successfully deleted instance %s", id)
	return nil
}

func instanceStateUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID string, size int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}
		log.Printf("[DEBUG] Updating instance.Volume : %+v", instance.Volume)
		if instance.Volume.Size == size {
			return instance, "UPDATED", nil
		}

		return instance, instance.Status, nil
	}
}

func instanceStateFlavorUpdateRefreshFunc(client *golangsdk.ServiceClient, instanceID string, _ string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := instances.Get(client, instanceID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return instance, "DELETED", nil
			}
			return nil, "", err
		}

		return instance, instance.Status, nil
	}
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error Updating OpenTelekomCloud rds client: %s ", err)
	}

	log.Printf("[DEBUG] Updating instances %s", d.Id())
	id := d.Id()

	if d.HasChange("volume") {
		var updateOpts instances.UpdateOps
		volume := make(map[string]interface{})
		volumeRaw := d.Get("volume").([]interface{})
		log.Printf("[DEBUG] volumeRaw: %+v", volumeRaw)
		if len(volumeRaw) == 1 {
			if m, ok := volumeRaw[0].(map[string]interface{}); ok {
				volume["size"] = m["size"].(int)
			}
		}
		log.Printf("[DEBUG] volume: %+v", volume)
		updateOpts.Volume = volume
		_, err = instances.UpdateVolumeSize(client, updateOpts, id).Extract()
		if err != nil {
			return fmterr.Errorf("error updating instance volume from result: %s ", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"ACTIVE"},
			Target:     []string{"UPDATED"},
			Refresh:    instanceStateUpdateRefreshFunc(client, id, updateOpts.Volume["size"].(int)),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf(
				"Error waiting for instance (%s) volume to be Updated: %s ",
				id, err)
		}
		log.Printf("[DEBUG] Successfully updated instance %s volume: %+v", id, volume)
	}

	if d.HasChange("flavorref") {
		var updateFlavorOpts instances.UpdateFlavorOps

		log.Printf("[DEBUG] Update flavorref: %s", d.Get("flavorref").(string))

		updateFlavorOpts.FlavorRef = d.Get("flavorref").(string)
		_, err = instances.UpdateFlavorRef(client, updateFlavorOpts, id).Extract()
		if err != nil {
			return fmterr.Errorf("error updating instance Flavor from result: %s ", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"MODIFYING"},
			Target:     []string{"ACTIVE"},
			Refresh:    instanceStateFlavorUpdateRefreshFunc(client, id, d.Get("flavorref").(string)),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      15 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmterr.Errorf(
				"Error waiting for instance (%s) flavor to be Updated: %s ",
				id, err)
		}
		log.Printf("[DEBUG] Successfully updated instance %s flavor: %s", id, d.Get("flavorref").(string))
	}

	if d.HasChange("backupstrategy") {
		var updatepolicyOpts instances.UpdatePolicyOps
		backupstrategyRaw := d.Get("backupstrategy").([]interface{})
		log.Printf("[DEBUG] backupstrategyRaw: %+v", backupstrategyRaw)
		if len(backupstrategyRaw) == 1 {
			if m, ok := backupstrategyRaw[0].(map[string]interface{}); ok {
				updatepolicyOpts.StartTime = m["starttime"].(string)
				updatepolicyOpts.KeepDays = m["keepdays"].(int)
			}
		}
		log.Printf("[DEBUG] updatepolicyOpts: %+v", updatepolicyOpts)
		_, err = instances.UpdatePolicy(client, updatepolicyOpts, id).Extract()
		if err != nil {
			return fmterr.Errorf("error updating instance policy from result: %s ", err)
		}

		log.Printf("[DEBUG] Successfully updated instance %s policy: %+v", id, updatepolicyOpts)
	}

	if d.HasChange("tag") {
		oraw, nraw := d.GetChange("tag")
		o := oraw.(map[string]interface{})
		n := nraw.(map[string]interface{})
		create, remove := diffTagsRDS(o, n)
		tagClient, err := config.RdsTagV1Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud rds tag client: %s ", err)
		}

		if len(remove) > 0 {
			for _, opts := range remove {
				err = tags.Delete(tagClient, id, opts).ExtractErr()
				if err != nil {
					log.Printf("[WARN] Error deleting tag(key/value) of instance:%s, err=%s", id, err)
				}
			}
		}
		if len(create) > 0 {
			for _, opts := range create {
				err = tags.Create(tagClient, id, opts).ExtractErr()
				if err != nil {
					log.Printf("[WARN] Error setting tag(key/value) of instance:%s, err=%s", id, err)
				}
			}
		}
	}

	log.Printf("[DEBUG] Successfully updated instance %s", id)
	d.SetId(id)
	return resourceInstanceRead(ctx, d, meta)
}

func diffTagsRDS(oldTags, newTags map[string]interface{}) ([]tags.CreateOptsBuilder, []tags.DeleteOptsBuilder) {
	var create []tags.CreateOptsBuilder
	var remove []tags.DeleteOptsBuilder
	for key, val := range oldTags {
		old, ok := newTags[key]
		if !ok || old.(string) != val.(string) {
			tagOpts := tags.DeleteOpts{
				Key: key,
			}
			remove = append(remove, tagOpts)
		}
	}
	for key, val := range newTags {
		old, ok := oldTags[key]
		if !ok || old.(string) != val.(string) {
			tagOpts := tags.CreateOpts{
				Key:   key,
				Value: val.(string),
			}
			create = append(create, tagOpts)
		}
	}
	return create, remove
}
