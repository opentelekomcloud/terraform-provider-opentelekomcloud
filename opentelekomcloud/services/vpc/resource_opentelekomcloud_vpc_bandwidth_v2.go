package vpc

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	legacyBandwidths "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/bandwidths"
	bandwidthsV1 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/bandwidths"
	bandwidthsV2 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v2/bandwidths"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const allGrantedEnterpriseProjects = "all_granted_eps"

func ResourceBandwidthV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBandwidthV2Create,
		ReadContext:   resourceBandwidthV2Read,
		UpdateContext: resourceBandwidthV2Update,
		DeleteContext: resourceBandwidthV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(
						regexp.MustCompile(`^[\w-.]+$`),
						"The value is a string that can contain letters, digits, underscores (_), hyphens (-), and periods (.).",
					),
				),
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publicip_info": bandwidthPublicIPInfoSchema(),
		},
	}
}

func resourceBandwidthV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	opts := bandwidthsV2.CreateOpts{
		Name:                d.Get("name").(string),
		Size:                d.Get("size").(int),
		EnterpriseProjectId: config.GetEnterpriseProjectID(d, "0"),
		PublicBorderGroup:   d.Get("public_border_group").(string),
	}
	bandwidth, err := bandwidthsV2.Create(client, opts)
	if err != nil {
		return fmterr.Errorf("error creating bandwidth: %w", err)
	}
	d.SetId(bandwidth.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceBandwidthV2Read(clientCtx, d, meta)
}

func resourceBandwidthV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	enterpriseProjectID := config.GetEnterpriseProjectID(d)
	if enterpriseProjectID == "" {
		enterpriseProjectID = allGrantedEnterpriseProjects
	}
	bandwidth, err := findBandwidthV1(client, d.Id(), enterpriseProjectID)
	if err != nil {
		return fmterr.Errorf("error reading bandwidth: %w", err)
	}
	if bandwidth == nil {
		d.SetId("")
		return nil
	}

	mErr := multierror.Append(
		d.Set("name", bandwidth.Name),
		d.Set("size", bandwidth.Size),
		d.Set("status", bandwidth.Status),
		d.Set("share_type", bandwidth.ShareType),
		d.Set("bandwidth_type", bandwidth.BandwidthType),
		d.Set("charge_mode", bandwidth.ChargeMode),
		d.Set("billing_info", bandwidth.BillingInfo),
		d.Set("tenant_id", bandwidth.TenantId),
		d.Set("enterprise_project_id", bandwidth.EnterpriseProjectID),
		d.Set("public_border_group", bandwidth.PublicBorderGroup),
		d.Set("created_at", bandwidth.CreatedAt),
		d.Set("updated_at", bandwidth.UpdatedAt),
		d.Set("publicip_info", flattenBandwidthPublicIPs(bandwidth.PublicipInfo)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting bandwidth fields: %w", err)
	}

	return nil
}

func resourceBandwidthV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	opts := legacyBandwidths.UpdateOpts{}
	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}
	if d.HasChange("size") {
		opts.Size = d.Get("size").(int)
	}
	_, err = legacyBandwidths.Update(client, d.Id(), opts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating bandwidth: %w", err)
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceBandwidthV2Read(clientCtx, d, meta)
}

func resourceBandwidthV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v2 client: %w", err)
	}

	if err := bandwidthsV2.Delete(client, d.Id()); err != nil {
		return fmterr.Errorf("error deleting bandwidth: %w", err)
	}

	return nil
}

func bandwidthPublicIPInfoSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ipv6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}},
	}
}

func findBandwidthV1(client *golangsdk.ServiceClient, id, enterpriseProjectID string) (*bandwidthsV1.BandWidth, error) {
	allBandwidths, err := bandwidthsV1.List(client, bandwidthsV1.ListOpts{
		EnterpriseProjectId: enterpriseProjectID,
	})
	if err != nil {
		return nil, err
	}
	for i := range allBandwidths {
		if allBandwidths[i].ID == id {
			return &allBandwidths[i], nil
		}
	}
	return nil, nil
}

func flattenBandwidthPublicIPs(publicIPs []bandwidthsV1.PublicIpinfo) []map[string]interface{} {
	result := make([]map[string]interface{}, len(publicIPs))
	for i, publicIP := range publicIPs {
		result[i] = map[string]interface{}{
			"id":           publicIP.PublicipId,
			"address":      publicIP.PublicipAddress,
			"ipv6_address": publicIP.Publicipv6Address,
			"ip_version":   publicIP.IPVersion,
			"type":         publicIP.PublicipType,
		}
	}
	return result
}
