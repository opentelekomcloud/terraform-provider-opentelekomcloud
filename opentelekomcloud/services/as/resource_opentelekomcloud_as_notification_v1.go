package as

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/autoscaling/v1/notifications"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

var validateASNotificationIdentifier = validation.StringMatch(
	regexp.MustCompile(`^[^/\s?#]+$`),
	"must be nonempty and must not contain whitespace, '/', '?' or '#'",
)

func ResourceASNotificationV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASNotificationV1Create,
		ReadContext:   resourceASNotificationV1Read,
		DeleteContext: resourceASNotificationV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceASNotificationV1Import,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scaling_group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateASNotificationIdentifier,
			},
			"topic_urn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateASNotificationIdentifier,
			},
			"topic_scene": {
				Type:     schema.TypeSet,
				Required: true,
				// The API documents adding scenes, not replacing an existing scene set.
				ForceNew: true,
				MinItems: 1,
				MaxItems: 5,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"SCALING_UP", "SCALING_UP_FAIL", "SCALING_DOWN",
						"SCALING_DOWN_FAIL", "SCALING_GROUP_ABNORMAL",
					}, false),
				},
			},
			"topic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceASNotificationV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}

	groupID := d.Get("scaling_group_id").(string)
	topicURN := d.Get("topic_urn").(string)
	_, err = notifications.Create(client, notifications.CreateOpts{
		TopicURN:   topicURN,
		TopicScene: common.ExpandToStringSlice(d.Get("topic_scene").(*schema.Set).List()),
	}, groupID)
	if err != nil {
		return fmterr.Errorf("error creating AS notification for group %s: %w", groupID, err)
	}
	d.SetId(groupID + "/" + topicURN)
	expectedScenes := schema.NewSet(schema.HashString, d.Get("topic_scene").(*schema.Set).List())

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"pending"},
		Target:       []string{"present"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 2 * time.Second,
		Refresh: func() (interface{}, string, error) {
			notification, err := findASNotificationV1(client, groupID, topicURN)
			if err != nil {
				return nil, "", err
			}
			if notification == nil {
				return nil, "pending", nil
			}
			actualScenes := schema.NewSet(schema.HashString, nil)
			for _, scene := range notification.TopicScene {
				actualScenes.Add(scene)
			}
			if !expectedScenes.Equal(actualScenes) {
				return notification, "pending", nil
			}
			return notification, "present", nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for AS notification %s to become available: %w", d.Id(), err)
	}
	return resourceASNotificationV1Read(common.CtxWithClient(ctx, client, keyClientV1), d, meta)
}

func resourceASNotificationV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	notification, err := findASNotificationV1(client, d.Get("scaling_group_id").(string), d.Get("topic_urn").(string))
	if err != nil {
		if asNotificationV1NotFound(err) {
			log.Printf("[WARN] AS notification %s no longer exists, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error reading AS notification %s", d.Id()))
	}
	if notification == nil {
		log.Printf("[WARN] AS notification %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("topic_urn", notification.TopicURN),
		d.Set("topic_name", notification.TopicName),
		d.Set("topic_scene", notification.TopicScene),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceASNotificationV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV1, func() (*golangsdk.ServiceClient, error) {
		return config.AutoscalingV1Client(config.GetRegion(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationV1Client, err)
	}
	groupID := d.Get("scaling_group_id").(string)
	topicURN := d.Get("topic_urn").(string)
	if err := notifications.Delete(client, groupID, topicURN); err != nil {
		if asNotificationV1NotFound(err) {
			log.Printf("[WARN] AS notification %s already deleted", d.Id())
			d.SetId("")
			return nil
		}
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting AS notification %s", d.Id()))
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"present"},
		Target:       []string{"deleted"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		PollInterval: 2 * time.Second,
		Refresh: func() (interface{}, string, error) {
			notification, err := findASNotificationV1(client, groupID, topicURN)
			if asNotificationV1NotFound(err) || (err == nil && notification == nil) {
				return true, "deleted", nil
			}
			if err != nil {
				return nil, "", err
			}
			return notification, "present", nil
		},
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmterr.Errorf("error waiting for AS notification %s to be deleted: %w", d.Id(), err)
	}
	d.SetId("")
	return nil
}

func asNotificationV1NotFound(err error) bool {
	var notFound golangsdk.ErrDefault404
	if errors.As(err, &notFound) {
		return true
	}
	var badRequest golangsdk.ErrDefault400
	if !errors.As(err, &badRequest) {
		return false
	}
	var response struct {
		ErrorCode string `json:"error_code"`
		Error     struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if json.Unmarshal(badRequest.Body, &response) != nil {
		return false
	}
	// AS reports a deleted parent group as HTTP 400 with AS.2007.
	return response.Error.Code == "AS.2007" || response.ErrorCode == "AS.2007"
}

func findASNotificationV1(client *golangsdk.ServiceClient, groupID, topicURN string) (*notifications.Notification, error) {
	items, err := notifications.List(client, groupID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.TopicURN == topicURN {
			return &item, nil
		}
	}
	return nil, nil
}

func resourceASNotificationV1Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("AS notification import ID must have format <scaling_group_id>/<topic_urn>")
	}
	for i, name := range []string{"scaling_group_id", "topic_urn"} {
		if _, errs := validateASNotificationIdentifier(parts[i], name); len(errs) != 0 {
			return nil, fmt.Errorf("invalid AS notification import ID: %w", errs[0])
		}
	}
	return common.ImportByPath("scaling_group_id", "topic_urn")(ctx, d, meta)
}
