# FULL Terraform OTC Example

This script will create the following resources (if enabled):
* Volumes
* Floating IPs
* Neutron Ports
* Instances
* Keypair
* Network
* Subnet
* Router
* Router Interface
* Loadbalancer
* Templates
* Security Group (Allow ICMP, 80/tcp, 22/tcp)
* Subscription

## Resource Creation

This example will, by default not create Volumes. This is to show how to enable resources via parameters. To enable Volume creation, set the **disk_\_size\_gb** variable to a value > 10.

## Available Variables

### Required

* **username** (your OTC username)
* **password** (your OTC password)
* **domain\_name** (your OTC domain name)
* You must have a **ssh\_pub\_key** file defined, or terraform will complain, see default path below.

### Optional
* **project** (this will prefix all your resources, _default=terraform_)
* **ssh\_pub\_key** (the path to the ssh public key you want to deploy, _default=~/.ssh/id\_rsa.pub_)
* **instance\_count** (affects the number of Floating IPs, Instances, Volumes and Ports, _default=1_)
* **flavor\_name** (flavor of the created instances, _default=s1.medium_)
* **image\_name** (image used for creating instances, _default=Standard\_CentOS\_7\_latest_)
* **disk\_size\_gb** (size of the volumes in gigabytes, _default=None_)
* **endpoint\_email** (The email endpoint for creating subscriptions, _default=mailtest@gmail.com)
* **endpoint\_sms** (The sms endpoint for creating subscriptions, _default=+8613600000000)
