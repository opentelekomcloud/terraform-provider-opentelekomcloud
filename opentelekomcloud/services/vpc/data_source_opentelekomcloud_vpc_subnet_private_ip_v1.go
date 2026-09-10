package vpc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v1/privateips"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceVpcSubnetPrivateIPV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcSubnetPrivateIPV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceVpcSubnetPrivateIPV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.VpcV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC v1 client: %w", err)
	}

	privateIP, err := privateips.Get(client, d.Get("id").(string))
	if err != nil {
		return fmterr.Errorf("error retrieving OpenTelekomCloud VPC subnet private IP: %w", err)
	}

	d.SetId(privateIP.ID)
	if err := d.Set("id", privateIP.ID); err != nil {
		return fmterr.Errorf("error setting OpenTelekomCloud VPC subnet private IP ID: %w", err)
	}
	if err := setVpcSubnetPrivateIPFields(d, privateIP, config.GetRegion(d)); err != nil {
		return fmterr.Errorf("error setting OpenTelekomCloud VPC subnet private IP fields: %w", err)
	}
	return nil
}
