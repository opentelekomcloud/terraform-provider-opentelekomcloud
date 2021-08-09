package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rts/v1/stacks"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

func TestAccOTCRTSStackV1_basic(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRTSStackV1Exists("opentelekomcloud_rts_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "name", "terraform_provider_stack"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "status", "CREATE_COMPLETE"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "disable_rollback", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "timeout_mins", "60"),
				),
			},
			{
				Config: testAccRTSStackV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRTSStackV1Exists("opentelekomcloud_rts_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "disable_rollback", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "timeout_mins", "50"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "status", "UPDATE_COMPLETE"),
				),
			},
		},
	})
}

func TestAccOTCRTSStackV1_timeout(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckOTCRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRTSStackV1Exists("opentelekomcloud_rts_stack_v1.stack_1", &stacks),
				),
			},
		},
	})
}

func testAccCheckOTCRTSStackV1Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	orchestrationClient, err := config.OrchestrationV1Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating RTS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rts_stack_v1" {
			continue
		}

		stack, err := stacks.Get(orchestrationClient, "terraform_provider_stack").Extract()

		if err == nil {
			if stack.Status != "DELETE_COMPLETE" {
				return fmt.Errorf("stack still exists")
			}
		}
	}

	return nil
}

func testAccCheckOTCRTSStackV1Exists(n string, stack *stacks.RetrievedStack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		orchestrationClient, err := config.OrchestrationV1Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating RTS Client : %s", err)
		}

		found, err := stacks.Get(orchestrationClient, "terraform_provider_stack").Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("stack not found")
		}

		*stack = *found

		return nil
	}
}

const testAccRTSStackV1_basic = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= true
  timeout_mins=60
  template_body = <<JSON
          {
    "outputs": {
      "str1": {
        "description": "The description of the nat server.",
        "value": {
          "get_resource": "random"
        }
      }
    },
    "heat_template_version": "2013-05-23",
    "description": "A HOT template that create a single server and boot from volume.",
    "parameters": {
      "key_name": {
        "type": "string",
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

}
`

const testAccRTSStackV1_update = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= false
  timeout_mins=50
  template_body = <<JSON
           {
    "outputs": {
      "str1": {
        "description": "The description of the nat server.",
        "value": {
          "get_resource": "random"
        }
      }
    },
    "heat_template_version": "2013-05-23",
    "description": "A HOT template that create a single server and boot from volume.",
    "parameters": {
      "key_name": {
        "type": "string",
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

}
`
const testAccRTSStackV1_timeout = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= true
  timeout_mins=60

  template_body = <<JSON
          {
    "outputs": {
      "str1": {
        "description": "The description of the nat server.",
        "value": {
          "get_resource": "random"
        }
      }
    },
    "heat_template_version": "2013-05-23",
    "description": "A HOT template that create a single server and boot from volume.",
    "parameters": {
      "key_name": {
        "type": "string",
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

  timeouts {
    create = "10m"
    delete = "10m"
  }
}
`
