package dcaas

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

func TestDirectConnectV2Datasource_basic(t *testing.T) {
	var directConnectName = fmt.Sprintf("dc-%s", acctest.RandString(5))
	const dc = "opentelekomcloud_direct_connect_v2.direct_connect"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectConnectV2Datasource_basic(directConnectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dc, "bandwidth", "100"),
				),
			},
		},
	})
}

func testAccDirectConnectV2Datasource_basic(directConnectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_direct_connect_v2" "direct_connect" {
  name          = "%s"
  port_type     = "1G"
  location      = "Biere"
  bandwidth     = 100
  provider_name = "OTC"
}

data "opentelekomcloud_direct_connect_v2" "direct_connect" {
  id = opentelekomcloud_direct_connect_v2.direct_connect.id
}
	`, directConnectName)
}

func TestDirectConnectV2Datasource_byName(t *testing.T) {
	var directConnectName = fmt.Sprintf("dc-%s", acctest.RandString(5))
	const dc = "data.opentelekomcloud_direct_connect_v2.by_name"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDirectConnectV2Datasource_byName(directConnectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dc, "name", directConnectName),
					resource.TestCheckResourceAttr(dc, "bandwidth", "100"),
					resource.TestCheckResourceAttrPair(dc, "id",
						"opentelekomcloud_direct_connect_v2.direct_connect", "id"),
				),
			},
		},
	})
}

func testAccDirectConnectV2Datasource_byName(directConnectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_direct_connect_v2" "direct_connect" {
  name          = "%s"
  port_type     = "1G"
  location      = "Biere"
  bandwidth     = 100
  provider_name = "OTC"
}

data "opentelekomcloud_direct_connect_v2" "by_name" {
  name = opentelekomcloud_direct_connect_v2.direct_connect.name
}
		`, directConnectName)
}

func TestDirectConnectV2Datasource_byNameMultiple(t *testing.T) {
	var directConnectName = fmt.Sprintf("dc-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDirectConnectV2Destroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDirectConnectV2Datasource_byNameMultiple(directConnectName),
				ExpectError: regexp.MustCompile(`your query returned more than one result: 2 direct connects match name .* \(IDs: [^)]+\)\. Please query by .id.`),
			},
		},
	})
}

func testAccDirectConnectV2Datasource_byNameMultiple(directConnectName string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_direct_connect_v2" "dc1" {
  name          = "%[1]s"
  port_type     = "1G"
  location      = "Biere"
  bandwidth     = 100
  provider_name = "OTC"
}

resource "opentelekomcloud_direct_connect_v2" "dc2" {
  name          = "%[1]s"
  port_type     = "1G"
  location      = "Biere"
  bandwidth     = 100
  provider_name = "OTC"
}

data "opentelekomcloud_direct_connect_v2" "by_name" {
  name = "%[1]s"
  depends_on = [
    opentelekomcloud_direct_connect_v2.dc1,
    opentelekomcloud_direct_connect_v2.dc2,
  ]
}
			`, directConnectName)
}

func TestDirectConnectV2Datasource_nameNotFound(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDirectConnectV2Datasource_nameNotFound(),
				ExpectError: regexp.MustCompile(`your query returned no results\. Please change your search criteria`),
			},
		},
	})
}

func testAccDirectConnectV2Datasource_nameNotFound() string {
	return fmt.Sprintf(`
data "opentelekomcloud_direct_connect_v2" "by_name" {
  name = "nonexistent-dc-%s"
}
	`, acctest.RandString(8))
}

func TestDirectConnectV2Datasource_idAndName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "opentelekomcloud_direct_connect_v2" "test" {
  id   = "any-value"
  name = "any-value"
}`,
				ExpectError: regexp.MustCompile("only one of `id` or `name` can be specified"),
			},
		},
	})
}

func TestDirectConnectV2Datasource_noFilter(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      `data "opentelekomcloud_direct_connect_v2" "test" {}`,
				ExpectError: regexp.MustCompile("one of `id` or `name` must be specified"),
			},
		},
	})
}
