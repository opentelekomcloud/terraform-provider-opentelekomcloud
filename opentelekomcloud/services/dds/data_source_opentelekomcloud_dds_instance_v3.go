package dds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceDdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsInstanceV3Read,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_encryption_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keep_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"db_username": {
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
		},
	}
}

func dataSourceDdsInstanceV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	ddsClient, err := config.DdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DDS client: %s", err)
	}

	listOpts := instances.ListInstanceOpts{
		Id:            d.Get("instance_id").(string),
		Name:          d.Get("name").(string),
		DataStoreType: d.Get("datastore_type").(string),
		VpcId:         d.Get("vpc_id").(string),
		SubnetId:      d.Get("subnet_id").(string),
	}

	instancesList, err := instances.List(ddsClient, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching DDS instance: %s", err)
	}
	if len(instancesList.Instances) < 1 {
		return fmterr.Errorf("your query returned no results. Please change your search criteria and try again")
	}

	if len(instancesList.Instances) > 1 {
		return fmterr.Errorf("your query returned more than one result. Please try a more specific search criteria")
	}

	ddsInstance := instancesList.Instances[0]

	d.SetId(ddsInstance.Id)
	mErr := multierror.Append(nil,
		d.Set("instance_id", ddsInstance.Id),
		d.Set("name", ddsInstance.Name),
		d.Set("region", config.GetRegion(d)),
		d.Set("vpc_id", ddsInstance.VpcId),
		d.Set("subnet_id", ddsInstance.SubnetId),
		d.Set("security_group_id", ddsInstance.SecurityGroupId),
		d.Set("disk_encryption_id", ddsInstance.DiskEncryptionId),
		d.Set("mode", ddsInstance.Mode),
		d.Set("db_username", ddsInstance.DbUserName),
		d.Set("status", ddsInstance.Status),
		d.Set("port", ddsInstance.Port),
		d.Set("pay_mode", ddsInstance.PayMode),
		d.Set("datastore_type", ddsInstance.DataStore.Type),
	)
	sslEnable := true
	if ddsInstance.Ssl == 0 {
		sslEnable = false
	}
	mErr = multierror.Append(
		mErr,
		d.Set("ssl", sslEnable),
	)
	datastoreList := make([]map[string]interface{}, 0, 1)
	datastore := map[string]interface{}{
		"type":           ddsInstance.DataStore.Type,
		"version":        ddsInstance.DataStore.Version,
		"storage_engine": ddsInstance.Engine,
	}
	datastoreList = append(datastoreList, datastore)
	if err = d.Set("datastore", datastoreList); err != nil {
		return fmterr.Errorf("error setting DDSv3 datastore opts: %s", err)
	}

	backupStrategyList := make([]map[string]interface{}, 0, 1)
	backupStrategy := map[string]interface{}{
		"start_time": ddsInstance.BackupStrategy.StartTime,
		"keep_days":  ddsInstance.BackupStrategy.KeepDays,
	}
	backupStrategyList = append(backupStrategyList, backupStrategy)
	if err = d.Set("backup_strategy", backupStrategyList); err != nil {
		return fmterr.Errorf("error setting DDSv3 backup_strategy opts: %s", err)
	}

	err = d.Set("nodes", flattenDdsInstanceV3Nodes(ddsInstance))
	if err != nil {
		return fmterr.Errorf("error setting nodes of DDSv3 instance: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}
