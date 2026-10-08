---
page_title: "oneuptime_incident_video_call Resource - oneuptime"
subcategory: "Incidents"
description: |-
  Video calls for an incident: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.
---

# oneuptime_incident_video_call (Resource)

Video calls for an incident: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.

## Example Usage

```terraform
resource "oneuptime_incident_video_call" "example" {
  incident_id = oneuptime_incident.example.id
}
```

## Schema

### Required

- `incident_id` (String) ID of the incident this call is for. The ID of a `oneuptime_incident`.

### Optional

- `join_url` (String) The link responders open to join the call. Set by OneUptime for a call started with a connection or a Slack huddle; required, as an https link, when you add a link of your own.
- `provider_value` (String) Where the call is held: Zoom, GoogleMeet, MicrosoftTeams, SlackHuddle, or CustomLink for a link a person provided. Taken from the connection when one is given.
- `title` (String) What the call is called. A meeting a provider creates is named for the incident.
- `video_call_connection_id` (String) ID of the Zoom, Google Meet, Microsoft Teams or meeting link connection to start the call with. Leave it out to add a link of your own, or to start the huddle of the incident's Slack channel. The ID of a `oneuptime_video_call_connection`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `external_meeting_id` (String) The provider's own id for the meeting: a Zoom meeting id, a Google Meet space name, a Microsoft Teams online meeting id or a Slack channel id.
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
- `workspace_notification_rule_id` (String) ID of the Slack or Microsoft Teams notification rule that started this call. The ID of a `oneuptime_workspace_notification_rule`.

## Import

Import an existing incident video call by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_incident_video_call.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_incident_video_call.example <id>
```
