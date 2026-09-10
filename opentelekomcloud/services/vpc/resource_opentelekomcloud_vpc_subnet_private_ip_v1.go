package vpc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/privateips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcSubnetPrivateIPV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcSubnetPrivateIPV1Create,
		ReadContext:   resourceVpcSubnetPrivateIPV1Read,
		DeleteContext: resourceVpcSubnetPrivateIPV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVpcSubnetPrivateIPV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	created, err := privateips.Create(client, privateips.CreateOpts{
		PrivateIPs: []privateips.PrivateIPRequest{{
			SubnetID:  d.Get("subnet_id").(string),
			IPAddress: d.Get("ip_address").(string),
		}},
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC subnet private IP: %w", err)
	}
	if err := validateVpcSubnetPrivateIPCreateResponse(created); err != nil {
		return fmterr.Errorf("%w", err)
	}

	d.SetId(created[0].ID)
	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVpcSubnetPrivateIPV1Read(clientCtx, d, config)
}

func resourceVpcSubnetPrivateIPV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	privateIP, err := privateips.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OpenTelekomCloud VPC subnet private IP")
	}

	if err := setVpcSubnetPrivateIPFields(d, privateIP, config.GetRegion(d)); err != nil {
		return fmterr.Errorf("error setting OpenTelekomCloud VPC subnet private IP fields: %w", err)
	}
	return nil
}

func resourceVpcSubnetPrivateIPV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.VpcV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	if err := privateips.Delete(client, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OpenTelekomCloud VPC subnet private IP")
	}

	d.SetId("")
	return nil
}

func setVpcSubnetPrivateIPFields(d *schema.ResourceData, privateIP *privateips.PrivateIP, region string) error {
	return multierror.Append(
		d.Set("region", region),
		d.Set("subnet_id", privateIP.SubnetID),
		d.Set("ip_address", privateIP.IPAddress),
		d.Set("status", privateIP.Status),
		d.Set("tenant_id", privateIP.TenantID),
		d.Set("device_owner", privateIP.DeviceOwner),
	).ErrorOrNil()
}

func validateVpcSubnetPrivateIPCreateResponse(created []privateips.PrivateIP) error {
	if len(created) != 1 {
		return fmt.Errorf("expected one VPC subnet private IP in create response, got %d", len(created))
	}
	return nil
}
