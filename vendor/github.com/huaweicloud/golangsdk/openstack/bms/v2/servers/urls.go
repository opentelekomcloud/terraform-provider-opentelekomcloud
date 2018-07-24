package servers

import (
	"github.com/huaweicloud/golangsdk"
)

func getURL(client *golangsdk.ServiceClient, server_id string) string {
	return client.ServiceURL("servers", server_id)
}

func listURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL("servers", "detail")
}
