// This set of code handles all functions required to configure networking
// on an opentelekomcloud_compute_instance_v2 resource.
//
// This is a complicated task because it's not possible to obtain all
// information in a single API call. In fact, it even traverses multiple
// OpenTelekomCloud services.
//
// The end result, from the user's point of view, is a structured set of
// understandable network information within the instance resource.

package ecs

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/tenantnetworks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/networks"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/ports"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

// InstanceNIC is a structured representation of a golangsdk servers.Server
// virtual NIC.
type InstanceNIC struct {
	FixedIPv4 string
	FixedIPv6 string
	MAC       string
}

// InstanceAddresses is a collection of InstanceNICs, grouped by the
// network name. An instance/server could have multiple NICs on the same
// network.
type InstanceAddresses struct {
	NetworkName  string
	InstanceNICs []InstanceNIC
}

// InstanceNetwork represents a collection of network information that a
// Terraform instance needs to satisfy all network information requirements.
type InstanceNetwork struct {
	UUID          string
	Name          string
	Port          string
	FixedIP       string
	AccessNetwork bool
}

// getAllInstanceNetworks loops through the networks defined in the Terraform
// configuration and structures that information into something standard that
// can be consumed by both OpenTelekomCloud and Terraform.
//
// This would be simple, except we have ensure both the network name and
// network ID have been determined. This isn't just for the convenience of a
// user specifying a human-readable network name, but the network information
// returned by an OpenTelekomCloud instance only has the network name set! So if a
// user specified a network ID, there's no way to correlate it to the instance
// unless we know both the name and ID.
//
// Not only that, but we have to account for two OpenTelekomCloud network services
// running: nova-network (legacy) and Neutron (current).
//
// In addition, if a port was specified, not all of the port information
// will be displayed, such as multiple fixed and floating IPs. This resource
// isn't currently configured for that type of flexibility. It's better to
// reference the actual port resource itself.
//
// So, let's begin the journey.
func getAllInstanceNetworks(d *schema.ResourceData, meta interface{}) ([]InstanceNetwork, error) {
	var instanceNetworks []InstanceNetwork

	networksRaw := d.Get("network").([]interface{})
	for _, v := range networksRaw {
		network := v.(map[string]interface{})
		networkID := network["uuid"].(string)
		networkName := network["name"].(string)
		portID := network["port"].(string)

		if networkID == "" && networkName == "" && portID == "" {
			return nil, fmt.Errorf("at least one of network.uuid, network.name, or network.port must be set")
		}

		// If a user specified both an ID and name, that makes things easy
		// since both name and ID are already satisfied. No need to query
		// further.
		if networkID != "" && networkName != "" {
			v := InstanceNetwork{
				UUID:          networkID,
				Name:          networkName,
				Port:          portID,
				FixedIP:       network["fixed_ip_v4"].(string),
				AccessNetwork: network["access_network"].(bool),
			}
			instanceNetworks = append(instanceNetworks, v)
			continue
		}

		// But if at least one of name or ID was missing, we have to query
		// for that other piece.
		//
		// Priority is given to a port since a network ID or name usually isn't
		// specified when using a port.
		//
		// Next priority is given to the network ID since it's guaranteed to be
		// an exact match.
		queryType := "name"
		queryTerm := networkName
		if networkID != "" {
			queryType = "id"
			queryTerm = networkID
		}
		if portID != "" {
			queryType = "port"
			queryTerm = portID
		}

		networkInfo, err := GetInstanceNetworkInfo(d, meta, queryType, queryTerm)
		if err != nil {
			return nil, err
		}

		v := InstanceNetwork{
			Port:          portID,
			FixedIP:       network["fixed_ip_v4"].(string),
			AccessNetwork: network["access_network"].(bool),
		}
		if networkInfo["uuid"] != nil {
			v.UUID = networkInfo["uuid"].(string)
		}
		if networkInfo["name"] != nil {
			v.Name = networkInfo["name"].(string)
		}

		instanceNetworks = append(instanceNetworks, v)
	}

	log.Printf("[DEBUG] getAllInstanceNetworks: %#v", instanceNetworks)
	return instanceNetworks, nil
}

// GetInstanceNetworkInfo will query for network information in order to make
// an accurate determination of a network's name and a network's ID.
//
// We will try to first query the Neutron network service and fall back to the
// legacy nova-network service if that fails.
//
// If OS_NOVA_NETWORK is set, query nova-network even if Neutron is available.
// This is to be able to explicitly test the nova-network API.
func GetInstanceNetworkInfo(d *schema.ResourceData, meta interface{}, queryType, queryTerm string) (map[string]interface{}, error) {
	config := meta.(*cfg.Config)

	if _, ok := os.LookupEnv("OS_NOVA_NETWORK"); !ok {
		networkClient, err := config.NetworkingV2Client(config.GetRegion(d))
		if err == nil {
			networkInfo, err := getInstanceNetworkInfoNeutron(networkClient, queryType, queryTerm)
			if err != nil {
				return nil, fmt.Errorf("error trying to get network information from the Network API: %w", err)
			}

			return networkInfo, nil
		}
	}

	log.Printf("[DEBUG] Unable to obtain a network client")

	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud compute client: %w", err)
	}

	networkInfo, err := getInstanceNetworkInfoNovaNet(computeClient, queryType, queryTerm)
	if err != nil {
		return nil, fmt.Errorf("error trying to get network information from the Nova API: %w", err)
	}

	return networkInfo, nil
}

// getInstanceNetworkInfoNovaNet will query the os-tenant-networks API for
// the network information.
func getInstanceNetworkInfoNovaNet(client *golangsdk.ServiceClient, queryType, queryTerm string) (map[string]interface{}, error) {
	// If somehow a port ended up here, we should just error out.
	if queryType == "port" {
		return nil, fmt.Errorf("unable to query a port (%s) using the Nova API", queryTerm)
	}

	allPages, err := tenantnetworks.List(client).AllPages()
	if err != nil {
		return nil, fmt.Errorf("an error occurred while querying the Nova API for network information: %w", err)
	}

	networkList, err := tenantnetworks.ExtractNetworks(allPages)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while querying the Nova API for network information: %w", err)
	}

	var networkFound bool
	var network tenantnetworks.Network

	for _, v := range networkList {
		if queryType == "id" && v.ID == queryTerm {
			networkFound = true
			network = v
			break
		}

		if queryType == "name" && v.Name == queryTerm {
			networkFound = true
			network = v
			break
		}
	}

	if networkFound {
		v := map[string]interface{}{
			"uuid": network.ID,
			"name": network.Name,
		}

		log.Printf("[DEBUG] getInstanceNetworkInfoNovaNet: %#v", v)
		return v, nil
	}

	return nil, fmt.Errorf("could not find any matching network for %s %s", queryType, queryTerm)
}

// getInstanceNetworkInfoNeutron will query the neutron API for the network
// information.
func getInstanceNetworkInfoNeutron(client *golangsdk.ServiceClient, queryType, queryTerm string) (map[string]interface{}, error) {
	// If a port was specified, use it to look up the network ID
	// and then query the network as if a network ID was originally used.
	if queryType == "port" {
		listOpts := ports.ListOpts{
			ID: queryTerm,
		}
		allPages, err := ports.List(client, listOpts).AllPages()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve networks from the Network API: %w", err)
		}

		allPorts, err := ports.ExtractPorts(allPages)
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve networks from the Network API: %w", err)
		}

		var port ports.Port
		switch len(allPorts) {
		case 0:
			return nil, fmt.Errorf("could not find any matching port for %s %s", queryType, queryTerm)
		case 1:
			port = allPorts[0]
		default:
			return nil, fmt.Errorf("more than one port found for %s %s", queryType, queryTerm)
		}

		queryType = "id"
		queryTerm = port.NetworkID
	}

	listOpts := networks.ListOpts{
		Status: "ACTIVE",
	}

	switch queryType {
	case "name":
		listOpts.Name = queryTerm
	default:
		listOpts.ID = queryTerm
	}

	allPages, err := networks.List(client, listOpts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve networks from the Network API: %w", err)
	}

	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve networks from the Network API: %w", err)
	}

	var network networks.Network
	switch len(allNetworks) {
	case 0:
		return nil, fmt.Errorf("could not find any matching network for %s %s", queryType, queryTerm)
	case 1:
		network = allNetworks[0]
	default:
		return nil, fmt.Errorf("more than one network found for %s %s", queryType, queryTerm)
	}

	v := map[string]interface{}{
		"uuid": network.ID,
		"name": network.Name,
	}

	log.Printf("[DEBUG] getInstanceNetworkInfoNeutron: %#v", v)
	return v, nil
}

// getInstanceAddresses parses a golangsdk server.Server's Address field into
// a structured InstanceAddresses struct.
func getInstanceAddresses(addresses map[string]interface{}) []InstanceAddresses {
	var allInstanceAddresses []InstanceAddresses

	for networkName, v := range addresses {
		instanceAddresses := InstanceAddresses{
			NetworkName: networkName,
		}

		for _, v := range v.([]interface{}) {
			instanceNIC := InstanceNIC{}
			var exists bool

			v := v.(map[string]interface{})
			if v, ok := v["OS-EXT-IPS-MAC:mac_addr"].(string); ok {
				instanceNIC.MAC = v
			}

			if v["OS-EXT-IPS:type"] == "fixed" || v["OS-EXT-IPS:type"] == nil {
				switch v["version"].(float64) {
				case 6:
					instanceNIC.FixedIPv6 = fmt.Sprintf("[%s]", v["addr"].(string))
				default:
					instanceNIC.FixedIPv4 = v["addr"].(string)
				}
			}

			// To associate IPv4 and IPv6 on the right NIC,
			// key on the mac address and fill in the blanks.
			for i, v := range instanceAddresses.InstanceNICs {
				if v.MAC == instanceNIC.MAC {
					exists = true
					if instanceNIC.FixedIPv6 != "" {
						instanceAddresses.InstanceNICs[i].FixedIPv6 = instanceNIC.FixedIPv6
					}
					if instanceNIC.FixedIPv4 != "" {
						instanceAddresses.InstanceNICs[i].FixedIPv4 = instanceNIC.FixedIPv4
					}
				}
			}

			if !exists {
				instanceAddresses.InstanceNICs = append(instanceAddresses.InstanceNICs, instanceNIC)
			}
		}

		allInstanceAddresses = append(allInstanceAddresses, instanceAddresses)
	}

	log.Printf("[DEBUG] Addresses: %#v", addresses)
	log.Printf("[DEBUG] allInstanceAddresses: %#v", allInstanceAddresses)

	return allInstanceAddresses
}

// expandInstanceNetworks takes network information found in []InstanceNetwork
// and builds a golangsdk []servers.Network for use in creating an Instance.
func expandInstanceNetworks(allInstanceNetworks []InstanceNetwork) []servers.Network {
	var networkList []servers.Network
	for _, v := range allInstanceNetworks {
		n := servers.Network{
			UUID:    v.UUID,
			Port:    v.Port,
			FixedIP: v.FixedIP,
		}
		networkList = append(networkList, n)
	}

	return networkList
}

// FlattenInstanceNetworks collects instance network information from different
// sources and aggregates it all together into a map array.
func FlattenInstanceNetworks(d *schema.ResourceData, meta interface{}) ([]map[string]interface{}, error) {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud compute client: %w", err)
	}

	server, err := servers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "server")
	}

	serverNICs, err := servers.GetNICs(computeClient, server.ID).Extract()
	if err != nil {
		return nil, err
	}

	allInstanceAddresses := getInstanceAddresses(server.Addresses)
	allInstanceNetworks, err := getAllInstanceNetworks(d, meta)
	if err != nil {
		return nil, err
	}

	var networkList []map[string]interface{}

	// If there were no instance networks returned, this means that there
	// was not a network specified in the Terraform configuration. When this
	// happens, the instance will be launched on a "default" network, if one
	// is available. If there isn't, the instance will fail to launch, so
	// this is a safe assumption at this point.
	if len(allInstanceNetworks) == 0 {
		for _, instanceAddresses := range allInstanceAddresses {
			for _, instanceNIC := range instanceAddresses.InstanceNICs {
				v := map[string]interface{}{
					"name":        instanceAddresses.NetworkName,
					"fixed_ip_v4": instanceNIC.FixedIPv4,
					"fixed_ip_v6": instanceNIC.FixedIPv6,
					"mac":         instanceNIC.MAC,
				}

				for _, nicObj := range serverNICs {
					if nicObj.MACAddress == instanceNIC.MAC {
						v["uuid"] = nicObj.NetID
					}
				}
				if v["uuid"] == "" {
					log.Printf("[WARN] Could not get default network uuid")
				}

				networkList = append(networkList, v)
			}
		}

		log.Printf("[DEBUG] flattenInstanceNetworks: %#v", networkList)
		return networkList, nil
	}

	// Loop through all networks and addresses, merge relevant address details.
	for _, instanceNetwork := range allInstanceNetworks {
		for _, instanceAddresses := range allInstanceAddresses {
			// Skip if instanceAddresses has no NICs
			if len(instanceAddresses.InstanceNICs) == 0 {
				continue
			}

			if instanceNetwork.Name == instanceAddresses.NetworkName {
				// Only use one NIC since it's possible the user defined another NIC
				// on this same network in another Terraform network block.
				instanceNIC := instanceAddresses.InstanceNICs[0]
				copy(instanceAddresses.InstanceNICs, instanceAddresses.InstanceNICs[1:])

				port := instanceNetwork.Port
				if port == "" {
					port = findPortByNetID(serverNICs, instanceNetwork.UUID)
				}

				v := map[string]interface{}{
					"name":           instanceAddresses.NetworkName,
					"fixed_ip_v4":    instanceNIC.FixedIPv4,
					"fixed_ip_v6":    instanceNIC.FixedIPv6,
					"mac":            instanceNIC.MAC,
					"uuid":           instanceNetwork.UUID,
					"port":           port,
					"access_network": instanceNetwork.AccessNetwork,
				}
				networkList = append(networkList, v)
			}
		}
	}

	log.Printf("[DEBUG] flattenInstanceNetworks: %#v", networkList)
	return networkList, nil
}

func findPortByNetID(nics []servers.NIC, network string) string {
	for _, nic := range nics {
		if nic.NetID == network {
			return nic.PortID
		}
	}
	return ""
}

// GetInstanceAccessAddresses determines the best IP address to communicate
// with the instance. It does this by looping through all networks and looking
// for a valid IP address. Priority is given to a network that was flagged as
// an access_network.
func GetInstanceAccessAddresses(_ *schema.ResourceData, networks []map[string]interface{}) (string, string) {
	var hostv4, hostv6 string

	// Loop through all networks
	// If the network has a valid fixed v4 or fixed v6 address
	// and hostv4 or hostv6 is not set, set hostv4/hostv6.
	// If the network is an "access_network" overwrite hostv4/hostv6.
	for _, n := range networks {
		var accessNetwork bool

		if an, ok := n["access_network"].(bool); ok && an {
			accessNetwork = true
		}

		if fixedIPv4, ok := n["fixed_ip_v4"].(string); ok && fixedIPv4 != "" {
			if hostv4 == "" || accessNetwork {
				hostv4 = fixedIPv4
			}
		}

		if fixedIPv6, ok := n["fixed_ip_v6"].(string); ok && fixedIPv6 != "" {
			if hostv6 == "" || accessNetwork {
				hostv6 = fixedIPv6
			}
		}
	}

	log.Printf("[DEBUG] OpenTelekomCloud Instance Network Access Addresses: %s, %s", hostv4, hostv6)

	return hostv4, hostv6
}
