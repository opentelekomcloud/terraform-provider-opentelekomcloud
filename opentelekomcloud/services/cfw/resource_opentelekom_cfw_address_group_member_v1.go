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

func ResourceCfwAddressGroupMemberV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWAddressGroupmemberV1Create,
		ReadContext:   resourceCFWAddressGroupMemberV1Read,
		DeleteContext: resourceCFWAddressGroupMemberV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("set_id", "address"),
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"set_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"address_type": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCFWAddressGroupmemberV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	setId := d.Get("set_id").(string)
	ipAddress := d.Get("address").(string)
	createOpts := group.AddGroupMemberOpts{
		SetId: setId,
		AddressItems: []group.AddressItem{
			{
				AddressType: InterfaceToIntPtr(d.Get("address_type")),
				Address:     ipAddress,
				Description: d.Get("description").(string),
			},
		},
	}

	_, err = group.AddGroupMember(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud CFW address group member from result: %w", err)
	}

	groupMember, err := group.GetGroupMember(client, setId, ipAddress)
	if err != nil {
		return fmterr.Errorf("unable to create OpenTelekomCloud CFW address group member: %w", err)
	}

	log.Printf("[DEBUG] Created CFW Address group member %s: %#v", groupMember.ItemID, groupMember)

	d.SetId(groupMember.ItemID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceCFWAddressGroupMemberV1Read(clientCtx, d, meta)
}

func resourceCFWAddressGroupMemberV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	setId := d.Get("set_id").(string)
	ipAddress := d.Get("address").(string)

	groupMember, err := group.GetGroupMember(client, setId, ipAddress)
	if err != nil {
		return fmterr.Errorf("error fetching CFW Address Group Member: %w", err)
	}
	d.SetId(groupMember.ItemID)

	log.Printf("[DEBUG] Retrieved instance %s: %#v", d.Id(), groupMember)

	mErr := multierror.Append(nil,
		d.Set("id", d.Id()),
		d.Set("name", groupMember.Name),
		d.Set("address", groupMember.Address),
		d.Set("address_type", groupMember.AddressType),
		d.Set("description", groupMember.Description),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCFWAddressGroupMemberV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud CFW Address Group Member %s", d.Id())

	err = group.DeleteGroupMember(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud CFW Address Group Member: %s", err)
	}

	d.SetId("")
	return nil
}
