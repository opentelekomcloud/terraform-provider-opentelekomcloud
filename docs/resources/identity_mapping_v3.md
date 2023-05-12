---
subcategory: "Identity and Access Management (IAM)"
---

Up-to-date reference of API arguments for IAM mapping you can get at
`https://docs.otc.t-systems.com/identity-access-management/api-ref/apis/federated_identity_authentication_management/mapping`.

# opentelekomcloud_identity_mapping_v3

-> You _must_ have security admin privileges in your OpenTelekomCloud
cloud to use this resource. Please refer to [User Management Model](https://docs.otc.t-systems.com/en-us/usermanual/iam/iam_01_0034.html).


## Example Usage

```hcl
resource "opentelekomcloud_identity_mapping_v3" "mapping" {
  mapping_id = "ACME"
  rules      = <<EOF
  [
    {
      "local":[
        {
          "user":{"name":"{0}"}
        },
        {
          "groups":"[\"admin\",\"manager\"]"
        }
      ],
      "remote":[
        {
          "type":"uid"
        }
      ]
    }
  ]
EOF
}
```

## Argument Reference

The following arguments are supported:

* `mapping_id` - (Required) The ID of the mapping. Changing this creates a new mapping.

* `rules` - (Required) Rules used to map federated users to local users.

-> For the full reference checkout [Syntax of Identity Conversion Rules](https://docs.otc.t-systems.com/en-us/usermanual/iam/en-us_topic_0079620340.html).

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

`links` - Resource links of an identity mapping.

## Import

Mappings can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_identity_mapping_v3.mapping ACME
```

