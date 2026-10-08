---
page_title: "oneuptime_alert_video_call Data Source - oneuptime"
subcategory: "Alerts"
description: |-
  Video calls for an alert: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.
---

# oneuptime_alert_video_call (Data Source)

Video calls for an alert: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one alert video call may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_alert_video_call" "example" {
  alert_id = oneuptime_alert.example.id
}

# Or by id:
data "oneuptime_alert_video_call" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `alert_id` (String) ID of the alert this call is for. The ID of a `oneuptime_alert`.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `external_meeting_id` (String) The provider's own id for the meeting: a Zoom meeting id, a Google Meet space name, a Microsoft Teams online meeting id or a Slack channel id.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `join_url` (String) The link responders open to join the call. Set by OneUptime for a call started with a connection or a Slack huddle; required, as an https link, when you add a link of your own.
- `provider_value` (String) Where the call is held: Zoom, GoogleMeet, MicrosoftTeams, SlackHuddle, or CustomLink for a link a person provided. Taken from the connection when one is given.
- `title` (String) What the call is called. A meeting a provider creates is named for the alert.
- `video_call_connection_id` (String) ID of the Zoom, Google Meet, Microsoft Teams or meeting link connection to start the call with. Leave it out to add a link of your own, or to start the huddle of the alert's Slack channel. The ID of a `oneuptime_video_call_connection`.
- `workspace_notification_rule_id` (String) ID of the Slack or Microsoft Teams notification rule that started this call. The ID of a `oneuptime_workspace_notification_rule`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.
