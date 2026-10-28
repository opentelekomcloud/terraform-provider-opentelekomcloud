package dms

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/diagnosis"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsRocketMQInstanceDiagnosisV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQInstanceDiagnosisV2Create,
		ReadContext:   resourceDmsRocketMQInstanceDiagnosisV2Read,
		DeleteContext: resourceDmsRocketMQInstanceDiagnosisV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"node_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"report_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"consumer_nums": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"abnormal_item_sum": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"faulted_node_sum": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"online": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"message_accumulation": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"subscription_consistency": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"duplicate_client_id": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"different_consumer_type": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumers_in_tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"consumers": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"tag_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"diagnosis_node_report_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_faulted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"abnormal_item_sum": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"message_accumulation": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dead_lock": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"deadlock_thread": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_pop": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"consume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDmsRocketMQInstanceDiagnosisV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	createOpts := diagnosis.CreateOpts{
		GroupName:  d.Get("group_name").(string),
		NodeIDList: common.ExpandToStringList(d.Get("node_id_list").([]interface{})),
	}
	created, err := diagnosis.Create(client, instanceID, createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS RocketMQ instance diagnosis task: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceID, created.ReportID))

	stateConf := &resource.StateChangeConf{
		Pending: []string{"diagnosing"},
		Target:  []string{"finished", "normal", "abnormal"},
		Refresh: func() (interface{}, string, error) {
			v, err := diagnosis.Get(client, created.ReportID)
			if err != nil {
				return nil, "", err
			}
			return v, v.Status, nil
		},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 5 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for DMS RocketMQ instance diagnosis task (%s) to finish: %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQInstanceDiagnosisV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQInstanceDiagnosisV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, reportID, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	v, err := diagnosis.Get(client, reportID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS RocketMQ instance diagnosis report")
	}
	if v.Status == "deleted" {
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(
		d.Set("instance_id", instanceID),
		d.Set("group_name", v.GroupName),
		d.Set("report_id", v.ReportID),
		d.Set("consumer_nums", v.ConsumerNums),
		d.Set("status", v.Status),
		d.Set("created_at", time.UnixMilli(v.CreatedAt).UTC().Format(time.RFC3339)),
		d.Set("abnormal_item_sum", v.AbnormalItemSum),
		d.Set("faulted_node_sum", v.FaultedNodeSum),
		d.Set("online", v.Online),
		d.Set("message_accumulation", v.MessageAccumulation),
		d.Set("subscription_consistency", v.SubscriptionConsistency),
		d.Set("duplicate_client_id", v.DuplicateClientID),
		d.Set("different_consumer_type", v.DifferentConsumerType),
		d.Set("subscriptions", flattenRocketMQDiagnosisSubscriptions(v.Subscriptions)),
		d.Set("diagnosis_node_report_list", flattenRocketMQDiagnosisNodeReports(v.NodeReports)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRocketMQDiagnosisSubscriptions(subscriptions []diagnosis.Subscription) []map[string]interface{} {
	res := make([]map[string]interface{}, len(subscriptions))
	for i, s := range subscriptions {
		tags := make([]map[string]interface{}, len(s.ConsumersInTags))
		for j, t := range s.ConsumersInTags {
			tags[j] = map[string]interface{}{
				"consumers": t.Consumers,
				"tag_name":  t.TagName,
			}
		}
		res[i] = map[string]interface{}{
			"topic_name":        s.TopicName,
			"consumers_in_tags": tags,
		}
	}
	return res
}

func flattenRocketMQDiagnosisNodeReports(reports []diagnosis.NodeReport) []map[string]interface{} {
	res := make([]map[string]interface{}, len(reports))
	for i, r := range reports {
		res[i] = map[string]interface{}{
			"node_id":              r.NodeID,
			"is_faulted":           r.IsFaulted,
			"abnormal_item_sum":    r.AbnormalItemSum,
			"message_accumulation": r.MessageAccumulation,
			"dead_lock":            r.DeadLock,
			"deadlock_thread":      r.DeadlockThread,
			"stack_id":             r.StackID,
			"is_pop":               r.IsPop,
			"consume_type":         r.ConsumeType,
		}
	}
	return res
}

func resourceDmsRocketMQInstanceDiagnosisV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, reportID, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = diagnosis.BatchDelete(client, instanceID, diagnosis.BatchDeleteOpts{ReportIDList: []string{reportID}})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DMS RocketMQ instance diagnosis report")
	}
	return nil
}
