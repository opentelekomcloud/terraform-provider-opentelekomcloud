package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/compute/v2/extensions/secgroups"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccComputeV2SecGroup_basic(t *testing.T) {
	var secGroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupBasicOrig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secGroup),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_update(t *testing.T) {
	var secGroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupBasicOrig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secGroup),
				),
			},
			{
				Config: testAccComputeV2SecGroupBasicUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secGroup),
					testAccCheckComputeV2SecGroupRuleCount(&secGroup, 2),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_groupID(t *testing.T) {
	var secgroup1, secGroup2, secgroup3 secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupGroupIDOrig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_2", &secGroup2),
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_3", &secgroup3),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secgroup1, &secgroup3),
				),
			},
			{
				Config: testAccComputeV2SecGroupGroupIDUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secgroup1),
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_2", &secGroup2),
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_3", &secgroup3),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secGroup2, &secgroup3),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_self(t *testing.T) {
	var secGroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupSelf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secGroup),
					testAccCheckComputeV2SecGroupGroupIDMatch(&secGroup, &secGroup),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_secgroup_v2.sg_1", "rule.3170486100.self", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_compute_secgroup_v2.sg_1", "rule.3170486100.from_group_id", ""),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_icmpZero(t *testing.T) {
	var secGroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupIcmpZero,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secGroup),
				),
			},
		},
	})
}

func TestAccComputeV2SecGroup_timeout(t *testing.T) {
	var secGroup secgroups.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeV2SecGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2SecGroupTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2SecGroupExists("opentelekomcloud_compute_secgroup_v2.sg_1", &secGroup),
				),
			},
		},
	})
}

func testAccCheckComputeV2SecGroupDestroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	computeClient, err := config.ComputeV2Client(env.OsRegionName)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_compute_secgroup_v2" {
			continue
		}

		_, err := secgroups.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("security group still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2SecGroupExists(n string, secgroup *secgroups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		computeClient, err := config.ComputeV2Client(env.OsRegionName)
		if err != nil {
			return fmt.Errorf("error creating OpenTelekomCloud compute client: %s", err)
		}

		found, err := secgroups.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("security group not found")
		}

		*secgroup = *found

		return nil
	}
}

func testAccCheckComputeV2SecGroupRuleCount(secgroup *secgroups.SecurityGroup, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(secgroup.Rules) != count {
			return fmt.Errorf("security group rule count does not match. Expected %d, got %d", count, len(secgroup.Rules))
		}

		return nil
	}
}

func testAccCheckComputeV2SecGroupGroupIDMatch(sg1, sg2 *secgroups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg2.Rules) == 1 {
			if sg1.Name != sg2.Rules[0].Group.Name || sg1.TenantID != sg2.Rules[0].Group.TenantID {
				return fmt.Errorf("%s was not correctly applied to %s", sg1.Name, sg2.Name)
			}
		} else {
			return fmt.Errorf("%s rule count is incorrect", sg2.Name)
		}

		return nil
	}
}

const (
	testAccComputeV2SecGroupBasicOrig = `
resource "opentelekomcloud_compute_secgroup_v2" "sg_1" {
  name        = "sg_1"
  description = "first test security group"
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
  rule {
    from_port   = 1
    to_port     = 65535
    ip_protocol = "udp"
    cidr        = "0.0.0.0/0"
  }
  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}
`

	testAccComputeV2SecGroupBasicUpdate = `
resource "opentelekomcloud_compute_secgroup_v2" "sg_1" {
  name        = "sg_1"
  description = "first test security group"
  rule {
    from_port   = 2200
    to_port     = 2200
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}
`

	testAccComputeV2SecGroupGroupIDOrig = `
resource "opentelekomcloud_compute_secgroup_v2" "sg_1" {
  name        = "sg_1"
  description = "first test security group"
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_secgroup_v2" "sg_2" {
  name        = "sg_2"
  description = "second test security group"
  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_secgroup_v2" "sg_3" {
  name        = "sg_3"
  description = "third test security group"
  rule {
    from_port     = 80
    to_port       = 80
    ip_protocol   = "tcp"
    from_group_id = opentelekomcloud_compute_secgroup_v2.sg_1.id
  }
}
`

	testAccComputeV2SecGroupGroupIDUpdate = `
resource "opentelekomcloud_compute_secgroup_v2" "sg_1" {
  name        = "sg_1"
  description = "first test security group"
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_secgroup_v2" "sg_2" {
  name        = "sg_2"
  description = "second test security group"
  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}

resource "opentelekomcloud_compute_secgroup_v2" "sg_3" {
  name        = "sg_3"
  description = "third test security group"
  rule {
    from_port     = 80
    to_port       = 80
    ip_protocol   = "tcp"
    from_group_id = opentelekomcloud_compute_secgroup_v2.sg_2.id
  }
}
`

	testAccComputeV2SecGroupSelf = `
resource "opentelekomcloud_compute_secgroup_v2" "sg_1" {
  name        = "sg_1"
  description = "first test security group"
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    self        = true
  }
}
`

	testAccComputeV2SecGroupIcmpZero = `
resource "opentelekomcloud_compute_secgroup_v2" "sg_1" {
  name        = "sg_1"
  description = "first test security group"
  rule {
    from_port   = 0
    to_port     = 0
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}
`

	testAccComputeV2SecGroupTimeout = `
resource "opentelekomcloud_compute_secgroup_v2" "sg_1" {
  name        = "sg_1"
  description = "first test security group"
  rule {
    from_port   = 0
    to_port     = 0
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }

  timeouts {
    delete = "5m"
  }
}
`
)
