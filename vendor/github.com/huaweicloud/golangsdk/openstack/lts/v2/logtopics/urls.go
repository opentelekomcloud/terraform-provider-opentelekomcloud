package logtopics

import "github.com/huaweicloud/golangsdk"

const (
	resourcePath = "log-topics"
	rootPath     = "log-groups"
)

// createURL will build the url of creation
func createURL(client *golangsdk.ServiceClient, groupId string) string {
	return client.ServiceURL(rootPath, groupId, resourcePath)
}

// deleteURL will build the url of deletion
func deleteURL(client *golangsdk.ServiceClient, groupId string, id string) string {
	return client.ServiceURL(rootPath, groupId, resourcePath, id)
}

// getURL will build the get url of get function
func getURL(client *golangsdk.ServiceClient, groupId string, id string) string {
	return client.ServiceURL(rootPath, groupId, resourcePath, id)
}
