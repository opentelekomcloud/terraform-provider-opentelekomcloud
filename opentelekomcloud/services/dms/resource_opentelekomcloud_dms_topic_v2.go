package dms

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/topics"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsTopicsV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsTopicsV2Create,
		ReadContext:   resourceDmsTopicsV2Read,
		DeleteContext: resourceDmsTopicsV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsTopicV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateName,
			},
			"partition": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 200),
			},
			"replication": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 3),
			},
			"sync_replication": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"retention_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 720),
			},
			"sync_message_flush": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"remain_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_partitions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDmsTopicsV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV2Client, err := common.ClientFromCtx(ctx, errCreationClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	createOpts := topics.CreateOpts{
		Name:             d.Get("name").(string),
		Partition:        d.Get("partition").(int),
		Replication:      d.Get("replication").(int),
		SyncReplication:  d.Get("sync_replication").(bool),
		RetentionTime:    d.Get("retention_time").(int),
		SyncMessageFlush: d.Get("sync_message_flush").(bool),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	v, err := topics.Create(DmsV2Client, d.Get("instance_id").(string), createOpts)
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud queue: %s", err)
	}

	d.SetId(v.Name)

	clientCtx := common.CtxWithClient(ctx, DmsV2Client, errCreationClientV2)
	return resourceDmsTopicsV2Read(clientCtx, d, meta)
}

func resourceDmsTopicsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)

	DmsV2Client, err := common.ClientFromCtx(ctx, errCreationClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	instanceID := d.Get("instance_id").(string)

	fTopic, listResponse, err := findDmsTopicV2(DmsV2Client, instanceID, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS topic")
	}
	if fTopic == nil {
		d.SetId("")
		return nil
	}

	var mErr *multierror.Error

	mErr = multierror.Append(
		mErr,
		d.Set("name", fTopic.Name),
		d.Set("partition", fTopic.Partition),
		d.Set("replication", fTopic.Replication),
		d.Set("sync_replication", fTopic.SyncReplication),
		d.Set("retention_time", fTopic.RetentionTime),
		d.Set("sync_message_flush", fTopic.SyncMessageFlush),
		d.Set("size", listResponse.Size),
		d.Set("remain_partitions", listResponse.RemainPartitions),
		d.Set("max_partitions", listResponse.MaxPartitions),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func findDmsTopicV2(client *golangsdk.ServiceClient, instanceID, topicName string) (*topics.Topic, *topics.ListResponse, error) {
	const limit = 50

	for offset := 0; ; offset += limit {
		response, err := topics.List(client, instanceID, topics.ListOpts{
			Offset: strconv.Itoa(offset),
			Limit:  strconv.Itoa(limit),
		})
		if err != nil {
			return nil, nil, err
		}

		for i := range response.Topics {
			if response.Topics[i].Name == topicName {
				return &response.Topics[i], response, nil
			}
		}

		if len(response.Topics) < limit || offset+limit >= response.Total {
			return nil, response, nil
		}
	}
}

func resourceDmsTopicsV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	DmsV2Client, err := common.ClientFromCtx(ctx, errCreationClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	deleteOpts := []string{d.Id()}
	response, err := topics.Delete(DmsV2Client, d.Get("instance_id").(string), deleteOpts)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud topic: %s", err)
	}
	if err := validateDmsTopicV2DeleteResponse(response, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Dms topic %s deactivated.", d.Id())
	d.SetId("")
	return nil
}

func validateDmsTopicV2DeleteResponse(response *topics.DeleteResponse, topicName string) error {
	if response == nil {
		return fmt.Errorf("DMS topic deletion response was empty for topic %q", topicName)
	}

	for _, result := range response.Topics {
		if result.Name == topicName {
			if !result.Success {
				return fmt.Errorf("DMS topic %q was not deleted", topicName)
			}
			return nil
		}
	}

	return fmt.Errorf("DMS topic deletion response did not contain topic %q", topicName)
}

func resourceDmsTopicV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<topic_name>', but '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
