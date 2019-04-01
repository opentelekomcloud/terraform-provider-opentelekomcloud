package opentelekomcloud

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
