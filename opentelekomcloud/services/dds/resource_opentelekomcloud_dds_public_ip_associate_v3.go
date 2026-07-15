package dds

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/instances"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceDdsPublicIpAssociateV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsPublicIpAssociateV3Create,
		ReadContext:   resourceDdsPublicIpAssociateV3Read,
		UpdateContext: resourceDdsPublicIpAssociateV3Update,
		DeleteContext: resourceDdsPublicIpAssociateV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDdsPublicIpAssociateV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	nodeId := d.Get("node_id").(string)
	d.SetId(nodeId)

	ip := d.Get("public_ip").(string)
	ipId := d.Get("public_ip_id").(string)

	jobId, err := instances.BindEIP(client, instances.BindEIPOpts{
		NodeId:     nodeId,
		PublicIp:   ip,
		PublicIpId: ipId,
	})

	if err != nil {
		return fmterr.Errorf("error binding public ip to DDS node: %w", err)
	}

	if err := waitForJobCompleted(client, 600, *jobId); err != nil {
		return diag.FromErr(err)
	}

	return resourceDdsPublicIpAssociateV3Read(ctx, d, meta)
}

func resourceDdsPublicIpAssociateV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	nodeId := d.Id()
	listOpts := instances.ListInstanceOpts{}
	instancesList, err := instances.List(client, listOpts)
	if err != nil {
		return fmterr.Errorf("error fetching DDS instances: %w", err)
	}

	var publicIp string
	found := false
	for _, instance := range instancesList.Instances {
		for _, group := range instance.Groups {
			for _, node := range group.Nodes {
				if node.Id == nodeId {
					found = true
					if node.PublicIP != "" {
						publicIp = node.PublicIP
					}
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		log.Printf("[WARN] DDS node (%s) was not found", nodeId)
		d.SetId("")
		return nil
	}

	if publicIp != "" {
		if err = d.Set("public_ip", publicIp); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceDdsPublicIpAssociateV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	nodeId := d.Id()

	if d.HasChange("public_ip") || d.HasChange("public_ip_id") {
		jobId, err := instances.UnBindEIP(client, nodeId)
		if err != nil {
			return fmterr.Errorf("error unbinding old ip: %w", err)
		}

		if err := waitForJobCompleted(client, 600, *jobId); err != nil {
			return diag.FromErr(err)
		}

		newIp := d.Get("public_ip").(string)
		newIpId := d.Get("public_ip_id").(string)

		jobId, err = instances.BindEIP(client, instances.BindEIPOpts{
			NodeId:     nodeId,
			PublicIp:   newIp,
			PublicIpId: newIpId,
		})
		if err != nil {
			return fmterr.Errorf("error binding new ip: %w", err)
		}

		if err := waitForJobCompleted(client, 600, *jobId); err != nil {
			return diag.FromErr(err)
		}
	}

	clientCtx := common.CtxWithClient(ctx, client, keyClientV3)
	return resourceDdsPublicIpAssociateV3Read(clientCtx, d, meta)
}

func resourceDdsPublicIpAssociateV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV3, func() (*golangsdk.ServiceClient, error) {
		return config.DdsV3Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV3Client, err)
	}

	nodeId := d.Id()

	log.Printf("[DEBUG] Unbinding public ip for DDS Node %s", nodeId)

	jobId, err := instances.UnBindEIP(client, nodeId)
	if err != nil {
		return fmterr.Errorf("error unbinding public ip from DDS node: %w", err)
	}

	if err := waitForJobCompleted(client, 600, *jobId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
