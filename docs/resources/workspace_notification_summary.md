---
page_title: "oneuptime_workspace_notification_summary Resource - oneuptime"
subcategory: "Slack, Teams & Video Calls"
description: |-
  Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams
---

# oneuptime_workspace_notification_summary (Resource)

Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams

## Example Usage

```terraform
resource "oneuptime_workspace_notification_summary" "example" {
  name                   = "Example workspace notification summary"
  workspace_type         = "Slack"
  summary_type           = "Incident"
  number_of_days_of_data = 42
  is_enabled             = true
  description            = "Managed by Terraform"
}
```

## Schema

### Required

- `is_enabled` (Boolean) Is this summary rule enabled?
- `name` (String) Name of the Summary Rule.
- `number_of_days_of_data` (Number) How many days of data to include in the summary.
- `summary_type` (String) Type of summary - Incident, Alert, Incident Episode, or Alert Episode. Allowed values: `Incident`, `Alert`, `Incident Episode`, `Alert Episode`.
- `workspace_type` (String) Type of Workspace - Slack, Microsoft Teams, etc. Allowed values: `Slack`, `MicrosoftTeams`, `Mattermost`.

### Optional

- `channel_names` (String) List of channel names to post the summary to. A JSON value: write it with `jsonencode()`.
- `description` (String) Description of the Summary Rule.
- `filter_condition` (String) How to combine filters - Any or All. Allowed values: `All`, `Any`.
- `filters` (String) Which items the summary includes. An owner condition, such as Incident Owner, names people and teams of this project by ID. A JSON value: write it with `jsonencode()`.
- `recurring_interval` (String) How often should the summary be sent? A JSON value: write it with `jsonencode()`.
- `send_first_report_at` (String) When should the first summary report be sent? Subsequent reports will follow the recurring interval from this date.
- `summary_items` (String) Checklist of items to include in the summary. A JSON value: write it with `jsonencode()`.
- `team_name` (String) The team the summary's channels are in: for Microsoft Teams, the team's id as picked from the connection's teams (a team's name, as summaries saved before held it, is matched too) - empty for the connection's only team; for Mattermost, the team's id - empty to look in the default team first, then in the bot's other teams. Not used for Slack.
- `timezone` (String) The IANA time zone the summary's schedule is read in, such as Europe/Berlin or America/New_York. The summary goes out at the same time of day there all year, also after the clocks change for daylight saving time. Left out when the summary is created, it is the time zone in the creator's profile, or UTC when no person creates it (an API key or a workflow). A summary without one is read in UTC.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `last_sent_at` (String) When the last summary was sent.
- `next_send_at` (String) When the next summary should be sent.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing workspace notification summary by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_workspace_notification_summary.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_workspace_notification_summary.example <id>
```
