package dns

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dns/v2/endpoints"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/subnets"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const dnsEndpointMaxIPs = 6

func ResourceDNSEndpointV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSEndpointV2Create,
		ReadContext:   resourceDNSEndpointV2Read,
		UpdateContext: resourceDNSEndpointV2Update,
		DeleteContext: resourceDNSEndpointV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[\w.-]{1,64}$`),
					"must be 1 to 64 characters long and contain only letters, digits, '_', '-' and '.'",
				),
			},
			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"inbound", "outbound"}, false),
			},
			"ip_addresses": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 2,
				MaxItems: dnsEndpointMaxIPs,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"ip_address_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_id": {
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
		},
	}
}

func resourceDNSEndpointV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	var ips []endpoints.IPAddressInfo
	for _, v := range d.Get("ip_addresses").([]interface{}) {
		m := v.(map[string]interface{})
		ips = append(ips, endpoints.IPAddressInfo{
			SubnetID: m["subnet_id"].(string),
			IP:       m["ip"].(string),
		})
	}

	ep, err := endpoints.Create(client, endpoints.CreateOpts{
		Name:        d.Get("name").(string),
		Direction:   d.Get("direction").(string),
		Region:      region,
		IPAddresses: ips,
	})
	if err != nil {
		return fmterr.Errorf("error creating DNS endpoint: %w", err)
	}
	d.SetId(ep.ID)

	if err := waitForDNSEndpointV2Active(ctx, client, ep.ID, len(ips), d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmterr.Errorf("error waiting for DNS endpoint %s to become ACTIVE: %w", ep.ID, err)
	}

	return resourceDNSEndpointV2Read(common.CtxWithClient(ctx, client, keyClientV2), d, meta)
}

func resourceDNSEndpointV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	ep, err := endpoints.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error reading DNS endpoint")
	}

	ips, err := endpoints.ListIPs(client, d.Id(), endpoints.ListIPsOpts{})
	if err != nil {
		return fmterr.Errorf("error listing IP addresses of DNS endpoint %s: %w", d.Id(), err)
	}

	// The API accepts VPC subnet IDs (neutron network IDs), but returns neutron subnet IDs.
	nwClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmterr.Errorf("error creating OpenTelekomCloud NetworkingV2 client: %w", err)
	}
	networkIDs := map[string]string{}
	for i, ip := range ips {
		if _, ok := networkIDs[ip.SubnetID]; !ok {
			s, err := subnets.Get(nwClient, ip.SubnetID).Extract()
			if err != nil {
				return fmterr.Errorf("error fetching subnet %s: %w", ip.SubnetID, err)
			}
			networkIDs[ip.SubnetID] = s.NetworkID
		}
		ips[i].SubnetID = networkIDs[ip.SubnetID]
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", ep.Name),
		d.Set("direction", ep.Direction),
		d.Set("ip_addresses", flattenDNSEndpointV2IPs(d.Get("ip_addresses").([]interface{}), ips)),
		d.Set("vpc_id", ep.VpcID),
		d.Set("status", ep.Status),
		d.Set("created_at", ep.CreateTime),
		d.Set("updated_at", ep.UpdateTime),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDNSEndpointV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if d.HasChange("name") {
		_, err := endpoints.Update(client, d.Id(), endpoints.UpdateOpts{Name: d.Get("name").(string)})
		if err != nil {
			return fmterr.Errorf("error updating DNS endpoint %s: %w", d.Id(), err)
		}
		// IP addresses are not changed yet, so the old count is expected.
		oldIPs, _ := d.GetChange("ip_addresses")
		ipCount := len(oldIPs.([]interface{}))
		if err := waitForDNSEndpointV2Active(ctx, client, d.Id(), ipCount, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmterr.Errorf("error waiting for DNS endpoint %s to become ACTIVE: %w", d.Id(), err)
		}
	}

	if d.HasChange("ip_addresses") {
		if err := updateDNSEndpointV2IPs(ctx, client, d); err != nil {
			return fmterr.Errorf("error updating IP addresses of DNS endpoint %s: %w", d.Id(), err)
		}
	}

	return resourceDNSEndpointV2Read(common.CtxWithClient(ctx, client, keyClientV2), d, meta)
}

func resourceDNSEndpointV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.DnsV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if err := endpoints.Delete(client, d.Id()); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DNS endpoint")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE", "PENDING_DELETE"},
		Target:       []string{"DELETED"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			ep, err := endpoints.Get(client, d.Id())
			if err != nil {
				var notFound golangsdk.ErrDefault404
				if errors.As(err, &notFound) {
					return true, "DELETED", nil
				}
				return nil, "", err
			}
			return ep, ep.Status, nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for DNS endpoint %s to be deleted: %w", d.Id(), err)
	}

	d.SetId("")
	return nil
}

// updateDNSEndpointV2IPs binds new and unbinds removed IP addresses one at a time,
// keeping the number of addresses within the 2..6 range allowed by the API.
func updateDNSEndpointV2IPs(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("ip_addresses")

	type boundIP struct{ id, subnetID string }
	bound := map[string]boundIP{}
	for _, v := range oldRaw.([]interface{}) {
		m := v.(map[string]interface{})
		bound[m["ip"].(string)] = boundIP{id: m["ip_address_id"].(string), subnetID: m["subnet_id"].(string)}
	}

	var toAdd []endpoints.IPAddressInfo
	kept := map[string]bool{}
	for _, v := range newRaw.([]interface{}) {
		m := v.(map[string]interface{})
		ip, subnetID := m["ip"].(string), m["subnet_id"].(string)
		b, ok := bound[ip]
		if ok && b.subnetID == subnetID && !kept[ip] {
			kept[ip] = true
			continue
		}
		if ok {
			// The address was carried over from state into an element whose subnet has changed.
			ip = ""
		}
		toAdd = append(toAdd, endpoints.IPAddressInfo{SubnetID: subnetID, IP: ip})
	}

	var toRemove []string
	for ip, b := range bound {
		if !kept[ip] {
			toRemove = append(toRemove, b.id)
		}
	}

	count := len(bound)
	for len(toAdd) > 0 || len(toRemove) > 0 {
		if len(toAdd) > 0 && count < dnsEndpointMaxIPs {
			if _, err := endpoints.AddIP(client, d.Id(), toAdd[0]); err != nil {
				return fmt.Errorf("error binding IP address in subnet %s: %w", toAdd[0].SubnetID, err)
			}
			toAdd = toAdd[1:]
			count++
		} else {
			if err := endpoints.RemoveIP(client, d.Id(), toRemove[0]); err != nil {
				return fmt.Errorf("error unbinding IP address %s: %w", toRemove[0], err)
			}
			toRemove = toRemove[1:]
			count--
		}
		if err := waitForDNSEndpointV2Active(ctx, client, d.Id(), count, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return err
		}
	}
	return nil
}

// waitForDNSEndpointV2Active waits until the endpoint and all of its ipCount IP addresses are ACTIVE.
func waitForDNSEndpointV2Active(ctx context.Context, client *golangsdk.ServiceClient, id string, ipCount int, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"ACTIVE"},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
		Refresh: func() (interface{}, string, error) {
			ep, err := endpoints.Get(client, id)
			if err != nil {
				return nil, "", err
			}
			if ep.Status == "ERROR" {
				return ep, "", fmt.Errorf("endpoint is in ERROR status")
			}
			if ep.Status != "ACTIVE" {
				return ep, "PENDING", nil
			}
			ips, err := endpoints.ListIPs(client, id, endpoints.ListIPsOpts{})
			if err != nil {
				return nil, "", err
			}
			if len(ips) != ipCount {
				return ep, "PENDING", nil
			}
			for _, ip := range ips {
				if ip.Status == "ERROR" {
					return ep, "", fmt.Errorf("IP address %s is in ERROR status: %s", ip.IP, ip.ErrorInfo)
				}
				if ip.Status != "ACTIVE" {
					return ep, "PENDING", nil
				}
			}
			return ep, "ACTIVE", nil
		},
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// flattenDNSEndpointV2IPs orders the IP addresses returned by the API like the prior list,
// so that elements without an explicit `ip` don't produce a diff. Addresses unknown to the
// prior list (e.g. on import) are appended in API order.
func flattenDNSEndpointV2IPs(prior []interface{}, ips []endpoints.IPAddress) []interface{} {
	used := make([]bool, len(ips))
	find := func(match func(endpoints.IPAddress) bool) int {
		for i, ip := range ips {
			if !used[i] && match(ip) {
				used[i] = true
				return i
			}
		}
		return -1
	}

	slots := make([]int, len(prior))
	// Explicit addresses go first, so that elements without one can't take them.
	for j, v := range prior {
		m, _ := v.(map[string]interface{})
		addr, _ := m["ip"].(string)
		slots[j] = -1
		if addr != "" {
			slots[j] = find(func(ip endpoints.IPAddress) bool { return ip.IP == addr })
		}
	}
	// The rest, including addresses that are gone, are matched by subnet.
	for j, v := range prior {
		m, _ := v.(map[string]interface{})
		subnetID, _ := m["subnet_id"].(string)
		if slots[j] == -1 {
			slots[j] = find(func(ip endpoints.IPAddress) bool { return ip.SubnetID == subnetID })
		}
	}

	res := make([]interface{}, 0, len(ips))
	add := func(ip endpoints.IPAddress) {
		res = append(res, map[string]interface{}{
			"subnet_id":     ip.SubnetID,
			"ip":            ip.IP,
			"ip_address_id": ip.ID,
		})
	}
	for _, i := range slots {
		if i >= 0 {
			add(ips[i])
		}
	}
	for i, ip := range ips {
		if !used[i] {
			add(ip)
		}
	}
	return res
}
