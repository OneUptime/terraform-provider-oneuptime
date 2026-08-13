---
page_title: "oneuptime_incident_grouping_rule Data Source - oneuptime"
subcategory: "Incidents"
description: |-
  Configure rules for automatically grouping related incidents into episodes
---

# oneuptime_incident_grouping_rule (Data Source)

Configure rules for automatically grouping related incidents into episodes Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_incident_grouping_rule" "by_name" {
  name = "example-incident_grouping_rule"
}

data "oneuptime_incident_grouping_rule" "by_id" {
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
- `description` (String) Description of this incident grouping rule.. Computed.
- `priority` (Number) Priority of this rule. Lower number = higher priority. Rules are evaluated in priority order... Computed.
- `is_enabled` (Bool) Whether this rule is enabled.. Computed.
- `match_criteria` (String) JSON object defining the criteria for matching incidents to this rule.. Computed.
- `monitors` (Set) Only group incidents from these monitors. Leave empty to match incidents from any monitor... Computed.
- `incident_severities` (Set) Only group incidents with these severities. Leave empty to match incidents of any severity... Computed.
- `incident_labels` (Set) Only group incidents that have at least one of these labels. Leave empty to match incidents regardless of incident labels... Computed.
- `monitor_labels` (Set) Only group incidents from monitors that have at least one of these labels. Leave empty to match incidents regardless of monitor labels... Computed.
- `incident_title_pattern` (String) Regular expression pattern to match incident titles. Leave empty to match any title. Example: 'CPU.*high' matches titles containing 'CPU' followed by 'high'... Computed.
- `incident_description_pattern` (String) Regular expression pattern to match incident descriptions. Leave empty to match any description... Computed.
- `monitor_name_pattern` (String) Regular expression pattern to match monitor names. Leave empty to match any monitor name. Example: 'prod-.*' matches monitors starting with 'prod-'... Computed.
- `monitor_description_pattern` (String) Regular expression pattern to match monitor descriptions. Leave empty to match any monitor description... Computed.
- `group_by_monitor` (Bool) When enabled, incidents from different monitors will be grouped into separate episodes. When disabled, incidents from any monitor can be grouped together... Computed.
- `group_by_severity` (Bool) When enabled, incidents with different severities will be grouped into separate episodes. When disabled, incidents of any severity can be grouped together... Computed.
- `group_by_incident_title` (Bool) When enabled, incidents with different titles will be grouped into separate episodes. When disabled, incidents with any title can be grouped together... Computed.
- `group_by_incident_labels` (Bool) When enabled, incidents with different sets of labels will be grouped into separate episodes (exact set match). When disabled, incident labels are ignored for grouping... Computed.
- `group_by_monitor_labels` (Bool) When enabled, incidents whose monitors have different sets of labels will be grouped into separate episodes (exact set match). When disabled, monitor labels are ignored for grouping... Computed.
- `enable_time_window` (Bool) Enable time-based grouping. When enabled, incidents are grouped within the specified time window. When disabled, all matching incidents are grouped into a single ongoing episode regardless of time... Computed.
- `time_window_minutes` (Number) Rolling time window in minutes. Incidents are grouped if they arrive within this gap from the last incident... Computed.
- `group_by_fields` (String) JSON object defining the fields to group incidents by (e.g., monitorId, severity).. Computed.
- `episode_title_template` (String) Template for generating episode titles. Supports placeholders like {{incidentSeverity}}, {{monitorName}}, {{incidentTitle}}, {{incidentDescription}}.. Computed.
- `episode_description_template` (String) Template for generating episode descriptions. Supports placeholders like {{incidentSeverity}}, {{monitorName}}, {{incidentTitle}}, {{incidentDescription}}.. Computed.
- `enable_resolve_delay` (Bool) Enable grace period before auto-resolving episode after all incidents resolve. Helps prevent rapid state changes during incident flapping... Computed.
- `resolve_delay_minutes` (Number) Grace period in minutes before auto-resolving an episode after all incidents are resolved.. Computed.
- `enable_reopen_window` (Bool) Enable reopening recently resolved episodes instead of creating new ones. Useful when related issues recur shortly after resolution... Computed.
- `reopen_window_minutes` (Number) Time window in minutes to reopen a recently resolved episode instead of creating a new one.. Computed.
- `enable_inactivity_timeout` (Bool) Enable auto-resolving episodes after a period of inactivity. Helps automatically close episodes when no new incidents arrive... Computed.
- `inactivity_timeout_minutes` (Number) Time in minutes after which an inactive episode will be auto-resolved.. Computed.
- `on_call_duty_policies` (Set) List of on-call duty policies to execute for episodes created by this rule... Computed.
- `default_assign_to_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `default_assign_to_team_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `episode_labels` (Set) Labels to automatically apply to episodes created by this rule... Computed.
- `episode_owner_users` (Set) Users to automatically add as owners to episodes created by this rule... Computed.
- `episode_owner_teams` (Set) Teams to automatically add as owners to episodes created by this rule... Computed.
- `episode_member_roles` (Set) Incident roles to display in the episode members form. Select the roles that can be assigned to episode members... Computed.
- `episode_member_role_assignments` (String) Users with specific incident roles to automatically add as members to episodes created by this rule. Each assignment includes a user ID and an incident role ID... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `show_episode_on_status_page` (Bool) Should episodes created by this rule be shown on the status page?.. Computed.
