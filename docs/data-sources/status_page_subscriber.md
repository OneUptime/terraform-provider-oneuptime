---
page_title: "oneuptime_status_page_subscriber Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Subscriber that subscribed to your status page
---

# oneuptime_status_page_subscriber (Data Source)

Subscriber that subscribed to your status page Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_subscriber" "by_name" {
  name = "example-status_page_subscriber"
}

data "oneuptime_status_page_subscriber" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `subscriber_email` (String) Email object. Computed.
- `subscriber_phone` (String) Phone object. Computed.
- `subscriber_webhook` (String) Webhook to ping when events happen on Status Page.. Computed.
- `slack_workspace_name` (String) Name of the Slack workspace for validation and identification.. Computed.
- `microsoft_teams_workspace_name` (String) Name of the Microsoft Teams workspace for validation and identification.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `is_subscription_confirmed` (Bool) Has subscriber confirmed their subscription? (for example, by clicking on a confirmation link in an email).. Computed.
- `is_unsubscribed` (Bool) Is Subscriber Unsubscribed?.. Computed.
- `send_you_have_subscribed_message` (Bool) Send You Have Subscribed Message when subscriber is created?.. Computed.
- `is_subscribed_to_all_resources` (Bool) Is Subscriber Subscribed to All Resources on this status page?.. Computed.
- `is_subscribed_to_all_event_types` (Bool) Is Subscriber Subscribed to All Event Types (like Incidents, Scheduled Events, Announcements) on this status page?.. Computed.
- `status_page_resources` (Set) Relation to Status Page Resources where this subscriber is subscribed to.. Computed.
- `status_page_event_types` (String) Which event types is the subscriber subscribed to (like Incidents, Scheduled Events, Announcements).. Computed.
- `internal_note` (String) Any notes or text you would like to add to this subscriber object. This is for internal use only... Computed.
