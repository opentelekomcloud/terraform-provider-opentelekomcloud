---
subcategory: "Cloud Container Engine (CCE)"
---

# opentelekomcloud_cce_addon_v3

Provides a cluster addon management.

## Example Usage

```hcl
variable "flavor_id" { }
variable "vpc_id" { }
variable "subnet_id" { }

resource opentelekomcloud_cce_cluster_v3 cluster_1 {
  name                    = "cce-cluster-1"
  cluster_type            = "VirtualMachine"
  flavor_id               = "cce.s1.small"
  vpc_id                  = var.vpc_id
  subnet_id               = var.subnet_id
  container_network_type  = "overlay_l2"
  kubernetes_svc_ip_range = "10.247.0.0/16"
}

resource opentelekomcloud_cce_addon_v3 addon {
  template_name    = "metrics-server"
  template_version = "1.0.3"
  cluster_id       = opentelekomcloud_cce_cluster_v3.cluster_1.id

  values {
    basic = {
      euleros_version = "2.5"
      rbac_enabled    = true
      swr_addr        = "100.125.7.25:20202"
      swr_user        = "hwofficial"
    }
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

    * `custom` - (Optional) Custom parameters of the add-on.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `name` - Installed add-on name.

* `description` - Installed add-on description
