package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/instances/lifecycle"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/instances/management"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/instances/specification"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceDmsDedicatedInstanceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaInstanceCreate,
		ReadContext:   resourceDmsKafkaInstanceRead,
		UpdateContext: resourceDmsKafkaInstanceUpdate,
		DeleteContext: resourceDmsKafkaInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

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
			"storage_spec_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ipv6_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arch_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"broker_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"new_tenant_ips": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"storage_space": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"access_user": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				RequiredWith: []string{
					"password",
				},
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
				RequiredWith: []string{
					"access_user",
				},
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
			"security_protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enabled_mechanisms": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"produce_reject", "time_base",
				}, false),
			},

			"ssl_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tags": common.TagsSchema(),
			"engine": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"partition_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_storage_space": {
				Type:     schema.TypeInt,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_spec_code": {
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
			"connector_node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"storage_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cert_replaced": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pod_connect_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_bandwidth": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_two_way_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dumping": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cross_vpc_accesses": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 3,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertised_ip": {
							Type:     schema.TypeString,
							Optional: true,
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

func updateCrossVpcAccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	newVal := d.Get("cross_vpc_accesses")
	var crossVpcAccessArr []map[string]interface{}

	instance, err := lifecycle.Get(client, d.Id())
	if err != nil {
		return fmt.Errorf("error getting DMS Kafka instance: %v", err)
	}

	crossVpcAccessArr, err = flattenCrossVpcInfo(instance.CrossVpcInfo)
	if err != nil {
		return fmt.Errorf("error retrieving details of the cross-VPC: %v", err)
	}

	newAccessArr := newVal.([]interface{})
	contentMap := make(map[string]string)
	for i, oldAccess := range crossVpcAccessArr {
		listenerIp := oldAccess["listener_ip"].(string)
		// If we configure the advertised ip as ["192.168.0.19", "192.168.0.8"], the length of new accesses is 2,
		// and the length of old accesses is always 3.
		if len(newAccessArr) > i {
			// Make sure the index is valid.
			newAccess := newAccessArr[i].(map[string]interface{})
			// Since the "advertised_ip" is already a definition in the schema, the key name must exist.
			if advIp, ok := newAccess["advertised_ip"].(string); ok && advIp != "" {
				contentMap[listenerIp] = advIp
				continue
			}
		}
		contentMap[listenerIp] = listenerIp
	}

	log.Printf("[DEBUG} Update Kafka cross-vpc contentMap: %#v", contentMap)

	retryFunc := func() (interface{}, bool, error) {
		updateRst, err := management.UpdateCrossVpc(client, d.Id(), management.CrossVpcUpdateOpts{
			Contents: contentMap,
		})
		retry, err := handleMultiOperationsError(err)
		return updateRst, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error updating advertised IP: %v", err)
	}
	updateRst := r.(*management.CrossVpc)

	if !updateRst.Success {
		failedIps := make([]string, 0, len(updateRst.Connections))
		for _, conn := range updateRst.Connections {
			if !conn.Success {
				failedIps = append(failedIps, conn.ListenersIp)
			}
		}
		return fmt.Errorf("failed to update the advertised IPs corresponding to some listener IPs (%v)", failedIps)
	}
	return nil
}

func resourceDmsKafkaInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	createOpts := lifecycle.CreateOpts{
		Name:                  d.Get("name").(string),
		Description:           d.Get("description").(string),
		Engine:                "kafka",
		EngineVersion:         d.Get("engine_version").(string),
		BrokerNum:             d.Get("broker_num").(int),
		StorageSpace:          d.Get("storage_space").(int),
		AccessUser:            d.Get("access_user").(string),
		VpcID:                 d.Get("vpc_id").(string),
		SecurityGroupID:       d.Get("security_group_id").(string),
		SubnetID:              d.Get("network_id").(string),
		ProductID:             d.Get("flavor_id").(string),
		ArchType:              d.Get("arch_type").(string),
		MaintainBegin:         d.Get("maintain_begin").(string),
		MaintainEnd:           d.Get("maintain_end").(string),
		RetentionPolicy:       d.Get("retention_policy").(string),
		StorageSpecCode:       d.Get("storage_spec_code").(string),
		SslEnable:             pointerto.Bool(d.Get("ssl_enable").(bool)),
		KafkaSecurityProtocol: d.Get("security_protocol").(string),
		SaslEnabledMechanisms: common.ExpandToStringList(d.Get("enabled_mechanisms").(*schema.Set).List()),
		IPv6Enable:            d.Get("ipv6_enable").(bool),
	}

	if zoneIDs, ok := d.GetOk("available_zones"); ok {
		createOpts.AvailableZones = common.ExpandToStringList(zoneIDs.([]interface{}))
	}

	// set tags
	if tagRaw := d.Get("tags").(map[string]interface{}); len(tagRaw) > 0 {
		createOpts.Tags = common.ExpandResourceTags(tagRaw)
	}
	log.Printf("[DEBUG] Create DMS Kafka instance options: %#v", createOpts)
	createOpts.Password = d.Get("password").(string)

	kafkaInstance, err := lifecycle.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating Kafka instance: %s", err)
	}
	instanceID := kafkaInstance.InstanceID

	var delayTime time.Duration = 300

	log.Printf("[INFO] Creating Kafka instance, ID: %s", instanceID)
	d.SetId(instanceID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      KafkaInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for Kafka instance (%s) to be ready: %s", instanceID, err)
	}

	if _, ok := d.GetOk("cross_vpc_accesses"); ok {
		if err = updateCrossVpcAccess(ctx, client, d); err != nil {
			return diag.Errorf("failed to update default advertised IP: %s", err)
		}
	}

	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func flattenCrossVpcInfo(str string) (result []map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[ERROR] Recover panic when flattening Cross-VPC structure: %#v \nCrossVpcInfo: %s", r, str)
			err = fmt.Errorf("faield to flattening Cross-VPC structure: %#v", r)
		}
	}()

	return unmarshalFlattenCrossVpcInfo(str)
}

func unmarshalFlattenCrossVpcInfo(crossVpcInfoStr string) ([]map[string]interface{}, error) {
	if crossVpcInfoStr == "" {
		return nil, nil
	}

	crossVpcInfos := make(map[string]interface{})
	err := json.Unmarshal([]byte(crossVpcInfoStr), &crossVpcInfos)
	if err != nil {
		return nil, fmt.Errorf("failed to Unmarshal CrossVpcInfo, crossVpcInfo: %s, error: %s", crossVpcInfoStr, err)
	}

	ipArr := make([]string, 0, len(crossVpcInfos))
	for ip := range crossVpcInfos {
		ipArr = append(ipArr, ip)
	}
	sort.Strings(ipArr) // Sort by listeners IP.

	result := make([]map[string]interface{}, len(crossVpcInfos))
	for i, ip := range ipArr {
		crossVpcInfo := crossVpcInfos[ip].(map[string]interface{})
		result[i] = map[string]interface{}{
			"listener_ip":   ip,
			"advertised_ip": crossVpcInfo["advertised_ip"],
			"port":          crossVpcInfo["port"],
			"port_id":       crossVpcInfo["port_id"],
		}
	}
	return result, nil
}

func resourceDmsKafkaInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)

	client, err := config.DmsV2Client(region)
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	v, err := lifecycle.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS Kafka instance")
	}
	log.Printf("[DEBUG] Get Kafka instance: %+v", v)

	crossVpcAccess, err := flattenCrossVpcInfo(v.CrossVpcInfo)
	if err != nil {
		return diag.Errorf("error parsing the cross-VPC information: %v", err)
	}

	partitionNum, _ := strconv.ParseInt(v.PartitionNum, 10, 64)
	mErr := multierror.Append(nil, err)

	mErr = multierror.Append(mErr,
		d.Set("region", config.GetRegion(d)),
		d.Set("flavor_id", v.ProductID),
		d.Set("name", v.Name),
		d.Set("description", v.Description),
		d.Set("engine", v.Engine),
		d.Set("engine_version", v.EngineVersion),
		d.Set("bandwidth", v.Specification),
		d.Set("storage_space", v.TotalStorageSpace),
		d.Set("partition_num", partitionNum),
		d.Set("vpc_id", v.VPCID),
		d.Set("security_group_id", v.SecurityGroupID),
		d.Set("network_id", v.SubnetID),
		d.Set("ipv6_enable", v.EnablePublicIP),
		d.Set("available_zones", v.AvailableZones),
		d.Set("broker_num", v.BrokerNum),
		d.Set("maintain_begin", v.MaintainBegin),
		d.Set("maintain_end", v.MaintainEnd),
		d.Set("retention_policy", v.RetentionPolicy),
		d.Set("dumping", v.ConnectorEnable),
		d.Set("storage_spec_code", v.StorageSpecCode),
		d.Set("used_storage_space", v.UsedStorageSpace),
		d.Set("connect_address", v.ConnectAddress),
		d.Set("port", v.Port),
		d.Set("status", v.Status),
		d.Set("resource_spec_code", v.ResourceSpecCode),
		d.Set("user_id", v.UserID),
		d.Set("user_name", v.UserName),
		d.Set("type", v.Type),
		d.Set("access_user", v.AccessUser),
		d.Set("cross_vpc_accesses", crossVpcAccess),
		d.Set("public_ip_address", v.PublicConnectAddress),
		d.Set("connector_node_num", v.ConnectorNodeNum),
		d.Set("storage_resource_id", v.StorageResourceID),
		d.Set("storage_type", v.StorageType),
		d.Set("created_at", v.CreatedAt),
		d.Set("cert_replaced", v.CertReplaced),
		d.Set("node_num", v.NodeNum),
		d.Set("pod_connect_address", v.PodConnectAddress),
		d.Set("public_bandwidth", v.PublicBandWidth),
		d.Set("ssl_two_way_enable", v.SSLTwoWayEnable),
	)

	// set tags
	if resourceTags, err := tags.Get(client, "kafka", d.Id()).Extract(); err == nil {
		tagMap := common.TagsToMap(resourceTags)
		if err = d.Set("tags", tagMap); err != nil {
			mErr = multierror.Append(mErr,
				fmt.Errorf("error saving tags to state for DMS kafka instance (%s): %s", d.Id(), err))
		}
	} else {
		log.Printf("[WARN] error fetching tags of DMS kafka instance (%s): %s", d.Id(), err)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("failed to set attributes for DMS kafka instance: %s", mErr)
	}

	return nil
}

func resourceDmsKafkaInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	var mErr *multierror.Error
	if d.HasChanges("name", "description", "maintain_begin", "maintain_end",
		"security_group_id", "retention_policy") {
		description := d.Get("description").(string)
		updateOpts := lifecycle.UpdateOpts{
			Description:     &description,
			MaintainBegin:   d.Get("maintain_begin").(string),
			MaintainEnd:     d.Get("maintain_end").(string),
			SecurityGroupID: d.Get("security_group_id").(string),
			RetentionPolicy: d.Get("retention_policy").(string),
		}

		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}

		retryFunc := func() (interface{}, bool, error) {
			_, err = lifecycle.Update(client, d.Id(), updateOpts)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error updating Kafka Instance: %s", err))
		}
	}

	if d.HasChanges("storage_space", "broker_num") {
		err = resizeKafkaInstance(ctx, d, meta)
		if err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("tags") {
		// update tags
		if err = common.UpdateResourceTags(client, d, "kafka", d.Id()); err != nil {
			mErr = multierror.Append(mErr, fmt.Errorf("error updating tags of Kafka instance: %s, err: %s",
				d.Id(), err))
		}
	}

	if d.HasChange("cross_vpc_accesses") {
		if err = updateCrossVpcAccess(ctx, client, d); err != nil {
			mErr = multierror.Append(mErr, err)
		}
	}

	if d.HasChange("password") {
		resetPasswordOpts := management.PasswordOpts{
			NewPassword: d.Get("password").(string),
		}
		retryFunc := func() (interface{}, bool, error) {
			err = management.ResetPassword(client, d.Id(), resetPasswordOpts)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 1 * time.Second,
			PollInterval: 10 * time.Second,
		})
		if err != nil {
			e := fmt.Errorf("error resetting password: %s", err)
			mErr = multierror.Append(mErr, e)
		}
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error while updating DMS Kafka instances, %s", mErr)
	}
	return resourceDmsKafkaInstanceRead(ctx, d, meta)
}

func resizeKafkaInstance(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error initializing DMS(v2) client: %s", err)
	}

	if d.HasChanges("broker_num") {
		operType := "horizontal"
		brokerNum := d.Get("broker_num").(int)

		resizeOpts := specification.IncreaseSpecOpts{
			OperType:        operType,
			NewBrokerNumber: &brokerNum,
			Engine:          "kafka",
		}
		oldNum, newNum := d.GetChange("broker_num")
		if v, ok := d.GetOk("new_tenant_ips"); ok {
			// precheck
			if len(v.([]interface{})) > newNum.(int)-oldNum.(int) {
				return fmt.Errorf("error resizing instance: the nums of new tenant IP must be less than the adding broker nums")
			}
			resizeOpts.TenantIps = common.ExpandToStringList(v.([]interface{}))
		}

		log.Printf("[DEBUG] Resize Kafka instance broker number options: %s", MarshalValue(resizeOpts))

		if err := doKafkaInstanceResize(ctx, d, client, resizeOpts); err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"BOUND"},
			Refresh:      kafkaInstanceBrokerNumberRefreshFunc(client, d.Id(), brokerNum),
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        10 * time.Second,
			PollInterval: 10 * time.Second,
		}
		if _, err = stateConf.WaitForStateContext(ctx); err != nil {
			return err
		}
	}

	if d.HasChanges("storage_space") {
		if err = resizeKafkaInstanceStorage(ctx, d, client); err != nil {
			return err
		}
	}

	return nil
}

func resizeKafkaInstanceStorage(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	newStorageSpace := d.Get("storage_space").(int)
	operType := "storage"
	resizeOpts := specification.IncreaseSpecOpts{
		OperType:        operType,
		NewStorageSpace: &newStorageSpace,
		Engine:          "kafka",
	}
	log.Printf("[DEBUG] Resize Kafka instance storage space options: %s", MarshalValue(resizeOpts))

	return doKafkaInstanceResize(ctx, d, client, resizeOpts)
}

func doKafkaInstanceResize(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, opts specification.IncreaseSpecOpts) error {
	retryFunc := func() (interface{}, bool, error) {
		_, err := specification.IncreaseSpec(client, d.Id(), opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("resize Kafka instance failed: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING", "EXTENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      kafkaResizeStateRefresh(client, d, opts.OperType),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        180 * time.Second,
		PollInterval: 15 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for instance (%s) to resize: %v", d.Id(), err)
	}
	return nil
}

func kafkaResizeStateRefresh(client *golangsdk.ServiceClient, d *schema.ResourceData, operType string) resource.StateRefreshFunc {
	storageSpace := d.Get("storage_space").(int)
	brokerNum := d.Get("broker_num").(int)

	return func() (interface{}, string, error) {
		v, err := lifecycle.Get(client, d.Id())
		if err != nil {
			return nil, "failed", err
		}

		if (operType == "storage" && v.TotalStorageSpace != storageSpace) || // expansion
			(operType == "horizontal" && v.BrokerNum != brokerNum) { // expand broker number
			return v, "PENDING", nil
		}

		return v, v.Status, nil
	}
}

func resourceDmsKafkaInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error initializing DMS Kafka(v2) client: %s", err)
	}

	retryFunc := func() (interface{}, bool, error) {
		err = lifecycle.Delete(client, d.Id())
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     KafkaInstanceStateRefreshFunc(client, d.Id()),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "failed to delete Kafka instance")
	}

	// Wait for the instance to delete before moving on.
	log.Printf("[DEBUG] Waiting for Kafka instance (%s) to be deleted", d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING", "ERROR"}, // Status may change to ERROR on deletion.
		Target:       []string{"DELETED"},
		Refresh:      KafkaInstanceStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        120 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DMS Kafka instance (%s) to be deleted: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] DMS Kafka instance %s has been deleted", d.Id())
	d.SetId("")
	return nil
}

func kafkaInstanceBrokerNumberRefreshFunc(client *golangsdk.ServiceClient, instanceID string, brokerNum int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := lifecycle.Get(client, instanceID)
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		if brokerNum == resp.BrokerNum && resp.CrossVpcInfo != "" {
			crossVpcInfoMap, err := flattenCrossVpcInfo(resp.CrossVpcInfo)
			if err != nil {
				return resp, "ParseError", err
			}

			if len(crossVpcInfoMap) == brokerNum {
				return resp, "BOUND", nil
			}
		}
		return resp, "PENDING", nil
	}
}

func KafkaInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := lifecycle.Get(client, instanceID)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "QUERY ERROR", err
		}

		return v, v.Status, nil
	}
}

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		// CBC.99003651: unsubscribe fail, another operation is being performed
		if errorCode.(string) == "DMS.00400026" || errorCode == "CBC.99003651" {
			return true, err
		}
	}
	return false, err
}
