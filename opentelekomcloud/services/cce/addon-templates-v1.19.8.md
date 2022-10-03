# Addon Templates

Addon support configuration input depending on addon type and version. This page contains description of addon arguments
for the cluster with available k8s version `v1.19.8`.

Up to date reference of addon arguments for your cluster you can get using API for listing CCE addon templates
at `https://<cluster_id>.cce.eu-de.otc.t-systems.com/api/v3/addontemplates`, where `<cluster_id>` is ID of the created
cluster.

Following addon templates exist in the addon template list:

- [`autoscaler`](#autoscaler)
- [`coredns`](#coredns)
- [`everest`](#everest)
- [`metrics-server`](#metrics-server)
- [`gpu-beta`](#gpu-beta)

All addons accept `basic` and some can accept `custom` input values.

In some regions can be impossible to create addon without `tenant_id` which equal to `project_id` of the region you're deploying into,
also be aware of that `swr_addr` can be also different in regions.

## Addon Inputs

### `autoscaler`

A component that automatically adjusts the size of a Kubernetes Cluster so that all pods have a place to run and there
are no unneeded nodes.
`template_version`: `1.19.1`

See the [OTC Cloud Container Engine Addon Documentation](https://docs.otc.t-systems.com/en-us/usermanual2/cce/cce_01_0154.html) for more information.

##### `basic`

```json
{
  "cceEndpoint": "https://cce.eu-de.otc.t-systems.com",
  "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
  "image_version": "1.19.1",
  "platform": "linux-amd64",
  "region": "eu-de",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "hwofficial"
}
```

##### `custom`

```json
{
  "cluster_id": "",
  "coresTotal": 32000,
  "expander": "priority",
  "logLevel": 4,
  "maxEmptyBulkDeleteFlag": 10,
  "maxNodeProvisionTime": 15,
  "maxNodesTotal": 1000,
  "memoryTotal": 128000,
  "scaleDownDelayAfterAdd": 10,
  "scaleDownDelayAfterDelete": 10,
  "scaleDownDelayAfterFailure": 3,
  "scaleDownEnabled": false,
  "scaleDownUnneededTime": 10,
  "scaleDownUtilizationThreshold": 0.5,
  "scaleUpCpuUtilizationThreshold": 1,
  "scaleUpMemUtilizationThreshold": 1,
  "scaleUpUnscheduledPodEnabled": true,
  "scaleUpUtilizationEnabled": true,
  "tenant_id": "",
  "unremovableNodeRecheckTimeout": 5
}
```

### `coredns`

CoreDNS is a DNS server that chains plugins and provides Kubernetes DNS Services.
`template_version`: `1.17.4`

See the [OTC Cloud Container Engine Addon Documentation](https://docs.otc.t-systems.com/en-us/usermanual2/cce/cce_01_0129.html) for more information.

##### `basic`

```json
{
  "cluster_ip": "10.247.3.10",
  "image_version": "1.17.4",
  "platform": "linux-amd64",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "hwofficial"
}
```

##### `custom`

```json
{
  "stub_domains": "",
  "upstream_nameservers": ""
}
```

### `everest`

Everest is a cloud native container storage system based on CSI, used to support cloud storages services for Kubernetes.
`template_version`: `1.2.2`

See the [OTC Cloud Container Engine Addon Documentation](https://docs.otc.t-systems.com/en-us/usermanual2/cce/cce_01_0066.html) for more information.

##### `basic`

```json
{
  "bms_url": "bms.eu-de.otc.t-systems.com",
  "controller_image_version": "1.2.2",
  "driver_image_version": "1.2.2",
  "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
  "evs_url": "evs.eu-de.otc.t-systems.com",
  "iam_url": "iam.eu-de.otc.t-systems.com",
  "ims_url": "ims.eu-de.otc.t-systems.com",
  "obs_url": "obs.eu-de.otc.t-systems.com",
  "platform": "linux-amd64",
  "sfs_turbo_url": "sfs_turbo.eu-de.otc.t-systems.com",
  "sfs_url": "sfs.eu-de.otc.t-systems.com",
  "supportHcs": false,
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "hwofficial"
}
```

##### `custom`

```json
{
  "cluster_id": "",
  "default_vpc_id": "",
  "project_id": ""
}
```

### `metrics-server`

Metrics Server is a cluster-level resource usage data aggregator.
`template_version`: `1.0.6`

See the [OTC Cloud Container Engine Addon Documentation](https://docs.otc.t-systems.com/en-us/usermanual2/cce/cce_01_0205.html) for more information.

##### `basic`

```json
{
  "image_version": "v0.3.7",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "hwofficial"
}
```

##### `custom`

_Not supported_

### `gpu-beta`

A device plugin for nvidia.com/gpu resource on nvidia driver.
`template_version`: `1.1.19`

See the [OTC Cloud Container Engine Addon Documentation](https://docs.otc.t-systems.com/en-us/usermanual2/cce/cce_01_0141.html) for more information.

##### `basic`

```json
{
  "device_version": "1.0.10",
  "driver_version": "1.1.15",
  "obs_url": "obs.eu-de.otc.t-systems.com",
  "region": "eu-de",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "hwofficial"
}
```

##### `custom`

```json
{
  "is_driver_from_nvidia": true,
  "nvidia_driver_download_url": ""
}
```
