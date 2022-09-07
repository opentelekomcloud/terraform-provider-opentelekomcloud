// Package bms This set of code handles all functions required to configure addresses
// on an opentelekomcloud_compute_bms_server_v2 datasource.
//
// This is a complicated task because it's not possible to obtain all
// information in a single API call. In fact, it even traverses multiple
// OpenTelekomCloud services.
//
// The end result, from the user's point of view, is a structured set of
// understandable network information within the instance resource.
package bms

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/servers"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/services/ecs"
)

// ServerNICS is a structured representation of a Gophercloud servers.Server
// virtual NIC.
type ServerNICS struct {
	IP      string
	MAC     string
	Type    string
	Version float64
}

// ServerAddress is a collection of ServerNICS, grouped by the
// network name. An instance/server could have multiple NICs on the same
// network.
type ServerAddress struct {
	NetworkName string
	ServerNICS  []ServerNICS
}

// ServerNetwork represents a collection of network information that a
// Terraform instance needs to satisfy all network information requirements.
type ServerNetwork struct {
	UUID          string
	Name          string
	Port          string
	FixedIP       string
	AccessNetwork bool
}

// getAllServerNetwork loops through the networks defined in the Terraform
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
func getAllServerNetwork(d *schema.ResourceData, meta interface{}) ([]ServerNetwork, error) {
	var serverNetworks []ServerNetwork

	networks := d.Get("network").([]interface{})
	for _, v := range networks {
		network := v.(map[string]interface{})
		networkID := network["uuid"].(string)
		networkName := network["name"].(string)
		portID := network["port"].(string)

		if networkID == "" && networkName == "" && portID == "" {
			return nil, fmt.Errorf(
				"at least one of network.uuid, network.name, or network.port must be set.")
		}

		// If a user specified both an ID and name, that makes things easy
		// since both name and ID are already satisfied. No need to query
		// further.
		if networkID != "" && networkName != "" {
			v := ServerNetwork{
				UUID:          networkID,
				Name:          networkName,
				Port:          portID,
				FixedIP:       network["fixed_ip_v4"].(string),
				AccessNetwork: network["access_network"].(bool),
			}
			serverNetworks = append(serverNetworks, v)
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

		networkInfo, err := ecs.GetInstanceNetworkInfo(d, meta, queryType, queryTerm)
		if err != nil {
			return nil, err
		}

		v := ServerNetwork{
			UUID:          networkInfo["uuid"].(string),
			Name:          networkInfo["name"].(string),
			Port:          portID,
			FixedIP:       network["fixed_ip_v4"].(string),
			AccessNetwork: network["access_network"].(bool),
		}

		serverNetworks = append(serverNetworks, v)
	}

	log.Printf("[DEBUG] getAllServerNetworks: %#v", serverNetworks)
	return serverNetworks, nil
}

// getInstanceAddresses parses a Gophercloud server.Server's Address field into
// a structured InstanceAddresses struct.
func getServerAddresses(addresses map[string]interface{}) []ServerAddress {
	var allServerAddresses []ServerAddress

	for networkName, v := range addresses {
		instanceAddresses := ServerAddress{
			NetworkName: networkName,
		}

		instanceNIC := ServerNICS{}
		for _, v := range v.([]interface{}) {
			v := v.(map[string]interface{})

			//	if v["OS-EXT-IPS:type"] == "fixed" {
			switch v["version"].(float64) {
			case 6:
				instanceNIC.IP = fmt.Sprintf("[%s]", v["addr"].(string))
			default:
				instanceNIC.IP = v["addr"].(string)
			}
			// }
			if v, ok := v["OS-EXT-IPS-MAC:mac_addr"].(string); ok {
				instanceNIC.MAC = v
			}

			instanceNIC.Type = v["OS-EXT-IPS:type"].(string)
			instanceNIC.Version = v["version"].(float64)

			instanceAddresses.ServerNICS = append(instanceAddresses.ServerNICS, instanceNIC)
		}
		allServerAddresses = append(allServerAddresses, instanceAddresses)
	}

	log.Printf("[DEBUG] Addresses: %#v", addresses)
	log.Printf("[DEBUG] allServerAddresses: %#v", allServerAddresses)

	return allServerAddresses
}

func expandBmsInstanceNetworks(allInstanceNetworks []ServerNetwork) []servers.Network {
	var networks []servers.Network
	for _, v := range allInstanceNetworks {
		n := servers.Network{
			UUID:    v.UUID,
			Port:    v.Port,
			FixedIP: v.FixedIP,
		}
		networks = append(networks, n)
	}

	return networks
}

// flattenInstanceNetworks collects instance network information from different
// sources and aggregates it all together into a map array.
func flattenServerNetwork(d *schema.ResourceData, meta interface{}) ([]map[string]interface{}, error) {
	config := meta.(*cfg.Config)
	computeClient, err := config.ComputeV2Client(config.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	server, err := servers.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "server")
	}

	allServerAddresses := getServerAddresses(server.Addresses)
	allServerNetworks, err := getAllServerNetwork(d, meta)
	if err != nil {
		return nil, err
	}

	networks := []map[string]interface{}{}

	// If there were no instance networks returned, this means that there
	// was not a network specified in the Terraform configuration. When this
	// happens, the instance will be launched on a "default" network, if one
	// is available. If there isn't, the instance will fail to launch, so
	// this is a safe assumption at this point.
	if len(allServerNetworks) == 0 {
		for _, instanceAddresses := range allServerAddresses {
			for _, instanceNIC := range instanceAddresses.ServerNICS {
				v := map[string]interface{}{
					"name":    instanceAddresses.NetworkName,
					"ip":      instanceNIC.IP,
					"mac":     instanceNIC.MAC,
					"type":    instanceNIC.Type,
					"version": instanceNIC.Version,
				}
				networks = append(networks, v)
			}
		}

		log.Printf("[DEBUG] flattenInstanceNetworks: %#v", networks)
		return networks, nil
	}

	// Loop through all networks and addresses, merge relevant address details.
	for _, instanceNetwork := range allServerNetworks {
		for _, instanceAddresses := range allServerAddresses {
			if instanceNetwork.Name == instanceAddresses.NetworkName {
				// Only use one NIC since it's possible the user defined another NIC
				// on this same network in another Terraform network block.
				instanceNIC := instanceAddresses.ServerNICS[0]
				copy(instanceAddresses.ServerNICS, instanceAddresses.ServerNICS[1:])
				v := map[string]interface{}{
					"name":           instanceAddresses.NetworkName,
					"ip":             instanceNIC.IP,
					"uuid":           instanceNetwork.UUID,
					"port":           instanceNetwork.Port,
					"access_network": instanceNetwork.AccessNetwork,
				}
				networks = append(networks, v)
			}
		}
	}

	log.Printf("[DEBUG] flattenServerNetworks: %#v", networks)
	return networks, nil
}
