package cce

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cce/v3/nodes"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCceNodesV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceNodesV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_pair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"extend_params": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"kms_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"billing_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"eip_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"eip_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceCceNodesV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.CceV3Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf(cceClientError, err)
	}

	listOpts := nodes.ListOpts{
		Uid:   d.Get("node_id").(string),
		Name:  d.Get("name").(string),
		Phase: d.Get("status").(string),
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("node_id"); ok {
		listOpts.Uid = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Phase = v.(string)
	}

	refinedNodes, err := nodes.List(client, d.Get("cluster_id").(string), listOpts)

	if err != nil {
		return fmterr.Errorf("unable to retrieve Nodes: %s", err)
	}

	if len(refinedNodes) < 1 {
		return fmterr.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNodes) > 1 {
		return fmterr.Errorf("your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Node := refinedNodes[0]

	var dataVolumes []map[string]interface{}
	for _, volume := range Node.Spec.DataVolumes {
		mapping := map[string]interface{}{
			"disk_size":     volume.Size,
			"volume_type":   volume.VolumeType,
			"extend_params": volume.ExtendParam,
			"kms_id":        volume.Metadata["__system__cmkid"],
		}
		dataVolumes = append(dataVolumes, mapping)
	}

	log.Printf("[DEBUG] Retrieved Nodes using given filter %s: %+v", Node.Metadata.Id, Node)
	d.SetId(Node.Metadata.Id)

	mErr := multierror.Append(
		d.Set("node_id", Node.Metadata.Id),
		d.Set("name", Node.Metadata.Name),
		d.Set("flavor_id", Node.Spec.Flavor),
		d.Set("availability_zone", Node.Spec.Az),
		d.Set("billing_mode", Node.Spec.BillingMode),
		d.Set("status", Node.Status.Phase),
		d.Set("data_volumes", dataVolumes),
		d.Set("disk_size", Node.Spec.RootVolume.Size),
		d.Set("volume_type", Node.Spec.RootVolume.VolumeType),
		d.Set("key_pair", Node.Spec.Login.SshKey),
		d.Set("charge_mode", Node.Spec.PublicIP.Eip.Bandwidth.ChargeMode),
		d.Set("bandwidth_size", Node.Spec.PublicIP.Eip.Bandwidth.Size),
		d.Set("share_type", Node.Spec.PublicIP.Eip.Bandwidth.ShareType),
		d.Set("ip_type", Node.Spec.PublicIP.Eip.IpType),
		d.Set("server_id", Node.Status.ServerID),
		d.Set("public_ip", Node.Status.PublicIP),
		d.Set("private_ip", Node.Status.PrivateIP),
		d.Set("eip_count", Node.Spec.PublicIP.Count),
		d.Set("eip_ids", Node.Spec.PublicIP.Ids),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
