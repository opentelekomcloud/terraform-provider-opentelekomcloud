package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/tms/v1/tags"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccTmsTag_basic(t *testing.T) {
	resourceName := "opentelekomcloud_tms_tags_v1.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckTmsV1TagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTmsTag_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmsV1TagExists("foo", "bar"),
					testAccCheckTmsV1TagExists("k", "v"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
				),
			},
			{
				Config: testTmsTag_Updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTmsV1TagExists("foo", "bar"),
					testAccCheckTmsV1TagExists("n", "p"),
					testAccCheckTmsV1TagExists("one", "two"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "3"),
				),
			},
		},
	})
}

func testAccCheckTmsV1TagDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	tmsClient, err := config.TmsV1Client()
	if err != nil {
		return fmt.Errorf("error creating opentelekomcloud tms client: %s", err)
	}

	tmsTags := map[string]string{"foo": "bar", "k": "v"}

	response, err := tags.Get(tmsClient).Extract()
	if err != nil {
		return err
	}

	for _, value := range response.Tags {
		if _, ok := tmsTags[value.Key]; ok {
			return fmt.Errorf("opentelekomcloud_tms_tags %s/%s still exists", value.Key, value)
		}
	}
	return nil
}

func testAccCheckTmsV1TagExists(key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := common.TestAccProvider.Meta().(*cfg.Config)
		tmsClient, err := config.TmsV1Client()
		if err != nil {
			return fmt.Errorf("error creating opentelekomcloud tms client: %s", err)
		}

		response, err := tags.Get(tmsClient).Extract()
		if err != nil {
			return err
		}
		for _, tValue := range response.Tags {
			if tValue.Key == key && tValue.Value == value {
				return nil
			}
		}
		return fmt.Errorf("opentelekomcloud_tms_tags %s/%s does not exist", key, value)
	}
}

const testTmsTag_basic = `
resource "opentelekomcloud_tms_tags_v1" "test" {
  tags {
    key   = "foo"
    value = "bar"
  }
  tags {
    key   = "k"
    value = "v"
  }
}
`

const testTmsTag_Updated = `
resource "opentelekomcloud_tms_tags_v1" "test" {
  tags {
    key   = "foo"
    value = "bar"
  }
  tags {
    key   = "n"
    value = "p"
  }
  tags {
    key   = "one"
    value = "two"
  }
}
`
