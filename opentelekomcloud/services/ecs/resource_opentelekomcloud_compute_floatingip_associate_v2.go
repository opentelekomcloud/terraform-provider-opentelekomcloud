package ecs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/floatingips"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	nfloatingips "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/layer3/floatingips"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceComputeFloatingIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeFloatingIPAssociateV2Create,
		ReadContext:   resourceComputeFloatingIPAssociateV2Read,
		DeleteContext: resourceComputeFloatingIPAssociateV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		DeprecationMessage: "Please use `opentelekomcloud_networking_floatingip_associate_v2` resource instead.",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"floating_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: common.SuppressComputedFixedWhenFloatingIp,
			},
		},
	}
}

func resourceComputeFloatingIPAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	floatingIP := d.Get("floating_ip").(string)
	fixedIP := d.Get("fixed_ip").(string)
	instanceId := d.Get("instance_id").(string)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: floatingIP,
		FixedIP:    fixedIP,
	}
	log.Printf("[DEBUG] Associate Options: %#v", associateOpts)

	err = floatingips.AssociateInstance(computeClient, instanceId, associateOpts).ExtractErr()
	if err != nil {
		return fmterr.Errorf("error associating Floating IP: %s", err)
	}

	// There's an API call to get this information, but it has been
	// deprecated. The Neutron API could be used, but I'm trying not
	// to mix service APIs. Therefore, a faux ID will be used.
	id := fmt.Sprintf("%s/%s/%s", floatingIP, instanceId, fixedIP)
	d.SetId(id)

	// This API call is synchronous, so Create won't return until the IP
	// is attached. No need to wait for a state.

	return resourceComputeFloatingIPAssociateV2Read(ctx, d, meta)
}

func resourceComputeFloatingIPAssociateV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	floatingIP, instanceId, fixedIP, err := ParseComputeFloatingIPAssociateId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Now check and see whether the floating IP still exists.
	// First try to do this by querying the Network API.
	networkEnabled := true
	networkClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		networkEnabled = false
	}

	var exists bool
	if networkEnabled {
		log.Printf("[DEBUG] Checking for Floating IP existence via Network API")
		exists, fixedIP, err = resourceComputeFloatingIPAssociateV2NetworkExists(networkClient, floatingIP)
	} else {
		log.Printf("[DEBUG] Checking for Floating IP existence via Compute API")
		exists, err = resourceComputeFloatingIPAssociateV2ComputeExists(computeClient, floatingIP)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		d.SetId("")
	}

	// Next, see if the instance still exists
	instance, err := servers.Get(computeClient, instanceId).Extract()
	if err != nil {
		if common.CheckDeleted(d, err, "instance") == nil {
			return nil
		}
	}

	// Finally, check and see if the floating ip is still associated with the instance.
	var associated bool
	for _, networkAddresses := range instance.Addresses {
		for _, element := range networkAddresses.([]interface{}) {
			address := element.(map[string]interface{})
			if (address["OS-EXT-IPS:type"] == "floating" && address["addr"] == floatingIP) ||
				(address["OS-EXT-IPS:type"] == "fixed" && address["addr"] == fixedIP) {
				associated = true
			}
		}
	}

	if !associated {
		d.SetId("")
	}

	// Set the attributes pulled from the composed resource ID
	mErr := multierror.Append(
		d.Set("floating_ip", floatingIP),
		d.Set("instance_id", instanceId),
		d.Set("fixed_ip", fixedIP),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceComputeFloatingIPAssociateV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	floatingIP := d.Get("floating_ip").(string)
	instanceId := d.Get("instance_id").(string)

	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: floatingIP,
	}
	log.Printf("[DEBUG] Disssociate Options: %#v", disassociateOpts)

	err = floatingips.DisassociateInstance(computeClient, instanceId, disassociateOpts).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "floating ip association")
	}

	return nil
}

func ParseComputeFloatingIPAssociateId(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("unable to determine floating ip association ID")
	}

	floatingIP := idParts[0]
	instanceId := idParts[1]
	fixedIP := idParts[2]

	return floatingIP, instanceId, fixedIP, nil
}

func resourceComputeFloatingIPAssociateV2NetworkExists(networkClient *golangsdk.ServiceClient, floatingIP string) (bool, string, error) {
	listOpts := nfloatingips.ListOpts{
		FloatingIP: floatingIP,
	}
	allPages, err := nfloatingips.List(networkClient, listOpts).AllPages()
	if err != nil {
		return false, "", err
	}

	allFips, err := nfloatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return false, "", err
	}

	if len(allFips) > 1 {
		return false, "", fmt.Errorf("there was a problem retrieving the floating IP")
	}

	if len(allFips) == 0 {
		return false, "", nil
	}

	return true, allFips[0].FixedIP, nil
}

func resourceComputeFloatingIPAssociateV2ComputeExists(computeClient *golangsdk.ServiceClient, floatingIP string) (bool, error) {
	// If the Network API isn't available, fall back to the deprecated Compute API.
	allPages, err := floatingips.List(computeClient).AllPages()
	if err != nil {
		return false, err
	}

	allFips, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return false, err
	}

	for _, f := range allFips {
		if f.IP == floatingIP {
			return true, nil
		}
	}

	return false, nil
}
