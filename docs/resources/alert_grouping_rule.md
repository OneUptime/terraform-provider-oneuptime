---
page_title: "oneuptime_alert_grouping_rule Resource - oneuptime"
subcategory: "Alerts"
description: |-
  Configure rules for automatically grouping related alerts into episodes
---

# oneuptime_alert_grouping_rule (Resource)

Configure rules for automatically grouping related alerts into episodes

## Example Usage

```terraform
resource "oneuptime_alert_grouping_rule" "example" {
  name        = "Example alert grouping rule"
  description = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Name of this alert grouping rule.

### Optional

- `alert_description_pattern` (String) Regular expression pattern to match alert descriptions. Leave empty to match any description.
- `alert_labels` (Set of String) Only group alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels. IDs of `oneuptime_label` resources.
- `alert_severities` (Set of String) Only group alerts with these severities. Leave empty to match alerts of any severity. IDs of `oneuptime_alert_severity` resources.
- `alert_title_pattern` (String) Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'.
- `criteria` (String) Versioned conditions that determine whether this rule matches a resource. A JSON value: write it with `jsonencode()`.
- `default_assign_to_team_id` (String) ID of defaultAssignToTeam. Kept for API compatibility: OneUptime does not show it anywhere. To make a team responsible for the episodes this rule opens, use episodeOwnerTeams. The ID of a `oneuptime_team`.
- `default_assign_to_user_id` (String) ID of defaultAssignToUser. Kept for API compatibility: OneUptime does not show it anywhere. To make someone responsible for the episodes this rule opens, use episodeOwnerUsers. The ID of a `oneuptime_user` (see the data source).
- `description` (String) Description of this alert grouping rule.
- `enable_inactivity_timeout` (Boolean) Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new alerts arrive. Defaults to `false`.
- `enable_reopen_window` (Boolean) Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution. Defaults to `false`.
- `enable_resolve_delay` (Boolean) Enable grace period before auto-resolving episode after all alerts resolve. Helps prevent rapid state changes during alert flapping. Defaults to `false`.
- `enable_time_window` (Boolean) Enable time-based grouping. When enabled, alerts are grouped within the specified time window. When disabled, all matching alerts are grouped into a single ongoing episode regardless of time. Defaults to `false`.
- `episode_description_template` (String) Template for generating episode descriptions. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}.
- `episode_labels` (Set of String) Labels to automatically apply to episodes created by this rule. IDs of `oneuptime_label` resources.
- `episode_owner_teams` (Set of String) Teams added as owners of every episode this rule opens, and notified like any owner. Each must be a team of the project. IDs of `oneuptime_team` resources.
- `episode_owner_users` (Set of String) Users added as owners of every episode this rule opens, and notified like any owner. Each must be a member of the project. IDs of `oneuptime_user` records.
- `episode_title_template` (String) Template for generating episode titles. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}.
- `group_by_alert_labels` (Boolean) When enabled, alerts with different sets of labels will be grouped into separate episodes (exact set match). When disabled, alert labels are ignored for grouping. Defaults to `false`.
- `group_by_alert_title` (Boolean) When enabled, alerts with different titles will be grouped into separate episodes. When disabled, alerts with any title can be grouped together. Defaults to `false`.
- `group_by_fields` (String) JSON object defining the fields to group alerts by (e.g., monitorId, severity). A JSON value: write it with `jsonencode()`.
- `group_by_monitor` (Boolean) When enabled, alerts from different monitors will be grouped into separate episodes. When disabled, alerts from any monitor can be grouped together. Defaults to `false`.
- `group_by_monitor_labels` (Boolean) When enabled, alerts whose monitors have different sets of labels will be grouped into separate episodes (exact set match). When disabled, monitor labels are ignored for grouping. Defaults to `false`.
- `group_by_severity` (Boolean) When enabled, alerts with different severities will be grouped into separate episodes. When disabled, alerts of any severity can be grouped together. Defaults to `false`.
- `inactivity_timeout_minutes` (Number) Time in minutes after which an inactive episode will be auto-resolved. Defaults to `60`.
- `is_enabled` (Boolean) Whether this rule is enabled. Defaults to `true`.
- `match_criteria` (String) JSON object defining the criteria for matching alerts to this rule. A JSON value: write it with `jsonencode()`.
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description.
- `monitor_labels` (Set of String) Only group alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels. IDs of `oneuptime_label` resources.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'.
- `monitors` (Set of String) Only group alerts from these monitors. Leave empty to match alerts from any monitor. IDs of `oneuptime_monitor` resources.
- `on_call_duty_policies` (Set of String) List of on-call duty policies to execute for episodes created by this rule. IDs of `oneuptime_on_call_policy` resources.
- `priority` (Number) Where this rule sits in the list. Rules are evaluated from the top of the list down, lowest number first. A new rule is added to the end of the list. Setting a number another one already has puts it in that place, and the ones in the way move one place along to make room. In the dashboard, drag the rows to reorder them.
- `reopen_window_minutes` (Number) Time window in minutes to reopen a recently resolved episode instead of creating a new one. Defaults to `0`.
- `resolve_delay_minutes` (Number) Grace period in minutes before auto-resolving an episode after all alerts are resolved. Defaults to `0`.
- `time_window_minutes` (Number) Rolling time window in minutes. Alerts are grouped if they arrive within this gap from the last alert. Defaults to `60`.

### Read-Only

- `created_at` (String) Date and Time when the object was created.
- `created_by_user_id` (String) User ID who created this object (if this object was created by a User). The ID of a `oneuptime_user` (see the data source).
- `id` (String) Unique identifier for the resource.
- `project_id` (String) ID of your OneUptime Project in which this object belongs. The ID of a `oneuptime_project`.
- `updated_at` (String) Date and Time when the object was updated.

## Import

Import an existing alert grouping rule by its id, with an `import` block (Terraform 1.5+, OpenTofu 1.6+):

```terraform
import {
  to = oneuptime_alert_grouping_rule.example
  id = "<id>"
}
```

or on the command line:

```shell
terraform import oneuptime_alert_grouping_rule.example <id>
```
