---
subcategory: "Anti-DDoS"
---

# opentelekomcloud_antiddos_v1

Anti-DDoS monitors the service traffic from the Internet to ECSs, ELB instances, and BMSs to detect attack traffic in real time. It then cleans attack traffic according to user-configured defense policies so that services run as normal.

~>
AntiDDoS protection for ElasticIP is provided by default and shouldn't be created.


## Example Usage

```hcl
variable "eip_id" {}

resource "opentelekomcloud_antiddos_v1" "myantiddos" {
  floating_ip_id         = var.eip_id
  enable_l7              = true
  traffic_pos_id         = 1
  http_request_pos_id    = 3
  cleaning_access_pos_id = 2
  app_type_id            = 0
}
```

## Argument Reference

The following arguments are supported:

* `enable_l7` - (Required) Specifies whether to enable L7 defense.

* `traffic_pos_id` - (Required) The position ID of traffic. The value ranges from 1 to 9.

* `http_request_pos_id` - (Required) The position ID of number of HTTP requests. The value ranges from 1 to 15.

* `cleaning_access_pos_id` - (Required) The position ID of access limit during cleaning. The value ranges from 1 to 8.

* `app_type_id` - (Required) The application type ID.

* `floating_ip_id` - (Required) The ID corresponding to the Elastic IP Address (EIP) of a user.

## Attributes Reference

All above argument parameters can be exported as attribute parameters.

## Import

Antiddos can be imported using the floating_ip_id, e.g.

```sh
terraform import opentelekomcloud_antiddos_v1.myantiddos c1881895-cdcb-4d23-96cb-032e6a3ee667
```
