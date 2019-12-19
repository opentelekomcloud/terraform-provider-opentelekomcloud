package loggroups

import "github.com/huaweicloud/golangsdk"

const rootPath = "log-groups"

// createURL will build the url of creation
func createURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootPath)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(rootPath, id)
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL(rootPath, id)
}
