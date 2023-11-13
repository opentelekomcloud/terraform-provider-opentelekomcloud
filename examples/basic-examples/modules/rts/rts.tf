
resource "opentelekomcloud_rts_stack_v1" "mystack" {
  name             = "${var.rts_name}_new"
  disable_rollback = true
  timeout_mins     = 60
  parameters = {
    "network_id"        = var.subnet_id
    "instance_type"     = var.instance_type
    "image_id"          = var.image_id
    "availability_zone" = var.availability_zone
  }
  template_body = <<STACK
  {
    "heat_template_version": "2016-04-08",
   "description": "Simple template to deploy",
    "parameters": {
        "image_id": {
           "type": "string",
            "description": "Image to be used for compute instance",
            "label": "Image ID"
        },
        "network_id": {
            "type": "string",
            "description": "The Network to be used",
            "label": "Network UUID"
       },
     "instance_type": {
          "type": "string",
           "description": "Type of instance (Flavor) to be used",
            "label": "Instance Type"
       },
	   "availability_zone":{
	       "type": "string",
           "description": "AZ",
            "label": "AZ Name"
	   }
    },
    "resources": {
        "my_instance": {
            "type": "OS::Nova::Server",
            "properties": {
                "image": {
                    "get_param": "image_id"
                },
                "flavor": {
                    "get_param": "instance_type"
                },
                "networks": [{
                    "network": {
                        "get_param": "network_id"
                    }
              }],
			  "availability_zone": {
				     "get_param": "availability_zone"
				}
            }
        }
    },
   "outputs":  {
      "InstanceIP":{
        "description": "Instance IP",
        "value": {  "get_attr": ["my_instance", "first_address"]  }
      }
    }
}
STACK
}

resource "opentelekomcloud_rts_software_config_v1" "myconfig" {
  name = "config_name"
}

resource "opentelekomcloud_rts_software_config_v1" "myconfig2" {
  #inputs = [
  #    {
  #        "default"= "null",
  #        "type" = "String",
  #        "name" = "foo",
  #        "description" = "null"
  #    }
  #]
  group = "script"
  name  = "a-config-we5zpvyu7b5o"
  # outputs  = [
  #     {
  #         "type" ="String",
  #         "name" = "result",
  #         "error_output" = false,
  #        "description" = "null"
  #    }
  # ]
  config  = "#!/bin/sh -x\necho \"Writing to /tmp/$bar\"\necho $foo > /tmp/hanmeina.txt"
  options = {}
}

resource "opentelekomcloud_rts_software_deployment_v1" "mydeployment1" {
  config_id = opentelekomcloud_rts_software_config_v1.myconfig.id
  server_id = "6ff183f8-a210-42ef-b6dd-2af0aaabc451"
}


resource "opentelekomcloud_rts_software_deployment_v1" "mydeployment2" {
  config_id = opentelekomcloud_rts_software_config_v1.myconfig2.id
  server_id = "e387232b-c8fb-4b05-9cb5-02d18c3ea939"
  status    = "COMPLETE"
  action    = "UPDATE"
  input_values = {
    "deploy_stdout" = "Writing to /tmp/baaaaa\nWritten to /tmp/baaaaa\n",
  }
  output_values = {
    "deploy_stdout"      = "Writing to /tmp/baaaaa\nWritten to /tmp/baaaaa\n",
    "deploy_stderr"      = "+ echo Writing to /tmp/baaaaa\n+ echo fooooo\n+ cat /tmp/baaaaa\n+ echo -n The file /tmp/baaaaa contains fooooo for server ec14c864-096e-4e27-bb8a-2c2b4dc6f3f5 during CREATE\n+ echo Written to /tmp/baaaaa\n+ echo Output to stderr\nOutput to stderr\n",
    "deploy_status_code" = 0,
    "result"             = "The file /tmp/baaaaa contains fooooo for server ec14c864-096e-4e27-bb8a-2c2b4dc6f3f5 during CREATE"

  }
  status_reason = "Outputs received"
  tenant_id     = "b730519ca7064da2a3233e86bee139e4"
}
