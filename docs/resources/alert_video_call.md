---
page_title: "oneuptime_alert_video_call Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Video calls for an alert: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.
---

# oneuptime_alert_video_call (Resource)

Video calls for an alert: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.

## Example Usage

```terraform
resource "oneuptime_alert_video_call" "example" {
  alert_id = oneuptime_alert.example.id
}
```

## Schema

### Required

- `alert_id` (String) ID of the alert this call is for. The ID of a `oneuptime_alert`.

### Optional

- `join_url` (String) The link responders open to join the call. Set by OneUptime for a call started with a connection or a Slack huddle; required, as an https link, when you add a link of your own.
- `provider_value` (String) Where the call is held: Zoom, GoogleMeet, MicrosoftTeams, SlackHuddle, or CustomLink for a link a person provided. Taken from the connection when one is given. Allowed values: `Zoom`, `GoogleMeet`, `MicrosoftTeams`, `SlackHuddle`, `CustomLink`.
- `title` (String) What the call is called. A meeting a provider creates is named for the alert.
- `video_call_connection_id` (String) ID of the Zoom, Google Meet, Microsoft Teams or meeting link connection to start the call with. Leave it out to add a link of your own, or to start the huddle of the alert's Slack channel. The ID of a `oneuptime_video_call_connection`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `external_meeting_id` (String) The provider's own id for the meeting: a Zoom meeting id, a Google Meet space name, a Microsoft Teams online meeting id or a Slack channel id.
- `id` (String) Unique identifier for the resource.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `workspace_notification_rule_id` (String) ID of the Slack or Microsoft Teams notification rule that started this call. The ID of a `oneuptime_workspace_notification_rule`.

## Import

Import an existing alert video call by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_video_call.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_video_call.example <id>
```
