package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/migration"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsRocketMQMigrationTaskV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQMigrationTaskV2Create,
		ReadContext:   resourceDmsRocketMQMigrationTaskV2Read,
		DeleteContext: resourceDmsRocketMQMigrationTaskV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		// Nested field names match the API request field names, see expandRocketMQMigrationMetadata.
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[A-Za-z][\w-]{3,63}$`),
					"must start with a letter, consist of 4 to 64 characters and contain only letters, digits, underscores (_) and hyphens (-)"),
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"rocketmq", "rabbitToRocket"}, false),
			},
			"overwrite": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"topic_configs": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"order": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"perm": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"read_queue_nums": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"write_queue_nums": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"topic_filter_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"SINGLE_TAG", "MULTI_TAG"}, false),
						},
						"topic_sys_flag": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"subscription_groups": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"consume_broadcast_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"consume_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"consume_from_min_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"notify_consumer_ids_changed_enable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"retry_max_times": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"retry_queue_nums": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"which_broker_when_consume_slow": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"vhosts": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"queues": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vhost": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"durable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"exchanges": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vhost": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"topic", "direct", "fanout", "headers"}, false),
						},
						"durable": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"bindings": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vhost": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"destination": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"destination_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"exchange", "queue"}, false),
						},
						"routing_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"start_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// expandRocketMQMigrationMetadata decodes nested blocks into the SDK structures,
// as the schema field names match their JSON names.
func expandRocketMQMigrationMetadata(d *schema.ResourceData, key string, out interface{}) error {
	b, err := json.Marshal(d.Get(key))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

func buildRocketMQMigrationTaskCreateOpts(d *schema.ResourceData) (migration.CreateOpts, error) {
	opts := migration.CreateOpts{
		Overwrite: d.Get("overwrite").(bool),
		Name:      d.Get("name").(string),
		Type:      d.Get("type").(string),
	}

	var topics []migration.TopicConfig
	var groups []migration.SubscriptionGroup
	mErr := multierror.Append(
		expandRocketMQMigrationMetadata(d, "topic_configs", &topics),
		expandRocketMQMigrationMetadata(d, "subscription_groups", &groups),
		expandRocketMQMigrationMetadata(d, "vhosts", &opts.Vhosts),
		expandRocketMQMigrationMetadata(d, "queues", &opts.Queues),
		expandRocketMQMigrationMetadata(d, "exchanges", &opts.Exchanges),
		expandRocketMQMigrationMetadata(d, "bindings", &opts.Bindings),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return opts, err
	}

	if len(topics) > 0 {
		opts.TopicConfigTable = make(map[string]migration.TopicConfig, len(topics))
		for _, t := range topics {
			opts.TopicConfigTable[t.TopicName] = t
		}
	}
	if len(groups) > 0 {
		opts.SubscriptionGroupTable = make(map[string]migration.SubscriptionGroup, len(groups))
		for _, g := range groups {
			opts.SubscriptionGroupTable[g.GroupName] = g
		}
	}
	return opts, nil
}

// getRocketMQMigrationTask looks the task up in the task list, because details of a rabbitToRocket task
// can't be queried without a query type and a missing task is reported with the same error.
func getRocketMQMigrationTask(client *golangsdk.ServiceClient, instanceID, taskID string) (*migration.Task, error) {
	const limit = 50
	for page := 1; ; page++ {
		res, err := migration.List(client, instanceID, migration.ListOpts{Offset: page, Limit: limit})
		if err != nil {
			return nil, err
		}
		for _, t := range res.Tasks {
			if t.ID == taskID {
				return &t, nil
			}
		}
		if len(res.Tasks) == 0 || page*limit >= res.Total {
			return nil, golangsdk.ErrDefault404{}
		}
	}
}

func resourceDmsRocketMQMigrationTaskV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)
	createOpts, err := buildRocketMQMigrationTaskCreateOpts(d)
	if err != nil {
		return diag.Errorf("error building DMS RocketMQ migration task options: %s", err)
	}

	created, err := migration.Create(client, instanceID, createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS RocketMQ migration task: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceID, created.TaskID))

	stateConf := &resource.StateChangeConf{
		Pending: []string{"creating", "starting"},
		Target:  []string{"finished"},
		Refresh: func() (interface{}, string, error) {
			v, err := getRocketMQMigrationTask(client, instanceID, created.TaskID)
			if err != nil {
				return nil, "", err
			}
			return v, v.Status, nil
		},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 5 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for DMS RocketMQ migration task (%s) to finish: %s", d.Id(), err)
	}

	clientCtx := common.CtxWithClient(ctx, client, dmsClientV2)
	return resourceDmsRocketMQMigrationTaskV2Read(clientCtx, d, meta)
}

func resourceDmsRocketMQMigrationTaskV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, taskID, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	v, err := getRocketMQMigrationTask(client, instanceID, taskID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS RocketMQ migration task")
	}

	mErr := multierror.Append(
		d.Set("instance_id", instanceID),
		d.Set("name", v.Name),
		d.Set("type", v.Type),
		d.Set("start_date", v.StartDate),
		d.Set("status", v.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQMigrationTaskV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID, taskID, err := parseRocketMQSubResourceID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Deleting a missing task succeeds as well.
	_, err = migration.BatchDelete(client, instanceID, migration.BatchDeleteOpts{TaskIDs: []string{taskID}})
	if err != nil {
		return diag.Errorf("error deleting DMS RocketMQ migration task (%s): %s", d.Id(), err)
	}
	return nil
}
