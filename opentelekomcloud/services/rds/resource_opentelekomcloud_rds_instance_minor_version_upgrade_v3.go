package rds

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceRdsInstanceMinorVersionUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsInstanceMinorVersionUpgradeCreate,
		ReadContext:   resourceRdsInstanceMinorVersionUpgradeRead,
		UpdateContext: resourceRdsInstanceMinorVersionUpgradeUpdate,
		DeleteContext: resourceRdsInstanceMinorVersionUpgradeDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delay": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRdsInstanceMinorVersionUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.RdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreateClient, err)
	}

	instanceID := d.Get("instance_id").(string)

	opts := instances.UpgradeDbVersionOpts{
		InstanceId: instanceID,
		Delay:      d.Get("delay").(bool),
	}

	jobId, err := instances.UpgradeDbVersion(client, opts)
	if err != nil {
		return fmterr.Errorf("error upgrading kernel minor version for instance (%s): %s", instanceID, err)
	}

	log.Printf("[DEBUG] Upgrade job ID: %s", *jobId)

	d.SetId(instanceID)

	return nil
}

func resourceRdsInstanceMinorVersionUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsInstanceMinorVersionUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRdsInstanceMinorVersionUpgradeDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Deleting kernel upgrade resource is not supported. This resource is only removed from the state.",
		},
	}
}
