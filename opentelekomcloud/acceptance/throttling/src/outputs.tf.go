package src

const Outputs = `
output "throttle-url" {
  value = ["https://${var.rancher_host}.${var.rancher_domain}", "https://${opentelekomcloud_networking_floatingip_v2.eip.address}"]
}

output "wireguard-server-ip" {
  value = var.deploy_wireguard ? opentelekomcloud_networking_floatingip_v2.wireguard[0].address : null
}

output "wireguard-server-port" {
  value = var.deploy_wireguard ? var.wg_server_port : null
}

output "wireguard-server-key" {
  value = var.deploy_wireguard ? var.wg_server_public_key : null
}

output "throttle-nodes" {
  value = [opentelekomcloud_compute_instance_v2.throttle-server-1.access_ip_v4, opentelekomcloud_compute_instance_v2.throttle-server-2.access_ip_v4]
}

output "bootstrap_password" {
  value = var.admin_password
}

output "vpc_id" {
  value = opentelekomcloud_vpc_v1.vpc.id
}
`
