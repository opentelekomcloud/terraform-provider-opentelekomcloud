package smn

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/opentelekomcloud/gophertelekomcloud/openstack/smn/v2/subscriptions"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
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
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					return common.ValidateStringList(v, k, []string{"email", "sms", "http", "https"})
				},
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
		return diag.Errorf("Error creating OpenTelekomCloud smn client: %s", err)
	}
	topicUrn := d.Get("topic_urn").(string)
	createOpts := subscriptions.CreateOps{
		Endpoint: d.Get("endpoint").(string),
		Protocol: d.Get("protocol").(string),
		Remark:   d.Get("remark").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	subscription, err := subscriptions.Create(client, createOpts, topicUrn).Extract()
	if err != nil {
		return diag.Errorf("Error getting subscription from result: %s", err)
	}
	log.Printf("[DEBUG] Create : subscription.SubscriptionUrn %s", subscription.SubscriptionUrn)
	if subscription.SubscriptionUrn != "" {
		d.SetId(subscription.SubscriptionUrn)
		d.Set("subscription_urn", subscription.SubscriptionUrn)
		return resourceSubscriptionRead(ctx, d, meta)
	}

	return diag.Errorf("Unexpected conversion error in resourceSubscriptionCreate.")
}

func resourceSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return diag.Errorf("Error creating OpenTelekomCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Deleting subscription %s", d.Id())

	id := d.Id()
	result := subscriptions.Delete(client, id)
	if result.Err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully deleted subscription %s", id)
	return nil
}

func resourceSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*cfg.Config)
	client, err := config.SmnV2Client(config.GetProjectName(d))
	if err != nil {
		return diag.Errorf("Error creating OpenTelekomCloud smn client: %s", err)
	}

	log.Printf("[DEBUG] Getting subscription %s", d.Id())

	id := d.Id()
	subscriptionsList, err := subscriptions.List(client).Extract()
	if err != nil {
		return diag.Errorf("Error Get subscriptionsList: %s", err)
	}
	log.Printf("[DEBUG] list : subscriptionsList %#v", subscriptionsList)
	for _, subscription := range subscriptionsList {
		if subscription.SubscriptionUrn == id {
			log.Printf("[DEBUG] subscription: %#v", subscription)
			d.Set("topic_urn", subscription.TopicUrn)
			d.Set("endpoint", subscription.Endpoint)
			d.Set("protocol", subscription.Protocol)
			d.Set("subscription_urn", subscription.SubscriptionUrn)
			d.Set("owner", subscription.Owner)
			d.Set("remark", subscription.Remark)
			d.Set("status", subscription.Status)
		}
	}

	log.Printf("[DEBUG] Successfully get subscription %s", id)
	return nil
}
