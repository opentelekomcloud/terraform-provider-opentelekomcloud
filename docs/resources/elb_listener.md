---
subcategory: "Classic Load Balancer (LB)"
---

# opentelekomcloud_elb_listener

Manages a classic loadbalancer listener resource within OpenTelekomCloud.

~>
Classic load balancers are no longer provided. Please use elastic load balancers instead.
This resource was replaced with `opentelekomcloud_lb_listener_v2`.

## Example Usage

```hcl
resource "opentelekomcloud_elb_loadbalancer" "elb" {
  name           = "elb"
  type           = "External"
  description    = "test elb"
  vpc_id         = "e346dc4a-d9a6-46f4-90df-10153626076e"
  admin_state_up = "true"
  bandwidth      = 5
}

resource "opentelekomcloud_elb_listener" "listener" {
  name             = "test-elb-listener"
  description      = "great listener"
  protocol         = "TCP"
  backend_protocol = "TCP"
  protocol_port    = 12345
  backend_port     = 8080
  lb_algorithm     = "roundrobin"
  loadbalancer_id  = opentelekomcloud_elb_loadbalancer.elb.id
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the load balancer name. The name is a string
  of 1 to 64 characters that consist of letters, digits, underscores (_), and
  hyphens (-).

* `description` - (Optional) Provides supplementary information about the listener.
  The value is a string of 0 to 128 characters and cannot be <>.

* `loadbalancer_id` - (Required) Specifies the ID of the load balancer to which
  the listener belongs.  Changing this creates a new elb listener.

* `protocol` - (Required) Specifies the listening protocol used for layer 4
  or 7. The value can be HTTP, TCP, HTTPS, or UDP.  Changing this creates a
  new elb listener.

* `protocol_port` - (Required) Specifies the listening port. The value ranges from 1
  to 65535.

* `backend_protocol` - (Required) Specifies the backend protocol. If the value
  of protocol is UDP, the value of this parameter can only be UDP. The value can
  be HTTP, TCP, or UDP.  Changing this creates a new elb listener.

* `backend_port` - (Required) Specifies the backend port. The value ranges from
  1 to 65535.

* `lb_algorithm` - (Required) Specifies the load balancing algorithm for the
  listener. The value can be roundrobin, leastconn, or source.

* `session_sticky` - (Optional) Specifies whether to enable sticky session.
  The value can be true or false. The Sticky session is enabled when the value
  is true, and is disabled when the value is false. If the value of protocol is
  HTTP, HTTPS, or TCP, and the value of lb_algorithm is not roundrobin, the value
  of this parameter can only be false.  Changing this creates a new elb listener.

* `session_sticky_type` - (Optional) Specifies the cookie processing method.
  The value is insert. insert indicates that the cookie is inserted by the load
  balancer. This parameter is valid when protocol is set to HTTP, and session_sticky
  to true. The default value is insert. This parameter is invalid when protocol
  is set to TCP or UDP, which means the parameter is empty.  Changing this creates
  a new elb listener.

* `cookie_timeout` - (Optional) Specifies the cookie timeout period (minutes).
  This parameter is valid when protocol is set to HTTP, session_sticky to true,
  and session_sticky_type to insert. This parameter is invalid when protocol is
  set to TCP or UDP. The value ranges from 1 to 1440.  Changing this creates a
  new elb listener.

* `tcp_timeout` - (Optional) Specifies the TCP timeout period (minutes). This
  parameter is valid when protocol is set to TCP. The value ranges from 1 to 5.

* `tcp_draining` - (Optional) Specifies whether to maintain the TCP connection
  to the backend ECS after the ECS is deleted. This parameter is valid when protocol
  is set to TCP. The value can be true or false.

* `tcp_draining_timeout` - (Optional) Specifies the timeout duration (minutes)
  for the TCP connection to the backend ECS after the ECS is deleted. This parameter
  is valid when protocol is set to TCP, and tcp_draining to true. The value ranges
  from 0 to 60.

* `certificate_id` - (Optional) Specifies the ID of the SSL certificate used
  for security authentication when HTTPS is used to make API calls. This parameter
  is mandatory if the value of protocol is HTTPS. The value can be obtained by
  viewing the details of the SSL certificate.  Changing this creates a new elb
  listener.

* `udp_timeout` - (Optional) Specifies the UDP timeout duration (minutes). This
  parameter is valid when protocol is set to UDP. The value ranges from 1 to 1440.

* `ssl_protocols` - (Optional) Specifies the SSL protocol standard supported
  by a tracker, which is used for enabling specified encryption protocols. This
  parameter is valid only when the value of protocol is set to HTTPS. The value
  is TLSv1.2 or TLSv1.2 TLSv1.1 TLSv1. The default value is TLSv1.2. Changing
  this creates a new elb listener.

* `ssl_ciphers` - (Optional) Specifies the cipher suite of an encryption protocol.
  This parameter is valid only when the value of protocol is set to HTTPS. The
  value is Default, Extended, or Strict. The default value is Default. The value
  can only be set to Extended if the value of ssl_protocols is set to TLSv1.2
  TLSv1.1 TLSv1.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `description` - See Argument Reference above.

* `loadbalancer_id` - See Argument Reference above.

* `protocol` - See Argument Reference above.

* `protocol_port` - See Argument Reference above.

* `backend_protocol` - See Argument Reference above.

* `backend_port` - See Argument Reference above.

* `lb_algorithm` - See Argument Reference above.

* `session_sticky` - See Argument Reference above.

* `session_sticky_type` - See Argument Reference above.

* `cookie_timeout` - See Argument Reference above.

* `tcp_timeout` - See Argument Reference above.

* `tcp_draining` - See Argument Reference above.

* `tcp_draining_timeout` - See Argument Reference above.

* `certificate_id` - See Argument Reference above.

* `udp_timeout` - See Argument Reference above.

* `ssl_protocols` - See Argument Reference above.

* `ssl_ciphers` - See Argument Reference above.

* `id` - Specifies the listener ID.

* `admin_state_up` - Specifies the status of the load balancer. Value range:
  * `false`: The load balancer is disabled.
  * `true`: The load balancer runs properly.
