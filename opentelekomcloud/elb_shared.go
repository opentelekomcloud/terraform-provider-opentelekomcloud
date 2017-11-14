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

const loadbalancerActiveTimeoutSeconds = 300
const loadbalancerDeleteTimeoutSeconds = 300

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
	fmt.Printf("[DEBUG] Waiting for loadbalancer %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceELBLoadBalancerRefreshFunc(networkingClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	fmt.Printf("[DEBUG] waitForELBLoadBalancer before WaitForState \n")
	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "DELETED":
				fmt.Printf("[DEBUG] waitForELBLoadBalancer DELETED\n")
				return nil
			default:
				fmt.Printf("[DEBUG] waitForELBLoadBalancer default \n")
				return fmt.Errorf("Error: loadbalancer %s not found: %s", id, err)
			}
		}
		fmt.Printf("Error waiting for loadbalancer %s to become %s: %s", id, target, err)
		return fmt.Errorf("Error waiting for loadbalancer %s to become %s: %s", id, target, err)
	}
	fmt.Printf("waiting for loadbalancer return nil")

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

func WaitForJobSuccess(client *gophercloud.ServiceClient, uri string, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		job := new(loadbalancer_elbs.JobStatus)
		_, err := client.Get("https://elb.eu-de.otc.t-systems.com"+uri, &job, nil)
		if err != nil {
			fmt.Printf("WaitForJobSuccess: err /%+v.\n", err)
			return false, err
		}
		fmt.Printf("JobStatus: %+v. uri=%s \n", job, uri)

		if job.Status == "SUCCESS" {
			fmt.Printf("JobStatus: SUCCESS\n")
			return true, nil
		}
		if job.Status == "FAIL" {
			err = fmt.Errorf("Job failed with code %s: %s.\n", job.ErrorCode, job.FailReason)
			fmt.Printf("JobStatus: Job failed with code %s: %s.\n", job.ErrorCode, job.FailReason)
			return false, err
		}

		return false, nil
	})
}

func GetJobEntity(client *gophercloud.ServiceClient, uri string, label string) (map[string]interface{}, error) {
	job := new(loadbalancer_elbs.JobStatus)
	_, err := client.Get("https://elb.eu-de.otc.t-systems.com"+uri, &job, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("JobStatus: %+v.\n", job)

	if job.Status == "SUCCESS" {
		if e := job.Entities[label]; e != nil {
			if m, ok := e.(map[string]interface{}); ok {
				return m, nil
			}
		}
	}

	return nil, nil
}

// WaitForLoadBalancerState will wait until a loadbalancer reaches a given state.
func WaitForLoadBalancerState(client *gophercloud.ServiceClient, lbID string, status int, secs int) error {
	return gophercloud.WaitFor(secs, func() (bool, error) {
		current, err := loadbalancer_elbs.Get(client, lbID).Extract()
		if err != nil {
			if httpStatus, ok := err.(gophercloud.ErrDefault404); ok {
				if httpStatus.Actual == 404 {
					//if status == "DELETED" {
					//	return true, nil
					//}
				}
			}
			return false, err
		}

		if current.AdminStateUp == status {
			return true, nil
		}

		return false, nil
	})
}
