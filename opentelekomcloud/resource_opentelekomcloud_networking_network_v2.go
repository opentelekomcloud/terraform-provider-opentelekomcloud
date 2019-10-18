package opentelekomcloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/provider"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/networks"
	"strings"
)

func suppressBooleanDiffs(k, old, new string, d *schema.ResourceData) bool {
	return true
}

func ExtractValSFromNid(s string) (string, string) {
	rgs := strings.Split(s, ":")
	if len(rgs) >= 2 {
		log.Printf("[DEBUG] ExtractValSFromNid: %s:%s from (%s)", rgs[0], rgs[1], s)
		return rgs[0], rgs[1]
	}
	log.Printf("[DEBUG] ExtractValSFromNid: true:'%s' from (%s)", s, s)
	return "true", s
}

func ExtractValFromNid(s string) (bool, string) {
	sasu, id := ExtractValSFromNid(s)
	asu, err := strconv.ParseBool(sasu)
	if err != nil {
		return true, id // Should never occur?
	}
	return asu, id
}

func FormatNidFromValS(asu string, id string) string {
	// Causing problems with instance network lookups right now
	//return fmt.Sprintf("%s:%s", asu, id)
	return fmt.Sprintf("%s", id)
}

func suppressAsuDiff(k, old, new string, d *schema.ResourceData) bool {
	_, id_old := ExtractValSFromNid(old)
	_, id_new := ExtractValSFromNid(new)
	return id_old == id_new
}

func resourceNetworkingNetworkV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingNetworkV2Create,
		Read:   resourceNetworkingNetworkV2Read,
		Update: resourceNetworkingNetworkV2Update,
		Delete: resourceNetworkingNetworkV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"admin_state_up": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Computed:     true,
				ValidateFunc: validateTrueOnly,
			},
			"shared": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"segments": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"physical_network": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"segmentation_id": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingNetworkV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	createOpts := NetworkCreateOpts{
		networks.CreateOpts{
			Name:     d.Get("name").(string),
			TenantID: d.Get("tenant_id").(string),
		},
		MapValueSpecs(d),
	}

	asuRaw := d.Get("admin_state_up").(string)
	asu := false
	if asuRaw != "" {
		asuT, err := strconv.ParseBool(asuRaw)
		if err != nil {
			return fmt.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
		}
		asu = asuT
		//asuFake := true
		//createOpts.AdminStateUp = &asuFake //&asu
	}

	sharedRaw := d.Get("shared").(string)
	if sharedRaw != "" {
		shared, err := strconv.ParseBool(sharedRaw)
		if err != nil {
			return fmt.Errorf("shared, if provided, must be either 'true' or 'false': %v", err)
		}
		createOpts.Shared = &shared
	}

	segments := resourceNetworkingNetworkV2Segments(d)

	n := &networks.Network{}
	if len(segments) > 0 {
		providerCreateOpts := provider.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			Segments:          segments,
		}
		log.Printf("[DEBUG] Create Options: %#v", providerCreateOpts)
		n, err = networks.Create(networkingClient, providerCreateOpts).Extract()
	} else {
		log.Printf("[DEBUG] Create Options: %#v", createOpts)
		n, err = networks.Create(networkingClient, createOpts).Extract()
	}

	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud Neutron network: %s", err)
	}
	d.SetId(FormatNidFromValS(strconv.FormatBool(asu), n.ID))

	log.Printf("[INFO] Network ID: %s", n.ID)

	log.Printf("[DEBUG] Waiting for Network (%s) to become available", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNetworkActive(networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(FormatNidFromValS(strconv.FormatBool(asu), n.ID))

	return resourceNetworkingNetworkV2Read(d, meta)
}

func resourceNetworkingNetworkV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	asu, id := ExtractValSFromNid(d.Id())
	n, err := networks.Get(networkingClient, id).Extract()
	if err != nil {
		return CheckDeleted(d, err, "network")
	}

	log.Printf("[DEBUG] Retrieved Network %s: %+v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("admin_state_up", asu)
	d.Set("shared", strconv.FormatBool(n.Shared))
	d.Set("tenant_id", n.TenantID)
	d.Set("region", GetRegion(d, config))

	d.SetId(FormatNidFromValS(asu, n.ID))
	return nil
}

func resourceNetworkingNetworkV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	var updateOpts networks.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	asu := false
	if d.HasChange("admin_state_up") {
		asuRaw := d.Get("admin_state_up").(string)
		if asuRaw != "" {
			asuT, err := strconv.ParseBool(asuRaw)
			if err != nil {
				return fmt.Errorf("admin_state_up, if provided, must be either 'true' or 'false'")
			}
			asu = asuT
			//asuFake := true
			//updateOpts.AdminStateUp = &asuFake //&asu
		}
	}
	if d.HasChange("shared") {
		sharedRaw := d.Get("shared").(string)
		if sharedRaw != "" {
			shared, err := strconv.ParseBool(sharedRaw)
			if err != nil {
				return fmt.Errorf("shared, if provided, must be either 'true' or 'false': %v", err)
			}
			updateOpts.Shared = &shared
		}
	}
	asu, id := ExtractValFromNid(d.Id())

	log.Printf("[DEBUG] Updating Network %s with options: %+v", id, updateOpts)

	_, err = networks.Update(networkingClient, id, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenTelekomCloud Neutron Network: %s", err)
	}

	d.SetId(FormatNidFromValS(strconv.FormatBool(asu), id))
	return resourceNetworkingNetworkV2Read(d, meta)
}

func resourceNetworkingNetworkV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud networking client: %s", err)
	}

	_, id := ExtractValFromNid(d.Id())
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkDelete(networkingClient, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting OpenTelekomCloud Neutron Network: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceNetworkingNetworkV2Segments(d *schema.ResourceData) (providerSegments []provider.Segment) {
	segments := d.Get("segments").([]interface{})
	for _, v := range segments {
		var segment provider.Segment
		segmentMap := v.(map[string]interface{})

		if v, ok := segmentMap["physical_network"].(string); ok {
			segment.PhysicalNetwork = v
		}

		if v, ok := segmentMap["network_type"].(string); ok {
			segment.NetworkType = v
		}

		if v, ok := segmentMap["segmentation_id"].(int); ok {
			segment.SegmentationID = v
		}

		providerSegments = append(providerSegments, segment)
	}
	return
}

func waitForNetworkActive(networkingClient *golangsdk.ServiceClient, networkId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := networks.Get(networkingClient, networkId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Neutron Network: %+v", n)
		if n.Status == "DOWN" || n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, n.Status, nil
	}
}

func waitForNetworkDelete(networkingClient *golangsdk.ServiceClient, networkId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete OpenTelekomCloud Network %s.\n", networkId)

		n, err := networks.Get(networkingClient, networkId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Network %s", networkId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = networks.Delete(networkingClient, networkId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted OpenTelekomCloud Network %s", networkId)
				return n, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return n, "ACTIVE", nil
				}
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] OpenTelekomCloud Network %s still active.\n", networkId)
		return n, "ACTIVE", nil
	}
}
