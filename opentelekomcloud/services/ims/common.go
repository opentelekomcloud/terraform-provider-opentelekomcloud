package ims

import (
	"fmt"
	"strings"
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
