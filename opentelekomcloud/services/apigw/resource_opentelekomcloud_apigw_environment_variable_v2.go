package apigw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	vars "github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/env_vars"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceAPIEnvironment2Variable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentVariableV2Create,
		ReadContext:   resourceEnvironmentVariableV2Read,
		DeleteContext: resourceEnvironmentVariableV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEnvironmentVariableV2ResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnvironmentVariableV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	opts := vars.CreateOpts{
		GatewayID:     d.Get("gateway_id").(string),
		GroupID:       d.Get("group_id").(string),
		EnvID:         d.Get("environment_id").(string),
		VariableName:  d.Get("name").(string),
		VariableValue: d.Get("value").(string),
	}
	resp, err := vars.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating OpenTelekomCloud APIGW environment variable: %s", err)
	}
	d.SetId(resp.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceEnvironmentVariableV2Read(clientCtx, d, meta)
}

func resourceEnvironmentVariableV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	resp, err := vars.Get(client, d.Get("gateway_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated environment variable")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("group_id", resp.GroupID),
		d.Set("environment_id", resp.EnvID),
		d.Set("name", resp.VariableName),
		d.Set("value", resp.VariableValue),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving OpenTelekomCloud APIGW environment variable (%s) fields: %s", d.Id(), mErr)
	}
	return nil
}

func resourceEnvironmentVariableV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}
	err = vars.Delete(client, d.Get("gateway_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("error deleting OpenTelekomCloud APIGW environment variable(%s): %s", d.Id(), err)
	}

	return nil
}

func resourceEnvironmentVariableV2ResourceImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<gateway_id>/<group_id>/<name>', but got '%s'",
			importedId)
	}

	gatewayId := parts[0]
	groupId := parts[1]
	mErr := multierror.Append(
		d.Set("gateway_id", gatewayId),
		d.Set("group_id", groupId),
	)
	if mErr.ErrorOrNil() != nil {
		return []*schema.ResourceData{d}, mErr
	}

	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.APIGWV2Client(config.GetRegion(d))
	})
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	variables, err := queryEnvironmentVariables(client, gatewayId, groupId)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error getting environment variables: %s", err)
	}

	variableName := parts[2]
	for _, variable := range variables {
		if variable.VariableName == variableName {
			d.SetId(variable.ID)
			return []*schema.ResourceData{d}, nil
		}
	}

	return []*schema.ResourceData{d}, fmt.Errorf("environment variable (%s) not found: %s", variableName, err)
}
