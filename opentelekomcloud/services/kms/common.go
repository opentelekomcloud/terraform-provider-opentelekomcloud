package kms

import (
	"fmt"
	"strings"
)

const (
	errCreationClient = "error creating OpenTelekomCloud KMSv1 client: %w"
	keyClientV1       = "kms-v1-client"
)

func ResourceKMSGrantV1ParseID(componentID string) (string, string, error) {
	parts := strings.Split(componentID, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unable to determine KMS Grant ID")
	}

	kmsID := parts[0]
	grantID := parts[1]

	return kmsID, grantID, nil
}
