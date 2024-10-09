package dms

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/pointerto"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/smart_connect"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsSmartConnectTaskV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsV2SmartConnectTaskCreate,
		ReadContext:   resourceDmsV2SmartConnectTaskRead,
		DeleteContext: resourceDmsV2SmartConnectTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsV2SmartConnectTaskImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topics": {
				Type:         schema.TypeSet,
				Optional:     true,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"topics", "topics_regex"},
			},
			"topics_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"start_later": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"destination_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"source_task": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"current_instance_alias": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"peer_instance_alias": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"peer_instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"peer_instance_address": {
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"security_protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"sasl_mechanism": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},

						// task configuration
						"direction": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"sync_consumer_offsets_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"replication_factor": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"task_num": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"rename_topic_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"provenance_header_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"consumer_strategy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"compression_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"topics_mapping": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"destination_task": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Kafka to OBS
						"access_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"secret_key": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"consumer_strategy": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"deliver_time_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"obs_bucket_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"partition_format": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"obs_path": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"destination_file_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"record_delimiter": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"store_keys": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsV2SmartConnectTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceId := d.Get("instance_id").(string)
	topics := d.Get("topics").(*schema.Set).List()

	createOpts := smart_connect.CreateTaskOpts{
		TaskName:    d.Get("task_name").(string),
		TopicsRegex: d.Get("topics_regex").(string),
		StartLater:  pointerto.Bool(d.Get("start_later").(bool)),
		SourceType:  d.Get("source_type").(string),
		SinkType:    d.Get("destination_type").(string),
		Topics:      changeListToStringWithCommasSplit(topics),
		SourceTask:  buildSourceTaskRequestBody(d.Get("source_task").([]interface{})),
		SinkTask:    buildSinkTaskRequestBody(d.Get("destination_task").([]interface{})),
	}

	createTaskResp, err := smart_connect.CreateTask(client, instanceId, createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS kafka smart connect task: %v", err)
	}

	d.SetId(createTaskResp.ID)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING", "WAITING"},
		Refresh:      DMSv2SmartConnectTaskStateRefreshFunc(client, instanceId, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the smart connect task (%s) to be done: %s", d.Id(), err)
	}

	return resourceDmsV2SmartConnectTaskRead(ctx, d, meta)
}

func buildSourceTaskRequestBody(rawParams []interface{}) *smart_connect.SmartConnectTaskSourceConfig {
	if len(rawParams) == 0 {
		return nil
	}
	params := rawParams[0].(map[string]interface{})
	rst := smart_connect.SmartConnectTaskSourceConfig{
		CurrentClusterName: params["current_instance_alias"].(string),
		ClusterName:        params["peer_instance_alias"].(string),
		InstanceId:         params["peer_instance_id"].(string),
		BootstrapServers: changeListToStringWithCommasSplit(
			params["peer_instance_address"].(*schema.Set).List()),
		SecurityProtocol:           params["security_protocol"].(string),
		UserName:                   params["user_name"].(string),
		Password:                   params["password"].(string),
		SaslMechanism:              params["sasl_mechanism"].(string),
		Direction:                  params["direction"].(string),
		SyncConsumerOffsetsEnabled: pointerto.Bool(params["sync_consumer_offsets_enabled"].(bool)),
		ReplicationFactor:          params["replication_factor"].(int),
		TaskNum:                    params["task_num"].(int),
		RenameTopicEnabled:         pointerto.Bool(params["rename_topic_enabled"].(bool)),
		ProvenanceHeaderEnabled:    pointerto.Bool(params["provenance_header_enabled"].(bool)),
		ConsumerStrategy:           params["consumer_strategy"].(string),
		CompressionType:            params["compression_type"].(string),
		TopicsMapping: changeListToStringWithCommasSplit(
			params["topics_mapping"].(*schema.Set).List()),
	}
	return &rst
}

func changeListToStringWithCommasSplit(params []interface{}) string {
	strArray := make([]string, 0, len(params))
	for _, param := range params {
		strArray = append(strArray, param.(string))
	}
	return strings.Join(strArray, ",")
}

func buildSinkTaskRequestBody(rawParams []interface{}) *smart_connect.SmartConnectTaskSinkConfig {
	if len(rawParams) == 0 {
		return nil
	}
	params := rawParams[0].(map[string]interface{})
	rst := smart_connect.SmartConnectTaskSinkConfig{
		ConsumerStrategy:    params["consumer_strategy"].(string),
		AccessKey:           params["access_key"].(string),
		SecretKey:           params["secret_key"].(string),
		ObsBucketName:       params["obs_bucket_name"].(string),
		PartitionFormat:     params["partition_format"].(string),
		DeliverTimeInterval: params["deliver_time_interval"].(int),
		ObsPath:             params["obs_path"].(string),
		RecordDelimiter:     params["record_delimiter"].(string),
		DestinationFileType: params["destination_file_type"].(string),
		StoreKeys:           pointerto.Bool(params["store_keys"].(bool)),
	}

	return &rst
}

func resourceDmsV2SmartConnectTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	getTask, err := smart_connect.GetTask(client, d.Get("instance_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS kafka smart connect task")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("task_name", getTask.TaskName),
		d.Set("topics", flattenStringWithCommaSplitToSlice(getTask.Topics)),
		d.Set("topics_regex", getTask.TopicsRegex),
		d.Set("source_type", getTask.SourceType),
		d.Set("destination_type", getTask.SinkType),
		d.Set("source_task", flattenSourceTaskResponse(d, *getTask.SourceTask)),
		d.Set("destination_task", flattenSinkTaskResponse(d, *getTask.SinkTask)),
		d.Set("created_at", strconv.FormatInt(getTask.CreateTime, 10)),
		d.Set("status", getTask.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSourceTaskResponse(d *schema.ResourceData, rawParams smart_connect.SmartConnectTaskSourceConfig) []map[string]interface{} {
	if reflect.DeepEqual(rawParams, smart_connect.SmartConnectTaskSourceConfig{}) {
		return nil
	}
	rst := make([]map[string]interface{}, 1)
	params := map[string]interface{}{
		"current_instance_alias":        rawParams.CurrentClusterName,
		"peer_instance_alias":           rawParams.ClusterName,
		"peer_instance_id":              rawParams.InstanceId,
		"peer_instance_address":         strings.Split(rawParams.BootstrapServers, ","),
		"security_protocol":             rawParams.SecurityProtocol,
		"user_name":                     rawParams.UserName,
		"password":                      d.Get("source_task.0.password"),
		"sasl_mechanism":                rawParams.SaslMechanism,
		"direction":                     rawParams.Direction,
		"sync_consumer_offsets_enabled": rawParams.SyncConsumerOffsetsEnabled,
		"replication_factor":            rawParams.ReplicationFactor,
		"task_num":                      rawParams.TaskNum,
		"rename_topic_enabled":          rawParams.RenameTopicEnabled,
		"provenance_header_enabled":     rawParams.ProvenanceHeaderEnabled,
		"consumer_strategy":             rawParams.ConsumerStrategy,
		"compression_type":              rawParams.CompressionType,
		"topics_mapping":                flattenStringWithCommaSplitToSlice(rawParams.TopicsMapping),
	}
	rst[0] = params
	return rst
}

func flattenSinkTaskResponse(d *schema.ResourceData, rawParams smart_connect.SmartConnectTaskSinkConfig) []map[string]interface{} {
	if reflect.DeepEqual(rawParams, smart_connect.SmartConnectTaskSinkConfig{}) {
		return nil
	}
	rst := make([]map[string]interface{}, 1)
	params := map[string]interface{}{
		"access_key":            d.Get("destination_task.0.access_key"),
		"secret_key":            d.Get("destination_task.0.secret_key"),
		"consumer_strategy":     rawParams.ConsumerStrategy,
		"destination_file_type": rawParams.DestinationFileType,
		"obs_bucket_name":       rawParams.ObsBucketName,
		"obs_path":              rawParams.ObsPath,
		"partition_format":      rawParams.PartitionFormat,
		"record_delimiter":      rawParams.RecordDelimiter,
		"deliver_time_interval": rawParams.DeliverTimeInterval,
		"store_keys":            rawParams.StoreKeys,
	}

	rst[0] = params
	return rst
}

func flattenStringWithCommaSplitToSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func resourceDmsV2SmartConnectTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	err = smart_connect.DeleteTask(client, d.Get("instance_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("error deleting DMS kafka smart connect task: %v", err)
	}

	return nil
}

// resourceDmsKafkav2SmartConnectTaskImportState is used to import an id with format <instance_id>/<task_id>
func resourceDmsV2SmartConnectTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<task_id>")
	}

	err := d.Set("instance_id", parts[0])
	if err != nil {
		return nil, err
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func DMSv2SmartConnectTaskStateRefreshFunc(client *golangsdk.ServiceClient, instanceID, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getTaskResp, err := smart_connect.GetTask(client, instanceID, taskID)
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		return getTaskResp, getTaskResp.Status, nil
	}
}
