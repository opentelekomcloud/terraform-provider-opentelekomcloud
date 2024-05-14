package vpc

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/helper/hashcode"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"
)

func DataSourceNetworkingPortIDsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingPortIDsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"device_owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, true),
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNetworkingPortIDsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.NetworkingV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV2Client, err)
	}

	listOpts := ports.ListOpts{}
	var listOptsBuilder ports.ListOptsBuilder

	if v, ok := d.GetOk("sort_key"); ok {
		listOpts.SortKey = v.(string)
	}
	if v, ok := d.GetOk("sort_direction"); ok {
		listOpts.SortDir = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}
	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}
	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}
	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}
	if v, ok := d.GetOk("device_owner"); ok {
		listOpts.DeviceOwner = v.(string)
	}
	if v, ok := d.GetOk("mac_address"); ok {
		listOpts.MACAddress = v.(string)
	}
	if v, ok := d.GetOk("device_id"); ok {
		listOpts.DeviceID = v.(string)
	}
	listOptsBuilder = listOpts

	allPages, err := ports.List(client, listOptsBuilder).AllPages()
	if err != nil {
		return diag.Errorf("unable to list OpenTelekomCloud ports: %s", err)
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return diag.Errorf("unable to retrieve OpenTelekomCloud: %s", err)
	}

	if len(allPorts) == 0 {
		log.Printf("[DEBUG] No ports in OpenTelekomCloud found")
	}

	portsList := make([]ports.Port, 0, len(allPorts))
	portIDs := make([]string, 0, len(allPorts))

	// Filter returned Fixed IPs by a "fixed_ip".
	if v, ok := d.GetOk("fixed_ip"); ok {
		for _, p := range allPorts {
			for _, ipObject := range p.FixedIPs {
				if v.(string) == ipObject.IPAddress {
					portsList = append(portsList, p)
				}
			}
		}
		if len(portsList) == 0 {
			log.Printf("[DEBUG] No ports in OpenTelekomCloud found after the 'fixed_ip' filter")
		}
	} else {
		portsList = allPorts
	}

	securityGroups := common.ExpandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	if len(securityGroups) > 0 {
		var sgPorts []ports.Port
		for _, p := range portsList {
			for _, sg := range p.SecurityGroups {
				if common.StrSliceContains(securityGroups, sg) {
					sgPorts = append(sgPorts, p)
				}
			}
		}
		if len(sgPorts) == 0 {
			log.Printf("[DEBUG] No ports in OpenTelekomCloud found after the 'security_group_ids' filter")
		}
		portsList = sgPorts
	}

	for _, p := range portsList {
		portIDs = append(portIDs, p.ID)
	}

	log.Printf("[DEBUG] Retrieved %d ports in OpenTelekomCloud: %+v", len(portsList), portsList)
	portID := fmt.Sprintf("%d", hashcode.String(strings.Join(portIDs, "")))
	d.SetId(portID)
	mErr := multierror.Append(nil,
		d.Set("ids", portIDs),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving port (%s) fields: %s", portID, err)
	}
	return nil
}
