---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_pods_v2"
sidebar_current: "docs-opentelekomcloud-data-source-cci-pods-v2"
description: |-
  Get the list of CCI v2 Pods within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI Pod you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_pods_v2

Use this data source to get the list of CCI v2 Pods under a namespace within OpenTelekomCloud.

## Example Usage

```hcl
variable "namespace" {}

data "opentelekomcloud_cci_pods_v2" "test" {
  namespace = var.namespace
}
```

### Query a Single Pod by Name

```hcl
variable "namespace" {}
variable "pod_name" {}

data "opentelekomcloud_cci_pods_v2" "test" {
  namespace = var.namespace
  name      = var.pod_name
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) Specifies the namespace to which the Pods belong.

* `name` - (Optional, String) Specifies the name of the Pod used to query a single Pod.
  If omitted, all Pods under the namespace are returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the Pods are queried.

* `pods` - The list of Pods. The [pods](#pods) structure is documented below.

<a name="pods"></a>
The `pods` block supports:

* `name` - The name of the Pod.

* `namespace` - The namespace to which the Pod belongs.

* `api_version` - The API version of the Pod.

* `kind` - The kind of the Pod.

* `annotations` - The annotations of the Pod.

* `labels` - The labels of the Pod.

* `creation_timestamp` - The creation timestamp of the Pod.

* `resource_version` - The resource version of the Pod.

* `uid` - The UID of the Pod.

* `finalizers` - The finalizers of the Pod.

* `active_deadline_seconds` - The active deadline of the Pod, in seconds.

* `affinity` - The affinity scheduling rules of the Pod. The [affinity](#affinity) structure is documented below.

* `containers` - The list of containers in the Pod. The [containers](#containers) structure is documented below.

* `dns_config` - The DNS configuration of the Pod. The [dns_config](#dns_config) structure is documented below.

* `dns_policy` - The DNS policy of the Pod.

* `ephemeral_containers` - The list of ephemeral containers in the Pod.
  The [containers](#containers) structure is documented below.

* `host_aliases` - The list of hosts and IPs injected into the Pod's `/etc/hosts`.
  The [host_aliases](#host_aliases) structure is documented below.

* `hostname` - The hostname of the Pod.

* `image_pull_secrets` - The list of references to secrets used for pulling images.
  The [image_pull_secrets](#image_pull_secrets) structure is documented below.

* `init_containers` - The list of init containers of the Pod.
  The [containers](#containers) structure is documented below.

* `node_name` - The name of the node the Pod is scheduled on.

* `overhead` - The resource overhead associated with running the Pod.

* `readiness_gates` - The list of additional pod readiness conditions.
  The [readiness_gates](#readiness_gates) structure is documented below.

* `restart_policy` - The restart policy of the Pod.

* `scheduler_name` - The name of the scheduler used to dispatch the Pod.

* `security_context` - The pod-level security attributes.
  The [security_context](#security_context) structure is documented below.

* `set_hostname_as_fqdn` - Whether the Pod's hostname is set to its FQDN.

* `share_process_namespace` - Whether a single process namespace is shared between all containers.

* `termination_grace_period_seconds` - The grace period in seconds before the Pod is terminated.

* `volumes` - The list of volumes used by the Pod. The [volumes](#volumes) structure is documented below.

* `status` - The status of the Pod. The [status](#status) structure is documented below.

<a name="affinity"></a>
The `affinity` block supports:

* `node_affinity` - The node affinity rules. The [node_affinity](#node_affinity) structure is documented below.

* `pod_anti_affinity` - The pod anti-affinity rules.
  The [pod_anti_affinity](#pod_anti_affinity) structure is documented below.

<a name="node_affinity"></a>
The `node_affinity` block supports:

* `required_during_scheduling_ignored_during_execution` - Node affinity requirements that must be met.
  The [node_selector](#node_selector) structure is documented below.

<a name="node_selector"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `node_selector_terms` - The list of node selector terms.
  The [node_selector_terms](#node_selector_terms) structure is documented below.

<a name="node_selector_terms"></a>
The `node_selector_terms` block supports:

* `match_expressions` - The list of label selector requirements.
  The [match_expressions](#match_expressions) structure is documented below.

<a name="match_expressions"></a>
The `match_expressions` block supports:

* `key` - The label key the selector applies to.

* `operator` - The operator of the selector.

* `values` - The list of values for the selector.

<a name="pod_anti_affinity"></a>
The `pod_anti_affinity` block supports:

* `preferred_during_scheduling_ignored_during_execution` - The soft anti-affinity preferences.
  The [weighted_pod_affinity_term](#weighted_pod_affinity_term) structure is documented below.

* `required_during_scheduling_ignored_during_execution` - The hard anti-affinity requirements.
  The [pod_affinity_term](#pod_affinity_term) structure is documented below.

<a name="weighted_pod_affinity_term"></a>
The `preferred_during_scheduling_ignored_during_execution` block supports:

* `weight` - The weight associated with matching the corresponding term.

* `pod_affinity_term` - The pod affinity term associated with this weight.
  The [pod_affinity_term](#pod_affinity_term) structure is documented below.

<a name="pod_affinity_term"></a>
The `pod_affinity_term` block supports:

* `label_selector` - The label query over a set of resources.
  The [label_selector](#label_selector) structure is documented below.

* `namespaces` - The list of namespaces the label selector applies to.

* `topology_key` - The key of the node label used to define a topology domain.

<a name="label_selector"></a>
The `label_selector` block supports:

* `match_labels` - The map of `{key,value}` pairs.

* `match_expressions` - The list of label selector requirements.
  The [match_expressions](#match_expressions) structure is documented below.

<a name="containers"></a>
The `containers` block supports:

* `name` - The name of the container.

* `image` - The container image name.

* `args` - The arguments to the entrypoint.

* `command` - The entrypoint array.

* `env` - The list of environment variables. The [env](#env) structure is documented below.

* `env_from` - The list of sources to populate environment variables from.
  The [env_from](#env_from) structure is documented below.

* `lifecycle` - Actions that the management system should take in response to container lifecycle events.
  The [lifecycle](#lifecycle) structure is documented below.

* `liveness_probe` - The liveness probe configuration. The [probe](#probe) structure is documented below.

* `readiness_probe` - The readiness probe configuration. The [probe](#probe) structure is documented below.

* `startup_probe` - The startup probe configuration. The [probe](#probe) structure is documented below.

* `ports` - The list of ports exposed by the container. The [ports](#ports) structure is documented below.

* `resources` - The compute resource requirements of the container.
  The [resources](#resources) structure is documented below.

* `security_context` - The security context of the container.
  The [container_security_context](#container_security_context) structure is documented below.

* `stdin` - Whether the container should allocate a buffer for stdin.

* `stdin_once` - Whether the stdin channel is closed after the first attach disconnects.

* `target_container_name` - The name of the target container for ephemeral containers.

* `termination_message_path` - The path at which the file for the container's termination message is written.

* `termination_message_policy` - The policy for determining the container's termination message.

* `tty` - Whether the container should be allocated a TTY.

* `working_dir` - The container's working directory.

* `volume_mounts` - The list of volume mounts within the container.
  The [volume_mounts](#volume_mounts) structure is documented below.

<a name="env"></a>
The `env` block supports:

* `name` - The name of the environment variable.

* `value` - The value of the environment variable.

<a name="env_from"></a>
The `env_from` block supports:

* `config_map_ref` - A reference to a ConfigMap from which to load environment variables.
  The [env_source](#env_source) structure is documented below.

* `prefix` - The prefix prepended to each key.

* `secret_ref` - A reference to a Secret from which to load environment variables.
  The [env_source](#env_source) structure is documented below.

<a name="env_source"></a>
The `config_map_ref` / `secret_ref` block supports:

* `name` - The name of the referenced ConfigMap or Secret.

* `optional` - Whether the ConfigMap or Secret must be defined.

<a name="lifecycle"></a>
The `lifecycle` block supports:

* `post_start` - The handler executed right after container is created.
  The [lifecycle_handler](#lifecycle_handler) structure is documented below.

* `pre_stop` - The handler executed right before a container is terminated.
  The [lifecycle_handler](#lifecycle_handler) structure is documented below.

<a name="lifecycle_handler"></a>
The `post_start` / `pre_stop` block supports:

* `exec` - The exec action to execute. The [exec](#exec) structure is documented below.

* `http_get` - The HTTP GET action to perform.
  The [http_get](#http_get) structure is documented below.

<a name="exec"></a>
The `exec` block supports:

* `command` - The command line to execute inside the container.

<a name="http_get"></a>
The `http_get` block supports:

* `host` - The host name to connect to.

* `http_headers` - The list of custom HTTP headers.
  The [http_headers](#http_headers) structure is documented below.

* `path` - The path to access on the HTTP server.

* `port` - The port to access on the container.

* `scheme` - The scheme to use for connecting to the host.

<a name="http_headers"></a>
The `http_headers` block supports:

* `name` - The header name.

* `value` - The header value.

<a name="probe"></a>
The `liveness_probe` / `readiness_probe` / `startup_probe` block supports:

* `exec` - The exec action to perform. The [exec](#exec) structure is documented below.

* `failure_threshold` - Minimum consecutive failures for the probe to be considered failed.

* `http_get` - The HTTP GET action to perform. The [http_get](#http_get) structure is documented below.

* `initial_delay_seconds` - Number of seconds before the first probe is initiated.

* `period_seconds` - How often to perform the probe, in seconds.

* `success_threshold` - Minimum consecutive successes for the probe to be considered successful.

* `termination_grace_period_seconds` - The grace period in seconds when the probe fails.

<a name="ports"></a>
The `ports` block supports:

* `container_port` - The port number exposed by the container.

* `name` - The name of the port.

* `protocol` - The protocol of the port.

<a name="resources"></a>
The `resources` block supports:

* `limits` - The maximum amount of compute resources allowed.

* `requests` - The minimum amount of compute resources required.

<a name="container_security_context"></a>
The container `security_context` block supports:

* `capabilities` - The list of capabilities to add or drop.
  The [capabilities](#capabilities) structure is documented below.

* `proc_mount` - The type of proc mount to use for the container.

* `read_only_root_file_system` - Whether the container has a read-only root filesystem.

* `run_as_group` - The GID to run the entrypoint of the container process.

* `run_as_non_root` - Whether the container must run as a non-root user.

* `run_as_user` - The UID to run the entrypoint of the container process.

<a name="capabilities"></a>
The `capabilities` block supports:

* `add` - The list of capabilities to add.

* `drop` - The list of capabilities to drop.

<a name="volume_mounts"></a>
The `volume_mounts` block supports:

* `name` - The name of the volume.

* `mount_path` - The path within the container at which the volume is mounted.

* `read_only` - Whether the mount is read-only.

* `sub_path` - A sub-path within the volume that should be mounted.

* `sub_path_expr` - An expanded sub-path within the volume.

* `extend_path_mode` - The extended path mode of the mount.

<a name="dns_config"></a>
The `dns_config` block supports:

* `nameservers` - The list of DNS name servers.

* `options` - The list of DNS resolver options. The [options](#options) structure is documented below.

* `searches` - The list of DNS search domains.

<a name="options"></a>
The `options` block supports:

* `name` - The name of the DNS option.

* `value` - The value of the DNS option.

<a name="host_aliases"></a>
The `host_aliases` block supports:

* `ip` - The IP address of the host alias.

* `hostnames` - The list of hostnames for the IP.

<a name="image_pull_secrets"></a>
The `image_pull_secrets` block supports:

* `name` - The name of the referenced Secret.

<a name="readiness_gates"></a>
The `readiness_gates` block supports:

* `condition_type` - The type of the pod condition.

<a name="security_context"></a>
The pod-level `security_context` block supports:

* `fs_group` - The GID that applies to all containers in the Pod.

* `fs_group_change_policy` - The policy for changing ownership and permission of the volume.

* `run_as_group` - The GID to run the entrypoint of the container process.

* `run_as_non_root` - Whether containers must run as a non-root user.

* `run_as_user` - The UID to run the entrypoint of the container process.

* `supplemental_groups` - The list of groups applied to the first process in each container.

* `sysctls` - The list of sysctls to set in the Pod. The [sysctls](#sysctls) structure is documented below.

<a name="sysctls"></a>
The `sysctls` block supports:

* `name` - The name of the sysctl.

* `value` - The value of the sysctl.

<a name="volumes"></a>
The `volumes` block supports:

* `name` - The name of the volume.

* `config_map` - A ConfigMap mounted as a volume.
  The [volume_config_map](#volume_config_map) structure is documented below.

* `nfs` - An NFS volume. The [nfs](#nfs) structure is documented below.

* `persistent_volume_claim` - A reference to a PersistentVolumeClaim.
  The [persistent_volume_claim](#persistent_volume_claim) structure is documented below.

* `secret` - A Secret mounted as a volume. The [volume_secret](#volume_secret) structure is documented below.

* `projected` - A projected volume combining multiple sources.
  The [projected](#projected) structure is documented below.

<a name="volume_config_map"></a>
The volume `config_map` block supports:

* `name` - The name of the referenced ConfigMap.

* `default_mode` - The default permission mode of the projected files.

* `items` - The list of keys to project. The [items](#items) structure is documented below.

* `optional` - Whether the ConfigMap must be defined.

<a name="items"></a>
The `items` block supports:

* `key` - The key to project.

* `path` - The relative path to map the key to.

* `mode` - The permission mode of the file.

<a name="nfs"></a>
The `nfs` block supports:

* `path` - The path exported by the NFS server.

* `server` - The hostname or IP of the NFS server.

* `read_only` - Whether the NFS volume is mounted read-only.

<a name="persistent_volume_claim"></a>
The `persistent_volume_claim` block supports:

* `claim_name` - The name of the referenced PersistentVolumeClaim.

* `read_only` - Whether the volume is mounted read-only.

<a name="volume_secret"></a>
The volume `secret` block supports:

* `secret_name` - The name of the referenced Secret.

* `default_mode` - The default permission mode of the projected files.

* `items` - The list of keys to project. The [items](#items) structure is documented below.

* `optional` - Whether the Secret must be defined.

<a name="projected"></a>
The `projected` block supports:

* `default_mode` - The default permission mode of the projected files.

* `sources` - The list of volume projection sources.
  The [projected_sources](#projected_sources) structure is documented below.

<a name="projected_sources"></a>
The `sources` block supports:

* `config_map` - A ConfigMap projected into the volume.
  The [projected_config_map](#projected_config_map) structure is documented below.

* `downward_api` - Downward API info projected into the volume.
  The [downward_api](#downward_api) structure is documented below.

* `secret` - A Secret projected into the volume.
  The [projected_secret](#projected_secret) structure is documented below.

<a name="projected_config_map"></a>
The projected `config_map` block supports:

* `name` - The name of the referenced ConfigMap.

* `items` - The list of keys to project. The [items](#items) structure is documented below.

* `optional` - Whether the ConfigMap must be defined.

<a name="downward_api"></a>
The `downward_api` block supports:

* `items` - The list of downward API items.
  The [downward_api_items](#downward_api_items) structure is documented below.

<a name="downward_api_items"></a>
The downward API `items` block supports:

* `path` - The relative path of the file to create.

* `mode` - The permission mode of the file.

* `field_ref` - A selection of a metadata field.
  The [field_ref](#field_ref) structure is documented below.

* `resource_field_ref` - A selection of a container resource field.
  The [resource_field_ref](#resource_field_ref) structure is documented below.

<a name="field_ref"></a>
The `field_ref` block supports:

* `api_version` - The API version of the selected metadata field.

* `field_path` - The path of the field to select.

<a name="resource_field_ref"></a>
The `resource_field_ref` block supports:

* `container_name` - The name of the container.

* `resource` - The name of the resource to select.

<a name="projected_secret"></a>
The projected `secret` block supports:

* `name` - The name of the referenced Secret.

* `items` - The list of keys to project. The [items](#items) structure is documented below.

* `optional` - Whether the Secret must be defined.

<a name="status"></a>
The `status` block supports:

* `phase` - The current phase of the Pod.

* `conditions` - The list of Pod conditions.
  The [conditions](#conditions) structure is documented below.

<a name="conditions"></a>
The `conditions` block supports:

* `type` - The type of the condition.

* `status` - The status of the condition.

* `last_probe_time` - The last time the condition was probed.

* `last_transition_time` - The last time the condition transitioned.

* `reason` - The reason for the condition's last transition.

* `message` - A human-readable message describing the transition.
