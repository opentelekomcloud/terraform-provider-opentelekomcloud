# Addon Templates

Addon support configuration input depending on addon type and version. This page contains description of addon arguments
for the cluster with available k8s version `v1.17.9`.

Up to date reference of addon arguments for your cluster you can get using API for listing CCE addon templates
at `https://<cluster_id>.cce.eu-de.otc.t-systems.com/api/v3/addontemplates`, where `<cluster_id>` is ID of the created
cluster.

Following addon templates exist in the addon template list:

- [`autoscaler`](#autoscaler)
- [`coredns`](#coredns)
- [`everest`](#everest)
- [`metrics-server`](#metrics-server)
- [`storage-driver`](#storage-driver)
- [`gpu-beta`](#gpu-beta)

All addons accept `basic` and some can accept `custom` input values.

## Addon Inputs

### `autoscaler`

A component that automatically adjusts the size of a Kubernetes Cluster so that all pods have a place to run and there
are no unneeded nodes.

##### `basic`

```json
{
  "cceEndpoint": "https://cce.eu-de.otc.t-systems.com",
  "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
  "euleros_version": "2.2.5",
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
  "maxEmptyBulkDeleteFlag": 10,
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
  "scaleUpUtilizationEnabled": false,
  "tenant_id": ""
}
```

### `coredns`

CoreDNS is a DNS server that chains plugins and provides Kubernetes DNS Services.

##### `basic`

```json
{
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

##### `basic`

```json
{
  "bms_url": "bms.eu-de.otc.t-systems.com",
  "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
  "euleros_version": "2.2.5",
  "evs_url": "evs.eu-de.otc.t-systems.com",
  "iam_url": "iam.eu-de.otc.t-systems.com",
  "ims_url": "ims.eu-de.otc.t-systems.com",
  "obs_url": "obs.eu-de.otc.t-systems.com",
  "platform": "linux-amd64",
  "sfs_turbo_url": "sfs_turbo.eu-de.otc.t-systems.com",
  "sfs_url": "sfs.eu-de.otc.t-systems.com",
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

##### `basic`

```json
{
  "euleros_version": "2.2.5",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "hwofficial"
}
```

##### `custom`

_Not supported_

### `storage-driver`

A Kubernetes FlexVolume Driver used to support storages services.

##### `basic`

```json
{
  "euleros_version": "2.2.5",
  "obs_url": "obs.eu-de.otc.t-systems.com",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "hwofficial"
}
```

##### `custom`

_Not supported_

### `gpu-beta`

A device plugin for nvidia.com/gpu resource on nvidia driver.

##### `basic`

```json
{
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
