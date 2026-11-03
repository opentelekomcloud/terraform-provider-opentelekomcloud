package dms

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dms/v2/rocketmq/messages"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDmsRocketMQDeadLetterResendV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQDeadLetterResendV2Create,
		ReadContext:   resourceDmsRocketMQDeadLetterResendV2Read,
		DeleteContext: resourceDmsRocketMQDeadLetterResendV2Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"msg_id_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"resend_results": rocketMQResendResultsSchema(),
		},
	}
}

func rocketMQResendResultsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"msg_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"error_code": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"error_message": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func flattenRocketMQResendResults(results []messages.ResendResult) []map[string]interface{} {
	res := make([]map[string]interface{}, len(results))
	for i, r := range results {
		res[i] = map[string]interface{}{
			"msg_id":        r.MsgID,
			"error_code":    r.ErrorCode,
			"error_message": r.ErrorMessage,
		}
	}
	return res
}

func resourceDmsRocketMQDeadLetterResendV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	res, err := messages.ResendDeadLetter(client, d.Get("instance_id").(string), messages.ResendDeadLetterOpts{
		Topic:     d.Get("topic").(string),
		MsgIDList: common.ExpandToStringList(d.Get("msg_id_list").([]interface{})),
	})
	if err != nil {
		return diag.Errorf("error resending DMS RocketMQ dead letter messages: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	// The API returns the resend results only in the response.
	return diag.FromErr(d.Set("resend_results", flattenRocketMQResendResults(res.ResendResults)))
}

func resourceDmsRocketMQDeadLetterResendV2Read(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsRocketMQDeadLetterResendV2Delete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The resource is only removed from the state, the dead letter messages remain resent.",
		},
	}
}
