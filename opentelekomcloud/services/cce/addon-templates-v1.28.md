# Addon Templates

Addon support configuration input depending on addon type and version. This page contains description of addon arguments
for the cluster with available k8s version `v1.28`.

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

#### For example
`swr_addr` in swiss region is `swr.eu-ch2.sc.otc.t-systems.com`

## Addon Inputs

### `autoscaler`

A component that automatically adjusts the size of a Kubernetes Cluster so that all pods have a place to run and there
are no unneeded nodes.
`template_version`: `1.28.17`

##### `basic`

```json
{
  "cceEndpoint": "https://cce.eu-de.otc.t-systems.com",
  "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
  "image_version": "1.28.17",
  "region": "eu-de",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "cce-addons"
}
```

##### `custom`

```json
{
  "cluster_id": "",
  "coresTotal": 32000,
  "expander": "priority,least-waste",
  "localVolumeNodeScalingEnabled": false,
  "logLevel": 4,
  "maxEmptyBulkDeleteFlag": 10,
  "maxNodeProvisionTime": 15,
  "maxNodesTotal": 1000,
  "memoryTotal": 128000,
  "multiAZBalance": false,
  "multiAZEnabled": false,
  "newEphemeralVolumesPodScaleUpDelay": 10,
  "node_match_expressions": [],
  "resetUnNeededWhenScaleUp": false,
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
  "skipNodesWithCustomControllerPods": true,
  "tenant_id": "",
  "tolerations": [
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "tolerationSeconds": 60
    },
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "tolerationSeconds": 60
    }
  ],
  "unremovableNodeRecheckTimeout": 5
}
```

### `coredns`

CoreDNS is a DNS server that chains plugins and provides Kubernetes DNS Services.
`template_version`: `1.29.4`

##### `basic`

```json
{
  "cluster_ip": "10.247.3.10",
  "image_version": "1.29.4",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "cce-addons"
}
```

##### `custom`

```json
{
  "annotations": {},
  "multiAZBalance": false,
  "multiAZEnabled": false,
  "node_match_expressions": [],
  "parameterSyncStrategy": "ensureConsistent",
  "servers": [
    {
      "plugins": [
        {
          "name": "bind",
          "parameters": "{$POD_IP}"
        },
        {
          "configBlock": "servfail 5s",
          "name": "cache",
          "parameters": 30
        },
        {
          "name": "errors"
        },
        {
          "name": "health",
          "parameters": "{$POD_IP}:8080"
        },
        {
          "name": "ready",
          "parameters": "{$POD_IP}:8081"
        },
        {
          "configBlock": "pods insecure\nfallthrough in-addr.arpa ip6.arpa",
          "name": "kubernetes",
          "parameters": "cluster.local in-addr.arpa ip6.arpa"
        },
        {
          "name": "loadbalance",
          "parameters": "round_robin"
        },
        {
          "name": "prometheus",
          "parameters": "{$POD_IP}:9153"
        },
        {
          "configBlock": "policy random",
          "name": "forward",
          "parameters": ". /etc/resolv.conf"
        },
        {
          "name": "reload"
        }
      ],
      "port": 5353,
      "zones": [
        {
          "zone": "."
        }
      ]
    }
  ],
  "stub_domains": {},
  "tolerations": [
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "tolerationSeconds": 60
    },
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "tolerationSeconds": 60
    }
  ],
  "upstream_nameservers": []
}
```

### `everest`

Everest is a cloud native container storage system based on CSI, used to support cloud storages services for Kubernetes.
`template_version`: `2.4.28`

##### `basic`

```json
{
  "bms_url": "bms.eu-de.otc.t-systems.com",
  "driver_init_image_version": "2.4.28",
  "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
  "everest_image_version": "2.4.28",
  "evs_url": "evs.eu-de.otc.t-systems.com",
  "iam_url": "iam.eu-de.otc.t-systems.com",
  "ims_url": "ims.eu-de.otc.t-systems.com",
  "obs_url": "obs.eu-de.otc.t-systems.com",
  "platform": "linux-amd64",
  "sfs30_url": "sfs3.eu-de.otc.t-systems.com",
  "sfs_turbo_url": "sfs-turbo.eu-de.otc.t-systems.com",
  "sfs_url": "sfs.eu-de.otc.t-systems.com",
  "supportHcs": false,
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "cce-addons"
}
```

##### `custom`

```json
{
  "annotations": {},
  "cluster_id": "",
  "cluster_name": "",
  "csi_attacher_detach_worker_threads": "60",
  "csi_attacher_worker_threads": "60",
  "default_vpc_id": "",
  "disable_auto_mount_secret": false,
  "enable_node_attacher": false,
  "flow_control": {},
  "multiAZBalance": false,
  "multiAZEnabled": false,
  "node_match_expressions": [],
  "number_of_reserved_disks": "6",
  "over_subscription": "80",
  "project_id": "",
  "tolerations": [
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "tolerationSeconds": 60
    },
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "tolerationSeconds": 60
    }
  ],
  "volume_attaching_flow_ctrl": "0"
}
```

### `metrics-server`

Metrics Server is a cluster-level resource usage data aggregator.
`template_version`: `1.3.60`

##### `basic`

```json
{
  "image_version": "v0.6.2",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "cce-addons"
}
```

##### `custom`

```json
{
  "annotations": {},
  "multiAZBalance": false,
  "multiAZEnabled": false,
  "node_match_expressions": [],
  "tolerations": [
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "tolerationSeconds": 60
    },
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "tolerationSeconds": 60
    }
  ]
}
```

### `gpu-beta`

A device plugin for nvidia.com/gpu resource on nvidia driver.
`template_version`: `2.6.4`

##### `basic`

```json
{
  "device_version": "2.6.4",
  "driver_version": "2.6.4",
  "obs_url": "obs.eu-de.otc.t-systems.com",
  "region": "eu-de",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "cce-addons"
}
```

##### `custom`

```json
{
  "annotations": {},
  "compatible_with_legacy_api": false,
  "component_schedulername": "default-scheduler",
  "disable_nvidia_gsp": true,
  "enable_fault_isolation": true,
  "enable_health_monitoring": true,
  "enable_metrics_monitoring": true,
  "enable_xgpu": false,
  "gpu_driver_config": {},
  "health_check_xids_v2": "74,79",
  "is_driver_from_nvidia": true,
  "metrics_delete_interval": 30000,
  "metrics_monitor_interval": 15000,
  "nvidia_driver_download_url": ""
}
```

### `npd`

Add-on for monitoring abnormal events of cluster nodes and connecting to a third-party monitoring platform.
`template_version`: `1.19.1`

##### `basic`

```json
{
  "image_version": "1.19.1",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "cce-addons"
}
```

##### `custom`

```json
{
  "annotations": {},
  "common": {},
  "feature_gates": "",
  "multiAZBalance": false,
  "multiAZEnabled": false,
  "node_match_expressions": [],
  "npc": {
    "maxTaintedNode": "10%"
  },
  "tolerations": [
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "tolerationSeconds": 60
    },
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "tolerationSeconds": 60
    }
  ]
}
```

### `volcano`

Batch processing platform based on Kubernetes.
`template_version`: `1.17.4`

##### `basic`

```json
{
  "ecsEndpoint": "https://ecs.eu-de.otc.t-systems.com",
  "platform": "linux-amd64",
  "swr_addr": "100.125.7.25:20202",
  "swr_user": "cce-addons"
}
```

##### `custom`

```json
{
  "annotations": {},
  "colocation_enable": "",
  "default_scheduler_conf": {
    "actions": "allocate, backfill",
    "metrics": {
      "interval": "30s",
      "type": ""
    },
    "tiers": [
      {
        "plugins": [
          {
            "name": "priority"
          },
          {
            "enableJobStarving": false,
            "enablePreemptable": false,
            "name": "gang"
          },
          {
            "name": "conformance"
          }
        ]
      },
      {
        "plugins": [
          {
            "enablePreemptable": false,
            "name": "drf"
          },
          {
            "name": "predicates"
          },
          {
            "name": "nodeorder"
          }
        ]
      },
      {
        "plugins": [
          {
            "name": "cce-gpu-topology-predicate"
          },
          {
            "name": "cce-gpu-topology-priority"
          },
          {
            "name": "xgpu"
          }
        ]
      },
      {
        "plugins": [
          {
            "name": "nodelocalvolume"
          },
          {
            "name": "nodeemptydirvolume"
          },
          {
            "name": "nodeCSIscheduling"
          },
          {
            "name": "networkresource"
          }
        ]
      }
    ]
  },
  "deschedulerPolicy": {
    "profiles": [
      {
        "name": "ProfileName",
        "pluginConfig": [
          {
            "args": {
              "nodeFit": true
            },
            "name": "DefaultEvictor"
          },
          {
            "args": {
              "evictableNamespaces": {
                "exclude": [
                  "kube-system"
                ]
              },
              "thresholds": {
                "cpu": 20,
                "memory": 20
              }
            },
            "name": "HighNodeUtilization"
          },
          {
            "args": {
              "evictableNamespaces": {
                "exclude": [
                  "kube-system"
                ]
              },
              "metrics": {
                "type": "prometheus_adaptor"
              },
              "targetThresholds": {
                "cpu": 80,
                "memory": 85
              },
              "thresholds": {
                "cpu": 30,
                "memory": 30
              }
            },
            "name": "LoadAware"
          }
        ],
        "plugins": {
          "balance": {
            "enabled": null
          }
        }
      }
    ]
  },
  "descheduler_enable": "",
  "deschedulingInterval": "10m",
  "enable_workload_balancer": false,
  "multiAZEnabled": false,
  "node_match_expressions": [],
  "oversubscription_ratio": 60,
  "recommendation_enable": "",
  "tolerations": [
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/not-ready",
      "operator": "Exists",
      "tolerationSeconds": 60
    },
    {
      "effect": "NoExecute",
      "key": "node.kubernetes.io/unreachable",
      "operator": "Exists",
      "tolerationSeconds": 60
    }
  ]
}
```
