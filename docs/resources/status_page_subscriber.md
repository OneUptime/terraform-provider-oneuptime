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
  status_page_id = oneuptime_status_page.example.id
}
```

## Schema

### Required

- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Optional

- `internal_note` (String) Any notes or text you would like to add to this subscriber object. This is for internal use only.
- `is_subscribed_to_all_event_types` (Boolean) Is Subscriber Subscribed to All Event Types (like Incidents, Maintenance Events, Announcements) on this status page? Defaults to `true`.
- `is_subscribed_to_all_resources` (Boolean) Is Subscriber Subscribed to All Resources on this status page? Defaults to `true`.
- `is_subscription_confirmed` (Boolean) Has subscriber confirmed their subscription? (for example, by clicking on a confirmation link in an email). OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent. Defaults to `false`.
- `is_unsubscribed` (Boolean) Is Subscriber Unsubscribed? OneUptime keeps it up to date once the record exists: the value given here is the one the record starts with, and a later change of it is never sent. Defaults to `false`.
- `language` (String) The language this subscriber is notified in. Empty means the status page's default language.
- `microsoft_teams_incoming_webhook_url` (String, Sensitive) Microsoft Teams incoming webhook URL to send notifications to Teams channel. Write-only: kept encrypted and never returned by the API.
- `microsoft_teams_workspace_name` (String) Name of the Microsoft Teams workspace for validation and identification.
- `send_you_have_subscribed_message` (Boolean) Send You Have Subscribed Message when subscriber is created? Defaults to `true`.
- `slack_incoming_webhook_url` (String, Sensitive) Slack incoming webhook URL to send notifications to Slack channel. Write-only: kept encrypted and never returned by the API.
- `slack_workspace_name` (String) Name of the Slack workspace for validation and identification.
- `status_page_event_types` (String) Which event types is the subscriber subscribed to (like Incidents, Maintenance Events, Announcements). A JSON value: write it with `jsonencode()`.
- `status_page_resources` (Set of String) Relation to Status Page Resources where this subscriber is subscribed to. IDs of `oneuptime_status_page_resource` resources.
- `subscriber_email` (String) Email address of the subscriber.
- `subscriber_phone` (String) Phone number of subscriber.
- `subscriber_webhook` (String, Sensitive) Webhook to ping when events happen on Status Page. Kept encrypted at rest.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `is_added_by_team` (Boolean) Whether your team added this subscriber (from the dashboard, with an API key or by a workflow) rather than the subscriber signing up on the status page. Set by OneUptime when the subscriber is created; any value sent for it is ignored.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `unsubscribed_at` (String) When this subscriber unsubscribed. Set by OneUptime when Is Unsubscribed is turned on, and cleared when it is turned off; any value sent for it is ignored.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing status page subscriber by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_status_page_subscriber.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_status_page_subscriber.example <id>
```
