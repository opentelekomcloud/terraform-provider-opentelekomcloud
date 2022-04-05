package smn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/topicattributes"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

func ResourceSMNTopicAttributeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSMNTopicAttributeV2Create,
		ReadContext:   resourceSMNTopicAttributeV2Read,
		DeleteContext: resourceSMNTopicAttributeV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: common.ImportByPath("topic_urn", "attribute_name"),
		},

		Schema: map[string]*schema.Schema{
			"topic_attribute": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: common.ValidateJsonString,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v.(string))
					return json
				},
			},
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attribute_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"access_policy",
				}, false),
			},
		},
	}
}

func resourceSMNTopicAttributeV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	topicURN := d.Get("topic_urn").(string)
	attributeName := d.Get("attribute_name").(string)

	createOpts := topicattributes.UpdateOpts{
		Value: d.Get("topic_attribute").(string),
	}

	if err := topicattributes.Update(client, topicURN, attributeName, createOpts).ExtractErr(); err != nil {
		return fmterr.Errorf("error updating SMN topic attribute: %w", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", topicURN, attributeName))
	return resourceSMNTopicAttributeV2Read(ctx, d, meta)
}

func resourceSMNTopicAttributeV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if err := setTopicURNAndAttributeName(d); err != nil {
		return diag.FromErr(err)
	}
	topicURN := d.Get("topic_urn").(string)
	attributeName := d.Get("attribute_name").(string)

	listOpts := topicattributes.ListOpts{
		Name: attributeName,
	}

	topicAttributesMap, err := topicattributes.List(client, topicURN, listOpts).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "topic attributes")
	}
	normalizePolicy, err := structure.NormalizeJsonString(topicAttributesMap["access_policy"])
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("topic_attribute", normalizePolicy); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSMNTopicAttributeV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	if err := setTopicURNAndAttributeName(d); err != nil {
		return diag.FromErr(err)
	}
	topicURN := d.Get("topic_urn").(string)
	attributeName := d.Get("attribute_name").(string)

	if err := topicattributes.Delete(client, topicURN, attributeName).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func setTopicURNAndAttributeName(d *schema.ResourceData) error {
	parts := strings.SplitN(d.Id(), "/", 2)
	mErr := multierror.Append(
		d.Set("topic_urn", parts[0]),
		d.Set("attribute_name", parts[1]),
	)
	return mErr.ErrorOrNil()
}
