---
page_title: "oneuptime_workspace_notification_summary Resource - oneuptime"
subcategory: "Other"
description: |-
  Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams
---

# oneuptime_workspace_notification_summary (Resource)

Recurring summary reports for incidents and alerts sent to Slack or Microsoft Teams

## Example Usage

```terraform
resource "oneuptime_workspace_notification_summary" "example" {
  name = "This is an example of longer text content that might be stored in this field."
  workspace_type = "This is an example of longer text content that might be stored in this field."
  summary_type = "Example short text"
  number_of_days_of_data = 42
  is_enabled = true
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of the Summary Rule..
- `workspace_type` (String) Type of Workspace - Slack, Microsoft Teams, etc...
- `summary_type` (String) Type of summary - Incident, Alert, Incident Episode, or Alert Episode..
- `number_of_days_of_data` (Number) How many days of data to include in the summary..
- `is_enabled` (Bool) Is this summary rule enabled?..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of the Summary Rule..
- `recurring_interval` (String) How often should the summary be sent?..
- `send_first_report_at` (String) A date time object..
- `timezone` (String) The IANA time zone the summary's schedule is read in, such as Europe/Berlin or America/New_York. The summary goes out at the same time of day there all year, also after the clocks change for daylight saving time. Left out when the summary is created, it is the time zone in the creator's profile, or UTC when no person creates it (an API key or a workflow). A summary without one is read in UTC...
- `channel_names` (String) List of channel names to post the summary to..
- `team_name` (String) Microsoft Teams team name (only for Microsoft Teams)..
- `summary_items` (String) Checklist of items to include in the summary..
- `filters` (String) Filter conditions for which items to include in the summary..
- `filter_condition` (String) How to combine filters - Any or All..
- `next_send_at` (String) A date time object..
- `last_sent_at` (String) A date time object..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..
- `deleted_by_user_id` (String) A unique identifier for an object, represented as a UUID..

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_workspace_notification_summary.example <id>
```
