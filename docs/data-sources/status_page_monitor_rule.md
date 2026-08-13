---
page_title: "oneuptime_status_page_monitor_rule Data Source - oneuptime"
subcategory: "Status Pages"
description: |-
  Configure rules that automatically add matching monitors to a status page group, instead of picking every monitor by hand
---

# oneuptime_status_page_monitor_rule (Data Source)

Configure rules that automatically add matching monitors to a status page group, instead of picking every monitor by hand Look up by `id` or by `name` (must match exactly one item).

## Example Usage

Look up by `name` (must match exactly one item) or by `id`:

```terraform
data "oneuptime_status_page_monitor_rule" "by_name" {
  name = "example-status_page_monitor_rule"
}

data "oneuptime_status_page_monitor_rule" "by_id" {
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
- `status_page_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `description` (String) Description of this status page monitor rule.. Computed.
- `is_enabled` (Bool) Whether this rule is enabled. A disabled rule removes the monitors it had added... Computed.
- `monitor_labels` (Set) Only match monitors that carry at least one of these labels. Leave empty to skip the label filter... Computed.
- `monitor_name_pattern` (String) Regex (case-insensitive) matched against the monitor name. Leave empty to skip the name filter. Use .* to match every monitor... Computed.
- `monitor_description_pattern` (String) Regex (case-insensitive) matched against the monitor description. Leave empty to skip the description filter... Computed.
- `status_page_group_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
- `show_current_status` (Bool) Show current status like offline, operational or degraded on the resources this rule adds... Computed.
- `show_uptime_percent` (Bool) Show uptime percent on the resources this rule adds to the status page... Computed.
- `uptime_percent_precision` (String) Precision of the uptime percent shown on the resources this rule adds.. Computed.
- `show_status_history_chart` (Bool) Show a 90 day uptime history on the resources this rule adds to the status page... Computed.
- `created_by_user_id` (String) A unique identifier for an object, represented as a UUID.. Computed.
