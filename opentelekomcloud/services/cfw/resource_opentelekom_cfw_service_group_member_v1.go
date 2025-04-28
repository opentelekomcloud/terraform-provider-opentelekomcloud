package cfw

import (
	"context"
	"fmt"
	"log"
	"strings"
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

func ResourceCfwServiceGroupMemberV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCFWServiceGroupmemberV1Create,
		ReadContext:   resourceCFWServiceGroupMemberV1Read,
		DeleteContext: resourceCFWServiceGroupMemberV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: customImportServiceGroupMember,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"set_id": {
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
			"protocol": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.IntInSlice([]int{
					-1, 1, 6, 17, 58,
				}),
			},
			"source_port": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dest_port": {
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
		},
	}
}

func resourceCFWServiceGroupmemberV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := group.AddGroupMemberOpts{
		SetId: strings.TrimSpace(d.Get("set_id").(string)),
		ServiceItems: []group.ServiceItem{
			{
				Protocol:    d.Get("protocol").(int),
				SourcePort:  d.Get("source_port").(string),
				DestPort:    d.Get("dest_port").(string),
				Description: d.Get("description").(string),
			},
		},
	}
	membersData, err := group.AddGroupMember(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error getting OpenTelekomCloud CFW service group member from result: %w", err)
	}

	d.SetId(membersData.Items[0].Id)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceCFWServiceGroupMemberV1Read(clientCtx, d, meta)
}

func resourceCFWServiceGroupMemberV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	groupMembers, err := group.ListGroupMembers(client, strings.TrimSpace(d.Get("set_id").(string)))
	if err != nil {
		return fmterr.Errorf("unable to fetch OpenTelekomCloud CFW service group member: %w", err)
	}
	for _, member := range groupMembers {
		if member.ItemID == d.Id() {
			log.Printf("[DEBUG] Retrieved instance %s: %#v", d.Id(), member)

			mErr := multierror.Append(nil,
				d.Set("id", d.Id()),
				d.Set("set_id", d.Get("set_id")),
				d.Set("protocol", member.Protocol),
				d.Set("source_port", member.SourcePort),
				d.Set("dest_port", member.DestPort),
				d.Set("description", member.Description),
			)
			return diag.FromErr(mErr.ErrorOrNil())

		}
	}
	return fmterr.Errorf("unable to find OpenTelekomCloud CFW service group member or group member does not exist: %w", err)
}

func resourceCFWServiceGroupMemberV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.CfwV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	log.Printf("[DEBUG] Deleting OpenTelekomCloud CFW Service Group Member %s", d.Id())

	err = group.DeleteGroupMember(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud CFW Service Group Member: %s", err)
	}

	d.SetId("")
	return nil
}

func customImportServiceGroupMember(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <set_id>/<id>")
	}
	d.Set("set_id", strings.TrimSpace(parts[0]))
	id := strings.TrimSpace(parts[1])
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
