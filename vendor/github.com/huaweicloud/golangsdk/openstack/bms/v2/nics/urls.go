package nics

import (
	"github.com/huaweicloud/golangsdk"
)

func listURL(client *golangsdk.ServiceClient, serverId string) string {
	return client.ServiceURL("servers", serverId, "os-interface")
}

func getURL(client *golangsdk.ServiceClient, serverId string, Id string) string {
	return client.ServiceURL("servers", serverId, "os-interface", Id)
}
