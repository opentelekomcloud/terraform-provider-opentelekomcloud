package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/deh/v1/hosts"
)

func TestAccOTCDedicatedHostV1_basic(t *testing.T) {
	var host hosts.Host

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCDeHV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDeHV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCDeHV1Exists("opentelekomcloud_deh_host_v1.deh1", &host),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "name", "test-deh-1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "auto_placement", "off"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "host_type", "h1"),
				),
			},
		},
	})
}

func TestAccOTCDedicatedHostV1_update(t *testing.T) {
	var host hosts.Host
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCDeHV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDeHV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCDeHV1Exists("opentelekomcloud_deh_host_v1.deh1", &host),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "name", "test-deh-1"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "auto_placement", "off"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "host_type", "h1"),
				),
			},
			resource.TestStep{
				Config: testAccDeHV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCDeHV1Exists("opentelekomcloud_deh_host_v1.deh1", &host),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "name", "test-deh-2"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "auto_placement", "on"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_deh_host_v1.deh1", "host_type", "h1"),
				),
			},
		},
	})
}

func TestAccOTCDedicatedHostV1_timeout(t *testing.T) {
	var host hosts.Host

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCDeHV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDeHV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCDeHV1Exists("opentelekomcloud_deh_host_v1.deh1", &host),
				),
			},
		},
	})
}

func testAccCheckOTCDeHV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dehClient, err := config.dehV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud deh client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_deh_host_v1" {
			continue
		}

		_, err := hosts.Get(dehClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Dedicated Host still exists")
		}
	}

	return nil
}

func testAccCheckOTCDeHV1Exists(n string, host *hosts.Host) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dehClient, err := config.dehV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud DeH client: %s", err)
		}

		found, err := hosts.Get(dehClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("DeH file not found")
		}

		*host = *found

		return nil
	}
}

var testAccDeHV1_basic = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
	 availability_zone= "%s"     
     auto_placement= "off"
     host_type= "h1"   
	name = "test-deh-1"
}
`, OS_AVAILABILITY_ZONE)

var testAccDeHV1_update = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
	 availability_zone= "%s"     
     auto_placement= "on"
     host_type= "h1"
	name = "test-deh-2"
}
`, OS_AVAILABILITY_ZONE)

var testAccDeHV1_timeout = fmt.Sprintf(`
resource "opentelekomcloud_deh_host_v1" "deh1" {
	 availability_zone= "%s"     
     auto_placement= "off"
     host_type= "h1"
	name = "test-deh-1"
  timeouts {
    create = "5m"
    delete = "5m"
  }
}`, OS_AVAILABILITY_ZONE)
