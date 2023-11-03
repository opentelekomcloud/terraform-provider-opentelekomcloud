package dcaas

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	dceg "github.com/opentelekomcloud/gophertelekomcloud/openstack/dcaas/v2/dc-endpoint-group"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDCEndpointGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDCEndpointGroupV2Create,
		DeleteContext: resourceDCEndpointGroupV2Delete,
		ReadContext:   resourceDCEndpointGroupV2Read,
		UpdateContext: resourceDCEndpointGroupV2Update,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceDCEndpointGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	createOpts := dceg.CreateOpts{
		Name:        d.Get("name").(string),
		TenantId:    d.Get("tenant_id").(string),
		Description: d.Get("description").(string),
		Endpoints:   GetEndpoints(d),
		Type:        d.Get("type").(string),
	}
	log.Printf("[DEBUG] DC endpoint group V2 createOpts: %+v", createOpts)

	created, err := dceg.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating DC endpoint group: %s", err)
	}
	d.SetId(created.ID)
	log.Printf("[DEBUG] DC endpoint group V2 created: %+v", created)
	return resourceDCEndpointGroupV2Read(ctx, d, meta)
}

func resourceDCEndpointGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	created, err := dceg.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error reading DC endpoint group: %s", err)
	}
	log.Printf("[DEBUG] DC endpoint group V2 read: %+v", created)

	mErr := multierror.Append(
		d.Set("name", created.Name),
		d.Set("tenant_id", created.TenantId),
		d.Set("description", created.Description),
		d.Set("endpoints", created.Endpoints),
		d.Set("type", created.Type),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDCEndpointGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	err = dceg.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting DC endpoint group: %s", err)
	}
	log.Printf("[DEBUG] DC endpoint group V2 deleted: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceDCEndpointGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DCaaSV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}
	var updateOpts dceg.UpdateOpts

	if d.HasChange("name") {
		newName := d.Get("name")
		updateOpts.Name = newName.(string)
	}
	if d.HasChange("description") {
		newDescription := d.Get("description")
		updateOpts.Description = newDescription.(string)
	}
	err = dceg.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating direct connect: %s", err)
	}
	return resourceDCEndpointGroupV2Read(ctx, d, meta)
}

func GetEndpoints(d *schema.ResourceData) []string {
	endpoints := make([]string, 0)
	for _, val := range d.Get("endpoints").([]interface{}) {
		endpoints = append(endpoints, val.(string))
	}
	return endpoints
}
