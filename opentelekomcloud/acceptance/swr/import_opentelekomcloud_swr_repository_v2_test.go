package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestSwrRepositoryV2_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testSwrRepositoryV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testSwrRepositoryV2Basic,
			},
			{
				ResourceName:      resourceRepoName1,
				ImportStateId:     fmt.Sprintf("%[1]s/%[1]s", name),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
