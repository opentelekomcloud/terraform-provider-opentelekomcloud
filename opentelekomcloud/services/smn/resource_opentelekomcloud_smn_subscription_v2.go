package smn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/subscriptions"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/fmterr"
)

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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
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

	subscription, err := subscriptions.Create(client, createOpts, topicUrn).Extract()
	if err != nil {
		return fmterr.Errorf("error creating subscription: %w", err)
	}

	d.SetId(subscription.SubscriptionUrn)

	return resourceSubscriptionRead(ctx, d, meta)
}

func resourceSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	log.Printf("[DEBUG] Getting subscription %s", d.Id())

	subscriptionsList, err := subscriptions.List(client).Extract()
	if err != nil {
		return fmterr.Errorf("error getting subscriptions: %w", err)
	}
	log.Printf("[DEBUG] list : subscriptionsList %#v", subscriptionsList)
	for _, subscription := range subscriptionsList {
		if subscription.SubscriptionUrn == d.Id() {
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
		}
	}

	log.Printf("[DEBUG] Successfully get subscription %s", d.Id())
	return nil
}

func resourceSubscriptionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return fmterr.Errorf(errCreationClient, err)
	}

	log.Printf("[DEBUG] Deleting subscription %s", d.Id())

	if err := subscriptions.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully deleted subscription %s", d.Id())
	return nil
}
