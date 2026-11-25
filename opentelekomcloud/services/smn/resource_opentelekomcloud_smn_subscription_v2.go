package smn

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	golangsdk "github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/subscriptions"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

const subscriptionPageSize = 100

func ResourceSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubscriptionCreate,
		ReadContext:   resourceSubscriptionRead,
		DeleteContext: resourceSubscriptionDelete,

		Schema: map[string]*schema.Schema{
			"topic_urn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"http", "https", "sms", "email",
				}, false),
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subscription_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_name": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "Parameter has been deprecated",
			},
		},
	}
}

func resourceSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	topicUrn := d.Get("topic_urn").(string)
	createOpts := subscriptions.CreateOpts{
		Endpoint: d.Get("endpoint").(string),
		Protocol: d.Get("protocol").(string),
		Remark:   d.Get("remark").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	subscription, err := subscriptions.Create(client, createOpts, topicUrn)
	if err != nil {
		return fmterr.Errorf("error creating subscription: %w", err)
	}

	d.SetId(subscription.SubscriptionUrn)

	clientCtx := common.CtxWithClient(ctx, client, keyClientV2)
	return resourceSubscriptionRead(clientCtx, d, meta)
}

func resourceSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	log.Printf("[DEBUG] Getting subscription %s", d.Id())

	subscription, err := findTopicSubscription(client, d.Get("topic_urn").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting subscriptions")
	}
	if subscription == nil {
		log.Printf("[WARN] SMN subscription %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] subscription: %#v", subscription)
	mErr := multierror.Append(
		d.Set("topic_urn", subscription.TopicUrn),
		d.Set("endpoint", subscription.Endpoint),
		d.Set("protocol", subscription.Protocol),
		d.Set("subscription_urn", subscription.SubscriptionUrn),
		d.Set("owner", subscription.Owner),
		d.Set("remark", subscription.Remark),
		d.Set("status", subscription.Status),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully get subscription %s", d.Id())
	return nil
}

func resourceSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := common.ClientFromCtx(ctx, keyClientV2, func() (*golangsdk.ServiceClient, error) {
		return config.SmnV2Client(config.GetProjectName(d))
	})
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	log.Printf("[DEBUG] Deleting subscription %s", d.Id())

	if err := subscriptions.Delete(client, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully deleted subscription %s", d.Id())
	return nil
}

// findTopicSubscription walks every page of the topic's subscriptions (the API returns at most
// subscriptionPageSize per request) and returns the one with the given URN, or nil.
func findTopicSubscription(client *golangsdk.ServiceClient, topicUrn, subscriptionUrn string) (*subscriptions.Subscription, error) {
	for offset := 0; ; offset += subscriptionPageSize {
		page, err := subscriptions.ListTopic(client, subscriptions.ListTopicOpts{
			TopicUrn: topicUrn,
			Offset:   strconv.Itoa(offset),
			Limit:    subscriptionPageSize,
		})
		if err != nil {
			return nil, err
		}
		for i := range page {
			if page[i].SubscriptionUrn == subscriptionUrn {
				return &page[i], nil
			}
		}
		if len(page) < subscriptionPageSize {
			return nil, nil
		}
	}
}
