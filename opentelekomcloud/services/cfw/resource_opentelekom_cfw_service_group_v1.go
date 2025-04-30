package cfw

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	group "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/servicegroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwServiceGroupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWServiceGroupV1Create,
		ReadContext:   resourceCFWServiceGroupV1Read,
		UpdateContext: resourceCFWServiceGroupV1Update,
		DeleteContext: resourceCFWServiceGroupV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_set_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceCFWServiceGroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := group.CreateOpts{
		ObjectID:    d.Get("object_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	serviceGroup, err := group.CreateServiceGroup(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud CFW service group from result: %w", err)
	}
	log.Printf("[DEBUG] Create CFW service group %s: %#v", serviceGroup.Id, serviceGroup)

	d.SetId(serviceGroup.Id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWServiceGroupV1Read(clientCtx, d, meta)
}

func resourceCFWServiceGroupV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	serviceGroup, err := group.GetServiceGroup(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching CFW Service Group: %w", err)
	}

	log.Printf("[DEBUG] Retrieved service group %s: %#v", d.Id(), serviceGroup)

	mErr := multierror.Append(nil,
		d.Set("id", d.Id()),
		d.Set("name", serviceGroup.Name),
		d.Set("description", serviceGroup.Description),
		d.Set("service_set_type", serviceGroup.ServiceSetType),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCFWServiceGroupV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	updateOpts := group.UpdateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	updatedServiceGroup, err := group.UpdateServiceGroup(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud CFW Service group %s: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] Updated CFW Service group '%s': %#v", updatedServiceGroup.Id, updatedServiceGroup)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWServiceGroupV1Read(clientCtx, d, meta)
}

func resourceCFWServiceGroupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud CFW Service Group %s", d.Id())

	err = group.DeleteServiceGroup(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud CFW Service Group: %s", err)
	}

	d.SetId("")
	return nil
}
