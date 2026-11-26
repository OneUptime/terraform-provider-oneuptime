---
page_title: "oneuptime_status_page_subscriber Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Subscriber that subscribed to your status page
---

# oneuptime_status_page_subscriber (Data Source)

Subscriber that subscribed to your status page

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one status page subscriber may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_status_page_subscriber" "example" {
  created_by_actor_type = "User"
}

# Or by id:
data "oneuptime_status_page_subscriber" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `internal_note` (String) Any notes or text you would like to add to this subscriber object. This is for internal use only.
- `is_added_by_team` (Boolean) Whether your team added this subscriber (from the dashboard, with an API key or by a workflow) rather than the subscriber signing up on the status page. Set by OneUptime when the subscriber is created; any value sent for it is ignored.
- `is_subscribed_to_all_event_types` (Boolean) Is Subscriber Subscribed to All Event Types (like Incidents, Maintenance Events, Announcements) on this status page?
- `is_subscribed_to_all_resources` (Boolean) Is Subscriber Subscribed to All Resources on this status page?
- `is_subscription_confirmed` (Boolean) Has subscriber confirmed their subscription? (for example, by clicking on a confirmation link in an email).
- `is_unsubscribed` (Boolean) Is Subscriber Unsubscribed?
- `language` (String) The language this subscriber is notified in. Empty means the status page's default language.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `microsoft_teams_workspace_name` (String) Name of the Microsoft Teams workspace for validation and identification.
- `send_you_have_subscribed_message` (Boolean) Send You Have Subscribed Message when subscriber is created?
- `slack_workspace_name` (String) Name of the Slack workspace for validation and identification.
- `status_page_id` (String) ID of your Status Page resource where this object belongs. The ID of a `oneuptime_status_page`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `status_page_event_types` (String) Which event types is the subscriber subscribed to (like Incidents, Maintenance Events, Announcements). A JSON value: read it with `jsondecode()`.
- `status_page_resources` (Set of String) Relation to Status Page Resources where this subscriber is subscribed to. IDs of `oneuptime_status_page_resource` resources.
- `subscriber_email` (String) Email address of the subscriber.
- `subscriber_phone` (String) Phone number of subscriber.
- `subscriber_webhook` (String, Sensitive) Webhook to ping when events happen on Status Page. Kept encrypted at rest.
- `unsubscribed_at` (String) When this subscriber unsubscribed. Set by OneUptime when Is Unsubscribed is turned on, and cleared when it is turned off; any value sent for it is ignored.
- `updated_at` (String) Date and Time when the object was updated.
