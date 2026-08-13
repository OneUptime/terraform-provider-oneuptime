---
page_title: "oneuptime_workspace_notification_summary Data Source - oneuptime"
subcategory: "Other"
description: |-
  Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams
---

# oneuptime_workspace_notification_summary (Data Source)

Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_workspace_notification_summary" "by_name" {
  name = "example-workspace_notification_summary"
}

data "oneuptime_workspace_notification_summary" "by_id" {
  id = "123e4567-e89b-12d3-a456-426614174000"
}
```

## Schema

- `id` (String) Look up by unique identifier. Exactly one of `id` or `name` must be set.. Computed.
- `name` (String) Look up by name. Exactly one of `id` or `name` must be set. Fails if the name does not match exactly one item.. Computed.
- `created_at` (String) A date time object.. Computed.
- `updated_at` (String) A date time object.. Computed.
- `deleted_at` (String) A date time object.. Computed.
- `version` (Number) Object version. Computed.
- `project_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of the Summary Rule.. Computed.
- `workspace_type` (String) Type of Workspace - Slack, Microsoft Teams, etc... Computed.
- `summary_type` (String) Type of summary - Incident, Alert, Incident Episode, or Alert Episode.. Computed.
- `recurring_interval` (String) How often should the summary be sent?.. Computed.
- `number_of_days_of_data` (Number) How many days of data to include in the summary.. Computed.
- `send_first_report_at` (String) A date time object.. Computed.
- `channel_names` (String) List of channel names to post the summary to.. Computed.
- `team_name` (String) Microsoft Teams team name (only for Microsoft Teams).. Computed.
- `summary_items` (String) Checklist of items to include in the summary.. Computed.
- `filters` (String) Filter conditions for which items to include in the summary.. Computed.
- `filter_condition` (String) How to combine filters - Any or All.. Computed.
- `next_send_at` (String) A date time object.. Computed.
- `last_sent_at` (String) A date time object.. Computed.
- `is_enabled` (Bool) Is this summary rule enabled?.. Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
