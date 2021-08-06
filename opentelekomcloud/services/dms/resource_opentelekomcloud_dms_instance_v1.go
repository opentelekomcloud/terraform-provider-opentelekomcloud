package dms

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v1/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsInstancesV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsInstancesV1Create,
		ReadContext:   resourceDmsInstancesV1Read,
		UpdateContext: resourceDmsInstancesV1Update,
		DeleteContext: resourceDmsInstancesV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"partition_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
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
			"connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDmsInstancesV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms instance client: %s", err)
	}

	ssl_enable := false
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		ssl_enable = true
	}
	createOpts := &instances.CreateOps{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Engine:          d.Get("engine").(string),
		EngineVersion:   d.Get("engine_version").(string),
		StorageSpace:    d.Get("storage_space").(int),
		Password:        d.Get("password").(string),
		AccessUser:      d.Get("access_user").(string),
		VPCID:           d.Get("vpc_id").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		SubnetID:        d.Get("subnet_id").(string),
		AvailableZones:  common.GetAllAvailableZones(d),
		ProductID:       d.Get("product_id").(string),
		MaintainBegin:   d.Get("maintain_begin").(string),
		MaintainEnd:     d.Get("maintain_end").(string),
		PartitionNum:    d.Get("partition_num").(int),
		Specification:   d.Get("specification").(string),
		StorageSpecCode: d.Get("storage_spec_code").(string),
		SslEnable:       ssl_enable,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := instances.Create(DmsV1Client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud instance: %s", err)
	}
	log.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    instancesV1StateRefreshFunc(DmsV1Client, v.InstanceID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"Error waiting for instance (%s) to become ready: %s",
			v.InstanceID, err)
	}

	// Store the instance ID now
	d.SetId(v.InstanceID)

	return resourceDmsInstancesV1Read(ctx, d, meta)
}

func resourceDmsInstancesV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms instance client: %s", err)
	}
	v, err := instances.Get(DmsV1Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Dms instance %s: %+v", d.Id(), v)

	d.SetId(v.InstanceID)
	mErr := multierror.Append(
		d.Set("name", v.Name),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
		d.Set("specification", v.Specification),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("description", v.Description),
		d.Set("instance_id", v.InstanceID),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("type", v.Type),
		d.Set("vpc_id", v.VPCID),
		d.Set("vpc_name", v.VPCName),
		d.Set("created_at", v.CreatedAt),
		d.Set("product_id", v.ProductID),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("security_group_name", v.SecurityGroupName),
		d.Set("subnet_id", v.SubnetID),
		d.Set("subnet_name", v.SubnetName),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("order_id", v.OrderID),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDmsInstancesV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud dms instance client: %s", err)
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
		maintainBegin := d.Get("maintain_begin").(string)
		updateOpts.MaintainBegin = maintainBegin
	}
	if d.HasChange("maintain_end") {
		maintainEnd := d.Get("maintain_end").(string)
		updateOpts.MaintainEnd = maintainEnd
	}
	if d.HasChange("security_group_id") {
		securityGroupID := d.Get("security_group_id").(string)
		updateOpts.SecurityGroupID = securityGroupID
	}

	err = instances.Update(DmsV1Client, d.Id(), updateOpts).Err
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Dms Instance: %s", err)
	}

	return resourceDmsInstancesV1Read(ctx, d, meta)
}

func resourceDmsInstancesV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV1Client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud dms instance client: %s", err)
	}

	_, err = instances.Get(DmsV1Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(common.CheckDeleted(d, err, "instance"))
	}

	err = instances.Delete(DmsV1Client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud instance: %s", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "RUNNING"},
		Target:     []string{"DELETED"},
		Refresh:    instancesV1StateRefreshFunc(DmsV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf(
			"Error waiting for instance (%s) to delete: %s",
			d.Id(), err)
	}

	log.Printf("[DEBUG] Dms instance %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func instancesV1StateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
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
