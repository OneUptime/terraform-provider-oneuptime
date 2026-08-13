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
  name = "Example short text"
  description = "This is an example of longer text content that might be stored in this field."
}
```

## Schema

### Required

- `name` (String) Name of this alert grouping rule..

### Optional

- `project_id` (String) A unique identifier for an object, represented as a UUID..
- `description` (String) Description of this alert grouping rule..
- `priority` (Number) Priority of this rule. Lower number = higher priority. Rules are evaluated in priority order...
- `is_enabled` (Bool) Whether this rule is enabled..
- `match_criteria` (String) JSON object defining the criteria for matching alerts to this rule..
- `monitors` (Set) Only group alerts from these monitors. Leave empty to match alerts from any monitor...
- `alert_severities` (Set) Only group alerts with these severities. Leave empty to match alerts of any severity...
- `alert_labels` (Set) Only group alerts that have at least one of these labels. Leave empty to match alerts regardless of alert labels...
- `monitor_labels` (Set) Only group alerts from monitors that have at least one of these labels. Leave empty to match alerts regardless of monitor labels...
- `alert_title_pattern` (String) Regular expression pattern to match alert titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'...
- `alert_description_pattern` (String) Regular expression pattern to match alert descriptions. Leave empty to match any description...
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'...
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description...
- `group_by_monitor` (Bool) When enabled, alerts from different monitors will be grouped into separate episodes. When disabled, alerts from any monitor can be grouped together...
- `group_by_severity` (Bool) When enabled, alerts with different severities will be grouped into separate episodes. When disabled, alerts of any severity can be grouped together...
- `group_by_alert_title` (Bool) When enabled, alerts with different titles will be grouped into separate episodes. When disabled, alerts with any title can be grouped together...
- `group_by_alert_labels` (Bool) When enabled, alerts with different sets of labels will be grouped into separate episodes (exact set match). When disabled, alert labels are ignored for grouping...
- `group_by_monitor_labels` (Bool) When enabled, alerts whose monitors have different sets of labels will be grouped into separate episodes (exact set match). When disabled, monitor labels are ignored for grouping...
- `enable_time_window` (Bool) Enable time-based grouping. When enabled, alerts are grouped within the specified time window. When disabled, all matching alerts are grouped into a single ongoing episode regardless of time...
- `time_window_minutes` (Number) Rolling time window in minutes. Alerts are grouped if they arrive within this gap from the last alert...
- `group_by_fields` (String) JSON object defining the fields to group alerts by (e.g., monitorId, severity)..
- `episode_title_template` (String) Template for generating episode titles. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}..
- `episode_description_template` (String) Template for generating episode descriptions. Supports placeholders like {{alertSeverity}}, {{monitorName}}, {{alertTitle}}, {{alertDescription}}..
- `enable_resolve_delay` (Bool) Enable grace period before auto-resolving episode after all alerts resolve. Helps prevent rapid state changes during alert flapping...
- `resolve_delay_minutes` (Number) Grace period in minutes before auto-resolving an episode after all alerts are resolved..
- `enable_reopen_window` (Bool) Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution...
- `reopen_window_minutes` (Number) Time window in minutes to reopen a recently resolved episode instead of creating a new one..
- `enable_inactivity_timeout` (Bool) Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new alerts arrive...
- `inactivity_timeout_minutes` (Number) Time in minutes after which an inactive episode will be auto-resolved..
- `on_call_duty_policies` (Set) List of on-call duty policies to execute for episodes created by this rule...
- `default_assign_to_user_id` (String) A unique identifier for an object, represented as a UUID..
- `default_assign_to_team_id` (String) A unique identifier for an object, represented as a UUID..
- `episode_labels` (Set) Labels to automatically apply to episodes created by this rule...
- `episode_owner_users` (Set) Users to automatically add as owners to episodes created by this rule...
- `episode_owner_teams` (Set) Teams to automatically add as owners to episodes created by this rule...
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID..

### Read-Only

- `id` (String) Unique identifier for the resource.
- `created_at` (String) A date time object..
- `updated_at` (String) A date time object..
- `deleted_at` (String) A date time object..
- `version` (Number) Object version.

## Import

Import is supported using the following syntax:

```shell
terraform import oneuptime_alert_grouping_rule.example <id>
```
