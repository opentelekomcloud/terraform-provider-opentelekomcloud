package ecs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/floatingips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceComputeFloatingIPV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeFloatingIPV2Create,
		ReadContext:   resourceComputeFloatingIPV2Read,
		UpdateContext: nil,
		DeleteContext: resourceComputeFloatingIPV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		DeprecationMessage: "Please use `opentelekomcloud_networking_floatingip_v2` or `opentelekomcloud_vpc_eip_v1` resources instead",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"pool": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "admin_external_net",
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeFloatingIPV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud ComputeV2 client: %s", err)
	}

	createOpts := &floatingips.CreateOpts{
		Pool: d.Get("pool").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newFip, err := floatingips.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating Floating IP: %s", err)
	}

	d.SetId(newFip.ID)

	return resourceComputeFloatingIPV2Read(ctx, d, meta)
}

func resourceComputeFloatingIPV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	fip, err := floatingips.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "floating ip")
	}

	log.Printf("[DEBUG] Retrieved Floating IP %s: %+v", d.Id(), fip)

	me := multierror.Append(
		d.Set("pool", fip.Pool),
		d.Set("instance_id", fip.InstanceID),
		d.Set("address", fip.IP),
		d.Set("fixed_ip", fip.FixedIP),
		d.Set("region", config.GetRegion(d)),
	)

	return diag.FromErr(me.ErrorOrNil())
}

func FloatingIPV2StateRefreshFunc(computeClient *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := floatingips.Get(computeClient, d.Id()).Extract()
		if err != nil {
			err = common.CheckDeleted(d, err, "Floating IP")
			if err != nil {
				return s, "", err
			} else {
				log.Printf("[DEBUG] Successfully deleted Floating IP %s", d.Id())
				return s, "DELETED", nil
			}
		}

		log.Printf("[DEBUG] Floating IP %s still active.\n", d.Id())
		return s, "ACTIVE", nil
	}
}

func resourceComputeFloatingIPV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete Floating IP %s.\n", d.Id())

	err = floatingips.Delete(computeClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    FloatingIPV2StateRefreshFunc(computeClient, d),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Floating IP: %s", err)
	}

	return nil
}
