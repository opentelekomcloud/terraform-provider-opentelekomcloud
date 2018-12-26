package services

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// List makes a request against the API to list services.
func List(client *golangsdk.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, listURL(client), func(r pagination.PageResult) pagination.Page {
		return ServicePage{pagination.SinglePageBase(r)}
	})
}
