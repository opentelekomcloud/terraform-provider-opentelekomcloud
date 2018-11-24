package resetstate

import (
	"github.com/huaweicloud/golangsdk"
)

func actionURL(client *golangsdk.ServiceClient, id string) string {
	return client.ServiceURL("servers", id, "action")
}
