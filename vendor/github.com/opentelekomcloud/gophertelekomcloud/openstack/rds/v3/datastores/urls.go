package datastores

import "github.com/opentelekomcloud/gophertelekomcloud"

func listURL(sc *golangsdk.ServiceClient, databasename string) string {
	return sc.ServiceURL("datastores", databasename)
}
