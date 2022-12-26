package rds

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"
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
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				MaxItems: 1,
				Set:      schema.HashString,
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
	job, err := instances.CreateReplica(client, *opts)
	if err != nil {
		return fmterr.Errorf("error creating read replica: %w", err)
	}
	d.SetId(job.Instance.Id)

	timeoutSeconds := d.Timeout(schema.TimeoutCreate).Seconds()
	err = instances.WaitForJobCompleted(client, int(timeoutSeconds), job.JobId)
	if err != nil {
		return fmterr.Errorf("error waiting for read replica to complete creation: %w", err)
	}

	if ip := getReplicaPublicIP(d); ip != "" {
		if err := resourceRdsReadReplicaV3Read(ctx, d, meta); err != nil {
			return err
		}
		nw, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return diag.FromErr(err)
		}
		subnetID, err := getSubnetSubnetID(d, config)
		if err != nil {
			return diag.FromErr(err)
		}
		if err := assignEipToInstance(nw, ip, getReplicaPrivateIP(d), subnetID); err != nil {
			log.Printf("[WARN] failed to assign public IP: %s", err)
		}
	}

	return resourceRdsReadReplicaV3Read(ctx, d, meta)
}

func getReplicaPublicIP(d *schema.ResourceData) string {
	ips := d.Get("public_ips").(*schema.Set)
	if ips.Len() == 0 {
		return ""
	}
	return ips.List()[0].(string)
}

func getReplicaPrivateIP(d *schema.ResourceData) string {
	return d.Get("private_ips").(*schema.Set).List()[0].(string)
}

func resourceRdsReadReplicaV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		setReplicaPrivateIPs(d, meta, replica.PrivateIps),
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

func setReplicaPrivateIPs(d *schema.ResourceData, meta interface{}, privateIPs []string) error {
	if len(privateIPs) == 0 {
		return nil
	}

	config := meta.(*cfg.Config)
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating networking client: %w", err)
	}
	listOpts := floatingips.ListOpts{
		FixedIP: privateIPs[0],
	}

	pages, err := floatingips.List(client, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("error listing floating IPs: %w", err)
	}
	floatingIPs, err := floatingips.ExtractFloatingIPs(pages)
	if err != nil {
		return fmt.Errorf("error listing floating IPs: %w", err)
	}
	addresses := make([]string, len(floatingIPs))
	for i, eip := range floatingIPs {
		addresses[i] = eip.FloatingIP
	}

	return d.Set("public_ips", addresses)
}

func resourceRdsReadReplicaV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	if d.HasChange("flavor_ref") {
		resizeOpts := instances.ResizeOpts{
			InstanceId: d.Id(),
			SpecCode:   d.Get("flavor_ref").(string),
		}

		_, err := instances.Resize(client, resizeOpts)
		if err != nil {
			return fmterr.Errorf("error resizing read replica: %w", err)
		}
	}

	if d.HasChange("public_ips") {
		nwClient, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating networking V2 client: %w", err)
		}
		oldPublicIps, newPublicIps := d.GetChange("public_ips")
		oldIPs := oldPublicIps.(*schema.Set)
		newIPs := newPublicIps.(*schema.Set)

		removeIPs := oldIPs.Difference(newIPs)
		addIPs := newIPs.Difference(oldIPs)

		for _, ip := range removeIPs.List() {
			err := unAssignEipFromInstance(nwClient, ip.(string)) // if it become 0, it was 1 before
			if err != nil {
				return diag.FromErr(err)
			}
		}

		privateIP := getReplicaPrivateIP(d)
		subnetID, err := getSubnetSubnetID(d, config)
		for _, ip := range addIPs.List() {
			if err != nil {
				return diag.FromErr(err)
			}
			if err := assignEipToInstance(nwClient, ip.(string), privateIP, subnetID); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceRdsReadReplicaV3Read(ctx, d, meta)
}

func resourceRdsReadReplicaV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.RdsV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	log.Printf("[DEBUG] Deleting Instance %s", d.Id())

	_, err = instances.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting read replica instance: %w", err)
	}

	d.SetId("")
	return nil
}
