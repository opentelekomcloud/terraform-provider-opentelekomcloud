---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_listener_v3"
sidebar_current: "docs-opentelekomcloud-resource-lb-listener-v3"
description: |-
  Manages a LB Listener resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DLB listener you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/listener)

# opentelekomcloud_lb_listener_v3

Manages a Dedicated LB listener resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = var.loadbalancer_id

  tags = {
    muh = "kuh"
  }
}
```

## Example Ip Address Group

```hcl
resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_2"
  description = "some interesting description 2"

  ip_list {
    ip          = "192.168.10.11"
    description = "one"
  }
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name            = "listener_1"
  description     = "some interesting description"
  loadbalancer_id = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol        = "HTTP"
  protocol_port   = 8080

  advanced_forwarding = true
  sni_match_algo      = "wildcard"

  insert_headers {
    forwarded_host = true
  }

  ip_group {
    id     = opentelekomcloud_lb_ipgroup_v3.group_1.id
    enable = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the listener name.

* `description` - (Optional, String) Provides supplementary information about the listener.

* `client_ca_tls_container_ref` - (Optional, String) Specifies the ID of the CA certificate used by the listener.

* `default_pool_id` - (Optional, String) Specifies the ID of the default backend server group. If there is no
  matched forwarding policy, requests are forwarded to the default backend server for processing.

* `default_tls_container_ref` - (Optional, String) Specifies the ID of the server certificate used by the listener.

* `http2_enable` - (Optional, Bool) Specifies whether to use HTTP/2. This parameter is available only for `HTTPS`
  listeners. If you configure this parameter for other types of listeners, it will not take effect. Enable
  HTTP/2 if you want the clients to use HTTP/2 to communicate with the load balancer.
  However, connections between the load balancer and backend servers use HTTP/1.x by default.

* `insert_headers` - (Optional, List) Specifies the HTTP header fields.
  * `forward_elb_ip` - (Optional, Bool) Specifies whether to transparently transmit the load balancer EIP
  to backend servers. If `forward_elb_ip` is set to `true`, the load balancer EIP will be stored in
  the HTTP header and passed to backend servers.
  * `forwarded_port` - (Optional, Bool) Specifies whether to transparently transmit the listening port of
  the load balancer to backend servers. If `forwarded_port` is set to `true`, the listening port of
  the load balancer will be stored in the HTTP header and passed to backend servers.
  * `forwarded_for_port` - (Optional, Bool) Specifies whether to transparently transmit the source port of
  the client to backend servers. If `forwarded_for_port` is set to `true`, the source port of the
  client will be stored in the HTTP header and passed to backend servers.
  * `forwarded_host` - (Optional, Bool) Specifies whether to rewrite the `X-Forwarded-Host` header.
  If `forwarded_host` is set to `true`, `X-Forwarded-Host` in the request header from the clients
  can be set to Host in the request header sent from the load balancer to backend servers.

* `loadbalancer_id` - (Required, ForceNew, String) Specifies the ID of the load balancer that the listener is added to.

* `protocol` - (Required, ForceNew, String) The protocol - can either be `TCP`, `HTTP`, `HTTPS` or `UDP`.
  Changing this creates a new Listener.

* `protocol_port` - (Required, ForceNew, Int) Specifies the port used by the listener. Changing this creates a new Listener.

* `sni_container_refs` - (Optional, List) Lists the IDs of SNI certificates (server certificates with domain names) used by the listener.
  Each SNI certificate can have up to 30 domain names, and each domain name in the SNI certificate must be unique.
  This parameter will be ignored and an empty array will be returned if the listener's protocol is not `HTTPS`.

* `tls_ciphers_policy` - (Optional, String) Specifies the security policy that will be used by the listener.
  This parameter is available only for `HTTPS` listeners. An error will be returned if the protocol
  of the listener is not `HTTPS`. Possible values are: `tls-1-0`, `tls-1-1`, `tls-1-0-inherit`, `tls-1-2`,
  `tls-1-2-strict`, `tls-1-2-fs`, `tls-1-0-with-1-3`, `tls-1-2-fs-with-1-3`, `hybrid-policy-1-0`, `tls-1-2-strict-no-cbc`.

* `member_retry_enable` - (Optional, Bool) Specifies whether to enable health check retries for backend servers.
  This parameter is available only for `HTTP` and `HTTPS` listeners. An error will be returned if you configure
  this parameter for `TCP` and `UDP` listeners.

* `keep_alive_timeout` - (Optional, Int) Specifies the idle timeout duration, in seconds.
  * For `TCP` listeners, the value ranges from `10` to `4000`, and the default value is `300`.
  * For `HTTP` and `HTTPS` listeners, the value ranges from `0` to `4000`, and the default value is `60`.
  * For `UDP` listeners, this parameter is not available. An error will be returned if you
  configure this parameter for `UDP` listeners.

* `client_timeout` - (Optional, Int) Specifies the timeout duration for waiting for a request from a client, in seconds.
  This parameter is available only for `HTTP` and `HTTPS` listeners. The value ranges from `1` to `300`, and
  the default value is `60`. An error will be returned if you configure this parameter for `TCP` and `UDP` listeners.

* `member_timeout` - (Optional, Int) Specifies the timeout duration for waiting for a request from a
  backend server, in seconds. This parameter is available only for `HTTP` and `HTTPS` listeners.
  The value ranges from `1` to `300`, and the default value is `60`. An error will be returned if
  you configure this parameter for `TCP` and `UDP` listeners.

* `tags` - (Optional, ForceNew, Map) Tags key/value pairs to associate with the loadbalancer listener.

* `advanced_forwarding` - (Optional, ForceNew, Bool) Specifies whether to enable advanced forwarding.
  If advanced forwarding is enabled, more flexible forwarding policies and rules are supported.
  The value can be `true` (enable advanced forwarding) or `false` (disable advanced forwarding),
  and the default value is `false`. Changing this creates a new Listener.

* `sni_match_algo` - (Optional, String) Specifies how wildcard domain name matches with the SNI certificates
  used by the listener.

* `security_policy_id` - (Optional, String) Specifies the ID of the custom security policy.

* `ip_group` - (Optional, List) Specifies the IP address group associated with the listener.
  * `id` - (Required, String) Specifies the ID of the IP address group associated with the listener.
    Specifies the ID of the IP address group associated with the listener.
    If `ip_list` in `opentelekomcloud_lb_ipgroup_v3` is set to an empty array `[]` and type to `whitelist`, no IP addresses are allowed to access the listener.
    If `ip_list` in `opentelekomcloud_lb_ipgroup_v3` is set to an empty array `[]` and type to `blacklist`, any IP address is allowed to access the listener.
  * `enable` - (Optional, Bool) Specifies whether to enable access control.
    `true`: Access control will be enabled.
    `false` (default): Access control will be disabled.
  * `type` - (Optional, String) Specifies how access to the listener is controlled.
    `white` (default): A whitelist will be configured. Only IP addresses in the whitelist can access the listener.
    `black`: A blacklist will be configured. IP addresses in the blacklist are not allowed to access the listener.

* `protection_status` - (Optional, String) Specifies the protection status. Value options: `nonProtection`
  (default, not protected), `consoleProtection` (modification protection is enabled on the console).

* `protection_reason` - (Optional, String) Specifies why the modification protection is enabled. Valid only
  when `protection_status` is set to `consoleProtection`. The value can contain a maximum of 255 Unicode
  characters, excluding angle brackets (`<>`).

* `access_log_customized_headers_config` - (Optional, List) Specifies the custom headers to be recorded in
  access logs. You can specify which headers to be or not to be recorded in the access logs of a load balancer.
  * `enable` - (Optional, Bool) Specifies whether to enable custom access log headers.
  * `include_headers` - (Optional, List) Specifies the custom headers that will be recorded in access logs.
    If this parameter is specified, only the specified headers will be recorded in access logs.
  * `exclude_headers` - (Optional, List) Specifies the headers that will not be recorded in access logs.
    If this parameter is specified, the specified headers will not be recorded in access logs.

## Attributes Reference

In addition, the following attributes are exported:

* `updated_at` - Indicates the update time.

* `created_at` - Indicates the creation time.

* `quic_config` - Specifies the QUIC configuration for the current listener. Valid only when `protocol` is
  set to `HTTPS`.
  * `quic_listener_id` - Specifies the ID of the QUIC listener.
  * `enable_quic_upgrade` - Specifies whether QUIC upgrade is enabled.

* `gzip_enable` - Specifies whether GZIP compression is enabled for the load balancer.

* `cps` - Specifies the maximum number of new connections that a listener can handle per second.

* `max_connections` - Specifies the maximum number of concurrent connections that a listener can handle per second.

* `nat64_enable` - Specifies whether to translate between IPv4 and IPv6 addresses. This option allows a client
  to access IPv4 or IPv6 backend servers by accessing the IPv4 or IPv6 address of a load balancer.

* `proxy_protocol_enable` - Specifies whether ProxyProtocol is enabled to pass the source IP addresses of the
  clients to backend servers.

* `tracing_config` - Specifies the open tracing configuration.
  * `tracing_enable` - Specifies whether open tracing is enabled.
  * `tracing_strategy` - Specifies the sampling mode. Value options: `full` (full sampling), `ratio`
    (sampling by ratio).
  * `tracing_sample` - Specifies the sampling rate. Values between `1` and `10000` represent sampling rates
    ranging from 0.01% to 100%.
  * `tracing_type` - Specifies the tracing type. Value: `W3CTraceContext`.

## Import

Listeners can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_lb_listener_v3.listener_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
