package rts

import (
	"context"
	"log"
	"reflect"
	"unsafe"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/stacks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/stacktemplates"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceRTSStackV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRTSStackV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"outputs": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"timeout_mins": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"disable_rollback": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"capabilities": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"notification_topics": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"template_body": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRTSStackV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud rts client: %s", err)
	}
	stackName := d.Get("name").(string)

	stack, err := stacks.Get(orchestrationClient, stackName).Extract()
	if err != nil {
		return fmterr.Errorf("unable to retrieve stack %s: %s", stackName, err)
	}

	log.Printf("[INFO] Retrieved Stack %s", stackName)
	d.SetId(stack.ID)

	mErr := multierror.Append(
		d.Set("disable_rollback", stack.DisableRollback),
		d.Set("parameters", stack.Parameters),
		d.Set("status_reason", stack.StatusReason),
		d.Set("name", stack.Name),
		d.Set("outputs", flattenStackOutputs(stack.Outputs)),
		d.Set("capabilities", stack.Capabilities),
		d.Set("notification_topics", stack.NotificationTopics),
		d.Set("timeout_mins", stack.Timeout),
		d.Set("status", stack.Status),
		d.Set("region", config.GetRegion(d)),
	)

	out, err := stacktemplates.Get(orchestrationClient, stack.Name, stack.ID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	sTemplate := BytesToString(out)
	template, err := normalizeStackTemplate(sTemplate)
	if err != nil {
		return fmterr.Errorf("template body contains an invalid JSON or YAML: %w", err)
	}
	mErr = multierror.Append(mErr, d.Set("template_body", template))

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func BytesToString(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	var s string
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = sliceHeader.Data
	sh.Len = sliceHeader.Len
	return s
}
