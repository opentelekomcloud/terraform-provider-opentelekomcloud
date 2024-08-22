package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/deh/v1/hosts"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const resourceHostName = "opentelekomcloud_deh_host_v1.deh1"

func TestAccDedicatedHostV1_basic(t *testing.T) {
	var host hosts.Host

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDeHV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeHV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeHV1Exists(resourceHostName, &host),
					resource.TestCheckResourceAttr(resourceHostName, "name", "test-deh-1"),
					resource.TestCheckResourceAttr(resourceHostName, "auto_placement", "off"),
					resource.TestCheckResourceAttr(resourceHostName, "host_type", "s3"),
					resource.TestCheckResourceAttr(resourceHostName, "tags.created_by", "terraform"),
					resource.TestCheckResourceAttr(resourceHostName, "tags.muh", "value-create"),
				),
			},
			{
				Config: testAccDeHV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeHV1Exists(resourceHostName, &host),
					resource.TestCheckResourceAttr(resourceHostName, "name", "test-deh-2"),
					resource.TestCheckResourceAttr(resourceHostName, "auto_placement", "on"),
					resource.TestCheckResourceAttr(resourceHostName, "host_type", "s3"),
					resource.TestCheckResourceAttr(resourceHostName, "tags.updated_by", "terraform"),
					resource.TestCheckResourceAttr(resourceHostName, "tags.muh", "value-update"),
				),
			},
		},
	})
}

func TestAccDedicatedHostV1_timeout(t *testing.T) {
	var host hosts.Host

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDeHV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeHV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeHV1Exists("opentelekomcloud_deh_host_v1.deh1", &host),
				),
			},
		},
	})
}

func testAccCheckDeHV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DehV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud deh client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_deh_host_v1" {
			continue
		}

		_, err := hosts.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("dedicated Host still exists")
		}
	}

	return nil
}

func testAccCheckDeHV1Exists(n string, host *hosts.Host) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DehV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud DeH client: %s", err)
		}

		found, err := hosts.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("deH file not found")
		}

		*host = *found

		return nil
	}
}

var testAccDeHV1Basic = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
  availability_zone = "%s"
  auto_placement    = "off"
  host_type         = "s3"
  name              = "test-deh-1"

  tags = {
    created_by = "terraform"
    muh = "value-create"
  }
}
`, env.OS_AVAILABILITY_ZONE)

var testAccDeHV1Update = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
  availability_zone = "%s"
  auto_placement    = "on"
  host_type         = "s3"
  name              = "test-deh-2"

  tags = {
    updated_by = "terraform"
    muh = "value-update"
  }
}
`, env.OS_AVAILABILITY_ZONE)

var testAccDeHV1Timeout = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
  availability_zone = "%s"
  auto_placement    = "off"
  host_type         = "h1"
  name              = "test-deh-1"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}`, env.OS_AVAILABILITY_ZONE)
