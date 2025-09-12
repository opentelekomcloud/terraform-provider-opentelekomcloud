---
subcategory: "Application Service Mesh (ASM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_asm_service_mesh_v1"
sidebar_current: "docs-opentelekomcloud-resource-asm-service-mesh-v1"
description: |-
  Manages an ASM Service Mesh v1 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for ASM service mesh you can get at
[documentation portal](https://docs.otc.t-systems.com/application-service-mesh/api-ref/api/service_mesh_apis/index.html)

# opentelekomcloud_asm_service_mesh_v1

Manages an ASM Service Mesh v1 resource within OpenTelekomCloud.

## Example Usage

### Basic ASM service mesh

```hcl
variable "cluster_id" {}
variable "node_id" {}

resource "opentelekomcloud_asm_service_mesh_v1" "mesh_1" {
  name    = "test-asm-service-mesh"
  type    = "InCluster"
  version = "1.18.7-r5"

  clusters {
    cluster_id         = var.cluster_id
    installation_nodes = [var.node_id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the service mesh name. The name must be 4 to 64 characters long, start with a lowercase letter and not end with a hyphen (-). Only lowercase letters, digits, and hyphens (-) are allowed.

* `type` - (Required, String, ForceNew) Specifies the service mesh type. Supported value: `InCluster` (service mesh with an in-cluster control plane). The value is `InCluster` for the service mesh of the Basic edition.

* `version` - (Required, String, ForceNew) Specifies the service mesh version.

* `clusters` - (Required, List, ForceNew) Specifies the cluster information in the service mesh. The [clusters](#clusters) structure is documented below.

* `ipv6_enable` - (Optional, Boolean, ForceNew) Specifies whether the service mesh supports IPv6.

* `proxy_config` - (Optional, List, ForceNew) Specifies the data plane configuration of the service mesh. The [proxy_config](#proxy_config) structure is documented below.

* `telemetry_config_tracing` - (Optional, List, ForceNew) Specifies the observability/tracing configuration, which is used to report traces in the service mesh. The [telemetry_config_tracing](#telemetry_config_tracing) structure is documented below.

<!-- ### CLUSTERS BLOCK ### -->

<a name="clusters"></a>
The `clusters` block supports:

* `cluster_id` - (Required, String, ForceNew) Specifies the CCE Cluster ID.

* `installation_nodes` - (Required, List[String], ForceNew) Specifies the node ID of nodes where service mesh components are to be installed.

* `injection_namespaces` - (Optional, List[String], ForceNew) Specifies the name of namespaces where sidecars are to be injected.

<!-- ### PROXY CONFIG BLOCK ### -->

<a name="proxy_config"></a>
The `proxy_config` block supports:

* `include_ip_ranges` - (Optional, String, ForceNew) Specifies the IP address ranges that will be included for outbound traffic redirection. Use commas (,) to separate the IP address ranges.

* `exclude_ip_ranges` - (Optional, String, ForceNew) Specifies the IP address ranges that will be excluded for outbound traffic redirection. Use commas (,) to separate the IP address ranges.

* `exclude_outbound_ports` - (Optional, String, ForceNew) Specifies the ports that will be excluded for outbound traffic redirection. Use commas (,) to separate the ports.

* `exclude_inbound_ports` - (Optional, String, ForceNew) Specifies the ports that will be excluded for inbound traffic redirection. Use commas (,) to separate the ports.

* `include_outbound_ports` - (Optional, String, ForceNew) Specifies the ports that will be included for outbound traffic redirection. Use commas (,) to separate the ports.

* `include_inbound_ports` - (Optional, String, ForceNew) Specifies the Ports that will be included for inbound traffic redirection. Use commas (,) to separate the ports.

<!-- ### TELEMETRY CONFIG BLOCK ### -->

<a name="telemetry_config_tracing"></a>
The `telemetry_config_tracing` block supports:

* `random_sampling_percentage` - (Optional, Float, ForceNew) Specifies the tracing sampling rate.

* `default_providers` - (Optional, List[String], ForceNew) Specifies the name of the default provider that tracing reports data to, which must match the name field in `extension_providers` or use the preset provider `apm-otel` of ASM. If `apm-otel` is used, ensure that APM 2.0 is supported in the current region and the service mesh version is later than 1.18.

* `extension_providers` - (Optional, List, ForceNew) Specifies the user-defined provider. Currently, Zipkin is supported. If you configure the Zipkin provider, ensure that the service mesh version is 1.15 or later. The structure is documented below:
    * `name` - (Optional, String, ForceNew) Specifies the provider name.
    * `zipkin_service_addr` - (Optional, String, ForceNew) Specifies the service address of Zipkin.
    * `zipkin_service_port` - (Optional, String, ForceNew) Specifies the service port of Zipkin.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Service Mesh ID.

* `cluster_ids` - Indicates the cluster id of CCE clusters associated with service mesh.

* `creation_timestamp` - Indicates the time when the service mesh was created..

* `status` - Indicates the service mesh status. The value can be: `Running`, `Creating`, `CreateFailed`, `Deleting`, `DeleteFailed`, `Upgrading`, `UpgradeFailed`, `RollingBack`, `RollbackFailed`.

## Import

ASM Service Mesh V1 resource can be imported using the service mesh ID, `id`, e.g.

```shell
terraform import opentelekomcloud_asm_service_mesh_v1.mesh_1 <id>
```

## Notes

But due to some attributes missing from the API response, it's required to ignore changes as below:

```hcl
resource "opentelekomcloud_asm_service_mesh_v1" "mesh_1" {
  # ...

  lifecycle {
    ignore_changes = [
      clusters,
    ]
  }
}
```
