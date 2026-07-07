package eps

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/eps/v1/projects"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceEnterpriseProjectV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnterpriseProjectCreate,
		ReadContext:   resourceEnterpriseProjectRead,
		UpdateContext: resourceEnterpriseProjectUpdate,
		DeleteContext: resourceEnterpriseProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		// Request and response parameters.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"poc", "prod"}, false),
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"skip_disable_on_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceEnterpriseProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := projects.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        d.Get("type").(string),
	}

	project, err := projects.Create(client, createOpts)

	if err != nil {
		return diag.Errorf("error creating enterprise project: %s", err)
	}

	d.SetId(project.ID)

	if !d.Get("enable").(bool) {
		if err := updateEnterpriseProjectEnable(client, d); err != nil {
			return diag.Errorf("error disabling enterprise project in create: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceEnterpriseProjectRead(clientCtx, d, meta)
}

func resourceEnterpriseProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	project, err := projects.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving enterprise project")
	}

	var enable bool
	if project.Status == 1 {
		enable = true
	}

	mErr := multierror.Append(nil,
		d.Set("name", project.Name),
		d.Set("description", project.Description),
		d.Set("type", project.Type),
		d.Set("status", project.Status),
		d.Set("enable", enable),
		d.Set("created_at", project.CreatedAt),
		d.Set("updated_at", project.UpdatedAt),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting enterprise project fields: %s", err)
	}

	return nil
}

func updateEnterpriseProjectEnable(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var actionOpts projects.ActionOpts
	if d.Get("enable").(bool) {
		actionOpts.Action = "enable"
	} else {
		actionOpts.Action = "disable"
	}

	err := projects.Action(client, d.Id(), actionOpts)
	return err
}

func resourceEnterpriseProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if d.HasChange("enable") {
		if err := updateEnterpriseProjectEnable(client, d); err != nil {
			return diag.Errorf("error enabling/disabling enterprise project in update: %s", err)
		}
	}

	if d.HasChanges("name", "description", "type") {
		updateOpts := projects.UpdateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Type:        d.Get("type").(string),
		}

		_, err = projects.Update(client, d.Id(), updateOpts)

		if err != nil {
			return diag.Errorf("error updating enterprise project: %s", err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceEnterpriseProjectRead(clientCtx, d, meta)
}

func resourceEnterpriseProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.EpsV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	if d.Get("skip_disable_on_destroy").(bool) {
		log.Printf("[DEBUG] skip disable on destroy for %s", d.Id())
		return nil
	}

	actionOpts := projects.ActionOpts{
		Action: "disable",
	}

	err = projects.Action(client, d.Id(), actionOpts)
	if err != nil {
		return diag.Errorf("error disabling enterprise project: %s", err)
	}

	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "The project is only disabled and removed from the state, but it remains in the cloud.",
		},
	}
}
