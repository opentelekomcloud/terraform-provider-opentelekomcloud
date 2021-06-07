package rds

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRdsReadReplicaV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsReadReplicaV3Create,
		ReadContext:   resourceRdsReadReplicaV3Read,
		UpdateContext: resourceRdsReadReplicaV3Update,
		DeleteContext: resourceRdsReadReplicaV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateName,
			},
			"replica_of_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"flavor_ref": {
				Type:     schema.TypeString,
				Required: true,
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
							Computed: true,
						},
						"disk_encryption_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"public_ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"db": {
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
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceRdsReadReplicaV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	opts := &instances.CreateReplicaOpts{
		Name:             d.Get("name").(string),
		ReplicaOfId:      d.Get("replica_of_id").(string),
		DiskEncryptionId: d.Get("volume.0.disk_encryption_id").(string),
		FlavorRef:        d.Get("flavor_ref").(string),
		Volume: &instances.Volume{
			Type: d.Get("volume.0.type").(string),
		},
		Region:           d.Get("region").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}
	job, err := instances.CreateReplica(client, opts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating read replica: %w", err)
	}
	d.SetId(job.Instance.Id)

	timeoutSeconds := d.Timeout(schema.TimeoutCreate).Seconds()
	err = instances.WaitForJobCompleted(client, int(timeoutSeconds), job.JobId)
	if err != nil {
		return fmterr.Errorf("error waiting for read replica to complete creation: %w", err)
	}

	return resourceRdsReadReplicaV3Read(ctx, d, meta)
}

func resourceRdsReadReplicaV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	replica, err := GetRdsInstance(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error finding RDS instance: %w", err)
	}
	if replica == nil {
		d.SetId("")
		return nil
	}

	var baseID string
	for _, inst := range replica.RelatedInstance {
		if inst.Type == "replica_of" {
			baseID = inst.Id
		}
	}

	az := ""
	if len(replica.Nodes) > 0 {
		az = replica.Nodes[0].AvailabilityZone
	}

	mErr := multierror.Append(nil,
		d.Set("name", replica.Name),
		d.Set("availability_zone", az),
		d.Set("flavor_ref", replica.FlavorRef),
		d.Set("replica_of_id", baseID),
		d.Set("security_group_id", replica.SecurityGroupId),
		d.Set("subnet_id", replica.SubnetId),
		d.Set("vpc_id", replica.VpcId),
		d.Set("private_ips", replica.PrivateIps),
		d.Set("region", replica.Region),
		d.Set("public_ips", replica.PublicIps),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting replica fields: %w", err)
	}

	volume := map[string]interface{}{
		"type":               replica.Volume.Type,
		"size":               replica.Volume.Size,
		"disk_encryption_id": replica.DiskEncryptionId,
	}
	if err = d.Set("volume", []interface{}{volume}); err != nil {
		return fmterr.Errorf("error setting replica volume: %w", err)
	}

	dbInfo := map[string]interface{}{
		"type":      replica.DataStore.Type,
		"version":   replica.DataStore.Version,
		"port":      replica.Port,
		"user_name": replica.DbUserName,
	}
	if err = d.Set("db", []interface{}{dbInfo}); err != nil {
		return fmterr.Errorf("error setting replica db info: %w", err)
	}

	return nil
}

func resourceRdsReadReplicaV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	if d.HasChange("flavor_ref") {
		resizeOpts := instances.ResizeFlavorOpts{
			ResizeFlavor: &instances.SpecCode{
				Speccode: d.Get("flavor_ref").(string),
			},
		}

		_, err := instances.Resize(client, resizeOpts, d.Id()).Extract()
		if err != nil {
			return fmterr.Errorf("error resizing read replica: %w", err)
		}
	}

	return resourceRdsReadReplicaV3Read(ctx, d, meta)
}

func resourceRdsReadReplicaV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	log.Printf("[DEBUG] Deleting Instance %s", d.Id())

	_, err = instances.Delete(client, d.Id()).Extract()
	if err != nil {
		return fmterr.Errorf("error deleting read replica instance: %w", err)
	}

	d.SetId("")
	return nil
}
