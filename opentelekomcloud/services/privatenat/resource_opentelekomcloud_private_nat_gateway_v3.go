package privatenat

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v3/privatenat/natgateway"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourcePrivateNatGatewayV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivateNatGatewayV3Create,
		ReadContext:   resourcePrivateNatGatewayV3Read,
		UpdateContext: resourcePrivateNatGatewayV3Update,
		DeleteContext: resourcePrivateNatGatewayV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Small", "Medium", "Large", "Extra-large",
				}, false),
			},
			"downlink_vpcs": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"virsubnet_id": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"ngport_ip_address": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
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
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
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
			"rule_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"transit_ip_pool_size_max": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePrivateNatGatewayV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	createOpts := natgateway.CreateGatewayOpts{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Spec:                d.Get("spec").(string),
		DownlinkVpcs:        getDownlinkVpcs(d),
		Tags:                getGatewayTags(d),
		EnterpriseProjectID: d.Get("enterprise_project_id").(string),
	}

	createResp, err := natgateway.Create(client, createOpts)
	if err != nil {
		return fmterr.Errorf("error creating Private NAT Gateway: %w", err)
	}
	d.SetId(createResp.Gateway.Id)
	log.Printf("Created Private NAT Gateway %s: %#v", d.Id(), createResp.Gateway)

	return resourcePrivateNatGatewayV3Read(ctx, d, meta)
}

func resourcePrivateNatGatewayV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	getResp, err := natgateway.Get(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error fetching private nat gateway : %w", err)
	}

	natGateway := getResp.Gateway

	mErr := multierror.Append(
		d.Set("id", d.Id()),
		d.Set("name", natGateway.Name),
		d.Set("description", natGateway.Description),
		d.Set("spec", natGateway.Spec),
		d.Set("tags", setDownlinkVpcs(natGateway.DownlinkVpcs)),
		d.Set("tags", setGatewayTags(natGateway.Tags)),
		d.Set("enterprise_project_id", natGateway.EnterpriseProjectID),
		d.Set("project_id", natGateway.ProjectId),
		d.Set("status", natGateway.Status),
		d.Set("created_at", natGateway.CreatedAt),
		d.Set("updated_at", natGateway.UpdatedAt),
		d.Set("rule_max", natGateway.RuleMax),
		d.Set("transit_ip_pool_size_max", natGateway.TransitIpPoolSizeMax),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePrivateNatGatewayV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	var updateOpts natgateway.UpdateGatewayOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("spec") {
		updateOpts.Spec = d.Get("spec").(string)
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = natgateway.Update(client, d.Id(), updateOpts)
	if err != nil {
		return fmterr.Errorf("error updating Private NAT Gateway: %w", err)
	}

	return resourcePrivateNatGatewayV3Read(ctx, d, meta)
}

func resourcePrivateNatGatewayV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.NatV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	err = natgateway.Delete(client, d.Id())
	if err != nil {
		return fmterr.Errorf("error deleting Private NAT Gateway: %w", err)
	}

	d.SetId("")
	return nil
}

func getDownlinkVpcs(d *schema.ResourceData) []natgateway.DownlinkVpcOption {
	vpcsInput := d.Get("downlink_vpcs").([]interface{})
	result := make([]natgateway.DownlinkVpcOption, 0, len(vpcsInput))

	for _, vpcInputRaw := range vpcsInput {
		vpcInput := vpcInputRaw.(map[string]interface{})
		downlinkVpc := natgateway.DownlinkVpcOption{
			VirSubnetID:     vpcInput["virsubnet_id"].(string),
			NgPortIPAddress: vpcInput["ngport_ip_address"].(string),
		}
		result = append(result, downlinkVpc)
	}
	return result
}

func setDownlinkVpcs(downlinkVpcsinResp []natgateway.DownlinkVpc) []map[string]interface{} {
	var downlinkVpcs []map[string]interface{}
	for _, downlinkVpcInResp := range downlinkVpcsinResp {
		downlinkVpc := map[string]interface{}{
			"virsubnet_id":      downlinkVpcInResp.VirSubnetID,
			"ngport_ip_address": downlinkVpcInResp.NgPortIPAddress,
			"vpc_id":            downlinkVpcInResp.VpcId,
		}
		downlinkVpcs = append(downlinkVpcs, downlinkVpc)
	}
	return downlinkVpcs
}

func getGatewayTags(d *schema.ResourceData) []natgateway.TagOptions {
	tagsInput := d.Get("tags").(map[string]interface{})
	result := make([]natgateway.TagOptions, 0, len(tagsInput))

	for key, value := range tagsInput {
		result = append(result, natgateway.TagOptions{
			Key:   key,
			Value: value.(string),
		})
	}
	return result
}

func setGatewayTags(tagsOutput []natgateway.Tag) map[string]interface{} {
	result := make(map[string]interface{})
	for _, tag := range tagsOutput {
		result[tag.Key] = tag.Value
	}
	return result
}
