package dataarts

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dataarts/v1.1/cluster"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDataArtsClusterV11() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataArtsClusterV11Create,
		ReadContext:   resourceDataArtsClusterV11Read,
		DeleteContext: resourceDataArtsClusterV11Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"auto_remind": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"phone_num": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"email": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"schedule_boot_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_schedule_boot_off": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"az": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"nics": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 2,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"net_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
						"flavor_ref": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"datastore_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"datastore_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"trial": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"schedule_off_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_auto_off": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"keep_backup_logs": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"flavor_id": {
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
						"role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_ip": {
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

func resourceDataArtsClusterV11Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV11, func() (*golangsdk.ServiceClient, error) {
		return config.DataArtsV11Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV11Client, err)
	}

	createOpts := cluster.CreateOpts{
		XLang:      "en",
		AutoRemind: d.Get("auto_remind").(bool),
		PhoneNum:   d.Get("phone_num").(string),
		Email:      d.Get("email").(string),
		Cluster: cluster.Cluster{
			ScheduleBootTime:  d.Get("schedule_boot_time").(string),
			IsScheduleBootOff: pointerto.Bool(d.Get("is_schedule_boot_off").(bool)),
			Instances:         getInstances(d),
			DataStore: &cluster.Datastore{
				Type:    d.Get("datastore_type").(string),
				Version: d.Get("datastore_version").(string),
			},
			ExtendedProperties: &cluster.ExtendedProp{
				WorkSpaceId: d.Get("workspace_id").(string),
				ResourceId:  d.Get("resource_id").(string),
				Trial:       d.Get("trial").(string),
			},
			ScheduleOffTime: d.Get("schedule_off_time").(string),
			VpcId:           d.Get("vpc_id").(string),
			Name:            d.Get("name").(string),
			IsAutoOff:       d.Get("is_auto_off").(bool),
		},
	}

	createResp, err := cluster.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating dataarts cluster: %w", err)
	}
	d.SetId(createResp.Id)

	err = WaitForClusterState(client, 1200, createResp.Id, "200")
	if err != nil {
		return fmterr.Errorf("error waiting for dataarts cluster to get ready: %w", err)
	}

	log.Printf("Created DataArts Cluster %s: %#v", d.Id(), createResp)

	return resourceDataArtsClusterV11Read(ctx, d, meta)
}

func resourceDataArtsClusterV11Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV11, func() (*golangsdk.ServiceClient, error) {
		return config.DataArtsV11Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV11Client, err)
	}

	getResp, err := cluster.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching dataarts cluster : %w", err)
	}

	mErr := multierror.Append(
		d.Set("is_schedule_boot_off", getResp.IsScheduleBootOff),
		d.Set("datastore_type", getResp.Datastore.Type),
		d.Set("datastore_version", getResp.Datastore.Version),
		d.Set("workspace_id", getResp.CustomerConfig.WorkSpaceId),
		d.Set("resource_id", getResp.CustomerConfig.ResourceId),
		d.Set("trial", getResp.CustomerConfig.Trial),
		d.Set("vpc_id", getResp.VpcId),
		d.Set("name", getResp.Name),
		d.Set("is_auto_off", getResp.IsAutoOff),
		d.Set("security_group_id", getResp.SecurityGroupId),
		d.Set("eip_id", getResp.EipId),
		d.Set("created_at", getResp.Created),
		d.Set("updated_at", getResp.Updated),
		d.Set("status", getResp.Status),
		d.Set("instances_detail", setInstances(getResp.Instances)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDataArtsClusterV11Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV11, func() (*golangsdk.ServiceClient, error) {
		return config.DataArtsV11Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV11Client, err)
	}

	_, err = cluster.Delete(client, d.Id(), cluster.DeleteOpts{
		KeepBackup: d.Get("keep_backup_logs").(int),
	})
	if err != nil {
		return fmterr.Errorf("error deleting dataarts cluster: %w", err)
	}

	d.SetId("")
	return nil
}

func getInstances(d *schema.ResourceData) []cluster.Instance {
	instancesInput := d.Get("instances").([]interface{})
	result := make([]cluster.Instance, 0, len(instancesInput))

	for _, instanceRaw := range instancesInput {
		instanceInput := instanceRaw.(map[string]interface{})
		instance := cluster.Instance{
			AZ:        instanceInput["az"].(string),
			Nics:      getNics(instanceInput["nics"].([]interface{})),
			FlavorRef: instanceInput["flavor_ref"].(string),
			Type:      instanceInput["type"].(string),
		}
		result = append(result, instance)
	}
	return result
}

func getNics(nicsInput []interface{}) []cluster.Nic {
	result := make([]cluster.Nic, 0, len(nicsInput))

	for _, nicRaw := range nicsInput {
		nicInput := nicRaw.(map[string]interface{})
		nic := cluster.Nic{
			SecurityGroupId: nicInput["security_group_id"].(string),
			NetId:           nicInput["net_id"].(string),
		}
		result = append(result, nic)
	}
	return result
}

func setInstances(instancesInResp []cluster.DetailedInstances) []map[string]interface{} {
	var result []map[string]interface{}
	for _, instanceInResp := range instancesInResp {
		instance := map[string]interface{}{
			"id":          instanceInResp.Id,
			"flavor_id":   instanceInResp.Flavor.Id,
			"volume_type": instanceInResp.Volume.Type,
			"volume_size": instanceInResp.Volume.Size,
			"role":        instanceInResp.Role,
			"public_ip":   instanceInResp.PublicIp,
			"private_ip":  instanceInResp.PrivateIp,
			"internal_ip": instanceInResp.InternalIp,
			"status":      instanceInResp.Status,
		}
		result = append(result, instance)

	}
	return result
}
