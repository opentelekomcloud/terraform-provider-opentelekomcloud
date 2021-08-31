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

const stackResourceName = "opentelekomcloud_rts_stack_v1.stack_1"

func TestAccRTSStackV1_basic(t *testing.T) {
	var stack stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackV1Exists(stackResourceName, &stack),
					resource.TestCheckResourceAttr(stackResourceName, "name", "terraform_provider_stack"),
					resource.TestCheckResourceAttr(stackResourceName, "status", "CREATE_COMPLETE"),
					resource.TestCheckResourceAttr(stackResourceName, "disable_rollback", "true"),
					resource.TestCheckResourceAttr(stackResourceName, "timeout_mins", "60"),
				),
			},
			{
				Config: testAccRTSStackV1Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackV1Exists(stackResourceName, &stack),
					resource.TestCheckResourceAttr(stackResourceName, "disable_rollback", "false"),
					resource.TestCheckResourceAttr(stackResourceName, "timeout_mins", "50"),
					resource.TestCheckResourceAttr(stackResourceName, "status", "UPDATE_COMPLETE"),
				),
			},
		},
	})
}

func TestAccRTSStackV1_timeout(t *testing.T) {
	var stack stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRTSStackV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRTSStackV1Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRTSStackV1Exists(stackResourceName, &stack),
				),
			},
		},
	})
}

func testAccCheckRTSStackV1Destroy(s *terraform.State) error {
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

func testAccCheckRTSStackV1Exists(n string, stack *stacks.RetrievedStack) resource.TestCheckFunc {
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

const (
	testAccRTSStackV1Basic = `
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
	testAccRTSStackV1Update = `
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
	testAccRTSStackV1Timeout = `
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
)
