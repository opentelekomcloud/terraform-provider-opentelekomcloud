---
subcategory: "Identity and Access Management (IAM)"
---

# opentelekomcloud_identity_agency_v3

Use this data source to get an agency information.

## Example Usage

```hcl
data "opentelekomcloud_identity_agency_v3" "agency_1" {
  name = "test-agency"
}
```

## Argument Reference

* `name` - (Optional) Name of the agency

* `domain_id` - (Optional) ID of the current domain.

* `trust_domain_id` - (Optional) ID of the delegated domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of an agency.

* `name` - Name of an agency.

* `trust_domain_name` - Name of the delegated domain.

* `description` - Description of an agency.

* `duration` - Validity period of an agency.
  The default value is `null`, indicating that the agency is permanently valid.

* `expire_time` - Expiration time of an agency.

* `create_time` - Time when an agency is created.
