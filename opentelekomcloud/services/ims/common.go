package ims

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/members"
)

const (
	errCreationClient = "error creating OpenTelekomCloud IMSv2 client: %w"
)

func ResourceImagesImageAccessV2ParseID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("unable to determine image share access ID")
	}

	imageID := idParts[0]
	memberID := idParts[1]

	return imageID, memberID, nil
}

func waitForImageRequestStatus(client *golangsdk.ServiceClient, imageID, memberID, status string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := members.Get(client, imageID, memberID).Extract()
		if err != nil {
			return nil, "", err
		}
		if status == n.Status {
			return n, n.Status, nil
		}

		return n, n.Status, nil
	}
}
