package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVirtualPrivateCloudV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVirtualPrivateCloudV1Create,
		ReadContext:   resourceVirtualPrivateCloudV1Read,
		UpdateContext: resourceVirtualPrivateCloudV1Update,
		DeleteContext: resourceVirtualPrivateCloudV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ // request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: common.ValidateName,
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: common.ValidateCIDR,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func addNetworkingTags(d *schema.ResourceData, config *cfg.Config, res string) error {
	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		taglist := common.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcV2Client, res, d.Id(), taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("error setting tags of VirtualPrivateCloud %s: %s", d.Id(), tagErr)
		}
	}
	return nil
}

func readNetworkingTags(d *schema.ResourceData, config *cfg.Config, resource string) error {
	client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %s", err)
	}
	resourceTags, err := tags.Get(client, resource, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("error fetching tags: %s", err)
	}

	tagMap := common.TagsToMap(resourceTags)
	if err := d.Set("tags", tagMap); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}
	return nil
}

func resourceVirtualPrivateCloudV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc client: %s", err)
	}

	createOpts := vpcs.CreateOpts{
		Name: d.Get("name").(string),
		CIDR: d.Get("cidr").(string),
	}

	n, err := vpcs.Create(vpcClient, createOpts).Extract()

	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC: %s", err)
	}
	d.SetId(n.ID)

	log.Printf("[INFO] Vpc ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcActive(vpcClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return fmterr.Errorf(
			"Error waiting for Vpc (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	if common.HasFilledOpt(d, "shared") {
		snat := d.Get("shared").(bool)
		updateOpts := vpcs.UpdateOpts{
			EnableSharedSnat: &snat,
		}
		_, err = vpcs.Update(vpcClient, d.Id(), updateOpts).Extract()
		if err != nil {
			log.Printf("[WARN] Error updating shared snat for OpenTelekomCloud Vpc: %s", err)
		}
	}

	if err := addNetworkingTags(d, config, "vpcs"); err != nil {
		return diag.FromErr(err)
	}

	return resourceVirtualPrivateCloudV1Read(ctx, d, meta)

}

func resourceVirtualPrivateCloudV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Vpc client: %s", err)
	}

	n, err := vpcs.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmterr.Errorf("error retrieving OpenTelekomCloud Vpc: %s", err)
	}

	d.Set("id", n.ID)
	d.Set("name", n.Name)
	d.Set("cidr", n.CIDR)
	d.Set("status", n.Status)
	d.Set("shared", n.EnableSharedSnat)
	d.Set("region", config.GetRegion(d))

	if err := readNetworkingTags(d, config, "vpcs"); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVirtualPrivateCloudV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud Vpc: %s", err)
	}

	var updateOpts vpcs.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("cidr") {
		updateOpts.CIDR = d.Get("cidr").(string)
	}
	if d.HasChange("shared") {
		snat := d.Get("shared").(bool)
		updateOpts.EnableSharedSnat = &snat
	}

	_, err = vpcs.Update(vpcClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud Vpc: %s", err)
	}

	// update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf("error creating OpenTelekomCloud networking client: %s", err)
		}

		tagErr := common.UpdateResourceTags(vpcV2Client, d, "vpcs", d.Id())
		if tagErr != nil {
			return fmterr.Errorf("error updating tags of VPC %s: %s", d.Id(), tagErr)
		}
	}

	return resourceVirtualPrivateCloudV1Read(ctx, d, meta)
}

func resourceVirtualPrivateCloudV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	config := meta.(*cfg.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud vpc: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcDelete(vpcClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud Vpc: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcActive(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "OK" {
			return n, "ACTIVE", nil
		}

		// If vpc status is other than Ok, send error
		if n.Status == "DOWN" {
			return nil, "", fmt.Errorf("Vpc status: '%s'", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForVpcDelete(vpcClient *golangsdk.ServiceClient, vpcId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := vpcs.Get(vpcClient, vpcId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc %s", vpcId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = vpcs.Delete(vpcClient, vpcId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud vpc %s", vpcId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}
