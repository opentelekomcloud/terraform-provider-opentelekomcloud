package opentelekomcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"

	"github.com/gophercloud/gophercloud"
	//"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/backendmember"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/healthcheck"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/listeners"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/elbaas/loadbalancer_elbs"
)

func waitForELBListener(networkingClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for listener %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceELBListenerRefreshFunc(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: listener %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for listener %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceELBListenerRefreshFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		listener, err := listeners.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		// The listener resource has no Status attribute, so a successful Get is the best we can do
		return listener, "ACTIVE", nil
	}
}

func waitForELBLoadBalancer(networkingClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for loadbalancer %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceELBLoadBalancerRefreshFunc(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: loadbalancer %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for loadbalancer %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceELBLoadBalancerRefreshFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		lb, err := loadbalancer_elbs.Get(networkingClient, id).Extract()
		if err != nil {
			fmt.Printf("@@@@@@@@@@@@@@@@ resourceELBLoadBalancerRefreshFunc  err %s \n", err)

			return nil, "", err
		}

		return lb, "" /*lb.ProvisioningStatus*/, nil
	}
}

func waitForELBBackend(networkingClient *gophercloud.ServiceClient, memberID string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for member %s to become %s.", memberID, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceELBBackendRefreshFunc(networkingClient, memberID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: member %s not found: %s", memberID, err)
			}
		}
		return fmt.Errorf("Error waiting for member %s to become %s: %s", memberID, target, err)
	}

	return nil
}

func resourceELBBackendRefreshFunc(networkingClient *gophercloud.ServiceClient, memberID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		//member, err := pools.GetMember(networkingClient, poolID, memberID).Extract()
		//if err != nil {
		//	return nil, "", err
		//}

		// The member resource has no Status attribute, so a successful Get is the best we can do
		//return member, "ACTIVE", nil
		return nil, "ACTIVE", nil
	}
}

func waitForELBHealth(networkingClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for monitor %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceELBHealthRefreshFunc(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("Error: monitor %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("Error waiting for monitor %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceELBHealthRefreshFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		health, err := healthcheck.Get(networkingClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		// The health resource has no Status attribute, so a successful Get is the best we can do
		return health, "ACTIVE", nil
	}
}
