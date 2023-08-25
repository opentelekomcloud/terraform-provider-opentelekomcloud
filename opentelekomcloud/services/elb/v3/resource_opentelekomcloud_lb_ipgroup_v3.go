package v3

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/elb/v3/ipgroups"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceIpGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpGroupV3Create,
		ReadContext:   resourceIpGroupV3Read,
		UpdateContext: resourceIpGroupV3Update,
		DeleteContext: resourceIpGroupV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 255),
			},
			"ip_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listeners": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

func getIpList(d *schema.ResourceData) []ipgroups.IpGroupOption {
	ipListRaw := d.Get("ip_list").(*schema.Set).List()
	var ipList []ipgroups.IpGroupOption

	for _, ip := range ipListRaw {
		ipRaw := ip.(map[string]interface{})

		ipList = append(ipList, ipgroups.IpGroupOption{
			Ip:          ipRaw["ip"].(string),
			Description: ipRaw["description"].(string),
		})
	}

	return ipList
}

func resourceIpGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	createOpts := ipgroups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		ProjectId:   d.Get("project_id").(string),
		IpList:      getIpList(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	ipGroup, err := ipgroups.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating LoadBalancer Ip Address Group: %w", err)
	}

	// If all has been successful, set the ID on the resource
	d.SetId(ipGroup.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceIpGroupV3Read(clientCtx, d, meta)
}

func resourceIpGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	ipGroup, err := ipgroups.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ipGroupV3")
	}
	log.Printf("[DEBUG] Retrieved Ip Group %s: %#v", d.Id(), ipGroup)

	var ipList []interface{}
	for _, v := range ipGroup.IpList {
		ipList = append(ipList, map[string]interface{}{
			"ip":          v.Ip,
			"description": v.Description,
		})
	}
	var listeners []string
	for _, v := range ipGroup.Listeners {
		listeners = append(listeners, v.ID)
	}
	mErr := multierror.Append(nil,
		d.Set("name", ipGroup.Name),
		d.Set("description", ipGroup.Description),
		d.Set("created_at", ipGroup.CreatedAt),
		d.Set("updated_at", ipGroup.UpdatedAt),
		d.Set("project_id", ipGroup.ProjectId),
		d.Set("ip_list", ipList),
		d.Set("listeners", listeners),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceIpGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	var updateOpts ipgroups.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("ip_list") {
		updateOpts.IpList = getIpList(d)
		if updateOpts.IpList == nil {
			err = batchDeleteAllIps(client, d)
			if err != nil {
				return diag.FromErr(err)
			}
			clientCtx := common.CtxWithClient(ctx, client, keyClient)
			return resourceIpGroupV3Read(clientCtx, d, meta)
		}
	}

	log.Printf("[DEBUG] Updating Ip Group %s with options: %#v", d.Id(), updateOpts)

	err = ipgroups.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating the Ip Group: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClient)
	return resourceIpGroupV3Read(clientCtx, d, meta)
}

func batchDeleteAllIps(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	ipGroup, err := ipgroups.Get(client, d.Id())
	if err != nil {
		return fmt.Errorf("error getting the Ip Group: %w", err)
	}
	var ipList []ipgroups.IpList
	for _, v := range ipGroup.IpList {
		ipList = append(ipList, ipgroups.IpList{
			Ip: v.Ip,
		})
	}
	_, err = ipgroups.DeleteIpFromList(client,
		ipGroup.ID,
		ipgroups.BatchDeleteOpts{IpList: ipList})
	if err != nil {
		return fmt.Errorf("error deleting the Ips from Ip Group: %w", err)
	}
	var updateOpts ipgroups.UpdateOpts
	if d.HasChanges("name", "description") {
		updateOpts.Name = d.Get("name").(string)
		updateOpts.Description = d.Get("description").(string)
		err = ipgroups.Update(client, d.Id(), updateOpts)
		if err != nil {
			return fmt.Errorf("error updating the Ip Group: %w", err)
		}
	}
	return nil
}

func resourceIpGroupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClient, func() (*golangsdk.ServiceClient, error) {
		return config.ElbV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(ErrCreateClient, err)
	}

	log.Printf("[DEBUG] Deleting Ip Group: %s", d.Id())
	if err := ipgroups.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("error deleting the Ip Group: %w", err)
	}

	return nil
}
