package asm

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/asm/v1/servicemesh"
)

const (
	errCreationV1Client = "error creating OpenTelekomCloud ASMv1 client: %w"
	keyClientV1         = "asm-v1-client"
)

func instanceStateRefreshFunc(client *golangsdk.ServiceClient, meshId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		meshList, err := servicemesh.List(client)
		if err != nil {
			return nil, "Error retrieving ASM v1 service mesh", err
		}
		if len(meshList) == 0 {
			return nil, "DELETED", nil
		}
		for _, mesh := range meshList {
			if mesh.Metadata.UID == meshId {
				return mesh, mesh.Status.Phase, nil
			}
		}
		return nil, "DELETED", nil
	}
}

func WaitForDeleteServiceMesh(client *golangsdk.ServiceClient, waitTime int, interval time.Duration, meshId string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	return golangsdk.WaitFor(waitTime, func() (bool, error) {
		meshList, err := servicemesh.List(client)
		if err != nil {
			return false, err
		}
		found := false

		for _, mesh := range meshList {
			if mesh.Metadata.UID == meshId {
				found = true
			}
		}
		if !found {
			return true, nil
		}

		time.Sleep(interval * time.Second)
		return false, nil
	})
}
