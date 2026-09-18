package er

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/er/v3/instance"
)

func TestEnterpriseProjectSchemas(t *testing.T) {
	resourceField := ResourceErInstanceV3().Schema["enterprise_project_id"]
	if !resourceField.Optional || !resourceField.Computed || !resourceField.ForceNew {
		t.Fatal("ER instance enterprise_project_id must be optional, computed, and force replacement")
	}

	instances := DataSourceErInstancesV3().Schema
	if !instances["enterprise_project_id"].Optional {
		t.Fatal("ER instances enterprise_project_id must be an optional filter")
	}
	nested := instances["instances"].Elem.(*schema.Resource).Schema["enterprise_project_id"]
	if !nested.Computed {
		t.Fatal("ER instances must export enterprise_project_id")
	}

	if !DataSourceErFlowLogsV3().Schema["enterprise_project_id"].Optional {
		t.Fatal("ER flow logs enterprise_project_id must be an optional filter")
	}
}

func TestBuildInstanceListOptsEnterpriseProject(t *testing.T) {
	d := schema.TestResourceDataRaw(t, DataSourceErInstancesV3().Schema, nil)
	opts := buildInstanceListOpts(d, "enterprise-project-id")
	if len(opts.EnterpriseProjectId) != 1 || opts.EnterpriseProjectId[0] != "enterprise-project-id" {
		t.Fatalf("unexpected enterprise project filter: %#v", opts.EnterpriseProjectId)
	}
}

func TestFlattenInstancesEnterpriseProject(t *testing.T) {
	result := flattenInstances([]instance.RouterInstance{{
		ID:                  "instance-id",
		EnterpriseProjectId: "enterprise-project-id",
	}})
	if result[0]["enterprise_project_id"] != "enterprise-project-id" {
		t.Fatalf("unexpected enterprise project value: %#v", result[0]["enterprise_project_id"])
	}
}
