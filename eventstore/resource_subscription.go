package eventstore

import (
  "context"
  "strings"

  "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
  "github.com/madedotcom/eventstore-client-go/eventstore"
)

func resourceSubscription() *schema.Resource {
  return &schema.Resource{
    CreateContext: resourceSubscriptionCreate,
    ReadContext:   resourceSubscriptionRead,
    UpdateContext: resourceSubscriptionUpdate,
    DeleteContext: resourceSubscriptionDelete,
    Importer: &schema.ResourceImporter{
       StateContext: schema.ImportStatePassthroughContext,
    },
    Schema: map[string]*schema.Schema{
      "stream_name": &schema.Schema{
        Type:     schema.TypeString,
        ForceNew: true,
        Required: true,
      },
      "subscription_name": &schema.Schema{
        Type:     schema.TypeString,
        ForceNew: true,
        Required: true,
      },
      "min_checkpoint_count": &schema.Schema{
        Type:      schema.TypeInt,
        Optional:  true,
        Default:   10,
      },
      "start_from": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  0,
      },
      "resolve_link_tos": {
        Type:     schema.TypeBool,
        Optional: true,
        Default:  true,
      },
      "read_batch_size": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  20,
      },
      "named_consumer_strategy": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
        Default:  "RoundRobin",
      },
      "extra_statistics": &schema.Schema{
        Type:     schema.TypeBool,
        Optional: true,
        Default:  false,
      },
      "max_retry_count": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  10,
      },
      "live_buffer_size": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  500,
      },
      "message_timeout_milliseconds": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  10000,
      },
      "max_checkpoint_count": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  500,
      },
      "max_subscriber_count": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  10,
      },
      "checkpoint_after_milliseconds": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  1000,
      },
      "buffer_size": &schema.Schema{
        Type:     schema.TypeInt,
        Optional: true,
        Default:  500,
      },
    },
  }
}

func resourceSubscriptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  client := m.(*eventstore.Client)

  subscription, err := client.CreateSubscription(
    d.Get("stream_name").(string),
    d.Get("subscription_name").(string),
    d.Get("min_checkpoint_count").(int),
    d.Get("start_from").(int),
    d.Get("resolve_link_tos").(bool),
    d.Get("read_batch_size").(int),
    d.Get("named_consumer_strategy").(string),
    d.Get("extra_statistics").(bool),
    d.Get("max_retry_count").(int),
    d.Get("live_buffer_size").(int),
    d.Get("message_timeout_milliseconds").(int),
    d.Get("max_checkpoint_count").(int),
    d.Get("max_subscriber_count").(int),
    d.Get("checkpoint_after_milliseconds").(int),
    d.Get("buffer_size").(int))

  if err != nil {
    return diag.FromErr(err)
  }

  d.SetId(subscription.StreamName + "/" + subscription.SubscriptionName)

  return diags
}

func resourceSubscriptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  client := m.(*eventstore.Client)

  idSlice := strings.Split(d.Id(), "/")

  subscription, err := client.GetSubscription(
    idSlice[0],
    idSlice[1])

  if err != nil {
    return diag.FromErr(err)
  }

  d.Set("stream_name", idSlice[0])
  d.Set("subscription_name", idSlice[1])

  d.Set("min_checkpoint_count", subscription.MinCheckPointCount)
  d.Set("start_from", subscription.StartFrom)
  d.Set("resolve_link_tos", subscription.ResolveLinkTos)
  d.Set("read_batch_size", subscription.ReadBatchSize)
  d.Set("named_consumer_strategy", subscription.NamedConsumerStrategy)
  d.Set("extra_statistics", subscription.ExtraStatistics)
  d.Set("max_retry_count", subscription.MaxRetryCount)
  d.Set("live_buffer_size", subscription.LiveBufferSize)
  d.Set("message_timeout_milliseconds", subscription.MessageTimeoutMilliseconds)
  d.Set("max_checkpoint_count", subscription.MaxCheckPointCount)
  d.Set("max_subscriber_count", subscription.MaxSubscriberCount)
  d.Set("checkpoint_after_milliseconds", subscription.CheckPointAfterMilliseconds)
  d.Set("buffer_size", subscription.BufferSize)

  return diags
}

func resourceSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  client := m.(*eventstore.Client)

  _, err := client.UpdateSubscription(
    d.Get("stream_name").(string),
    d.Get("subscription_name").(string),
    d.Get("min_checkpoint_count").(int),
    d.Get("start_from").(int),
    d.Get("resolve_link_tos").(bool),
    d.Get("read_batch_size").(int),
    d.Get("named_consumer_strategy").(string),
    d.Get("extra_statistics").(bool),
    d.Get("max_retry_count").(int),
    d.Get("live_buffer_size").(int),
    d.Get("message_timeout_milliseconds").(int),
    d.Get("max_checkpoint_count").(int),
    d.Get("max_subscriber_count").(int),
    d.Get("checkpoint_after_milliseconds").(int),
    d.Get("buffer_size").(int))

  if err != nil {
    return diag.FromErr(err)
  }  

  return resourceSubscriptionRead(ctx, d, m)
}

func resourceSubscriptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
  // Warning or errors can be collected in a slice type
  var diags diag.Diagnostics

  client := m.(*eventstore.Client)

  success, err := client.DeleteSubscription(d.Get("stream_name").(string), d.Get("subscription_name").(string))

  if err != nil {
    return diag.FromErr(err)
  }

  if !success {
    return diag.Errorf("Subscription Delete Failed")
  }

  d.SetId("")

  return diags
}
