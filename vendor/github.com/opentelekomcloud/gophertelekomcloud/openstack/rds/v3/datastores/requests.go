package datastores

import (
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/pagination"
)

func List(client *golangsdk.ServiceClient, databasesname string) pagination.Pager {
	url := listURL(client, databasesname)

	pageRdsList := pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return DataStoresPage{pagination.SinglePageBase(r)}
	})

	rdsheader := map[string]string{"Content-Type": "application/json"}
	pageRdsList.Headers = rdsheader
	return pageRdsList
}
