package taurusdb

import (
	"encoding/json"
	"fmt"

	"github.com/jmespath/go-jmespath"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

var (
	retryErrCodes = map[string]struct{}{
		"DBS.200019":   {},
		"DBS.201014":   {},
		"DBS.201015":   {},
		"DBS.200047":   {},
		"DBS.05000084": {},
	}
)

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			return true, err
		}
	}
	return false, err
}
