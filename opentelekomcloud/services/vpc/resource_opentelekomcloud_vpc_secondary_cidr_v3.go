package vpc

import (
	"bytes"
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	VpcV3 "github.com/opentelekomcloud/gophertelekomcloud/openstack/vpc/v3/vpcs"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceVpcSecondaryCidrV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVpcSecondaryCidrV3Create,
		ReadContext:   resourceVpcSecondaryCidrV3Read,
		UpdateContext: resourceVpcSecondaryCidrV3Update,
		DeleteContext: resourceVpcSecondaryCidrV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidrs": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
			},
		},
	}
}

func resourceVpcSecondaryCidrV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	vpcID := d.Get("vpc_id").(string)
	cidrs := common.ExpandToStringListBySet(d.Get("cidrs").(*schema.Set))

	if _, err := VpcV3.AddSecondaryCidr(client, vpcID, VpcV3.CidrOpts{
		Vpc: &VpcV3.AddExtendCidrOption{ExtendCidrs: cidrs},
	}); err != nil {
		return fmterr.Errorf("error adding secondary CIDRs to VPC %s: %w", vpcID, err)
	}

	d.SetId(vpcID)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceVpcSecondaryCidrV3Read(clientCtx, d, meta)
}

func resourceVpcSecondaryCidrV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	vpcInfo, err := VpcV3.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "vpc_secondary_cidr_v3")
	}

	mErr := multierror.Append(
		d.Set("region", config.GetRegion(d)),
		d.Set("vpc_id", d.Id()),
		d.Set("cidrs", vpcInfo.SecondaryCidrs),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmterr.Errorf("error setting vpc_secondary_cidr_v3 fields: %w", err)
	}

	return nil
}

func resourceVpcSecondaryCidrV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	removed, added := common.GetSetChanges(d, "cidrs")

	if removed.Len() > 0 {
		if _, err := VpcV3.RemoveSecondaryCidr(client, d.Id(), VpcV3.CidrOpts{
			Vpc: &VpcV3.AddExtendCidrOption{ExtendCidrs: common.ExpandToStringListBySet(removed)},
		}); err != nil {
			return fmterr.Errorf("error removing secondary CIDRs from VPC %s: %w", d.Id(), err)
		}
	}

	if added.Len() > 0 {
		if _, err := VpcV3.AddSecondaryCidr(client, d.Id(), VpcV3.CidrOpts{
			Vpc: &VpcV3.AddExtendCidrOption{ExtendCidrs: common.ExpandToStringListBySet(added)},
		}); err != nil {
			return fmterr.Errorf("error adding secondary CIDRs to VPC %s: %w", d.Id(), err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceVpcSecondaryCidrV3Read(clientCtx, d, meta)
}

func resourceVpcSecondaryCidrV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	cidrs := common.ExpandToStringListBySet(d.Get("cidrs").(*schema.Set))
	if len(cidrs) == 0 {
		return nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"IN_USE"},
		Target:     []string{"REMOVED"},
		Refresh:    waitForVpcSecondaryCidrRemoved(client, d.Id(), cidrs),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error removing secondary CIDRs from VPC %s: %w", d.Id(), err)
	}

	return nil
}

// cidrInUseErrCode is returned by OTC while a subnet carved from the secondary CIDR is
// still being released. It arrives as HTTP 400, not 409.
const cidrInUseErrCode = "VPC.0602"

// waitForVpcSecondaryCidrRemoved retries the removal while OTC reports the CIDR as still
// in use, which happens while a subnet carved from it is still being released.
func waitForVpcSecondaryCidrRemoved(client *golangsdk.ServiceClient, vpcID string, cidrs []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vpc, err := VpcV3.RemoveSecondaryCidr(client, vpcID, VpcV3.CidrOpts{
			Vpc: &VpcV3.AddExtendCidrOption{ExtendCidrs: cidrs},
		})
		if err != nil {
			switch e := err.(type) {
			case golangsdk.ErrDefault404:
				return struct{}{}, "REMOVED", nil
			case golangsdk.ErrDefault400:
				if bytes.Contains(e.Body, []byte(cidrInUseErrCode)) {
					return struct{}{}, "IN_USE", nil
				}
				return nil, "", err
			case golangsdk.ErrDefault409:
				return struct{}{}, "IN_USE", nil
			default:
				return nil, "", err
			}
		}

		return vpc, "REMOVED", nil
	}
}
