---
subcategory: "Cloud Container Engine (CCE)"
---

Up-to-date reference of API arguments for CCE addons you can get at
`https://docs.otc.t-systems.com/cloud-container-engine/api-ref/apis/add-on_management/`.

# opentelekomcloud_cce_addon_v3

Manages a V3 CCE Addon resource within OpenTelekomCloud.

## Example Usage

### Basic addon setting

```hcl
variable "flavor_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource "opentelekomcloud_cce_cluster_v3" "cluster_1" {
  name                    = "cce-cluster-1"
  cluster_type            = "VirtualMachine"
  flavor_id               = var.flavor_id
  vpc_id                  = var.vpc_id
  subnet_id               = var.subnet_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
  cluster_version         = "v1.17.9-r0"
}

resource "opentelekomcloud_cce_addon_v3" "addon" {
  template_name    = "metrics-server"
  template_version = "1.0.6"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic = {
      "image_version" : "v0.3.7",
      "swr_addr" : "100.125.7.25:20202",
      "swr_user" : "hwofficial"
    }
    custom = {}
  }
}
```

### CCE addon setting with flavor in json format

```hcl
variable "flavor_id" {}
variable "vpc_id" {}
variable "subnet_id" {}

resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "%s"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id               = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  cluster_version         = "v1.25"
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource "opentelekomcloud_cce_addon_v3" "autoscaler" {
  template_name    = "autoscaler"
  template_version = "1.25.7"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic = {
      "cceEndpoint" : "https://cce.eu-de.otc.t-systems.com",
      "ecsEndpoint" : "https://ecs.eu-de.otc.t-systems.com",
      "image_version" : "1.25.7",
      "region" : "eu-de",
      "swr_addr" : "100.125.7.25:20202",
      "swr_user" : "cce-addons"
    }
    custom = {
      "cluster_id" : opentelekomcloud_cce_cluster_v3.cluster_1.id,
      "coresTotal" : 32000,
      "expander" : "priority",
      "logLevel" : 4,
      "maxEmptyBulkDeleteFlag" : 10,
      "maxNodeProvisionTime" : 15,
      "maxNodesTotal" : 1000,
      "memoryTotal" : 128000,
      "scaleDownDelayAfterAdd" : 10,
      "scaleDownDelayAfterDelete" : 11,
      "scaleDownDelayAfterFailure" : 3,
      "scaleDownEnabled" : true,
      "scaleDownUnneededTime" : 10,
      "scaleDownUtilizationThreshold" : 0.5,
      "scaleUpCpuUtilizationThreshold" : 1,
      "scaleUpMemUtilizationThreshold" : 1,
      "scaleUpUnscheduledPodEnabled" : true,
      "scaleUpUtilizationEnabled" : true,
      "tenant_id" : data.opentelekomcloud_identity_project_v3.project.id,
      "unremovableNodeRecheckTimeout" : 5
    }
    flavor = <<EOF
      {
        "description": "Has only one instance",
        "name": "Single",
        "replicas": 1,
        "resources": [
          {
            "limitsCpu": "1000m",
            "limitsMem": "1000Mi",
            "name": "autoscaler",
            "requestsCpu": "500m",
            "requestsMem": "500Mi"
          }
        ]
      }
	EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required) Name of the add-on template to be installed, for example, `coredns`.

* `template_version` - (Required) Version number of the add-on to be installed or upgraded, for example, `v1.0.0`.

* `cluster_id` - (Required) ID of cluster to install the add-on on.

* `values` - (Required) Parameters of the template to be installed or upgraded.

    * `basic` - (Required) Basic add-on information.

    * `custom` - (Required) Custom parameters of the add-on.

    * `flavor` - (Optional) Specifies the json string vary depending on the add-on.

Arguments which can be passed to the `basic` and `custom` addon parameters depends on the addon type and version.
For more detailed description of addons for k8s version `v1.17.9` see [addons description](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/blob/devel/opentelekomcloud/services/cce/addon-templates-v1.17.9.md).
For more detailed description of addons for k8s version `v1.19.8` see [addons description](https://github.com/opentelekomcloud/terraform-provider-opentelekomcloud/blob/devel/opentelekomcloud/services/cce/addon-templates-v1.19.8.md).

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `name` - Installed add-on name.

* `description` - Installed add-on description


## Import

CCE addons can be imported using the `cluster_id/addon_id`, e.g.

```sh
terraform import opentelekomcloud_cce_addon_v3.autoscaler c1881895-cdcb-4d23-96cb-032e6a3ee667/ea257959-eeb1-4c10-8d33-26f0409a755d
```
