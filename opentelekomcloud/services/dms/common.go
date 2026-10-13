package dms

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	errCreationClient   = "error creating OpenTelekomCloud DMSv1 client: %w"
	errCreationClientV2 = "error creating OpenTelekomCloud DMSv2 client: %w"
	dmsClientV2         = "dms-v2-client"
	dmsV1Deprecated     = "This resource API is deprecated. Please use opentelekomcloud_dms_topic_v2 instead."
)

func MarshalValue(i interface{}) string {
	if i == nil {
		return ""
	}

	jsonRaw, err := json.Marshal(i)
	if err != nil {
		log.Printf("[WARN] failed to marshal %#v: %s", i, err)
		return ""
	}

	return strings.Trim(string(jsonRaw), `"`)
}

// parseRocketMQSubResourceID splits ID of a RocketMQ topic or consumer group in format <instance_id>/<name>.
func parseRocketMQSubResourceID(id string) (string, string, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid ID format, want '<instance_id>/<name>', but '%s'", id)
	}
	return parts[0], parts[1], nil
}
