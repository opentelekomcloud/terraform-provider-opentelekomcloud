package rds

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmespath/go-jmespath"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
)

const (
	errCreateClient = "error creating RDSv3 client: %w"
	keyClientV3     = "rdsv3-client"
)

var (
	retryErrCodes = map[string]struct{}{
		"DBS.201202": {},
		"DBS.200011": {},
		"DBS.200018": {},
		"DBS.200019": {},
		"DBS.200047": {},
		"DBS.200611": {},
		"DBS.200080": {},
		"DBS.200463": {}, // create replica instance
		"DBS.201015": {},
		"DBS.201206": {},
		"DBS.212033": {}, // http response code is 403
		"DBS.280011": {},
		"DBS.280343": {},
		"DBS.280816": {},
	}
)

func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok && errCode.Actual == 409 {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	if errCode, ok := err.(golangsdk.ErrDefault403); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("errCode||error_code", apiError)
		if errorCodeErr != nil {
			return false, fmt.Errorf("error parse errorCode from response body: %s", errorCodeErr)
		}

		if _, ok = retryErrCodes[errorCode.(string)]; ok {
			// The operation failed to execute and needs to be executed again, because other operations are
			// currently in progress.
			return true, err
		}
	}
	// Operation execution failed due to some resource or server issues, no need to try again.
	return false, err
}

func checkMinorVersion(dbInfo map[string]interface{}) bool {
	// returns true if version is not minor
	version, ok := dbInfo["version"].(string)
	if !ok {
		return true
	}
	parts := strings.SplitN(version, ".", 3)

	return len(parts) <= 2
}
