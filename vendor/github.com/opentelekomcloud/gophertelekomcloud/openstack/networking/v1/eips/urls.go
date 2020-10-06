package eips

import "github.com/opentelekomcloud/gophertelekomcloud"

const resourcePath = "publicips"

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(client.ProjectID, resourcePath)
}

func resourceURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(client.ProjectID, resourcePath, id)
}
