package ecs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/servergroups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceComputeServerGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeServerGroupV2Create,
		ReadContext:   resourceComputeServerGroupV2Read,
		UpdateContext: nil,
		DeleteContext: resourceComputeServerGroupV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: validatePolicy,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeServerGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	createOpts := ServerGroupCreateOpts{
		servergroups.CreateOpts{
			Name:     d.Get("name").(string),
			Policies: resourceServerGroupPoliciesV2(d),
		},
		common.MapValueSpecs(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newSG, err := servergroups.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating ServerGroup: %s", err)
	}

	d.SetId(newSG.ID)

	return resourceComputeServerGroupV2Read(ctx, d, meta)
}

func resourceComputeServerGroupV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	sg, err := servergroups.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "server group")
	}

	log.Printf("[DEBUG] Retrieved ServerGroup %s: %+v", d.Id(), sg)

	mErr := multierror.Append(
		d.Set("name", sg.Name),
		d.Set("policies", sg.Policies),
		d.Set("members", sg.Members),
		d.Set("region", config.GetRegion(d)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceComputeServerGroupV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	log.Printf("[DEBUG] Deleting ServerGroup %s", d.Id())
	if err := servergroups.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return fmterr.Errorf("error deleting ServerGroup: %s", err)
	}

	return nil
}

func resourceServerGroupPoliciesV2(d *schema.ResourceData) []string {
	rawPolicies := d.Get("policies").([]interface{})
	policies := make([]string, len(rawPolicies))
	for i, raw := range rawPolicies {
		policies[i] = raw.(string)
	}
	return policies
}

func validatePolicy(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	policyData := d.Get("policies").([]interface{})
	if policyData[0] != "anti-affinity" {
		return fmt.Errorf("only 'anti-affinity' policies are supported")
	}
	return nil
}
