package credentials

import (
	"strings"

	"github.com/opentelekomcloud/gophertelekomcloud"
)

const (
	rootPath        = "OS-CREDENTIAL"
	credentialsPath = "credentials"
)

func broken30Url(serviceUrl string) string {
	return strings.Replace(serviceUrl, "v3/", "v3.0/", 1)
}

func listURL(client *golangsdk.ServiceClient) string {
	return broken30Url(client.ServiceURL(rootPath, credentialsPath))
}

func getURL(client *golangsdk.ServiceClient, credID string) string {
	return broken30Url(client.ServiceURL(rootPath, credentialsPath, credID))
}

func createURL(client *golangsdk.ServiceClient) string {
	return broken30Url(client.ServiceURL(rootPath, credentialsPath))
}

func updateURL(client *golangsdk.ServiceClient, credID string) string {
	return broken30Url(client.ServiceURL(rootPath, credentialsPath, credID))
}

func deleteURL(client *golangsdk.ServiceClient, credID string) string {
	return broken30Url(client.ServiceURL(rootPath, credentialsPath, credID))
}
