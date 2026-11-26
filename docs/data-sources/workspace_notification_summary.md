---
page_title: "oneuptime_workspace_notification_summary Data Source - oneuptime"
subcategory: "Slack, Teams & Video Calls"
description: |-
  Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams
---

# oneuptime_workspace_notification_summary (Data Source)

Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams

Look one up by `id`, or by any of its other arguments: every argument you set must match, and exactly one workspace notification summary may match them all - none, or more than one, is an error rather than an empty or arbitrary result.

## Example Usage

```terraform
data "oneuptime_workspace_notification_summary" "example" {
  name = "Example workspace notification summary"
}

# Or by id:
data "oneuptime_workspace_notification_summary" "by_id" {
  id = "<id>"
}
```

## Schema

### Optional

- `created_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that created this object. Empty when a person or OneUptime itself created it.
- `created_by_actor_type` (String) What created this object: User, ApiKey, Workflow, McpClient (for a person), MasterAdmin or System (OneUptime). Empty if made before it was recorded.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of the Summary Rule.
- `filter_condition` (String) How to combine filters - Any or All.
- `id` (String) Look up by unique identifier. Leave unset to look up by the other arguments instead.
- `is_enabled` (Boolean) Is this summary rule enabled?
- `last_modified_by_actor_id` (String) ID of the API key, workflow or MCP client authorization that last changed this object. Empty when a person last changed it.
- `last_modified_by_actor_type` (String) What last changed this object: User, ApiKey, Workflow, McpClient (for a person) or MasterAdmin. Empty until it is changed after it is created.
- `last_modified_by_user_id` (String) ID of the person who last changed this object, or the member an MCP client changed it for. Empty when an API key or a workflow last changed it. The ID of a `oneuptime_user` (see the data source).
- `name` (String) Name of the Summary Rule.
- `number_of_days_of_data` (Number) How many days of data to include in the summary.
- `summary_type` (String) Type of summary - Incident, Alert, Incident Episode, or Alert Episode.
- `team_name` (String) The team the summary's channels are in: for Microsoft Teams, the team's id as picked from the connection's teams (a team's name, as summaries saved before held it, is matched too) - empty for the connection's only team; for Mattermost, the team's id - empty to look in the default team first, then in the bot's other teams. Not used for Slack.
- `timezone` (String) The IANA time zone the summary's schedule is read in, such as Europe/Berlin or America/New_York. The summary goes out at the same time of day there all year, also after the clocks change for daylight saving time. Left out when the summary is created, it is the time zone in the creator's profile, or UTC when no person creates it (an API key or a workflow). A summary without one is read in UTC.
- `workspace_type` (String) Type of Workspace - Slack, Microsoft Teams, etc.

### Read-Only

- `channel_names` (String) List of channel names to post the summary to. A JSON value: read it with `jsondecode()`.
- `created_at` (String) Date and Time when the object was created.
- `filters` (String) Which items the summary includes. An owner condition, such as Incident Owner, names people and teams of this project by ID. A JSON value: read it with `jsondecode()`.
- `last_modified_at` (String) When a person, an API key, a workflow or an MCP client last changed this object. OneUptime's own updates do not count. Empty until it is changed.
- `last_sent_at` (String) When the last summary was sent.
- `next_send_at` (String) When the next summary should be sent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `recurring_interval` (String) How often should the summary be sent? A JSON value: read it with `jsondecode()`.
- `send_first_report_at` (String) When should the first summary report be sent? Subsequent reports will follow the recurring interval from this date.
- `summary_items` (String) Checklist of items to include in the summary. A JSON value: read it with `jsondecode()`.
- `updated_at` (String) Date and Time when the object was updated.
