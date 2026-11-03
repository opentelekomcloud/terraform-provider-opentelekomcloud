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

func ResourceDmsRocketMQConsumptionVerifyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQConsumptionVerifyV2Create,
		ReadContext:   resourceDmsRocketMQConsumptionVerifyV2Read,
		DeleteContext: resourceDmsRocketMQConsumptionVerifyV2Delete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
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

func resourceDmsRocketMQConsumptionVerifyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, dmsClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DmsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClientV2, err)
	}

	res, err := messages.VerifyConsumption(client, d.Get("instance_id").(string), messages.VerifyConsumptionOpts{
		Group:     d.Get("group").(string),
		Topic:     d.Get("topic").(string),
		ClientID:  d.Get("client_id").(string),
		MsgIDList: common.ExpandToStringList(d.Get("msg_id_list").([]interface{})),
	})
	if err != nil {
		return diag.Errorf("error verifying DMS RocketMQ message consumption: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	// The API returns the verification results only in the response.
	return diag.FromErr(d.Set("resend_results", flattenRocketMQResendResults(res.ResendResults)))
}

func resourceDmsRocketMQConsumptionVerifyV2Read(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsRocketMQConsumptionVerifyV2Delete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting action resource is not supported. The resource is only removed from the state, the messages remain re-delivered to the client.",
		},
	}
}
