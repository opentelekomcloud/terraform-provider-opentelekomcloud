package cfw

import (
	"fmt"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	cfwjob "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v3/job"
)

const (
	errCreationV1Client = "error creating OpenTelekomCloud CFWv1 client: %w"
	errCreationV2Client = "error creating OpenTelekomCloud CFWv2 client: %w"
	errCreationV3Client = "error creating OpenTelekomCloud CFWv3 client: %w"
	keyClientV1         = "cfw-v1-client"
	keyClientV2         = "cfw-v2-client"
	keyClientV3         = "cfw-v3-client"
)

func WaitForJobCompleted(client *golangsdk.ServiceClient, waitTime int, interval time.Duration, jobID string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	return golangsdk.WaitFor(waitTime, func() (bool, error) {
		job, err := cfwjob.Get(client, jobID)
		if err != nil {
			return false, err
		}

		if job.Status == "Success" {
			return true, nil
		}
		if job.Status == "Failed" {
			err = fmt.Errorf("job %s failed", job.Id)
			return false, err
		}

		time.Sleep(interval * time.Second)
		return false, nil
	})
}
