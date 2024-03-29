package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/identity/v3/policies"
	acc "github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const (
	resourceRoleName = "opentelekomcloud_identity_role_v3.role"
)

func TestAccIdentityRoleV3_basic(t *testing.T) {
	var iamRole policies.Policy

	roleName := "custom-role" + acctest.RandString(10)
	roleName2 := "custom-role" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			acc.TestAccPreCheckAdminOnly(t)
		},

		ProviderFactories: acc.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityRoleV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleV3_basic(roleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleV3Exists(resourceRoleName, &iamRole),
					resource.TestCheckResourceAttr(resourceRoleName, "description", "role"),
					resource.TestCheckResourceAttr(resourceRoleName, "display_name", roleName),
					resource.TestCheckResourceAttr(resourceRoleName, "display_layer", "domain"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.effect", "Allow"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.0", "obs:bucket:GetBucketAcl"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.resource.0", "obs:*:*:bucket:*"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.1.effect", "Allow"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.1.action.0", "obs:bucket:HeadBucket"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.1.action.1", "obs:bucket:ListBucketMultipartUploads"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.1.action.2", "obs:bucket:ListBucket"),
				),
			},
			{
				Config: testAccIdentityRoleV3_update(roleName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleV3Exists(resourceRoleName, &iamRole),
					resource.TestCheckResourceAttr(resourceRoleName, "description", "updated role#1"),
					resource.TestCheckResourceAttr(resourceRoleName, "display_name", roleName2),
					resource.TestCheckResourceAttr(resourceRoleName, "display_layer", "project"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.effect", "Deny"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.0", "evs:volumeTags:list"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.1", "evs:transfers:list"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.2", "evs:snapshots:list"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.3", "evs:volumes:list"),
				),
			},
			{
				Config: testAccIdentityRoleV3_update_2(roleName2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityRoleV3Exists(resourceRoleName, &iamRole),
					resource.TestCheckResourceAttr(resourceRoleName, "description", "updated role#1"),
					resource.TestCheckResourceAttr(resourceRoleName, "display_name", roleName2),
					resource.TestCheckResourceAttr(resourceRoleName, "display_layer", "domain"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.effect", "Allow"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.0", "obs:bucket:ListBucketVersions"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.1", "obs:bucket:GetBucketAcl"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.2", "obs:bucket:GetBucketNotification"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.action.3", "obs:bucket:GetBucketWebsite"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.resource.0", "OBS:*:*:bucket:test-bucket"),
					resource.TestCheckResourceAttr(resourceRoleName, "statement.0.resource.1", "OBS:*:*:object:your_object"),
				),
			},
		},
	})
}

func TestAccIdentityV3Role_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.TestAccPreCheck(t)
			acc.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acc.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIdentityV3UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleV3_import(acctest.RandString(10)),
			},

			{
				ResourceName:      resourceRoleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIdentityRoleV3_basic(val string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_role_v3" "role" {
  description   = "role"
  display_name  = "%s"
  display_layer = "domain"
  statement {
    effect    = "Allow"
    action    = ["obs:bucket:GetBucketAcl"]
    resource  = ["obs:*:*:bucket:*"]
    condition = <<EOF
	    {
	      "StringStartWith": {
	          "g:ProjectName": [
	              "eu-de"
	          ]
	      },
	      "StringNotEqualsIgnoreCase": {
	          "g:ServiceName": [
	              "iam"
	          ]
	    }
  }
EOF
  }
  statement {
    effect = "Allow"
    action = [
      "obs:bucket:HeadBucket",
      "obs:bucket:ListBucketMultipartUploads",
      "obs:bucket:ListBucket"
    ]
  }
}`, val)
}

func testAccIdentityRoleV3_update(val string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_role_v3" "role" {
  description   = "updated role#1"
  display_name  = "%s"
  display_layer = "project"
  statement {
    effect = "Deny"
    action = ["evs:volumeTags:list",
      "evs:transfers:list",
      "evs:snapshots:list",
      "evs:volumes:list"
    ]
  }
}`, val)
}

func testAccIdentityRoleV3_update_2(val string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_role_v3" "role" {
  description   = "updated role#1"
  display_name  = "%s"
  display_layer = "domain"
  statement {
    effect = "Allow"
    action = ["obs:bucket:ListBucketVersions",
      "obs:bucket:GetBucketAcl",
      "obs:bucket:GetBucketNotification",
      "obs:bucket:GetBucketWebsite"
    ]
    resource = ["OBS:*:*:bucket:test-bucket",
      "OBS:*:*:object:your_object"
    ]
    condition = <<EOF
	    {
	      "StringStartWith": {
	          "g:ProjectName": [
	              "eu-de"
	          ]
	      }
  }
EOF
  }
}`, val)
}

func testAccCheckIdentityRoleV3Destroy(s *terraform.State) error {
	config := acc.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.IdentityV30Client()
	if err != nil {
		return fmt.Errorf("error creating sdk client, err=%s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_identity_role_v3" {
			continue
		}

		_, err = policies.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IAM role still exists")
		}
	}

	return nil
}

func testAccCheckIdentityRoleV3Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := acc.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.IdentityV30Client()
		if err != nil {
			return fmt.Errorf("error creating sdk client, err=%s", err)
		}

		found, err := policies.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		policy = found

		return nil
	}
}

func testAccIdentityRoleV3_import(val string) string {
	return fmt.Sprintf(`
resource "opentelekomcloud_identity_role_v3" "role" {
  description   = "role"
  display_name  = "custom_role%s"
  display_layer = "domain"
  statement {
    effect = "Allow"
    action = ["ecs:*:list*"]
  }
}`, val)
}
