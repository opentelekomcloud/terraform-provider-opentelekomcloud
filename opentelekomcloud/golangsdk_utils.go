package opentelekomcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func chooseECSV1Client(d *schema.ResourceData, config *Config) (*golangsdk.ServiceClient, error) {
	return config.loadECSV1Client(GetRegion(d, config))
}

func chooseCESClient(d *schema.ResourceData, config *Config) (*golangsdk.ServiceClient, error) {
	return config.loadCESClient(GetRegion(d, config))
}
