package dataarts

import (
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dataarts/v1.1/cluster"
)

const (
	errCreationV11Client = "error creating OpenTelekomCloud DataArtsV11 client: %w"
	keyClientV11         = "dataarts-v11-client"
)

func WaitForClusterState(client *golangsdk.ServiceClient, secs int, instanceID string, status string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	return golangsdk.WaitFor(secs, func() (bool, error) {
		resp, err := cluster.Get(client, instanceID)
		if err != nil {
			return false, err
		}

		if resp.Status == status {
			return true, nil
		}

		return false, nil
	})
}
