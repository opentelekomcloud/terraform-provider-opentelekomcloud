package cci

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/pod"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func DataSourceCCIPodsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCIPodsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pods": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"finalizers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"active_deadline_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"affinity": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_affinity": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     podsNodeAffinitySchema(),
									},
									"pod_anti_affinity": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     podsAntiAffinitySchema(),
									},
								},
							},
						},
						"containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersSchema(),
						},
						"dns_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsDNSConfigSchema(),
						},
						"dns_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ephemeral_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersSchema(),
						},
						"host_aliases": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hostnames": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hostname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_pull_secrets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"init_containers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersSchema(),
						},
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"overhead": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"readiness_gates": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"restart_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_context": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsSecurityContextSchema(),
						},
						"set_hostname_as_fqdn": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"share_process_namespace": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"termination_grace_period_seconds": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volumes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsVolumesSchema(),
						},
						"status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsStatusSchema(),
						},
					},
				},
			},
		},
	}
}

func podsStatusSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"phase": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_probe_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_transition_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func podsVolumesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"projected": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSchema(),
			},
			"config_map": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesConfigMapSchema(),
			},
			"nfs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesNfsSchema(),
			},
			"persistent_volume_claim": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesPersistentVolumeClaimSchema(),
			},
			"secret": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesSecretSchema(),
			},
		},
	}
}

func podsVolumesSecretSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"secret_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func podsVolumesKeyToPathSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func podsVolumesPersistentVolumeClaimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"claim_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func podsVolumesNfsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func podsVolumesConfigMapSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func podsVolumesProjectedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSourcesSchema(),
			},
			"default_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func podsVolumesProjectedSourcesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_map": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSourcesConfigMapSchema(),
			},
			"downward_api": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSourcesDownwardAPISchema(),
			},
			"secret": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesProjectedSourcesSecretSchema(),
			},
		},
	}
}

func podsVolumesProjectedSourcesSecretSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func podsVolumesProjectedSourcesDownwardAPISchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsDownwardAPIFileSchema(),
			},
		},
	}
}

func podsDownwardAPIFileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_field_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func podsVolumesProjectedSourcesConfigMapSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsVolumesKeyToPathSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func podsSecurityContextSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fs_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fs_group_change_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"run_as_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"run_as_non_root": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"run_as_user": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"supplemental_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sysctls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func podsDNSConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nameservers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"searches": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func podsNodeAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsNodeSelectorSchema(),
			},
		},
	}
}

func podsNodeSelectorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_selector_terms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsNodeSelectorTermSchema(),
			},
		},
	}
}

func podsNodeSelectorTermSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsNodeSelectorRequirementSchema(),
			},
		},
	}
}

func podsNodeSelectorRequirementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func podsAntiAffinitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"preferred_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsWeightedPodAffinityTermSchema(),
			},
			"required_during_scheduling_ignored_during_execution": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsAffinityTermSchema(),
			},
		},
	}
}

func podsAffinityTermSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"label_selector": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsLabelSelectorSchema(),
			},
			"namespaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"topology_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func podsLabelSelectorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"match_labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"match_expressions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsMatchExpressionsSchema(),
			},
		},
	}
}

func podsMatchExpressionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func podsWeightedPodAffinityTermSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"pod_affinity_term": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsAffinityTermSchema(),
			},
			"weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func podsContainersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"args": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"command": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"env": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"env_from": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersEnvFromSchema(),
			},
			"image": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lifecycle": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"post_start": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersLifecycleHandlerSchema(),
						},
						"pre_stop": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     podsContainersLifecycleHandlerSchema(),
						},
					},
				},
			},
			"liveness_probe": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersProbeSchema(),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"readiness_probe": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersProbeSchema(),
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limits": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"requests": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"security_context": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersSecurityContextSchema(),
			},
			"startup_probe": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersProbeSchema(),
			},
			"stdin": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"stdin_once": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"target_container_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"termination_message_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"termination_message_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tty": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"working_dir": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"volume_mounts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"extend_path_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mount_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sub_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_path_expr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func podsContainersSecurityContextSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"capabilities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"add": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"drop": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"proc_mount": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_only_root_file_system": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"run_as_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"run_as_non_root": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"run_as_user": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func podsContainersProbeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerExecSchema(),
			},
			"failure_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"http_get": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerHttpGetActionSchema(),
			},
			"initial_delay_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"period_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"success_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func podsContainersLifecycleHandlerSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerExecSchema(),
			},
			"http_get": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersLifecycleHandlerHttpGetActionSchema(),
			},
		},
	}
}

func podsContainersLifecycleHandlerHttpGetActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_headers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheme": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func podsContainersLifecycleHandlerExecSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"command": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func podsContainersEnvFromSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_map_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersEnvSourceSchema(),
			},
			"prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_ref": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podsContainersEnvSourceSchema(),
			},
		},
	}
}

func podsContainersEnvSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceCCIPodsV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	region := config.GetRegion(d)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(region)
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)

	podList := make([]pod.Pod, 0)
	if name, ok := d.GetOk("name"); ok {
		resp, err := pod.Get(client, ns, name.(string))
		if err != nil {
			return diag.Errorf("error getting the CCI v2 Pod (%s/%s) from the server: %s", ns, name.(string), err)
		}
		podList = append(podList, *resp)
	} else {
		resp, err := pod.List(client, ns, pod.ListOpts{})
		if err != nil {
			return diag.Errorf("error querying CCI v2 Pods under namespace %s: %s", ns, err)
		}
		podList = resp
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("pods", flattenPods(podList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPods(podList []pod.Pod) []map[string]interface{} {
	if len(podList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(podList))
	for _, p := range podList {
		result = append(result, map[string]interface{}{
			"name":                             p.Metadata.Name,
			"namespace":                        p.Metadata.Namespace,
			"api_version":                      p.APIVersion,
			"kind":                             p.Kind,
			"annotations":                      p.Metadata.Annotations,
			"labels":                           p.Metadata.Labels,
			"creation_timestamp":               p.Metadata.CreationTimestamp,
			"resource_version":                 p.Metadata.ResourceVersion,
			"uid":                              p.Metadata.UID,
			"finalizers":                       p.Metadata.Finalizers,
			"active_deadline_seconds":          derefInt64(p.Spec.ActiveDeadlineSeconds),
			"affinity":                         flattenAffinity(p.Spec.Affinity),
			"containers":                       flattenContainers(p.Spec.Containers),
			"dns_config":                       flattenDNSConfig(p.Spec.DNSConfig),
			"dns_policy":                       p.Spec.DNSPolicy,
			"ephemeral_containers":             flattenEphemeralContainers(p.Spec.EphemeralContainers),
			"host_aliases":                     flattenHostAliases(p.Spec.HostAliases),
			"hostname":                         p.Spec.Hostname,
			"image_pull_secrets":               flattenLocalObjectReferences(p.Spec.ImagePullSecrets),
			"init_containers":                  flattenContainers(p.Spec.InitContainers),
			"node_name":                        p.Spec.NodeName,
			"overhead":                         p.Spec.Overhead,
			"readiness_gates":                  flattenReadinessGates(p.Spec.ReadinessGates),
			"restart_policy":                   p.Spec.RestartPolicy,
			"scheduler_name":                   p.Spec.SchedulerName,
			"security_context":                 flattenPodSecurityContext(p.Spec.SecurityContext),
			"set_hostname_as_fqdn":             derefBool(p.Spec.SetHostnameAsFQDN),
			"share_process_namespace":          derefBool(p.Spec.ShareProcessNamespace),
			"termination_grace_period_seconds": derefInt64(p.Spec.TerminationGracePeriodSeconds),
			"volumes":                          flattenVolumes(p.Spec.Volumes),
			"status":                           flattenPodStatusAttr(p.Status),
		})
	}
	return result
}
