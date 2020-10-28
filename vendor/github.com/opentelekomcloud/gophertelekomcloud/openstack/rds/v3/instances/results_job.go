package instances

import (
	"fmt"
	"github.com/opentelekomcloud/gophertelekomcloud"
)

type DeleteInstanceRdsResponse struct {
	JobId string `json:"job_id"`
}

type EnlargeVolumeResp struct {
	JobId string `json:"job_id"`
}

type RestartRdsResponse struct {
	JobId string `json:"job_id"`
}

type SingleToHaResponse struct {
	JobId string `json:"job_id"`
}

type ResizeFlavor struct {
	JobId string `json:"job_id"`
}

type CreateRds struct {
	Instance Instance `json:"instance"`
	JobId    string   `json:"job_id"`
	OrderId  string   `json:"order_id"`
}

func WaitForJobCompleted(client *golangsdk.ServiceClient, secs int, jobID string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	return golangsdk.WaitFor(secs, func() (bool, error) {
		job := new(golangsdk.RDSJobStatus)

		requestOpts := &golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}}
		_, err := jobClient.Get(jobURL(jobClient.ResourceBase, jobID), job, requestOpts)
		if err != nil {
			return false, err
		}

		if job.Job.Status == "Completed" {
			return true, nil
		}
		if job.Job.Status == "Failed" {
			err = fmt.Errorf("Job failed %s.\n", job.Job.Status)
			return false, err
		}

		return false, nil
	})
}

func WaitForStateAvailable(client *golangsdk.ServiceClient, secs int, instanceID string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	return golangsdk.WaitFor(secs, func() (bool, error) {
		job := new(golangsdk.JsonRDSInstanceStatus)

		requestOpts := &golangsdk.RequestOpts{MoreHeaders: map[string]string{"Content-Type": "application/json"}}
		_, err := jobClient.Get(detailsURL(jobClient.ResourceBase, instanceID), job, requestOpts)
		if err != nil {
			return false, err
		}

		if job.Instances[0].Status == "ACTIVE" {
			return true, nil
		}
		if job.Instances[0].Status == "FAILED" {
			err = fmt.Errorf("Job failed %s.\n", job.Instances[0].Status)
			return false, err
		}

		return false, nil
	})
}

func jobURL(endpoint string, jobID string) string {
	return fmt.Sprintf("%sjobs?id=%s", endpoint, jobID)
}

func detailsURL(endpoint string, instanceID string) string {
	return fmt.Sprintf("%sinstances?id=%s", endpoint, instanceID)
}
