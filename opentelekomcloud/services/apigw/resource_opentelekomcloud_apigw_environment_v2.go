package apigw

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/apigw/v2/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func ResourceAPIEnvironmentv2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEnvironmentResourceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,64}$`),
					"The maximum length is 64 characters. "+
						"Only letters, digits and underscores (_) are allowed"),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := env.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		GatewayID:   d.Get("instance_id").(string),
	}

	resp, err := env.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating dedicated environment: %s", err)
	}
	d.SetId(resp.ID)

	return resourceEnvironmentRead(ctx, d, meta)
}

func GetEnvironment(client *golangsdk.ServiceClient, instanceId,
	envId string) (*env.EnvResp, error) {
	envs, err := env.List(client, env.ListOpts{GatewayID: instanceId})
	if err != nil {
		return nil, err
	}
	for _, v := range envs {
		if v.ID == envId {
			return &v, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceEnvironmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	resp, err := GetEnvironment(client, instanceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated environment")
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("created_at", resp.CreateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving dedicated environment fields: %s", err)
	}
	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opt := env.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		GatewayID:   d.Get("instance_id").(string),
		EnvID:       d.Id(),
	}
	_, err = env.Update(client, opt)
	if err != nil {
		return diag.Errorf("error updating dedicated environment (%s): %s", opt.EnvID, err)
	}

	return resourceEnvironmentRead(ctx, d, meta)
}

func resourceEnvironmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	err = env.Delete(client, instanceId, d.Id())
	if err != nil {
		return diag.Errorf("error deleting dedicated environment from the instance (%s): %s", instanceId, err)
	}

	return nil
}

func resourceEnvironmentResourceImportState(_ context.Context, d *schema.ResourceData,
	meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<name>")
	}
	config := meta.(*cfg.Config)
	client, err := config.APIGWV2Client(config.GetRegion(d))
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = parts[0]
		name       = parts[1]

		opt = env.ListOpts{
			Name:      name,
			GatewayID: instanceId,
		}
	)
	resp, err := env.List(client, opt)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error retrieving environment: %s", err)
	}

	d.SetId(resp[0].ID)
	return []*schema.ResourceData{d}, d.Set("instance_id", instanceId)
}
