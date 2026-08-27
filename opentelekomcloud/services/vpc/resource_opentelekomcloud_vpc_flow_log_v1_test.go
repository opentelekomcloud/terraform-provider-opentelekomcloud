package vpc

import "testing"

func TestResourceVpcFlowLogV1Schema(t *testing.T) {
	resource := ResourceVpcFlowLogV1()

	if !resource.Schema["log_topic_id"].ForceNew {
		t.Fatal("log_topic_id must force replacement")
	}

	indexEnabled := resource.Schema["index_enabled"]
	if !indexEnabled.Optional || indexEnabled.Computed || !indexEnabled.ForceNew || indexEnabled.Default != false {
		t.Fatal("index_enabled must default to false and force replacement")
	}

	enabled := resource.Schema["enabled"]
	if !enabled.Optional || enabled.Computed || enabled.ForceNew || enabled.Default != true {
		t.Fatal("enabled must default to true and be updatable")
	}

	for _, field := range []string{"status", "tenant_id", "created_at", "updated_at"} {
		if !resource.Schema[field].Computed || resource.Schema[field].Optional {
			t.Fatalf("%s must be computed only", field)
		}
	}
}
