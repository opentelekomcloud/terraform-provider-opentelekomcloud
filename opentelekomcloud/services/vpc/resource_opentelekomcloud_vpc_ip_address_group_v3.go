package vpc

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/addressgroup"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcIPAddressGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcIPAddressGroupV3Create,
		ReadContext:   resourceVpcIPAddressGroupV3Read,
		UpdateContext: resourceVpcIPAddressGroupV3Update,
		DeleteContext: resourceVpcIPAddressGroupV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"max_capacity": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip_set": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip_extra_set": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"remarks": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"project_id": {
				Type:     schema.TypeString,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcIPAddressGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := addressgroup.CreateOpts{
		AddressGroup: addressgroup.AddressGroupOptions{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			IpVersion:           d.Get("ip_version").(int),
			EnterpriseProjectId: d.Get("enterprise_project_id").(string),
			MaxCapacity:         d.Get("max_capacity").(int),
			IpSet:               common.ExpandToStringList(d.Get("ip_set").([]interface{})),
			IpExtraSet:          getIpExtraSet(d),
		},
	}

	addressGroup, err := addressgroup.Create(client, createOpts)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(addressGroup.ID)

	log.Printf("[DEBUG] OpenTelekomCloud VPC IP Address Group `%s` created: %#v", d.Id(), addressGroup)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceVpcIPAddressGroupV3Read(clientCtx, d, meta)
}

func resourceVpcIPAddressGroupV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	addressGroup, err := addressgroup.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "OpenTelekomCloud VPC IP Address Group")
	}

	ipExtraSet := make([]interface{}, 0, len(addressGroup.IPExtraSet))
	for _, v := range addressGroup.IPExtraSet {
		ipExtraSet = append(ipExtraSet, map[string]interface{}{
			"ip":      v.IP,
			"remarks": v.Remark,
		})
	}

	mErr := multierror.Append(nil,
		d.Set("name", addressGroup.Name),
		d.Set("description", addressGroup.Description),
		d.Set("ip_version", addressGroup.IPVersion),
		d.Set("enterprise_project_id", addressGroup.EnterpriseProjectID),
		d.Set("max_capacity", addressGroup.MaxCapacity),
		d.Set("ip_set", addressGroup.IPSet),
		d.Set("ip_extra_set", ipExtraSet),
		d.Set("project_id", addressGroup.TenantID),
		d.Set("created_at", addressGroup.CreatedAt),
		d.Set("updated_at", addressGroup.UpdatedAt),
		d.Set("status", addressGroup.Status),
		d.Set("status_message", addressGroup.StatusMessage),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVpcIPAddressGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	updateOpts := addressgroup.UpdateOpts{
		AddressGroup: addressgroup.UpdateAddressGroupOptions{},
	}

	if d.HasChange("name") {
		updateOpts.AddressGroup.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updateOpts.AddressGroup.Description = d.Get("description").(string)
	}

	if d.HasChange("ip_set") {
		updateOpts.AddressGroup.IpSet = common.ExpandToStringList(d.Get("ip_set").([]interface{}))
	}

	if d.HasChange("ip_extra_set") {
		updateOpts.AddressGroup.IpExtraSet = getIpExtraSet(d)
	}

	log.Printf("[DEBUG] Updating VPC IP Address Group %s with options: %#v", d.Id(), updateOpts)
	_, err = addressgroup.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud VPC IP Address Group: %s", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceVpcIPAddressGroupV3Read(clientCtx, d, meta)
}

func resourceVpcIPAddressGroupV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = addressgroup.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud VPC IP Address Group: %s", err)
	}

	d.SetId("")
	return nil
}

func getIpExtraSet(d *schema.ResourceData) []addressgroup.IpExtraSetOption {
	ipExtraSetRaw := d.Get("ip_extra_set").([]interface{})
	ipExtraSet := make([]addressgroup.IpExtraSetOption, 0, len(ipExtraSetRaw))
	for _, v := range ipExtraSetRaw {
		item := v.(map[string]interface{})
		ipExtraSet = append(ipExtraSet, addressgroup.IpExtraSetOption{
			Ip:      item["ip"].(string),
			Remarks: item["remarks"].(string),
		})
	}
	return ipExtraSet
}
