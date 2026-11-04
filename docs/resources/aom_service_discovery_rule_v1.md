---
subcategory: "Application Operations Management (AOM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_aom_service_discovery_rule_v1"
sidebar_current: "docs-opentelekomcloud-resource-aom-service-discovery-rule-v1"
description: |-
  Manages an AOM service discovery rule resource within T-Cloud Public (formerly OpenTelekomCloud).
---

# opentelekomcloud_aom_service_discovery_rule_v1

Manages an AOM service discovery rule resource within T-Cloud Public (formerly OpenTelekomCloud).

Up-to-date reference of API arguments for AOM monitoring you can get at
[documentation portal](https://docs.otc.t-systems.com/application-operations-management-2/api-ref/apis/monitoring/index.html)

-> **Note:** A maximum of 100 service discovery rules can be added to a project.

## Example Usage

```hcl
resource "opentelekomcloud_aom_service_discovery_rule_v1" "test" {
  name        = "service-discovery-rule"
  enabled     = true
  description = "discover Java processes"

  spec {
    detect_log      = false
    is_default_rule = false
    is_detect       = false
    priority        = 9999
    app_type        = "Java"
    log_file_fix    = ["log", "trace", "out"]

    discovery_rule {
      check_type    = "cmdLine"
      check_mode    = "contain"
      check_content = ["my-application-process"]
    }

    name_rule {
      app_name_rule {
        name_type = "str"
        args      = ["-svc"]
      }

      application_name_rule {
        name_type = "str"
        args      = ["-app"]
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the service discovery rule.
  A maximum of 64 characters are allowed. The name must start with a lowercase
  letter, cannot end with a hyphen (-), and only digits, lowercase letters,
  and hyphens are allowed.

* `enabled` - (Required, Bool) Whether the rule is enabled.

* `spec` - (Required, List) The details of the service discovery rule.
  The [spec](#spec) structure is documented below.

* `description` - (Optional, String) The custom description of the rule.

<a name="spec"></a>The `spec` block supports:

* `detect_log` - (Required, Bool) Whether to enable log collection.

* `is_default_rule` - (Required, Bool) Whether the rule is the default one.

* `is_detect` - (Required, Bool) Whether the scenario is a pre-check scenario.
  No rules will be saved in the pre-check scenario.

* `priority` - (Required, Int) The rule priority. Value range: `1` to `9999`.
  Default value: `9999`.

* `log_file_fix` - (Required, List) The log file suffixes.
  Valid values are `log`, `trace`, and `out`.

* `discovery_rule` - (Required, List) The process match conditions.
  If the list contains multiple conditions, only the processes that meet all
  the conditions will be matched.
  The `discovery_rule` structure is documented below.

* `name_rule` - (Required, List) The naming rules for discovered services and
  applications. The `name_rule` structure is documented below.

* `log_path_rule` - (Optional, List) The log path configuration rules.
  The `log_path_rule` structure is documented below.

* `app_type` - (Optional, String) The service type, used only for rule
  classification and UI display. Any value can be entered, for example,
  `Java` or `Python`.

* `attr_list` - (Optional, List) The attribute list. Currently not used.
  Valid values are `cmdLine` and `env`.

* `data_source` - (Optional, String) The data source.

The `discovery_rule` block supports:

* `check_type` - (Required, String) The match type.
  Valid values are `cmdLine`, `env`, and `scope`.

* `check_mode` - (Required, String) The match condition.
  Valid values are `contain` and `equals`.
  For `check_type` `cmdLine` or `env`, use `contain`; for `scope`, use `equals`.

* `check_content` - (Required, List) The matched values.
  For `cmdLine`: the process must contain the given parameter.
  For `env`: pairs in the format `["k1", "v1"]`, the process must contain the
  environment variable `k1` with value `v1`.
  For `scope`: host IDs the rule applies to; if empty, all project nodes.

The `name_rule` block supports:

* `app_name_rule` - (Required, List) The service name extraction rules.
  Documented under the name extraction rules below.

* `application_name_rule` - (Required, List) The application name extraction
  rules, with the same format as `app_name_rule`.

The name extraction rules (`app_name_rule`, `application_name_rule`) support:

* `name_type` - (Required, String) The value type.
  Valid values are `cmdLineHash`, `cmdLine`, `env`, and `str`.
  For `cmdLine`, `args` is `["start", "end"]` and the characters between
  `start` and `end` in the command are extracted.
  For `env`, `args` is `["aa"]` and the environment variable `aa` is extracted.
  For `str`, `args` is `["fix"]` and the name is suffixed with `fix`.
  For `cmdLineHash`, `args` is `["0001"]` and `value` is `["ser"]`, meaning
  the name is `ser` when the startup command is `0001`.

* `args` - (Required, List) The input values.

* `value` - (Optional, List) The name, mandatory only when `name_type` is
  `cmdLineHash`.

The `log_path_rule` block supports:

* `name_type` - (Required, String) The value type, which can be `cmdLineHash`.

* `args` - (Required, List) The commands.

* `value` - (Required, List) The log paths.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the ID of the service discovery rule.

* `event_name` - The rule event name
  (fixed value: `aom_inventory_rules_event`).

* `project_id` - The project ID the rule belongs to.

* `create_time` - The creation time returned by the API.

## Import

AOM service discovery rules can be imported using the rule ID, e.g.

```
terraform import opentelekomcloud_aom_service_discovery_rule_v1.test <rule_id>
```
