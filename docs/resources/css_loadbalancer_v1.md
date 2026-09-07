---
subcategory: "Cloud Search Service (CSS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_css_loadbalancer_v1"
sidebar_current: "docs-opentelekomcloud-resource-css-loadbalancer-v1"
description: |-
  Manages a CSS load balancing resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CSS load balancing you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-search-service/api-ref/apis/load_balancing)

# opentelekomcloud_css_loadbalancer_v1

Manages a CSS configuration of loadbalancer.

## Example Usage

```hcl
resource "opentelekomcloud_css_loadbalancer_v1" "css_lb" {
  cluster_id = "terraform-test-cluster"
  elb_id     = var.elb_id
  agency     = "css_upgrade_agency"

  listener {
    protocol       = "HTTPS"
    protocol_port  = 443
    server_cert_id = var.server_cert_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) ID of the CSS cluster.

* `agency` - (Required) The agency used by CSS to access ELB.

* `elb_id` - (Required) ELBv3 ID.

* `listener` - (Optional) configure load balancing listeners for a cluster. Structure is documented below.

The `listener` block supports:

* `protocol` - (Required) Protocol type. HTTP and HTTPS are supported.

* `protocol_port` - (Required) Port.

* `server_cert_id` - (Optional) Server certificate ID. This parameter is mandatory when protocol is set to HTTPS.

* `ca_cert_id` - (Optional) CA certificate ID. This parameter is mandatory when protocol is set to HTTPS and bidirectional authentication is used.

* `type` - (Optional) Type: searchTool indicates that the load balancer is enabled or disabled for Elasticsearch/OpenSearch. viewTool indicates that the load balancer is enabled or disabled for Kibana/OpenSearch Dashboards. The default value is searchTool.
