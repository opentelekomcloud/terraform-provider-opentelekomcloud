package opentelekomcloud

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
)

func checkCssClusterV1ExtendClusterFinished(data interface{}) bool {
	instances, err := navigateValue(data, []string{"instances"}, nil)
	if err != nil {
		return false
	}
	if v, ok := instances.([]interface{}); ok {
		if len(v) == 0 {
			return false
		}
		for _, item := range v {
			status, err := navigateValue(item, []string{"status"}, nil)
			if err != nil {
				return false
			}
			if s, ok := status.(string); !ok || "200" != s {
				return false
			}
		}
		return true
	}
	return false
}

func expandCssClusterV1ExtendClusterNodeNum(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	t, _ := navigateValue(d, []string{"terraform_resource_data"}, nil)
	rd := t.(*schema.ResourceData)

	oldv, newv := rd.GetChange("expect_node_num")
	v := newv.(int) - oldv.(int)
	if v < 0 {
		return 0, fmt.Errorf("it only supports extending nodes")
	}
	return v, nil
}

func expandRdsInstanceV3CreateRegion(d interface{}, arrayIndex map[string]int) (interface{}, error) {
	return navigateValue(d, []string{"region"}, arrayIndex)
}
