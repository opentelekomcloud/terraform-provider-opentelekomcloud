package privatenat

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/transitip"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourcePrivateNatTransitIpV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateNatTransitIpV3Create,
		ReadContext:   resourcePrivateNatTransitIpV3Read,
		DeleteContext: resourcePrivateNatTransitIpV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"virsubnet_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
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
			"gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePrivateNatTransitIpV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := transitip.CreateTransitIpOpts{
		VirSubnetID:         d.Get("virsubnet_id").(string),
		IpAddress:           d.Get("ip_address").(string),
		Tags:                getTransitIpTags(d),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
	}

	createResp, err := transitip.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Private NAT Transit IP: %w", err)
	}
	d.SetId(createResp.TransitIp.Id)
	log.Printf("Created Private NAT TransitIp %s: %#v", d.Id(), createResp.TransitIp)

	return resourcePrivateNatTransitIpV3Read(ctx, d, meta)
}

func resourcePrivateNatTransitIpV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	getResp, err := transitip.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching private NAT Transit IP : %w", err)
	}

	natTransitIp := getResp.TransitIp

	mErr := multierror.Append(
		d.Set("id", d.Id()),
		d.Set("virsubnet_id", natTransitIp.VirSubnetID),
		d.Set("ip_address", natTransitIp.IpAddress),
		d.Set("tags", setTransitIpTags(natTransitIp.Tags)),
		d.Set("enterprise_project_id", natTransitIp.EnterpriseProjectID),
		d.Set("project_id", natTransitIp.ProjectId),
		d.Set("status", natTransitIp.Status),
		d.Set("created_at", natTransitIp.CreatedAt),
		d.Set("updated_at", natTransitIp.UpdatedAt),
		d.Set("network_interface_id", natTransitIp.NetworkInterfaceId),
		d.Set("gateway_id", natTransitIp.GatewayId),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePrivateNatTransitIpV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = transitip.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting Private NAT Transit Ip: %w", err)
	}

	d.SetId("")
	return nil
}

func getTransitIpTags(d *schema.ResourceData) []transitip.TagOptions {
	tagsInput := d.Get("tags").(map[string]interface{})
	result := make([]transitip.TagOptions, 0, len(tagsInput))

	for key, value := range tagsInput {
		result = append(result, transitip.TagOptions{
			Key:   key,
			Value: value.(string),
		})
	}
	return result
}

func setTransitIpTags(tagsOutput []transitip.Tag) map[string]interface{} {
	result := make(map[string]interface{})
	for _, tag := range tagsOutput {
		result[tag.Key] = tag.Value
	}
	return result
}
