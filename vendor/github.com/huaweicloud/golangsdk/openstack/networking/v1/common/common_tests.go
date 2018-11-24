package common

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/testhelper/client"
)

const TokenID = client.TokenID

// Fake project id to use.
const ProjectID = "85636478b0bd8e67e89469c7749d4127"

func ServiceClient() *golangsdk.ServiceClient {
	sc := client.ServiceClient()
	sc.ResourceBase = sc.Endpoint + "v1/"
	sc.ProjectID = ProjectID
	return sc
}
