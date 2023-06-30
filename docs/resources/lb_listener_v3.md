---
subcategory: "Dedicated Load Balancer (DLB)"
---

Up-to-date reference of API arguments for DLB listener you can get at
`https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/listener`.

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

* `name` - (Optional) Specifies the listener name.

* `description` - (Optional) Provides supplementary information about the listener.

* `client_ca_tls_container_ref` - (Optional) Specifies the ID of the CA certificate used by the listener.

* `default_pool_id` - (Optional) Specifies the ID of the default backend server group. If there is no
  matched forwarding policy, requests are forwarded to the default backend server for processing.

* `default_tls_container_ref` - (Optional) Specifies the ID of the server certificate used by the listener.

* `http2_enable` - (Optional) Specifies whether to use HTTP/2. This parameter is available only for `HTTPS`
  listeners. If you configure this parameter for other types of listeners, it will not take effect. Enable
  HTTP/2 if you want the clients to use HTTP/2 to communicate with the load balancer.
  However, connections between the load balancer and backend servers use HTTP/1.x by default.

* `insert_headers` - (Optional) Specifies the HTTP header fields.
  * `forward_elb_ip` - (Optional) Specifies whether to transparently transmit the load balancer EIP
  to backend servers. If `forward_elb_ip` is set to `true`, the load balancer EIP will be stored in
  the HTTP header and passed to backend servers.
  * `forwarded_port` - (Optional) Specifies whether to transparently transmit the listening port of
  the load balancer to backend servers. If `forwarded_port` is set to `true`, the listening port of
  the load balancer will be stored in the HTTP header and passed to backend servers.
  * `forwarded_for_port` - (Optional) Specifies whether to transparently transmit the source port of
  the client to backend servers. If `forwarded_for_port` is set to `true`, the source port of the
  client will be stored in the HTTP header and passed to backend servers.
  * `forwarded_host` - (Optional) Specifies whether to rewrite the `X-Forwarded-Host` header.
  If `forwarded_host` is set to `true`, `X-Forwarded-Host` in the request header from the clients
  can be set to Host in the request header sent from the load balancer to backend servers.

* `loadbalancer_id` - (Required) Specifies the ID of the load balancer that the listener is added to.

* `protocol` - (Required) The protocol - can either be `TCP`, `HTTP`, `HTTPS` or `UDP`.
  Changing this creates a new Listener.

* `protocol_port` - (Required) Specifies the port used by the listener. Changing this creates a new Listener.

* `sni_container_refs` - (Optional) Lists the IDs of SNI certificates (server certificates with domain names) used by the listener.
  Each SNI certificate can have up to 30 domain names, and each domain name in the SNI certificate must be unique.
  This parameter will be ignored and an empty array will be returned if the listener's protocol is not `HTTPS`.

* `tls_ciphers_policy` - (Optional) Specifies the security policy that will be used by the listener.
  This parameter is available only for `HTTPS` listeners. An error will be returned if the protocol
  of the listener is not `HTTPS`. Possible values are: `tls-1-0`, `tls-1-1`, `tls-1-2`, `tls-1-2-strict`,
  `tls-1-2-fs`, `tls-1-0-with-1-3`, `tls-1-2-fs-with-1-3`.

* `member_retry_enable` - (Optional) Specifies whether to enable health check retries for backend servers.
  This parameter is available only for `HTTP` and `HTTPS` listeners. An error will be returned if you configure
  this parameter for `TCP` and `UDP` listeners.

* `keep_alive_timeout` - (Optional) Specifies the idle timeout duration, in seconds.
  * For `TCP` listeners, the value ranges from `10` to `4000`, and the default value is `300`.
  * For `HTTP` and `HTTPS` listeners, the value ranges from `0` to `4000`, and the default value is `60`.
  * For `UDP` listeners, this parameter is not available. An error will be returned if you
  configure this parameter for `UDP` listeners.

* `client_timeout` - (Optional) Specifies the timeout duration for waiting for a request from a client, in seconds.
  This parameter is available only for `HTTP` and `HTTPS` listeners. The value ranges from `1` to `300`, and
  the default value is `60`. An error will be returned if you configure this parameter for `TCP` and `UDP` listeners.

* `member_timeout` - (Optional) Specifies the timeout duration for waiting for a request from a
  backend server, in seconds. This parameter is available only for `HTTP` and `HTTPS` listeners.
  The value ranges from `1` to `300`, and the default value is `60`. An error will be returned if
  you configure this parameter for `TCP` and `UDP` listeners.

* `tags` - (Optional) Tags key/value pairs to associate with the loadbalancer listener.

* `advanced_forwarding` - (Optional) Specifies whether to enable advanced forwarding.
  If advanced forwarding is enabled, more flexible forwarding policies and rules are supported.
  The value can be `true` (enable advanced forwarding) or `false` (disable advanced forwarding),
  and the default value is `false`. Changing this creates a new Listener.

* `sni_match_algo` - (Optional) Specifies how wildcard domain name matches with the SNI certificates
  used by the listener.

* `security_policy_id` - (Optional) Specifies the ID of the custom security policy.

* `ip_group` - (Optional) Specifies the IP address group associated with the listener.
  * `id` - (Required) Specifies the ID of the IP address group associated with the listener.
  * `enable` - (Optional) Specifies whether to enable access control.
    `true` (default): Access control will be enabled.
    `false`: Access control will be disabled.
  * `type` - (Optional) Specifies how access to the listener is controlled.
    `white` (default): A whitelist will be configured. Only IP addresses in the whitelist can access the listener.
    `black`: A blacklist will be configured. IP addresses in the blacklist are not allowed to access the listener.

## Attributes Reference

In addition, the following attributes are exported:

* `updated_at` - Indicates the update time.

* `createa_at` - Indicates the creation time.

## Import

Listeners can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_lb_listener_v3.listener_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
