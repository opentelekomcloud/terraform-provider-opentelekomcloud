package privatenat

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/dnatrules"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourcePrivateNatDnatRuleV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateNatDnatRuleV3Create,
		ReadContext:   resourcePrivateNatDnatRuleV3Read,
		UpdateContext: resourcePrivateNatDnatRuleV3Update,
		DeleteContext: resourcePrivateNatDnatRuleV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_ip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal_service_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transit_service_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
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
		},
	}
}

func resourcePrivateNatDnatRuleV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	networkInterfaceId := d.Get("network_interface_id").(string)
	privateIpAddress := d.Get("private_ip_address").(string)
	if networkInterfaceId == "" && privateIpAddress == "" {
		return fmterr.Errorf("Error: One of network_interface_id or private_ip_address must be specified")
	}
	createOpts := dnatrules.CreatePrivateDnatOpts{
		Description:         d.Get("description").(string),
		TransitIpId:         d.Get("transit_ip_id").(string),
		NetworkInterfaceId:  networkInterfaceId,
		GatewayId:           d.Get("gateway_id").(string),
		PrivateIpAddress:    privateIpAddress,
		Protocol:            d.Get("protocol").(string),
		InternalServicePort: d.Get("internal_service_port").(string),
		TransitServicePort:  d.Get("transit_service_port").(string),
	}

	createResp, err := dnatrules.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Private NAT DNAT rule: %w", err)
	}
	d.SetId(createResp.DnatRule.Id)
	log.Printf("Created Private NAT DNAT rule %s: %#v", d.Id(), createResp.DnatRule)

	return resourcePrivateNatDnatRuleV3Read(ctx, d, meta)
}

func resourcePrivateNatDnatRuleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	getResp, err := dnatrules.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching private NAT DNAT rule : %w", err)
	}

	natDnatRule := getResp.DnatRule

	mErr := multierror.Append(
		d.Set("id", d.Id()),
		d.Set("description", natDnatRule.Description),
		d.Set("transit_ip_id", natDnatRule.TransitIpId),
		d.Set("network_interface_id", natDnatRule.NetworkInterfaceId),
		d.Set("gateway_id", natDnatRule.GatewayId),
		d.Set("private_ip_address", natDnatRule.PrivateIpAddress),
		d.Set("protocol", natDnatRule.Protocol),
		d.Set("internal_service_port", natDnatRule.InternalServicePort),
		d.Set("transit_service_port", natDnatRule.TransitServicePort),
		d.Set("project_id", natDnatRule.ProjectId),
		d.Set("type", natDnatRule.Type),
		d.Set("enterprise_project_id", natDnatRule.EnterpriseProjectId),
		d.Set("created_at", natDnatRule.CreatedAt),
		d.Set("updated_at", natDnatRule.UpdatedAt),
		d.Set("status", natDnatRule.Status),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePrivateNatDnatRuleV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	networkInterfaceId := d.Get("network_interface_id").(string)
	privateIpAddress := d.Get("private_ip_address").(string)
	if networkInterfaceId == "" && privateIpAddress == "" {
		return fmterr.Errorf("Error: One of network_interface_id or private_ip_address must be specified")
	}
	updateOpts := dnatrules.UpdatePrivateDnatOpts{
		Description:         d.Get("description").(string),
		TransitIpId:         d.Get("transit_ip_id").(string),
		NetworkInterfaceId:  networkInterfaceId,
		PrivateIpAddress:    privateIpAddress,
		Protocol:            d.Get("protocol").(string),
		InternalServicePort: d.Get("internal_service_port").(string),
		TransitServicePort:  d.Get("transit_service_port").(string),
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = dnatrules.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating Private NAT DNAT rule: %w", err)
	}

	return resourcePrivateNatGatewayV3Read(ctx, d, meta)
}

func resourcePrivateNatDnatRuleV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = dnatrules.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting Private NAT DNAT rule: %w", err)
	}

	d.SetId("")
	return nil
}
