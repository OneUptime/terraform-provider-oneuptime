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
  incident_id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

### Required

- `incident_id` (String) A unique identifier for an object, represented as a UUID..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `provider` (String) Where the call is held: Zoom, GoogleMeet, MicrosoftTeams, SlackHuddle, or CustomLink for a link a person provided. Taken from the connection when one is given...
- `video_call_connection_id` (String) A unique identifier for an object, represented as a UUID..
- `title` (String) What the call is called. A meeting a provider creates is named for the incident...
- `join_url` (String) The link responders open to join the call. Set by OneUptime for a call started with a connection or a Slack huddle; required, as an https link, when you add a link of your own...

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `external_meeting_id` (String) The provider's own id for the meeting: a Zoom meeting id, a Google Meet space name, a Microsoft Teams online meeting id or a Slack channel id...
- `workspace_notification_rule_id` (String) A unique identifier for an object, represented as a UUID..
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_incident_video_call.example <id>
```
