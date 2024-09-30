package dms

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
)

const (
	errCreationClient   = "error creating OpenTelekomCloud DMSv1 client: %w"
	errCreationClientV2 = "error creating OpenTelekomCloud DMSv2 client: %w"
	dmsClientV2         = "dms-v2-client"
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

// StrSliceContainsAnother checks whether a string slice (b) contains another string slice (s).
func StrSliceContainsAnother(b []string, s []string) bool {
	// The empty set is the subset of any set.
	if len(s) < 1 {
		return true
	}
	for _, v := range s {
		if !common.StrSliceContains(b, v) {
			return false
		}
	}
	return true
}
