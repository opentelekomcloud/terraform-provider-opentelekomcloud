package cfw

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	list "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/blackwhitelist"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCfwBlacklistWhitelistRuleV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWBlacklistWhitelistRuleV1Create,
		ReadContext:   resourceCFWBlacklistWhitelistRuleV1Read,
		UpdateContext: resourceCFWBlacklistWhitelistRuleV1Update,
		DeleteContext: resourceCFWBlacklistWhitelistRuleV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: customImportBlacklistWhitelistRule,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				StateFunc: func(v interface{}) string {
					return strings.TrimSpace(v.(string))
				},
			},
			"list_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					4, 5,
				}),
			},
			"direction": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"address_type": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					0, 1,
				}),
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					-1, 1, 6, 17, 58,
				}),
			},
			"port": {
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
		},
	}
}

func resourceCFWBlacklistWhitelistRuleV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := list.CreateOpts{
		ObjectID:    d.Get("object_id").(string),
		ListType:    d.Get("list_type").(int),
		Direction:   InterfaceToIntPtr(d.Get("direction")),
		AddressType: InterfaceToIntPtr(d.Get("address_type")),
		Address:     d.Get("address").(string),
		Protocol:    d.Get("protocol").(int),
		Port:        d.Get("port").(string),
		Description: d.Get("description").(string),
	}

	blacklistWhitelistRule, err := list.CreateBlacklistOrWhitelistRule(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud CFW blacklist or whitelist rule from result: %w", err)
	}
	log.Printf("[DEBUG] Create CFW blacklist or whitelist rule %s: %#v", blacklistWhitelistRule.Id, blacklistWhitelistRule)

	d.SetId(blacklistWhitelistRule.Id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWBlacklistWhitelistRuleV1Read(clientCtx, d, meta)
}

func resourceCFWBlacklistWhitelistRuleV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	objectId := d.Get("object_id").(string)
	listType := d.Get("list_type").(int)
	address := d.Get("address").(string)
	blacklistWhitelistRule, err := list.GetBlacklistOrWhitelistRule(client, objectId, listType, address)
	if err != nil {
		return fmterr.Errorf("error fetching CFW blacklist or whitelist rule: %w", err)
	}
	d.SetId(blacklistWhitelistRule.ListId)

	log.Printf("[DEBUG] Retrieved blacklist or whitelist rule %s: %#v", d.Id(), blacklistWhitelistRule)

	mErr := multierror.Append(nil,
		d.Set("id", d.Id()),
		d.Set("direction", blacklistWhitelistRule.Direction),
		d.Set("address_type", blacklistWhitelistRule.AddressType),
		d.Set("address", blacklistWhitelistRule.Address),
		d.Set("protocol", blacklistWhitelistRule.Protocol),
		d.Set("port", blacklistWhitelistRule.Port),
		d.Set("description", blacklistWhitelistRule.Description),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCFWBlacklistWhitelistRuleV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	updateOpts := list.UpdateOpts{
		Direction:   InterfaceToIntPtr(d.Get("direction")),
		AddressType: InterfaceToIntPtr(d.Get("address_type")),
		Address:     d.Get("address").(string),
		Protocol:    d.Get("protocol").(int),
		Port:        d.Get("port").(string),
		Description: d.Get("description").(string),
	}

	updatedBlacklistWhitelistRule, err := list.UpdateBlacklistOrWhitelistRule(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud CFW Blacklist or Whitelist rule %s: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] Updated CFW Blacklist or Whitelist rule '%s': %#v", updatedBlacklistWhitelistRule.Id, updatedBlacklistWhitelistRule)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceCFWBlacklistWhitelistRuleV1Read(clientCtx, d, meta)
}

func resourceCFWBlacklistWhitelistRuleV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud CFW Blacklist or Whitelist rule %s", d.Id())

	err = list.DeleteBlacklistOrWhitelistRule(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud CFW Blacklist or Whitelist rule: %s", err)
	}

	d.SetId("")
	return nil
}

func customImportBlacklistWhitelistRule(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <object_id>/<list_type>/<address>")
	}

	listType, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error converting list_type to int: %s", err)
	}
	mErr := multierror.Append(nil,
		d.Set("object_id", strings.TrimSpace(parts[0])),
		d.Set("list_type", listType),
		d.Set("address", parts[2]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
