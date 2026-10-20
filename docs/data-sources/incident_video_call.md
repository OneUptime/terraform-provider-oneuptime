---
page_title: "oneuptime_incident_video_call Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Video calls for an incident: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.
---

# oneuptime_incident_video_call (Data Source)

Video calls for an incident: a Zoom, Google Meet or Microsoft Teams meeting, a Slack huddle, or a meeting link.

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one incident video call may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_incident_video_call" "example" {
  incident_id = oneuptime_incident.example.id
}

# Or by id:
data "oneuptime_incident_video_call" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `call_summary_status` (String) Where OneUptime AI's reading of the transcript stands: Pending while it reads, Ready, Failed (Call Summary Status Message says why), or Skipped when AI is off for the project.
- `call_summary_status_message` (String) Why OneUptime AI did not read the transcript, or could not finish reading it.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `external_meeting_id` (String) The provider's own id for the meeting: a Zoom meeting id, a Google Meet space name, a Microsoft Teams online meeting id or a Slack channel id.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `incident_id` (String) ID of the incident this call is for. The ID of a `oneuptime_incident`.
- `join_url` (String) The link responders open to join the call. Set by OneUptime for a call started with a connection or a Slack huddle; required, as an https link, when you add a link of your own.
- `provider_value` (String) Where the call is held: Zoom, GoogleMeet, MicrosoftTeams, SlackHuddle, or CustomLink for a link a person provided. Taken from the connection when one is given.
- `root_cause_suggestion_reviewed_by_user_id` (String) ID of the person who accepted or dismissed the root cause suggested from the call. The ID of a `oneuptime_user` (see the data source).
- `root_cause_suggestion_status` (String) What a person did with the root cause OneUptime AI suggested from the call: Suggested while it waits, Accepted once used as the incident's root cause, Dismissed. Empty when nothing was suggested.
- `title` (String) What the call is called. A meeting a provider creates is named for the incident.
- `transcript` (String) What was said on the call, one remark a line: how far into the call (HH:MM:SS, empty when unknown), a tab, who spoke (empty when unknown), a tab, and the words. Read from Zoom or Google Meet, or added by a person, through the call's transcript routes.
- `transcript_added_by_user_id` (String) ID of the person who added the transcript. The ID of a `oneuptime_user` (see the data source).
- `transcript_fetch_reason` (String) Why the provider had no transcript to give when OneUptime last asked: NotReady (still preparing it), NotTranscribed (nothing recorded or transcribed yet), NoMeeting, Expired, PermissionMissing, ConnectionProblem or NotSupported. Empty when the provider did not answer.
- `transcript_fetch_status` (String) Whether OneUptime gets the call's transcript from Zoom or Google Meet by itself, as it does for the calls it starts while Record and transcribe incident calls is on: Waiting while it asks the provider (again at Transcript Fetch Next Attempt At), Stopped once it no longer asks (Transcript Fetch Reason says why). Empty for a call it does not ask about, and once it has the transcript.
- `transcript_source` (String) Where the transcript came from: Provider when OneUptime read it from Zoom or Google Meet, Upload when a person added it.
- `video_call_connection_id` (String) ID of the Zoom, Google Meet, Microsoft Teams or meeting link connection to start the call with. Leave it out to add a link of your own, or to start the huddle of the incident's Slack channel. The ID of a `oneuptime_video_call_connection`.
- `workspace_notification_rule_id` (String) ID of the Slack or Microsoft Teams notification rule that started this call. The ID of a `oneuptime_workspace_notification_rule`.

### Read-Only

- `call_summary` (String) What OneUptime AI read in the transcript: the root cause the call points at (a suggestion, or that the call did not discuss it enough), what was found, tried, decided and left to do, and open questions - each with the transcript lines it rests on. A JSON value: write it with `jsonencode()`.
- `call_summary_updated_at` (String) When the Call Summary Status last changed: when OneUptime AI started reading the transcript, or finished.
- `created_at` (String) Date and Time when the object was created.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `root_cause_suggestion_reviewed_at` (String) When a person accepted or dismissed the root cause suggested from the call.
- `transcript_added_at` (String) When the transcript was added to the call.
- `transcript_fetch_checked_at` (String) When OneUptime last asked the provider for the call's transcript.
- `transcript_fetch_next_attempt_at` (String) When OneUptime next asks the provider for the call's transcript, while the Transcript Fetch Status is Waiting.
- `transcript_info` (String) How long the transcript is - remarks, speakers, characters and how many seconds it spans - and the name of the file it was added from, so the call can be described without reading the transcript. A JSON value: write it with `jsonencode()`.
- `updated_at` (String) Date and Time when the object was updated.
