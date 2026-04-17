---
subcategory: "Application Service Mesh (ASM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_asm_service_mesh_v1"
sidebar_current: "docs-opentelekomcloud-data-source-asm-service-mesh-v1"
description: |-
  Manages an ASM Service Mesh v1 data source within OpenTelekomCloud.
---

Up-to-date reference of API arguments for ASM service mesh you can get at
[documentation portal](https://docs.otc.t-systems.com/application-service-mesh/api-ref/api/service_mesh_apis/index.html)

# opentelekomcloud_asm_service_mesh_v1

Manages an ASM Service Mesh v1 data source within OpenTelekomCloud.

## Example Usage

### List all ASM service meshes

```hcl
data "opentelekomcloud_asm_service_mesh_v1" "mesh_1" {}
```

### Get ASM service mesh using ID

```hcl
variable "mesh_id" {}

data "opentelekomcloud_asm_service_mesh_v1" "mesh_1" {
  id = var.mesh_id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) Specifies the service mesh ID.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `service_meshes` - The service mesh list. The structure is documented below:
  
The `service_meshes` block supports:

* `name` - Indicates the service mesh name.

* `type` - Indicates the service mesh type.

* `version` - Indicates the service mesh version.

* `ipv6_enable` - Indicates whether the service mesh supports IPv6.

* `cluster_ids` - Indicates the cluster id of CCE clusters associated with service mesh.

* `creation_timestamp` - Indicates the time when the service mesh was created..

* `status` - Indicates the service mesh status. The value can be: `Running`, `Creating`, `CreateFailed`, `Deleting`, `DeleteFailed`, `Upgrading`, `UpgradeFailed`, `RollingBack`, `RollbackFailed`.

* `proxy_config` - Indicates the data plane configuration of the service mesh. The structure is documented below.
    * `include_ip_ranges` - Indicates the IP address ranges that will be included for outbound traffic redirection.
    * `exclude_ip_ranges` - Indicates the IP address ranges that will be excluded for outbound traffic redirection.
    * `exclude_outbound_ports` - Indicates the ports that will be excluded for outbound traffic redirection.
    * `exclude_inbound_ports` - Indicates the ports that will be excluded for inbound traffic redirection.
    * `include_outbound_ports` - Indicates the ports that will be included for outbound traffic redirection.
    * `include_inbound_ports` - Indicates the Ports that will be included for inbound traffic redirection.

* `telemetry_config_tracing` - Indicates the observability/tracing configuration, which is used to report traces in the service mesh. The structure is documented below.

The `telemetry_config_tracing` block supports:

* `random_sampling_percentage` - Indicates the tracing sampling rate.
* `default_providers` - Indicates the name of the default provider that tracing reports data to.
* `extension_providers` - Indicates the user-defined provider. Currently, Zipkin is supported.The structure is documented below:
    * `name` - Indicates the provider name.
    * `zipkin_service_addr` - Indicates the service address of Zipkin.
    * `zipkin_service_port` - Indicates the service port of Zipkin.
