---
subcategory: "Elastic Load Balancer (ELB)"
---

# opentelekomcloud_lb_listener_v2

Manages an Enhanced LB listener resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  tags = {
    muh = "kuh"
  }
}
```

## Argument Reference

The following arguments are supported:

* `protocol` - (Required) The protocol - can either be `TCP`, `HTTP`, `HTTPS` or `TERMINATED_HTTPS`.
  Changing this creates a new Listener.

* `protocol_port` - (Required) The port on which to listen for client traffic.
  Changing this creates a new Listener.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
  the Listener.  Only administrative users can specify a tenant UUID
  other than their own. Changing this creates a new Listener.

* `loadbalancer_id` - (Required) The load balancer on which to provision this
  Listener. Changing this creates a new Listener.

* `name` - (Optional) Human-readable name for the Listener. Does not have
  to be unique.

* `default_pool_id` - (Optional) The ID of the default pool with which the
  Listener is associated. Changing this creates a new Listener.

* `description` - (Optional) Human-readable description for the Listener.

* `http2_enable`- (Optional) `true` to enable HTTP/2 mode of ELB.
  HTTP/2 is disabled by default if not set.

* `default_tls_container_ref` - (Optional) Specifies the ID of a certificate container of type `server`
  used by the listener. The value contains a maximum of 128 characters. The default value is `null`.
  This parameter is **required** when protocol is set to `TERMINATED_HTTPS`.
  See [here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
  for more information.

* `client_ca_tls_container_ref`  (Optional) Specifies the ID of a certificate container of type `client`
  used by the listener. The value contains a maximum of 128 characters. The default value is `null`.
  The loadbalancer only establishes a TLS connection if the client presents a certificate delivered by
  the client CA whose certificate is registered in the referenced certificate container. The option is
  effective only in conjunction with `TERMINATED_HTTPS`.

* `sni_container_refs` - (Optional) Lists the IDs of SNI certificates (server certificates with a domain name) used
  by the listener. If the parameter value is an empty list, the SNI feature is disabled.
  The default value is `[]`. It only works in conjunction with `TERMINATED_HTTPS`.

* `tls_ciphers_policy`- (Optional) Controls the TLS version used. Supported values are `tls-1-0`, `tls-1-1`,
  `tls-1-2` and `tls-1-2-strict`. If not set, the loadbalancer uses `tls-1-0`. See
  [here](https://docs.otc.t-systems.com/api/elb/elb_zq_jt_0001.html) for details about the supported cipher
  suites. The option is effective only in conjunction with `TERMINATED_HTTPS`.

* `transparent_client_ip_enable` - (Optional) Specifies whether to pass source IP addresses of the clients to
  backend servers. The value is always `true` for `HTTP` and `HTTPS` listeners. For `TCP` and `UDP` listeners the
  value can be `true` or `false` with `false` by default.

->
  If the load balancer is a Dedicated Load Balancer, `transparent_client_ip_enable` is always `true`

* `admin_state_up` - (Optional) The administrative state of the Listener.
  A valid value is `true` (UP) or `false` (DOWN).

* `tags` - (Optional) Tags key/value pairs to associate with the loadbalancer listener.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the Listener.

* `protocol` - See Argument Reference above.

* `protocol_port` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `name` - See Argument Reference above.

* `default_port_id` - See Argument Reference above.

* `description` - See Argument Reference above.

* `http2_enable` - See Argument Reference above.

* `default_tls_container_ref` - See Argument Reference above.

* `client_ca_tls_container_ref` - See Argument Reference above.

* `sni_container_refs` - See Argument Reference above.

* `tls_ciphers_policy` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.

* `tags` - See Argument Reference above.

## Import

Listeners can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_lb_listener_v2.listener_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
