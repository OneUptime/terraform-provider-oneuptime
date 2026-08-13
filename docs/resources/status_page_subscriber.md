---
page_title: "oneuptime_status_page_subscriber Resource - oneuptime"
subcategory: "Status Pages"
description: |-
  Subscriber that subscribed to your status page
---

# oneuptime_status_page_subscriber (Resource)

Subscriber that subscribed to your status page

## Example Usage

```terraform
resource "oneuptime_status_page_subscriber" "example" {
  status_page_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `status_page_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `subscriber_email` (String) Email object.
- `subscriber_phone` (String) Phone object.
- `subscriber_webhook` (String) Webhook to ping when events happen on Status Page..
- `slack_incoming_webhook_url` (String) Slack incoming webhook URL to send notifications to Slack channel..
- `slack_workspace_name` (String) Name of the Slack workspace for validation and identification..
- `microsoft_teams_incoming_webhook_url` (String) Microsoft Teams incoming webhook URL to send notifications to Teams channel..
- `microsoft_teams_workspace_name` (String) Name of the Microsoft Teams workspace for validation and identification..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `is_subscription_confirmed` (Bool) Has subscriber confirmed their subscription? (for example, by clicking on a confirmation link in an email)..
- `subscription_confirmation_token` (String) Token used to confirm subscription. This is a random token that is sent to the subscriber's email address to confirm their subscription...
- `is_unsubscribed` (Bool) Is Subscriber Unsubscribed?..
- `send_you_have_subscribed_message` (Bool) Send You Have Subscribed Message when subscriber is created?..
- `is_subscribed_to_all_resources` (Bool) Is Subscriber Subscribed to All Resources on this status page?..
- `is_subscribed_to_all_event_types` (Bool) Is Subscriber Subscribed to All Event Types (like Incidents, Scheduled Events, Announcements) on this status page?..
- `status_page_resources` (Set) Relation to Status Page Resources where this subscriber is subscribed to..
- `status_page_event_types` (String) Which event types is the subscriber subscribed to (like Incidents, Scheduled Events, Announcements)..
- `internal_note` (String) Any notes or text you would like to add to this subscriber object. This is for internal use only...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_status_page_subscriber.example <id>
```
