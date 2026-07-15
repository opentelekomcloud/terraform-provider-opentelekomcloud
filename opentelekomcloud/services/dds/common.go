package dds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmespath/go-jmespath"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/dds/v3/job"
)

const (
	errCreationV3Client = "error creating OpenTelekomCloud DDSv3 client: %w"
	keyClientV3         = "dds-v3-client"
)

var (
	retryErrCodes = map[string]struct{}{
		"DBS.200019":   {}, // An operation that conflicts with the current operation is in progress.
		"DBS.201014":   {},
		"DBS.201015":   {},
		"DBS.239037":   {},
		"DBS.201000":   {}, // ssl
		"DBS.00010009": {}, // Instance's status is not available for this operation.
	}
)

// The DDS instance is limited to only one operation at a time.
// In addition to locking and waiting between multiple operations, a retry method is required to ensure that the
// request can be executed correctly.
func handleMultiOperationsError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}
	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, fmt.Errorf("unmarshal the response body failed: %s", jsonErr)
		}

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
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

		errorCode, errorCodeErr := jmespath.Search("error_code", apiError)
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

func waitForJobCompleted(client *golangsdk.ServiceClient, secs int, jobID string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	return golangsdk.WaitFor(secs, func() (bool, error) {
		job, err := job.Get(client, jobID)
		if err != nil {
			return false, err
		}

		if job.Status == "Completed" {
			return true, nil
		}
		if job.Status == "Failed" {
			err = fmt.Errorf("Job failed %s.\n", job.Status)
			return false, err
		}

		time.Sleep(5 * time.Second)
		return false, nil
	})
}
