package ims

import (
	"fmt"
	"strings"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/imageservice/v2/members"
)

const (
	errCreationClient = "error creating OpenTelekomCloud IMSv2 client: %w"
)

func resourceImagesImageAccessV2ParseID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("unable to determine image share access ID")
	}

	imageID := idParts[0]
	memberID := idParts[1]

	return imageID, memberID, nil
}

func resourceImagesImageAccessV2DetectMemberID(client *golangsdk.ServiceClient, imageID string) (string, error) {
	allPages, err := members.List(client, imageID).AllPages()
	if err != nil {
		return "", fmt.Errorf("unable to list image members: %w", err)
	}
	allMembers, err := members.ExtractMembers(allPages)
	if err != nil {
		return "", fmt.Errorf("unable to extract image members: %w", err)
	}
	if len(allMembers) == 0 {
		return "", fmt.Errorf("no members found for the %q image", imageID)
	}
	if len(allMembers) > 1 {
		return "", fmt.Errorf("too many members found for the %q image, please specify the member_id explicitly", imageID)
	}
	return allMembers[0].MemberID, nil
}
