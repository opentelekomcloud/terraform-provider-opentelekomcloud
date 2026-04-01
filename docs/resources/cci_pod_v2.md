---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_pod_v2"
sidebar_current: "docs-opentelekomcloud-resource-cci-pod-v2"
description: |-
  Manages a CCI v2 Pod resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI pod you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_pod_v2

Manages a CCI v2 Pod resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_cci_pod_v2" "test" {
  depends_on = [opentelekomcloud_cci_network_v2.example]

  namespace = opentelekomcloud_cci_namespace_v2.example.name
  name      = "my-pod"

  annotations = {
    "resource.cci.io/pod-size-specs" = "2.00_4.0"
    "resource.cci.io/instance-type"  = "general-computing"
  }

  containers {
    name  = "nginx"
    image = "swr.eu-de.otc.t-systems.com/cce-1.25/nginx:1.25.3-alpine"

    resources {
      limits = {
        cpu    = 2
        memory = "4G"
      }

      requests = {
        cpu    = 2
        memory = "4G"
      }
    }
  }

  image_pull_secrets {
    name = "imagepull-secret"
  }
}
```

## Argument Reference

The following arguments are supported:

~> **NOTE:** The Pod API only allows updating `containers[*].image`, `initContainers[*].image`,
`active_deadline_seconds`, `tolerations`, and `termination_grace_period_seconds`.
All other spec fields are marked as ForceNew and will trigger resource re-creation if changed.

* `namespace` - (Required, String, ForceNew) Specifies the namespace of the CCI pod.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the CCI pod.
  Changing this creates a new resource.

* `containers` - (Required, List) Specifies the containers of the CCI pod.
  Only `image` can be updated in-place; changing other container fields will trigger re-creation.
  The [containers](#containers) structure is documented below.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI pod.

* `labels` - (Optional, Map) Specifies the labels of the CCI pod.

* `active_deadline_seconds` - (Optional, Int) Specifies the active deadline seconds for the pod.

* `affinity` - (Optional, List, ForceNew) Specifies the affinity scheduling rules of the CCI pod.
  Changing this creates a new resource.
  The [affinity](#affinity) structure is documented below.

* `dns_config` - (Optional, List, ForceNew) Specifies the DNS configuration of the pod.
  Changing this creates a new resource.
  The [dns_config](#dns_config) structure is documented below.

* `dns_policy` - (Optional, String, ForceNew) Specifies the DNS policy of the pod.
  Valid values are `ClusterFirst`, `ClusterFirstWithHostNet`, `Default`, or `None`.
  Changing this creates a new resource.

* `ephemeral_containers` - (Optional, List, ForceNew) Specifies the ephemeral containers of the CCI pod.
  Changing this creates a new resource.
  The [ephemeral_containers](#ephemeral_containers) structure is documented below.

* `host_aliases` - (Optional, List, ForceNew) Specifies the host aliases of the CCI pod.
  Changing this creates a new resource.
  The [host_aliases](#host_aliases) structure is documented below.

* `hostname` - (Optional, String, ForceNew) Specifies the hostname of the pod.
  Changing this creates a new resource.

* `image_pull_secrets` - (Optional, List, ForceNew) Specifies the image pull secrets of the pod.
  Changing this creates a new resource.
  The [image_pull_secrets](#image_pull_secrets) structure is documented below.

* `init_containers` - (Optional, List) Specifies the init containers of the CCI pod.
  Only `image` can be updated in-place; changing other fields will trigger re-creation.
  The [init_containers](#containers) structure is documented below.

* `node_name` - (Optional, String, ForceNew) Specifies the node name of the CCI pod.
  Changing this creates a new resource.

* `overhead` - (Optional, Map, ForceNew) Specifies the overhead resources of the CCI pod.
  Changing this creates a new resource.

* `readiness_gates` - (Optional, List, ForceNew) Specifies the readiness gates of the CCI pod.
  Changing this creates a new resource.
  The [readiness_gates](#readiness_gates) structure is documented below.

* `restart_policy` - (Optional, String, ForceNew) The restart policy for all containers within the pod.
  Valid values are `Always`, `Never`, or `OnFailure`.
  Changing this creates a new resource.

* `scheduler_name` - (Optional, String, ForceNew) Specifies the scheduler name of the CCI pod.
  Changing this creates a new resource.

* `security_context` - (Optional, List, ForceNew) Specifies the pod-level security context.
  Changing this creates a new resource.
  The [security_context](#pod_security_context) structure is documented below.

* `set_hostname_as_fqdn` - (Optional, Bool, ForceNew) Specifies whether the pod hostname is configured
  as the pod FQDN. Changing this creates a new resource.

* `share_process_namespace` - (Optional, Bool, ForceNew) Specifies whether to share a single process namespace
  between all containers in a pod. Changing this creates a new resource.

* `termination_grace_period_seconds` - (Optional, Int) Specifies the grace period in seconds before the pod
  is forcefully terminated.

* `volumes` - (Optional, List, ForceNew) Specifies the volumes of the CCI pod.
  Changing this creates a new resource.
  The [volumes](#volumes) structure is documented below.

<a name="containers"></a>
The `containers` and `init_containers` block supports:

* `name` - (Required, String) Specifies the name of the container.

* `image` - (Optional, String) Specifies the image name of the container.

* `args` - (Optional, List) Specifies the arguments to the entrypoint of the container.

* `command` - (Optional, List) Specifies the entrypoint command of the container.

* `env` - (Optional, List) Specifies the environment variables of the container.
  The [env](#containers_env) structure is documented below.

* `env_from` - (Optional, List) Specifies the sources to populate environment variables.
  The [env_from](#containers_env_from) structure is documented below.

* `lifecycle` - (Optional, List) Specifies the lifecycle hooks of the container.
  The [lifecycle](#containers_lifecycle) structure is documented below.

* `liveness_probe` - (Optional, List) Specifies the liveness probe of the container.
  The [liveness_probe](#containers_probe) structure is documented below.

* `readiness_probe` - (Optional, List) Specifies the readiness probe of the container.
  The [readiness_probe](#containers_probe) structure is documented below.

* `startup_probe` - (Optional, List) Specifies the startup probe of the container.
  The [startup_probe](#containers_probe) structure is documented below.

* `ports` - (Optional, List) Specifies the ports exposed by the container.
  The [ports](#containers_ports) structure is documented below.

* `resources` - (Optional, List) Specifies the compute resources of the container.
  The [resources](#containers_resources) structure is documented below.

* `security_context` - (Optional, List) Specifies the security context of the container.
  The [security_context](#containers_security_context) structure is documented below.

* `stdin` - (Optional, Bool) Specifies whether this container should allocate a buffer for stdin.

* `stdin_once` - (Optional, Bool) Specifies whether the container runtime should close the stdin channel
  after it has been opened by a single attach.

* `termination_message_path` - (Optional, String) Specifies the path at which the file to which the
  container's termination message will be written.

* `termination_message_policy` - (Optional, String) Specifies how the termination message should be populated.

* `tty` - (Optional, Bool) Specifies whether this container should allocate a TTY for itself.

* `working_dir` - (Optional, String) Specifies the working directory of the container.

* `volume_mounts` - (Optional, List) Specifies the volume mounts of the container.
  The [volume_mounts](#containers_volume_mounts) structure is documented below.

<a name="ephemeral_containers"></a>
The `ephemeral_containers` block supports:

* `name` - (Required, String) Specifies the name of the ephemeral container.

* `image` - (Optional, String) Specifies the image name of the ephemeral container.

* `args` - (Optional, List) Specifies the arguments to the entrypoint.

* `command` - (Optional, List) Specifies the entrypoint command.

* `env` - (Optional, List) Specifies the environment variables.
  The [env](#containers_env) structure is documented below.

* `env_from` - (Optional, List) Specifies the sources to populate environment variables.
  The [env_from](#containers_env_from) structure is documented below.

* `security_context` - (Optional, List) Specifies the security context.
  The [security_context](#containers_security_context) structure is documented below.

* `stdin` - (Optional, Bool) Specifies whether to allocate a buffer for stdin.

* `stdin_once` - (Optional, Bool) Specifies whether the runtime should close the stdin channel.

* `target_container_name` - (Optional, String) Specifies the name of the target container
  to which the ephemeral container will be attached.

* `termination_message_path` - (Optional, String) Specifies the termination message path.

* `termination_message_policy` - (Optional, String) Specifies the termination message policy.

* `tty` - (Optional, Bool) Specifies whether to allocate a TTY.

* `volume_mounts` - (Optional, List) Specifies the volume mounts.
  The [volume_mounts](#containers_volume_mounts) structure is documented below.

* `working_dir` - (Optional, String) Specifies the working directory.

<a name="containers_volume_mounts"></a>
The `volume_mounts` block supports:

* `mount_path` - (Required, String) Specifies the path within the container at which the volume should be mounted.

* `name` - (Required, String) Specifies the name of the volume to mount. Must match the name of a volume.

* `read_only` - (Optional, Bool) Specifies whether the volume is mounted as read-only. Defaults to `false`.

* `sub_path` - (Optional, String) Specifies the sub-path inside the volume to mount.

* `sub_path_expr` - (Optional, String) Specifies the expanded sub-path using environment variables.

* `extend_path_mode` - (Optional, String) Specifies the extend path mode of the volume mount.

<a name="containers_env"></a>
The `env` block supports:

* `name` - (Required, String) Specifies the name of the environment variable.

* `value` - (Optional, String) Specifies the value of the environment variable.

<a name="containers_env_from"></a>
The `env_from` block supports:

* `config_map_ref` - (Optional, List) Specifies the ConfigMap to select from.
  The [config_map_ref](#containers_env_source) structure is documented below.

* `prefix` - (Optional, String) Specifies an optional identifier to prepend to each key.

* `secret_ref` - (Optional, List) Specifies the Secret to select from.
  The [secret_ref](#containers_env_source) structure is documented below.

<a name="containers_env_source"></a>
The `config_map_ref` and `secret_ref` block supports:

* `name` - (Optional, String) Specifies the name of the referent.

* `optional` - (Optional, Bool) Specifies whether the referent must be defined.

<a name="containers_lifecycle"></a>
The `lifecycle` block supports:

* `post_start` - (Optional, List) Specifies the handler called after a container is created.
  The [post_start](#containers_lifecycle_handler) structure is documented below.

* `pre_stop` - (Optional, List) Specifies the handler called before a container is terminated.
  The [pre_stop](#containers_lifecycle_handler) structure is documented below.

<a name="containers_lifecycle_handler"></a>
The `post_start` and `pre_stop` block supports:

* `exec` - (Optional, List) Specifies the exec-based handler.
  The [exec](#exec_action) structure is documented below.

* `http_get` - (Optional, List) Specifies the HTTP GET-based handler.
  The [http_get](#http_get_action) structure is documented below.

<a name="containers_probe"></a>
The `liveness_probe`, `readiness_probe`, and `startup_probe` block supports:

* `exec` - (Optional, List) Specifies the exec-based probe action.
  The [exec](#exec_action) structure is documented below.

* `http_get` - (Optional, List) Specifies the HTTP GET-based probe action.
  The [http_get](#http_get_action) structure is documented below.

* `failure_threshold` - (Optional, Int) Specifies the minimum consecutive failures for the probe to be
  considered failed after having succeeded.

* `initial_delay_seconds` - (Optional, Int) Specifies the number of seconds after the container has started
  before probes are initiated.

* `period_seconds` - (Optional, Int) Specifies how often (in seconds) to perform the probe.

* `success_threshold` - (Computed, Int) The minimum consecutive successes for the probe to be considered
  successful after having failed.

* `termination_grace_period_seconds` - (Optional, Int) Specifies the grace period in seconds before
  the pod is forcefully terminated when the probe fails.

<a name="exec_action"></a>
The `exec` block supports:

* `command` - (Optional, List) Specifies the command line to execute inside the container.

<a name="http_get_action"></a>
The `http_get` block supports:

* `host` - (Optional, String) Specifies the hostname to connect to.

* `http_headers` - (Optional, List) Specifies the custom headers to set in the request.
  The [http_headers](#http_headers) structure is documented below.

* `path` - (Optional, String) Specifies the path to access on the HTTP server.

* `port` - (Required, String) Specifies the port to access on the container.

* `scheme` - (Optional, String) Specifies the scheme to use for connecting to the host.

<a name="http_headers"></a>
The `http_headers` block supports:

* `name` - (Required, String) Specifies the name of the custom HTTP header.

* `value` - (Required, String) Specifies the value of the custom HTTP header.

<a name="containers_ports"></a>
The `ports` block supports:

* `container_port` - (Required, Int) Specifies the port number to expose on the pod's IP address.

* `name` - (Optional, String) Specifies the name of the port.

* `protocol` - (Optional, String) Specifies the protocol for the port. Valid values are `TCP`, `UDP`.

<a name="containers_resources"></a>
The `resources` block supports:

* `limits` - (Optional, Map) Specifies the maximum amount of compute resources allowed.

* `requests` - (Optional, Map) Specifies the minimum amount of compute resources required.

<a name="containers_security_context"></a>
The `security_context` block supports:

* `capabilities` - (Optional, List) Specifies the capabilities to add/drop.
  The [capabilities](#capabilities) structure is documented below.

* `proc_mount` - (Optional, String) Specifies the type of proc mount to use for the container.

* `read_only_root_file_system` - (Optional, Bool) Specifies whether this container has a read-only root filesystem.

* `run_as_group` - (Optional, Int) Specifies the GID to run the entrypoint of the container process.

* `run_as_non_root` - (Optional, Bool) Specifies that the container must run as a non-root user.

* `run_as_user` - (Optional, Int) Specifies the UID to run the entrypoint of the container process.

<a name="capabilities"></a>
The `capabilities` block supports:

* `add` - (Optional, List) Specifies the list of capabilities to add.

* `drop` - (Optional, List) Specifies the list of capabilities to drop.

<a name="affinity"></a>
The `affinity` block supports:

* `node_affinity` - (Optional, List) Specifies the node affinity scheduling rules.
  The [node_affinity](#affinity_node_affinity) structure is documented below.

* `pod_anti_affinity` - (Optional, List) Specifies the pod anti-affinity scheduling rules.
  The [pod_anti_affinity](#affinity_pod_anti_affinity) structure is documented below.

<a name="affinity_node_affinity"></a>
The `node_affinity` block supports:

* `required_during_scheduling_ignored_during_execution` - (Optional, List) Specifies the required node selector.
  The [required_during_scheduling_ignored_during_execution](#node_affinity_required) structure is documented below.

<a name="node_affinity_required"></a>
The `required_during_scheduling_ignored_during_execution` block supports:

* `node_selector_terms` - (Required, List) Specifies the list of node selector terms.
  The [node_selector_terms](#node_selector_terms) structure is documented below.

<a name="node_selector_terms"></a>
The `node_selector_terms` block supports:

* `match_expressions` - (Optional, List) Specifies the list of node selector requirements by node's labels.
  The [match_expressions](#match_expressions) structure is documented below.

<a name="affinity_pod_anti_affinity"></a>
The `pod_anti_affinity` block supports:

* `preferred_during_scheduling_ignored_during_execution` - (Optional, List) Specifies the preferred scheduling terms.
  The [preferred_during_scheduling_ignored_during_execution](#pod_anti_affinity_preferred) structure is documented below.

* `required_during_scheduling_ignored_during_execution` - (Optional, List) Specifies the required scheduling terms.
  The [required_during_scheduling_ignored_during_execution](#pod_affinity_term) structure is documented below.

<a name="pod_anti_affinity_preferred"></a>
The `preferred_during_scheduling_ignored_during_execution` block supports:

* `weight` - (Required, Int) Specifies the weight associated with matching the corresponding pod affinity term.

* `pod_affinity_term` - (Required, List) Specifies the pod affinity term.
  The [pod_affinity_term](#pod_affinity_term) structure is documented below.

<a name="pod_affinity_term"></a>
The `pod_affinity_term` and `required_during_scheduling_ignored_during_execution` block supports:

* `topology_key` - (Required, String) Specifies the topology key.

* `label_selector` - (Optional, List) Specifies the label selector.
  The [label_selector](#label_selector) structure is documented below.

* `namespaces` - (Optional, List) Specifies the namespaces to match against.

<a name="label_selector"></a>
The `label_selector` block supports:

* `match_labels` - (Optional, Map) Specifies the map of key-value pairs for matching.

* `match_expressions` - (Optional, List) Specifies the list of label selector requirements.
  The [match_expressions](#match_expressions) structure is documented below.

<a name="match_expressions"></a>
The `match_expressions` block supports:

* `key` - (Required, String) Specifies the label key that the selector applies to.

* `operator` - (Required, String) Specifies the operator. Valid values are `In`, `NotIn`, `Exists`, `DoesNotExist`.

* `values` - (Optional, List) Specifies the list of string values.

<a name="dns_config"></a>
The `dns_config` block supports:

* `nameservers` - (Optional, List) Specifies the list of DNS name servers.

* `options` - (Optional, List) Specifies the list of DNS resolver options.
  The [options](#dns_config_options) structure is documented below.

* `searches` - (Optional, List) Specifies the list of DNS search domains.

<a name="dns_config_options"></a>
The `options` block supports:

* `name` - (Optional, String) Specifies the name of the option.

* `value` - (Optional, String) Specifies the value of the option.

<a name="host_aliases"></a>
The `host_aliases` block supports:

* `hostnames` - (Optional, List) Specifies the list of hostnames for the IP address.

* `ip` - (Optional, String) Specifies the IP address of the host.

<a name="image_pull_secrets"></a>
The `image_pull_secrets` block supports:

* `name` - (Optional, String) Specifies the name of the image pull secret.

<a name="readiness_gates"></a>
The `readiness_gates` block supports:

* `condition_type` - (Required, String) Specifies the condition type of the readiness gate.

<a name="pod_security_context"></a>
The `security_context` block supports:

* `fs_group` - (Optional, Int) Specifies the GID applied to all containers in the pod.

* `fs_group_change_policy` - (Optional, String) Specifies the behavior of changing ownership
  and permission of the volume. Valid values are `Always`, `OnRootMismatch`.

* `run_as_group` - (Optional, Int) Specifies the GID to run the entrypoint of the container process.

* `run_as_non_root` - (Optional, Bool) Specifies that the container must run as a non-root user.

* `run_as_user` - (Optional, Int) Specifies the UID to run the entrypoint of the container process.

* `supplemental_groups` - (Optional, List) Specifies the list of supplemental groups applied to the pod.

* `sysctls` - (Optional, List) Specifies the list of namespaced sysctls.
  The [sysctls](#security_context_sysctls) structure is documented below.

<a name="security_context_sysctls"></a>
The `sysctls` block supports:

* `name` - (Required, String) Specifies the name of the sysctl.

* `value` - (Required, String) Specifies the value of the sysctl.

<a name="volumes"></a>
The `volumes` block supports:

* `name` - (Required, String) Specifies the name of the volume.

* `config_map` - (Optional, List) Specifies the ConfigMap volume source.
  The [config_map](#volumes_config_map) structure is documented below.

* `nfs` - (Optional, List) Specifies the NFS volume source.
  The [nfs](#volumes_nfs) structure is documented below.

* `persistent_volume_claim` - (Optional, List) Specifies the PersistentVolumeClaim volume source.
  The [persistent_volume_claim](#volumes_persistent_volume_claim) structure is documented below.

* `projected` - (Optional, List) Specifies the projected volume source.
  The [projected](#volumes_projected) structure is documented below.

* `secret` - (Optional, List) Specifies the Secret volume source.
  The [secret](#volumes_secret) structure is documented below.

<a name="volumes_config_map"></a>
The `config_map` block supports:

* `name` - (Optional, String) Specifies the name of the ConfigMap.

* `default_mode` - (Optional, Int) Specifies the default file mode bits for the volume.

* `items` - (Optional, List) Specifies the list of key-to-path mappings.
  The [items](#key_to_path) structure is documented below.

* `optional` - (Optional, Bool) Specifies whether the ConfigMap must be defined.

<a name="volumes_nfs"></a>
The `nfs` block supports:

* `path` - (Required, String) Specifies the NFS server path.

* `server` - (Required, String) Specifies the NFS server hostname or IP address.

* `read_only` - (Optional, Bool) Specifies whether the NFS volume is read-only.

<a name="volumes_persistent_volume_claim"></a>
The `persistent_volume_claim` block supports:

* `claim_name` - (Required, String) Specifies the name of the PersistentVolumeClaim.

* `read_only` - (Optional, Bool) Specifies whether the volume is read-only.

<a name="volumes_projected"></a>
The `projected` block supports:

* `default_mode` - (Optional, Int) Specifies the default file mode bits for the volume.

* `sources` - (Optional, List) Specifies the list of volume projection sources.
  The [sources](#volumes_projected_sources) structure is documented below.

<a name="volumes_projected_sources"></a>
The `sources` block supports:

* `config_map` - (Optional, List) Specifies the ConfigMap projection.
  The [config_map](#volumes_projected_sources_config_map) structure is documented below.

* `downward_api` - (Optional, List) Specifies the Downward API projection.
  The [downward_api](#volumes_projected_sources_downward_api) structure is documented below.

* `secret` - (Optional, List) Specifies the Secret projection.
  The [secret](#volumes_projected_sources_secret) structure is documented below.

<a name="volumes_projected_sources_config_map"></a>
The `config_map` block supports:

* `name` - (Optional, String) Specifies the name of the ConfigMap.

* `items` - (Optional, List) Specifies the key-to-path mappings.
  The [items](#key_to_path) structure is documented below.

* `optional` - (Optional, Bool) Specifies whether the ConfigMap must be defined.

<a name="volumes_projected_sources_downward_api"></a>
The `downward_api` block supports:

* `items` - (Optional, List) Specifies the list of downward API volume files.
  The [items](#downward_api_items) structure is documented below.

<a name="downward_api_items"></a>
The `items` block supports:

* `path` - (Required, String) Specifies the relative path of the file to create.

* `mode` - (Optional, Int) Specifies the file mode bits for the file.

* `field_ref` - (Optional, List) Specifies the object field selector.
  The [field_ref](#downward_api_field_ref) structure is documented below.

* `resource_field_ref` - (Optional, List) Specifies the container resource field selector.
  The [resource_field_ref](#downward_api_resource_field_ref) structure is documented below.

<a name="downward_api_field_ref"></a>
The `field_ref` block supports:

* `field_path` - (Required, String) Specifies the field path to select.

* `api_version` - (Optional, String) Specifies the API version of the field path.

<a name="downward_api_resource_field_ref"></a>
The `resource_field_ref` block supports:

* `resource` - (Required, String) Specifies the resource to select.

* `container_name` - (Optional, String) Specifies the name of the container.

<a name="volumes_projected_sources_secret"></a>
The `secret` block supports:

* `name` - (Optional, String) Specifies the name of the Secret.

* `items` - (Optional, List) Specifies the key-to-path mappings.
  The [items](#key_to_path) structure is documented below.

* `optional` - (Optional, Bool) Specifies whether the Secret must be defined.

<a name="key_to_path"></a>
The `items` block supports:

* `key` - (Required, String) Specifies the key to project.

* `path` - (Required, String) Specifies the relative path of the file to map the key to.

* `mode` - (Optional, Int) Specifies the file mode bits for the file.

<a name="volumes_secret"></a>
The `secret` block supports:

* `secret_name` - (Optional, String) Specifies the name of the Secret.

* `default_mode` - (Optional, Int) Specifies the default file mode bits for the volume.

* `items` - (Optional, List) Specifies the list of key-to-path mappings.
  The [items](#key_to_path) structure is documented below.

* `optional` - (Optional, Bool) Specifies whether the Secret must be defined.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<namespace>/<name>`.

* `region` - The region of the CCI pod.

* `api_version` - The API version of the CCI pod.

* `kind` - The kind of the CCI pod.

* `creation_timestamp` - The creation timestamp of the CCI pod.

* `resource_version` - The resource version of the CCI pod.

* `finalizers` - The finalizers of the CCI pod.

* `uid` - The UID of the CCI pod.

* `status` - The status of the CCI pod.
  The [status](#attr_status) structure is documented below.

<a name="attr_status"></a>
The `status` block contains:

* `phase` - The phase of the CCI pod. Valid values include `Pending`, `Running`, `Succeeded`, `Failed`, `Unknown`.

* `conditions` - The conditions of the CCI pod.
  The [conditions](#attr_status_conditions) structure is documented below.

<a name="attr_status_conditions"></a>
The `conditions` block contains:

* `type` - The type of the pod condition.

* `status` - The status of the condition.

* `last_probe_time` - The last time the condition was probed.

* `last_transition_time` - The last time the condition transitioned from one status to another.

* `reason` - The reason for the condition's last transition.

* `message` - A human-readable message indicating details about the transition.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

