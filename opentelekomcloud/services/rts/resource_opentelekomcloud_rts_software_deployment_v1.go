package rts

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/softwaredeployment"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRtsSoftwareDeploymentV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRtsSoftwareDeploymentV1Create,
		ReadContext:   resourceRtsSoftwareDeploymentV1Read,
		UpdateContext: resourceRtsSoftwareDeploymentV1Update,
		DeleteContext: resourceRtsSoftwareDeploymentV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"input_values": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"status_reason": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"output_values": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceInputValuesV1(d *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range d.Get("input_values").(map[string]interface{}) {
		m[key] = val.(string)
	}

	return m
}

func resourceOutputValuesV1(d *schema.ResourceData) map[string]interface{} {
	m := make(map[string]interface{})
	for key, val := range d.Get("output_values").(map[string]interface{}) {
		m[key] = val.(string)
	}

	return m
}

func resourceRtsSoftwareDeploymentV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Orchestration client: %s", err)
	}

	createOpts := softwaredeployment.CreateOpts{
		Action:       d.Get("action").(string),
		ConfigId:     d.Get("config_id").(string),
		ServerId:     d.Get("server_id").(string),
		StatusReason: d.Get("status_reason").(string),
		Status:       d.Get("status").(string),
		TenantId:     d.Get("tenant_id").(string),
		InputValues:  resourceInputValuesV1(d),
	}

	n, err := softwaredeployment.Create(orchestrationClient, createOpts).Extract()

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud RTS Software Deployment: %s", err)
	}

	d.SetId(n.Id)

	log.Printf("[INFO] Software Deployment ID: %s", n.Id)

	return resourceRtsSoftwareDeploymentV1Read(ctx, d, meta)
}

func resourceRtsSoftwareDeploymentV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Orchestration client: %s", err)
	}

	n, err := softwaredeployment.Get(orchestrationClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmterr.Errorf("error retrieving OpenTelekomCloud RTS Software Deployment: %s", err)
	}

	mErr := multierror.Append(
		d.Set("config_id", n.ConfigId),
		d.Set("status", n.Status),
		d.Set("status_reason", n.StatusReason),
		d.Set("server_id", n.ServerId),
		d.Set("output_values", n.OutputValues),
		d.Set("input_values", n.InputValues),
		d.Set("action", n.Action),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRtsSoftwareDeploymentV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud orchestration client: %s", err)
	}

	var updateOpts softwaredeployment.UpdateOpts

	updateOpts.ConfigId = d.Get("config_id").(string)
	updateOpts.OutputValues = resourceOutputValuesV1(d)

	if d.HasChange("status") {
		updateOpts.Status = d.Get("status").(string)
	}
	if d.HasChange("action") {
		updateOpts.Action = d.Get("action").(string)
	}
	if d.HasChange("status_reason") {
		updateOpts.StatusReason = d.Get("status_reason").(string)
	}
	if d.HasChange("input_values") {
		updateOpts.InputValues = resourceInputValuesV1(d)
	}

	_, err = softwaredeployment.Update(orchestrationClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud RTS Software Deployment: %s", err)
	}

	return resourceRtsSoftwareDeploymentV1Read(ctx, d, meta)
}

func resourceRtsSoftwareDeploymentV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Orchestration: %s", err)
	}

	err = softwaredeployment.Delete(orchestrationClient, d.Id()).ExtractErr()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[INFO] Successfully deleted OpenTelekomCloud RTS Software Deployment %s", d.Id())
		}
		if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
			if errCode.Actual == 409 {
				log.Printf("[INFO] Error deleting OpenTelekomCloud RTS Software Deployment %s", d.Id())
			}
		}
		log.Printf("[INFO] Successfully deleted OpenTelekomCloud RTS Software Deployment %s", d.Id())
	}

	d.SetId("")
	return nil
}
