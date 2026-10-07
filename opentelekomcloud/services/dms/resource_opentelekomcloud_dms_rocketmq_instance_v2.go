package dms

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/configs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/instances"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/specification"
	rmqtags "github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/tags"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const rocketMQEngine = "reliability"

func ResourceDmsRocketMQInstanceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQInstanceV2Create,
		ReadContext:   resourceDmsRocketMQInstanceV2Read,
		UpdateContext: resourceDmsRocketMQInstanceV2Update,
		DeleteContext: resourceDmsRocketMQInstanceV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		CustomizeDiff: validateRocketMQSpecChange,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"available_zones": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"broker_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tls_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"arch_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enable_acl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_publicip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"publicip_id": {
				Type:         schema.TypeSet,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				RequiredWith: []string{"enable_publicip"},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": common.TagsSchema(),
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_begin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"maintain_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"publicip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"namesrv_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"broker_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_namesrv_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_broker_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"grpc_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_grpc_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cross_vpc_accesses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertised_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// validateRocketMQSpecChange rejects the plans the specification change API can't apply:
// it has no broker number parameter, an instance can only be scaled via flavor_id,
// and the storage can only be expanded.
func validateRocketMQSpecChange(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.Id() == "" {
		return nil
	}
	if d.HasChange("broker_num") && !d.HasChange("flavor_id") {
		return fmt.Errorf("`broker_num` can't be changed separately, the instance can only be scaled via `flavor_id`")
	}
	if oldValue, newValue := d.GetChange("storage_space"); newValue.(int) < oldValue.(int) {
		return fmt.Errorf("`storage_space` can't be decreased from %d to %d GB", oldValue, newValue)
	}
	return nil
}

func resourceDmsRocketMQInstanceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	createOpts := instances.CreateOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Engine:              rocketMQEngine,
		EngineVersion:       d.Get("engine_version").(string),
		StorageSpace:        d.Get("storage_space").(int),
		VpcID:               d.Get("vpc_id").(string),
		SubnetID:            d.Get("subnet_id").(string),
		SecurityGroupID:     d.Get("security_group_id").(string),
		AvailableZones:      common.ExpandToStringListBySet(d.Get("available_zones").(*schema.Set)),
		ProductID:           d.Get("flavor_id").(string),
		StorageSpecCode:     d.Get("storage_spec_code").(string),
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
		BrokerNum:           d.Get("broker_num").(int),
		ArchType:            d.Get("arch_type").(string),
		TlsMode:             d.Get("tls_mode").(string),
	}
	if v, ok := d.GetOk("ssl_enable"); ok {
		createOpts.SslEnable = pointerto.Bool(v.(bool))
	}
	if v, ok := d.GetOk("enable_acl"); ok {
		createOpts.EnableACL = pointerto.Bool(v.(bool))
	}
	if d.Get("enable_publicip").(bool) {
		createOpts.EnablePublicIP = pointerto.Bool(true)
		createOpts.PublicIpID = strings.Join(common.ExpandToStringListBySet(d.Get("publicip_id").(*schema.Set)), ",")
	}
	log.Printf("[DEBUG] Create DMS RocketMQ instance options: %#v", createOpts)

	created, err := instances.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS RocketMQ instance: %s", err)
	}
	d.SetId(created.InstanceID)
	log.Printf("[INFO] Creating DMS RocketMQ instance, ID: %s", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      rocketMQInstanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        60 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for DMS RocketMQ instance (%s) to be ready: %s", d.Id(), err)
	}

	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		if err = rmqtags.Create(client, d.Id(), common.ExpandResourceTags(tagRaw)); err != nil {
			return diag.Errorf("error setting tags of DMS RocketMQ instance (%s): %s", d.Id(), err)
		}
	}

	if v := d.Get("configs").(*schema.Set).List(); len(v) > 0 {
		if err = updateRocketMQConfigs(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate), v); err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQInstanceV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQInstanceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	v, err := instances.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS RocketMQ instance")
	}
	log.Printf("[DEBUG] Get DMS RocketMQ instance: %+v", v)

	crossVpcAccess, err := flattenCrossVpcInfo(v.CrossVpcInfo)
	if err != nil {
		return diag.Errorf("error parsing the cross-VPC information: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", v.Name),
		d.Set("description", v.Description),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
		d.Set("storage_space", v.TotalStorageSpace),
		d.Set("vpc_id", v.VpcID),
		d.Set("subnet_id", v.SubnetID),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("available_zones", v.AvailableZones),
		d.Set("flavor_id", v.ProductID),
		d.Set("storage_spec_code", v.StorageSpecCode),
		d.Set("broker_num", v.BrokerNum),
		d.Set("ssl_enable", v.SslEnable),
		d.Set("tls_mode", v.TlsMode),
		d.Set("arch_type", v.ArchType),
		d.Set("enable_acl", v.EnableACL),
		d.Set("enable_publicip", v.EnablePublicIP),
		d.Set("publicip_id", splitRocketMQPublicIPs(v.PublicIpID)),
		d.Set("enterprise_project_id", v.EnterpriseProjectID),
		d.Set("status", v.Status),
		d.Set("type", v.Type),
		d.Set("specification", v.Specification),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("publicip_address", v.PublicIpAddress),
		d.Set("ipv6_enable", v.IPv6Enable),
		d.Set("namesrv_address", v.NameSrvAddress),
		d.Set("broker_address", v.BrokerAddress),
		d.Set("public_namesrv_address", v.PublicNameSrvAddress),
		d.Set("public_broker_address", v.PublicBrokerAddress),
		d.Set("grpc_address", v.GrpcAddress),
		d.Set("public_grpc_address", v.PublicGrpcAddress),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("created_at", v.CreatedAt),
		d.Set("cross_vpc_accesses", crossVpcAccess),
	)

	// Only the configurations managed by the resource are tracked:
	// the instance has a lot of them, all with default values.
	if names := rocketMQConfigNames(d.Get("configs").(*schema.Set).List()); len(names) > 0 {
		if rocketMQConfigs, err := getRocketMQConfigs(client, d.Id(), names); err == nil {
			mErr = multierror.Append(mErr, d.Set("configs", rocketMQConfigs))
		} else {
			log.Printf("[WARN] error fetching configs of DMS RocketMQ instance (%s): %s", d.Id(), err)
		}
	}

	// A RocketMQ instance can't have more than 20 tags.
	if resourceTags, err := rmqtags.Get(client, d.Id(), rmqtags.ListOpts{Limit: 20}); err == nil {
		mErr = multierror.Append(mErr, d.Set("tags", common.TagsToMap(resourceTags.Tags)))
	} else {
		log.Printf("[WARN] error fetching tags of DMS RocketMQ instance (%s): %s", d.Id(), err)
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting DMS RocketMQ instance fields: %s", err)
	}

	return nil
}

func resourceDmsRocketMQInstanceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}
	timeout := d.Timeout(schema.TimeoutUpdate)

	if d.HasChanges("name", "description", "security_group_id", "enable_acl", "enterprise_project_id") {
		description := d.Get("description").(string)
		updateOpts := instances.UpdateOpts{
			Description:     &description,
			SecurityGroupID: d.Get("security_group_id").(string),
		}
		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}
		if d.HasChange("enable_acl") {
			updateOpts.EnableACL = pointerto.Bool(d.Get("enable_acl").(bool))
		}
		if d.HasChange("enterprise_project_id") {
			updateOpts.EnterpriseProjectID = d.Get("enterprise_project_id").(string)
		}
		if err = updateRocketMQInstance(ctx, client, d.Id(), timeout, updateOpts); err != nil {
			return diag.Errorf("error updating DMS RocketMQ instance (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		if err = updateRocketMQTags(client, d); err != nil {
			return diag.Errorf("error updating tags of DMS RocketMQ instance (%s): %s", d.Id(), err)
		}
	}

	publicIPChanged := d.HasChanges("enable_publicip", "publicip_id")
	if oldEnabled, _ := d.GetChange("enable_publicip"); publicIPChanged && oldEnabled.(bool) {
		log.Printf("[DEBUG] Unbinding EIPs from DMS RocketMQ instance (%s)", d.Id())
		if err = setRocketMQPublicIP(ctx, client, d.Id(), timeout, false, ""); err != nil {
			return diag.Errorf("error disabling public access of DMS RocketMQ instance (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("flavor_id") {
		newProductID := d.Get("flavor_id").(string)
		resizeOpts := specification.ResizeOpts{
			OperType:     "horizontal",
			NewProductID: newProductID,
		}
		// EIPs are mandatory for a public instance, they are left bound unless they are being changed.
		if !publicIPChanged && d.Get("enable_publicip").(bool) {
			resizeOpts.PublicIPID = strings.Join(common.ExpandToStringListBySet(d.Get("publicip_id").(*schema.Set)), ",")
		}
		err = resizeRocketMQInstance(ctx, client, d.Id(), timeout, resizeOpts, func(v *instances.Instance) bool {
			return v.ProductID == newProductID
		})
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// The new flavor can change the total storage, so it is compared with the actual value.
	if d.HasChanges("storage_space", "flavor_id") {
		v, err := instances.Get(client, d.Id())
		if err != nil {
			return diag.Errorf("error retrieving DMS RocketMQ instance (%s): %s", d.Id(), err)
		}
		newStorageSpace := d.Get("storage_space").(int)
		if v.TotalStorageSpace < newStorageSpace {
			resizeOpts := specification.ResizeOpts{
				OperType:        "storage",
				NewStorageSpace: newStorageSpace,
			}
			err = resizeRocketMQInstance(ctx, client, d.Id(), timeout, resizeOpts, func(v *instances.Instance) bool {
				return v.TotalStorageSpace == newStorageSpace
			})
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if publicIPChanged && d.Get("enable_publicip").(bool) {
		log.Printf("[DEBUG] Binding EIPs to DMS RocketMQ instance (%s)", d.Id())
		publicIPs := strings.Join(common.ExpandToStringListBySet(d.Get("publicip_id").(*schema.Set)), ",")
		if err = setRocketMQPublicIP(ctx, client, d.Id(), timeout, true, publicIPs); err != nil {
			return diag.Errorf("error enabling public access of DMS RocketMQ instance (%s): %s", d.Id(), err)
		}
	}

	if d.HasChange("configs") {
		oldRaw, newRaw := d.GetChange("configs")
		changed := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set)).List()
		if len(changed) > 0 {
			if err = updateRocketMQConfigs(ctx, client, d.Id(), timeout, changed); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQInstanceV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQInstanceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	retryFunc := func() (interface{}, bool, error) {
		err := instances.Delete(client, d.Id())
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rocketMQInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DMS RocketMQ instance")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING", "ERROR"}, // Status may change to ERROR on deletion.
		Target:       []string{"DELETED"},
		Refresh:      rocketMQInstanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        60 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for DMS RocketMQ instance (%s) to be deleted: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] DMS RocketMQ instance %s has been deleted", d.Id())
	d.SetId("")
	return nil
}

func updateRocketMQInstance(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration,
	opts instances.UpdateOpts) error {
	retryFunc := func() (interface{}, bool, error) {
		err := instances.Update(client, id, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rocketMQInstanceStateRefreshFunc(client, id),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	return err
}

// setRocketMQPublicIP binds (enabled is true) or unbinds the EIPs and waits until the instance reports the change.
func setRocketMQPublicIP(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration,
	enabled bool, publicIPs string) error {
	opts := instances.UpdateOpts{
		EnablePublicIP: pointerto.Bool(enabled),
		PublicIpID:     publicIPs,
	}
	if err := updateRocketMQInstance(ctx, client, id, timeout, opts); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"DONE"},
		Refresh: rocketMQInstanceConditionRefreshFunc(client, id, func(v *instances.Instance) bool {
			return v.EnablePublicIP == enabled
		}),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resizeRocketMQInstance(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration,
	opts specification.ResizeOpts, done func(*instances.Instance) bool) error {
	log.Printf("[DEBUG] Resize DMS RocketMQ instance options: %s", MarshalValue(opts))

	retryFunc := func() (interface{}, bool, error) {
		_, err := specification.Resize(client, id, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rocketMQInstanceStateRefreshFunc(client, id),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error resizing DMS RocketMQ instance (%s): %s", id, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DONE"},
		Refresh:      rocketMQInstanceConditionRefreshFunc(client, id, done),
		Timeout:      timeout,
		Delay:        60 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for DMS RocketMQ instance (%s) to be resized: %s", id, err)
	}
	return nil
}

func updateRocketMQTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("tags")
	oldMap, newMap := oldRaw.(map[string]interface{}), newRaw.(map[string]interface{})

	removed := make(map[string]interface{})
	for k, v := range oldMap {
		if newValue, ok := newMap[k]; !ok || newValue != v {
			removed[k] = v
		}
	}
	added := make(map[string]interface{})
	for k, v := range newMap {
		if oldValue, ok := oldMap[k]; !ok || oldValue != v {
			added[k] = v
		}
	}

	if len(removed) > 0 {
		if err := rmqtags.Delete(client, d.Id(), common.ExpandResourceTags(removed)); err != nil {
			return err
		}
	}
	if len(added) > 0 {
		if err := rmqtags.Create(client, d.Id(), common.ExpandResourceTags(added)); err != nil {
			return err
		}
	}
	return nil
}

func updateRocketMQConfigs(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration,
	rawConfigs []interface{}) error {
	opts := configs.UpdateOpts{
		Configs: make([]configs.UpdateConfig, len(rawConfigs)),
	}
	for i, raw := range rawConfigs {
		c := raw.(map[string]interface{})
		opts.Configs[i] = configs.UpdateConfig{
			Name:  c["name"].(string),
			Value: c["value"].(string),
		}
	}

	retryFunc := func() (interface{}, bool, error) {
		err := configs.Update(client, id, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rocketMQInstanceStateRefreshFunc(client, id),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      timeout,
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating configs of DMS RocketMQ instance (%s): %s", id, err)
	}
	return nil
}

func getRocketMQConfigs(client *golangsdk.ServiceClient, id string, names []string) ([]map[string]interface{}, error) {
	resp, err := configs.List(client, id, configs.ListOpts{})
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, c := range resp.Configs {
		if common.StrSliceContains(names, c.Name) {
			result = append(result, map[string]interface{}{
				"name":  c.Name,
				"value": c.Value,
			})
		}
	}
	return result, nil
}

func rocketMQConfigNames(rawConfigs []interface{}) []string {
	names := make([]string, len(rawConfigs))
	for i, raw := range rawConfigs {
		names[i] = raw.(map[string]interface{})["name"].(string)
	}
	return names
}

func splitRocketMQPublicIPs(publicIPs string) []string {
	if publicIPs == "" {
		return nil
	}
	return strings.Split(publicIPs, ",")
}

func rocketMQInstanceStateRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, id)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "QUERY ERROR", err
		}
		return v, v.Status, nil
	}
}

// rocketMQInstanceConditionRefreshFunc reports DONE once the instance is running and done returns true.
func rocketMQInstanceConditionRefreshFunc(client *golangsdk.ServiceClient, id string,
	done func(*instances.Instance) bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := instances.Get(client, id)
		if err != nil {
			return nil, "QUERY ERROR", err
		}
		switch v.Status {
		case "ERROR", "EXTENDEDFAILED":
			return v, v.Status, fmt.Errorf("unexpected DMS RocketMQ instance status: %s", v.Status)
		case "RUNNING":
			if done(v) {
				return v, "DONE", nil
			}
		}
		return v, "PENDING", nil
	}
}
