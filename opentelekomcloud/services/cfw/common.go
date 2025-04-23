package cfw

import (
	"fmt"
	"time"

	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	cfwmanagementv1 "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v1/management"
	cfwjob "github.com/opentelekomcloud/gophertelekomcloud/openstack/cfw/v3/job"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
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

func WaitForDeleteFirewall(client *golangsdk.ServiceClient, waitTime int, interval time.Duration, firewallID string) error {
	jobClient := *client
	jobClient.ResourceBase = jobClient.Endpoint

	return golangsdk.WaitFor(waitTime, func() (bool, error) {
		firewallList, err := cfwmanagementv1.List(client, cfwmanagementv1.ListOpts{
			Limit:  1024,
			Offset: 0,
		})
		if err != nil {
			return false, err
		}
		found := false

		for _, fw := range firewallList {
			if fw.FwInstanceId == firewallID {
				found = true
			}
		}
		if !found {
			return true, nil
		}

		time.Sleep(interval * time.Second)
		return false, nil
	})
}

func InterfaceToIntPtr(i interface{}) *int {
	v, ok := i.(int)
	if !ok {
		panic(fmterr.Errorf(`interfaceToIntPtr: value is not of type int: %#v`, i))
	}
	return &v
}
