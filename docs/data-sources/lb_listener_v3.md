---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_listener_v3

Use this data source to get the info about an existing ELBv3 listener.

## Example Usage

```hcl
data "opentelekomcloud_lb_listener_v3" "listener" {
  loadbalancer_id = var.loadbalancer_id
  name            = "https_listener"
}
```

## Argument Reference

The following arguments are supported:
* `id` - (Optional) Specifies the listener ID.

* `name` - (Optional) Specifies the listener name.

* `description` - (Optional) Provides supplementary information about the listener.

* `client_ca_tls_container_ref` - (Optional) Specifies the ID of the CA certificate used by the listener.

* `default_pool_id` - (Optional) Specifies the ID of the default backend server group.

* `default_tls_container_ref` - (Optional) Specifies the ID of the server certificate used by the listener.

* `loadbalancer_id` - (Optional) Specifies the ID of the load balancer that the listener is added to.

* `protocol` - (Optional) The protocol - can either be `TCP`, `HTTP`, `HTTPS` or `UDP`.

* `protocol_port` - (Optional) Specifies the port used by the listener. Changing this creates a new Listener.

* `tls_ciphers_policy`- (Optional) Specifies the TLS version used.

* `keep_alive_timeout` - (Optional) Specifies the idle timeout duration, in seconds.

* `client_timeout` - (Optional) Specifies the timeout duration for waiting for a request from a client, in seconds.
  This parameter is available only for `HTTP` and `HTTPS` listeners. The value ranges from `1` to `300`, and
  the default value is `60`. An error will be returned if you configure this parameter for `TCP` and `UDP` listeners.

* `member_timeout` - (Optional) Specifies the timeout duration for waiting for a request from a
  backend server, in seconds. This parameter is available only for `HTTP` and `HTTPS` listeners.
  The value ranges from `1` to `300`, and the default value is `60`. An error will be returned if
  you configure this parameter for `TCP` and `UDP` listeners.

* `member_address` - (Optional) Specifies the private IP address bound to the backend server.
  This parameter is used only as a query condition and is not included in the response.

* `member_device_id` - (Optional) Specifies the ID of the cloud server that serves as a backend server.
  This parameter is used only as a query condition and is not included in the response.

## Attributes Reference

In addition, the following attributes are exported:

* `insert_headers` - Specifies the HTTP header fields.
    * `forward_elb_ip` - Specifies whether to transparently transmit the load balancer EIP
      to backend servers. If `forward_elb_ip` is set to `true`, the load balancer EIP will be stored in
      the HTTP header and passed to backend servers.
    * `forwarded_port` - Specifies whether to transparently transmit the listening port of
      the load balancer to backend servers. If `forwarded_port` is set to `true`, the listening port of
      the load balancer will be stored in the HTTP header and passed to backend servers.
    * `forwarded_for_port` - Specifies whether to transparently transmit the source port of
      the client to backend servers. If `forwarded_for_port` is set to `true`, the source port of the
      client will be stored in the HTTP header and passed to backend servers.
    * `forwarded_host` - Specifies whether to rewrite the `X-Forwarded-Host` header.
      If `forwarded_host` is set to `true`, `X-Forwarded-Host` in the request header from the clients
      can be set to Host in the request header sent from the load balancer to backend servers.

* `project_id` - Specifies the project ID.

* `member_retry_enable` - Specifies whether to enable health check retries for backend servers.

* `sni_container_refs` - Lists the IDs of SNI certificates (server certificates with domain names) used by the listener.

* `advanced_forwarding` - Specifies whether to enable advanced forwarding.

* `sni_match_algo` - Specifies how wildcard domain name matches with the SNI certificates
  used by the listener.

* `security_policy_id` - Specifies the ID of the custom security policy.

* `ip_group` - Specifies the IP address group associated with the listener.

* `tags` - Tags key/value pairs to associate with the loadbalancer listener.

* `http2_enable` - Specifies whether to use HTTP/2.

* `updated_at` - Indicates the update time.

* `created_at` - Indicates the creation time.
