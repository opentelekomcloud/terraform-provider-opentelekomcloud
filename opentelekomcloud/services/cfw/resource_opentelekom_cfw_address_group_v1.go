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

	group "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/addressgroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwAddressGroupV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWAddressGroupV1Create,
		ReadContext:   resourceCFWAddressGroupV1Read,
		UpdateContext: resourceCFWAddressGroupV1Update,
		DeleteContext: resourceCFWAddressGroupV1Delete,

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
			"address_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_set_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceCFWAddressGroupV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		AddressType: InterfaceToIntPtr(d.Get("address_type")),
	}

	addressGroup, err := group.CreateAddressGroup(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud CFW address group from result: %w", err)
	}
	log.Printf("[DEBUG] Create CFW address group %s: %#v", addressGroup.Id, addressGroup)

	d.SetId(addressGroup.Id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceCFWAddressGroupV1Read(clientCtx, d, meta)
}

func resourceCFWAddressGroupV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	addressGroup, err := group.GetAddressGroup(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching CFW Address Group: %w", err)
	}

	log.Printf("[DEBUG] Retrieved address group %s: %#v", d.Id(), addressGroup)

	mErr := multierror.Append(nil,
		d.Set("id", d.Id()),
		d.Set("name", addressGroup.Name),
		d.Set("description", addressGroup.Description),
		d.Set("address_type", addressGroup.AddressType),
		d.Set("address_set_type", addressGroup.AddressSetType),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCFWAddressGroupV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	updatedAddressGroup, err := group.UpdateAddressGroup(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud CFW Address group %s: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] Updated CFW Address group '%s': %#v", updatedAddressGroup.Id, updatedAddressGroup)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWAddressGroupV1Read(clientCtx, d, meta)
}

func resourceCFWAddressGroupV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud CFW Address Group %s", d.Id())

	err = group.DeleteAddressGroup(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud CFW Address Group: %s", err)
	}

	d.SetId("")
	return nil
}
