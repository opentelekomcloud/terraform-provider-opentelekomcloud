---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_api_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-api-v2"
description: |-
  Manages a APIGW API resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for API Gateway API service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/api_management/index.html)

# opentelekomcloud_apigw_api_v2

Provides an API gateway API resource.

## Example Usage

```hcl
variable "vpc_id" {}
variable "subnet_id" {}
variable "secgroup_id" {}

resource "opentelekomcloud_apigw_gateway_v2" "gateway" {
  name                            = "my_gw"
  spec_id                         = "BASIC"
  vpc_id                          = var.vpc_id
  subnet_id                       = var.subnet_id
  security_group_id               = var.secgroup_id
  availability_zones              = ["eu-de-01", "eu-de-02"]
  description                     = "test gateway 2"
  ingress_bandwidth_size          = 5
  ingress_bandwidth_charging_mode = "bandwidth"
  maintain_begin                  = "02:00:00"
}

resource "opentelekomcloud_apigw_environment_v2" "env" {
  name        = "my_env"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"
}

resource "opentelekomcloud_apigw_group_v2" "group" {
  name        = "my_group"
  instance_id = opentelekomcloud_apigw_gateway_v2.gateway.id
  description = "test description"

  environment {
    variable {
      name  = "test-name"
      value = "test-value"
    }
    environment_id = opentelekomcloud_apigw_environment_v2.env.id
  }
}

resource "opentelekomcloud_apigw_api_v2" "api" {
  gateway_id                   = opentelekomcloud_apigw_gateway_v2.gateway.id
  group_id                     = opentelekomcloud_apigw_group_v2.group.id
  name                         = "my_api"
  type                         = "Public"
  request_protocol             = "HTTP"
  request_method               = "GET"
  request_uri                  = "/user_info/{user_age}"
  security_authentication_type = "APP"
  match_mode                   = "EXACT"
  success_response             = "Success response"
  failure_response             = "Failed response"
  description                  = "Created by script"

  request_params {
    name     = "user_age"
    type     = "NUMBER"
    location = "PATH"
    required = true
    maximum  = 200
    minimum  = 0
  }
  request_params {
    name        = "X-TEST-ENUM"
    type        = "STRING"
    location    = "HEADER"
    maximum     = 20
    minimum     = 10
    sample      = "ACC_TEST_XXX"
    passthrough = true
    enumeration = "ACC_TEST_A,ACC_TEST_B"
  }

  backend_params {
    type     = "REQUEST"
    name     = "userAge"
    location = "PATH"
    value    = "user_age"
  }

  http {
    url_domain       = "opentelekomcloud.my.com"
    request_uri      = "/getUserAge/{userAge}"
    request_method   = "GET"
    request_protocol = "HTTP"
    timeout          = 30000
    retry_count      = 1
  }

  http_policy {
    url_domain       = "opentelekomcloud.my.com"
    name             = "my_policy1"
    request_protocol = "HTTP"
    request_method   = "GET"
    effective_mode   = "ANY"
    request_uri      = "/getUserAge/{userAge}"
    timeout          = 30000
    retry_count      = 1

    backend_params {
      type     = "REQUEST"
      name     = "userAge"
      location = "PATH"
      value    = "user_age"
    }
    backend_params {
      type              = "SYSTEM"
      name              = "%[2]s"
      location          = "HEADER"
      value             = "serverName"
      system_param_type = "internal"
    }

    conditions {
      origin     = "param"
      param_name = "user_age"
      type       = "EXACT"
      value      = "28"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `region` - (Optional, String, ForceNew) Specifies the region where the API is located.
  If omitted, the provider-level region will be used. Changing this will create a new API resource.

* `gateway_id` - (Required, String, ForceNew) Specifies an ID of the APIG dedicated instance to which the API belongs
  to. Changing this will create a new API resource.

* `group_id` - (Required, String) Specifies an ID of the APIG group to which the API belongs to.

* `type` - (Required, String) Specifies the API type.
  The valid values are `Public` and `Private`.

* `version` - (Optional, String) Specifies the API version.

* `name` - (Required, String) Specifies the API name.
  The valid length is limited from can contain `3` to `255`, only Chinese and English letters, digits and
  following special characters are allowed: `-_./（()）:：、`.
  The name must start with a digit, Chinese or English letter.

* `request_method` - (Required, String) Specifies the request method of the API.
  The valid values are `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`, `OPTIONS` and `ANY`.

* `request_uri` - (Required, String) Specifies the request address, which can contain a maximum of `512` characters,
  the request parameters enclosed with brackets ({}).
  + The address can contain special characters, such as asterisks (*), percent signs (%), hyphens (-), and
    underscores (_) and must comply with URI specifications.
  + The address can contain environment variables, each starting with a letter and consisting of `3` to `32` characters.

  Only letters, digits, hyphens (-), and underscores (_) are allowed in environment variables.

* `request_protocol` - (Required, String) Specifies the request protocol of the API.
  The valid values are `HTTP`, `HTTPS` and `BOTH`.

* `security_authentication_type` - (Optional, String) Specifies the security authentication mode of the API request.
  The valid values are `NONE`, `APP`, `IAM` and `AUTHORIZER`, defaults to `NONE`.

* `security_authentication_enabled` - (Optional, Bool) Specifies whether the authentication of the application code is enabled.
  The application code must located in the header when `security_authentication_enabled` is `true`.

* `authorizer_id` - (Optional, String) Specifies the ID of the authorizer to which the API request used.
  It is Required when `security_authentication_type` is `AUTHORIZER`.

* `body_description` - (Optional, String) Specifies the description of the API request body, which can be an example
  request body, media type or parameters.
  The request body does not exceed `20,480` characters.

* `cors` - (Optional, Bool) Specifies whether CORS is supported, defaults to `false`.

* `description` - (Optional, String) Specifies the API description.
  The description contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

* `match_mode` - (Optional, String) Specifies the route matching mode.
  The valid values are `EXACT` and `PREFIX`, defaults to `EXACT`.

* `success_response` - (Optional, String) Specifies the example response for a successful request.
  The response contains a maximum of `20,480` characters.

* `failure_response` - (Optional, String) Specifies the example response for a failure request.
  The response contains a maximum of `20,480` characters.

* `response_id` - (Optional, String) Specifies the APIG group response ID.

* `tags` - (Optional, List) Tags. Use letters, digits, and special characters `(-*#%.:_)` and start with a letter.
  By default, `10` tags are supported. To increase the quota, contact technical support to modify the `API_TAG_NUM_LIMIT` configuration.

* `request_params` - (Optional, List) Specifies the configurations of the front-end parameters.
  The [object](#apigw_api_request_params) structure is documented below.

* `backend_params` - (Optional, List) Specifies the configurations of the backend parameters.
  The [object](#apigw_api_backend_params) structure is documented below.

* `mock` - (Optional, List, ForceNew) Specifies the mock backend details.
  The [object](#apigw_api_mock) structure is documented below.
  Changing this will create a new API resource.

* `func_graph` - (Optional, List, ForceNew) Specifies the function graph backend details.
  The [object](#apigw_api_func_graph) structure is documented below.
  Changing this will create a new API resource.

* `http` - (Optional, List, ForceNew) Specifies the web backend details.
  The [object](#apigw_api_http) structure is documented below. Changing this will create a new API resource.

* `mock_policy` - (Optional, List) Specifies the Mock policy backends.
  The maximum blocks of the policy is 5.
  The [object](#apigw_api_mock_policy) structure is documented below.

* `func_graph_policy` - (Optional, List) Specifies the Mock policy backends.
  The maximum blocks of the policy is 5.
  The [object](#apigw_api_func_graph_policy) structure is documented below.

* `http_policy` - (Optional, List) Specifies the example response for a failed request.
  The maximum blocks of the policy is 5.
  The [object](#apigw_api_http_policy) structure is documented below.

<a name="apigw_api_request_params"></a>
The `request_params` block supports:

* `name` - (Required, String) Specifies the request parameter name.
  The valid length is limited from can contain `1` to `32`, only letters, digits, hyphens (-), underscores (_) and
  periods (.) are allowed.
  If Location is specified as `HEADER` and `security_authentication_type` is specified as `APP`, the parameter name
  cannot be `Authorization` (case-insensitive) and cannot contain underscores.

* `required` - (Optional, Bool) Specifies whether the request parameter is required.

* `passthrough` - (Optional, Bool) Specifies whether to transparently transfer the parameter.

* `enumeration` - (Optional, String) Specifies the enumerated value(s).
  Use commas to separate multiple enumeration values, such as `VALUE_A,VALUE_B`.

* `location` - (Optional, String) Specifies the location of the request parameter.
  The valid values are `PATH`, `QUERY` and `HEADER`, defaults to `PATH`.

* `type` - (Optional, String) Specifies the request parameter type.
  The valid values are `STRING` and `NUMBER`, defaults to `STRING`.

* `maximum` - (Optional, Int) Specifies the maximum value or size of the request parameter.

* `minimum` - (Optional, Int) Specifies the minimum value or size of the request parameter.

-> For string type, The `maximum` and `minimum` means size. For number type, they means value.

* `sample` - (Optional, String) Specifies the example value of the request parameter.
  The example contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

* `default` - (Optional, String) Specifies the default value of the request parameter.
  The value contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

* `description` - (Optional, String) Specifies the description of the request parameter.
  The description contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

<a name="apigw_api_backend_params"></a>
The `backend_params` block supports:

* `type` - (Required, String) Specifies the backend parameter type.
  The valid values are `REQUEST`, `CONSTANT` and `SYSTEM`.

* `name` - (Required, String) Specifies the backend parameter name, which contain of 1 to 32 characters and start with a
  letter. Only letters, digits, hyphens (-), underscores (_) and periods (.) are allowed. The parameter name is not
  case-sensitive. It cannot start with `x-apig-` or `x-sdk-` and cannot be `x-stage`. If the location is specified as
  `HEADER`, the name cannot contain underscores.

* `location` - (Required, String) Specifies the location of the backend parameter.
  The valid values are `PATH`, `QUERY` and `HEADER`.

* `value` - (Required, String) Specifies the request parameter name corresponding to the back-end request parameter.

* `description` - (Optional, String) Specifies the description of the constant or system parameter.
  The description contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

* `system_param_type` - (Optional, String) Specifies the type of the system parameter.
  The valid values are `frontend`, `backend` and `internal`, defaults to `internal`.

<a name="apigw_api_mock"></a>
The `mock` block supports:

* `response` - (Required, String) Specifies the response of the backend policy.
  The description contains a maximum of `2,048` characters and the angle brackets (< and >) are not allowed.

  -> **NOTE:**  Mock enables APIG to return a response without sending the request to the backend. This is useful for
  testing APIs when the backend is not available.

* `authorizer_id` - (Optional, String) Specifies the ID of the backend custom authorization.

* `description` - (Optional, String) Specifies the description of the constant or system parameter.
  The description contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

* `version` - (Optional, String) Specifies the mock version. It cannot exceed 64 characters.

<a name="apigw_api_func_graph"></a>
The `func_graph` block supports:

* `function_urn` - (Required, String) Specifies the URN of the FunctionGraph function.

* `version` - (Required, String) Specifies the function version.

* `timeout` - (Optional, Int) Specifies the timeout for API requests to backend service.
  The valid value is range form `1` to `600,000`, defaults to `5,000`.

* `invocation_type` - (Optional, String) Specifies the invocation type.
  The valid values are `async` and `sync`, defaults to `sync`.

* `network_type` - (Optional, String) Function network architecture.
  The valid values are `VPC` and `NON-VPC`, defaults to `NON-VPC`.

* `authorizer_id` - (Optional, String) Specifies the ID of the backend custom authorization.

* `description` - (Optional, String) Specifies the description of the constant or system parameter.
  The description contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

<a name="apigw_api_http"></a>
The `http` block supports:

* `authorizer_id` - (Optional, String) Specifies the ID of the backend custom authorization.

* `url_domain` - (Optional, String) Backend service address.
  A backend service address consists of a domain name or IP address and a port number,
  with not more than `255` characters. It must be in the format `Host name:Port number`,
  for example, `apig.example.com:7443`. If the port number is not specified, the default `HTTPS` port `443` or the default `HTTP` port `80` is used.

* `request_protocol` - (Optional, String) Specifies the backend request protocol.
  The valid values are `HTTP` and `HTTPS`, defaults to `HTTPS`.

* `description` - (Optional, String) Specifies the description of the constant or system parameter.
  The description contains a maximum of `255` characters and the angle brackets (< and >) are not allowed.

* `version` - (Required, String) Specifies the function version.

* `request_method` - (Optional, String) Specifies the backend request method of the API.
  The valid values are `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`, `OPTIONS` and `ANY`.

* `request_uri` - (Required, String) Specifies the backend request address, which can contain a maximum of `512` characters and
  must comply with URI specifications.
  + The address can contain request parameters enclosed with brackets ({}).
  + The address can contain special characters, such as asterisks (*), percent signs (%), hyphens (-) and
    underscores (_) and must comply with URI specifications.
  + The address can contain environment variables, each starting with a letter and consisting of `3` to `32` characters.
    Only letters, digits, hyphens (-), and underscores (_) are allowed in environment variables.

* `timeout` - (Optional, Int) Specifies the timeout for API requests to backend service, the unit is **ms**.
  The valid value ranges from `1` to `600,000`, defaults to `5,000`.

* `ssl_enable` - (Optional, Bool) Specifies whether to enable two-way authentication, defaults to **false**.

* `retry_count` - (Optional, Int) Specifies the number of retry attempts to request the backend service.
  The valid value ranges from `-1` to `10`, defaults to `-1`.
  `-1` indicates that idempotent APIs will retry once and non-idempotent APIs will not retry.
  `POST` and `PATCH` are not-idempotent.
  `GET`, `HEAD`, `PUT`, `OPTIONS` and `DELETE` are idempotent.

  -> When the (web) backend uses the channel, the `retry_count` must be less than the number of available backend
  servers in the channel.

* `vpc_channel_proxy_host` - (Optional, String) Specifies the proxy host header.
  The host header can be customized for requests to be forwarded to cloud servers through the VPC channel.
  By default, the original host header of the request is used.

* `vpc_channel_id` - (Optional, String) Specifies the VPC channel ID. This parameter and `url_domain` are
  alternative.

<a name="apigw_api_mock_policy"></a>
The `mock_policy` block supports:

* `effective_mode` - (Optional, String) Specifies the effective mode of the backend policy.
  The valid values are `ALL` and `ANY`, defaults to `ANY`.

* `response` - (Optional, String) Specifies the response of the backend policy.
  The description contains a maximum of `2,048` characters and the angle brackets (< and >) are not allowed.

* `name` - (Required, String) Specifies the backend policy name.
  The valid length is limited from can contain `3` to `64`, only letters, digits and underscores (_) are allowed.

* `conditions` - (Required, List) Specifies an array of one or more policy conditions.
  Up to five conditions can be set.
  The [object](#apigw_api_conditions) structure is documented below.

* `backend_params` - (Optional, List) Specifies an array of one or more backend parameters.
  The maximum of request parameters is `50`.
  The [object](#apigw_api_backend_params) structure is documented above.

* `authorizer_id` - (Optional, String) Specifies the ID of the backend custom authorization.

<a name="apigw_api_func_graph_policy"></a>
The `func_graph_policy` block supports:

* `function_urn` - (Required, String) Specifies the URN of the FunctionGraph function.

* `invocation_mode` - (Optional, String) Specifies the invocation mode of the FunctionGraph function.
  The valid values are `async` and `sync`, defaults to `sync`.

* `network_type` - (Optional, String) Function network architecture.
  The valid values are `VPC` and `NON-VPC`, defaults to `NON-VPC`.

* `version` - (Optional, String) Specifies the version of the FunctionGraph function.

* `timeout` - (Optional, Int) Specifies the timeout for API requests to backend service, the unit is `ms`.
  The valid value ranges from `1` to `600,000`, defaults to `5,000`.

* `effective_mode` - (Optional, String) Specifies the effective mode of the backend policy.
  The valid values are `ALL` and `ANY`, defaults to `ANY`.

* `name` - (Required, String) Specifies the backend policy name.
  The valid length is limited from can contain `3` to `64`, only letters, digits and underscores (_) are allowed.

* `conditions` - (Required, List) Specifies an array of one or more policy conditions.
  Up to five conditions can be set.
  The [object](#apigw_api_conditions) structure is documented below.

* `backend_params` - (Optional, List) Specifies the configuration list of the backend parameters.
  The maximum of request parameters is `50`.
  The [object](#apigw_api_backend_params) structure is documented above.

* `authorizer_id` - (Optional, String) Specifies the ID of the backend custom authorization.

<a name="apigw_api_http_policy"></a>
The `http_policy` block supports:

* `url_domain` - (Optional, String) Specifies the backend service address.
  The value which consists of a domain name or IP address, and a port number, with not more than `255` characters.
  The backend service address must be in the format "{host name}:{Port number}", for example, `apig.example.com:7443`.
  If the port number is not specified, the default HTTPS port `443`, or the default HTTP port `80` is used.
  The backend service address can contain environment variables, each starting with a letter and consisting of `3` to
  `32` characters. Only letters, digits, hyphens (-), and underscores (_) are allowed.

* `request_protocol` - (Optional, String) Specifies the backend request protocol. The valid values are `HTTP` and
  `HTTPS`, defaults to `HTTPS`.

* `request_method` - (Optional, String) Specifies the backend request method of the API.
  The valid values are `GET`, `POST`, `PUT`, `DELETE`, `HEAD`, `PATCH`, `OPTIONS` and `ANY`.

* `request_uri` - (Required, String) Specifies the backend request address, which can contain a maximum of `512` characters and
  must comply with URI specifications.
  + The address can contain request parameters enclosed with brackets ({}).
  + The address can contain special characters, such as asterisks (*), percent signs (%), hyphens (-) and
    underscores (_) and must comply with URI specifications.
  + The address can contain environment variables, each starting with a letter and consisting of `3` to `32` characters.
    Only letters, digits, hyphens (-), and underscores (_) are allowed in environment variables.

* `timeout` - (Optional, Int) Specifies the timeout, in ms, which allowed for APIGW to request the backend service. The
  valid value is range from `1` to `600,000`, defaults to `5,000`.

* `retry_count` - (Optional, Int) Specifies the number of retry attempts to request the backend service.
  The valid value ranges from `-1` to `10`, defaults to `-1`.
  `-1` indicates that idempotent APIs will retry once and non-idempotent APIs will not retry.
  `POST` and `PATCH` are not-idempotent.
  `GET`, `HEAD`, `PUT`, `OPTIONS` and `DELETE` are idempotent.

  -> When the (web) backend uses the channel, the `retry_count` must be less than the number of available backend
  servers in the channel.

* `effective_mode` - (Optional, String) Specifies the effective mode of the backend policy. The valid values are `ALL`
  and `ANY`, defaults to `ANY`.

* `name` - (Required, String) Specifies the backend policy name.
  The valid length is limited from can contain `3` to `64`, only letters, digits and underscores (_) are allowed.

* `backend_params` - (Optional, List) Specifies an array of one or more backend parameters. The maximum of request
  parameters is 50. The [object](#apigw_api_backend_params) structure is documented above.

* `conditions` - (Required, List) Specifies an array of one or more policy conditions.
  Up to five conditions can be set.
  The [object](#apigw_api_conditions) structure is documented below.

* `vpc_channel_proxy_host` - (Optional, String) Specifies the proxy host header.
  The host header can be customized for requests to be forwarded to cloud servers through the VPC channel.
  By default, the original host header of the request is used.

* `vpc_channel_id` - (Optional, String) Specifies the VPC channel ID.
  This parameter and `url_domain` are alternative.

* `authorizer_id` - (Optional, String) Specifies the ID of the backend custom authorization.

<a name="apigw_api_conditions"></a>
The `conditions` block supports:

* `param_name` - (Optional, String) Specifies the request parameter name.
  This parameter is required if the policy type is `param`. The valid values are `user_age` and `X-TEST-ENUM`.

* `type` - (Optional, String) Specifies the condition type of the backend policy.
  The valid values are `EXACT`, `ENUM` and `PATTERN`, defaults to `EXACT`.

* `origin` - (Optional, String) Specifies the backend policy type.
  The valid values are `param`, `source`, defaults to `source`.

* `value` - (Required, String) Specifies the value of the backend policy.
  For a condition with the input parameter source:
  + If the condition type is `ENUM`, separate condition values with commas.
  + If the condition type is `PATTERN`, enter a regular expression compatible with PERL.

  For a condition with the Source IP address source, enter IPv4 addresses and separate them with commas. The CIDR
  address format is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API.

* `registered_at` - Time when the API is registered.

* `updated_at` - Time when the API was last modified.

## Import

API can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_apigw_api_v2.api "774438a28a574ac8a496325d1bf51807"
```
