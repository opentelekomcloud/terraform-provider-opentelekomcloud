package ims

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/image/v2/members"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/ims/v2/images"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	errCreationClient = "error creating OpenTelekomCloud IMSv2 client: %w"
)

func GetImageByName(d *schema.ResourceData, cfg *cfg.Config, name string) (string, error) {
	client, err := cfg.ImageV2Client(cfg.GetRegion(d))
	if err != nil {
		return "", fmt.Errorf("error creating IMSv2 client: %w", err)
	}

	img, err := images.ListImages(client, images.ListImagesOpts{
		Name:  d.Get("image_name").(string),
		Limit: 1,
	})
	if err != nil {
		return "", fmt.Errorf("error extracting images: %w", err)
	}

	if len(img) < 1 {
		return "", fmt.Errorf("no image matching name: %s", name)
	}

	return img[0].Id, nil
}

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
		n, err := members.Get(client, members.MemberOpts{
			ImageId:  imageID,
			MemberId: memberID,
		})

		if err != nil {
			return nil, "", err
		}

		if status == n.Status {
			return n, n.Status, nil
		}

		return n, n.Status, nil
	}
}
