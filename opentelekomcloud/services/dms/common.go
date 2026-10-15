package dms

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
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

// convertRocketMQNotFound converts a RocketMQ 400 error with one of the given error codes into a 404 error,
// as the RocketMQ API reports missing sub-resources this way.
func convertRocketMQNotFound(err error, errCodes ...string) error {
	var err400 golangsdk.ErrDefault400
	if !errors.As(err, &err400) {
		return err
	}
	var apiErr struct {
		ErrorCode string `json:"error_code"`
	}
	if json.Unmarshal(err400.Body, &apiErr) == nil && slices.Contains(errCodes, apiErr.ErrorCode) {
		return golangsdk.ErrDefault404{}
	}
	return err
}
