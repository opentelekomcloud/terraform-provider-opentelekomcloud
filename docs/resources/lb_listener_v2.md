---
subcategory: "Elastic Load Balance (ELB)"
---

# opentelekomcloud_lb_listener_v2

Manages an Enhanced LB listener resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_listener_v2" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
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

* `default_tls_container_ref` - (Optional) Specifies the ID of the server certificate used by the listener.
  The value contains a maximum of 128 characters. The default value is `null`.
  This parameter is **required** when protocol is set to `TERMINATED_HTTPS`.
  See [here](https://wiki.openstack.org/wiki/Network/LBaaS/docs/how-to-create-tls-loadbalancer)
  for more information.

* `sni_container_refs` - (Optional) Lists the IDs of SNI certificates (server certificates with a domain name) used
  by the listener. If the parameter value is an empty list, the SNI feature is disabled.
  The default value is `[]`. This is **required** if the protocol is `TERMINATED_HTTPS`.

* `admin_state_up` - (Optional) The administrative state of the Listener.
  A valid value is `true` (UP) or `false` (DOWN).

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the Listener.

* `protocol` - See Argument Reference above.

* `protocol_port` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `name` - See Argument Reference above.

* `default_port_id` - See Argument Reference above.

* `description` - See Argument Reference above.

* `default_tls_container_ref` - See Argument Reference above.

* `sni_container_refs` - See Argument Reference above.

* `admin_state_up` - See Argument Reference above.
