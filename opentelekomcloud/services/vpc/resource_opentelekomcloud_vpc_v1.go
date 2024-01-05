package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/common/tags"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"
	VpcV3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/vpcs"

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

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: common.ValidateName,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"secondary_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsCIDR,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
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

func addSecondaryCidr(d *schema.ResourceData, config *cfg.Config) error {
	vpcV3Client, err := config.NetworkingV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(errCreationV3Client, err)
	}
	cidr := d.Get("secondary_cidr").(string)
	cidrOpts := VpcV3.CidrOpts{
		Vpc: &VpcV3.AddExtendCidrOption{
			ExtendCidrs: []string{cidr}},
	}
	_, err = VpcV3.AddSecondaryCidr(vpcV3Client, d.Id(), cidrOpts)
	if err != nil {
		return fmt.Errorf("error setting secondary cidr of VirtualPrivateCloud %s: %s", d.Id(), cidr)
	}
	return nil
}

func updateSecondaryCidr(d *schema.ResourceData, config *cfg.Config) error {
	vpcV3Client, err := config.NetworkingV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(errCreationV3Client, err)
	}
	vpcV3Get, err := VpcV3.Get(vpcV3Client, d.Id())
	if err != nil {
		return fmt.Errorf("error fetching vpc: %s", err)
	}
	cidrOpts := VpcV3.CidrOpts{
		Vpc: &VpcV3.AddExtendCidrOption{
			ExtendCidrs: []string{vpcV3Get.SecondaryCidrs[0]}},
	}
	_, err = VpcV3.RemoveSecondaryCidr(vpcV3Client, d.Id(), cidrOpts)
	if err != nil {
		return fmt.Errorf("error removing old secondary cidr of VirtualPrivateCloud %s: %s", d.Id(), vpcV3Get.SecondaryCidrs[0])
	}
	cidr := d.Get("secondary_cidr").(string)
	cidrAddOpts := VpcV3.CidrOpts{
		Vpc: &VpcV3.AddExtendCidrOption{
			ExtendCidrs: []string{cidr}},
	}
	_, err = VpcV3.AddSecondaryCidr(vpcV3Client, d.Id(), cidrAddOpts)
	if err != nil {
		return fmt.Errorf("error setting new secondary cidr of VirtualPrivateCloud %s: %s", d.Id(), cidr)
	}
	return nil
}

func readSecondaryCidr(d *schema.ResourceData, config *cfg.Config) error {
	vpcV3Client, err := config.NetworkingV3Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(errCreationV3Client, err)
	}
	vpcV3Get, err := VpcV3.Get(vpcV3Client, d.Id())
	if err != nil {
		return fmt.Errorf("error fetching vpc: %s", err)
	}

	if err := d.Set("secondary_cidr", vpcV3Get.SecondaryCidrs[0]); err != nil {
		return fmt.Errorf("error setting secondary cidr: %s", err)
	}
	return nil
}

func addNetworkingTags(d *schema.ResourceData, config *cfg.Config, res string) error {
	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmt.Errorf(errCreationV2Client, err)
		}

		tagList := common.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(vpcV2Client, res, d.Id(), tagList).ExtractErr(); tagErr != nil {
			return fmt.Errorf("error setting tags of VirtualPrivateCloud %s: %s", d.Id(), tagErr)
		}
	}
	return nil
}

func readNetworkingTags(d *schema.ResourceData, config *cfg.Config, resource string) error {
	vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmt.Errorf(errCreationV2Client, err)
	}
	resourceTags, err := tags.Get(vpcV2Client, resource, d.Id()).Extract()
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
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	createOpts := vpcs.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		CIDR:        d.Get("cidr").(string),
	}

	n, err := vpcs.Create(client, createOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud VPC: %w", err)
	}

	d.SetId(n.ID)

	log.Printf("[INFO] VPC ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForVpcActive(client, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error waiting for VPC (%s) to become ACTIVE: %w", n.ID, err)
	}

	if common.HasFilledOpt(d, "shared") {
		snat := d.Get("shared").(bool)
		updateOpts := vpcs.UpdateOpts{
			EnableSharedSnat: &snat,
		}
		_, err = vpcs.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			log.Printf("[WARN] Error updating shared SNAT for OpenTelekomCloud VPC: %s", err)
		}
	}

	if err := addNetworkingTags(d, config, "vpcs"); err != nil {
		return diag.FromErr(err)
	}
	if _, ok := d.GetOk("secondary_cidr"); ok {
		if err := addSecondaryCidr(d, config); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(n.ID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVirtualPrivateCloudV1Read(clientCtx, d, meta)
}

func resourceVirtualPrivateCloudV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	n, err := vpcs.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "vpc")
	}

	mErr := multierror.Append(
		d.Set("name", n.Name),
		d.Set("description", n.Description),
		d.Set("cidr", n.CIDR),
		d.Set("status", n.Status),
		d.Set("shared", n.EnableSharedSnat),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	if err := readNetworkingTags(d, config, "vpcs"); err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk("secondary_cidr"); ok {
		if err := readSecondaryCidr(d, config); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceVirtualPrivateCloudV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	var updateOpts vpcs.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("cidr") {
		updateOpts.CIDR = d.Get("cidr").(string)
	}
	if d.HasChange("shared") {
		snat := d.Get("shared").(bool)
		updateOpts.EnableSharedSnat = &snat
	}

	_, err = vpcs.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmterr.Errorf("error updating OpenTelekomCloud VPC: %s", err)
	}

	// update tags
	if d.HasChange("tags") {
		vpcV2Client, err := config.NetworkingV2Client(config.GetRegion(d))
		if err != nil {
			return fmterr.Errorf(errCreationV2Client, err)
		}

		tagErr := common.UpdateResourceTags(vpcV2Client, d, "vpcs", d.Id())
		if tagErr != nil {
			return fmterr.Errorf("error updating tags of VPC %s: %w", d.Id(), tagErr)
		}
	}

	// update secondary cidr
	if d.HasChange("secondary_cidr") {
		if err := updateSecondaryCidr(d, config); err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV1)
	return resourceVirtualPrivateCloudV1Read(clientCtx, d, meta)
}

func resourceVirtualPrivateCloudV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcDelete(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmterr.Errorf("error deleting OpenTelekomCloud VPC: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcActive(client *golangsdk.ServiceClient, vpcID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := vpcs.Get(client, vpcID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "OK" {
			return n, "ACTIVE", nil
		}

		// If vpc status is other than Ok, send error
		if n.Status == "DOWN" {
			return nil, "", fmt.Errorf("VPC status: %s", n.Status)
		}

		return n, n.Status, nil
	}
}

func waitForVpcDelete(client *golangsdk.ServiceClient, vpcID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := vpcs.Get(client, vpcID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud VPC %s", vpcID)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		if err = vpcs.Delete(client, vpcID).ExtractErr(); err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted OpenTelekomCloud VPC %s", vpcID)
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
