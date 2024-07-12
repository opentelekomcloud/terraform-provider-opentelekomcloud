---
subcategory: "Autoscaling"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_as_configuration_v1"
sidebar_current: "docs-opentelekomcloud-resource-as-configuration-v1"
description: |-
  Manages a AS Configuration resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for AS configuration you can get at
[documentation portal](https://docs.otc.t-systems.com/auto-scaling/api-ref/apis/as_configurations)

# opentelekomcloud_as_configuration_v1

Manages a V1 AS Configuration resource within OpenTelekomCloud.

## Example Usage

### Basic AS Configuration

```hcl
resource "opentelekomcloud_as_configuration_v1" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor = var.flavor
    image  = var.image_id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }

    key_name  = var.keyname
    user_data = file("userdata.txt")
  }
}
```

### AS Configuration With Encrypted Data Disk

```hcl
resource "opentelekomcloud_as_configuration_v1" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor = var.flavor
    image  = var.image_id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    disk {
      size        = 100
      volume_type = "SATA"
      disk_type   = "DATA"
      kms_id      = var.kms_id
    }

    key_name  = var.keyname
    user_data = file("userdata.txt")
  }
}
```

### AS Configuration With User Data and Metadata

```hcl
resource "opentelekomcloud_as_configuration_v1" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor = var.flavor
    image  = var.image_id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name  = var.keyname
    user_data = file("userdata.txt")
    metadata = {
      some_key = "some_value"
    }
  }
}
```

`user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

### AS Configuration uses the existing instance specifications as the template

```hcl
resource "opentelekomcloud_as_configuration_v1" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config = {
    instance_id = "4579f2f5-cbe8-425a-8f32-53dcb9d9053a"
    key_name    = var.keyname
  }
}
```

### AS Configuration With Security Groups

```hcl
resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name        = "acc-test-sg-1"
  description = "Security group for AS config tf test"
}

resource "opentelekomcloud_as_configuration_v1" "my_as_config" {
  scaling_configuration_name = "my_as_config"

  instance_config {
    flavor = var.flavor
    image  = var.image_id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
    key_name  = var.keyname
    user_data = file("userdata.txt")
    security_groups = [
      opentelekomcloud_compute_secgroup_v2.secgroup_1.id
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V1 AutoScaling client. If omitted, the
  `region` argument of the provider is used. Changing this creates a new group.

* `scaling_configuration_name` - (Required) The name of the AS configuration. The name can contain letters,
  digits, underscores(_), and hyphens(-), and cannot exceed 64 characters.

* `instance_config` - (Required) The information about instance configurations. The instance_config
  dictionary data structure is documented below.

The `instance_config` block supports:

* `instance_id` - (Optional) When using the existing instance specifications as the template to
  create AS configurations, specify this argument. In this case, `flavor`, `image`,
  and `disk` arguments do not take effect. If the `instance_id` argument is not specified,
  `flavor`, `image`, and `disk` arguments are mandatory.

* `flavor` - (Optional) The flavor ID. By default, it tries to get from env by `OS_FLAVOR_ID`.

* `image` - (Optional) The image ID. By default, it tries to get from env by `OS_IMAGE_ID`.

* `disk` - (Optional) The disk group information. System disks are mandatory and data disks are optional.
  The dick structure is described below.

* `key_name` - (Required) The name of the SSH key pair used to log in to the instance.

* `user_data` - (Optional) The user data to provide when launching the instance.
  The file content must be encoded with Base64.

* `personality` - (Optional) Customize the personality of an instance by
  defining one or more files and their contents. The personality structure
  is described below.

* `public_ip` - (Optional) The elastic IP address of the instance. The public_ip structure
  is described below.

* `metadata` - (Optional) Metadata key/value pairs to make available from
  within the instance.

* `security_groups` - (Optional) An array of one or more security group IDs
  to associate with the AS configuration.

The `disk` block supports:

* `size` - (Required) The disk size. The unit is GB. The system disk size ranges from `4` to `32768` and must
  be greater than or equal to the minimum size (`min_disk` value) of the system disk specified in the image.
  The data disk size ranges from `10` to `32768`.

* `volume_type` - (Required) Specifies the ECS system disk type. The disk type must match the available disk type.
  * `SATA`: common I/O disk type.
  * `SAS`: high I/O disk type.
  * `SSD`: ultra-high I/O disk type.
  * `co-p1`: high I/O (performance-optimized I) disk type.
  * `uh-l1`: ultra-high I/O (latency-optimized) disk type.

->For HANA, `HL1`, and `HL2` ECSs, use `co-p1` and `uh-l1` disks. For other ECSs, do not use `co-p1` or `uh-l1` disks.

* `disk_type` - (Required) Specifies a disk type. The options are as follows:
  * `DATA`: indicates a data disk.
  * `SYS`: indicates a system disk.

* `kms_id` - (Optional) The Encryption KMS ID of the data disk.

The `personality` block supports:

* `path` - (Required) The absolute path of the destination file.

* `contents` - (Required) The content of the injected file, which must be encoded with base64.

The `public_ip` block supports:

* `eip` - (Required) The configuration parameter for creating an elastic IP address
  that will be automatically assigned to the instance. The eip structure is described below.

The `eip` block supports:

* `ip_type` - (Required) The IP address type. The system only supports `5_bgp` and `5_mailbgp`.

* `bandwidth` - (Required) The bandwidth information. The structure is described below.

The `bandwidth` block supports:

* `size` - (Required) The bandwidth (Mbit/s). The value range is 1 to 500.

->The specific range may vary depending on the configuration in each region. You can see the bandwidth range of
  each region on the management console. The minimum unit is 1 Mbit/s if the allowed bandwidth size ranges from
  0 to 300 Mbit/s. The minimum unit is 50 Mbit/s if the allowed bandwidth size ranges 300 Mbit/s to 500 Mbit/s.

* `share_type` - (Required) The bandwidth sharing type. The system only supports `PER`.

* `charging_mode` - (Required) The bandwidth charging mode. The system only supports `traffic`.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `id` - Specifies the AS configuration ID.

* `scaling_configuration_name` - Specifies the AS configuration name.

* `instance_config` - Specifies the information about instance configurations.

The `instance_config` block supports:

* `flavor` - Specifies the ECS flavor ID.

* `image` - Specifies the image ID.

* `key_name` - See Argument Reference above.

* `user_data` - See Argument Reference above.

* `region` - See Argument Reference above.
