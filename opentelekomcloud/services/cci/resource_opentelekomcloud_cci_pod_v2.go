package cci

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/cci/v2/pod"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceCCIPodV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCciPodV2Create,
		ReadContext:   resourceCciPodV2Read,
		UpdateContext: resourceCciPodV2Update,
		DeleteContext: resourceCciPodV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"active_deadline_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"affinity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"required_during_scheduling_ignored_during_execution": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"node_selector_terms": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"match_expressions": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"operator": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																		"values": {
																			Type:     schema.TypeList,
																			Optional: true,
																			Elem:     &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"pod_anti_affinity": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"preferred_during_scheduling_ignored_during_execution": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"weight": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"pod_affinity_term": {
													Type:     schema.TypeList,
													Required: true,
													MaxItems: 1,
													Elem:     podAffinityTermElem(),
												},
											},
										},
									},
									"required_during_scheduling_ignored_during_execution": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     podAffinityTermElem(),
									},
								},
							},
						},
					},
				},
			},
			"containers": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     podContainerElem(),
			},
			"dns_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nameservers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"options": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"searches": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"dns_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ephemeral_containers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     podEphemeralContainerElem(),
			},
			"host_aliases": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostnames": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"ip": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"image_pull_secrets": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"init_containers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     podContainerElem(),
			},
			"node_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"overhead": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"readiness_gates": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"restart_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scheduler_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_context": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fs_group": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"fs_group_change_policy": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"run_as_group": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"run_as_non_root": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"run_as_user": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"supplemental_groups": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"sysctls": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"set_hostname_as_fqdn": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"share_process_namespace": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"volumes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     podVolumeElem(),
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kind": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finalizers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func podAffinityTermElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topology_key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"label_selector": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match_labels": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match_expressions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},
									"values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"namespaces": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func podContainerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"args": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"command": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"env": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"env_from": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_map_ref": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podEnvSourceElem(),
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"secret_ref": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podEnvSourceElem(),
						},
					},
				},
			},
			"lifecycle": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"post_start": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podLifecycleHandlerElem(),
						},
						"pre_stop": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podLifecycleHandlerElem(),
						},
					},
				},
			},
			"liveness_probe": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podProbeElem(),
			},
			"readiness_probe": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podProbeElem(),
			},
			"startup_probe": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     podProbeElem(),
			},
			"ports": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"limits": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"requests": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"security_context": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capabilities": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"add": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"drop": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"proc_mount": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"read_only_root_file_system": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"run_as_group": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"run_as_non_root": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"run_as_user": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"stdin": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"stdin_once": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"termination_message_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"termination_message_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tty": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"working_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_mounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"sub_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sub_path_expr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"extend_path_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func podEphemeralContainerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"args": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"command": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"env": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"env_from": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_map_ref": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podEnvSourceElem(),
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"secret_ref": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem:     podEnvSourceElem(),
						},
					},
				},
			},
			"security_context": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capabilities": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"add": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"drop": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"proc_mount": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"read_only_root_file_system": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"run_as_group": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"run_as_non_root": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"run_as_user": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"stdin": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"stdin_once": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"target_container_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"termination_message_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"termination_message_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tty": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"volume_mounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"sub_path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"sub_path_expr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"extend_path_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"working_dir": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func podEnvSourceElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"optional": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func podLifecycleHandlerElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"http_get": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"http_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func podProbeElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exec": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"http_get": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"http_headers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scheme": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"failure_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"initial_delay_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"success_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"termination_grace_period_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func podVolumeElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_map": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"default_mode": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"items": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     podKeyToPathElem(),
						},
						"optional": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"nfs": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"server": {
							Type:     schema.TypeString,
							Required: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"persistent_volume_claim": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"claim_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"read_only": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"secret": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"secret_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"default_mode": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"items": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     podKeyToPathElem(),
						},
						"optional": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"projected": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_mode": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"sources": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_map": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"items": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem:     podKeyToPathElem(),
												},
												"optional": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"downward_api": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"items": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"path": {
																Type:     schema.TypeString,
																Required: true,
															},
															"mode": {
																Type:     schema.TypeInt,
																Optional: true,
																Computed: true,
															},
															"field_ref": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																MaxItems: 1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"api_version": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"field_path": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																	},
																},
															},
															"resource_field_ref": {
																Type:     schema.TypeList,
																Optional: true,
																Computed: true,
																MaxItems: 1,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"container_name": {
																			Type:     schema.TypeString,
																			Optional: true,
																			Computed: true,
																		},
																		"resource": {
																			Type:     schema.TypeString,
																			Required: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
									"secret": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Optional: true,
													Computed: true,
												},
												"items": {
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem:     podKeyToPathElem(),
												},
												"optional": {
													Type:     schema.TypeBool,
													Optional: true,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func podKeyToPathElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// CRUD

func resourceCciPodV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	createOpts := pod.CreateOpts{
		Pod: buildPodFromResourceData(d),
	}

	resp, err := pod.Create(client, ns, createOpts)
	if err != nil {
		return diag.Errorf("error creating CCI v2 Pod: %s", err)
	}

	if resp.Metadata.Namespace == "" || resp.Metadata.Name == "" {
		return diag.Errorf("unable to find namespace or CCI v2 Pod name from API response")
	}
	d.SetId(resp.Metadata.Namespace + "/" + resp.Metadata.Name)

	err = waitForPodRunning(ctx, client, resp.Metadata.Namespace, resp.Metadata.Name, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCciPodV2Read(ctx, d, meta)
}

func resourceCciPodV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	resp, err := pod.Get(client, ns, name)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying CCI v2 Pod")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("namespace", resp.Metadata.Namespace),
		d.Set("name", resp.Metadata.Name),
		d.Set("api_version", resp.APIVersion),
		d.Set("kind", resp.Kind),
		d.Set("annotations", resp.Metadata.Annotations),
		d.Set("labels", resp.Metadata.Labels),
		d.Set("creation_timestamp", resp.Metadata.CreationTimestamp),
		d.Set("resource_version", resp.Metadata.ResourceVersion),
		d.Set("finalizers", resp.Metadata.Finalizers),
		d.Set("uid", resp.Metadata.UID),
		d.Set("active_deadline_seconds", derefInt64(resp.Spec.ActiveDeadlineSeconds)),
		d.Set("affinity", flattenAffinity(resp.Spec.Affinity)),
		d.Set("containers", flattenContainers(resp.Spec.Containers)),
		d.Set("dns_config", flattenDNSConfig(resp.Spec.DNSConfig)),
		d.Set("dns_policy", resp.Spec.DNSPolicy),
		d.Set("ephemeral_containers", flattenEphemeralContainers(resp.Spec.EphemeralContainers)),
		d.Set("host_aliases", flattenHostAliases(resp.Spec.HostAliases)),
		d.Set("hostname", resp.Spec.Hostname),
		d.Set("image_pull_secrets", flattenLocalObjectReferences(resp.Spec.ImagePullSecrets)),
		d.Set("init_containers", flattenContainers(resp.Spec.InitContainers)),
		d.Set("node_name", resp.Spec.NodeName),
		d.Set("overhead", resp.Spec.Overhead),
		d.Set("readiness_gates", flattenReadinessGates(resp.Spec.ReadinessGates)),
		d.Set("restart_policy", resp.Spec.RestartPolicy),
		d.Set("scheduler_name", resp.Spec.SchedulerName),
		d.Set("security_context", flattenPodSecurityContext(resp.Spec.SecurityContext)),
		d.Set("set_hostname_as_fqdn", derefBool(resp.Spec.SetHostnameAsFQDN)),
		d.Set("share_process_namespace", derefBool(resp.Spec.ShareProcessNamespace)),
		d.Set("termination_grace_period_seconds", derefInt64(resp.Spec.TerminationGracePeriodSeconds)),
		d.Set("volumes", flattenVolumes(resp.Spec.Volumes)),
		d.Set("status", flattenPodStatusAttr(resp.Status)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCciPodV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	p := buildPodFromResourceData(d)
	p.Metadata.ResourceVersion = d.Get("resource_version").(string)

	updateOpts := pod.UpdateOpts{
		Pod: p,
	}

	_, err = pod.Update(client, ns, name, updateOpts)
	if err != nil {
		return diag.Errorf("error updating CCI v2 Pod: %s", err)
	}

	return resourceCciPodV2Read(ctx, d, meta)
}

func resourceCciPodV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, cciClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.CciV2Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationCciClient, err)
	}

	ns := d.Get("namespace").(string)
	name := d.Get("name").(string)

	_, err = pod.Delete(client, pod.DeleteOpts{
		NameSpace: ns,
		PodName:   name,
		Body:      pod.DeleteBody{},
	})
	if err != nil {
		return diag.Errorf("error deleting CCI v2 Pod: %s", err)
	}

	err = waitForPodDeleted(ctx, client, ns, name, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Waiters

func waitForPodRunning(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Running"},
		Refresh:      refreshPodStatus(client, ns, name),
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CCI v2 Pod to become running: %s", err)
	}
	return nil
}

func waitForPodDeleted(ctx context.Context, client *golangsdk.ServiceClient, ns, name string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Deleted"},
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
		Delay:        10 * time.Second,
		Refresh: func() (interface{}, string, error) {
			_, err := pod.Get(client, ns, name)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					log.Printf("[DEBUG] successfully deleted CCI v2 Pod: %s", name)
					return "", "Deleted", nil
				}
				return nil, "ERROR", err
			}
			return nil, "Pending", nil
		},
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CCI v2 Pod to be deleted: %s", err)
	}
	return nil
}

func refreshPodStatus(client *golangsdk.ServiceClient, ns, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := pod.Get(client, ns, name)
		if err != nil {
			return nil, "ERROR", err
		}
		if resp.Status.Phase == "Running" {
			return resp, "Running", nil
		}
		return resp, "Pending", nil
	}
}

// Build helpers: ResourceData -> structs

func buildPodFromResourceData(d *schema.ResourceData) pod.Pod {
	p := pod.Pod{
		APIVersion: "cci/v2",
		Kind:       "Pod",
		Metadata: pod.ObjectMeta{
			Name:        d.Get("name").(string),
			Namespace:   d.Get("namespace").(string),
			Annotations: expandStringMap(d.Get("annotations")),
			Labels:      expandStringMap(d.Get("labels")),
		},
		Spec: pod.PodSpec{
			Containers:     expandContainers(d.Get("containers").([]interface{})),
			InitContainers: expandContainers(d.Get("init_containers").([]interface{})),
			DNSPolicy:      d.Get("dns_policy").(string),
			Hostname:       d.Get("hostname").(string),
			NodeName:       d.Get("node_name").(string),
			RestartPolicy:  d.Get("restart_policy").(string),
			SchedulerName:  d.Get("scheduler_name").(string),
		},
	}

	if v, ok := d.GetOk("active_deadline_seconds"); ok {
		val := int64(v.(int))
		p.Spec.ActiveDeadlineSeconds = &val
	}
	if v, ok := d.GetOk("termination_grace_period_seconds"); ok {
		val := int64(v.(int))
		p.Spec.TerminationGracePeriodSeconds = &val
	}
	if v, ok := d.GetOk("set_hostname_as_fqdn"); ok {
		val := v.(bool)
		p.Spec.SetHostnameAsFQDN = &val
	}
	if v, ok := d.GetOk("share_process_namespace"); ok {
		val := v.(bool)
		p.Spec.ShareProcessNamespace = &val
	}
	if v, ok := d.GetOk("overhead"); ok {
		p.Spec.Overhead = expandStringMap(v)
	}
	if v, ok := d.GetOk("affinity"); ok {
		p.Spec.Affinity = expandAffinity(v.([]interface{}))
	}
	if v, ok := d.GetOk("dns_config"); ok {
		p.Spec.DNSConfig = expandDNSConfig(v.([]interface{}))
	}
	if v, ok := d.GetOk("ephemeral_containers"); ok {
		p.Spec.EphemeralContainers = expandEphemeralContainers(v.([]interface{}))
	}
	if v, ok := d.GetOk("host_aliases"); ok {
		p.Spec.HostAliases = expandHostAliases(v.([]interface{}))
	}
	if v, ok := d.GetOk("image_pull_secrets"); ok {
		p.Spec.ImagePullSecrets = expandLocalObjectReferences(v.([]interface{}))
	}
	if v, ok := d.GetOk("readiness_gates"); ok {
		p.Spec.ReadinessGates = expandReadinessGates(v.([]interface{}))
	}
	if v, ok := d.GetOk("security_context"); ok {
		p.Spec.SecurityContext = expandPodSecurityContext(v.([]interface{}))
	}
	if v, ok := d.GetOk("volumes"); ok {
		p.Spec.Volumes = expandVolumes(v.([]interface{}))
	}

	return p
}

func expandContainers(raw []interface{}) []pod.Container {
	if len(raw) == 0 {
		return nil
	}
	containers := make([]pod.Container, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		c := pod.Container{
			Name:                     m["name"].(string),
			Image:                    m["image"].(string),
			Args:                     expandStringList(m["args"].([]interface{})),
			Command:                  expandStringList(m["command"].([]interface{})),
			Stdin:                    m["stdin"].(bool),
			StdinOnce:                m["stdin_once"].(bool),
			TTY:                      m["tty"].(bool),
			WorkingDir:               m["working_dir"].(string),
			TerminationMessagePath:   m["termination_message_path"].(string),
			TerminationMessagePolicy: m["termination_message_policy"].(string),
		}
		if env, ok := m["env"]; ok {
			c.Env = expandEnvVars(env.([]interface{}))
		}
		if envFrom, ok := m["env_from"]; ok {
			c.EnvFrom = expandEnvFromSources(envFrom.([]interface{}))
		}
		if ports, ok := m["ports"]; ok {
			c.Ports = expandContainerPorts(ports.([]interface{}))
		}
		if res, ok := m["resources"]; ok {
			c.Resources = expandResourceRequirements(res.([]interface{}))
		}
		if sc, ok := m["security_context"]; ok {
			c.SecurityContext = expandSecurityContext(sc.([]interface{}))
		}
		if lc, ok := m["lifecycle"]; ok {
			c.Lifecycle = expandLifecycle(lc.([]interface{}))
		}
		if lp, ok := m["liveness_probe"]; ok {
			c.LivenessProbe = expandProbe(lp.([]interface{}))
		}
		if rp, ok := m["readiness_probe"]; ok {
			c.ReadinessProbe = expandProbe(rp.([]interface{}))
		}
		if sp, ok := m["startup_probe"]; ok {
			c.StartupProbe = expandProbe(sp.([]interface{}))
		}
		if vm, ok := m["volume_mounts"]; ok {
			c.VolumeMounts = expandVolumeMounts(vm.([]interface{}))
		}
		containers[i] = c
	}
	return containers
}

func expandEphemeralContainers(raw []interface{}) []pod.EphemeralContainer {
	if len(raw) == 0 {
		return nil
	}
	containers := make([]pod.EphemeralContainer, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		c := pod.EphemeralContainer{
			Name:                     m["name"].(string),
			Image:                    m["image"].(string),
			Args:                     expandStringList(m["args"].([]interface{})),
			Command:                  expandStringList(m["command"].([]interface{})),
			Stdin:                    m["stdin"].(bool),
			StdinOnce:                m["stdin_once"].(bool),
			TTY:                      m["tty"].(bool),
			WorkingDir:               m["working_dir"].(string),
			TargetContainerName:      m["target_container_name"].(string),
			TerminationMessagePath:   m["termination_message_path"].(string),
			TerminationMessagePolicy: m["termination_message_policy"].(string),
		}
		if env, ok := m["env"]; ok {
			c.Env = expandEnvVars(env.([]interface{}))
		}
		if envFrom, ok := m["env_from"]; ok {
			c.EnvFrom = expandEnvFromSources(envFrom.([]interface{}))
		}
		if sc, ok := m["security_context"]; ok {
			c.SecurityContext = expandSecurityContext(sc.([]interface{}))
		}
		if vm, ok := m["volume_mounts"]; ok {
			c.VolumeMounts = expandVolumeMounts(vm.([]interface{}))
		}
		containers[i] = c
	}
	return containers
}

func expandEnvVars(raw []interface{}) []pod.EnvVar {
	if len(raw) == 0 {
		return nil
	}
	envs := make([]pod.EnvVar, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		envs[i] = pod.EnvVar{
			Name:  m["name"].(string),
			Value: m["value"].(string),
		}
	}
	return envs
}

func expandEnvFromSources(raw []interface{}) []pod.EnvFromSource {
	if len(raw) == 0 {
		return nil
	}
	sources := make([]pod.EnvFromSource, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		s := pod.EnvFromSource{
			Prefix: m["prefix"].(string),
		}
		if cmRef, ok := m["config_map_ref"]; ok {
			list := cmRef.([]interface{})
			if len(list) > 0 && list[0] != nil {
				ref := list[0].(map[string]interface{})
				s.ConfigMapRef = &pod.ConfigMapEnvSource{
					Name: ref["name"].(string),
				}
				if opt, ok := ref["optional"]; ok {
					val := opt.(bool)
					s.ConfigMapRef.Optional = &val
				}
			}
		}
		if sRef, ok := m["secret_ref"]; ok {
			list := sRef.([]interface{})
			if len(list) > 0 && list[0] != nil {
				ref := list[0].(map[string]interface{})
				s.SecretRef = &pod.SecretEnvSource{
					Name: ref["name"].(string),
				}
				if opt, ok := ref["optional"]; ok {
					val := opt.(bool)
					s.SecretRef.Optional = &val
				}
			}
		}
		sources[i] = s
	}
	return sources
}

func expandContainerPorts(raw []interface{}) []pod.ContainerPort {
	if len(raw) == 0 {
		return nil
	}
	ports := make([]pod.ContainerPort, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		ports[i] = pod.ContainerPort{
			ContainerPort: int32(m["container_port"].(int)),
			Name:          m["name"].(string),
			Protocol:      m["protocol"].(string),
		}
	}
	return ports
}

func expandResourceRequirements(raw []interface{}) pod.ResourceRequirements {
	if len(raw) == 0 || raw[0] == nil {
		return pod.ResourceRequirements{}
	}
	m := raw[0].(map[string]interface{})
	return pod.ResourceRequirements{
		Limits:   expandStringMap(m["limits"]),
		Requests: expandStringMap(m["requests"]),
	}
}

func expandSecurityContext(raw []interface{}) *pod.SecurityContext {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	sc := &pod.SecurityContext{
		ProcMount: m["proc_mount"].(string),
	}
	if v, ok := m["read_only_root_file_system"]; ok && v.(bool) {
		val := v.(bool)
		sc.ReadOnlyRootFilesystem = &val
	}
	if v, ok := m["run_as_group"]; ok && v.(int) != 0 {
		val := int64(v.(int))
		sc.RunAsGroup = &val
	}
	if v, ok := m["run_as_non_root"]; ok && v.(bool) {
		val := v.(bool)
		sc.RunAsNonRoot = &val
	}
	if v, ok := m["run_as_user"]; ok && v.(int) != 0 {
		val := int64(v.(int))
		sc.RunAsUser = &val
	}
	if caps, ok := m["capabilities"]; ok {
		list := caps.([]interface{})
		if len(list) > 0 && list[0] != nil {
			cm := list[0].(map[string]interface{})
			sc.Capabilities = &pod.Capabilities{
				Add:  expandStringList(cm["add"].([]interface{})),
				Drop: expandStringList(cm["drop"].([]interface{})),
			}
		}
	}
	return sc
}

func expandLifecycle(raw []interface{}) *pod.Lifecycle {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	lc := &pod.Lifecycle{}
	if ps, ok := m["post_start"]; ok {
		lc.PostStart = expandLifecycleHandler(ps.([]interface{}))
	}
	if ps, ok := m["pre_stop"]; ok {
		lc.PreStop = expandLifecycleHandler(ps.([]interface{}))
	}
	return lc
}

func expandLifecycleHandler(raw []interface{}) *pod.LifecycleHandler {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	h := &pod.LifecycleHandler{}
	if exec, ok := m["exec"]; ok {
		list := exec.([]interface{})
		if len(list) > 0 && list[0] != nil {
			em := list[0].(map[string]interface{})
			h.Exec = &pod.ExecAction{
				Command: expandStringList(em["command"].([]interface{})),
			}
		}
	}
	if httpGet, ok := m["http_get"]; ok {
		list := httpGet.([]interface{})
		if len(list) > 0 && list[0] != nil {
			h.HTTPGet = expandHTTPGetAction(list[0].(map[string]interface{}))
		}
	}
	return h
}

func expandHTTPGetAction(m map[string]interface{}) *pod.HTTPGetAction {
	action := &pod.HTTPGetAction{
		Host:   m["host"].(string),
		Path:   m["path"].(string),
		Port:   m["port"].(string),
		Scheme: m["scheme"].(string),
	}
	if headers, ok := m["http_headers"]; ok {
		for _, h := range headers.([]interface{}) {
			hm := h.(map[string]interface{})
			action.HTTPHeaders = append(action.HTTPHeaders, pod.HTTPHeader{
				Name:  hm["name"].(string),
				Value: hm["value"].(string),
			})
		}
	}
	return action
}

func expandProbe(raw []interface{}) *pod.Probe {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	p := &pod.Probe{
		FailureThreshold:    int32(m["failure_threshold"].(int)),
		InitialDelaySeconds: int32(m["initial_delay_seconds"].(int)),
		PeriodSeconds:       int32(m["period_seconds"].(int)),
		SuccessThreshold:    int32(m["success_threshold"].(int)),
	}
	if v, ok := m["termination_grace_period_seconds"]; ok && v.(int) != 0 {
		val := int64(v.(int))
		p.TerminationGracePeriodSeconds = &val
	}
	if exec, ok := m["exec"]; ok {
		list := exec.([]interface{})
		if len(list) > 0 && list[0] != nil {
			em := list[0].(map[string]interface{})
			p.Exec = &pod.ExecAction{
				Command: expandStringList(em["command"].([]interface{})),
			}
		}
	}
	if httpGet, ok := m["http_get"]; ok {
		list := httpGet.([]interface{})
		if len(list) > 0 && list[0] != nil {
			p.HTTPGet = expandHTTPGetAction(list[0].(map[string]interface{}))
		}
	}
	return p
}

func expandVolumeMounts(raw []interface{}) []pod.VolumeMount {
	if len(raw) == 0 {
		return nil
	}
	mounts := make([]pod.VolumeMount, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		mounts[i] = pod.VolumeMount{
			MountPath:      m["mount_path"].(string),
			Name:           m["name"].(string),
			ReadOnly:       m["read_only"].(bool),
			SubPath:        m["sub_path"].(string),
			SubPathExpr:    m["sub_path_expr"].(string),
			ExtendPathMode: m["extend_path_mode"].(string),
		}
	}
	return mounts
}

func expandAffinity(raw []interface{}) *pod.Affinity {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	a := &pod.Affinity{}

	if na, ok := m["node_affinity"]; ok {
		list := na.([]interface{})
		if len(list) > 0 && list[0] != nil {
			nam := list[0].(map[string]interface{})
			a.NodeAffinity = &pod.NodeAffinity{}
			if req, ok := nam["required_during_scheduling_ignored_during_execution"]; ok {
				reqList := req.([]interface{})
				if len(reqList) > 0 && reqList[0] != nil {
					rm := reqList[0].(map[string]interface{})
					ns := &pod.NodeSelector{}
					if terms, ok := rm["node_selector_terms"]; ok {
						for _, t := range terms.([]interface{}) {
							tm := t.(map[string]interface{})
							term := pod.NodeSelectorTerm{}
							if me, ok := tm["match_expressions"]; ok {
								for _, e := range me.([]interface{}) {
									em := e.(map[string]interface{})
									term.MatchExpressions = append(term.MatchExpressions, pod.NodeSelectorRequirement{
										Key:      em["key"].(string),
										Operator: em["operator"].(string),
										Values:   expandStringList(em["values"].([]interface{})),
									})
								}
							}
							ns.NodeSelectorTerms = append(ns.NodeSelectorTerms, term)
						}
					}
					a.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = ns
				}
			}
		}
	}

	if paa, ok := m["pod_anti_affinity"]; ok {
		list := paa.([]interface{})
		if len(list) > 0 && list[0] != nil {
			pam := list[0].(map[string]interface{})
			a.PodAntiAffinity = &pod.PodAntiAffinity{}
			if pref, ok := pam["preferred_during_scheduling_ignored_during_execution"]; ok {
				for _, p := range pref.([]interface{}) {
					pm := p.(map[string]interface{})
					wt := pod.WeightedPodAffinityTerm{
						Weight: int32(pm["weight"].(int)),
					}
					if pat, ok := pm["pod_affinity_term"]; ok {
						patList := pat.([]interface{})
						if len(patList) > 0 && patList[0] != nil {
							wt.PodAffinityTerm = expandPodAffinityTerm(patList[0].(map[string]interface{}))
						}
					}
					a.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
						a.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution, wt)
				}
			}
			if req, ok := pam["required_during_scheduling_ignored_during_execution"]; ok {
				for _, r := range req.([]interface{}) {
					rm := r.(map[string]interface{})
					a.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
						a.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
						expandPodAffinityTerm(rm))
				}
			}
		}
	}

	return a
}

func expandPodAffinityTerm(m map[string]interface{}) pod.PodAffinityTerm {
	t := pod.PodAffinityTerm{
		TopologyKey: m["topology_key"].(string),
		Namespaces:  expandStringList(m["namespaces"].([]interface{})),
	}
	if ls, ok := m["label_selector"]; ok {
		list := ls.([]interface{})
		if len(list) > 0 && list[0] != nil {
			t.LabelSelector = expandLabelSelector(list[0].(map[string]interface{}))
		}
	}
	return t
}

func expandLabelSelector(m map[string]interface{}) *pod.LabelSelector {
	ls := &pod.LabelSelector{}
	if ml, ok := m["match_labels"]; ok {
		ls.MatchLabels = expandStringMap(ml)
	}
	if me, ok := m["match_expressions"]; ok {
		for _, e := range me.([]interface{}) {
			em := e.(map[string]interface{})
			ls.MatchExpressions = append(ls.MatchExpressions, pod.LabelSelectorRequirement{
				Key:      em["key"].(string),
				Operator: em["operator"].(string),
				Values:   expandStringList(em["values"].([]interface{})),
			})
		}
	}
	return ls
}

func expandDNSConfig(raw []interface{}) *pod.PodDNSConfig {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	dnsCfg := &pod.PodDNSConfig{
		Nameservers: expandStringList(m["nameservers"].([]interface{})),
		Searches:    expandStringList(m["searches"].([]interface{})),
	}
	if opts, ok := m["options"]; ok {
		for _, o := range opts.([]interface{}) {
			om := o.(map[string]interface{})
			dnsCfg.Options = append(dnsCfg.Options, pod.PodDNSConfigOption{
				Name:  om["name"].(string),
				Value: om["value"].(string),
			})
		}
	}
	return dnsCfg
}

func expandHostAliases(raw []interface{}) []pod.HostAlias {
	if len(raw) == 0 {
		return nil
	}
	aliases := make([]pod.HostAlias, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		aliases[i] = pod.HostAlias{
			IP:        m["ip"].(string),
			Hostnames: expandStringList(m["hostnames"].([]interface{})),
		}
	}
	return aliases
}

func expandLocalObjectReferences(raw []interface{}) []pod.LocalObjectReference {
	if len(raw) == 0 {
		return nil
	}
	refs := make([]pod.LocalObjectReference, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		refs[i] = pod.LocalObjectReference{Name: m["name"].(string)}
	}
	return refs
}

func expandReadinessGates(raw []interface{}) []pod.PodReadinessGate {
	if len(raw) == 0 {
		return nil
	}
	gates := make([]pod.PodReadinessGate, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		gates[i] = pod.PodReadinessGate{ConditionType: m["condition_type"].(string)}
	}
	return gates
}

func expandPodSecurityContext(raw []interface{}) *pod.PodSecurityContext {
	if len(raw) == 0 || raw[0] == nil {
		return nil
	}
	m := raw[0].(map[string]interface{})
	sc := &pod.PodSecurityContext{
		FSGroupChangePolicy: m["fs_group_change_policy"].(string),
	}
	if v, ok := m["fs_group"]; ok && v.(int) != 0 {
		val := int64(v.(int))
		sc.FSGroup = &val
	}
	if v, ok := m["run_as_group"]; ok && v.(int) != 0 {
		val := int64(v.(int))
		sc.RunAsGroup = &val
	}
	if v, ok := m["run_as_non_root"]; ok && v.(bool) {
		val := v.(bool)
		sc.RunAsNonRoot = &val
	}
	if v, ok := m["run_as_user"]; ok && v.(int) != 0 {
		val := int64(v.(int))
		sc.RunAsUser = &val
	}
	if v, ok := m["supplemental_groups"]; ok {
		for _, g := range v.([]interface{}) {
			sc.SupplementalGroups = append(sc.SupplementalGroups, int64(g.(int)))
		}
	}
	if v, ok := m["sysctls"]; ok {
		for _, s := range v.([]interface{}) {
			sm := s.(map[string]interface{})
			sc.Sysctls = append(sc.Sysctls, pod.Sysctl{
				Name:  sm["name"].(string),
				Value: sm["value"].(string),
			})
		}
	}
	return sc
}

func expandVolumes(raw []interface{}) []pod.Volume {
	if len(raw) == 0 {
		return nil
	}
	volumes := make([]pod.Volume, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		vol := pod.Volume{
			Name: m["name"].(string),
		}
		if cm, ok := m["config_map"]; ok {
			list := cm.([]interface{})
			if len(list) > 0 && list[0] != nil {
				cmm := list[0].(map[string]interface{})
				vol.ConfigMap = &pod.ConfigMapVolumeSource{
					Name: cmm["name"].(string),
				}
				if dm, ok := cmm["default_mode"]; ok && dm.(int) != 0 {
					val := int32(dm.(int))
					vol.ConfigMap.DefaultMode = &val
				}
				if opt, ok := cmm["optional"]; ok {
					val := opt.(bool)
					vol.ConfigMap.Optional = &val
				}
				if items, ok := cmm["items"]; ok {
					vol.ConfigMap.Items = expandKeyToPath(items.([]interface{}))
				}
			}
		}
		if nfs, ok := m["nfs"]; ok {
			list := nfs.([]interface{})
			if len(list) > 0 && list[0] != nil {
				nm := list[0].(map[string]interface{})
				vol.NFS = &pod.NFSVolumeSource{
					Path:     nm["path"].(string),
					Server:   nm["server"].(string),
					ReadOnly: nm["read_only"].(bool),
				}
			}
		}
		if pvc, ok := m["persistent_volume_claim"]; ok {
			list := pvc.([]interface{})
			if len(list) > 0 && list[0] != nil {
				pm := list[0].(map[string]interface{})
				vol.PersistentVolumeClaim = &pod.PersistentVolumeClaimVolumeSource{
					ClaimName: pm["claim_name"].(string),
					ReadOnly:  pm["read_only"].(bool),
				}
			}
		}
		if sec, ok := m["secret"]; ok {
			list := sec.([]interface{})
			if len(list) > 0 && list[0] != nil {
				sm := list[0].(map[string]interface{})
				vol.Secret = &pod.SecretVolumeSource{
					SecretName: sm["secret_name"].(string),
				}
				if dm, ok := sm["default_mode"]; ok && dm.(int) != 0 {
					val := int32(dm.(int))
					vol.Secret.DefaultMode = &val
				}
				if opt, ok := sm["optional"]; ok {
					val := opt.(bool)
					vol.Secret.Optional = &val
				}
				if items, ok := sm["items"]; ok {
					vol.Secret.Items = expandKeyToPath(items.([]interface{}))
				}
			}
		}
		if proj, ok := m["projected"]; ok {
			list := proj.([]interface{})
			if len(list) > 0 && list[0] != nil {
				pm := list[0].(map[string]interface{})
				vol.Projected = &pod.ProjectedVolumeSource{}
				if dm, ok := pm["default_mode"]; ok && dm.(int) != 0 {
					val := int32(dm.(int))
					vol.Projected.DefaultMode = &val
				}
				if sources, ok := pm["sources"]; ok {
					vol.Projected.Sources = expandVolumeProjections(sources.([]interface{}))
				}
			}
		}
		volumes[i] = vol
	}
	return volumes
}

func expandVolumeProjections(raw []interface{}) []pod.VolumeProjection {
	if len(raw) == 0 {
		return nil
	}
	projections := make([]pod.VolumeProjection, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		vp := pod.VolumeProjection{}
		if cm, ok := m["config_map"]; ok {
			list := cm.([]interface{})
			if len(list) > 0 && list[0] != nil {
				cmm := list[0].(map[string]interface{})
				vp.ConfigMap = &pod.ConfigMapProjection{
					Name: cmm["name"].(string),
				}
				if opt, ok := cmm["optional"]; ok {
					val := opt.(bool)
					vp.ConfigMap.Optional = &val
				}
				if items, ok := cmm["items"]; ok {
					vp.ConfigMap.Items = expandKeyToPath(items.([]interface{}))
				}
			}
		}
		if da, ok := m["downward_api"]; ok {
			list := da.([]interface{})
			if len(list) > 0 && list[0] != nil {
				dam := list[0].(map[string]interface{})
				vp.DownwardAPI = &pod.DownwardAPIProjection{}
				if items, ok := dam["items"]; ok {
					vp.DownwardAPI.Items = expandDownwardAPIVolumeFiles(items.([]interface{}))
				}
			}
		}
		if sec, ok := m["secret"]; ok {
			list := sec.([]interface{})
			if len(list) > 0 && list[0] != nil {
				sm := list[0].(map[string]interface{})
				vp.Secret = &pod.SecretProjection{
					Name: sm["name"].(string),
				}
				if opt, ok := sm["optional"]; ok {
					val := opt.(bool)
					vp.Secret.Optional = &val
				}
				if items, ok := sm["items"]; ok {
					vp.Secret.Items = expandKeyToPath(items.([]interface{}))
				}
			}
		}
		projections[i] = vp
	}
	return projections
}

func expandDownwardAPIVolumeFiles(raw []interface{}) []pod.DownwardAPIVolumeFile {
	if len(raw) == 0 {
		return nil
	}
	files := make([]pod.DownwardAPIVolumeFile, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		f := pod.DownwardAPIVolumeFile{
			Path: m["path"].(string),
		}
		if mode, ok := m["mode"]; ok && mode.(int) != 0 {
			val := int32(mode.(int))
			f.Mode = &val
		}
		if fr, ok := m["field_ref"]; ok {
			list := fr.([]interface{})
			if len(list) > 0 && list[0] != nil {
				frm := list[0].(map[string]interface{})
				f.FieldRef = &pod.ObjectFieldSelector{
					APIVersion: frm["api_version"].(string),
					FieldPath:  frm["field_path"].(string),
				}
			}
		}
		if rfr, ok := m["resource_field_ref"]; ok {
			list := rfr.([]interface{})
			if len(list) > 0 && list[0] != nil {
				rfrm := list[0].(map[string]interface{})
				f.ResourceFieldRef = &pod.ResourceFieldSelector{
					ContainerName: rfrm["container_name"].(string),
					Resource:      rfrm["resource"].(string),
				}
			}
		}
		files[i] = f
	}
	return files
}

func expandKeyToPath(raw []interface{}) []pod.KeyToPath {
	if len(raw) == 0 {
		return nil
	}
	items := make([]pod.KeyToPath, len(raw))
	for i, v := range raw {
		m := v.(map[string]interface{})
		ktp := pod.KeyToPath{
			Key:  m["key"].(string),
			Path: m["path"].(string),
		}
		if mode, ok := m["mode"]; ok && mode.(int) != 0 {
			val := int32(mode.(int))
			ktp.Mode = &val
		}
		items[i] = ktp
	}
	return items
}

// Flatten helpers: structs -> ResourceData

func flattenAffinity(a *pod.Affinity) []map[string]interface{} {
	if a == nil {
		return nil
	}
	result := map[string]interface{}{
		"node_affinity":     flattenNodeAffinity(a.NodeAffinity),
		"pod_anti_affinity": flattenPodAntiAffinity(a.PodAntiAffinity),
	}
	return []map[string]interface{}{result}
}

func flattenNodeAffinity(na *pod.NodeAffinity) []map[string]interface{} {
	if na == nil {
		return nil
	}
	result := map[string]interface{}{
		"required_during_scheduling_ignored_during_execution": flattenNodeSelector(na.RequiredDuringSchedulingIgnoredDuringExecution),
	}
	return []map[string]interface{}{result}
}

func flattenNodeSelector(ns *pod.NodeSelector) []map[string]interface{} {
	if ns == nil {
		return nil
	}
	terms := make([]map[string]interface{}, len(ns.NodeSelectorTerms))
	for i, t := range ns.NodeSelectorTerms {
		exprs := make([]map[string]interface{}, len(t.MatchExpressions))
		for j, e := range t.MatchExpressions {
			exprs[j] = map[string]interface{}{
				"key":      e.Key,
				"operator": e.Operator,
				"values":   e.Values,
			}
		}
		terms[i] = map[string]interface{}{
			"match_expressions": exprs,
		}
	}
	return []map[string]interface{}{
		{"node_selector_terms": terms},
	}
}

func flattenPodAntiAffinity(paa *pod.PodAntiAffinity) []map[string]interface{} {
	if paa == nil {
		return nil
	}

	preferred := make([]map[string]interface{}, len(paa.PreferredDuringSchedulingIgnoredDuringExecution))
	for i, p := range paa.PreferredDuringSchedulingIgnoredDuringExecution {
		preferred[i] = map[string]interface{}{
			"weight":            p.Weight,
			"pod_affinity_term": flattenPodAffinityTermSingle(p.PodAffinityTerm),
		}
	}

	required := make([]map[string]interface{}, len(paa.RequiredDuringSchedulingIgnoredDuringExecution))
	for i, r := range paa.RequiredDuringSchedulingIgnoredDuringExecution {
		required[i] = flattenPodAffinityTermMap(r)
	}

	return []map[string]interface{}{
		{
			"preferred_during_scheduling_ignored_during_execution": preferred,
			"required_during_scheduling_ignored_during_execution":  required,
		},
	}
}

func flattenPodAffinityTermSingle(t pod.PodAffinityTerm) []map[string]interface{} {
	return []map[string]interface{}{flattenPodAffinityTermMap(t)}
}

func flattenPodAffinityTermMap(t pod.PodAffinityTerm) map[string]interface{} {
	result := map[string]interface{}{
		"topology_key": t.TopologyKey,
		"namespaces":   t.Namespaces,
	}
	if t.LabelSelector != nil {
		result["label_selector"] = flattenLabelSelector(t.LabelSelector)
	}
	return result
}

func flattenLabelSelector(ls *pod.LabelSelector) []map[string]interface{} {
	if ls == nil {
		return nil
	}
	exprs := make([]map[string]interface{}, len(ls.MatchExpressions))
	for i, e := range ls.MatchExpressions {
		exprs[i] = map[string]interface{}{
			"key":      e.Key,
			"operator": e.Operator,
			"values":   e.Values,
		}
	}
	return []map[string]interface{}{
		{
			"match_labels":      ls.MatchLabels,
			"match_expressions": exprs,
		},
	}
}

func flattenContainers(containers []pod.Container) []map[string]interface{} {
	if len(containers) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(containers))
	for i, c := range containers {
		result[i] = map[string]interface{}{
			"name":                       c.Name,
			"image":                      c.Image,
			"args":                       c.Args,
			"command":                    c.Command,
			"env":                        flattenEnvVars(c.Env),
			"env_from":                   flattenEnvFromSources(c.EnvFrom),
			"lifecycle":                  flattenLifecycle(c.Lifecycle),
			"liveness_probe":             flattenProbe(c.LivenessProbe),
			"readiness_probe":            flattenProbe(c.ReadinessProbe),
			"startup_probe":              flattenProbe(c.StartupProbe),
			"ports":                      flattenContainerPorts(c.Ports),
			"resources":                  flattenResourceRequirements(c.Resources),
			"security_context":           flattenSecurityContext(c.SecurityContext),
			"stdin":                      c.Stdin,
			"stdin_once":                 c.StdinOnce,
			"termination_message_path":   c.TerminationMessagePath,
			"termination_message_policy": c.TerminationMessagePolicy,
			"tty":                        c.TTY,
			"working_dir":                c.WorkingDir,
			"volume_mounts":              flattenVolumeMounts(c.VolumeMounts),
		}
	}
	return result
}

func flattenEphemeralContainers(containers []pod.EphemeralContainer) []map[string]interface{} {
	if len(containers) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(containers))
	for i, c := range containers {
		result[i] = map[string]interface{}{
			"name":                       c.Name,
			"image":                      c.Image,
			"args":                       c.Args,
			"command":                    c.Command,
			"env":                        flattenEnvVars(c.Env),
			"env_from":                   flattenEnvFromSources(c.EnvFrom),
			"security_context":           flattenSecurityContext(c.SecurityContext),
			"stdin":                      c.Stdin,
			"stdin_once":                 c.StdinOnce,
			"target_container_name":      c.TargetContainerName,
			"termination_message_path":   c.TerminationMessagePath,
			"termination_message_policy": c.TerminationMessagePolicy,
			"tty":                        c.TTY,
			"working_dir":                c.WorkingDir,
			"volume_mounts":              flattenVolumeMounts(c.VolumeMounts),
		}
	}
	return result
}

func flattenEnvVars(envs []pod.EnvVar) []map[string]interface{} {
	if len(envs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(envs))
	for i, e := range envs {
		result[i] = map[string]interface{}{
			"name":  e.Name,
			"value": e.Value,
		}
	}
	return result
}

func flattenEnvFromSources(sources []pod.EnvFromSource) []map[string]interface{} {
	if len(sources) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(sources))
	for i, s := range sources {
		m := map[string]interface{}{
			"prefix": s.Prefix,
		}
		if s.ConfigMapRef != nil {
			m["config_map_ref"] = []map[string]interface{}{
				{"name": s.ConfigMapRef.Name, "optional": s.ConfigMapRef.Optional},
			}
		}
		if s.SecretRef != nil {
			m["secret_ref"] = []map[string]interface{}{
				{"name": s.SecretRef.Name, "optional": s.SecretRef.Optional},
			}
		}
		result[i] = m
	}
	return result
}

func flattenLifecycle(lc *pod.Lifecycle) []map[string]interface{} {
	if lc == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"post_start": flattenLifecycleHandler(lc.PostStart),
			"pre_stop":   flattenLifecycleHandler(lc.PreStop),
		},
	}
}

func flattenLifecycleHandler(h *pod.LifecycleHandler) []map[string]interface{} {
	if h == nil {
		return nil
	}
	result := map[string]interface{}{}
	if h.Exec != nil {
		result["exec"] = []map[string]interface{}{
			{"command": h.Exec.Command},
		}
	}
	if h.HTTPGet != nil {
		result["http_get"] = []map[string]interface{}{flattenHTTPGetAction(h.HTTPGet)}
	}
	return []map[string]interface{}{result}
}

func flattenHTTPGetAction(a *pod.HTTPGetAction) map[string]interface{} {
	headers := make([]map[string]interface{}, len(a.HTTPHeaders))
	for i, h := range a.HTTPHeaders {
		headers[i] = map[string]interface{}{
			"name":  h.Name,
			"value": h.Value,
		}
	}
	return map[string]interface{}{
		"host":         a.Host,
		"path":         a.Path,
		"port":         a.Port,
		"scheme":       a.Scheme,
		"http_headers": headers,
	}
}

func flattenProbe(p *pod.Probe) []map[string]interface{} {
	if p == nil {
		return nil
	}
	result := map[string]interface{}{
		"failure_threshold":                p.FailureThreshold,
		"initial_delay_seconds":            p.InitialDelaySeconds,
		"period_seconds":                   p.PeriodSeconds,
		"success_threshold":                p.SuccessThreshold,
		"termination_grace_period_seconds": p.TerminationGracePeriodSeconds,
	}
	if p.Exec != nil {
		result["exec"] = []map[string]interface{}{
			{"command": p.Exec.Command},
		}
	}
	if p.HTTPGet != nil {
		result["http_get"] = []map[string]interface{}{flattenHTTPGetAction(p.HTTPGet)}
	}
	return []map[string]interface{}{result}
}

func flattenContainerPorts(ports []pod.ContainerPort) []map[string]interface{} {
	if len(ports) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(ports))
	for i, p := range ports {
		result[i] = map[string]interface{}{
			"container_port": p.ContainerPort,
			"name":           p.Name,
			"protocol":       p.Protocol,
		}
	}
	return result
}

func flattenResourceRequirements(r pod.ResourceRequirements) []map[string]interface{} {
	if r.Limits == nil && r.Requests == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"limits":   r.Limits,
			"requests": r.Requests,
		},
	}
}

func flattenSecurityContext(sc *pod.SecurityContext) []map[string]interface{} {
	if sc == nil {
		return nil
	}
	result := map[string]interface{}{
		"proc_mount":                 sc.ProcMount,
		"read_only_root_file_system": sc.ReadOnlyRootFilesystem,
		"run_as_group":               sc.RunAsGroup,
		"run_as_non_root":            sc.RunAsNonRoot,
		"run_as_user":                sc.RunAsUser,
	}
	if sc.Capabilities != nil {
		result["capabilities"] = []map[string]interface{}{
			{
				"add":  sc.Capabilities.Add,
				"drop": sc.Capabilities.Drop,
			},
		}
	}
	return []map[string]interface{}{result}
}

func flattenVolumeMounts(mounts []pod.VolumeMount) []map[string]interface{} {
	if len(mounts) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(mounts))
	for i, m := range mounts {
		result[i] = map[string]interface{}{
			"mount_path":       m.MountPath,
			"name":             m.Name,
			"read_only":        m.ReadOnly,
			"sub_path":         m.SubPath,
			"sub_path_expr":    m.SubPathExpr,
			"extend_path_mode": m.ExtendPathMode,
		}
	}
	return result
}

func flattenDNSConfig(cfg *pod.PodDNSConfig) []map[string]interface{} {
	if cfg == nil {
		return nil
	}
	opts := make([]map[string]interface{}, len(cfg.Options))
	for i, o := range cfg.Options {
		opts[i] = map[string]interface{}{
			"name":  o.Name,
			"value": o.Value,
		}
	}
	return []map[string]interface{}{
		{
			"nameservers": cfg.Nameservers,
			"options":     opts,
			"searches":    cfg.Searches,
		},
	}
}

func flattenHostAliases(aliases []pod.HostAlias) []map[string]interface{} {
	if len(aliases) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(aliases))
	for i, a := range aliases {
		result[i] = map[string]interface{}{
			"ip":        a.IP,
			"hostnames": a.Hostnames,
		}
	}
	return result
}

func flattenLocalObjectReferences(refs []pod.LocalObjectReference) []map[string]interface{} {
	if len(refs) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(refs))
	for i, r := range refs {
		result[i] = map[string]interface{}{
			"name": r.Name,
		}
	}
	return result
}

func flattenReadinessGates(gates []pod.PodReadinessGate) []map[string]interface{} {
	if len(gates) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(gates))
	for i, g := range gates {
		result[i] = map[string]interface{}{
			"condition_type": g.ConditionType,
		}
	}
	return result
}

func flattenPodSecurityContext(sc *pod.PodSecurityContext) []map[string]interface{} {
	if sc == nil {
		return nil
	}
	sysctls := make([]map[string]interface{}, len(sc.Sysctls))
	for i, s := range sc.Sysctls {
		sysctls[i] = map[string]interface{}{
			"name":  s.Name,
			"value": s.Value,
		}
	}
	return []map[string]interface{}{
		{
			"fs_group":               sc.FSGroup,
			"fs_group_change_policy": sc.FSGroupChangePolicy,
			"run_as_group":           sc.RunAsGroup,
			"run_as_non_root":        sc.RunAsNonRoot,
			"run_as_user":            sc.RunAsUser,
			"supplemental_groups":    sc.SupplementalGroups,
			"sysctls":                sysctls,
		},
	}
}

func flattenVolumes(volumes []pod.Volume) []map[string]interface{} {
	if len(volumes) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(volumes))
	for i, v := range volumes {
		m := map[string]interface{}{
			"name": v.Name,
		}
		if v.ConfigMap != nil {
			m["config_map"] = []map[string]interface{}{
				{
					"name":         v.ConfigMap.Name,
					"default_mode": v.ConfigMap.DefaultMode,
					"items":        flattenKeyToPath(v.ConfigMap.Items),
					"optional":     v.ConfigMap.Optional,
				},
			}
		}
		if v.NFS != nil {
			m["nfs"] = []map[string]interface{}{
				{
					"path":      v.NFS.Path,
					"server":    v.NFS.Server,
					"read_only": v.NFS.ReadOnly,
				},
			}
		}
		if v.PersistentVolumeClaim != nil {
			m["persistent_volume_claim"] = []map[string]interface{}{
				{
					"claim_name": v.PersistentVolumeClaim.ClaimName,
					"read_only":  v.PersistentVolumeClaim.ReadOnly,
				},
			}
		}
		if v.Secret != nil {
			m["secret"] = []map[string]interface{}{
				{
					"secret_name":  v.Secret.SecretName,
					"default_mode": v.Secret.DefaultMode,
					"items":        flattenKeyToPath(v.Secret.Items),
					"optional":     v.Secret.Optional,
				},
			}
		}
		if v.Projected != nil {
			m["projected"] = []map[string]interface{}{
				{
					"default_mode": v.Projected.DefaultMode,
					"sources":      flattenVolumeProjections(v.Projected.Sources),
				},
			}
		}
		result[i] = m
	}
	return result
}

func flattenVolumeProjections(projections []pod.VolumeProjection) []map[string]interface{} {
	if len(projections) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(projections))
	for i, p := range projections {
		m := map[string]interface{}{}
		if p.ConfigMap != nil {
			m["config_map"] = []map[string]interface{}{
				{
					"name":     p.ConfigMap.Name,
					"items":    flattenKeyToPath(p.ConfigMap.Items),
					"optional": p.ConfigMap.Optional,
				},
			}
		}
		if p.DownwardAPI != nil {
			m["downward_api"] = []map[string]interface{}{
				{
					"items": flattenDownwardAPIVolumeFiles(p.DownwardAPI.Items),
				},
			}
		}
		if p.Secret != nil {
			m["secret"] = []map[string]interface{}{
				{
					"name":     p.Secret.Name,
					"items":    flattenKeyToPath(p.Secret.Items),
					"optional": p.Secret.Optional,
				},
			}
		}
		result[i] = m
	}
	return result
}

func flattenDownwardAPIVolumeFiles(files []pod.DownwardAPIVolumeFile) []map[string]interface{} {
	if len(files) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(files))
	for i, f := range files {
		m := map[string]interface{}{
			"path": f.Path,
			"mode": f.Mode,
		}
		if f.FieldRef != nil {
			m["field_ref"] = []map[string]interface{}{
				{
					"api_version": f.FieldRef.APIVersion,
					"field_path":  f.FieldRef.FieldPath,
				},
			}
		}
		if f.ResourceFieldRef != nil {
			m["resource_field_ref"] = []map[string]interface{}{
				{
					"container_name": f.ResourceFieldRef.ContainerName,
					"resource":       f.ResourceFieldRef.Resource,
				},
			}
		}
		result[i] = m
	}
	return result
}

func flattenKeyToPath(items []pod.KeyToPath) []map[string]interface{} {
	if len(items) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"key":  item.Key,
			"path": item.Path,
			"mode": item.Mode,
		}
	}
	return result
}

func flattenPodStatusAttr(status pod.PodStatus) []map[string]interface{} {
	conditions := make([]map[string]interface{}, len(status.Conditions))
	for i, c := range status.Conditions {
		conditions[i] = map[string]interface{}{
			"type":                 c.Type,
			"status":               c.Status,
			"last_probe_time":      c.LastProbeTime,
			"last_transition_time": c.LastTransitionTime,
			"reason":               c.Reason,
			"message":              c.Message,
		}
	}
	return []map[string]interface{}{
		{
			"phase":      status.Phase,
			"conditions": conditions,
		},
	}
}

func derefInt64(p *int64) int {
	if p == nil {
		return 0
	}
	return int(*p)
}

func derefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}
