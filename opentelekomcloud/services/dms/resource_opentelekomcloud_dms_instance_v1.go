package dms

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

		DeprecationMessage: "Support will be discontinued in favor of DMS v2. " +
			"Please use `opentelekomcloud_dms_instance_v2` resource instead",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(4, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[\w\-.]+$`),
						"Only lowercase letters, digits, periods (.), underscores (_), and hyphens (-) are allowed.",
					),
					validation.StringDoesNotMatch(
						regexp.MustCompile(`_{3,}?|\.{2,}?|-{2,}?`),
						"Periods, underscores, and hyphens cannot be placed next to each other. A maximum of two consecutive underscores are allowed.",
					),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"kafka",
				}, false),
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				ForceNew:  true,
				Optional:  true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
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
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
			"partition_num": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"produce_reject", "time_base",
				}, false),
			},
			"specification": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"user_name": {
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
				Type:     schema.TypeInt,
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
		},
	}
}

func resourceDmsInstancesV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	sslEnable := false
	if d.Get("access_user").(string) != "" || d.Get("password").(string) != "" {
		sslEnable = true
	}
	createOpts := &instances.CreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Engine:          d.Get("engine").(string),
		EngineVersion:   d.Get("engine_version").(string),
		StorageSpace:    d.Get("storage_space").(int),
		Password:        d.Get("password").(string),
		AccessUser:      d.Get("access_user").(string),
		VpcID:           d.Get("vpc_id").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		SubnetID:        d.Get("subnet_id").(string),
		AvailableZones:  common.GetAllAvailableZones(d),
		ProductID:       d.Get("product_id").(string),
		MaintainBegin:   d.Get("maintain_begin").(string),
		MaintainEnd:     d.Get("maintain_end").(string),
		PartitionNum:    d.Get("partition_num").(int),
		Specification:   d.Get("specification").(string),
		StorageSpecCode: d.Get("storage_spec_code").(string),
		RetentionPolicy: d.Get("retention_policy").(string),
		SslEnable:       sslEnable,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := instances.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud DMSv1 instance: %w", err)
	}
	log.Printf("[INFO] instance ID: %s", v.InstanceID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"RUNNING"},
		Refresh:    instancesV1StateRefreshFunc(client, v.InstanceID),
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

	return resourceDmsInstancesV1Read(ctx, d, meta)
}

func resourceDmsInstancesV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}
	v, err := instances.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS instance")
	}

	log.Printf("[DEBUG] DMS instance %s: %+v", d.Id(), v)

	mErr := multierror.Append(
		d.Set("name", v.Name),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
		d.Set("specification", v.Specification),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("storage_space", v.TotalStorageSpace),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("description", v.Description),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("type", v.Type),
		d.Set("vpc_id", v.VpcID),
		d.Set("vpc_name", v.VpcName),
		d.Set("created_at", v.CreatedAt),
		d.Set("product_id", v.ProductID),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("security_group_name", v.SecurityGroupName),
		d.Set("subnet_id", v.SubnetID),
		d.Set("subnet_name", v.SubnetName),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("access_user", v.AccessUser),
		d.Set("order_id", v.OrderID),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("retention_policy", v.RetentionPolicy),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDmsInstancesV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV1Client(config.GetRegion(d))
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

	if err := instances.Update(client, d.Id(), updateOpts).Err; err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud DMSv1 Instance: %s", err)
	}

	return resourceDmsInstancesV1Read(ctx, d, meta)
}

func resourceDmsInstancesV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	_, err = instances.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "instance")
	}

	err = instances.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud DMSv1 instance: %w", err)
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for instance (%s) to delete", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING", "RUNNING"},
		Target:     []string{"DELETED"},
		Refresh:    instancesV1StateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for instance (%s) to delete: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] DMS instance %s deactivated.", d.Id())
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
