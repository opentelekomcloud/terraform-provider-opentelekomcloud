package eps

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/eps/v1/projects"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

var actionNonUpdatableParams = []string{
	"enterprise_project_id",
	"action",
}

func ResourceEnterpriseProjectActionV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseProjectActionV1Create,
		UpdateContext: resourceEnterpriseProjectActionV1Update,
		ReadContext:   resourceEnterpriseProjectActionV1Read,
		DeleteContext: resourceEnterpriseProjectActionV1Delete,

		CustomizeDiff: common.FlexibleForceNew(actionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
		},
	}
}

func resourceEnterpriseProjectActionV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}
	epsId := d.Get("enterprise_project_id").(string)
	action := d.Get("action").(string)
	createOpts := projects.ActionOpts{
		Action: action,
	}

	err = projects.Action(client, epsId, createOpts)
	if err != nil {
		return diag.Errorf("unable to %s the enterprise project (%s): %s", action, epsId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceEnterpriseProjectActionV1Read(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEnterpriseProjectActionV1Update(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEnterpriseProjectActionV1Delete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the enterprise project. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
