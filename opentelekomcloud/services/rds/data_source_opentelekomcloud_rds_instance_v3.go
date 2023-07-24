package rds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRdsInstanceV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsInstanceV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"datastore_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"availability_zone": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"flavor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ha": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"timezone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_encryption_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
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
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role": {
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
			"private_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"public_ips": {
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

func dataSourceRdsInstanceV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	listOpts := instances.ListOpts{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		DataStoreType: d.Get("datastore_type").(string),
		VpcId:         d.Get("vpc_id").(string),
		SubnetId:      d.Get("subnet_id").(string),
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("id"); ok {
		listOpts.Id = v.(string)
	}

	instancesList, err := instances.List(client, listOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(instancesList.Instances) < 1 {
		return common.DataSourceTooFewDiag
	}

	if len(instancesList.Instances) > 1 {
		return common.DataSourceTooManyDiag
	}

	rdsInstance := instancesList.Instances[0]

	d.SetId(rdsInstance.Id)

	mErr := multierror.Append(nil,
		d.Set("name", rdsInstance.Name),
		d.Set("region", rdsInstance.Region),
		d.Set("port", rdsInstance.Port),
		d.Set("created", rdsInstance.Created),
		d.Set("vpc_id", rdsInstance.VpcId),
		d.Set("subnet_id", rdsInstance.SubnetId),
		d.Set("security_group_id", rdsInstance.SecurityGroupId),
		d.Set("flavor", rdsInstance.FlavorRef),
		d.Set("timezone", rdsInstance.TimeZone),
		d.Set("db_username", rdsInstance.DbUserName),
		d.Set("status", rdsInstance.Status),
		d.Set("created", rdsInstance.Created),
		d.Set("disk_encryption_id", rdsInstance.DiskEncryptionId),
		d.Set("datastore_type", rdsInstance.DataStore.Type),
		d.Set("datastore_version", rdsInstance.DataStore.Version),
		d.Set("updated", rdsInstance.Updated),
		d.Set("volume_type", rdsInstance.Volume.Type),
		d.Set("volume_size", rdsInstance.Volume.Size),
		d.Set("timezone", rdsInstance.TimeZone),
		d.Set("private_ips", rdsInstance.PrivateIps),
		d.Set("public_ips", rdsInstance.PublicIps),
	)

	// backup
	backup := make([]map[string]interface{}, 1)
	backup[0] = map[string]interface{}{
		"start_time": rdsInstance.BackupStrategy.StartTime,
		"keep_days":  rdsInstance.BackupStrategy.KeepDays,
	}
	if err = d.Set("backup_strategy", backup); err != nil {
		return fmterr.Errorf("error setting RDSv3 datastore opts: %s", err)
	}

	// nodes
	nodes := make([]map[string]interface{}, len(rdsInstance.Nodes))
	for i, v := range rdsInstance.Nodes {
		nodes[i] = map[string]interface{}{
			"id":                v.Id,
			"name":              v.Name,
			"role":              v.Role,
			"status":            v.Status,
			"availability_zone": v.AvailabilityZone,
		}
	}

	if err = d.Set("nodes", nodes); err != nil {
		return fmterr.Errorf("error setting RDSv3 datastore opts: %s", err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting RDSv3 rdsInstance backup fields: %w", err)
	}

	tagMap := common.TagsToMap(rdsInstance.Tags)
	mErr = multierror.Append(mErr, d.Set("tags", tagMap))

	return nil
}
