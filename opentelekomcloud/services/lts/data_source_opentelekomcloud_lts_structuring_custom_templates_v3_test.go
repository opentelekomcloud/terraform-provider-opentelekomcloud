package lts

import (
	"testing"

	cloud_structuring "github.com/opentelekomcloud/gophertelekomcloud/openstack/lts/v2/cloud-structuring"
)

func TestDataSourceLtsStructuringCustomTemplatesV3Schema(t *testing.T) {
	dataSource := DataSourceLtsStructuringCustomTemplatesV3()

	if !dataSource.Schema["id"].Optional {
		t.Fatal("id must be an optional filter")
	}
	if !dataSource.Schema["templates"].Computed {
		t.Fatal("templates must be computed")
	}
}

func TestFlattenLtsStructuringCustomTemplates(t *testing.T) {
	templates := []cloud_structuring.StructTemplateModel{{
		ID:        "template-id",
		ProjectId: "project-id",
		Name:      "template-name",
		Type:      "json",
		DemoLog:   `{"message":"value"}`,
		DemoLabel: "label",
		CreatedAt: 1_700_000_000_000,
		DemoFields: []cloud_structuring.DemoField{{
			Name:            "message",
			Content:         "value",
			Type:            "string",
			IsAnalysis:      false,
			Index:           0,
			Relation:        "root",
			UserDefinedName: "alias",
		}},
		TagFields: []cloud_structuring.TagFieldNew{{
			Name:       "host",
			Content:    "host-1",
			Type:       "string",
			IsAnalysis: true,
			Index:      0,
		}},
		Rule: &cloud_structuring.TemplateRule{
			Type:  "json",
			Param: `{"layers":1}`,
		},
	}}

	result, ids := flattenLtsStructuringCustomTemplates(templates)

	if len(result) != 1 || len(ids) != 1 || ids[0] != "template-id" {
		t.Fatalf("unexpected flattened result: %#v, %#v", result, ids)
	}
	if result[0]["template_name"] != "template-name" {
		t.Fatalf("unexpected template name: %#v", result[0]["template_name"])
	}
	demoFields := result[0]["demo_fields"].([]map[string]interface{})
	if len(demoFields) != 1 || demoFields[0]["user_defined_name"] != "alias" {
		t.Fatalf("unexpected demo fields: %#v", demoFields)
	}
	if demoFields[0]["is_analysis"] != false || demoFields[0]["index"] != 0 {
		t.Fatalf("meaningful zero values were not preserved: %#v", demoFields[0])
	}
	tagFields := result[0]["tag_fields"].([]map[string]interface{})
	if len(tagFields) != 1 || tagFields[0]["field_name"] != "host" {
		t.Fatalf("unexpected tag fields: %#v", tagFields)
	}
	rule := result[0]["rule"].([]map[string]interface{})
	if len(rule) != 1 || rule[0]["type"] != "json" {
		t.Fatalf("unexpected rule: %#v", rule)
	}
}

func TestFlattenLtsStructuringCustomTemplatesSortsByID(t *testing.T) {
	templates := []cloud_structuring.StructTemplateModel{
		{ID: "template-b"},
		{ID: "template-a"},
	}

	result, ids := flattenLtsStructuringCustomTemplates(templates)

	if ids[0] != "template-a" || ids[1] != "template-b" {
		t.Fatalf("template IDs are not sorted: %#v", ids)
	}
	if result[0]["id"] != "template-a" || result[1]["id"] != "template-b" {
		t.Fatalf("flattened templates are not sorted: %#v", result)
	}
	if templates[0].ID != "template-b" {
		t.Fatalf("input templates were mutated: %#v", templates)
	}
}

func TestFlattenLtsStructuringCustomTemplateRuleNil(t *testing.T) {
	if result := flattenLtsStructuringCustomTemplateRule(nil); result != nil {
		t.Fatalf("expected nil rule, got %#v", result)
	}
}
